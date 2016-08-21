package merit

import "sort"

// Calculate receives a merit order and computes which producers are running in
// each frame, and at what level of production, in order to meet demand.
func Calculate(order Order) {
	sort.Sort(order.Dispatchables)

	for frame := 0; frame < 8760; frame++ {
		calculateFrame(frame, order)
	}
}

func calculateFrame(frame int, order Order) {
	remaining := order.DemandAt(frame)

	for _, producer := range order.AlwaysOns {
		remaining -= producer.LoadAt(frame)
	}

	for _, producer := range order.Dispatchables {
		maxLoad := producer.TotalCapacity()

		if maxLoad < remaining {
			producer.SetLoadAt(frame, maxLoad)
		} else if remaining > 0 {
			producer.SetLoadAt(frame, remaining)
			break // All demand is assigned.
		} else {
			break // All demand is assigned. TODO Can we ever get here?
		}

		remaining -= maxLoad
	}
}
