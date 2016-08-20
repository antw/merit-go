package merit

import "testing"

func TestDispatchableTotalCapacity(t *testing.T) {
	tests := []struct {
		capacity, units, wants float64
	}{
		{10.5, 1, 10.5},
		{0.55, 2, 1.1},
		{0.0, 2, 0.0},
		{2.0, 0, 0},
	}

	for _, test := range tests {
		dis := Dispatchable{Capacity: test.capacity, Units: test.units}

		if actual := dis.TotalCapacity(); actual != test.wants {
			t.Errorf("Dispatchable{C: %f, U: %f} total capacity = %f, wants %f",
				test.capacity, test.units, actual, 10.5)
		}
	}
}

func TestDispatchableSetLoadAt(t *testing.T) {
	tests := []struct {
		frame  int
		amount float64
	}{
		{0, 5.0},
		{2, 2.0},
		{8759, 8.0},
	}

	for _, test := range tests {
		dis := Dispatchable{}
		dis.SetLoadAt(test.frame, test.amount)

		if set := dis.load[test.frame]; set != test.amount {
			t.Errorf("Dispatchable.SetLoadAt(%d, %f) set %f",
				test.frame, test.amount, set)
		}
	}
}

func TestDispatchableSetLoadAtBadIndex(t *testing.T) {
	dis := Dispatchable{}

	if dis.SetLoadAt(8760, 5.0) == nil {
		t.Errorf("Dispatchable.SetLoadAt(8760) should return an error")
	}
}
