package output

import (
	"fmt"
	"io"

	"github.com/beekman/cut-calculator/internal/model"
)

type TextWriter struct{}

func (t *TextWriter) Write(w io.Writer, plan model.CutPlan) error {
	for i, r := range plan.Results {
		origin := "on hand"
		if !r.Stock.OnHand {
			origin = "purchased"
		}
		fmt.Fprintf(w, "Stock #%d — %.4g\" board (%s)\n", i+1, r.Stock.Length, origin)

		for j, cut := range r.Cuts {
			a := r.Assignments[j]
			fmt.Fprintf(w, "  Cut %d: %.4g\"  → Piece %s\n", j+1, a.Length, cut.Label)
		}

		if r.WasteLength > 0 {
			fmt.Fprintf(w, "  Offcut: %.4g\"\n", r.WasteLength)
		}
		fmt.Fprintln(w)
	}

	fmt.Fprintf(w, "Overall waste: %.1f%%\n", plan.WastePct)

	if len(plan.Purchased) > 0 {
		fmt.Fprintln(w, "\nAdditional stock to purchase:")
		for _, p := range plan.Purchased {
			fmt.Fprintf(w, "  • %.4g\" board\n", p.Length)
		}
	} else {
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
