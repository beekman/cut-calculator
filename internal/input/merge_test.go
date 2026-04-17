package input

import (
	"testing"

	"github.com/beekman/cut-calculator/internal/model"
)

func fileConfig() *Config {
	return &Config{
		Kerf:         0.125,
		Rotate:       true,
		OutputFormat: model.OutputText,
		Stock: []model.StockPiece{
			{Length: 96, Count: 2},
		},
		Need: []model.RequiredPiece{
			{Label: "A", Length: 36, Count: 1},
			{Label: "B", Length: 24, Count: 2},
		},
	}
}

func emptyFlags() *Config {
	return &Config{Rotate: true}
}

func TestMerge_FlagsOverrideKerf(t *testing.T) {
	flags := emptyFlags()
	flags.Kerf = 0.25

	out := Merge(fileConfig(), flags)
	if out.Kerf != 0.25 {
		t.Errorf("kerf: got %v, want 0.25", out.Kerf)
	}
}

func TestMerge_FileKerfWhenFlagsZero(t *testing.T) {
	out := Merge(fileConfig(), emptyFlags())
	if out.Kerf != 0.125 {
		t.Errorf("kerf: got %v, want 0.125", out.Kerf)
	}
}

func TestMerge_NoRotateFlagDisablesRotate(t *testing.T) {
	flags := emptyFlags()
	flags.Rotate = false

	out := Merge(fileConfig(), flags)
	if out.Rotate {
		t.Error("rotate should be false when --no-rotate flag is set")
	}
}

func TestMerge_RotateTruePreservesFileValue(t *testing.T) {
	file := fileConfig()
	file.Rotate = false
	flags := emptyFlags() // flags.Rotate = true means --no-rotate was NOT passed

	out := Merge(file, flags)
	if out.Rotate {
		t.Error("file rotate=false should be preserved when --no-rotate not passed")
	}
}

func TestMerge_FlagsOverrideOutputFormat(t *testing.T) {
	flags := emptyFlags()
	flags.OutputFormat = model.OutputJSON

	out := Merge(fileConfig(), flags)
	if out.OutputFormat != model.OutputJSON {
		t.Errorf("output: got %q, want json", out.OutputFormat)
	}
}

func TestMerge_FileOutputWhenFlagsEmpty(t *testing.T) {
	out := Merge(fileConfig(), emptyFlags())
	if out.OutputFormat != model.OutputText {
		t.Errorf("output: got %q, want text", out.OutputFormat)
	}
}

func TestMerge_StockCombined(t *testing.T) {
	flags := emptyFlags()
	flags.Stock = []model.StockPiece{{Length: 48, Count: 1}}

	out := Merge(fileConfig(), flags)
	if len(out.Stock) != 2 {
		t.Fatalf("stock len: got %d, want 2", len(out.Stock))
	}
	if out.Stock[0].Length != 96 {
		t.Errorf("stock[0].Length: got %v, want 96", out.Stock[0].Length)
	}
	if out.Stock[1].Length != 48 {
		t.Errorf("stock[1].Length: got %v, want 48", out.Stock[1].Length)
	}
}

func TestMerge_NeedCombinedAndRelabeled(t *testing.T) {
	flags := emptyFlags()
	flags.Need = []model.RequiredPiece{
		{Label: "X", Length: 12, Count: 1},
	}

	out := Merge(fileConfig(), flags)
	if len(out.Need) != 3 {
		t.Fatalf("need len: got %d, want 3", len(out.Need))
	}
	wantLabels := []string{"A", "B", "C"}
	for i, p := range out.Need {
		if p.Label != wantLabels[i] {
			t.Errorf("need[%d].Label: got %q, want %q", i, p.Label, wantLabels[i])
		}
	}
}

func TestMerge_NeedFlagsOnlyStillRelabeled(t *testing.T) {
	file := &Config{
		Kerf:         0,
		Rotate:       true,
		OutputFormat: model.OutputText,
		Stock:        []model.StockPiece{{Length: 96, Count: 1}},
	}
	flags := emptyFlags()
	flags.Need = []model.RequiredPiece{
		{Label: "Z", Length: 20, Count: 1},
		{Label: "Y", Length: 10, Count: 1},
	}

	out := Merge(file, flags)
	wantLabels := []string{"A", "B"}
	for i, p := range out.Need {
		if p.Label != wantLabels[i] {
			t.Errorf("need[%d].Label: got %q, want %q", i, p.Label, wantLabels[i])
		}
	}
}
