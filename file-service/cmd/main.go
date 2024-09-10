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
	"github.com/murilo-bracero/raspstore/file-service/internal/infra/bootstrap"
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

	config := config.New("config/config.yaml")

	slog.Info("Bootstrapping Application")

	err := bootstrap.Run(ctx, config)

	if err != nil {
		slog.Error("Could not bootstrap the application", "err", err)
	}

	conn, err := db.NewSqliteDatabaseConnection(config)

	if err != nil {
		slog.Error("could not connect to database", "error", err)
		os.Exit(1)
	}

	defer conn.Close()

	fileRepo := repository.NewFilesRepository(ctx, conn.Db())

	txFileRepo := repository.NewTxFilesRepository(ctx, conn.Db())

	loginPAMUseCase := usecase.NewLoginPAMUseCase(config)

	createFileUseCase := usecase.NewCreateFileUseCase(config, fileRepo)

	updateFileUseCase := usecase.NewUpdateFileUseCase(txFileRepo)

	fileFacade := facade.NewFileFacade(fileRepo)

	fileSystemFacade := facade.NewFileSystemFacade(config)

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

	server.StartApiServer(&server.ApiServerParams{
		Config:            config,
		CreateFileUseCase: createFileUseCase,
		FileFacade:        fileFacade,
		FileSystemFacade:  fileSystemFacade,
		LoginPAMUseCase:   loginPAMUseCase,
		UpdateFileUseCase: updateFileUseCase,
	})
}
