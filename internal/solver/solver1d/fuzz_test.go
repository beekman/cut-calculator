package solver1d_test

import (
	"testing"

	"github.com/beekman/cut-calculator/internal/model"
	"github.com/beekman/cut-calculator/internal/solver/solver1d"
)

// FuzzSolve1D tests solver invariants against arbitrary inputs.
// Invariants checked per solve call:
//   - assigned + unfit == total expanded need (no silent drops)
//   - each stock result is not over-allocated
//   - WastePct is non-negative
//   - if any on-hand stock is unused, Purchased is empty
func FuzzSolve1D(f *testing.F) {
	// seed corpus covering common cases
	f.Add(96.0, 0.125, 36.0, 2, 22.0, 1, true)
	f.Add(96.0, 0.0, 48.0, 2, 0.0, 0, true)
	f.Add(48.0, 0.0625, 24.0, 3, 12.0, 2, false)
	f.Add(120.0, 0.0, 60.0, 1, 30.0, 4, true)
	f.Add(96.0, 0.125, 97.0, 1, 0.0, 0, true) // piece larger than stock

	f.Fuzz(func(t *testing.T, stockLen, kerf, needLen1 float64, needCount1 int, needLen2 float64, needCount2 int, onHand bool) {
		// skip degenerate values — these are not valid inputs
		if stockLen <= 0 || stockLen > 500 {
			return
		}
		if kerf < 0 || kerf >= stockLen {
			return
		}
		if needLen1 <= 0 || needLen1 > 500 {
			return
		}
		if needCount1 < 1 || needCount1 > 20 {
			return
		}

		stock := []model.StockPiece{{Length: stockLen, Count: 3, OnHand: onHand}}
		need := []model.RequiredPiece{{Label: "A", Length: needLen1, Count: needCount1}}

		// optionally add a second need piece type
		if needLen2 > 0 && needLen2 <= 500 && needCount2 >= 1 && needCount2 <= 20 {
			need = append(need, model.RequiredPiece{Label: "B", Length: needLen2, Count: needCount2})
		}

		plan, _ := solver1d.New().Solve(stock, need, kerf)

		checkInvariants1D(t, stock, need, plan, kerf)
	})
}

func checkInvariants1D(t *testing.T, stock []model.StockPiece, need []model.RequiredPiece, plan model.CutPlan, kerf float64) {
	t.Helper()

	// invariant: assigned + unfit == total expanded need
	totalNeed := expandedCount(need)
	assigned := 0
	for _, r := range plan.Results {
		assigned += len(r.Assignments)
	}
	if assigned+len(plan.Unfit) != totalNeed {
		t.Errorf("piece accounting: assigned(%d) + unfit(%d) = %d, want %d",
			assigned, len(plan.Unfit), assigned+len(plan.Unfit), totalNeed)
	}

	// invariant: no stock result over-allocated
	for i, r := range plan.Results {
		used := 0.0
		for j, a := range r.Assignments {
			used += a.Length
			if j > 0 {
				used += kerf
			}
		}
		if used > r.Stock.Length+1e-9 {
			t.Errorf("result[%d] over-allocated: used %.6f > stock %.6f", i, used, r.Stock.Length)
		}
	}

	// invariant: WastePct non-negative
	if plan.WastePct < -1e-9 {
		t.Errorf("WastePct negative: %.6f", plan.WastePct)
	}

	// invariant: on-hand stock used before purchased (if all on-hand, Purchased empty)
	allOnHand := true
	for _, s := range stock {
		if !s.OnHand {
			allOnHand = false
			break
		}
	}
	if allOnHand && len(plan.Purchased) > 0 && len(plan.Unfit) == 0 {
		t.Errorf("on-hand stock available but Purchased is non-empty")
	}
}

func expandedCount(need []model.RequiredPiece) int {
	n := 0
	for _, p := range need {
		c := p.Count
		if c <= 0 {
			c = 1
		}
		n += c
	}
	return n
}
