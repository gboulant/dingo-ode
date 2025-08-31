package solver

import (
	"errors"
)

// Function defines the function F of an ODE system, i.e. the function that
// implements the time derivative of the state of the system:
//
// dX/dt = F(t,X).
//
// Note that even if it represents a mathematic concept, it should be able to
// raise an error to inform for example that a forbiden math operation is
// realized (division by zero, square root of a negative number, etc).
type Function func(t float64, X []float64) (dXdt []float64, err error)

// Solver is the interface to be implemented by the ODE solvers. An ODE solver
// tries to solve the problem defined by (1) a 1-degree dynamic system dX/dt=f(X,t)
// where X is a set of variables that characterize the state of the system
// (coordinates in the phase diagram) and (2) the initial conditions defined by
// (t0, X0).
type Solver interface {
	// Solve solves the system defined by the Function f, from initial
	// conditions (t0,X0), with a step size of h, and stopping the process when
	// the stop handler return true. The Solve function returns the number of
	// iterations and a non nil error if that occurs.
	Solve(f Function, t0 float64, X0 []float64, h float64, c Controller, r Recorder) (uint64, error)
	// Result returns the values of t and X obtained at the end of the solving process
	Result() (t float64, X []float64)
}

// Iteration defines a function that implements an iteration step of a standard solver
type Iteration func(f Function, tn float64, Xn []float64, h float64) ([]float64, error)

// StandardSolver implements the interface Solver with a standard Solve
// implementation. The StandardSolver apply an Iteration function (to be
// defined) at each step of the Solve function.
type StandardSolver struct {
	t         float64
	X         []float64
	iteration Iteration
}

// Solve implements the Solver interface for the StandarSolver
func (solver *StandardSolver) Solve(f Function, t0 float64, X0 []float64, h float64, c Controller, r Recorder) (uint64, error) {
	if f == nil {
		return 0, errors.New("ERR: the function f is not defined")
	}
	if r == nil {
		r = &RecorderNone{} // Record no intermediate iteration
	}

	tm := t0
	Xm := X0
	r.Record(tm, Xm)

	var nbIterations uint64 = 0

	for {
		Xn, err := solver.iteration(f, tm, Xm, h)
		if err != nil {
			return nbIterations, err
		}
		tn := tm + h
		r.Record(tn, Xn)

		stop, err := c(tn, Xn)
		if err != nil {
			return nbIterations, err
		}
		if stop {
			break
		}

		Xm = Xn
		tm = tn
		nbIterations++
	}

	solver.t = tm
	solver.X = Xm

	return nbIterations, nil
}

// Result implements the Solver interface
func (solver *StandardSolver) Result() (t float64, X []float64) {
	return solver.t, solver.X
}
