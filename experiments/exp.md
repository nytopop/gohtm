# Info
This file contains ideas for in-progress and upcoming research directions.

# Feedback impl
If we have two regions r1 & r2, each learning a character stream, and a region r3 learning activity of r1 & r2.

r1 [a, b, c, d]
r2 [e, f, g, h]
r3 [a:e, b:f, c:g, d:h]

Feedback should allow for r3 to disambiguate sequences like the following:

r1: [a, b, c, d]
r2: [e, f, g, h]
-reset-
r1: [a, b, x, z]
r2: [t, g, z, y]
-reset-

When we receive:

r1: [a, b]
r2: [e, f]

r3 should bias r1 into predicting [c, d] rather than [x, z]

Likewise, when we receive:

r1: [a, b]
r2: [t, g]

r3 should bias r1 into predicting [x, z] rather than [c, d]

# Zero pooler test
Feed encoded vector directly into temporal memory.
Analyze predictive performance, perf metrics.

# TM performance benchmark
Analyze decay rate of iteration speed.
Analyze correlations between col count, cell count, synapse count
CPU profile.

# Sensory-motor test / impl
Waiting until extended temporal memory is implemented
Build a simulated environment like:

|-------|
|   |   |
|   O   |
|     x |
|-------|

- current rotation (body centric)
- sensor data (object centric)
- motor command

# JSON region specs
test JSON format region specification with a constructor.

# Region Serialization
Train a region,
serialize all stateful components to bytestream with gob

- encoder
- spatial pooler
- temporal memory
- temporal pooler
- classifier

# Sensory-motor action integration
- classification -> prediction -> optimization
- l0 : universal learner (sp)
- l1 : universal learning predictor (sp + tm)
- l2 : universal learning predictive optimizer (sp + tm + sm)
- predict an action that leads to x result
fn(s sensor) -> prediction
fn(x result, p prediction) -> action

# Spatio-Temporal scoping
We can use region hierarchies to realize two key properties in an htm network shaped as a binary tree: spatial and temporal scoping.

k : depth
s : spatial scope
t : temporal scope

With a scoping rate of 2^k, the _spatial scope_ of any region is { 2^k }. E.g., if k = 2, that region's input pooling process will isolate features distributed across 4 inputs on the network.

The spatial abstraction level of a region can be given as k; each level of the hierarchy isolates features in the next lower level's pooled representation. This property persists with any network topology.

With a scoping rate of 2^k, the _temporal scope_ of any region is { 2^k }. E.g., if k = 2, that region will only process inputs where { i % 4 = 0 }. In the interims where no inputs are processed, the feedback from the last processed input is persisted.
