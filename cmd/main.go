package main

import (
	"context"
	"errors"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"time"

	"github.com/caarlos0/env"
	"github.com/rendau/s3-uploader/internal/core"
	"github.com/rendau/s3-uploader/internal/httphandler"
	"github.com/rendau/s3-uploader/internal/ost"
	"github.com/rs/cors"
)

type Conf struct {
	Debug       bool   `env:"DEBUG" envDefault:"false"`
	HttpPort    string `env:"HTTP_PORT" envDefault:"80"`
	HttpCors    bool   `env:"HTTP_CORS" envDefault:"false"`
	OstUrl      string `env:"OST_URL"`
	OstKeyId    string `env:"OST_KEY_ID"`
	OstKey      string `env:"OST_KEY"`
	OstSecure   bool   `env:"OST_SECURE" envDefault:"false"`
	OstBucket   string `env:"OST_BUCKET"`
	UrlTemplate string `env:"URL_TEMPLATE"`
}

type App struct {
	conf Conf
	ost  *ost.St
	core *core.Core

	httpServer *http.Server

	exitCode int
}

func (a *App) Init() {
	// config
	{
		if err := env.Parse(&a.conf); err != nil {
			panic(err)
		}
	}

	// logger
	{
		if !a.conf.Debug {
			logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
			slog.SetDefault(logger)
		}
	}

	// ost
	{
		var err error
		a.ost, err = ost.New(a.conf.OstUrl, a.conf.OstKeyId, a.conf.OstKey, a.conf.OstSecure)
		a.errAssert(err, "fail to ost.New")
	}

	// core
	{
		var err error
		a.core, err = core.NewCore(a.conf.OstBucket, a.conf.UrlTemplate, a.ost)
		a.errAssert(err, "fail to core.NewCore")
	}

	// http server
	{
		router := httphandler.NewHttpHandler(a.core)

		// router

		// add cors middleware
		if a.conf.HttpCors {
			router = cors.New(cors.Options{
				AllowOriginFunc: func(origin string) bool { return true },
				AllowedMethods: []string{
					http.MethodGet,
					http.MethodPut,
					http.MethodPatch,
					http.MethodPost,
					http.MethodDelete,
				},
				AllowedHeaders: []string{
					"Accept",
					"Content-Type",
					"X-Requested-With",
					"Authorization",
				},
				AllowCredentials: true,
				MaxAge:           604800,
			}).Handler(router)
		}

		// add recover middleware
		router = func(h http.Handler) http.Handler {
			return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				defer func() {
					if recErr := recover(); recErr != nil {
						buf := make([]byte, 1<<20)
						stackLen := runtime.Stack(buf, false)
						slog.Error("Recovered from panic", slog.Any("err", recErr), slog.String("stack", string(buf[:stackLen])))
						w.WriteHeader(http.StatusInternalServerError)
					}
				}()
				h.ServeHTTP(w, r)
			})
		}(router)

		// server
		a.httpServer = &http.Server{
			Addr:              ":" + a.conf.HttpPort,
			Handler:           router,
			ReadHeaderTimeout: 2 * time.Second,
			ReadTimeout:       time.Minute,
			MaxHeaderBytes:    300 * 1024,
		}
	}
}

func (a *App) Start() {
	slog.Info("start")

	// http server
	{
		go func() {
			err := a.httpServer.ListenAndServe()
			if err != nil && !errors.Is(err, http.ErrServerClosed) {
				a.errAssert(err, "grpc-server stopped")
			}
		}()
		slog.Info("http-server started " + a.httpServer.Addr)
	}

	signalCtx, signalCtxCancel := signal.NotifyContext(context.Background(), os.Interrupt)
	defer signalCtxCancel()

	// wait signal
	<-signalCtx.Done()
}

func (a *App) Stop() {
	slog.Info("stopping")
}

func (a *App) WaitJobs() {
	slog.Info("waiting jobs")
}

func (a *App) errAssert(err error, msg string) {
	if err != nil {
		if msg != "" {
			err = fmt.Errorf("%s: %w", msg, err)
		}
		slog.Error(err.Error())
		os.Exit(1)
	}
}

func (a *App) Exit() {
	slog.Info("exit")
	os.Exit(a.exitCode)
}

func main() {
	app := App{}
	app.Init()
	app.Start()
	app.Stop()
	app.WaitJobs()
	app.Exit()
}
