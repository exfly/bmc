package cmd

import (
	"bytes"
	"context"
	"io"
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/exfly/bmc/pkg/debug"
	"github.com/pkg/errors"
)

func execCmd(ctx context.Context, bin string, args ...string) (string, error) {
	envs := os.Environ()

	debug := debug.Init("exec_cmd")

	debug("[DEBUG] %s %s", bin, strings.Join(args, " "))

	cmd := exec.CommandContext(
		ctx,
		bin,
		args...,
	)

	buf := bytes.NewBuffer(nil)
	var multWriter io.Writer = buf

	cmd.Env = envs
	cmd.Stdin = os.Stdin
	cmd.Stdout = multWriter
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return "", errors.Wrap(err, "")
	}

	return buf.String(), nil
}

func execCmdInteract(ctx context.Context, dir string, bin string, args ...string) error {
	envs := os.Environ()

	log.Println("env", envs)

	debug := debug.Init("exec_cmd")

	debug("[DEBUG] %s %s", bin, strings.Join(args, " "))

	cmd := exec.CommandContext(
		ctx,
		bin,
		args...,
	)

	cmd.Dir = dir
	cmd.Env = envs
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	if err := cmd.Run(); err != nil {
		return errors.Wrap(err, "")
	}

	return nil
}

func fileExists(filename string) bool {
	info, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}
	return !info.IsDir()
}
