package merit

import (
	"math"
	"testing"
)

func round(num float64) int {
	return int(num + math.Copysign(0.5, num))
}

func toFixed(num float64, precision int) float64 {
	output := math.Pow(10, float64(precision))
	return float64(round(num*output)) / output
}

func TestCalculateOneDispatchable(t *testing.T) {
	disp := Dispatchable{Capacity: 0.5, Units: 3.0}
	cons := Consumer{Profile: [8760]float64{0.2, 0.4, 1.0}, TotalDemand: 2.0}

	order := Order{
		Consumers:     []*Consumer{&cons},
		Dispatchables: DispatchableList{&disp},
	}

	Calculate(order)

	tests := []struct {
		frame int
		want  float64
	}{
		{0, 0.4},
		{1, 0.8},
		{2, 1.5}, // Limited by dispatchable capacity.
		{3, 0.0},
	}

	for _, test := range tests {
		if load := disp.LoadAt(test.frame); load != test.want {
			t.Errorf("Calculate assigned dispatchable load %f, want %f",
				load, test.want)
		}
	}
}

// Asserts that the cheaper dispatchable is used first.
func TestCalculateTwoDispatchables(t *testing.T) {
	d1 := Dispatchable{Cost: 2.0, Capacity: 0.5, Units: 2.0}
	d2 := Dispatchable{Cost: 1.0, Capacity: 0.5, Units: 2.0}

	cons := Consumer{Profile: [8760]float64{0.2, 0.4, 0.8}, TotalDemand: 2.0}

	order := Order{
		Consumers:     []*Consumer{&cons},
		Dispatchables: DispatchableList{&d1, &d2},
	}

	Calculate(order)

	tests := []struct {
		frame int
		want1 float64
		want2 float64
	}{
		{0, 0.0, 0.4},
		{1, 0.0, 0.8},
		{2, 0.6, 1.0},
		{3, 0.0, 0.0},
	}

	for _, test := range tests {
		if load := d1.LoadAt(test.frame); toFixed(load, 10) != test.want1 {
			t.Errorf("Calculate assigned dispatchable1 load %f, want %f",
				load, test.want1)
		}

		if load := d2.LoadAt(test.frame); toFixed(load, 10) != test.want2 {
			t.Errorf("Calculate assigned dispatchable2 load %f, want %f",
				load, test.want2)
		}
	}
}

func TestCalculateOneAOOneDisp(t *testing.T) {
	ao := AlwaysOn{Profile: [8760]float64{0.5, 0.5, 0.5}, TotalProduction: 1.0}
	disp := Dispatchable{Key: "only", Capacity: 0.5, Units: 3.0}
	cons := Consumer{Profile: [8760]float64{0.2, 0.4, 1.0}, TotalDemand: 2.0}

	order := Order{
		Consumers:     []*Consumer{&cons},
		AlwaysOns:     []*AlwaysOn{&ao},
		Dispatchables: DispatchableList{&disp},
	}

	Calculate(order)

	tests := []struct {
		frame int
		want  float64
	}{
		{0, 0.0},
		{1, 0.3},
		{2, 1.5},
		{3, 0.0},
	}

	for _, test := range tests {
		if load := disp.LoadAt(test.frame); toFixed(load, 10) != test.want {
			t.Errorf("Calculate assigned dispatchable load %d = %f, want %f",
				test.frame, load, test.want)
		}
	}
}

func BenchmarkCalculate(b *testing.B) {
	ao := AlwaysOn{Profile: [8760]float64{0.5, 0.5, 0.5}, TotalProduction: 1.0}

	d1 := Dispatchable{Capacity: 0.5, Units: 3.0}
	d2 := Dispatchable{Capacity: 0.5, Units: 3.0}
	d3 := Dispatchable{Capacity: 0.5, Units: 3.0}
	d4 := Dispatchable{Capacity: 0.5, Units: 3.0}

	cons := Consumer{Profile: [8760]float64{0.2, 0.4, 1.0}, TotalDemand: 8.0}

	order := Order{
		Consumers:     []*Consumer{&cons},
		AlwaysOns:     []*AlwaysOn{&ao},
		Dispatchables: DispatchableList{&d1, &d2, &d3, &d4},
	}

	for i := 0; i < b.N; i++ {
		Calculate(order)
	}
}
