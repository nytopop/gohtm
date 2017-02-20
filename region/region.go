// Package region provides an implementation agnostic interface for linking
// and interacting with computational units of the htm algorithm.
package region

/* Region Types

SensorRegion | Bottom level of hierarchy, com with higher regions
  Algs: encoder, classifier, sp, tm, tp (?)
  Inputs : sensor data, topDown
  Outputs: bottomUp

MidRegion | Mid levels of hierarchy, com with lower & higher regions
  Algs: sp, tm, tp (?)
  Inputs : bottomUp, topDown
  Outputs: buttomUp, topDown

TopRegion | Top level of hierarchy, com with lower regions
  Algs: sp, tm, tp (?)
  Inputs: bottomUp
  Outputs: topDown

TODO
  MultiEncoder
  Classifier
  Temporal Pooler
  Temporal Memory w/ apical dendrites
*/
