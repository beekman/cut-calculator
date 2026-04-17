package joiner

import (
	"testing"

	"github.com/beekman/cut-calculator/internal/model"
)

func TestHasJoinGroups(t *testing.T) {
	none := []model.RequiredPiece{{Label: "A", Length: 12}}
	if HasJoinGroups(none) {
		t.Error("expected false for pieces without join groups")
	}

	some := []model.RequiredPiece{
		{Label: "A", Length: 12},
		{Label: "B", Length: 6, JoinGroup: "g1"},
	}
	if !HasJoinGroups(some) {
		t.Error("expected true when at least one piece has a join group")
	}
}

func TestBuildCombinedNeed_1D(t *testing.T) {
	need := []model.RequiredPiece{
		{Label: "A", Length: 36, JoinGroup: "g1"},
		{Label: "B", Length: 24, JoinGroup: "g1"},
		{Label: "C", Length: 12},
	}
	out, meta := BuildCombinedNeed(need)

	if len(out) != 2 {
		t.Fatalf("expected 2 pieces, got %d", len(out))
	}
	// ungrouped C should be first
	if out[0].Label != "C" {
		t.Errorf("expected C first, got %q", out[0].Label)
	}
	combined := out[1]
	if combined.Label != "A+B" {
		t.Errorf("expected label A+B, got %q", combined.Label)
	}
	if combined.Length != 60 {
		t.Errorf("expected length 60, got %v", combined.Length)
	}
	m := meta["A+B"]
	if m.Axis != "length" {
		t.Errorf("expected axis 'length', got %q", m.Axis)
	}
	if len(m.Divisions) != 1 || m.Divisions[0] != 36 {
		t.Errorf("expected divisions [36], got %v", m.Divisions)
	}
	if len(m.Labels) != 2 || m.Labels[0] != "A" || m.Labels[1] != "B" {
		t.Errorf("expected labels [A B], got %v", m.Labels)
	}
}

func TestBuildCombinedNeed_2D_sameWidth(t *testing.T) {
	need := []model.RequiredPiece{
		{Label: "A", Width: 12, Height: 8, JoinGroup: "g1"},
		{Label: "B", Width: 12, Height: 6, JoinGroup: "g1"},
	}
	out, meta := BuildCombinedNeed(need)

	if len(out) != 1 {
		t.Fatalf("expected 1 piece, got %d", len(out))
	}
	c := out[0]
	if c.Width != 12 || c.Height != 14 {
		t.Errorf("expected 12×14, got %.4g×%.4g", c.Width, c.Height)
	}
	m := meta[c.Label]
	if m.Axis != "height" {
		t.Errorf("expected axis 'height', got %q", m.Axis)
	}
	if len(m.Divisions) != 1 || m.Divisions[0] != 8 {
		t.Errorf("expected divisions [8], got %v", m.Divisions)
	}
}

func TestBuildCombinedNeed_2D_sameHeight(t *testing.T) {
	need := []model.RequiredPiece{
		{Label: "A", Width: 10, Height: 8, JoinGroup: "g1"},
		{Label: "B", Width: 6, Height: 8, JoinGroup: "g1"},
	}
	out, meta := BuildCombinedNeed(need)

	c := out[0]
	if c.Width != 16 || c.Height != 8 {
		t.Errorf("expected 16×8, got %.4g×%.4g", c.Width, c.Height)
	}
	m := meta[c.Label]
	if m.Axis != "width" {
		t.Errorf("expected axis 'width', got %q", m.Axis)
	}
	if len(m.Divisions) != 1 || m.Divisions[0] != 10 {
		t.Errorf("expected divisions [10], got %v", m.Divisions)
	}
}

func TestBuildCombinedNeed_2D_incompatible(t *testing.T) {
	need := []model.RequiredPiece{
		{Label: "A", Width: 10, Height: 8, JoinGroup: "g1"},
		{Label: "B", Width: 6, Height: 5, JoinGroup: "g1"},
	}
	out, _ := BuildCombinedNeed(need)
	// incompatible dims → kept as individuals
	if len(out) != 2 {
		t.Errorf("expected 2 individual pieces, got %d", len(out))
	}
}

func TestBuildCombinedNeed_singleMember(t *testing.T) {
	need := []model.RequiredPiece{
		{Label: "A", Length: 36, JoinGroup: "g1"},
	}
	out, _ := BuildCombinedNeed(need)
	if len(out) != 1 || out[0].Label != "A" {
		t.Errorf("single-member group should be kept as-is, got %v", out)
	}
}

func TestBetterPlan_fewerUnfit(t *testing.T) {
	a := model.CutPlan{Unfit: []model.RequiredPiece{{}}, WastePct: 5}
	b := model.CutPlan{Unfit: nil, WastePct: 10}
	if got := BetterPlan(a, b); got.WastePct != 10 {
		t.Error("should prefer b (fewer unfit)")
	}
}

func TestBetterPlan_lowerWaste(t *testing.T) {
	a := model.CutPlan{WastePct: 5}
	b := model.CutPlan{WastePct: 10}
	if got := BetterPlan(a, b); got.WastePct != 5 {
		t.Error("should prefer a (lower waste)")
	}
}

func TestEnrichPlan(t *testing.T) {
	plan := &model.CutPlan{
		Results: []model.StockResult{
			{Assignments: []model.Assignment{
				{RequiredLabel: "A+B", Length: 60},
			}},
		},
	}
	meta := map[string]JoinMeta{
		"A+B": {Labels: []string{"A", "B"}, Divisions: []float64{36}, Axis: "length"},
	}
	EnrichPlan(plan, meta)
	a := plan.Results[0].Assignments[0]
	if a.JoinAxis != "length" || len(a.JoinLabels) != 2 || len(a.JoinDivisions) != 1 {
		t.Errorf("enrichment failed: %+v", a)
	}
}
