package main

import (
	"context"
	"flag"
	"os"
	"os/exec"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

func main() {
	err := run()
	if err != nil {
		logrus.Fatal(err)
	}
}

func run() error {
	timeout, interval, command, args, err := parseArguments()
	if err != nil {
		return err
	}
	return travisWait(timeout, interval, command, args...)
}

func parseArguments() (time.Duration, time.Duration, string, []string, error) {
	timeoutString := flag.String("timeout", "20m", "Timeout for this command")
	intervalString := flag.String("interval", "1m", "The interval at which to print keep-alive messages")

	flag.Parse()
	args := flag.Args()

	timeout, err := time.ParseDuration(*timeoutString)
	if err != nil {
		return timeout, 0, "", nil, errors.Wrap(err, "could not parse timeout as a duration")
	}
	interval, err := time.ParseDuration(*intervalString)
	if err != nil {
		return timeout, interval, "", nil, errors.Wrap(err, "could not parse interval as duration")
	}

	if len(args) < 1 {
		return timeout, interval, "", nil, errors.New("could not parse command to run")
	}

	return timeout, interval, args[0], args[1:], nil
}

func travisWait(timeout time.Duration, interval time.Duration, command string, args ...string) error {
	ticker := time.NewTicker(interval)
	go func() {
		for t := range ticker.C {
			logrus.Infof("go-travis-wait waiting at %s...", t.Format(time.RFC1123Z))
		}
	}()
	defer ticker.Stop()

	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer cancel()

	cmd := exec.CommandContext(ctx, command, args...)
	// Redirect output
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr

	err := cmd.Run()
	if ctx.Err() == context.DeadlineExceeded {
		return errors.New("timeout exceeded, shutting down")
	}
	return errors.Wrap(err, "non-zero exit code")
}
