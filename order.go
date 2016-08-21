package merit

// Order contains information about the participants in the merit order.
type Order struct {
	Consumers     []*Consumer
	AlwaysOns     []*AlwaysOn
	Dispatchables DispatchableList
	PriceSetters  []*Dispatchable
}

// NewOrder creates and returns new merit order. Prefer this over creating an
// Order{} directly.
func NewOrder() Order {
	return Order{PriceSetters: make([]*Dispatchable, 8760)}
}

// DemandAt returns the total demand for energy in frame.
func (o *Order) DemandAt(frame int) float64 {
	var sum float64

	for _, consumer := range o.Consumers {
		sum += consumer.LoadAt(frame)
	}

	return sum
}

// AddConsumer adds a Consumer to the merit order.
func (o *Order) AddConsumer(c *Consumer) {
	o.Consumers = append(o.Consumers, c)
}

// AddAlwaysOn adds an AlwaysOn producer to the merit order.
func (o *Order) AddAlwaysOn(a *AlwaysOn) {
	o.AlwaysOns = append(o.AlwaysOns, a)
}

// AddDispatchable adds a Dispatchable producer to the merit order.
func (o *Order) AddDispatchable(d *Dispatchable) {
	o.Dispatchables = append(o.Dispatchables, d)
}
