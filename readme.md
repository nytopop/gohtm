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
- [x] Scalar encoder
- [ ] Audio encoder
- [ ] Vision encoder
- [ ] Random distributed scalar encoder
- [x] Spatial Pooler
- [x] Temporal Memory
- [ ] Temporal Pooler
- [ ] Classifier
- [ ] Tests
- [ ] Visualization

## Experiments & research directions

### Networks
- [ ] Spec out a network definition language. Code generation? 
- [ ] First in Last out stack for processing

### Temporal memory
- [ ] get some benchmark sequences for testing prediction accuracy, etc
- [ ] figure out what to do when we hit the limit on cellular objects. More recent data is preferable, online learning and all...
- [ ] fix the anomaly calculation

## HTM Components
### Encoder
An encoder converts some quantity of sensory information into a semi-sparse vectorized representation for processing in HTM. Encoders should filter out any irrelevant information and only include relevant semantic information about the input in their output vector.

### Spatial Pooler
A spatial pooler converts a semi-sparse vector representation of some quantity of sensory information into a fixed sparsity vector. Every bit in the output vector should correspond to a column of cells in temporal memory.

### Temporal Memory
Temporal Memory learns variable order sequences.

### Temporal Pooler
### Classifier
