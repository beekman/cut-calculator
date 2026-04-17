package model

type StockPiece struct {
	Length float64
	Count  int
	OnHand bool
}

type RequiredPiece struct {
	Label  string
	Length float64
	Count  int
}

type Cut struct {
	Position float64
	Label    string
}

type Assignment struct {
	StockIndex    int
	RequiredLabel string
	Length        float64
}

type StockResult struct {
	Stock       StockPiece
	Assignments []Assignment
	Cuts        []Cut
	WasteLength float64
}

type CutPlan struct {
	Results   []StockResult
	WastePct  float64
	Purchased []StockPiece
	Unfit     []RequiredPiece
}

type OutputFormat string

const (
	OutputText OutputFormat = "text"
	OutputJSON OutputFormat = "json"
	OutputASCII OutputFormat = "ascii"
)
