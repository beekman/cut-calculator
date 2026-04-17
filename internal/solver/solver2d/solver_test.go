package solver2d

import (
	"testing"

	"github.com/beekman/cut-calculator/internal/model"
)

func sheet(w, h float64, onHand bool) model.StockPiece {
	return model.StockPiece{Width: w, Height: h, Count: 1, OnHand: onHand}
}

func piece(label string, w, h float64, count int) model.RequiredPiece {
	return model.RequiredPiece{Label: label, Width: w, Height: h, Count: count}
}

// --- reference cases ---

func TestSolve_SinglePieceFitsExactly(t *testing.T) {
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{piece("A", 48, 96, 1)},
		0, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Unfit) != 0 {
		t.Errorf("expected 0 unfit, got %d", len(plan.Unfit))
	}
	if len(plan.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(plan.Results))
	}
	if plan.Results[0].WasteArea != 0 {
		t.Errorf("expected 0 waste, got %.2f", plan.Results[0].WasteArea)
	}
}

func TestSolve_FourEqualQuarters(t *testing.T) {
	// 48×96 sheet → four 24×48 pieces (exact fit, no kerf)
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{piece("A", 24, 48, 4)},
		0, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Unfit) != 0 {
		t.Errorf("expected 0 unfit, got %d: %v", len(plan.Unfit), plan.Unfit)
	}
	if len(plan.Results[0].Assignments) != 4 {
		t.Errorf("expected 4 assignments, got %d", len(plan.Results[0].Assignments))
	}
}

func TestSolve_KerfReducesUsableArea(t *testing.T) {
	// 48×96 with 0.125" kerf; two 48×47.9375 pieces should fit but two 48×48 shouldn't
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{piece("A", 48, 47.9375, 2)},
		0.125, false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Unfit) != 0 {
		t.Errorf("expected 0 unfit with kerf, got %d", len(plan.Unfit))
	}
}

func TestSolve_PieceTooLarge(t *testing.T) {
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{piece("A", 60, 96, 1)},
		0, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Unfit) != 1 {
		t.Errorf("expected 1 unfit, got %d", len(plan.Unfit))
	}
}

func TestSolve_OnHandUsedFirst(t *testing.T) {
	plan, err := New().Solve(
		[]model.StockPiece{
			sheet(48, 96, false), // purchasable
			sheet(48, 96, true),  // on hand
		},
		[]model.RequiredPiece{piece("A", 24, 48, 1)},
		0, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Results) != 1 {
		t.Fatalf("expected 1 result, got %d", len(plan.Results))
	}
	if !plan.Results[0].Stock.OnHand {
		t.Error("expected on-hand sheet to be used first")
	}
	if len(plan.Purchased) != 0 {
		t.Errorf("expected 0 purchased, got %d", len(plan.Purchased))
	}
}

func TestSolve_SpillsToSecondSheet(t *testing.T) {
	// two 48×96 sheets; need more than one sheet can hold
	plan, err := New().Solve(
		[]model.StockPiece{{Width: 48, Height: 96, Count: 2, OnHand: true}},
		[]model.RequiredPiece{piece("A", 48, 48, 3)}, // 3 pieces, each half a sheet
		0, false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Unfit) != 0 {
		t.Errorf("expected 0 unfit, got %d", len(plan.Unfit))
	}
	if len(plan.Results) != 2 {
		t.Errorf("expected 2 sheets used, got %d", len(plan.Results))
	}
}

// --- rotation ---

func TestSolve_RotationAllowsFit(t *testing.T) {
	// piece is 12×48; sheet is 48×12; only fits if rotated
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 12, true)},
		[]model.RequiredPiece{piece("A", 12, 48, 1)},
		0, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Unfit) != 0 {
		t.Errorf("expected 0 unfit with rotation, got %d", len(plan.Unfit))
	}
	if !plan.Results[0].Assignments[0].Rotated {
		t.Error("expected piece to be marked as rotated")
	}
}

func TestSolve_NoRotateRespected(t *testing.T) {
	// same as above but rotation disabled — piece should not fit
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 12, true)},
		[]model.RequiredPiece{piece("A", 12, 48, 1)},
		0, false,
	)
	if err != nil {
		t.Fatal(err)
	}
	if len(plan.Unfit) != 1 {
		t.Errorf("expected 1 unfit when rotation disabled, got %d", len(plan.Unfit))
	}
}

// --- invariants ---

func TestSolve_NoOverlap(t *testing.T) {
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{
			piece("A", 24, 48, 2),
			piece("B", 12, 24, 4),
		},
		0, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	for _, result := range plan.Results {
		checkNoOverlap(t, result.Assignments)
	}
}

func TestSolve_AllAssignmentsInBounds(t *testing.T) {
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{
			piece("A", 24, 48, 2),
			piece("B", 10, 20, 3),
		},
		0.125, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	for ri, result := range plan.Results {
		sw := result.Stock.Width
		sh := result.Stock.Height
		for ai, a := range result.Assignments {
			if a.OffsetX < 0 || a.OffsetY < 0 ||
				a.OffsetX+a.Width > sw || a.OffsetY+a.Height > sh {
				t.Errorf("result[%d] assignment[%d]: out of bounds (%.2f,%.2f)+(%.2f×%.2f) in %.2f×%.2f",
					ri, ai, a.OffsetX, a.OffsetY, a.Width, a.Height, sw, sh)
			}
		}
	}
}

func TestSolve_WasteNonNegative(t *testing.T) {
	plan, err := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{piece("A", 12, 12, 6)},
		0.125, true,
	)
	if err != nil {
		t.Fatal(err)
	}
	if plan.WastePct < 0 {
		t.Errorf("waste pct should be non-negative, got %.2f", plan.WastePct)
	}
	for _, result := range plan.Results {
		if result.WasteArea < -0.001 {
			t.Errorf("WasteArea should be non-negative, got %.4f", result.WasteArea)
		}
	}
}

func TestSolve_ModeFlagIsTwo(t *testing.T) {
	plan, _ := New().Solve(
		[]model.StockPiece{sheet(48, 96, true)},
		[]model.RequiredPiece{piece("A", 12, 24, 1)},
		0, true,
	)
	if plan.Mode != 2 {
		t.Errorf("expected Mode=2, got %d", plan.Mode)
	}
}

// checkNoOverlap verifies no two assignments on the same sheet overlap.
func checkNoOverlap(t *testing.T, assignments []model.Assignment) {
	t.Helper()
	for i := 0; i < len(assignments); i++ {
		for j := i + 1; j < len(assignments); j++ {
			a, b := assignments[i], assignments[j]
			if rectsOverlap(a.OffsetX, a.OffsetY, a.Width, a.Height,
				b.OffsetX, b.OffsetY, b.Width, b.Height) {
				t.Errorf("assignments %d and %d overlap: (%.2f,%.2f)+(%.2f×%.2f) vs (%.2f,%.2f)+(%.2f×%.2f)",
					i, j, a.OffsetX, a.OffsetY, a.Width, a.Height,
					b.OffsetX, b.OffsetY, b.Width, b.Height)
			}
		}
	}
}

func rectsOverlap(x1, y1, w1, h1, x2, y2, w2, h2 float64) bool {
	return x1 < x2+w2 && x1+w1 > x2 &&
		y1 < y2+h2 && y1+h1 > y2
}
