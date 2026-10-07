package steps

type ControlFlow interface {
	Run(env map[string]any, runBody func([]RawStep) error) error
}
