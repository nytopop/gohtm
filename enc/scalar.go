package enc

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

func (s *ScalarEncoder) Encode(d interface{}) []bool {
	out := make([]bool, s.Bits)
	i := s.P.Buckets * (d.(int) - s.P.Min) / s.Range
	for j := 0; j < s.P.Active; j++ {
		if (i + j) < (len(out) - 1) {
			out[i+j] = true
		} else {
			if s.P.Wrap {
				// TODO wraparound
			}
		}
	}
	return out
}

func (s *ScalarEncoder) Decode(sv []bool) interface{} {
	for i, v := range sv {
		if v {
			return (i+1)*s.Range/s.P.Buckets - s.P.Min
		}
	}
	return 0
}
