package merit

import "fmt"

// Dispatchable describes a source of energy along with its cost, capacity,
// units, and other information required to determine its role in the merit
// order.
type Dispatchable struct {
	Key      string
	Cost     float64
	Capacity float64
	Units    float64
	load     [8760]float64
}

// TotalCapacity returns the total amount of energy which may be produced by the
// producer in each hour in kWh.
func (d *Dispatchable) TotalCapacity() float64 {
	return d.Capacity * d.Units
}

// SetLoadAt assigns a load to the dispatchable in the chosen frame. The amount
// should not exceed the total capacity, but SetLoadAt does not assert that this
// is the case.
func (d *Dispatchable) SetLoadAt(frame int, amount float64) error {
	if frame > cap(d.load)-1 {
		return fmt.Errorf(
			"Dispatchable.SetLoadAt: Cannot assign to out of range index %d",
			frame)
	}

	d.load[frame] = amount
	return nil
}

// LoadAt returns the load of the dispatchable in frame. May return nil if no
// load is yet assigned.
func (d *Dispatchable) LoadAt(frame int) float64 {
	return d.load[frame]
}
