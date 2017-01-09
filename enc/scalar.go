package enc

import "github.com/nytopop/gohtm/vec"

type ScalarEncoderParams struct {
	Buckets  int
	Min, Max int
	Active   int
	Wrap     bool
}

// NewScalarEncoderParams returns a default param set.
func NewScalarEncoderParams() ScalarEncoderParams {
	return ScalarEncoderParams{
		Buckets: 64,
		Min:     0,
		Max:     63,
		Active:  10,
		Wrap:    false,
	}
}

// ScalarEncoder is a linearly derived scalar encoder.
type ScalarEncoder struct {
	P     ScalarEncoderParams
	Bits  int
	Range int
}

func NewScalarEncoder(p ScalarEncoderParams) *ScalarEncoder {
	return &ScalarEncoder{
		P:     p,
		Bits:  p.Buckets + p.Active - 1,
		Range: p.Max - p.Min, // should be Absolute value
	}
}

func (s *ScalarEncoder) Encode(d interface{}) vec.SparseBinaryVector {
	out := vec.NewSparseBinaryVector(s.Bits)
	i := s.P.Buckets * (d.(int) - s.P.Min) / s.Range
	for j := 0; j < s.P.Active; j++ {
		out.Set(i+j, true)
	}
	return out
}

func (s *ScalarEncoder) Decode(sv vec.SparseBinaryVector) interface{} {
	for i, v := range sv.Dense() {
		if v {
			return (i+1)*s.Range/s.P.Buckets - s.P.Min
		}
	}
	return 0
}
