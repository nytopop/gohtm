# go-htm 
[![GoDoc](https://godoc.org/github.com/nytopop/gohtm?status.svg)](https://godoc.org/github.com/nytopop/gohtm) [![Build Status](https://travis-ci.org/nytopop/gohtm.svg?branch=master)](https://travis-ci.org/nytopop/gohtm)

An implementation of the HTM algorithm.

In progress.

## Overview

HTM (Hierarchical Temporal Memory) is an online, unsupervised learning algorithm pioneered by the folks over at [Numenta](http://numenta.org/). This is not an up to spec implementation; gohtm diverges somewhat from the original.

My goals for this project are geared more towards experimentation and generalized machine intelligence development, rather than anomaly detection and vector prediction. I'm experimenting with some different approaches toward the storage and manipulation of synaptic connectivity networks, primarily to the ultimate goal of efficient scaling across multiple machines with large(r) networks of HTM regions, as well as coordination amongst them.

Further, I'm devoting attention to vision oriented problems and coordination of full multi-sensory modalities into global ensemble representations.

## Roadmap

- [x] Encoder Base
- [x] Scalar Encoder
- [x] Random Distributed Scalar Encoder
- [x] Spatial Pooler
- [ ] Temporal Memory
- [ ] Temporal Pooler
- [ ] Classifier
- [ ] Tests
- [ ] Visualization

## Experiments & research directions

### Async processing

If we can simulate each neuron in its own goroutine, this should simplify overall algorithmic processing. With fixed sparsity (0.02) representations, a reactive, event driven processing model should radically lower the CPU cost of iterating the algorithm - as each neuron will only be process inputs/outputs if it is signalled by some other neuron.

In this model, each cell has a single broadcast channel, which other cells can connect to to receive messages. If synaptic weight is high enough to be connected, messages are successfully processed by the postsynaptic cell.

### Vision / Audio
- [ ] Figure out how to code for intensity

pixel (r, g, b) --> [00000111]

hierarchy layers
Edge detect --> oriented edge detect --> 


### Networks
- [ ] Spec out a network definition language. Code generation? 
