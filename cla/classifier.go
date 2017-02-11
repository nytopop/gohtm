// Package cla provides implementation agnostic classifiers for gohtm.
package cla

/* Classifier functionality

Receive copy of input SDR paired with the encoded vector.
Receive predicted SDR, and sort previous SDRs by overlap score.

Previous inputs:
1111000000000000 : 0
0001111000000000 : 1
0000001111000000 : 2
0000000001111000 : 3
0000000000001111 : 4

Receives:

0111100000000000

Outputs :
0 : 3 overlap
1 : 2 overlap
2 : 0 overlap ...

This provides a pseudo-probability function for predictions
As in, x input sdr is probably referring to z previous SDR.
_This is not a prediction of likelihood!!!_

count how many synapses are active/matching for confidence value
*/

// Classifier asdf
type Classifier interface {
	Associate(active, vector []bool)
	Classify(prediction []bool)
}
