package model

type StockPiece struct {
	Length float64
	Width  float64
	Height float64
	Count  int
	OnHand bool
}

type RequiredPiece struct {
	Label  string
	Length float64
	Width  float64
	Height float64
	Count  int
}

type Cut struct {
	Position float64
	Label    string
	Axis     string // 2D only: "x" or "y"
}

type Assignment struct {
	StockIndex    int
	RequiredLabel string
	Length        float64
	Width         float64
	Height        float64
	OffsetX       float64
	OffsetY       float64
	Rotated       bool
}

type StockResult struct {
	Stock       StockPiece
	Assignments []Assignment
	Cuts        []Cut
	WasteLength float64
	WasteArea   float64 // 2D only
}

type CutPlan struct {
	Mode      int // 1 or 2
	Results   []StockResult
	WastePct  float64
	Purchased []StockPiece
	Unfit     []RequiredPiece
}

type OutputFormat string

const (
	OutputText  OutputFormat = "text"
	OutputJSON  OutputFormat = "json"
	OutputASCII OutputFormat = "ascii"
)
