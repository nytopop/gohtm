package enc

// AudioEncoder encodes a sound sample by performing FFT, and splitting
// frequency domains into intensity buckets.
type AudioEncoder struct {
}

func NewAudioEncoder() *AudioEncoder {
	return &AudioEncoder{}
}

func (a *AudioEncoder) Encode(d interface{}) []bool {
	return []bool{}
}

func (a *AudioEncoder) Decode(sv []bool) interface{} {
	return 0
}
