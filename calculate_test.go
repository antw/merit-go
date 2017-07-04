package merit

import (
	"math"
	"math/rand"
	"testing"
	"time"
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

	order := NewOrder()
	order.AddConsumer(&cons)
	order.AddDispatchable(&disp)

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

	order := NewOrder()
	order.AddConsumer(&cons)
	order.AddDispatchable(&d1)
	order.AddDispatchable(&d2)

	Calculate(order)

	tests := []struct {
		frame  int
		want1  float64
		want2  float64
		wantPS *Dispatchable
	}{
		{0, 0.0, 0.4, &d2},
		{1, 0.0, 0.8, &d2},
		{2, 0.6, 1.0, &d1},
		{3, 0.0, 0.0, &d2},
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

		if setter := order.PriceSetters[test.frame]; setter != test.wantPS {
			t.Errorf("Calculate assigned price setter in frame %d = %f, want %f",
				test.frame, setter.Cost, test.wantPS.Cost)
		}
	}
}

// Asserts that the concurrent calculator assigns a load in every frame.
func TestCalculateParallel(t *testing.T) {
	seed := time.Now().UTC().UnixNano()
	t.Logf("Random seed: %d", seed)
	rand.Seed(seed)

	var consumption [8760]float64

	for i := 0; i < 8760; i++ {
		consumption[i] = rand.Float64() + 0.1
	}

	makeOrder := func() Order {
		disp := Dispatchable{Capacity: 0.5, Units: 3.0}
		cons := Consumer{Profile: consumption, TotalDemand: 2.0}

		order := NewOrder()
		order.AddConsumer(&cons)
		order.AddDispatchable(&disp)

		return order
	}

	serial := makeOrder()
	parallel := makeOrder()

	Calculate(serial)
	CalculateParallel(parallel, 4)

	sDisp := serial.Dispatchables[0]
	pDisp := parallel.Dispatchables[0]

	for i := 0; i < 8760; i++ {
		pLoad := pDisp.LoadAt(i)

		if pLoad == 0 {
			t.Errorf("Parallel dispatchable load zero in frame %d", i)
		}

		if sLoad := sDisp.LoadAt(i); pLoad != sLoad {
			t.Errorf("Parallel dispatchable load in frame %d = %f, want %f",
				i, pLoad, sLoad)
		}
	}
}

func TestCalculateOneAOOneDisp(t *testing.T) {
	ao := AlwaysOn{Profile: [8760]float64{0.5, 0.5, 0.5}, TotalProduction: 1.0}
	disp := Dispatchable{Key: "only", Capacity: 0.5, Units: 3.0}
	cons := Consumer{Profile: [8760]float64{0.2, 0.4, 1.0}, TotalDemand: 2.0}

	order := NewOrder()
	order.AddConsumer(&cons)
	order.AddAlwaysOn(&ao)
	order.AddDispatchable(&disp)

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

// Creates a merit order for use in benchmarking with n dispatchables.
func benchmarkOrder(n int) Order {
	var demand [8760]float64
	var always [8760]float64

	for i := 0; i < 8760; i++ {
		demand[i] = rand.Float64()
		always[i] = rand.Float64()
	}

	ao := AlwaysOn{Profile: always, TotalProduction: 1.0}
	cons := Consumer{Profile: demand, TotalDemand: 1.5 * float64(n)}

	order := NewOrder()

	order.AddConsumer(&cons)
	order.AddAlwaysOn(&ao)

	for i := 0; i < n; i++ {
		order.AddDispatchable(&Dispatchable{Capacity: 0.5, Units: 3.0})
	}

	return order
}

func BenchmarkCalculate(b *testing.B) {
	b.StopTimer()

	// Using a predefined seed for predictable results.
	rand.Seed(10)
	order := benchmarkOrder(40)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		Calculate(order)
	}
}

func BenchmarkCalculateParallel(b *testing.B) {
	b.StopTimer()

	// Using a predefined seed for predictable results.
	rand.Seed(10)
	order := benchmarkOrder(40)

	b.StartTimer()

	for i := 0; i < b.N; i++ {
		CalculateParallel(order, 4)
	}
}
