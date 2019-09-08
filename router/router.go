package router

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"os/signal"
	"time"

	"github.com/labstack/echo"
	mw "github.com/labstack/echo/middleware"
	"go.elastic.co/apm"
	"go.elastic.co/apm/module/apmechov4"

	. "github.com/IntouchOpec/base-go-echo/conf"
	"github.com/IntouchOpec/base-go-echo/router/api"
	"github.com/IntouchOpec/base-go-echo/router/channel"
	"github.com/hb-go/echo-web/middleware/metrics/prometheus"
	"github.com/hb-go/echo-web/middleware/opentracing"
	"github.com/hb-go/echo-web/middleware/pprof"
	"github.com/hb-go/echo-web/module/log"
)

type (
	Host struct {
		Echo *echo.Echo
	}
)

func InitRoutes() map[string]*Host {
	// Hosts
	hosts := make(map[string]*Host)

	// hosts[Conf.Server.DomainWeb] = &Host{web.Routers()}
	hosts["localhost"] = &Host{api.Routers()}
	hosts["942b3ab4.ngrok.io"] = &Host{channel.Routers()}

	return hosts
}

func RunSubdomains(confFilePath string) {

	log.SetLevel(GetLogLvl())

	// Server
	e := echo.New()

	// pprof
	e.Pre(pprof.Serve())

	e.Pre(mw.RemoveTrailingSlash())

	// Elastic APM
	// Requires APM Server 6.5.0 or newer
	apm.DefaultTracer.Service.Name = Conf.Opentracing.ServiceName
	apm.DefaultTracer.Service.Version = Conf.App.Version
	e.Use(apmechov4.Middleware(
		apmechov4.WithRequestIgnorer(func(request *http.Request) bool {
			return false
		}),
	))

	// OpenTracing
	otCtf := opentracing.Configuration{
		Disabled: Conf.Opentracing.Disable,
		Type:     opentracing.TracerType(Conf.Opentracing.Type),
	}
	if closer := otCtf.InitGlobalTracer(
		opentracing.ServiceName(Conf.Opentracing.ServiceName),
		opentracing.Address(Conf.Opentracing.Address),
	); closer != nil {
		defer closer.Close()
	}

	e.Logger.SetLevel(GetLogLvl())

	// Metrics
	if !Conf.Metrics.Disable {
		e.Use(prometheus.MetricsFunc(
			prometheus.Namespace("echo_web"),
		))
	}

	// Secure, XSS/CSS HSTS
	e.Use(mw.SecureWithConfig(mw.DefaultSecureConfig))
	e.Use(mw.MethodOverride())

	// CORS
	e.Use(mw.CORSWithConfig(mw.CORSConfig{
		// AllowOrigins: []string{"http://" + Conf.Server.DomainWeb, "http://" + Conf.Server.DomainApi},
		AllowHeaders: []string{echo.HeaderOrigin, echo.HeaderContentType, echo.HeaderAcceptEncoding, echo.HeaderAuthorization},
	}))
	hosts := InitRoutes()
	e.Any("/*", func(c echo.Context) (err error) {
		req := c.Request()
		res := c.Response()
		fmt.Println(req.Host, "========req.Host==========")

		u, _err := url.Parse(c.Scheme() + "://" + req.Host)
		if _err != nil {
			e.Logger.Errorf("Request URL parse error:%v", _err)
		}

		host := hosts[u.Hostname()]
		if host == nil {
			e.Logger.Info("Host not found")
			err = echo.ErrNotFound
		} else {
			host.Echo.ServeHTTP(res, req)
		}

		return
	})

	if !Conf.Server.Graceful {
		e.Logger.Fatal(e.Start(Conf.Server.Addr))
	} else {
		// Graceful Shutdown
		// Start server
		go func() {
			if err := e.Start(Conf.Server.Addr); err != nil {
				e.Logger.Errorf("Shutting down the server with error:%v", err)
			}
		}()

		// Wait for interrupt signal to gracefully shutdown the server with
		// a timeout of 10 seconds.
		quit := make(chan os.Signal)
		signal.Notify(quit, os.Interrupt)
		<-quit
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		if err := e.Shutdown(ctx); err != nil {
			e.Logger.Fatal(err)
		}
	}
}
