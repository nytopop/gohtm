package main

/* Temporal Memory */

type TemporalParams struct {
	numColumns int
}

func NewTemporalParams() TemporalParams {
	return TemporalParams{
		numColumns: 2048,
	}
}

type TMColumn struct {
}

type TemporalMemory struct {
	// state
	cols []TMColumn

	// params
	numColumns int
}

func NewTemporalMemory(p TemporalParams) TemporalMemory {
	tm := TemporalMemory{
		numColumns: p.numColumns,
	}

	tm.cols = make([]TMColumn, p.numColumns)

	return tm
}
