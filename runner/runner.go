package runner

import (
	"errors"
	"os/exec"
)

type Runner struct {
	Command string
	Args    []string
	Output  string
}

func New(cmd string, args ...string) *Runner {
	return &Runner{
		Command: cmd,
		Args:    args,
	}
}

func (r *Runner) Execute() error {
	if r.Command == "" {
		return errors.New("command cannot be empty")
	}

	cmd := exec.Command(r.Command, r.Args...)
	
	output, err := cmd.CombinedOutput()
	r.Output = string(output)
	if err != nil {
		return err
	}
	
	return nil
}
