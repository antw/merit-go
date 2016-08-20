// Package merit computes ranks sources of energy by to cost and assigns load
// to producers accordingly.
package merit

type participant interface {
	LoadAt(int) float64
}
