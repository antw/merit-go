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
		produced := producer.LoadAt(frame)

		if produced > remaining {
			produced -= remaining
			remaining = 0

			if produced > 0 {
				for _, flex := range order.Flexibles {
					produced -= flex.AssignExcessAt(frame, produced)

					// If there is no energy remaining to be assigned, we can exit
					// early and - as an added bonus - prevent assigning tiny
					// negatives resulting from floating point errors. This would
					// otherwise mess up technologies which have a Reserve whose
					// volume is 0.0.
					if produced <= 0.0 {
						break
					}
				}
			}
		} else if produced < remaining {
			// The producer is providing less energy than remaining demand. Take
			// it all and continue with the next producer.
			remaining -= produced
		}

		if remaining < 0 {
			remaining = 0
		}
	}

	for _, producer := range order.Flexibles {
		maxLoad := producer.AvailableAt(frame)

		if remaining > 0 && maxLoad < remaining {
			producer.SetLoadAt(frame, maxLoad)
		} else {
			// remaining is less than 0 if always-on supply exceeds demand.
			if remaining > 0 {
				producer.SetLoadAt(frame, remaining)
			}

			if len(order.Dispatchables) > 0 {
				// TODO Add a test! Panics without the above conditional when
				// there are no dispatchables.
				order.PriceSetters[frame] = order.Dispatchables[0]
			}

			break // All demand is assigned.
		}

		remaining -= maxLoad
	}

	for _, producer := range order.Dispatchables {
		maxLoad := producer.TotalCapacity()

		if maxLoad < remaining {
			producer.SetLoadAt(frame, maxLoad)
		} else {
			// remaining is less than 0 if always-on supply exceeds demand.
			if remaining > 0 {
				producer.SetLoadAt(frame, remaining)
			}

			order.PriceSetters[frame] = producer
			break // All demand is assigned.
		}

		remaining -= maxLoad
	}
}
