package solver1d

import (
	"math"

	"github.com/beekman/cut-calculator/internal/model"
)

// knapsack finds the subset of pieces that maximizes used length within stockLength.
// Each piece costs (length + kerf); one extra kerf is added to capacity to cancel
// the unnecessary kerf charged to the first piece.
// Net constraint: sum(lengths) + (n-1)*kerf ≤ stockLength.
func knapsack(stockLength float64, pieces []model.RequiredPiece, kerf float64) []model.RequiredPiece {
	n := len(pieces)
	if n == 0 {
		return nil
	}

	const scale = 1000
	cap := int(math.Round((stockLength + kerf) * scale))

	dp := make([]int, cap+1)
	choice := make([][]bool, n)
	for i := range choice {
		choice[i] = make([]bool, cap+1)
	}

	weights := make([]int, n)
	for i, p := range pieces {
		weights[i] = int(math.Round((p.Length + kerf) * scale))
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
