package region

/* Sensor region
inputs : topDown | sensor datagram
outputs: bottomUp | probability distribution
*/
type Sensor struct {
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
