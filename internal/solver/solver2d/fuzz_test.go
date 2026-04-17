package solver2d

import (
	"testing"

	"github.com/beekman/cut-calculator/internal/model"
)

// FuzzSolve2D tests 2D solver invariants against arbitrary inputs.
// Invariants checked per solve call:
//   - assigned + unfit == total expanded need (no silent drops)
//   - no two assignments on the same sheet overlap
//   - all assignments fit within their sheet bounds
//   - WastePct is non-negative
func FuzzSolve2D(f *testing.F) {
	// seed corpus
	f.Add(48.0, 96.0, 0.0, 24.0, 48.0, 4, true, true)
	f.Add(48.0, 96.0, 0.125, 12.0, 12.0, 6, true, false)
	f.Add(24.0, 24.0, 0.0, 12.0, 12.0, 4, true, true)
	f.Add(48.0, 96.0, 0.0, 30.0, 48.0, 2, true, true) // rotation needed

	f.Fuzz(func(t *testing.T, sheetW, sheetH, kerf, pieceW, pieceH float64, count int, onHand, rotate bool) {
		if sheetW <= 0 || sheetW > 500 {
			return
		}
		if sheetH <= 0 || sheetH > 500 {
			return
		}
		if kerf < 0 || kerf >= sheetW || kerf >= sheetH {
			return
		}
		if pieceW <= 0 || pieceW > 500 {
			return
		}
		if pieceH <= 0 || pieceH > 500 {
			return
		}
		if count < 1 || count > 20 {
			return
		}

		stock := []model.StockPiece{{Width: sheetW, Height: sheetH, Count: 3, OnHand: onHand}}
		need := []model.RequiredPiece{{Label: "A", Width: pieceW, Height: pieceH, Count: count}}

		plan, _ := New().Solve(stock, need, kerf, rotate)

		checkInvariants2D(t, need, plan)
	})
}

func checkInvariants2D(t *testing.T, need []model.RequiredPiece, plan model.CutPlan) {
	t.Helper()

	// invariant: assigned + unfit == total expanded need
	totalNeed := 0
	for _, p := range need {
		c := p.Count
		if c <= 0 {
			c = 1
		}
		totalNeed += c
	}
	assigned := 0
	for _, r := range plan.Results {
		assigned += len(r.Assignments)
	}
	if assigned+len(plan.Unfit) != totalNeed {
		t.Errorf("piece accounting: assigned(%d) + unfit(%d) = %d, want %d",
			assigned, len(plan.Unfit), assigned+len(plan.Unfit), totalNeed)
	}

	// invariant: no overlap and all in bounds per sheet
	for i, r := range plan.Results {
		sw := r.Stock.Width
		sh := r.Stock.Height

		for j, a := range r.Assignments {
			// bounds check
			if a.OffsetX < -1e-9 || a.OffsetY < -1e-9 ||
				a.OffsetX+a.Width > sw+1e-9 || a.OffsetY+a.Height > sh+1e-9 {
				t.Errorf("result[%d] assignment[%d] out of bounds: (%.4f,%.4f)+(%.4f×%.4f) in %.4f×%.4f",
					i, j, a.OffsetX, a.OffsetY, a.Width, a.Height, sw, sh)
			}

			// overlap check against all later assignments on same sheet
			for k := j + 1; k < len(r.Assignments); k++ {
				b := r.Assignments[k]
				if overlaps(a, b) {
					t.Errorf("result[%d]: assignments %d and %d overlap", i, j, k)
				}
			}
		}
	}

	// invariant: WastePct non-negative
	if plan.WastePct < -1e-9 {
		t.Errorf("WastePct negative: %.6f", plan.WastePct)
	}
}

func overlaps(a, b model.Assignment) bool {
	return a.OffsetX < b.OffsetX+b.Width-1e-9 &&
		a.OffsetX+a.Width > b.OffsetX+1e-9 &&
		a.OffsetY < b.OffsetY+b.Height-1e-9 &&
		a.OffsetY+a.Height > b.OffsetY+1e-9
}
