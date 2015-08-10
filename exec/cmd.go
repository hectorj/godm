package exec

import (
	"bytes"
	"fmt"
	"os/exec"
	"strings"
)

type Result interface {
	GetError() error
	GetStdout() []byte
	GetStderr() []byte
}

type result struct {
	err    error
	cmd    string
	stdout []byte
	stderr []byte
}

func (r result) GetError() error {
	if r.err != nil {
		return fmt.Errorf("Failure : %q\n== Output\n%s\n== Err\n%s\nError : %v", r.cmd, string(r.stdout), string(r.stderr), r.err)
	} else {
		return nil
	}
}

func (r result) GetStdout() []byte {
	return r.stdout
}
func (r result) GetStderr() []byte {
	return r.stderr
}

func Cmd(dir string, name string, args ...string) Result {
	r := result{
		cmd: name + " " + strings.Join(args, " "),
	}

	cmd := exec.Command(name, args...)
	cmd.Dir = dir

	bufout := &bytes.Buffer{}
	buferr := &bytes.Buffer{}

	cmd.Stdout = bufout
	cmd.Stderr = buferr

	r.err = cmd.Run()

	r.stdout = bufout.Bytes()
	r.stderr = buferr.Bytes()
	return r
}
