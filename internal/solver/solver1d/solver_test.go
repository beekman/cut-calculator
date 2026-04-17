package solver1d_test

import (
	"testing"

	"github.com/beekman/cut-calculator/internal/model"
	"github.com/beekman/cut-calculator/internal/solver/solver1d"
)

func TestSolve_ReferenceCases(t *testing.T) {
	cases := []struct {
		name          string
		stock         []model.StockPiece
		need          []model.RequiredPiece
		kerf          float64
		wantUnfit     int
		wantPurchased int
		wantMaxWaste  float64 // upper bound on waste %, -1 to skip
	}{
		{
			name:         "exact fit two halves",
			stock:        []model.StockPiece{{Length: 96, Count: 1, OnHand: true}},
			need:         []model.RequiredPiece{{Label: "A", Length: 48, Count: 2}},
			kerf:         0,
			wantUnfit:    0,
			wantMaxWaste: 0,
		},
		{
			name:         "three thirds",
			stock:        []model.StockPiece{{Length: 96, Count: 1, OnHand: true}},
			need:         []model.RequiredPiece{{Label: "A", Length: 32, Count: 3}},
			kerf:         0,
			wantUnfit:    0,
			wantMaxWaste: 0,
		},
		{
			name:          "three thirds with kerf forces second board",
			stock:         []model.StockPiece{{Length: 96, Count: 1, OnHand: true}, {Length: 96, OnHand: false}},
			need:          []model.RequiredPiece{{Label: "A", Length: 32, Count: 3}},
			kerf:          0.125,
			wantUnfit:     0,
			wantPurchased: 1,
			wantMaxWaste:  -1, // ~50% is expected: 2 boards needed for 3 pieces
		},
		{
			name:         "on-hand used before purchasable",
			stock:        []model.StockPiece{{Length: 96, Count: 3, OnHand: true}, {Length: 96, Count: 10}},
			need:         []model.RequiredPiece{{Label: "A", Length: 36, Count: 4}, {Label: "B", Length: 48, Count: 2}},
			kerf:         0.125,
			wantUnfit:    0,
			wantPurchased: 0,
			wantMaxWaste: 25,
		},
		{
			name:         "offcut used optimally",
			stock:        []model.StockPiece{{Length: 96, Count: 1, OnHand: true}, {Length: 48, Count: 1, OnHand: true}},
			need:         []model.RequiredPiece{{Label: "A", Length: 48, Count: 2}},
			kerf:         0,
			wantUnfit:    0,
			wantPurchased: 0,
			wantMaxWaste: 0,
		},
		{
			name:        "piece exactly equals stock",
			stock:       []model.StockPiece{{Length: 48, Count: 1, OnHand: true}},
			need:        []model.RequiredPiece{{Label: "A", Length: 48, Count: 1}},
			kerf:        0.125,
			wantUnfit:   0,
			wantMaxWaste: 0,
		},
	}

	s := solver1d.New()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := s.Solve(tc.stock, tc.need, tc.kerf)
			if err != nil {
				t.Fatalf("Solve returned error: %v", err)
			}
			if len(plan.Unfit) != tc.wantUnfit {
				t.Errorf("unfit: got %d, want %d (pieces: %v)", len(plan.Unfit), tc.wantUnfit, plan.Unfit)
			}
			if len(plan.Purchased) != tc.wantPurchased {
				t.Errorf("purchased: got %d, want %d", len(plan.Purchased), tc.wantPurchased)
			}
			if tc.wantMaxWaste >= 0 && plan.WastePct > tc.wantMaxWaste {
				t.Errorf("waste: got %.2f%%, want ≤ %.2f%%", plan.WastePct, tc.wantMaxWaste)
			}
		})
	}
}

func TestSolve_Invariants(t *testing.T) {
	cases := []struct {
		name  string
		stock []model.StockPiece
		need  []model.RequiredPiece
		kerf  float64
	}{
		{
			name:  "standard job",
			stock: []model.StockPiece{{Length: 96, Count: 3, OnHand: true}, {Length: 96, Count: 5}},
			need:  []model.RequiredPiece{{Label: "A", Length: 36, Count: 4}, {Label: "B", Length: 48, Count: 2}},
			kerf:  0.125,
		},
		{
			name:  "many small pieces",
			stock: []model.StockPiece{{Length: 96, Count: 2, OnHand: true}},
			need:  []model.RequiredPiece{{Label: "A", Length: 12, Count: 8}},
			kerf:  0.125,
		},
	}

	s := solver1d.New()
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			plan, err := s.Solve(tc.stock, tc.need, tc.kerf)
			if err != nil {
				t.Fatalf("Solve: %v", err)
			}

			// Invariant: waste is non-negative and ≤ 100%
			if plan.WastePct < 0 || plan.WastePct > 100 {
				t.Errorf("waste out of range: %.2f%%", plan.WastePct)
			}

			// Invariant: no stock piece is over-allocated
			for _, r := range plan.Results {
				used := 0.0
				for _, a := range r.Assignments {
					used += a.Length
				}
				if len(r.Assignments) > 1 {
					used += float64(len(r.Assignments)-1) * tc.kerf
				}
				if used > r.Stock.Length+1e-9 {
					t.Errorf("stock over-allocated: used %.4f > stock %.4f", used, r.Stock.Length)
				}
			}

			// Invariant: on-hand stock exhausted before purchased stock is used
			onHandRemaining := false
			for _, r := range plan.Results {
				if r.Stock.OnHand && len(r.Assignments) == 0 {
					onHandRemaining = true
				}
			}
			if onHandRemaining && len(plan.Purchased) > 0 {
				t.Error("purchased stock used while on-hand stock has unused capacity")
			}
		})
	}
}

func TestSolve_EdgeCases(t *testing.T) {
	s := solver1d.New()

	t.Run("all pieces fit in on-hand stock", func(t *testing.T) {
		plan, _ := s.Solve(
			[]model.StockPiece{{Length: 96, Count: 3, OnHand: true}},
			[]model.RequiredPiece{{Label: "A", Length: 36, Count: 4}, {Label: "B", Length: 48, Count: 2}},
			0.125,
		)
		if len(plan.Purchased) != 0 {
			t.Errorf("expected no purchases, got %d", len(plan.Purchased))
		}
		if len(plan.Unfit) != 0 {
			t.Errorf("expected no unfit pieces, got %v", plan.Unfit)
		}
	})

	t.Run("duplicate required sizes are separate pieces", func(t *testing.T) {
		plan, _ := s.Solve(
			[]model.StockPiece{{Length: 96, Count: 2, OnHand: true}},
			[]model.RequiredPiece{{Label: "A", Length: 48, Count: 3}},
			0,
		)
		assigned := 0
		for _, r := range plan.Results {
			assigned += len(r.Assignments)
		}
		total := assigned + len(plan.Unfit)
		if total != 3 {
			t.Errorf("expected 3 total pieces accounted for, got %d", total)
		}
	})

	t.Run("stock smaller than all required pieces is skipped", func(t *testing.T) {
		plan, _ := s.Solve(
			[]model.StockPiece{
				{Length: 10, Count: 1, OnHand: true},
				{Length: 96, Count: 1, OnHand: true},
			},
			[]model.RequiredPiece{{Label: "A", Length: 48, Count: 1}},
			0,
		)
		if len(plan.Unfit) != 0 {
			t.Errorf("piece should be assigned to 96\" board, got unfit: %v", plan.Unfit)
		}
	})
}
