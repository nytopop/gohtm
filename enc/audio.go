package enc

// AudioEncoder encodes a sound sample by performing FFT, and splitting
// frequency domains into intensity buckets.
type AudioEncoder struct {
}

// NewAudioEncoder asdf
func NewAudioEncoder() *AudioEncoder {
	return &AudioEncoder{}
}

// Encode asdf
func (a *AudioEncoder) Encode(d interface{}) []bool {
	return []bool{}
}

// Decode asdf
func (a *AudioEncoder) Decode(sv []bool) interface{} {
	return 0
}
