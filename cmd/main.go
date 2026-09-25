package main

import (
	"fmt"
	"log/slog"
	"os"
	"project/internal/engine"
	_ "project/internal/logger"
)

func main() {
	// get filepath from command args
	filepath := os.Args[1]

	// create config from file
	cfg, err := engine.NewConfig(filepath)
	if err != nil {
		slog.Error(err.Error(), err)
	}
	slog.Info(fmt.Sprintf("Successfully created config from %s", filepath))

	rt, err := engine.NewRuntime()
	if err != nil {
		slog.Error(err.Error(), err)
	}
	slog.Info("Successfully created runtime")

	eng := engine.NewEngine(cfg, rt)
	res := eng.RunEngine()
	for _, v := range res {
		slog.Info(fmt.Sprintf("Host: %s:%s, Err: %w", v.Host, v.Err))
	}
	slog.Info("Successfully ran workflow")
}
