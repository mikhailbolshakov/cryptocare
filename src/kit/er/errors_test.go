package er

import (
	"context"
	"fmt"
	kitContext "github.com/mikhailbolshakov/cryptocare/src/kit/context"
	"github.com/stretchr/testify/assert"
	"testing"
)

// by default prints in format "code: message"
func Test_Error(t *testing.T) {
	e := New("ERR-123", "%s happened", "shit")
	fmt.Println(e)
}

// use Is function to assert to *AppError
func Test_WithStack(t *testing.T) {
	e := New("ERR-123", "%s happened", "shit")
	if appErr, ok := Is(e); ok {
		fmt.Println(appErr.WithStack())
		return
	}
	t.Fatal()
}

// you can get error interface with stack
func Test_WithStackError(t *testing.T) {
	e := New("ERR-123", "%s happened", "shit")
	if appErr, ok := Is(e); ok {
		fmt.Println(appErr.WithStackErr())
		return
	}
	t.Fatal()
}

// use Is function to split code and message
func Test_CodeAndMessageSplit(t *testing.T) {
	e := New("ERR-123", "%s happened", "shit")
	if appErr, ok := Is(e); ok {
		fmt.Printf("%s %s", appErr.Code(), appErr.Message())
		return
	}
	t.Fatal()
}

func Test_Wrap(t *testing.T) {
	originalErr := fmt.Errorf("original issue")
	e := Wrap(originalErr, "ERR-123", "%s happened", "shit")
	if appErr, ok := Is(e); ok {
		fmt.Println(appErr.WithStackErr())
		return
	}
	t.Fatal()
}

func Test_TwoWrappers(t *testing.T) {
	originalErr := fmt.Errorf("original issue")
	e := Wrap(originalErr, "ERR-123", "%s happened", "shit")
	e2 := Wrap(e, "ERR-124", "very bad %s happened", "shit")
	if appErr, ok := Is(e2); ok {
		fmt.Println(appErr.WithStackErr())
		return
	}
	t.Fatal()
}

func Test_NewWithBuilder_WhenContext(t *testing.T) {
	e := WithBuilder("ERR-123", "%s happens", "shit").
		C(kitContext.NewRequestCtx().Test().ToContext(context.Background())).
		Err()
	if appErr, ok := Is(e); ok {
		ctx := appErr.Fields()["ctx"]
		assert.NotEmpty(t, ctx)
		cl := ctx.(map[string]interface{})["_ctx.cl"]
		assert.NotEmpty(t, "test", cl)
		return
	}
	t.Fatal()
}

func Test_NewWithBuilder_WhenEmptyContext(t *testing.T) {
	e := WithBuilder("ERR-123", "%s happens", "shit").C(context.Background()).Err()
	if appErr, ok := Is(e); ok {
		fmt.Println(appErr)
		return
	}
	t.Fatal()
}

func Test_NewWithBuilder_WhenFields(t *testing.T) {
	e := WithBuilder("ERR-123", "%s happens", "shit").
		F(FF{"f": "v"}).
		Err()
	if appErr, ok := Is(e); ok {
		fmt.Println(appErr.WithStack())
		assert.NotEmpty(t, appErr.fields)
		assert.Equal(t, appErr.fields["f"], "v")
		return
	}
	t.Fatal()
}

func Test_NewWithBuilder_WrapWithFields(t *testing.T) {
	originalErr := WithBuilder("ERR-123", "%s happened", "shit").
		F(FF{"f": "v", "f2": "v2"}).
		Err()

	e := Wrap(originalErr, "ERR-124", "%s happens", "shit2")
	if appErr, ok := Is(e); ok {
		fmt.Println(appErr.WithStack())
		assert.True(t, ok)
		assert.NotEmpty(t, appErr.fields)
		assert.Equal(t, appErr.fields["f"], "v")
		assert.Equal(t, appErr.fields["f2"], "v2")
		return
	}
	t.Fatal()
}
