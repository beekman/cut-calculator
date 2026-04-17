package output_test

import (
	"bytes"
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/beekman/cut-calculator/internal/model"
	"github.com/beekman/cut-calculator/internal/output"
)

var update = flag.Bool("update", false, "regenerate golden files")

// plan1D is a fixed 1D cut plan used for all golden file comparisons.
// Stock: 96" board (on hand). Need: A(36")×2, B(22")×1. Kerf: 0.125".
var plan1D = model.CutPlan{
	Mode: 1,
	Results: []model.StockResult{
		{
			Stock: model.StockPiece{Length: 96, Count: 1, OnHand: true},
			Assignments: []model.Assignment{
				{StockIndex: 0, RequiredLabel: "A", Length: 36},
				{StockIndex: 0, RequiredLabel: "A", Length: 36},
				{StockIndex: 0, RequiredLabel: "B", Length: 22},
			},
			Cuts: []model.Cut{
				{Position: 0, Label: "A"},
				{Position: 36.125, Label: "A"},
				{Position: 72.25, Label: "B"},
			},
			WasteLength: 1.75,
		},
	},
	WastePct:  2.0833,
	Purchased: nil,
	Unfit:     nil,
}

// plan2D is a fixed 2D cut plan: four 24×48 pieces filling a 48×96 sheet exactly.
var plan2D = model.CutPlan{
	Mode: 2,
	Results: []model.StockResult{
		{
			Stock: model.StockPiece{Width: 48, Height: 96, Count: 1, OnHand: true},
			Assignments: []model.Assignment{
				{StockIndex: 0, RequiredLabel: "A", Width: 24, Height: 48, OffsetX: 0, OffsetY: 0},
				{StockIndex: 0, RequiredLabel: "A", Width: 24, Height: 48, OffsetX: 24, OffsetY: 0},
				{StockIndex: 0, RequiredLabel: "A", Width: 24, Height: 48, OffsetX: 0, OffsetY: 48},
				{StockIndex: 0, RequiredLabel: "A", Width: 24, Height: 48, OffsetX: 24, OffsetY: 48},
			},
			WasteArea: 0,
		},
	},
	WastePct:  0,
	Purchased: nil,
	Unfit:     nil,
}

func TestGolden(t *testing.T) {
	cases := []struct {
		name   string
		format model.OutputFormat
		plan   model.CutPlan
	}{
		{"simple_1d", model.OutputASCII, plan1D},
		{"simple_1d", model.OutputText, plan1D},
		{"simple_1d", model.OutputJSON, plan1D},
		{"simple_2d", model.OutputASCII, plan2D},
		{"simple_2d", model.OutputText, plan2D},
		{"simple_2d", model.OutputJSON, plan2D},
	}

	for _, tc := range cases {
		ext := map[model.OutputFormat]string{
			model.OutputASCII: "ascii.txt",
			model.OutputText:  "text.txt",
			model.OutputJSON:  "json",
		}[tc.format]
		name := tc.name + "." + ext

		t.Run(name, func(t *testing.T) {
			var buf bytes.Buffer
			w := output.New(tc.format)
			if err := w.Write(&buf, tc.plan); err != nil {
				t.Fatalf("Write: %v", err)
			}
			got := buf.Bytes()

			path := filepath.Join("testdata", "golden", name)

			if *update {
				if err := os.WriteFile(path, got, 0644); err != nil {
					t.Fatalf("writing golden file: %v", err)
				}
				t.Logf("updated %s", path)
				return
			}

			want, err := os.ReadFile(path)
			if err != nil {
				t.Fatalf("reading golden file %s: %v\n(run with -update to generate)", path, err)
			}
			if !bytes.Equal(got, want) {
				t.Errorf("output mismatch for %s\ngot:\n%s\nwant:\n%s", name, got, want)
			}
		})
	}
}
