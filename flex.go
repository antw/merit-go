package merit

import "fmt"

type Flexlike interface {
	AssignExcessAt(int, float64) float64
	SetLoadAt(int, float64) error
	AvailableAt(int) float64
}

type Flex struct {
	Key      string
	Capacity float64
	Units    float64
	load     [8760]float64
}

// TotalCapacity returns the total amount of energy which may be produced or
// consumed by the flex in each hour in kWh.
func (f *Flex) TotalCapacity() float64 {
	return f.Capacity * f.Units
}

func (f *Flex) AvailableAt(frame int) float64 {
	return 0.0
}

// AssignExcessAt informs the flex that an amount of energy has been assigned to
// be stored or otherwise used by the participant.
func (f *Flex) AssignExcessAt(frame int, amount float64) float64 {
	// TODO Assert valid frame.

	input_capacity := f.TotalCapacity() + f.load[frame]

	if amount > input_capacity {
		amount = input_capacity
	}

	f.load[frame] = f.load[frame] - amount

	return amount
}

// SetLoadAt assigns a load to the flex in the chosen frame. The amount should
// not exceed the total capacity, but SetLoadAt does not assert that this is the
// case.
func (f *Flex) SetLoadAt(frame int, amount float64) error {
	if frame > cap(f.load)-1 {
		return fmt.Errorf(
			"Flex.SetLoadAt: Cannot assign to out of range index %d",
			frame)
	}

	f.load[frame] = amount
	return nil
}

// LoadAt returns the load of the dispatchable in the frame. May return nil if
// no load is yet assigned.
func (f *Flex) LoadAt(frame int) float64 {
	return f.load[frame]
}

type Storage struct {
	Flex
	reserve reserve
}

func (s *Storage) AssignExcessAt(frame int, amount float64) float64 {
	input_cap := s.TotalCapacity() + s.load[frame]

	if amount > input_cap {
		amount = input_cap
	}

	// TODO adjust for input efficiency

	stored := s.reserve.Add(frame, amount)
	s.load[frame] = s.load[frame] - stored

	return stored
}

// SetLoadAt assigns a load to the flex in the chosen frame. The amount should
// not exceed the total capacity, but SetLoadAt does not assert that this is the
// case.
func (s *Storage) SetLoadAt(frame int, amount float64) error {
	// if frame > cap(f.load)-1 {
	// 	return fmt.Errorf(
	// 		"Flex.SetLoadAt: Cannot assign to out of range index %d",
	// 		frame)
	// }

	s.load[frame] = s.reserve.Take(frame, amount)
	return nil
}

func (s *Storage) AvailableAt(frame int) float64 {
	available := s.reserve.At(frame)
	capacity := s.Capacity * s.Units

	if available > capacity {
		return capacity
	}

	return available
}
