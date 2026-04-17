package input

import (
	"os"
	"path/filepath"
	"testing"
)

func writeTemp(t *testing.T, content string) string {
	t.Helper()
	f, err := os.CreateTemp(t.TempDir(), "*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	if _, err := f.WriteString(content); err != nil {
		t.Fatal(err)
	}
	f.Close()
	return f.Name()
}

func TestParseFile_Basic1D(t *testing.T) {
	yaml := `
kerf: 0.125
stock:
  - length: 96
    count: 2
need:
  - length: 36
  - length: 24
    count: 3
`
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}

	if cfg.Kerf != 0.125 {
		t.Errorf("kerf: got %v, want 0.125", cfg.Kerf)
	}
	if len(cfg.Stock) != 1 {
		t.Fatalf("stock len: got %d, want 1", len(cfg.Stock))
	}
	if cfg.Stock[0].Length != 96 || cfg.Stock[0].Count != 2 {
		t.Errorf("stock[0]: got %+v", cfg.Stock[0])
	}
	if len(cfg.Need) != 2 {
		t.Fatalf("need len: got %d, want 2", len(cfg.Need))
	}
	if cfg.Need[0].Length != 36 || cfg.Need[0].Count != 1 {
		t.Errorf("need[0]: got %+v", cfg.Need[0])
	}
	if cfg.Need[1].Length != 24 || cfg.Need[1].Count != 3 {
		t.Errorf("need[1]: got %+v", cfg.Need[1])
	}
}

func TestParseFile_LabelsSequential(t *testing.T) {
	yaml := `
stock:
  - length: 96
need:
  - length: 10
  - length: 20
  - length: 30
`
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}

	want := []string{"A", "B", "C"}
	for i, p := range cfg.Need {
		if p.Label != want[i] {
			t.Errorf("need[%d].Label: got %q, want %q", i, p.Label, want[i])
		}
	}
}

func TestParseFile_RotateDefault(t *testing.T) {
	yaml := "stock:\n  - length: 96\nneed:\n  - length: 10\n"
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Rotate {
		t.Error("rotate should default to true when not set in file")
	}
}

func TestParseFile_RotateExplicitFalse(t *testing.T) {
	yaml := "rotate: false\nstock:\n  - length: 96\nneed:\n  - length: 10\n"
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Rotate {
		t.Error("rotate should be false when explicitly set to false in file")
	}
}

func TestParseFile_OnHand(t *testing.T) {
	yaml := `
stock:
  - length: 48
    on_hand: true
  - length: 96
need:
  - length: 10
`
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}
	if !cfg.Stock[0].OnHand {
		t.Error("stock[0] should be on_hand")
	}
	if cfg.Stock[1].OnHand {
		t.Error("stock[1] should not be on_hand")
	}
}

func TestParseFile_CountZeroDefaultsToOne(t *testing.T) {
	yaml := "stock:\n  - length: 96\nneed:\n  - length: 10\n"
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Stock[0].Count != 1 {
		t.Errorf("stock count: got %d, want 1", cfg.Stock[0].Count)
	}
	if cfg.Need[0].Count != 1 {
		t.Errorf("need count: got %d, want 1", cfg.Need[0].Count)
	}
}

func TestParseFile_MissingFile(t *testing.T) {
	_, err := ParseFile(filepath.Join(t.TempDir(), "nonexistent.yaml"))
	if err == nil {
		t.Error("expected error for missing file")
	}
}

func TestParseFile_InvalidYAML(t *testing.T) {
	_, err := ParseFile(writeTemp(t, ":: invalid yaml ::"))
	if err == nil {
		t.Error("expected error for invalid YAML")
	}
}

func TestParseFile_RepeatFields(t *testing.T) {
	yaml := `
stock:
  - width: 27
    height: 240
    on_hand: true
    repeat_distance: 24
    repeat_axis: height
need:
  - width: 27
    height: 96
`
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}
	s := cfg.Stock[0]
	if s.RepeatDistance != 24 {
		t.Errorf("repeat_distance: got %v, want 24", s.RepeatDistance)
	}
	if s.RepeatAxis != "height" {
		t.Errorf("repeat_axis: got %q, want %q", s.RepeatAxis, "height")
	}
}

func TestParseFile_JoinGroup(t *testing.T) {
	yaml := `
stock:
  - length: 96
    count: 2
need:
  - length: 48
    join_group: panel-left
  - length: 48
    join_group: panel-left
  - length: 36
`
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.Need[0].JoinGroup != "panel-left" {
		t.Errorf("need[0].JoinGroup: got %q, want %q", cfg.Need[0].JoinGroup, "panel-left")
	}
	if cfg.Need[1].JoinGroup != "panel-left" {
		t.Errorf("need[1].JoinGroup: got %q, want %q", cfg.Need[1].JoinGroup, "panel-left")
	}
	if cfg.Need[2].JoinGroup != "" {
		t.Errorf("need[2].JoinGroup: got %q, want empty", cfg.Need[2].JoinGroup)
	}
}

func TestParseFile_OutputFormatDefault(t *testing.T) {
	yaml := "stock:\n  - length: 96\nneed:\n  - length: 10\n"
	cfg, err := ParseFile(writeTemp(t, yaml))
	if err != nil {
		t.Fatal(err)
	}
	if cfg.OutputFormat != "" {
		t.Errorf("output format: got %q, want empty (factory decides default)", cfg.OutputFormat)
	}
}
