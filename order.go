package merit

// Order contains information about the participants in the merit order.
type Order struct {
	Consumers     []Consumer
	AlwaysOns     []AlwaysOn
	Dispatchables DispatchableList
}

// DemandAt returns the total demand for energy in frame.
func (o *Order) DemandAt(frame int) float64 {
	var sum float64

	for _, consumer := range o.Consumers {
		sum += consumer.LoadAt(frame)
	}

	return sum
}
