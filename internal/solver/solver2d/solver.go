package solver2d

import (
	"sort"

	"github.com/beekman/cut-calculator/internal/model"
)

type Solver struct{}

func New() *Solver { return &Solver{} }

func (s *Solver) Solve(stock []model.StockPiece, need []model.RequiredPiece, kerf float64, rotate bool) (model.CutPlan, error) {
	sheets := expandStock(stock)
	pieces := expandNeed(need)

	// on-hand first, then largest area
	sort.Slice(sheets, func(i, j int) bool {
		if sheets[i].OnHand != sheets[j].OnHand {
			return sheets[i].OnHand
		}
		return sheetArea(sheets[i]) > sheetArea(sheets[j])
	})

	// largest area first for better packing
	sort.Slice(pieces, func(i, j int) bool {
		return pieces[i].Width*pieces[i].Height > pieces[j].Width*pieces[j].Height
	})

	remaining := make([]model.RequiredPiece, len(pieces))
	copy(remaining, pieces)

	var results []model.StockResult
	var purchased []model.StockPiece
	var totalSheetArea, totalPlacedArea float64

	for _, sheet := range sheets {
		if len(remaining) == 0 {
			break
		}

		fr := region{0, 0, sheet.Width, sheet.Height}
		placements := packSheet(remaining, []region{fr}, kerf, rotate)
		if len(placements) == 0 {
			continue
		}

		result, placedArea := buildResult(len(results), sheet, placements)
		results = append(results, result)

		remaining = removeByLabel(remaining, placements)
		totalSheetArea += sheetArea(sheet)
		totalPlacedArea += placedArea
		if !sheet.OnHand {
			purchased = append(purchased, sheet)
		}
	}

	var wastePct float64
	if totalSheetArea > 0 {
		wastePct = (totalSheetArea - totalPlacedArea) / totalSheetArea * 100
	}

	return model.CutPlan{
		Mode:      2,
		Results:   results,
		WastePct:  wastePct,
		Purchased: purchased,
		Unfit:     remaining,
	}, nil
}

func buildResult(idx int, sheet model.StockPiece, placements []placed) (model.StockResult, float64) {
	var assignments []model.Assignment
	var placedArea float64

	for _, pl := range placements {
		assignments = append(assignments, model.Assignment{
			StockIndex:    idx,
			RequiredLabel: pl.p.Label,
			Width:         pl.w,
			Height:        pl.h,
			OffsetX:       pl.x,
			OffsetY:       pl.y,
			Rotated:       pl.rotated,
		})
		placedArea += pl.w * pl.h
	}

	wasteArea := sheetArea(sheet) - placedArea
	return model.StockResult{
		Stock:       sheet,
		Assignments: assignments,
		WasteArea:   wasteArea,
	}, placedArea
}

func expandStock(stock []model.StockPiece) []model.StockPiece {
	var out []model.StockPiece
	for _, s := range stock {
		n := s.Count
		if n <= 0 {
			n = 1
		}
		for i := 0; i < n; i++ {
			out = append(out, model.StockPiece{Width: s.Width, Height: s.Height, Count: 1, OnHand: s.OnHand})
		}
	}
	return out
}

func expandNeed(need []model.RequiredPiece) []model.RequiredPiece {
	var out []model.RequiredPiece
	for _, n := range need {
		for i := 0; i < n.Count; i++ {
			out = append(out, model.RequiredPiece{Label: n.Label, Width: n.Width, Height: n.Height, Count: 1})
		}
	}
	return out
}

func sheetArea(s model.StockPiece) float64 { return s.Width * s.Height }

// removeByLabel removes one occurrence of each placed piece from remaining.
func removeByLabel(remaining []model.RequiredPiece, placements []placed) []model.RequiredPiece {
	out := make([]model.RequiredPiece, len(remaining))
	copy(out, remaining)
	for _, pl := range placements {
		for i, r := range out {
			if r.Label == pl.p.Label && r.Width == pl.p.Width && r.Height == pl.p.Height {
				out = append(out[:i], out[i+1:]...)
				break
			}
		}
	}
	return out
}
