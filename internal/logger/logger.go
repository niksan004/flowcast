package logger

import (
	"fmt"
	"log/slog"
	"os"
)

func init() {
	opts := &slog.HandlerOptions{
		Level: slog.LevelInfo,
	}
	handler := slog.NewTextHandler(os.Stdout, opts)
	slog.SetDefault(slog.New(handler))
}

func StringifyStruct(s any) string {
	return fmt.Sprintf("%#v", s)
}
