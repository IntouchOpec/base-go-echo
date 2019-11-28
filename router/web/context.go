package web

import (
	"strconv"
	"sync"

	"github.com/IntouchOpec/base-go-echo/middleware/session"
	"github.com/IntouchOpec/base-go-echo/model"
	"github.com/IntouchOpec/base-go-echo/module/auth"
	"github.com/labstack/echo"
)

type Context struct {
	echo.Context

	Account model.Account `json:"account"`
}

type Queryparams struct {
	Page  string `json:"page"`
	Limit string `json:"limit"`
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

func SetPagination(queryPar map[string][]string) (int, int) {
	var page int
	var limit int = 10
	var err error
	if len(queryPar["limit"]) == 0 {
		limit = 10
	}

	if len(queryPar["page"]) == 0 {
		page = 0
		return 0, 10
	}

	page, err = strconv.Atoi(queryPar["page"][0])
	if err != nil {
		page = 0
	}
	limit, err = strconv.Atoi(queryPar["limit"][0])
	if err != nil {
		limit = 10
		return page, limit
	}

	return page, limit
}

type Pagination struct {
	Page      int
	Previous  bool
	Next      bool
	StartPage int
	List      []int
	Record    int
	Offset    int
}

func MakePagination(total, page, limit int) Pagination {
	previous := false
	next := false
	countPage := total / limit
	startPage := 1
	var list []int

	if page > 1 {
		previous = true
	}

	if page > countPage {
		next = true
	}
	startPage = page + 1
	if startPage <= 2 {
		startPage = 1
	}
	for index := startPage; index < startPage+5; index++ {
		list = append(list, index)

		if countPage+1 <= index {
			break
		}
	}
	record := total - limit*page
	offset := limit * page
	if record > 10 {
		record = 10
	}
	return Pagination{Page: page + 1, Previous: previous, Next: next, StartPage: startPage, List: list,
		Record: record,
		Offset: offset}
}
