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
