package output

import (
	"fmt"
	"io"
	"math"
	"strings"

	"github.com/beekman/cut-calculator/internal/model"
)

type ASCIIWriter struct{}

func (a *ASCIIWriter) Write(w io.Writer, plan model.CutPlan) error {
	if plan.Mode == 2 {
		return a.write2D(w, plan)
	}
	return a.write1D(w, plan)
}

// --- 1D ---

func (a *ASCIIWriter) write1D(w io.Writer, plan model.CutPlan) error {
	for i, r := range plan.Results {
		origin := "on hand"
		if !r.Stock.OnHand {
			origin = "purchased"
		}
		fmt.Fprintf(w, "Stock #%d (%.4g\") — %s\n", i+1, r.Stock.Length, origin)
		fmt.Fprintln(w, bar1D(r))
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "Overall waste: %.1f%%\n", plan.WastePct)

	if len(plan.Purchased) > 0 {
		fmt.Fprintln(w, "\nAdditional stock to purchase:")
		for _, p := range plan.Purchased {
			fmt.Fprintf(w, "  • %.4g\" board\n", p.Length)
		}
	} else if len(plan.Results) > 0 {
		fmt.Fprintln(w, "Nothing to purchase — all pieces fit within on-hand stock.")
	}

	if len(plan.Unfit) > 0 {
		fmt.Fprintln(w, "\nCould not fit:")
		for _, u := range plan.Unfit {
			fmt.Fprintf(w, "  • Piece %s (%.4g\")\n", u.Label, u.Length)
		}
	}
	return nil
}

// bar1D renders a stock result as a schematic cut bar.
// Example: |--A(36")--|--B(48")--|  [waste: 10"]
// With pattern repeat: |--A(36")--|~~(12")~~|--A(36")--|  [waste: 12"]
func bar1D(r model.StockResult) string {
	var sb strings.Builder
	sb.WriteRune('|')
	for i, a := range r.Assignments {
		// show alignment gap before this piece when repeat is active
		if i > 0 && r.Stock.RepeatDistance > 0 && i < len(r.Cuts) {
			prevEnd := r.Cuts[i-1].Position + r.Assignments[i-1].Length
			gap := r.Cuts[i].Position - prevEnd
			if gap > 0.001 {
				fmt.Fprintf(&sb, "~~(%.4g\")~~|", gap)
			}
		}
		fmt.Fprintf(&sb, "--%s(%.4g\")--|", a.RequiredLabel, a.Length)
	}
	if r.WasteLength > 0.001 {
		fmt.Fprintf(&sb, "  [waste: %.4g\"]", r.WasteLength)
	}
	return sb.String()
}

// --- 2D ---

const (
	maxGridW = 56
	maxGridH = 24
)

func (a *ASCIIWriter) write2D(w io.Writer, plan model.CutPlan) error {
	for i, r := range plan.Results {
		origin := "on hand"
		if !r.Stock.OnHand {
			origin = "purchased"
		}
		fmt.Fprintf(w, "Sheet #%d (%.4g\" × %.4g\") — %s\n",
			i+1, r.Stock.Width, r.Stock.Height, origin)

		grid, scale := newSheetGrid(r.Stock)
		renderRepeatLines(grid, r.Stock, scale)
		renderAssignments(grid, r.Assignments, scale)
		printGrid(w, grid)

		pct := 0.0
		if a := r.Stock.Width * r.Stock.Height; a > 0 {
			pct = r.WasteArea / a * 100
		}
		fmt.Fprintf(w, "waste: %.1f%%\n\n", pct)
	}

	fmt.Fprintf(w, "Overall waste: %.1f%%\n", plan.WastePct)

	if len(plan.Purchased) > 0 {
		fmt.Fprintln(w, "\nAdditional stock to purchase:")
		for _, p := range plan.Purchased {
			fmt.Fprintf(w, "  • %.4g\" × %.4g\" sheet\n", p.Width, p.Height)
		}
	} else if len(plan.Results) > 0 {
		fmt.Fprintln(w, "Nothing to purchase — all pieces fit within on-hand stock.")
	}

	if len(plan.Unfit) > 0 {
		fmt.Fprintln(w, "\nCould not fit:")
		for _, u := range plan.Unfit {
			fmt.Fprintf(w, "  • Piece %s (%.4g\" × %.4g\")\n", u.Label, u.Width, u.Height)
		}
	}
	return nil
}

