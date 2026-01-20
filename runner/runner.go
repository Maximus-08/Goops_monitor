package runner

import (
	"errors"
	"log"
	"os/exec"
)

type Runner struct {
	Command string
	Args    []string
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
	log.Printf("Executing command: %s", r.Command)
	
	output, err := cmd.CombinedOutput()
	if err != nil {
		return err
	}
	
	log.Printf("Output: %s", string(output))
	return nil
}
