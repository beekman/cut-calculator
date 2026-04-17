package output

import (
	"encoding/json"
	"io"

	"github.com/beekman/cut-calculator/internal/model"
)

type JSONWriter struct{}

func (j *JSONWriter) Write(w io.Writer, plan model.CutPlan) error {
	enc := json.NewEncoder(w)
	enc.SetIndent("", "  ")
	return enc.Encode(plan)
}
