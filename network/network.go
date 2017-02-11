// Package network provides an implementation agnostic interface for
// building, linking, and managing regions.
package network

// Network asdf
type Network interface {
	Serialize() []byte
}

// Essentially, we need to store a graph representation of the network.
