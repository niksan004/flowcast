package main

import (
	"fmt"
	"os"
	"project/internal/engine"
	"project/internal/logger"
)

func main() {
	// get filepath from command args
	filepath := os.Args[1]

	// create config from file
	cfg, err := engine.NewConfig(filepath)
	if err != nil {
		logger.Error(err.Error(), err)
	}
	logger.Info(fmt.Sprintf("Successfully created config from %s", filepath))

	rt, err := engine.NewRuntime()
	if err != nil {
		logger.Error(err.Error(), err)
	}
	logger.Info("Successfully created runtime")

	eng := engine.NewEngine(cfg, rt)
	res := eng.RunEngine()
	for _, v := range res {
		logger.Info(fmt.Sprintf("Host: %s, Err: %w", v.Host, v.Err))
	}
	logger.Info("Successfully ran workflow")
}
