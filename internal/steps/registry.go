package steps

type StepExecutor interface {
	Execute() error
}

type StepFactory func(data map[string]any) (StepExecutor, error)

var Registry = map[string]StepFactory{
	"echo": parseEcho,
	"http": parseHTTP,
}
