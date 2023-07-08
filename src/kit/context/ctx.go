package context

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/mikhailbolshakov/cryptocare/src/kit"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/grpc/metadata"
)

const (
	CallerTypeRest   = "rest"
	CallerTypeTest   = "test"
	CallerTypeJob    = "job"
	CallerTypeQueue  = "queue"
	CallerTypeWs     = "ws"
	CallerTypeWebRtc = "webrtc"
)

type requestContextKey struct{}

type RequestContext struct {
	// Rid request ID
	Rid string `json:"_ctx.rid,omitempty" mapstructure:"_ctx.rid"`
	// Sid session ID
	Sid string `json:"_ctx.sid,omitempty" mapstructure:"_ctx.sid"`
	// Uid user ID
	Uid string `json:"_ctx.uid,omitempty" mapstructure:"_ctx.uid"`
	// Un username
	Un string `json:"_ctx.un,omitempty" mapstructure:"_ctx.un"`
	// Caller who is calling
	Caller string `json:"_ctx.cl,omitempty" mapstructure:"_ctx.cl"`
	// Roles list of roles
	Roles []string `json:"_ctx.rl,omitempty" mapstructure:"_ctx.rl"`
}

func NewRequestCtx() *RequestContext {
	return &RequestContext{}
}

func (r *RequestContext) GetRequestId() string {
	return r.Rid
}

func (r *RequestContext) GetSessionId() string {
	return r.Sid
}

func (r *RequestContext) GetUserId() string {
	return r.Uid
}

func (r *RequestContext) GetCaller() string {
	return r.Caller
}

func (r *RequestContext) GetRoles() []string {
	return r.Roles
}

func (r *RequestContext) GetUsername() string {
	return r.Un
}

func (r *RequestContext) Empty() *RequestContext {
	return &RequestContext{}
}

func (r *RequestContext) WithRequestId(requestId string) *RequestContext {
	r.Rid = requestId
	return r
}

func (r *RequestContext) WithNewRequestId() *RequestContext {
	r.Rid = kit.NewId()
	return r
}

func (r *RequestContext) WithSessionId(sessionId string) *RequestContext {
	r.Sid = sessionId
	return r
}

func (r *RequestContext) Rest() *RequestContext {
	r.Caller = CallerTypeRest
	return r
}

func (r *RequestContext) Webrtc() *RequestContext {
	r.Caller = CallerTypeWebRtc
	return r
}

func (r *RequestContext) Test() *RequestContext {
	r.Caller = CallerTypeTest
	return r
}

func (r *RequestContext) Job() *RequestContext {
	r.Caller = CallerTypeJob
	return r
}

func (r *RequestContext) Queue() *RequestContext {
	r.Caller = CallerTypeQueue
	return r
}

func (r *RequestContext) Ws() *RequestContext {
	r.Caller = CallerTypeWs
	return r
}

func (r *RequestContext) WithCaller(caller string) *RequestContext {
	r.Caller = caller
	return r
}

func (r *RequestContext) WithUser(userId, username string) *RequestContext {
	r.Uid = userId
	r.Un = username
	return r
}

func (r *RequestContext) WithRoles(roles ...string) *RequestContext {
	r.Roles = roles
	return r
}

func (r *RequestContext) ToContext(parent context.Context) context.Context {
	if parent == nil {
		parent = context.Background()
	}
	return context.WithValue(parent, requestContextKey{}, r)
}

func (r *RequestContext) ToMap() map[string]interface{} {
	return map[string]interface{}{
		"_ctx.rid": r.Rid,
		"_ctx.sid": r.Sid,
		"_ctx.uid": r.Uid,
		"_ctx.un":  r.Un,
		"_ctx.cl":  r.Caller,
		"_ctx.rl":  r.Roles,
	}
}

func Request(context context.Context) (*RequestContext, bool) {
	if r, ok := context.Value(requestContextKey{}).(*RequestContext); ok {
		return r, true
	}
	return &RequestContext{}, false
}

func MustRequest(context context.Context) (*RequestContext, error) {
	if r, ok := context.Value(requestContextKey{}).(*RequestContext); ok {
		return r, nil
	}
	return &RequestContext{}, errors.New("context is invalid")
}

func FromContextToGrpcMD(ctx context.Context) (metadata.MD, bool) {
	if r, ok := Request(ctx); ok {
		rm, _ := json.Marshal(*r)
		return metadata.Pairs("rq-bin", string(rm)), true
	}
	return metadata.Pairs(), false
}

func FromGrpcMD(ctx context.Context, md metadata.MD) context.Context {

	if rqb, ok := md["rq-bin"]; ok {
		if len(rqb) > 0 {
			rm := []byte(rqb[0])
			rq := &RequestContext{}
			_ = json.Unmarshal(rm, rq)
			return context.WithValue(ctx, requestContextKey{}, rq)
		}
	}
	return ctx
}

func FromMap(ctx context.Context, mp map[string]interface{}) (context.Context, error) {
	var r *RequestContext
	err := mapstructure.Decode(mp, &r)
	if err != nil {
		return nil, err
	}
	return r.ToContext(ctx), nil
}

func Copy(ctx context.Context) context.Context {
	if r, ok := Request(ctx); ok {
		ct, err := FromMap(context.TODO(), r.ToMap())
		if err != nil {
			return ctx
		}

		return ct
	}

	return ctx
}
