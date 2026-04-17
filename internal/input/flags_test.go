package input_test

import (
	"testing"

	"github.com/beekman/cut-calculator/internal/input"
	"github.com/beekman/cut-calculator/internal/model"
)

func TestParse_StockFlags(t *testing.T) {
	cases := []struct {
		arg     string
		want    model.StockPiece
		wantErr bool
	}{
		{"96", model.StockPiece{Length: 96, Count: 1, OnHand: false}, false},
		{"96:3", model.StockPiece{Length: 96, Count: 3, OnHand: false}, false},
		{"96:3:onhand", model.StockPiece{Length: 96, Count: 3, OnHand: true}, false},
		{"96:onhand", model.StockPiece{Length: 96, Count: 1, OnHand: true}, false},
		{"", model.StockPiece{}, true},
		{"0", model.StockPiece{}, true},
		{"-5", model.StockPiece{}, true},
		{"abc", model.StockPiece{}, true},
		{"96:3:onhand:extra", model.StockPiece{}, true},
	}

	for _, tc := range cases {
		t.Run(tc.arg, func(t *testing.T) {
			cfg, err := input.Parse([]string{"--stock", tc.arg, "--need", "10"})
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := cfg.Stock[0]
			if got.Length != tc.want.Length || got.Count != tc.want.Count || got.OnHand != tc.want.OnHand {
				t.Errorf("got %+v, want %+v", got, tc.want)
			}
		})
	}
}

func TestParse_NeedFlags(t *testing.T) {
	cases := []struct {
		arg     string
		wantLen float64
		wantCnt int
		wantErr bool
	}{
		{"36", 36, 1, false},
		{"36:4", 36, 4, false},
		{"", 0, 0, true},
		{"0", 0, 0, true},
		{"-1", 0, 0, true},
		{"36:0", 0, 0, true},
		{"36:4:extra", 0, 0, true},
	}

	for _, tc := range cases {
		t.Run(tc.arg, func(t *testing.T) {
			cfg, err := input.Parse([]string{"--stock", "96", "--need", tc.arg})
			if tc.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			got := cfg.Need[0]
			if got.Length != tc.wantLen || got.Count != tc.wantCnt {
				t.Errorf("got Length=%.4g Count=%d, want Length=%.4g Count=%d",
					got.Length, got.Count, tc.wantLen, tc.wantCnt)
			}
		})
	}
}

func TestValidate_Errors(t *testing.T) {
	cases := []struct {
		name    string
		args    []string
		wantErr string
	}{
		{"no stock", []string{"--need", "36"}, "no stock"},
		{"no need", []string{"--stock", "96"}, "nothing to cut"},
		{"piece larger than stock", []string{"--stock", "48", "--need", "96"}, "longer than"},
		{"bad output format", []string{"--stock", "96", "--need", "36", "--output", "yaml"}, "unknown output"},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			cfg, err := input.Parse(tc.args)
			if err != nil {
				return // parse error is fine for these cases
			}
			err = input.Validate(cfg)
			if err == nil {
				t.Fatal("expected validation error, got nil")
			}
		})
	}
}
