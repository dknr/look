package main

import (
	"bytes"
	"context"
	"flag"
	"fmt"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

func main() {
	interval := flag.Duration("n", 2*time.Second, "interval seconds or suffix (ms/s/m/h)")
	intervalP := flag.Float64("interval", 0, "interval seconds or suffix (ms/s/m/h)")
	errexit := flag.Bool("e", false, "exit if exit code is non-zero")
	chgexit := flag.Bool("g", false, "exit when output changes")
	noTitle := flag.Bool("t", false, "don't display the header title")
	help := flag.Bool("h", false, "show this help")
	flag.Parse()

	if *help {
		flag.Usage()
		return
	}

	var ival time.Duration
	if *intervalP > 0 {
		ival = time.Duration(*intervalP * float64(time.Second))
	} else {
		ival = *interval
	}

	if ival <= 0 {
		fmt.Fprintln(os.Stderr, "interval must be positive")
		os.Exit(1)
	}

	cmdArgs := flag.Args()
	if len(cmdArgs) == 0 {
		fmt.Fprintln(os.Stderr, "no command specified")
		flag.Usage()
		os.Exit(1)
	}

	cmdStr := strings.Join(cmdArgs, " ")
	if !*noTitle {
		fmt.Printf("Every %.4g%s: %s\n", ival.Seconds(), durationSuffix(ival), cmdStr)
	}

	ctx, cancel := contextWithSignals()
	defer cancel()

	var prevOutput []byte
	for {
		select {
		case <-ctx.Done():
			fmt.Println()
			return
		default:
		}

		buf := &bytes.Buffer{}
		cmd := exec.Command("sh", "-c", cmdStr)
		cmd.Stdout = buf
		cmd.Stderr = os.Stderr

		err := cmd.Run()
		exitCode := 0
		if err != nil {
			if exitError, ok := err.(*exec.ExitError); ok {
				exitCode = exitError.ExitCode()
			} else {
				fmt.Fprintf(os.Stderr, "look: %v\n", err)
				os.Exit(1)
			}
		}

		currentOutput := buf.Bytes()

		if *chgexit && len(prevOutput) > 0 && !bytes.Equal(prevOutput, currentOutput) {
			fmt.Print(buf.String())
			return
		}

		if *errexit && exitCode != 0 {
			fmt.Print(buf.String())
			return
		}

		// Double buffer: clear screen, then write all output at once
		fmt.Printf("\033[2J\033[H")
		fmt.Print(buf.String())

		prevOutput = currentOutput

		select {
		case <-ctx.Done():
			fmt.Println()
			return
		case <-time.After(ival):
		}
	}
}

func durationSuffix(d time.Duration) string {
	if d >= time.Hour {
		return "h"
	}
	if d >= time.Minute {
		return "m"
	}
	if d >= time.Second {
		return "s"
	}
	return "ms"
}

func contextWithSignals() (context.Context, context.CancelFunc) {
	ctx, cancel := context.WithCancel(context.Background())

	go func() {
		sig := make(chan os.Signal, 1)
		signal.Notify(sig, syscall.SIGINT, syscall.SIGTERM)
		<-sig
		cancel()
	}()

	return ctx, cancel
}
