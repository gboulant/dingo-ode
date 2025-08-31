// Package system is a wrapper around the solver package (which is consider as
// the core package) to make it easier to setup and execute a solving problem in
// most standard cases.
package system

/*
The package system provides (1) an interface named
System that can be used to modelise the dynamical system and (2) a solver
SystemSolver that helps you with the solving of a System, given an initial
condition, the time step and the time max. A SystemSolver uses the core
package solver with a predefined setup (solver method of type RK4, controller
of type StopAtTime(tmax), recorder of type TimeSeries). The SystemSolver
provides you with helper fiunctions to postprocess the resulting timeseries,
either by dumping the timeseries into a CSV file, or by ploting the
associated curves (this last feature requires an external python library
matplotlib).

Of course, the interface System and the wrapper SystemSolver are optional to
use the ode solver, but they could help to process the standard use cases.

In addition, and to illustrate the feature of ode, the package contains
predefined systems of well-known academic dynamical systems. Some
implementations are based on the core API (spring, lorenz, laser), and other
are implemented using the System interface (watertank, volterra).
*/

import (
	"fmt"
	"log"

	"github.com/gboulant/dingo-ode/solver"
)

// System defines an interface to manipulate a dynamical system governed by an
// ordinary differential equation.
type System interface {
	// F should implement the rate fonction of the dynamical system: dX/dt = F(X,t)
	F(t float64, X []float64) ([]float64, error)
	// GetDefaultInput should returns a default set of input parameters for the Solve function
	GetDefaultInput() (t0 float64, X0 []float64, step float64, tmax float64)
}

// SystemSolver is a tool that helps the setup and execution of a diego solver
// on a specified System. It solves the initial value problem with a stop
// condition of type stopAtTime, and with a timeseries recorder.
type SystemSolver struct {
	system   System
	solver   solver.Solver
	recorder solver.RecorderTimeSeries
}

// NewSystemSolver creates an instance of a SystemSolver for the specified
// System. The SystemSolver is a wrapper around the core diego solver, with a
// predefined setup (solver method of type RK4, controller of type
// StopAtTime(tmax), recorder of type TimeSeries).
func NewSystemSolver(system System) SystemSolver {
	s := SystemSolver{
		system: system,
		solver: solver.NewRK4Solver(),
	}
	return s
}

// Solve executes the Solve function of the solver of the SystemSolver
func (s *SystemSolver) Solve(t0 float64, X0 []float64, h, tmax float64) error {
	controller := solver.StopAtTime(tmax)
	n, err := s.solver.Solve(s.system.F, t0, X0, h, controller, &s.recorder)
	log.Printf("DBG: number of iterations: %d\n", n)
	return err
}

// Series returns a pointer to the current timeseries recorded by the recorder
// of the SystemSolver.
func (s SystemSolver) Series() *solver.TimeSeries {
	return &(s.recorder.Series)
}

// SaveTimeseries saves the current timeseries in the specified csv file with
// the specified names for data columns
func (s SystemSolver) SaveTimeseries(csvpath string, names []string) error {
	return s.Series().ToCSVwithNames(csvpath, names)
}

// PlotTimeSeries plots the current timeseries and assign the specified names to
// the curves. WARN: this function uses an external python library (matplotlib)
// that should be installed on your system.
func (s SystemSolver) PlotTimeSeries(names []string, multi bool) error {
	csvpath := "/tmp/diegodata.csv"
	s.SaveTimeseries(csvpath, names)

	pynames := pystring(names)
	pymulti := pybool(multi)
	plotter := NewPlotter()
	lines := []string{
		fmt.Sprintf("plot.timeseries(csvpath='%s',names=%s,multi=%s)", csvpath, pynames, pymulti),
	}

	scriptpath := "/tmp/outplot.py"
	err := plotter.Create(scriptpath, lines)
	if err != nil {
		return err
	}
	return plotter.Execute(scriptpath)
}
