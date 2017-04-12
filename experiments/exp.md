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

Deps: extended TM

# Zero pooler test
Feed encoded vector directly into temporal memory.
Analyze predictive performance, perf metrics.

# TM performance benchmark
Analyze decay rate of iteration speed.
Analyze correlations between col count, cell count, synapse count
CPU profile.

# Sensory-motor test / impl
Build a simulated environment like:

|-------|
|   |   |
|   O   |
|     x |
|-------|

- current rotation (body centric)
- sensor data (object centric)
- motor command

Deps: extended TM

# JSON region specs
test JSON format region specification with a constructor.

# Region Serialization
Train a region,
serialize all stateful components to bytestream with JSON

- [x] encoder
- [x] spatial pooler
- [x] temporal memory
- [ ] region
- [ ] network
- [ ] classifier

# Sensory-motor action integration
- classification -> prediction -> optimization
- l0 : universal learner (sp)
- l1 : universal learning predictor (sp + tm)
- l2 : universal learning predictive optimizer (sp + tm + sm)
- predict an action that leads to x result
fn(s sensor) -> prediction
fn(x result, p prediction) -> action

// generate a state prediction
fn(s sensor) -> prediction

// generate a desired transform; we have state and goal
fn(g goal, p prediction) -> transform
fn() -> motor

current aState; desired bState
tform = transform(aState -> bState)
action = action(tform)

Deps: extended TM
