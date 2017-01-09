package enc

import "github.com/nytopop/gohtm/vec"

/* Encoder Design Guidelines
1. Semantically similar data should result in SDRs with overlapping active bits.
2. The same input should always produce the same SDR as output.
3. The output should have the same dimensionality (total number of bits) for all inputs.
4. The output should have similar sparsity for all inputs and have enough one-bits to handle noise and subsampling.
*/

// Encoder is an interface for all sparse encoders. A valid Encoder
// should implement a Decode and an Encode method, persistent
// state is not required.
type Encoder interface {
	Encode(interface{}) vec.SparseBinaryVector
	Decode(vec.SparseBinaryVector) interface{}
}
