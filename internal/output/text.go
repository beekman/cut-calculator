package output

import (
	"fmt"
	"io"

	"github.com/beekman/cut-calculator/internal/model"
)

type TextWriter struct{}

func (t *TextWriter) Write(w io.Writer, plan model.CutPlan) error {
	is2D := plan.Mode == 2

	for i, r := range plan.Results {
		origin := "on hand"
		if !r.Stock.OnHand {
			origin = "purchased"
		}

		if is2D {
			fmt.Fprintf(w, "Sheet #%d — %.4g\" × %.4g\" (%s)\n", i+1, r.Stock.Width, r.Stock.Height, origin)
			for _, a := range r.Assignments {
				line := fmt.Sprintf("  Piece %s: %.4g\" × %.4g\" at (%.4g\", %.4g\")",
					a.RequiredLabel, a.Width, a.Height, a.OffsetX, a.OffsetY)
				if a.Rotated {
					line += " [rotated]"
				}
				fmt.Fprintln(w, line)
			}
			if r.WasteArea > 0 {
				fmt.Fprintf(w, "  Waste area: %.4g sq in\n", r.WasteArea)
			}
		} else {
			fmt.Fprintf(w, "Stock #%d — %.4g\" board (%s)\n", i+1, r.Stock.Length, origin)
			for j, cut := range r.Cuts {
				a := r.Assignments[j]
				fmt.Fprintf(w, "  Cut %d: %.4g\"  → Piece %s\n", j+1, a.Length, cut.Label)
			}
			if r.WasteLength > 0 {
				fmt.Fprintf(w, "  Offcut: %.4g\"\n", r.WasteLength)
			}
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "Overall waste: %.1f%%\n", plan.WastePct)

	if len(plan.Purchased) > 0 {
		fmt.Fprintln(w, "\nAdditional stock to purchase:")
		for _, p := range plan.Purchased {
			if is2D {
				fmt.Fprintf(w, "  • %.4g\" × %.4g\" sheet\n", p.Width, p.Height)
			} else {
				fmt.Fprintf(w, "  • %.4g\" board\n", p.Length)
			}
		}
	} else {
		fmt.Fprintln(w, "Nothing to purchase — all pieces fit within on-hand stock.")
	}

	if len(plan.Unfit) > 0 {
		fmt.Fprintln(w, "\nCould not fit:")
		for _, u := range plan.Unfit {
			if is2D {
				fmt.Fprintf(w, "  • Piece %s (%.4g\" × %.4g\")\n", u.Label, u.Width, u.Height)
			} else {
				fmt.Fprintf(w, "  • Piece %s (%.4g\")\n", u.Label, u.Length)
			}
		}
	}

	return nil
}
