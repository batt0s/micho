package logging

import (
	"io"
	"log"
	"log/slog"
	"os"
	"path/filepath"
)

var AuditLog *slog.Logger

func Init() error {
	logDir := "logs"
	if err := os.MkdirAll(logDir, 0755); err != nil {
		return err
	}

	logFile, err := os.OpenFile(filepath.Join(logDir, "audit.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	handler := slog.NewJSONHandler(logFile, &slog.HandlerOptions{
		Level:     slog.LevelInfo,
		AddSource: true,
	})

	AuditLog = slog.New(handler)

	AuditLog.Info("Audit Logger initialized successfully")

	appLogFile, err := os.OpenFile(filepath.Join(logDir, "micho.log"), os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		return err
	}

	multiWriter := io.MultiWriter(os.Stdout, appLogFile)
	log.SetOutput(multiWriter)
	log.SetFlags(log.LstdFlags | log.Lshortfile)

	return nil
}

func Record(slug, action, status string, err error) {
	if err == nil {
		AuditLog.Info("Action",
			slog.String("user", slug),
			slog.String("action", action),
			slog.String("status", status),
		)
	} else {
		AuditLog.Error("Action",
			slog.String("user", slug),
			slog.String("action", action),
			slog.String("status", status),
			slog.String("error", err.Error()),
		)
	}
}
