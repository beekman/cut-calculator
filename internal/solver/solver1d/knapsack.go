package solver1d

import (
	"math"

	"github.com/beekman/cut-calculator/internal/model"
)

// knapsack finds the subset of pieces that maximizes used length within stockLength.
//
// Without repeat: each piece costs (length + kerf); capacity gets +kerf to cancel the
// kerf charged to the last piece. Net: sum(lengths) + (n-1)*kerf ≤ stockLength.
//
// With repeat: each piece costs ceil((length+kerf)/repeatDist)*repeatDist — the full
// repeat cell(s) it consumes. Capacity is stockLength with no adjustment (conservative:
// may leave a small gap at end, never over-commits).
func knapsack(stockLength float64, pieces []model.RequiredPiece, kerf, repeatDist float64) []model.RequiredPiece {
	n := len(pieces)
	if n == 0 {
		return nil
	}

	const scale = 1000
	var cap int
	if repeatDist > 0 {
		cap = int(math.Round(stockLength * scale))
	} else {
		cap = int(math.Round((stockLength + kerf) * scale))
	}

	dp := make([]int, cap+1)
	choice := make([][]bool, n)
	for i := range choice {
		choice[i] = make([]bool, cap+1)
	}

	weights := make([]int, n)
	for i, p := range pieces {
		if repeatDist > 0 {
			cells := math.Ceil((p.Length + kerf) / repeatDist)
			weights[i] = int(math.Round(cells * repeatDist * scale))
		} else {
			weights[i] = int(math.Round((p.Length + kerf) * scale))
		}
	}

	for i := range pieces {
		w := weights[i]
		for c := cap; c >= w; c-- {
			if dp[c-w]+w > dp[c] {
				dp[c] = dp[c-w] + w
				choice[i][c] = true
			}
		}
	}

	var chosen []model.RequiredPiece
	remaining := cap
	for i := n - 1; i >= 0; i-- {
		if choice[i][remaining] {
			chosen = append(chosen, pieces[i])
			remaining -= weights[i]
		}
	}

	return chosen
}
