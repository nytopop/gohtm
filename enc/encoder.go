/*
Package enc provides an implementation agnostic
interface for encoding and decoding data into
semi-sparse distributed representations consumable
by a spatial pooling algorithm.

Encoders also return a bucket index that can be
passed back through the algorithm for a decoded
representation (+/- encoder resolution).
*/
package enc

/* Encoder Design Guidelines
1. Semantically similar data should result in SDRs with overlapping active bits.
2. The same input should always produce the same SDR as output.
3. The output should have the same dimensionality (total number of bits) for all inputs.
4. The output should have similar sparsity for all inputs and have enough one-bits to handle noise and subsampling.
*/

// Encoder is an interface for all sparse encoders.
type Encoder interface {
	Encode(interface{}) ([]bool, int)
	Decode([]bool) interface{}
}
