package input

import (
	"errors"
	"fmt"

	"github.com/beekman/cut-calculator/internal/model"
)

func Validate(cfg *Config) error {
	if len(cfg.Stock) == 0 {
		return errors.New("no stock provided: use --stock to specify at least one piece")
	}
	if len(cfg.Need) == 0 {
		return errors.New("nothing to cut: use --need to specify required pieces")
	}
	if cfg.Kerf < 0 {
		return fmt.Errorf("kerf must be >= 0, got %v", cfg.Kerf)
	}

	switch cfg.OutputFormat {
	case model.OutputText, model.OutputJSON, model.OutputASCII, "":
	default:
		return fmt.Errorf("unknown output format %q: use text, json, or ascii", cfg.OutputFormat)
	}

	maxStock := maxStockLength(cfg.Stock)
	for _, need := range cfg.Need {
		if need.Length > maxStock {
			return fmt.Errorf(
				"required piece %s (%.4g\") is longer than the largest stock piece (%.4g\")",
				need.Label, need.Length, maxStock,
			)
		}
	}

	return nil
}

// DetectMode returns 1 for 1D input, 2 for 2D input.
// 2D mode is triggered when any stock or required piece has non-zero Width or Height.
func DetectMode(cfg *Config) int {
	for _, s := range cfg.Stock {
		if s.Width > 0 || s.Height > 0 {
			return 2
		}
	}
	for _, n := range cfg.Need {
		if n.Width > 0 || n.Height > 0 {
			return 2
		}
	}
	return 1
}

func maxStockLength(stock []model.StockPiece) float64 {
	var max float64
	for _, s := range stock {
		if s.Length > max {
			max = s.Length
		}
	}
	return max
}
