/*
Package region provides an implementation agnostic
interface for building and using gohtm regions.
*/
package region

// Region ...
type Region interface {
}

/* Sensor region
inputs : topDown | sensor datagram
outputs: bottomUp | probability distribution
*/
type Input struct {
}

/* Bi-directional region
inputs : bottomUp | topDown
outputs: bottomUp | topDown
*/
type Bidir struct {
}

/* Uni-directional region
inputs : bottomUp
outputs: topDown
*/
type Unidir struct {
}

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
*/
