package merit

// AlwaysOn is a type of energy producer which has a fixed amount of production
// in each frame. This amount is produced whether there is demand or not;
// therefore energy from AlwaysOn producers is consumed before using
// dispatchables.
type AlwaysOn struct {
	Key             string
	Profile         [8760]float64
	TotalProduction float64
}

// LoadAt returns the load of the dispatchable in frame. May return nil if no
// load is yet assigned.
func (a *AlwaysOn) LoadAt(frame int) float64 {
	return a.Profile[frame] * a.TotalProduction
}
