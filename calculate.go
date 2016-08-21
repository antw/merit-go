package merit

import (
	"sort"
	"sync"
)

// Calculate receives a merit order and computes which producers are running in
// each frame, and at what level of production, in order to meet demand.
func Calculate(order Order) {
	sort.Sort(order.Dispatchables)

	for frame := 0; frame < 8760; frame++ {
		calculateFrame(frame, order)
	}
}

// CalculateParallel receives a merit order and computes the batches of frames
// in goroutines. Not suitable for merit orders which use electricity storage.
func CalculateParallel(order Order, batches int) {
	var wg sync.WaitGroup

	batchSize := 8760 / batches

	for i := 0; i < batches; i++ {
		wg.Add(1)

		go func(i int) {
			calculateFrameBatch(batchSize*i, (batchSize * (i + 1)), order)
			wg.Done()
		}(i)
	}

	wg.Wait()
}

func calculateFrameBatch(start, end int, order Order) {
	for frame := start; frame < end; frame++ {
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
			order.PriceSetters[frame] = producer
			break // All demand is assigned.
		} else {
			order.PriceSetters[frame] = producer
			break // All demand is assigned. TODO Can we ever get here?
		}

		remaining -= maxLoad
	}
}
