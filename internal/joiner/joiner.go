package joiner

import (
	"strings"

	"github.com/beekman/cut-calculator/internal/model"
)

// JoinMeta holds rendering info for a combined piece.
type JoinMeta struct {
	Labels    []string
	Divisions []float64
	Axis      string // "length" (1D), "width" or "height" (2D)
}

// HasJoinGroups reports whether any piece has a JoinGroup set.
func HasJoinGroups(need []model.RequiredPiece) bool {
	for _, p := range need {
		if p.JoinGroup != "" {
			return true
		}
	}
	return false
}

// BuildCombinedNeed replaces each join group with a single combined piece.
// Returns the modified need list and metadata for enriching assignments later.
// Groups with incompatible dimensions are kept as individual pieces.
func BuildCombinedNeed(need []model.RequiredPiece) ([]model.RequiredPiece, map[string]JoinMeta) {
	groups := map[string][]model.RequiredPiece{}
	var groupOrder []string
	var ungrouped []model.RequiredPiece

	for _, p := range need {
		if p.JoinGroup == "" {
			ungrouped = append(ungrouped, p)
		} else {
			if _, exists := groups[p.JoinGroup]; !exists {
				groupOrder = append(groupOrder, p.JoinGroup)
			}
			groups[p.JoinGroup] = append(groups[p.JoinGroup], p)
		}
	}

	meta := map[string]JoinMeta{}
	out := append([]model.RequiredPiece{}, ungrouped...)

	for _, gname := range groupOrder {
		members := groups[gname]
		combined, m, ok := combineGroup(members)
		if !ok {
			out = append(out, members...)
			continue
		}
		out = append(out, combined)
		meta[combined.Label] = m
	}

	return out, meta
}

// EnrichPlan annotates assignments that correspond to combined pieces with
// JoinLabels, JoinDivisions, and JoinAxis.
func EnrichPlan(plan *model.CutPlan, meta map[string]JoinMeta) {
	for i := range plan.Results {
		for j := range plan.Results[i].Assignments {
			a := &plan.Results[i].Assignments[j]
			if m, ok := meta[a.RequiredLabel]; ok {
				a.JoinLabels = m.Labels
				a.JoinDivisions = m.Divisions
				a.JoinAxis = m.Axis
			}
		}
	}
}

// BetterPlan returns whichever plan has fewer unfit pieces; tie-breaks on lower waste.
func BetterPlan(a, b model.CutPlan) model.CutPlan {
	if len(a.Unfit) < len(b.Unfit) {
		return a
	}
	if len(b.Unfit) < len(a.Unfit) {
		return b
	}
	if a.WastePct <= b.WastePct {
		return a
	}
	return b
}

// combineGroup merges members of a join group into one combined piece.
// Returns false if combining is not possible (single member, or incompatible 2D dims).
func combineGroup(members []model.RequiredPiece) (model.RequiredPiece, JoinMeta, bool) {
	if len(members) < 2 {
		return model.RequiredPiece{}, JoinMeta{}, false
	}

	labels := make([]string, len(members))
	for i, m := range members {
		labels[i] = m.Label
	}
	combinedLabel := strings.Join(labels, "+")

	// 1D case: no width or height set
	if members[0].Width == 0 && members[0].Height == 0 {
		pos := 0.0
		divisions := make([]float64, len(members)-1)
		for i, m := range members {
			pos += m.Length
			if i < len(members)-1 {
				divisions[i] = pos
			}
		}
		return model.RequiredPiece{
			Label:  combinedLabel,
			Length: pos,
			Count:  1,
		}, JoinMeta{Labels: labels, Divisions: divisions, Axis: "length"}, true
	}

	// 2D: same width → join along height axis (stack vertically)
	sameW, sameH := true, true
	for _, m := range members[1:] {
		if m.Width != members[0].Width {
			sameW = false
		}
		if m.Height != members[0].Height {
			sameH = false
		}
	}

	if sameW {
		totalH := 0.0
		divisions := make([]float64, len(members)-1)
		for i, m := range members {
			totalH += m.Height
			if i < len(members)-1 {
				divisions[i] = totalH
			}
		}
		return model.RequiredPiece{
			Label:  combinedLabel,
			Width:  members[0].Width,
			Height: totalH,
			Count:  1,
		}, JoinMeta{Labels: labels, Divisions: divisions, Axis: "height"}, true
	}

	// 2D: same height → join along width axis (side by side)
	if sameH {
		totalW := 0.0
		divisions := make([]float64, len(members)-1)
		for i, m := range members {
			totalW += m.Width
			if i < len(members)-1 {
				divisions[i] = totalW
			}
		}
		return model.RequiredPiece{
			Label:  combinedLabel,
			Width:  totalW,
			Height: members[0].Height,
			Count:  1,
		}, JoinMeta{Labels: labels, Divisions: divisions, Axis: "width"}, true
	}

	return model.RequiredPiece{}, JoinMeta{}, false
}
