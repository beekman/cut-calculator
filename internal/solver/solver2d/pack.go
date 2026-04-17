package solver2d

import (
	"math"

	"github.com/beekman/cut-calculator/internal/model"
)

// region is a free rectangular area on a sheet, with absolute coordinates.
type region struct{ x, y, w, h float64 }

func (r region) area() float64 { return r.w * r.h }
func (r region) valid() bool   { return r.w > 0 && r.h > 0 }

// placed records where a required piece was positioned on a sheet.
type placed struct {
	p       model.RequiredPiece
	x, y    float64
	w, h    float64 // actual dimensions used (may be swapped if rotated)
	rotated bool
}

type orientation struct {
	w, h    float64
	rotated bool
}

// packSheet greedily places as many pieces as possible into free regions.
// Pieces must already be sorted largest-area-first by the caller.
func packSheet(pieces []model.RequiredPiece, free []region, kerf float64, rotate bool, repeatDist float64, repeatAxis string) []placed {
	var placements []placed
	for _, p := range pieces {
		if pl, newFree := placeOne(p, free, kerf, rotate, repeatDist, repeatAxis); pl != nil {
			placements = append(placements, *pl)
			free = newFree
		}
	}
	return placements
}

// placeOne tries to fit piece p into the first valid free region.
// It tries both split orientations and keeps the one with the larger max free rect.
// When repeatDist > 0, each region is snapped to the next boundary on repeatAxis
// before fitting; sub-regions after splits are snapped likewise.
func placeOne(p model.RequiredPiece, free []region, kerf float64, rotate bool, repeatDist float64, repeatAxis string) (*placed, []region) {
	for i, fr := range free {
		snapped := snapRegion(fr, repeatDist, repeatAxis)
		if !snapped.valid() {
			continue
		}
		others := withoutIdx(free, i)
		for _, o := range orientations(p, rotate) {
			if o.w > snapped.w || o.h > snapped.h {
				continue
			}

			sub1H, sub2H := hSplit(snapped, o.w, o.h, kerf)
			sub1V, sub2V := vSplit(snapped, o.w, o.h, kerf)

			freesH := validOnly(snapRegions(concat(others, sub1H, sub2H), repeatDist, repeatAxis))
			freesV := validOnly(snapRegions(concat(others, sub1V, sub2V), repeatDist, repeatAxis))

			var chosen []region
			if maxArea(freesH) >= maxArea(freesV) {
				chosen = freesH
			} else {
				chosen = freesV
			}
			return &placed{p, snapped.x, snapped.y, o.w, o.h, o.rotated}, chosen
		}
	}
	return nil, nil
}

// snapRegion advances the starting coordinate of fr on repeatAxis to the next
// repeat boundary, shrinking the available dimension accordingly.
// A zero repeatDist is a no-op.
func snapRegion(fr region, repeatDist float64, repeatAxis string) region {
	if repeatDist <= 0 {
		return fr
	}
	switch repeatAxis {
	case "height":
		snappedY := math.Ceil(fr.y/repeatDist) * repeatDist
		gap := snappedY - fr.y
		return region{x: fr.x, y: snappedY, w: fr.w, h: fr.h - gap}
	case "width":
		snappedX := math.Ceil(fr.x/repeatDist) * repeatDist
		gap := snappedX - fr.x
		return region{x: snappedX, y: fr.y, w: fr.w - gap, h: fr.h}
	}
	return fr
}

func snapRegions(rects []region, repeatDist float64, repeatAxis string) []region {
	if repeatDist <= 0 {
		return rects
	}
	out := make([]region, len(rects))
	for i, r := range rects {
		out[i] = snapRegion(r, repeatDist, repeatAxis)
	}
	return out
}

// hSplit splits region fr after placing a piece (pw×ph) at its top-left.
// Horizontal cut runs across full width at y = fr.y + ph + kerf.
//   FR1: to the right of the piece, same row height
//   FR2: below the cut, full sheet width
func hSplit(fr region, pw, ph, kerf float64) (region, region) {
	fr1 := region{x: fr.x + pw + kerf, y: fr.y, w: fr.w - pw - kerf, h: ph}
	fr2 := region{x: fr.x, y: fr.y + ph + kerf, w: fr.w, h: fr.h - ph - kerf}
	return fr1, fr2
}

// vSplit splits region fr after placing a piece (pw×ph) at its top-left.
// Vertical cut runs across full height at x = fr.x + pw + kerf.
//   FR1: to the right of the piece, full column height
//   FR2: below the piece, same column width
func vSplit(fr region, pw, ph, kerf float64) (region, region) {
	fr1 := region{x: fr.x + pw + kerf, y: fr.y, w: fr.w - pw - kerf, h: fr.h}
	fr2 := region{x: fr.x, y: fr.y + ph + kerf, w: pw, h: fr.h - ph - kerf}
	return fr1, fr2
}

func orientations(p model.RequiredPiece, rotate bool) []orientation {
	base := orientation{p.Width, p.Height, false}
	if !rotate || p.Width == p.Height {
		return []orientation{base}
	}
	return []orientation{base, {p.Height, p.Width, true}}
}

func withoutIdx(rects []region, idx int) []region {
	out := make([]region, 0, len(rects)-1)
	out = append(out, rects[:idx]...)
	return append(out, rects[idx+1:]...)
}

func concat(a []region, extra ...region) []region {
	out := make([]region, len(a), len(a)+len(extra))
	copy(out, a)
	return append(out, extra...)
}

func validOnly(rects []region) []region {
	var out []region
	for _, r := range rects {
		if r.valid() {
			out = append(out, r)
		}
	}
	return out
}

func maxArea(rects []region) float64 {
	var m float64
	for _, r := range rects {
		if a := r.area(); a > m {
			m = a
		}
	}
	return m
}
