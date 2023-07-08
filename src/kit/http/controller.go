package http

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	kitContext "github.com/mikhailbolshakov/cryptocare/src/kit/context"
	"github.com/mikhailbolshakov/cryptocare/src/kit/er"
	"github.com/mikhailbolshakov/cryptocare/src/kit/log"
	"io"
	"mime"
	"mime/multipart"
	"net/http"
	"strconv"
	"strings"
	"time"
)

const (
	SortRequestMissingFirst = "first"
	SortRequestMissingLast  = "last"
)

type SortRequest struct {
	Field   string `json:"field"`
	Asc     bool   `json:"asc"`
	Missing string `json:"missing"`
}

type PagingRequest struct {
	Size   int            `json:"size"`
	Index  int            `json:"index"`
	SortBy []*SortRequest `json:"sortBy"`
}

type PagingResponse struct {
	Total int `json:"total"`
	Index int `json:"index"`
}

// Error is a HTTP error object returning to clients in case of error
type Error struct {
	Code           string                 `json:"code,omitempty"`    // Code is error code provided by error producer
	Type           string                 `json:"type,omitempty"`    // Type is error type (panic, system, business)
	Message        string                 `json:"message"`           // Message is error description
	TranslationKey string                 `json:"translationKey"`    // TranslationKey is error code translation key
	Details        map[string]interface{} `json:"details,omitempty"` // Details is additional info provided by error producer
}

func (e *Error) Error() string {
	return fmt.Sprintf("%s:%s", e.Code, e.Message)
}

const (
	Me = "me" // Me can be used in URL whenever userId is expected. When encountered, userId from the session context is used
)

var EmptyOkResponse = struct {
	Status string `json:"status"`
}{
	Status: "OK",
}

// Controller is a base controller interface
type Controller interface {
	// MyClientProfile returns true if client requests his own client's data
	MyClientProfile(ctx context.Context, r *http.Request) (bool, error)
	// MyUser returns true if current user requests his own data
	MyUser(ctx context.Context, r *http.Request) (bool, error)
}

// BaseController is a base controller implementation
type BaseController struct {
	Logger log.CLoggerFunc
}

var MediaContentTypes = [...]string{
	"image/jpeg",
	"image/png",
	"image/bmp",
	"image/gif",
	"image/tiff",
	"video/avi",
	"video/mpeg",
	"video/mp4",
	"audio/mpeg",
	"audio/wav",
}

type ResponseContentOpts struct {
	Filename     string
	ContentType  string
	ContentSize  int
	Download     bool
	ModifiedTime time.Time
}

func (c *BaseController) RespondContent(w http.ResponseWriter, r *http.Request, opts ResponseContentOpts, file []byte) {

	w.Header().Set("Cache-Control", "private, no-cache")

	if opts.ContentSize > 0 {
		contentSizeStr := strconv.Itoa(opts.ContentSize)
		w.Header().Set("Content-Length", contentSizeStr)
	}

	if opts.ContentType == "" {
		opts.ContentType = "application/octet-stream"
	}
	w.Header().Set("Content-Type", opts.ContentType)

	if !opts.Download {
		isMedia := false
		for _, mct := range MediaContentTypes {
			if strings.HasPrefix(opts.ContentType, mct) {
				isMedia = true
				break
			}
		}
		opts.Download = !isMedia
	}

	if opts.Download {
		w.Header().Set("Content-Disposition", "attachment;filename=\""+opts.Filename+"\"; filename*=UTF-8''"+opts.Filename)
	} else {
		w.Header().Set("Content-Disposition", "inline;filename=\""+opts.Filename+"\"; filename*=UTF-8''"+opts.Filename)
	}

	http.ServeContent(w, r, opts.Filename, opts.ModifiedTime, bytes.NewReader(file))

}

