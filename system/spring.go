package system

/*

The spring system is a simple example of the harmonic oscillator model.
The simplest model (without damping and no gravity) is:

 m*x'' = -k*x

Where x is the position of the mass m by respect to the equilibrium position
(position held when the force is null)

The equation can be rewriten as a one-degree 2D system:

 x' = v
 v' = -x*k/m

If we consider a damping proportionnal to the velocity (viscous damping), the
system becomes:

 m*x'' = -k*x - a*v

 Which can be rewriten to:

 x' = v
 v' = -x*k/m - v*a/m

where:

- (x,v) defines the state of the system (position from equilibrium and velocity)
- k is the spring stiffness
- m is the mass attached to the spring
- a is the damping rate

The general form of this equation is dX/dt = f(X,t) where X=(x,v) and f is
defined by:

 fx = dx/dt = v
 fv = dv/dt = -x*k/m - v*a/m.

The analytic resolution of this equation can be found at:

http://res-nlp.univ-lemans.fr/NLP_C_M01_G04/co/module_NLP_C_M01_G04_3.html

*/

import (
	"fmt"
	"log"
	"math"

	"github.com/gboulant/dingo-ode/solver"
)

// -------------------------------------------------------------------
// DEMO01: illustrates the basic usage of the package ode/solver.

// DemoSpring01 illustrates the basic usage of the package ode/solver. The basic
// usage consists in (1) defining the system by implementing the function f (in
// dX/dt=f(X,t)), (2) specifying the initial conditions t0 and X0 (x0,v0 in the
// spring example), (3) specifying the solving conditions (step size h and stop
// condition), for finally (4) select a solver (Euler, RK2, RK4) and run the
// Solve function with the parameters (1), (2) and (3). This example uses the
// Euler method for solver.
func DemoSpring01(postpro bool) error {

	var k, m, a float64

	k = 2.0
	m = 1.0
	a = 0.1

	f := func(t float64, X []float64) ([]float64, error) {
		x := X[0]
		v := X[1]
		dx := v
		dv := -x*k/m - v*a/m
		return []float64{dx, dv}, nil
	}

	x := 0.5
	v := 0.0
	X0 := []float64{x, v}
	t0 := 0.0
	h := 0.01 // The euler method requires a finer step than RK methods
	tmax := 60.000

	algo := solver.NewEulerSolver()
	var recorder solver.RecorderTimeSeries
	n, err := algo.Solve(f, t0, X0, h, solver.StopAtTime(tmax), &recorder)
	if err != nil {
		return err
	}
	log.Printf("Problem solved in %d iterations\n", n)

	t, X := algo.Result()
	x = X[0]
	v = X[1]
	log.Printf("t: %.2f, x: %.4f, v: %.4f\n", t, x, v)

	// Postprocessing the result
	timeseries := recorder.Series
	csvpath := "out.spring01_data.csv"
	timeseries.ToCSVwithNames(csvpath, []string{"x", "v"})

	plotter := NewPlotter()
	lines := []string{
		fmt.Sprintf("plot.timeseries(csvpath='%s',names=['x','v'])", csvpath),
		fmt.Sprintf("plot.diagram2D(csvpath='%s',xname='x',yname='v')", csvpath),
	}
	scriptpath := "out.spring01_plot.py"
	err = plotter.Create(scriptpath, lines)
	if postpro {
		err = plotter.Execute(scriptpath)
	}

	return err
}

// -------------------------------------------------------------------
// DEMO02: illustrates how to organise the ODE system using a structure

// SpringSystem defines the dynamical system modelling the damped spring
type SpringSystem struct {
	k, m, a float64
}

// f implements the function f of the spring system (in dX/dt = f(X,t))
func (dynsys SpringSystem) f(t float64, X []float64) ([]float64, error) {
	x := X[0]
	v := X[1]
	dx := v
	dv := -x*dynsys.k/dynsys.m - v*dynsys.a/dynsys.m
	return []float64{dx, dv}, nil
}

// DemoSpring02 is a rewrite of DemoSpring01, but using the SpringSystem and the
// RK2 solver method. Note that the RK2 does not require a step size as fine as
// require by the Euler method
func DemoSpring02(postpro bool) error {
	x0 := 0.5
	v0 := 0.0
	X0 := []float64{x0, v0}
	t0 := 0.0
	h := 0.1
	tmax := 60.0

	dynsys := SpringSystem{
		k: 2.0,
		m: 1.0,
		a: 0.1,
	}

	algo := solver.NewRK2Solver()
	var recorder solver.RecorderTimeSeries
	n, err := algo.Solve(dynsys.f, t0, X0, h, solver.StopAtTime(tmax), &recorder)
	if err != nil {
		return err
	}
	log.Printf("Problem solved in %d iterations\n", n)
	t, X := algo.Result()
	x := X[0]
	v := X[1]
	log.Printf("t: %.2f, x: %.4f, v: %.4f\n", t, x, v)

	// Postprocessing the result
	timeseries := recorder.Series
	csvpath := "out.spring02_data.csv"
	timeseries.ToCSVwithNames(csvpath, []string{"x", "v"})

	plotter := NewPlotter()
	lines := []string{
		fmt.Sprintf("plot.timeseries(csvpath='%s',names=['x','v'])", csvpath),
		fmt.Sprintf("plot.diagram2D(csvpath='%s',xname='x',yname='v')", csvpath),
	}
	scriptpath := "out.spring02_plot.py"
	err = plotter.Create(scriptpath, lines)
	if postpro {
		err = plotter.Execute(scriptpath)
	}

	return err
}

