package main

import (
	"flag"
	"fmt"
	"log/slog"
	"os"
	"project/internal/engine"
	_ "project/internal/logger"
)

var privSshKey string

func init() {
	const (
		defaultPrivSshKey = "id_ed25519"
		usage             = `private ssh key used to establish ssh connection to remote hosts;
if key name contains '/' it will be used as a path,
if not an ssh key with the given name will be searched for in ~/.ssh`
	)

	flag.StringVar(&privSshKey, "ssh_key", defaultPrivSshKey, usage)
	flag.StringVar(&privSshKey, "k", defaultPrivSshKey, usage+" (shorthand)")
}

func run() error {
	if flag.NArg() == 0 {
		flag.Usage()
		return fmt.Errorf("not enough arguments")
	}

	// get workflow filepath from command args
	filepath := flag.Arg(0)

	// create config from file
	cfg, err := engine.NewConfig(filepath)
	if err != nil {
		return fmt.Errorf("failed to create config from %s: %w", filepath, err)
	}
	slog.Info("successfully created config from", "filepath", filepath)

	rt, err := engine.NewRuntime(privSshKey)
	if err != nil {
		return fmt.Errorf("failed to create runtime: %w", err)
	}
	slog.Info("successfully created runtime")

	eng := engine.NewEngine(cfg, rt)
	res := eng.RunEngine()
	failedHost := false
	for _, v := range res {
		if v.Err != nil {
			failedHost = true
			slog.Error("host finished", "host", v.Host, "err", v.Err)
		} else {
			slog.Info("host finished", "host", v.Host)
		}
	}

	if failedHost {
		return fmt.Errorf("a host has failed")
	}
	slog.Info("successfully ran workflow")

	return nil
}

func main() {
	// parse flags
	flag.Parse()

	// run program
	if err := run(); err != nil {
		slog.Error("run failed", "err", err)
		os.Exit(1)
	}
}
