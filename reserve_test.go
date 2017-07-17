package merit

import (
	"sync"
	"testing"
)

func testDecay(t *testing.T, res *reserve, frame int, expected float64) {
	if res.DecayAt(frame) != expected {
		t.Errorf("Reserve.DecayAt(%d) = %f, wants %f",
			frame, res.DecayAt(frame), expected)
	}
}

func testAt(t *testing.T, res *reserve, frame int, expected float64) {
	if res.At(frame) != expected {
		t.Errorf("Reserve.At(%d) = %f, wants %f",
			frame, res.At(frame), expected)
	}
}

func TestReserveStartsEmpty(t *testing.T) {
	res := NewReserveWithoutDecay(2.0)

	for frame := range []int{0, 1, 8759} {
		testAt(t, &res, frame, 0.0)
	}
}

func TestReserveCarriesPreviousValues(t *testing.T) {
	res := NewReserveWithoutDecay(2.0)

	res.Add(0, 2.0)

	for frame := 0; frame < 3; frame++ {
		testAt(t, &res, frame, 2.0)
	}
}

func TestReserveCarriesWithDecay(t *testing.T) {
	res := NewReserve(2.0, func(frame int, stored float64) float64 {
		return stored * 0.25
	})

	res.Add(0, 2.0)

	tests := []struct {
		frame  int
		stored float64
		decay  float64
	}{
		{0, 2.0, 0.0},
		{1, 1.5, 0.5},
		{2, 1.125, 0.375},
	}

	for _, test := range tests {
		testDecay(t, &res, test.frame, test.decay)
		testAt(t, &res, test.frame, test.stored)
	}
}

func TestDecayAfterTake(t *testing.T) {
	res := NewReserve(2.0, func(frame int, stored float64) float64 {
		return 0.5
	})

	res.Add(0, 2.0)

	testDecay(t, &res, 1, 0.5)

	res.Take(1, 1.0)

	testDecay(t, &res, 1, 0.5)
	testDecay(t, &res, 2, 0.5)
}

func TestDecayFixedAmount(t *testing.T) {
	res := NewReserve(10.0, func(frame int, stored float64) float64 {
		return 2.0
	})

	res.Add(0, 3.0)

	tests := []struct {
		frame  int
		decay  float64
		stored float64
	}{
		{0, 0.0, 3.0},
		{1, 2.0, 1.0},
		{2, 1.0, 0.0},
		{3, 0.0, 0.0},
	}

	for _, test := range tests {
		testDecay(t, &res, test.frame, test.decay)
		testAt(t, &res, test.frame, test.stored)
	}
}

func TestDecayByFrame(t *testing.T) {
	res := NewReserve(10.0, func(frame int, stored float64) float64 {
		if frame%2 == 0 {
			return 2.0
		}

		return 0.0
	})

	res.Add(0, 5.0)

	tests := []struct {
		frame  int
		decay  float64
		stored float64
	}{
		{0, 0.0, 5.0},
		{1, 0.0, 5.0},
		{2, 2.0, 3.0},
		{3, 0.0, 3.0},
		{4, 2.0, 1.0},
		{5, 0.0, 1.0},
		{6, 1.0, 0.0},
	}

	for _, test := range tests {
		testDecay(t, &res, test.frame, test.decay)
		testAt(t, &res, test.frame, test.stored)
	}
}

func TestReserveTakeAll(t *testing.T) {
	res := NewReserveWithoutDecay(10.0)

	res.Add(0, 1.0)

	if taken := res.Take(1, 1.0); taken != 1.0 {
		t.Errorf("Reserve.Take(1, 1.0) = %f, wants %f", taken, 1.0)
	}

	testAt(t, &res, 0, 1.0)
	testAt(t, &res, 1, 0.0)
}

func TestReserveTakeSome(t *testing.T) {
	res := NewReserveWithoutDecay(10.0)

	res.Add(0, 5.0)

	if taken := res.Take(1, 2.0); taken != 2.0 {
		t.Errorf("Reserve.Take(1, 2.0) = %f, wants %f", taken, 2.0)
	}

	testAt(t, &res, 1, 3.0)
}

func TestReserveTakeTooMuch(t *testing.T) {
	res := NewReserveWithoutDecay(10.0)

	res.Add(0, 5.0)

	if taken := res.Take(1, 7.0); taken != 5.0 {
		t.Errorf("Reserve.Take(1, 7.0) = %f, wants %f", taken, 5.0)
	}

	testAt(t, &res, 1, 0.0)
}

func TestVolumeLimit(t *testing.T) {
	res := NewReserveWithoutDecay(2.0)

	tests := []struct {
		amount float64
		added  float64
		stored float64
	}{
		{1.0, 1.0, 1.0},
		{1.0, 1.0, 2.0},
		{1.0, 0.0, 2.0},
		{1.0, 0.0, 2.0},
	}

	for _, test := range tests {
		if added := res.Add(0, test.amount); added != test.added {
			t.Errorf("Reserve.Add(0, %f) = %f, wants %f",
				test.amount, added, test.added)
		}

		testAt(t, &res, 0, test.stored)
	}
}

func BenchmarkReserveAdd(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		res := NewReserveWithoutDecay(8760.0)
		res.Set(0, 0)
		b.StartTimer()

		for frame := 0; frame < 8760; frame++ {
			res.Add(frame, 1.0)
		}
	}
}

func BenchmarkReserveTake(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		res := NewReserveWithoutDecay(8760.0)
		res.Set(0, 8760.0)
		b.StartTimer()

		for frame := 0; frame < 8760; frame++ {
			res.Take(frame, 1.0)
		}
	}
}

func BenchmarkReserveNilDecay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		res := NewReserveWithoutDecay(20000.0)
		res.Set(0, 20000.0)
		b.StartTimer()

		for frame := 0; frame < 8760; frame++ {
			res.At(frame)
			res.At(frame)
			res.At(frame)
		}
	}
}

func TestReserveNilDecay(t *testing.T) {
	res := NewReserveWithoutDecay(20000.0)
	res.Set(0, 20000.0)

	for frame := 0; frame < 10; frame++ {
		res.At(frame)
		res.At(frame)
		res.At(frame)
	}
}

func BenchmarkReserveDecay(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		res := NewReserve(
			20000.0,
			func(frame int, stored float64) float64 { return 2.0 },
		)

		res.Set(0, 20000.0)
		b.StartTimer()

		for frame := 0; frame < 8760; frame++ {
			res.At(frame)
		}
	}
}

func BenchmarkReserve(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		res := NewReserve(
			10.0,
			func(frame int, stored float64) float64 { return 2.0 },
		)
		res.Set(0, 5.0)
		b.StartTimer()

		for frame := 0; frame < 8760; frame++ {
			res.Add(frame, 1.0)
			res.Take(frame, 1.0)
			res.At(frame)
		}
	}
}

func BenchmarkReserveParallel(b *testing.B) {
	for i := 0; i < b.N; i++ {

		b.StopTimer()
		res := NewReserve(
			10.0,
			func(frame int, stored float64) float64 { return 2.0 },
		)
		res.Set(0, 5.0)
		b.StartTimer()

		var wg sync.WaitGroup

		batchSize := 8760 / 4

		for i := 0; i < 4; i++ {
			wg.Add(1)

			start := batchSize * i
			end := batchSize * (i + 1)

			go func(i int) {
				for frame := start; frame < end; frame++ {
					res.Add(frame, 1.0)
					res.Take(frame, 1.0)
					res.At(frame)
				}
				wg.Done()
			}(i)
		}

		wg.Wait()
	}
}
