package main

import (
	"flag"
	"fmt"
	"log/slog"
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

func main() {
	// parse flags
	flag.Parse()

	// get workflow filepath from command args
	filepath := flag.Arg(0)

	// create config from file
	cfg, err := engine.NewConfig(filepath)
	if err != nil {
		slog.Error(err.Error(), err)
	}
	slog.Info(fmt.Sprintf("Successfully created config from %s", filepath))

	rt, err := engine.NewRuntime(privSshKey)
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
