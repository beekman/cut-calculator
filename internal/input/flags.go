package input

import (
	"flag"
	"fmt"
	"strconv"
	"strings"

	"github.com/beekman/cut-calculator/internal/model"
)

type Config struct {
	Stock        []model.StockPiece
	Need         []model.RequiredPiece
	Kerf         float64
	OutputFormat model.OutputFormat
}

type stringSlice []string

func (s *stringSlice) String() string { return strings.Join(*s, ", ") }
func (s *stringSlice) Set(v string) error {
	*s = append(*s, v)
	return nil
}

func Parse(args []string) (*Config, error) {
	fs := flag.NewFlagSet("cut-calculator", flag.ContinueOnError)

	var stockFlags stringSlice
	var needFlags stringSlice
	var kerf float64
	var outputFormat string

	fs.Var(&stockFlags, "stock", "Stock piece: LENGTH[:COUNT][:onhand] e.g. 96:3:onhand")
	fs.Var(&needFlags, "need", "Required piece: LENGTH[:COUNT] e.g. 36:4")
	fs.Float64Var(&kerf, "kerf", 0, "Blade kerf width in inches")
	fs.StringVar(&outputFormat, "output", "text", "Output format: text, json")

	if err := fs.Parse(args); err != nil {
		return nil, err
	}

	cfg := &Config{
		Kerf:         kerf,
		OutputFormat: model.OutputFormat(outputFormat),
	}

	for _, s := range stockFlags {
		p, err := parseStock(s)
		if err != nil {
			return nil, fmt.Errorf("invalid --stock %q: %w", s, err)
		}
		cfg.Stock = append(cfg.Stock, p)
	}

	labels := labelSeq()
	for _, n := range needFlags {
		p, err := parseNeed(n, labels)
		if err != nil {
			return nil, fmt.Errorf("invalid --need %q: %w", n, err)
		}
		cfg.Need = append(cfg.Need, p)
	}

	return cfg, nil
}

func parseStock(s string) (model.StockPiece, error) {
	parts := strings.Split(s, ":")
	if len(parts) > 3 {
		return model.StockPiece{}, fmt.Errorf("too many fields")
	}

	length, err := parsePositiveFloat(parts[0])
	if err != nil {
		return model.StockPiece{}, err
	}

	p := model.StockPiece{Length: length, Count: 1}

	for _, part := range parts[1:] {
		if part == "onhand" {
			p.OnHand = true
			continue
		}
		n, err := strconv.Atoi(part)
		if err != nil || n < 1 {
			return model.StockPiece{}, fmt.Errorf("invalid count %q", part)
		}
		p.Count = n
	}

	return p, nil
}

func parseNeed(s string, next func() string) (model.RequiredPiece, error) {
	parts := strings.Split(s, ":")
	if len(parts) > 2 {
		return model.RequiredPiece{}, fmt.Errorf("too many fields")
	}

	length, err := parsePositiveFloat(parts[0])
	if err != nil {
		return model.RequiredPiece{}, err
	}

	p := model.RequiredPiece{Length: length, Count: 1, Label: next()}

	if len(parts) == 2 {
		n, err := strconv.Atoi(parts[1])
		if err != nil || n < 1 {
			return model.RequiredPiece{}, fmt.Errorf("invalid count %q", parts[1])
		}
		p.Count = n
	}

	return p, nil
}

func parsePositiveFloat(s string) (float64, error) {
	if s == "" {
		return 0, fmt.Errorf("empty value")
	}
	v, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, fmt.Errorf("not a number: %q", s)
	}
	if v <= 0 {
		return 0, fmt.Errorf("must be greater than zero, got %v", v)
	}
	return v, nil
}

func labelSeq() func() string {
	i := 0
	return func() string {
		label := string(rune('A' + i%26))
		if i >= 26 {
			label = string(rune('A'+i/26-1)) + label
		}
		i++
		return label
	}
}
