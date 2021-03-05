package simpleCommand

/*
The simpleCommand library is free software: you can redistribute it and/or modify it under the terms of the Apache License.
*/

import (
	"bytes"
	"io"
	"os"
	"os/exec"
	"strconv"
	"syscall"
	"time"
)

const (
	ErrCode        int = -1
	ErrCodeProcess     = -2
)

const (
	StatNotStarted = "not started"
	StatRunning    = "in progress"
	StatCompleted  = "completed"
	StatTimeout    = "completed with timeout"
)

type Command struct {
	cmd       *exec.Cmd
	status    string
	exitCode  int
	output    *bytes.Buffer
	isTimeout bool
	timeout   time.Duration
	exitC     chan error
}

func (c *Command) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
	c.exitC = make(chan error, 1)
}

func (c *Command) SetTimeoutWithSecond(n int64) {
	c.SetTimeout(time.Duration(n) * time.Second)
}

func (c Command) hasTimeout() bool {
	return c.timeout > 0
}

func (c Command) Output() string {
	if c.output == nil {
		return ""
	}

	return c.output.String()
}

func (c Command) ExitCode() int {
	return c.exitCode
}

func (c Command) Status() string {
	return c.status
}

func (c *Command) combinedOutput(out, err io.Writer) {
	if c.output == nil {
		c.output = new(bytes.Buffer)
	}

	c.cmd.Stdout = c.output
	c.cmd.Stderr = c.output

	if out != nil && out != c.output {
		c.cmd.Stdout = io.MultiWriter(out, c.output)
	}

	if err != nil && err != c.output {
		c.cmd.Stderr = io.MultiWriter(err, c.output)
	}
}

func (c *Command) errHandle(err error) (int, string, error) {
	if err == nil {
		panic("only handle error case")
	}

	if exitError, ok := err.(*exec.ExitError); ok {
		c.exitCode = exitError.ExitCode()
	} else {
		c.exitCode = ErrCode
	}

	return c.exitCode, c.Output(), err
}

func (c *Command) run(out, err io.Writer) (int, string, error) {
	c.combinedOutput(out, err)

	if err := c.cmd.Start(); err != nil {
		return ErrCode, "", err
	}

	c.status = StatRunning
	defer func() {
		if c.isTimeout {
			c.status = StatTimeout
		} else {
			c.status = StatCompleted
		}
	}()

	if c.hasTimeout() {
		timer := time.NewTimer(c.timeout)
		defer timer.Stop()

		go func() {
			c.exitC <- c.cmd.Wait()
		}()

		select {
		case err := <-c.exitC:
			if !timer.Stop() {
				select {
				case <-timer.C:
				default:
				}
			}

			if err != nil {
				// process finished with error
				return c.errHandle(err)
			}
		case <-timer.C:
			// terminate process after a timeout
			c.isTimeout = true
			c.exitCode = ErrCodeProcess

			// process killed as timeout reached
			return c.exitCode, c.Output(), c.cmd.Process.Kill()
		}
	} else {
		if err := c.cmd.Wait(); err != nil {
			return c.errHandle(err)
		}
	}

	// process finished successfully
	ws := c.cmd.ProcessState.Sys().(syscall.WaitStatus)
	c.exitCode = ws.ExitStatus()

	return c.exitCode, c.Output(), nil
}

func (c *Command) Run() (int, string, error) {
	return c.run(nil, nil)
}

func (c *Command) RunWithOutput() (int, string, error) {
	return c.run(os.Stdout, os.Stderr)
}

func (c Command) String() string {
	buf := new(bytes.Buffer)
	buf.WriteString("Exec the command: " +
		c.cmd.String() + "\n")
	buf.WriteString("Status          : " +
		c.Status() + "\n")
	buf.WriteString("ExitCode        : " +
		strconv.FormatInt(int64(c.ExitCode()), 10) + "\n")
	buf.WriteString("Output          : " +
		c.Output() + "\n")

	return buf.String()
}

func New(name string, arg ...string) *Command {
	return &Command{
		cmd:    exec.Command(name, arg...),
		output: new(bytes.Buffer),
		status: StatNotStarted,
	}
}
