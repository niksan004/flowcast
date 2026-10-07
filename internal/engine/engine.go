package engine

import (
	"bytes"
	"fmt"
	"github.com/expr-lang/expr"
	"golang.org/x/crypto/ssh"
	"log/slog"
	"project/internal/agentbin"
	"project/internal/logger"
	"project/internal/protocol"
	"project/internal/steps"
	"sync"
)

type Engine struct {
	cfg *Config
	rt  *Runtime
}

func NewEngine(cfg *Config, rt *Runtime) *Engine {
	return &Engine{
		cfg: cfg,
		rt:  rt,
	}
}

// execute steps for each host
func runSteps(stepsSlice []steps.RawStep, client *ssh.Client, env map[string]any, log *slog.Logger) error {
	for _, step := range stepsSlice {
		// render templates in Data
		renderedData, err := renderData(step.Data, env)
		if err != nil {
			return fmt.Errorf("failed to render data for %s: %w", step.Name, err)
		}

		// get specific step(StepExecutor) e.g. EchoStep
		exec, err := steps.Decode(step.Name, renderedData)
		if err != nil {
			return fmt.Errorf("error while decoding step %s: %w", step.Name, err)
		}

		log.Info("{Executing} " + logger.StringifyStruct(exec))

		// determine how to execute next steps; normal step or control flow
		switch exec := exec.(type) {
		// check if there is a control flow structure
		case steps.ControlFlow:
			if err := exec.Run(env, func(nextSteps []steps.RawStep) error {
				return runSteps(nextSteps, client, env, log)
			}); err != nil {
				return fmt.Errorf("error while executing control flow step %s: %w", step.Name, err)
			}
		// normal step
		default:
			// create session for step
			sesh, err := client.NewSession()
			if err != nil {
				return fmt.Errorf("error while creating session for step %s: %w", step.Name, err)
			}

			// execute step
			returnVals, err := exec.Execute(sesh)
			log.Info("return values: " + logger.StringifyStruct(returnVals))
			if err != nil {
				sesh.Close()
				return fmt.Errorf("failed to execute step %s: %w", step.Name, err)
			}
			sesh.Close()

			// save workflow variables if there are any
			if len(step.SaveAs) != 0 {
				// temporarily add return vals to env for evaluating save_as
				env["result"] = returnVals
				// evaluate save_as
				for varName, path := range step.SaveAs {
					val, err := expr.Eval(path, env)
					if err != nil {
						return fmt.Errorf("failed to evaluate save_as condition for %s: %w", step.Name, err)
					}
					env[varName] = val
				}
				// remove return vals
				delete(env, "result")
			}
			log.Info("env: " + logger.StringifyStruct(env))
		}
	}

	return nil
}

func (eng *Engine) SetupAndRunSteps(host Host) error {
	// gloal environment for the workflow
	env := map[string]any{}

	// per-host logger
	hostLog := slog.With("host", fmt.Sprintf("%s:%s", host.Ip, host.Port))

	// client used for ssh
	conn, err := eng.rt.SshCl.Connect(host.Ip, host.User, host.Port)
	if err != nil {
		return fmt.Errorf("error while trying to connect to host: %w", err)
	}
	hostLog.Info("{Connected to} " + logger.StringifyStruct(host))
	defer conn.Close()

	// sftp connection to copy agent to remote
	if err := eng.rt.SftpCl.Upload(conn, bytes.NewReader(agentbin.Linux_amd64), protocol.AgentRemotePath, 0755); err != nil {
		return fmt.Errorf("error while tring to build sftp connection: %w", err)
	}

	return runSteps(eng.cfg.wf.Steps, conn, env, hostLog)
}

type EngineResult struct {
	Host Host
	Err  error
}

func (eng *Engine) RunEngine() []EngineResult {
	var wg sync.WaitGroup
	resultCh := make(chan EngineResult, len(eng.cfg.inv.Hosts))

	// iterate through hosts
	for _, host := range eng.cfg.inv.Hosts {
		// run steps concurrently
		wg.Add(1)
		go func(h Host) {
			defer wg.Done()
			err := eng.SetupAndRunSteps(h)
			resultCh <- EngineResult{h, err}
		}(host)
	}

	// blocks until all routines call Done()
	wg.Wait()

	results := make([]EngineResult, 0, len(eng.cfg.inv.Hosts))
	for range eng.cfg.inv.Hosts {
		results = append(results, <-resultCh)
	}

	return results
}