// GetUploadFileMultipartContent it parse body for multipart content disposition
// it expects the only one part with the following structure:
//-----------------------------4562559108110960722260982980
//Content-Disposition: form-data; name="files"; filename="my-file.jpg"
//Content-Type: image/jpeg
//....
//.....
func (c *BaseController) GetUploadFileMultipartContent(ctx context.Context, r *http.Request) (io.Reader, string, error) {

	// parse form
	if r.Form == nil {
		err := r.ParseForm()
		if err != nil {
			return nil, "", ErrHttpMultipartParseForm(err, ctx)
		}
	}
	if r.ContentLength == 0 {
		return nil, "", ErrHttpMultipartEmptyContent(ctx)
	}

	// get content type from header
	contentType := r.Header.Get("Content-Type")
	if contentType == "" {
		return nil, "", ErrHttpMultipartNotMultipart(ctx)
	}

	// parse mime type
	mediaType, params, err := mime.ParseMediaType(contentType)
	if err != nil {
		return nil, "", ErrHttpMultipartParseMediaType(err, ctx)
	}
	if mediaType != "multipart/form-data" {
		return nil, "", ErrHttpMultipartWrongMediaType(ctx, mediaType)
	}

	// identify boundary
	boundary, ok := params["boundary"]
	if !ok {
		return nil, "", ErrHttpMultipartMissingBoundary(ctx)
	}

	// create a new reader
	mr := multipart.NewReader(r.Body, boundary)

	// go through all parts
	for {

		// take next part
		part, err := mr.NextPart()
		if err != nil {
			if err == io.EOF {
				// if we get here, we haven't found any useful parts, so it's wrong format
				return nil, "", ErrHttpMultipartEofReached(ctx)
			} else {
				return nil, "", ErrHttpMultipartNext(err, ctx)
			}
		}

		// check found part
		if part.FormName() == "file" {
			filename := part.FileName()
			if filename == "" {
				return nil, "", ErrHttpMultipartFilename(ctx)
			}
			// return first part
			return part, filename, nil
		} else {
			return nil, "", ErrHttpMultipartFormNameFileExpected(ctx)
		}

	}
}

func (c *BaseController) RespondJson(w http.ResponseWriter, httpStatus int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatus)
	_, _ = w.Write(response)
}

func (c *BaseController) RespondError(w http.ResponseWriter, err error) {

	httpErr := &Error{}
	httpStatus := http.StatusInternalServerError

	// check if this is an app error
	if appErr, ok := er.Is(err); ok {
		httpErr.Code = appErr.Code()
		httpErr.Message = appErr.Message()
		httpErr.TranslationKey = "errors.app.code." + strings.ReplaceAll(strings.ToLower(appErr.Code()), "-", ".")
		httpErr.Details = appErr.Fields()
		httpErr.Type = appErr.Type()
		if httpSt := appErr.HttpStatus(); httpSt != nil {
			httpStatus = int(*httpSt)
		}
	} else {
		httpErr.Message = err.Error()
	}
	if c.Logger != nil {
		c.Logger().Cmp("api").Pr("rest").E(err).St().Err()
	}
	c.RespondJson(w, httpStatus, httpErr)
}

func (c *BaseController) RespondWithStatus(w http.ResponseWriter, status int, payload interface{}) {
	c.RespondJson(w, status, payload)
}

func (c *BaseController) RespondOK(w http.ResponseWriter, payload interface{}) {
	c.RespondJson(w, http.StatusOK, payload)
}

func (c *BaseController) DecodeRequest(r *http.Request, ctx context.Context, body interface{}) error {
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(body); err != nil {
		return ErrHttpDecodeRequest(err, ctx)
	}
	return nil
}

func (c *BaseController) Var(r *http.Request, ctx context.Context, varName string, allowEmpty bool) (string, error) {
	if val, ok := mux.Vars(r)[varName]; ok {
		if !allowEmpty && val == "" {
			return "", ErrHttpUrlVarEmpty(ctx, varName)
		}
		return val, nil
	} else {
		return "", ErrHttpUrlVar(ctx, varName)
	}
}

func (c *BaseController) VarUUID(r *http.Request, ctx context.Context, varName string, allowEmpty bool) (string, error) {
	valStr, err := c.Var(r, ctx, varName, allowEmpty)
	if err != nil {
		return "", err
	}
	if allowEmpty && valStr == "" {
		return "", nil
	}
	err = kit.ValidateUUIDs(valStr)
	if err != nil {
		return "", ErrHttpUrlVarInvalidUUID(ctx, varName)
	}
	return valStr, nil
}

func (c *BaseController) CurrentUser(ctx context.Context) (uid string, un string, err error) {
	if rCtx, ok := kitContext.Request(ctx); ok {
		if rCtx.Un != "" && rCtx.Uid != "" {
			return rCtx.Uid, rCtx.Un, nil
		} else {
			return "", "", ErrHttpCurrentUser(ctx)
		}
	} else {
		return "", "", ErrHttpCurrentUser(ctx)
	}
}

func (c *BaseController) UserIdVar(r *http.Request, ctx context.Context, varName string) (string, error) {
	val, err := c.Var(r, ctx, varName, false)
	if err != nil {
		return "", err
	}
	// if current user
	if val == Me {
		if uid, _, err := c.CurrentUser(ctx); err != nil {
			return "", err
		} else {
			return uid, nil
		}
	}
	// validate UUID
	err = kit.ValidateUUIDs(val)
	if err != nil {
		return "", ErrHttpUrlVarInvalidUUID(ctx, val)
	}
	return val, nil
}

