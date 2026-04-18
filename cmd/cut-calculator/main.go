package main

import (
	"fmt"
	"os"

	"github.com/beekman/cut-calculator/internal/input"
	"github.com/beekman/cut-calculator/internal/joiner"
	"github.com/beekman/cut-calculator/internal/model"
	"github.com/beekman/cut-calculator/internal/output"
	"github.com/beekman/cut-calculator/internal/solver/solver1d"
	"github.com/beekman/cut-calculator/internal/solver/solver2d"
)

func main() {
	flags, err := input.Parse(os.Args[1:])
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	cfg := flags
	if flags.FilePath != "" {
		fileCfg, err := input.ParseFile(flags.FilePath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: %v\n", err)
			os.Exit(1)
		}
		cfg = input.Merge(fileCfg, flags)
	}

	if err := input.Validate(cfg); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}

	solve := func(need []model.RequiredPiece) (model.CutPlan, error) {
		if input.DetectMode(cfg) == 2 {
			return solver2d.New().Solve(cfg.Stock, need, cfg.Kerf, cfg.Rotate)
		}
		return solver1d.New().Solve(cfg.Stock, need, cfg.Kerf)
	}

	plan, solveErr := solve(cfg.Need)
	if solveErr != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", solveErr)
		os.Exit(1)
	}

	// try combined join groups if present; keep whichever plan wastes less
	if joiner.HasJoinGroups(cfg.Need) {
		combinedNeed, meta := joiner.BuildCombinedNeed(cfg.Need)
		if planCombined, err := solve(combinedNeed); err == nil {
			joiner.EnrichPlan(&planCombined, meta)
			plan = joiner.BetterPlan(planCombined, plan)
		}
	}

	w := output.New(cfg.OutputFormat)
	if err := w.Write(os.Stdout, plan); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
