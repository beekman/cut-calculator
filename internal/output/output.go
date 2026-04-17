package output

import (
	"io"

	"github.com/beekman/cut-calculator/internal/model"
)

type Writer interface {
	Write(w io.Writer, plan model.CutPlan) error
}

func New(format model.OutputFormat) Writer {
	switch format {
	case model.OutputJSON:
		return &JSONWriter{}
	default:
		return &TextWriter{}
	}
}
