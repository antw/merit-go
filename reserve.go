package merit

import "math"

type decayFunc func(frame int, stored float64) float64

type reserve struct {
	Volume float64
	decay  decayFunc
	store  [8760]float64
}

func NewReserve(volume float64, decay decayFunc) reserve {
	var store [8760]float64

	for i := 1; i < 8760; i++ {
		// Set all store values except the first to -1, indicating that the
		// value has not yet been computed.
		store[i] = -1
	}

	return reserve{
		Volume: volume,
		decay:  decay,
		store:  store,
	}
}

func NewReserveWithoutDecay(volume float64) reserve {
	return NewReserve(
		volume,
		func(int, float64) float64 { return 0.0 },
	)
}

// At returns how much energy is stored in the reserve at the end of the given
// frame. If the technology to which the reserve is attached is still being
// calculated, the energy stored may be subject to change.
func (r *reserve) At(frame int) float64 {
	if r.store[frame] == -1 {
		previous := r.store[frame-1]

		if previous == -1 {
			// Previous reserve value has not been computed yet. For easy-of-use
			// in parallel calculations, we just assume it to be zero.
			r.store[frame] = 0.0
		} else {
			r.store[frame] = previous - r.DecayAt(frame)
		}
	}

	return r.store[frame]
}

// Set sets the amount in the reserve for the chosen frame. Ignores volume
// constraints and assumes you will check this yourself.
func (r *reserve) Set(frame int, amount float64) {
	r.store[frame] = amount
}

// Add adds the given amount of energy in the chosen frame, ensuring that the
// amount stored does not exceed the volume of the reserve.
//
// Returns the amount of energy which was actually added; note that this may be
// less than the amount parameter.
func (r *reserve) Add(frame int, amount float64) float64 {
	if amount <= 0 {
		return 0
	}

	stored := r.At(frame)

	if stored+amount > r.Volume {
		amount = r.Volume - stored
	}

	r.Set(frame, stored+amount)

	return amount
}

// Take removes the desired amount of energy from the reserve.
//
// Returns the amount of energy subtracted from the reserve. This may be less
// than asked for if insufficient was stored.
func (r *reserve) Take(frame int, amount float64) float64 {
	if amount <= 0 {
		return 0
	}

	stored := r.At(frame)

	if stored > amount {
		r.Set(frame, r.At(frame)-amount)
		return amount
	}

	r.Set(frame, 0.0)
	return stored
}

// DecayAt returns how much energy decayed in the reserve at the beginning of
// the chosen frame.
func (r *reserve) DecayAt(frame int) float64 {
	if frame == 0 {
		return 0.0
	}

	stored := r.At(frame - 1)
	decay := r.decay(frame, stored)

	return math.Min(stored, decay)
}
