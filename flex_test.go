package merit

import "testing"

func TestStorageStartsEmpty(t *testing.T) {
	storage := Storage{
		Flex: Flex{
			Key:      "abc",
			Capacity: 10.0,
			Units:    2.0,
		},
		reserve: NewReserveWithoutDecay(50.0),
	}

	for frame := range []int{0, 1, 8759} {
		if load := storage.LoadAt(frame); load != 0 {
			t.Errorf("Storage.LoadAt(%d) = %f, want 0.0", frame, load)
		}
	}
}

func TestAssigningExcess(t *testing.T) {
	storage := Storage{
		Flex: Flex{
			Key:      "abc",
			Capacity: 10.0,
			Units:    2.0,
		},
		reserve: NewReserveWithoutDecay(50.0),
	}

	tests := []struct {
		amount float64
		expect float64
		total  float64
	}{
		{2.0, 2.0, 2.0},
		{10.0, 10.0, 12.0},
		{10.0, 8.0, 20.0},
		{10.0, 0.0, 20.0},
	}

	for _, test := range tests {
		assigned := storage.AssignExcessAt(0, test.amount)

		if assigned != test.expect {
			t.Errorf("Storage.AssignExcessAt(0, %f) = %f, want %f",
				test.amount, assigned, test.expect)
		}

		if load := storage.LoadAt(0); load != -test.total {
			t.Errorf("Storage.LoadAt(0, %f) = %f, want %f",
				test.amount, load, -test.total)
		}

		if stored := storage.reserve.At(0); stored != test.total {
			t.Errorf("Storage.reserve.At(0, %f) = %f, want %f",
				test.amount, stored, test.total)
		}
	}
}

func TestReserveCapacity(t *testing.T) {
	storage := Storage{
		Flex: Flex{
			Key:      "abc",
			Capacity: 100.0,
			Units:    2.0,
		},
		reserve: NewReserveWithoutDecay(50.0),
	}

	tests := []struct {
		amount float64
		expect float64
		total  float64
	}{
		{2.0, 2.0, 2.0},
		{10.0, 10.0, 12.0},
		{10.0, 10.0, 22.0},
		{50.0, 28.0, 50.0},
	}

	for _, test := range tests {
		assigned := storage.AssignExcessAt(0, test.amount)

		if assigned != test.expect {
			t.Errorf("Storage.AssignExcessAt(0, %f) = %f, want %f",
				test.amount, assigned, test.expect)
		}

		if load := storage.LoadAt(0); load != -test.total {
			t.Errorf("Storage.LoadAt(0, %f) = %f, want %f",
				test.amount, load, -test.total)
		}

		if stored := storage.reserve.At(0); stored != test.total {
			t.Errorf("Storage.reserve.At(0, %f) = %f, want %f",
				test.amount, stored, test.total)
		}
	}
}

func BenchStorage(b *testing.B) {
	for i := 0; i < b.N; i++ {
		b.StopTimer()
		stor := Storage{
			Flex: Flex{
				Key:      "abc",
				Capacity: 10.0,
			},
			reserve: NewReserveWithoutDecay(10.0),
		}
		stor.AssignExcessAt(0, 5.0)
		b.StartTimer()

		for frame := 0; frame < 8760; frame++ {
			// stor.AssignExcessAt(frame, 1.0)
			stor.SetLoadAt(frame, 5.0)
			// stor.AvailableAt(frame)
		}
	}
}
