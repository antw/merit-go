package merit

// Consumer is a merit order participant which uses energy. Dispatchable and
// AlwaysOn plants will be used to satisfy the demand of Consumers.
type Consumer struct {
	Key         string
	Profile     [8760]float64
	TotalDemand float64
}

// LoadAt returns the energy used by the Consumer in frame.
func (c *Consumer) LoadAt(frame int) float64 {
	return c.Profile[frame] * c.TotalDemand
}
