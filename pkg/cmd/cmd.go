package cmd

import (
	"io"
	"os"
	"os/exec"
)

type Cmd struct {
	nameOrPath string
	workDir    string
	args       []string
	in         io.Reader
	out        io.Writer
	err        io.Writer
}

/*
type saveOutput struct {
	savedOutput []byte
}

func (so *saveOutput) Write(p []byte) (n int, err error) {
	so.savedOutput = append(so.savedOutput, p...)
	return os.Stdout.Write(p)
}
*/

func New(nameOrPath string, workDir string, args ...string) *Cmd {
	// var so saveOutput
	return &Cmd{
		nameOrPath: nameOrPath,
		workDir:    workDir,
		args:       args,
		in:         os.Stdin,
		out:        os.Stdout,
		err:        os.Stderr,
	}
}

func (c *Cmd) Execute() error {
	cmd := exec.Command(c.nameOrPath, c.args...)
	cmd.Dir = c.workDir
	cmd.Stdin = c.in
	cmd.Stdout = c.out
	cmd.Stderr = c.err
	if err := cmd.Start(); err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}
	return nil
}