func (c *BaseController) UserNameVar(r *http.Request, ctx context.Context, varName string) (string, error) {
	val, err := c.Var(r, ctx, varName, false)
	if err != nil {
		return "", err
	}
	if val == Me {
		if _, un, err := c.CurrentUser(ctx); err != nil {
			return "", err
		} else {
			return un, nil
		}
	}
	return val, nil
}

func (c *BaseController) FormVal(r *http.Request, ctx context.Context, name string, allowEmpty bool) (string, error) {
	val := r.FormValue(name)
	if !allowEmpty && val == "" {
		return "", ErrHttpUrlFormVarEmpty(ctx, name)
	}
	return val, nil
}

func (c *BaseController) FormValUUID(r *http.Request, ctx context.Context, name string, allowEmpty bool) (string, error) {
	valStr, err := c.FormVal(r, ctx, name, allowEmpty)
	if err != nil {
		return "", err
	}
	if allowEmpty && valStr == "" {
		return "", nil
	}
	err = kit.ValidateUUIDs(valStr)
	if err != nil {
		return "", ErrHttpUrlVarInvalidUUID(ctx, name)
	}
	return valStr, nil
}

func (c *BaseController) FormValInt(r *http.Request, ctx context.Context, name string, allowEmpty bool) (*int, error) {
	valStr, err := c.FormVal(r, ctx, name, allowEmpty)
	if err != nil {
		return nil, err
	}
	if allowEmpty && valStr == "" {
		return nil, nil
	}
	valInt, err := strconv.Atoi(valStr)
	if err != nil {
		return nil, ErrHttpUrlFormVarNotInt(err, ctx, name)
	}
	return &valInt, nil
}

func (c *BaseController) FormValFloat(r *http.Request, ctx context.Context, name string, allowEmpty bool) (*float64, error) {
	valStr, err := c.FormVal(r, ctx, name, allowEmpty)
	if err != nil {
		return nil, err
	}
	if allowEmpty && valStr == "" {
		return nil, nil
	}
	valFloat, err := strconv.ParseFloat(valStr, 64)
	if err != nil {
		return nil, ErrHttpUrlFormVarNotFloat(err, ctx, name)
	}
	return &valFloat, nil
}

func (c *BaseController) FormValBool(r *http.Request, ctx context.Context, name string, allowEmpty bool) (*bool, error) {
	valStr, err := c.FormVal(r, ctx, name, allowEmpty)
	if err != nil {
		return nil, err
	}
	if allowEmpty && valStr == "" {
		return nil, nil
	}
	b, err := strconv.ParseBool(valStr)
	if err != nil {
		return nil, ErrHttpUrlFormVarNotBool(err, ctx, name)
	}
	return &b, nil
}

// FormValTime parses URL form value and checks for time in RFC3339 format(UTC)
func (c *BaseController) FormValTime(r *http.Request, ctx context.Context, name string, allowEmpty bool) (*time.Time, error) {
	valStr, err := c.FormVal(r, ctx, name, allowEmpty)
	if err != nil {
		return nil, err
	}
	if allowEmpty && valStr == "" {
		return nil, nil
	}
	valTime, err := time.Parse(time.RFC3339, valStr)
	if err != nil {
		return nil, ErrHttpUrlFormVarNotTime(err, ctx, name)
	}
	return &valTime, nil
}

// FormSort parses URL form value with sort parameters and make a slice of special objects
func (c *BaseController) FormSort(r *http.Request, ctx context.Context, name string, allowEmpty bool) ([]*SortRequest, error) {
	valStr, err := c.FormVal(r, ctx, name, allowEmpty)
	if err != nil {
		return nil, err
	}
	if allowEmpty && valStr == "" {
		return nil, nil
	}
	return ParseSortBy(ctx, valStr)
}

// FormPaging parses URL form value for paging params. Allows specifying max page size
func (c *BaseController) FormPaging(r *http.Request, ctx context.Context, maxPageSize *int) (size *int, index *int, err error) {
	size, err = c.FormValInt(r, ctx, "size", true)
	if err != nil {
		return
	}
	index, err = c.FormValInt(r, ctx, "index", true)
	if err != nil {
		return
	}
	if maxPageSize != nil && size != nil && *size > *maxPageSize {
		err = ErrHttpUrlMaxPageSizeExceeded(ctx, *maxPageSize)
		return
	}
	return
}

// MyUser returns true if current user requests his own data
func (c *BaseController) MyUser(ctx context.Context, r *http.Request) (bool, error) {
	currentUid, _, err := c.CurrentUser(ctx)
	if err != nil {
		return false, err
	}
	uid, err := c.UserIdVar(r, ctx, "userId")
	return currentUid == uid && err == nil, nil
}
