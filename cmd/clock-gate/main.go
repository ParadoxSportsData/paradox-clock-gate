// Package main is the clock-gate CLI entry point for temporal NFL game-state queries.
package main

import (
	"flag"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/ParadoxSportsData/paradox-clock-gate/internal/gate"
	"github.com/ParadoxSportsData/paradox-clock-gate/internal/ingestion"
	"github.com/ParadoxSportsData/paradox-clock-gate/internal/matrix"
	"github.com/ParadoxSportsData/paradox-clock-gate/internal/presenter"
)

func main() {
	// Serve subcommand: clock-gate serve [--port N] [--data DIR]
	if len(os.Args) > 1 && os.Args[1] == "serve" {
		runServe(os.Args[2:])
		return
	}

	tick := flag.Int("tick", -1, "elapsed seconds since kickoff to query")
	format := flag.String("format", "text", "output format: text or json")
	list := flag.Bool("list", false, "list game files in the given directory")
	flag.Usage = func() {
		fmt.Fprintln(os.Stderr, "Usage: clock-gate --tick <seconds> [--format text|json] <game-file>")
		fmt.Fprintln(os.Stderr, "       clock-gate --list <directory>")
		flag.PrintDefaults()
	}
	flag.Parse()

	args := flag.Args()

	if *list {
		if len(args) < 1 {
			fmt.Fprintln(os.Stderr, "error: --list requires a directory argument")
			os.Exit(1)
		}
		if err := listGames(args[0]); err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		return
	}

	if *tick < 0 {
		fmt.Fprintln(os.Stderr, "error: --tick is required (non-negative integer)")
		flag.Usage()
		os.Exit(1)
	}
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "error: game file path required")
		flag.Usage()
		os.Exit(1)
	}

	header, plays, err := ingestion.ParseFile(args[0])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	m := matrix.Compile(plays, header)

	if err := gate.Validate(*tick, m.Meta.MaxTick); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	gs := m.States[*tick]

	switch *format {
	case "json":
		fmt.Println(presenter.RenderJSON(gs, m.Meta, m.Arena))
	default:
		fmt.Print(presenter.RenderText(gs, m.Meta, m.Arena))
	}
}

func listGames(dir string) error {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return fmt.Errorf("read directory %s: %w", dir, err)
	}
	for _, e := range entries {
		if e.IsDir() || !strings.HasSuffix(e.Name(), ".json") {
			continue
		}
		gameID := strings.TrimSuffix(e.Name(), ".json")
		fmt.Printf("%s\t%s\n", gameID, filepath.Join(dir, e.Name()))
	}
	return nil
}
