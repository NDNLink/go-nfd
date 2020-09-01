package runningstat

/*
#include "../../csrc/core/running-stat.h"
*/
import "C"
import (
	"encoding/json"
	"math"
)

// Snapshot contains a snapshot of RunningStat reading.
type Snapshot struct {
	v runningStat
}

// Count returns number of inputs.
func (s Snapshot) Count() uint64 {
	return s.v.I
}

// Len returns number of samples.
func (s Snapshot) Len() uint64 {
	return s.v.N
}

// Min returns minimum value, if enabled.
func (s Snapshot) Min() float64 {
	if s.v.N == 0 {
		return math.NaN()
	}
	return s.v.Min
}

// Max returns maximum value, if enabled.
func (s Snapshot) Max() float64 {
	if s.v.N == 0 {
		return math.NaN()
	}
	return s.v.Max
}

// Mean returns mean value.
func (s Snapshot) Mean() float64 {
	if s.v.N == 0 {
		return math.NaN()
	}
	return s.v.M1
}

// Variance returns variance of samples.
func (s Snapshot) Variance() float64 {
	if s.v.N <= 1 {
		return math.NaN()
	}
	return s.v.M2 / float64(s.v.N-1)
}

// Stdev returns standard deviation of samples.
func (s Snapshot) Stdev() float64 {
	return math.Sqrt(s.Variance())
}

// Add combines stats with another instance.
func (s Snapshot) Add(o Snapshot) (sum Snapshot) {
	if s.v.I == 0 {
		return o
	} else if o.v.I == 0 {
		return s
	}
	sum.v.I = s.v.I + o.v.I
	sum.v.N = s.v.N + o.v.N
	sum.v.Min = math.Min(s.v.Min, o.v.Min)
	sum.v.Max = math.Max(s.v.Max, o.v.Max)
	aN := float64(s.v.N)
	bN := float64(o.v.N)
	cN := float64(sum.v.N)
	delta := o.v.M1 - s.v.M1
	delta2 := delta * delta
	sum.v.M1 = (aN*s.v.M1 + bN*o.v.M1) / cN
	sum.v.M2 = s.v.M2 + o.v.M2 + delta2*aN*bN/cN
	return
}

// Sub subtracts stats in another instance.
func (s Snapshot) Sub(o Snapshot) (diff Snapshot) {
	diff.v.I = s.v.I - o.v.I
	diff.v.N = s.v.N - o.v.N
	diff.v.Min = math.NaN()
	diff.v.Max = math.NaN()
	cN := float64(s.v.N)
	aN := float64(o.v.N)
	bN := float64(diff.v.N)
	diff.v.M1 = (cN*s.v.M1 - aN*o.v.M1) / bN
	delta := o.v.M1 - diff.v.M1
	delta2 := delta * delta
	diff.v.M2 = s.v.M2 - o.v.M2 - delta2*aN*bN/cN
	return diff
}

// Scale multiplies every number by a ratio.
func (s Snapshot) Scale(ratio float64) (o Snapshot) {
	o = s
	o.v.Min *= ratio
	o.v.Max *= ratio
	o.v.M1 *= ratio
	o.v.M2 *= ratio * ratio
	return o
}

// MarshalJSON implements json.Marshaler interface.
func (s Snapshot) MarshalJSON() ([]byte, error) {
	m := make(map[string]interface{})
	m["count"] = s.Count()
	m["len"] = s.Len()
	m["m1"] = s.v.M1
	m["m2"] = s.v.M2

	addUnlessNaN := func(key string, value float64) {
		if !math.IsNaN(value) {
			m[key] = value
		}
	}
	addUnlessNaN("min", s.Min())
	addUnlessNaN("max", s.Max())
	addUnlessNaN("mean", s.Mean())
	addUnlessNaN("variance", s.Variance())
	addUnlessNaN("stdev", s.Stdev())
	return json.Marshal(m)
}

// UnmarshalJSON implements json.Unmarshaler interface.
func (s *Snapshot) UnmarshalJSON(p []byte) (e error) {
	m := make(map[string]interface{})
	if e = json.Unmarshal(p, &m); e != nil {
		return e
	}

	readNum := func(key string) float64 {
		if i, ok := m[key]; ok {
			if v, ok := i.(float64); ok {
				return v
			}
		}
		return math.NaN()
	}
	s.v.I = uint64(readNum("count"))
	s.v.N = uint64(readNum("len"))
	s.v.Min = readNum("min")
	s.v.Max = readNum("max")
	s.v.M1 = readNum("m1")
	s.v.M2 = readNum("m2")
	return nil
}
