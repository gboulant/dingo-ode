package system

import (
	"fmt"
	"math"

	"github.com/gboulant/dingo-ode/solver"
)

/*
The laser dynamics can exhibit chaos. We use here the rate equations
model that describes the ligth intensity L and the population inversion
D:

 L' = D-1-m*cos(w*t)
 D' = g*(a-D*(1+exp(L)))

We set Z=w*t to have an autonomous system:

 L' = D-1-m*cos(Z)
 D' = g*(a-D*(1+exp(L)))
 Z' = w

*/

// LaserSystem defines the dynamical system modelling the laser in the
// approximation of the rates equations balancing the light intensity L with the
// population inversion D.
type LaserSystem struct {
	m, g, w, a float64
}

// f implements the function f of the laser system (in dX/dt = f(X,t))
func (dynsys LaserSystem) f(t float64, X []float64) ([]float64, error) {
	L := X[0]
	D := X[1]
	Z := X[2]

	dL := D - 1 - dynsys.m*math.Cos(Z)
	dD := dynsys.g * (dynsys.a - D*(1+math.Exp(L)))
	dZ := dynsys.w
	return []float64{dL, dD, dZ}, nil
}

var configurations = map[string]LaserSystem{
	"chaos": {
		m: 2.5e-2,
		g: 1e-3,
		w: 1e-2,
		a: 1.1,
	},
	"1T_LONG_TRANSIENT": {
		m: 2.5e-2,
		g: 1e-3,
		w: 1e-1,
		a: 1.1,
	},
}

// DemoLaser simulates the evolution of a laser light intensity
func DemoLaser(postpro bool) error {

	dynsys := configurations["chaos"]

	L0 := 1.0
	D0 := 1.0
	Z0 := 0.0
	t0 := 0.0
	X0 := []float64{L0, D0, Z0}

	T := 2 * math.Pi / dynsys.w
	tmax := 60 * T
	h := T / 40

	algo := solver.NewRK4Solver()
	var recorder solver.RecorderTimeSeries
	_, err := algo.Solve(dynsys.f, t0, X0, h, solver.StopAtTime(tmax), &recorder)
	if err != nil {
		return err
	}

	// Postprocessing the result
	timeseries := recorder.Series
	csvpath := "out.laser01_data.csv"
	timeseries.ToCSVwithNames(csvpath, []string{"L", "D", "Z"})

	plotter := NewPlotter()
	lines := []string{
		fmt.Sprintf("plot.timeseries(csvpath='%s',names=['L','D'],multi=True)", csvpath),
	}
	scriptpath := "out.laser01_plot.py"
	err = plotter.Create(scriptpath, lines)
	if postpro {
		err = plotter.Execute(scriptpath)
	}

	return err
}

// DemoLaserFirstReturnMap draw a first return map of the chaotic dynamic
func DemoLaserFirstReturnMap(postpro bool) error {

	dynsys := configurations["chaos"]

	L0 := 1.0
	D0 := 1.0
	Z0 := 0.0
	t0 := 0.0
	X0 := []float64{L0, D0, Z0}

	T := 2 * math.Pi / dynsys.w
	h := T / 40

	// We first let the transient behaviour stabilized on the attractor
	algo := solver.NewRK4Solver()
	_, err := algo.Solve(dynsys.f, t0, X0, h, solver.StopAtTime(100*T), nil)
	if err != nil {
		return err
	}

	// We consider the end state as the starting state for the following iteration
	_, X := algo.Result()
	var timeseries solver.TimeSeries
	timeseries.Append(solver.NewTimeData(0, X))

	for i := 0; i < 10000; i++ {
		_, err := algo.Solve(dynsys.f, 0, X, h, solver.StopAtTime(T), nil)
		if err != nil {
			return err
		}
		_, X = algo.Result()
		timeseries.Append(solver.NewTimeData(float64(i+1), X))
	}

	// Postprocessing the result
	csvpath := "out.laser02_data.csv"
	timeseries.ToCSVwithNames(csvpath, []string{"L", "D", "Z"})

	plotter := NewPlotter()
	lines := []string{
		fmt.Sprintf("plot.firstReturnMap(csvpath='%s',xname='L',yname='D')", csvpath),
	}
	scriptpath := "out.laser02_plot.py"
	err = plotter.Create(scriptpath, lines)
	if postpro {
		err = plotter.Execute(scriptpath)
	}

	return err
}
