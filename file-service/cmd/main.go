package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/facade"
	"github.com/murilo-bracero/raspstore/file-service/internal/application/usecase"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/config"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/db"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/repository"
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/server"
)

func main() {

	ctx := context.Background()

	if err := godotenv.Load(); err != nil {
		slog.Warn("Could not load .env file. Using system variables instead")
	}

	config := config.NewConfig("config/config.yaml")

	slog.Info("creating required folders")

	if err := os.MkdirAll(config.Storage.Path+"/internal", os.ModePerm); err != nil {
		slog.Error("could not create required internal folder", "error", err)
		os.Exit(1)
	}

	if err := os.MkdirAll(config.Storage.Path+"/storage", os.ModePerm); err != nil {
		slog.Error("could not create required storage folder", "error", err)
		os.Exit(1)
	}

	conn, err := db.NewSqliteDatabaseConnection(config)

	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}

	defer conn.Close()

	fileRepo := repository.NewFilesRepository(ctx, conn.Db())

	txFileRepo := repository.NewTxFilesRepository(ctx, conn.Db())

	useCases := usecase.InitUseCases(config, fileRepo, txFileRepo)

	fileFacade := facade.NewFileFacade(fileRepo)

	if err != nil {
		slog.Error("Error initializing database", "err", err)
	}

	sigc := make(chan os.Signal, 1)

	signal.Notify(sigc, syscall.SIGTERM, syscall.SIGABRT, syscall.SIGQUIT, syscall.SIGINT)

	go func() {
		<-sigc

		if err := conn.Close(); err != nil {
			slog.Error("could not close db connection", "err", err)
			os.Exit(1)
		}

		slog.Info("Process finished")
		os.Exit(0)
	}()

	slog.Info("Bootstraping servers")
	server.StartApiServer(config, fileFacade, useCases)
}
