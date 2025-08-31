package solver

import "fmt"

// Recorder is the interface to be implemented by the data recorders. A recorder is
// a dataset that can be used by a Solver to record the iteration states of the
// solving process, i.e. the sequence of values (tn,Xn).
type Recorder interface {
	Record(t float64, X []float64)
}

// RecorderNone defines a Recorder that records nothing (default if no recorder
// is specified when calling the Solve function of a Solver, i.e. if recorder=nil)
type RecorderNone struct{}

// Record implements the Recorder interface
func (recorder *RecorderNone) Record(t float64, X []float64) {
	// Do nothing
}

// RecorderLogger defines a Recorder that prints the values on standard output
type RecorderLogger struct{}

// Record implements the Recorder interface
func (recorder *RecorderLogger) Record(t float64, X []float64) {
	fmt.Printf("t: %.2f, X0: %.4f, X1: %.4f\n", t, X[0], X[1])
}

// RecorderTimeSeries defines a Recorder that registers all values in a TimeSeries
type RecorderTimeSeries struct {
	Series TimeSeries
}

// Record implements the Recorder interface so that the TimeSeries can be used
// as a Recorder of a Solve process.
func (recorder *RecorderTimeSeries) Record(t float64, X []float64) {
	data := TimeData{t, X}
	recorder.Series = append(recorder.Series, data)
}
