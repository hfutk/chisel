package main

import (
	"fmt"
	"io"
	"os"
	"strings"
)

// argsStdinEnv enables supplying chisel's command line arguments through
// standard input. When this environment variable is set (to any value other
// than "0"), all of stdin is read during package init, parsed into
// arguments and spliced into os.Args before any flag/subcommand parsing
// happens in main().
//
// Parsing: if the input contains newlines, each non-empty trimmed line is
// taken as one argument; otherwise the single line is split on whitespace.
//
// This is intentionally opt-in: chisel can tunnel data through
// stdin/stdout (stdio remotes, `ssh -o ProxyCommand` mode). Only enable
// this when stdin is not carrying tunnel data.
const argsStdinEnv = "CHISEL_ARGS_STDIN"

func init() {
	v, ok := os.LookupEnv(argsStdinEnv)
	if !ok || v == "" || v == "0" {
		return
	}
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "chisel: reading arguments from stdin: %v\n", err)
		return
	}
	var args []string
	if strings.Contains(string(data), "\n") {
		for _, line := range strings.Split(string(data), "\n") {
			if line = strings.TrimSpace(line); line != "" {
				args = append(args, line)
			}
		}
	} else {
		args = strings.Fields(string(data))
	}
	if len(args) == 0 {
		return
	}
	// stdin args go first; original CLI args are appended after so that
	// repeated flags on the real command line may still override them.
	prog := os.Args[0]
	extra := os.Args[1:]
	os.Args = append([]string{prog}, args...)
	os.Args = append(os.Args, extra...)
}