// getAnalyticSolution returns a function that represents the analytic solution
// for the position x in the case where the initial conditions are x=x0 and v=0.
func getAnalyticSolution(k, m, a, x0 float64) func(t float64) float64 {
	w0 := math.Sqrt(k / m)
	w := math.Sqrt(w0*w0 - a*a/(4*m*m))
	xanalytic := func(t float64) float64 {
		r := x0 * math.Exp(-a*t/(2*m)) * (math.Cos(w*t) + a*math.Sin(w*t)/(2*m))
		return r
	}
	return xanalytic
}

// DemoSpring03 creates a CSV data file containing the x timeseries computed
// from the analytic solution of the equation in the case where the initial
// conditions are x=x0 and v=0.
func DemoSpring03(postpro bool) error {
	dynsys := SpringSystem{
		k: 1.4,
		m: 1.0,
		a: 0.1,
	}

	x0 := 0.5
	v0 := 0.0
	X0 := []float64{x0, v0}
	t0 := 0.0
	h := 0.01
	tmax := 60.0

	algo := solver.NewEulerSolver()
	var recorder solver.RecorderTimeSeries
	algo.Solve(dynsys.f, t0, X0, h, solver.StopAtTime(tmax), &recorder)
	timeseriesEuler := recorder.Series.Clone()
	timeseriesEuler.ToCSV("out.spring03_simulation_euler.csv")

	algo = solver.NewRK2Solver()
	recorder.Series.Clear()
	algo.Solve(dynsys.f, t0, X0, h, solver.StopAtTime(tmax), &recorder)
	timeseriesRK2 := recorder.Series.Clone()
	timeseriesRK2.ToCSV("out.spring03_simulation_rk2.csv")

	algo = solver.NewRK4Solver()
	recorder.Series.Clear()
	algo.Solve(dynsys.f, t0, X0, h, solver.StopAtTime(tmax), &recorder)
	timeseriesRK4 := recorder.Series.Clone()
	timeseriesRK4.ToCSV("out.spring03_simulation_rk4.csv")

	xanalytic := getAnalyticSolution(dynsys.k, dynsys.m, dynsys.a, x0)
	var atimeseries solver.TimeSeries
	for t := t0; t < tmax+h/2; t += h {
		x := xanalytic(t)
		atimeseries.Append(solver.NewTimeData(t, []float64{x}))
	}
	atimeseries.ToCSV("out.spring03_analytic.csv")

	var mtimeseries solver.TimeSeries
	for i := 0; i < len(atimeseries); i++ {
		ta := atimeseries[i].GetTime()
		xa := atimeseries[i].GetState()[0]

		te := timeseriesEuler[i].GetTime()
		if te != ta {
			return fmt.Errorf("ERR: the euler time (%.4f) is different than the analytic time (%.4f)", te, ta)
		}
		xe := timeseriesEuler[i].GetState()[0]
		de := (xe - xa) * (xe - xa)

		t2 := timeseriesRK2[i].GetTime()
		if t2 != ta {
			return fmt.Errorf("ERR: the RK2 time (%.4f) is different than the analytic time (%.4f)", t2, ta)
		}
		x2 := timeseriesRK2[i].GetState()[0]
		d2 := (x2 - xa) * (x2 - xa)

		t4 := timeseriesRK4[i].GetTime()
		if t4 != ta {
			return fmt.Errorf("ERR: the RK4 time (%.4f) is different than the analytic time (%.4f)", t4, ta)
		}
		x4 := timeseriesRK4[i].GetState()[0]
		d4 := (x4 - xa) * (x4 - xa)

		mtimeseries.Append(solver.NewTimeData(ta, []float64{xa, xe, x2, x4, de, d2, d4}))
	}

	// Postprocessing the result
	names := []string{"xref", "xeuler", "xrk2", "xrk4", "deuler", "drk2", "drk4"}
	csvpath := "out.spring03_data.csv"
	mtimeseries.ToCSVwithNames(csvpath, names)

	plotter := NewPlotter()
	lines := []string{
		fmt.Sprintf("plot.timeseries('%s',['xref', 'xeuler', 'xrk2', 'xrk4'])", csvpath),
		fmt.Sprintf("plot.timeseries('%s',['deuler', 'drk2', 'drk4'])", csvpath),
		fmt.Sprintf("plot.timeseries('%s',['drk2', 'drk4'])", csvpath),
	}
	scriptpath := "out.spring03_plot.py"
	err := plotter.Create(scriptpath, lines)
	if postpro {
		err = plotter.Execute(scriptpath)
	}
	return err
}
