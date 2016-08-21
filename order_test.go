package merit

import "testing"

var consumerOne = Consumer{
	Profile:     [8760]float64{0: 0.4, 1: 0.2},
	TotalDemand: 20.0,
}

var consumerTwo = Consumer{
	Profile:     [8760]float64{0: 0.2, 1: 0.1},
	TotalDemand: 10.0,
}

func TestOrderDemandAt(t *testing.T) {
	order := Order{Consumers: []*Consumer{&consumerOne, &consumerTwo}}

	tests := []struct {
		frame int
		want  float64
	}{
		{0, 10.0},
		{1, 5.0},
		{2, 0.0},
	}

	for _, test := range tests {
		if demand := order.DemandAt(test.frame); demand != test.want {
			t.Errorf("Order.DemandAt(%d) = %f, want %f",
				test.frame, demand, test.want)
		}
	}
}

func TestEmptyOrderDemandAt(t *testing.T) {
	order := Order{}

	if demand := order.DemandAt(0); demand != 0.0 {
		t.Errorf("Order.DemandAt(0) = %f, want 0.0", demand)
	}
}
