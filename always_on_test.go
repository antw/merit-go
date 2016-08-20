package merit

import "testing"

func TestAlwaysOnLoadAt(t *testing.T) {
	alwaysOn := AlwaysOn{
		Profile:         [8760]float64{0: 0.1, 1: 0.3},
		TotalProduction: 20.0,
	}

	tests := []struct {
		frame int
		want  float64
	}{
		{0, 2.0},
		{1, 6.0},
		{2, 0.0}, // Zero value in profile.
	}

	for _, test := range tests {
		if load := alwaysOn.LoadAt(test.frame); load != test.want {
			t.Errorf("AlwaysOn.LoadAt(%d) = %f, want %f",
				test.frame, load, test.want)
		}
	}
}
