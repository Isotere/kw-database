package start

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Isotere/logs/logger"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/render"

	"github.com/Isotere/kw-database/apps/server/internal/api/handler/echo"
	"github.com/Isotere/kw-database/pkg/dotenv"
)

const (
	success = 0
	fail    = 1

	shutdownTimeoutSec = 30
)

func Handle(host string, port int) {
	os.Exit(run(host, port))
}

func run(host string, port int) (exitCode int) {
	dotenv.Load()

	log, err := initLogger()
	if err != nil {
		fmt.Println(err.Error())
		return fail
	}
	defer log.Close()

	defer func() {
		if panicErr := recover(); panicErr != nil {
			log.Error(context.Background(), panicErr)
			exitCode = fail
		}
	}()

	ctx, ctxCancel := context.WithCancel(context.Background())
	defer ctxCancel()

	router := initRouter(log)

	router.Get("/_echo", echo.New().Handle)

	// Запуск сервера http с grace-full shutdown
	// ************************************************************************

	serverAddress := fmt.Sprintf("%s:%d", host, port)
	server := &http.Server{Addr: serverAddress, Handler: router, ReadHeaderTimeout: time.Millisecond * 200}

	sig := make(chan os.Signal, 1)
	signal.Notify(sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	go func() {
		<-sig

		// Shutdown signal with grace period of 30 seconds
		shutdownCtx, shutdownCtxCancel := context.WithTimeout(ctx, shutdownTimeoutSec*time.Second)
		defer shutdownCtxCancel()

		log.Info("Starting server shutdown...")
		go func() {
			// Даем серверу какое-то время для завершения всех запущенных процессов в стандартном режиме
			<-shutdownCtx.Done()
			if errors.Is(shutdownCtx.Err(), context.DeadlineExceeded) {
				log.Fatal("graceful shutdown timed out.. forcing exit.")
			}
		}()

		// Trigger graceful shutdown
		errInner := server.Shutdown(shutdownCtx)
		if errInner != nil {
			log.Fatal(errInner)
		}

		ctxCancel()

		log.Info("Server shutdown finished")
	}()

	log.Info("Server started on address: ", serverAddress)

	err = server.ListenAndServe()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		log.Fatal(err)
	}

	<-ctx.Done()

	log.Info("Server stopped")

	// ************************************************************************

	return success
}

func initRouter(log *logger.Logger) *chi.Mux {
	router := chi.NewRouter()

	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.Recoverer)

	router.Use(LoggerMiddleware(log))

	router.Use(render.SetContentType(render.ContentTypeJSON))

	return router
}

func initLogger() (*logger.Logger, error) {
	if dotenv.GetCurrentEnv() != dotenv.Production {
		return logger.New(logger.LogLevelProd)
	}

	return logger.New(logger.LogLevelDevel)
}
