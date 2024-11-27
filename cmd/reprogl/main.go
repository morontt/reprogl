package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"runtime/debug"
	"syscall"
	"time"

	"github.com/xelbot/reverse"
	"xelbot.com/reprogl/container"
	"xelbot.com/reprogl/middlewares"
	"xelbot.com/reprogl/services/oauth"
	"xelbot.com/reprogl/views"
)

func main() {
	var err error
	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime|log.Lmicroseconds)
	errorLog := log.New(os.Stderr, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)

	infoLog.Printf("Version: %s (tag: %s)", container.Version, container.GetBuildTag())
	infoLog.Printf("Build time: %s", container.BuildTimeRFC1123())
	infoLog.Printf("Go version: %s", runtime.Version())

	handleError(container.Load("config/app.ini"), errorLog)

	app := &container.Application{
		ErrorLog: errorLog,
		InfoLog:  infoLog,
	}

	handleError(app.SetupDatabase(), errorLog)
	app.SetupMetrics()

	router := getRoutes(app)
	handler := middlewares.Session(router, infoLog)
	handler = middlewares.Clacks().Middleware(handler)
	handler = middlewares.CDN(handler)
	handler = middlewares.Recover(handler, app)
	handler = middlewares.Track(handler, app)
	handler = middlewares.AccessLog(handler, app)
	handler = middlewares.ResponseWrapper(handler)

	cfg := container.GetConfig()

	urlGenerator := func(routeName string, absoluteURL bool, pairs ...string) string {
		var url string
		url, err = reverse.Get(routeName, pairs...)
		if err != nil {
			errorLog.Printf("[urlGenerator] URL generation error for: %s", routeName)
			url = "/error"
		}

		var prefix string
		if absoluteURL {
			prefix = "https://" + cfg.Host
		}

		return prefix + url
	}

	container.SetURLGenerator(urlGenerator)
	handleError(views.LoadViewSet(), errorLog)

	oauth.SetLogger(infoLog)

	mux := http.NewServeMux()
	mux.Handle("/metrics", getPrometheusMetricsHandler(app.Metrics.Collectors()...))
	mux.Handle("/", handler)

	server := &http.Server{
		Handler:      mux,
		Addr:         ":" + cfg.Port,
		WriteTimeout: 10 * time.Second,
		ReadTimeout:  10 * time.Second,
	}

	go func() {
		infoLog.Printf("Starting server on %s port", cfg.Port)
		err = server.ListenAndServe()
		if !errors.Is(err, http.ErrServerClosed) {
			handleError(err, errorLog)
		}
	}()

	exit := make(chan os.Signal, 1)
	signal.Notify(exit, os.Interrupt, syscall.SIGTERM)

	<-exit

	infoLog.Print("Shutting down...")
	err = server.Shutdown(context.Background())
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Server stopped")

	err = app.Stop()
	if err != nil {
		errorLog.Fatal(err)
	}
	infoLog.Print("Application stopped")
}

func handleError(err error, logger *log.Logger) {
	if err != nil {
		logger.Println(fmt.Sprintf("%s", debug.Stack()))
		logger.Fatal(err)
	}
}
