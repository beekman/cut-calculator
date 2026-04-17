package input

import (
	"fmt"
	"os"

	"gopkg.in/yaml.v3"

	"github.com/beekman/cut-calculator/internal/model"
)

type yamlFile struct {
	Kerf   float64      `yaml:"kerf"`
	Rotate *bool        `yaml:"rotate"`
	Stock  []yamlStock  `yaml:"stock"`
	Need   []yamlNeed   `yaml:"need"`
}

type yamlStock struct {
	Length         float64 `yaml:"length"`
	Width          float64 `yaml:"width"`
	Height         float64 `yaml:"height"`
	Count          int     `yaml:"count"`
	OnHand         bool    `yaml:"on_hand"`
	RepeatDistance float64 `yaml:"repeat_distance"`
	RepeatAxis     string  `yaml:"repeat_axis"`
}

type yamlNeed struct {
	Length    float64 `yaml:"length"`
	Width     float64 `yaml:"width"`
	Height    float64 `yaml:"height"`
	Count     int     `yaml:"count"`
	JoinGroup string  `yaml:"join_group"`
}

func ParseFile(path string) (*Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("reading file %q: %w", path, err)
	}

	var f yamlFile
	if err := yaml.Unmarshal(data, &f); err != nil {
		return nil, fmt.Errorf("parsing YAML in %q: %w", path, err)
	}

	cfg := &Config{
		Kerf:   f.Kerf,
		Rotate: true,
	}
	if f.Rotate != nil {
		cfg.Rotate = *f.Rotate
	}

	for _, s := range f.Stock {
		p := model.StockPiece{
			Length:         s.Length,
			Width:          s.Width,
			Height:         s.Height,
			Count:          s.Count,
			OnHand:         s.OnHand,
			RepeatDistance: s.RepeatDistance,
			RepeatAxis:     s.RepeatAxis,
		}
		if p.Count == 0 {
			p.Count = 1
		}
		cfg.Stock = append(cfg.Stock, p)
	}

	labels := labelSeq()
	for _, n := range f.Need {
		p := model.RequiredPiece{
			Label:     labels(),
			Length:    n.Length,
			Width:     n.Width,
			Height:    n.Height,
			Count:     n.Count,
			JoinGroup: n.JoinGroup,
		}
		if p.Count == 0 {
			p.Count = 1
		}
		cfg.Need = append(cfg.Need, p)
	}

	return cfg, nil
}
