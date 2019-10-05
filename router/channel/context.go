package channel

import (
	"sync"

	ot "github.com/IntouchOpec/base-go-echo/middleware/opentracing"
	"github.com/IntouchOpec/base-go-echo/middleware/session"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
	"github.com/opentracing/opentracing-go"
)

type Context struct {
	echo.Context
}

func (ctx *Context) Auth() auth.Auth {
	return auth.Default(ctx)
}

type (
	HandlerFunc func(*Context) error
)

func (c *Context) reset() {
	c.Context = nil
}

func (ctx *Context) Session() session.Session {
	return session.Default(ctx)
}

func (ctx *Context) OpenTracingSpan() opentracing.Span {
	return ot.Default(ctx)
}

var (
	ctxPool = sync.Pool{
		New: func() interface{} {
			return &Context{}
		},
	}
)

func NewContext() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			ctx := ctxPool.Get().(*Context)
			defer func() {
				ctx.reset()
				ctxPool.Put(ctx)
			}()

			ctx.Context = c
			return next(ctx)
		}
	}
}
