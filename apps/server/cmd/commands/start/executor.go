package start

import (
	"context"
	"errors"
	"fmt"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/Isotere/logs/logger"

	"github.com/Isotere/kw-database/apps/server/internal/handler"
	"github.com/Isotere/kw-database/pkg/dotenv"
	"github.com/Isotere/kw-database/pkg/tcp"
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

	// ************************************************************************

	// Запуск сервера tcp с grace-full shutdown
	// ************************************************************************

	server := tcp.NewServer(log, host, port)

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

		server.Shutdown()
		ctxCancel()

		log.Info("Server shutdown finished")
	}()

	h := handler.New(log)

	err = server.Listen(ctx, h.Handle)
	if err != nil {
		log.WithError("tcp server error", err)
		return fail
	}

	<-ctx.Done()

	log.Info("Server stopped")

	// ************************************************************************

	return success
}

func initLogger() (*logger.Logger, error) {
	if dotenv.GetCurrentEnv() != dotenv.Production {
		return logger.New(logger.LogLevelProd)
	}

	return logger.New(logger.LogLevelDevel)
}
