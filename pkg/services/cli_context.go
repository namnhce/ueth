package services

type CLIContext interface {
	String(s string) string
	Int(s string) int
	Float64(s string) float64
	Bool(s string) bool
}
