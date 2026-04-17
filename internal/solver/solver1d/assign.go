package solver1d

import (
	"math"
	"sort"

	"github.com/beekman/cut-calculator/internal/model"
)

// assign distributes required pieces across stock using on-hand stock first.
// It returns a CutPlan and any pieces that could not be assigned.
func assign(stock []model.StockPiece, pieces []model.RequiredPiece, kerf float64) (model.CutPlan, []model.RequiredPiece) {
	// sort: on-hand first, then by length descending
	sort.Slice(stock, func(i, j int) bool {
		if stock[i].OnHand != stock[j].OnHand {
			return stock[i].OnHand
		}
		return stock[i].Length > stock[j].Length
	})

	// sort pieces longest first for better packing
	sort.Slice(pieces, func(i, j int) bool {
		return pieces[i].Length > pieces[j].Length
	})

	remaining := make([]model.RequiredPiece, len(pieces))
	copy(remaining, pieces)

	var results []model.StockResult
	var purchased []model.StockPiece
	totalUsed := 0.0
	totalRequired := 0.0

	for _, s := range stock {
		if len(remaining) == 0 {
			break
		}
		chosen := knapsack(s.Length, remaining, kerf, s.RepeatDistance)
		if len(chosen) == 0 {
			continue
		}

		result, used := buildResult(len(results), s, chosen, kerf, s.RepeatDistance)
		results = append(results, result)
		remaining = remove(remaining, chosen)

		totalUsed += s.Length
		totalRequired += used
		if !s.OnHand {
			purchased = append(purchased, s)
		}
	}

	var wastePct float64
	if totalUsed > 0 {
		wastePct = (totalUsed - totalRequired) / totalUsed * 100
	}

	return model.CutPlan{
		Results:   results,
		WastePct:  wastePct,
		Purchased: purchased,
	}, remaining
}

func buildResult(idx int, s model.StockPiece, chosen []model.RequiredPiece, kerf, repeatDist float64) (model.StockResult, float64) {
	var cuts []model.Cut
	var assignments []model.Assignment
	pos := 0.0
	usedLength := 0.0

	for i, p := range chosen {
		cuts = append(cuts, model.Cut{Position: pos, Label: p.Label})
		assignments = append(assignments, model.Assignment{
			StockIndex:    idx,
			RequiredLabel: p.Label,
			Length:        p.Length,
		})
		usedLength += p.Length
		pos += p.Length
		if i < len(chosen)-1 {
			pos += kerf
			if repeatDist > 0 {
				pos = math.Ceil(pos/repeatDist) * repeatDist
			}
		}
	}

	waste := s.Length - pos

	return model.StockResult{
		Stock:       s,
		Assignments: assignments,
		Cuts:        cuts,
		WasteLength: waste,
	}, usedLength
}

// remove returns pieces with the chosen subset removed (first occurrence of each).
func remove(pieces []model.RequiredPiece, chosen []model.RequiredPiece) []model.RequiredPiece {
	remaining := make([]model.RequiredPiece, len(pieces))
	copy(remaining, pieces)

	for _, c := range chosen {
		for i, r := range remaining {
			if r.Label == c.Label && r.Length == c.Length {
				remaining = append(remaining[:i], remaining[i+1:]...)
				break
			}
		}
	}
	return remaining
}
