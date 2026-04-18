package solver1d

import "github.com/beekman/cut-calculator/internal/model"

type Solver struct{}

func New() *Solver { return &Solver{} }

func (s *Solver) Solve(stock []model.StockPiece, need []model.RequiredPiece, kerf float64) (model.CutPlan, error) {
	// expand counted stock and need into flat slices
	pieces := expandNeed(need)
	inventory := expandStock(stock)

	plan, unfit := assign(inventory, pieces, kerf)
	plan.Mode = 1
	plan.Unfit = unfit
	return plan, nil
}

func expandNeed(need []model.RequiredPiece) []model.RequiredPiece {
	var out []model.RequiredPiece
	for _, n := range need {
		for i := 0; i < n.Count; i++ {
			out = append(out, model.RequiredPiece{Label: n.Label, Length: n.Length, Count: 1})
		}
	}
	return out
}

func expandStock(stock []model.StockPiece) []model.StockPiece {
	var out []model.StockPiece
	for _, s := range stock {
		n := s.Count
		if n <= 0 {
			n = 1
		}
		for i := 0; i < n; i++ {
			out = append(out, model.StockPiece{
				Length:         s.Length,
				Count:          1,
				OnHand:         s.OnHand,
				RepeatDistance: s.RepeatDistance,
			})
		}
	}
	return out
}