// newSheetGrid allocates a character grid sized to fit the sheet, auto-scaled to terminal width.
func newSheetGrid(sheet model.StockPiece) ([][]rune, float64) {
	sx := float64(maxGridW) / sheet.Width
	sy := float64(maxGridH) / sheet.Height
	scale := sx
	if sy < sx {
		scale = sy
	}
	if scale > 2.0 {
		scale = 2.0
	}

	gw := int(math.Round(sheet.Width*scale)) + 1
	gh := int(math.Round(sheet.Height*scale)) + 1

	grid := make([][]rune, gh)
	for i := range grid {
		row := make([]rune, gw)
		for j := range row {
			row[j] = ' '
		}
		grid[i] = row
	}

	drawRect(grid, 0, 0, gw-1, gh-1)
	return grid, scale
}

func renderAssignments(grid [][]rune, assignments []model.Assignment, scale float64) {
	for _, a := range assignments {
		gx1 := int(math.Round(a.OffsetX * scale))
		gy1 := int(math.Round(a.OffsetY * scale))
		gx2 := int(math.Round((a.OffsetX + a.Width) * scale))
		gy2 := int(math.Round((a.OffsetY + a.Height) * scale))

		drawRect(grid, gx1, gy1, gx2, gy2)

		interior := gx2 - gx1 - 1
		if interior > 0 {
			if gy1+1 < gy2 {
				setStr(grid, gx1+1, gy1+1, a.RequiredLabel, interior)
			}
			if gy1+2 < gy2 {
				dims := fmt.Sprintf("%.4g×%.4g", a.Width, a.Height)
				setStr(grid, gx1+1, gy1+2, dims, interior)
			}
		}
	}
}

// renderRepeatLines draws boundary markers at each repeat interval before piece borders.
// Piece borders drawn afterwards naturally overwrite or merge at intersections.
func renderRepeatLines(grid [][]rune, sheet model.StockPiece, scale float64) {
	if sheet.RepeatDistance <= 0 {
		return
	}
	switch sheet.RepeatAxis {
	case "height":
		for y := sheet.RepeatDistance; y < sheet.Height; y += sheet.RepeatDistance {
			gy := int(math.Round(y * scale))
			for x := 0; x < len(grid[0]); x++ {
				setCell(grid, x, gy, '-')
			}
		}
	case "width":
		for x := sheet.RepeatDistance; x < sheet.Width; x += sheet.RepeatDistance {
			gx := int(math.Round(x * scale))
			for y := 0; y < len(grid); y++ {
				setCell(grid, gx, y, '|')
			}
		}
	}
}

// drawRect draws a rectangle border using +, -, | characters.
// Intersecting borders produce + corners automatically via setCell.
func drawRect(grid [][]rune, x1, y1, x2, y2 int) {
	for x := x1; x <= x2; x++ {
		setCell(grid, x, y1, '-')
		setCell(grid, x, y2, '-')
	}
	for y := y1; y <= y2; y++ {
		setCell(grid, x1, y, '|')
		setCell(grid, x2, y, '|')
	}
	// explicit corners in case earlier passes left stray chars
	setCell(grid, x1, y1, '+')
	setCell(grid, x1, y2, '+')
	setCell(grid, x2, y1, '+')
	setCell(grid, x2, y2, '+')
}

// setCell writes ch to (x,y) with priority rules:
//   + beats everything; | + - intersection becomes +; no other overwrites.
func setCell(grid [][]rune, x, y int, ch rune) {
	if y < 0 || y >= len(grid) || x < 0 || x >= len(grid[0]) {
		return
	}
	cur := grid[y][x]
	switch {
	case ch == '+':
		grid[y][x] = '+'
	case cur == '+':
		// + is never overwritten
	case cur == ' ':
		grid[y][x] = ch
	case (cur == '|' && ch == '-') || (cur == '-' && ch == '|'):
		grid[y][x] = '+'
	}
}

func setStr(grid [][]rune, x, y int, s string, maxW int) {
	runes := []rune(s)
	if len(runes) > maxW {
		runes = runes[:maxW]
	}
	for i, ch := range runes {
		if y >= 0 && y < len(grid) && x+i >= 0 && x+i < len(grid[0]) {
			grid[y][x+i] = ch
		}
	}
}

func printGrid(w io.Writer, grid [][]rune) {
	for _, row := range grid {
		fmt.Fprintln(w, string(row))
	}
}
