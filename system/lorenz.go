package system

import (
	"fmt"
	"log"

	"github.com/gboulant/dingo-ode/solver"
)

/*
Lorenz system:

 x' = sigma (y-x)
 y' = x (rho - z) - y
 z' = xy - beta z
*/

// LorenzSystem defines the dynamical system modelling the damped spring
type LorenzSystem struct {
	rho   float64
	sigma float64
	beta  float64
}

// f implements the function f of the spring system (in dX/dt = f(X,t))
func (dynsys LorenzSystem) f(t float64, X []float64) ([]float64, error) {
	x := X[0]
	y := X[1]
	z := X[2]
	dx := dynsys.sigma * (y - x)
	dy := x*(dynsys.rho-z) - y
	dz := x*y - dynsys.beta*z
	return []float64{dx, dy, dz}, nil
}

// DemoLorenz illustrates a system exhibiting a chaotic behavior. The orbit of
// the system can be drawn in the 3D phase space to display the Lorenz
// attractor. This example use the RK4 solver.
func DemoLorenz(postpro bool) error {
	dynsys := LorenzSystem{
		rho:   28.0,
		sigma: 10.0,
		beta:  8.0 / 3.0,
	}

	x := 1.0
	y := 1.0
	z := 1.0
	X0 := []float64{x, y, z}
	t0 := 0.0
	h := 0.01
	tmax := 100.0

	algo := solver.NewRK4Solver()

	var recorder solver.RecorderTimeSeries
	n, err := algo.Solve(dynsys.f, t0, X0, h, solver.StopAtTime(tmax), &recorder)
	if err != nil {
		return err
	}
	log.Printf("Problem solved in %d iterations\n", n)
	t, X := algo.Result()
	x = X[0]
	y = X[1]
	z = X[2]
	log.Printf("t: %.2f, x: %.4f, y: %.4f, z: %.4f\n", t, x, y, z)

	// Postprocessing the result
	timeseries := recorder.Series
	csvpath := "out.lorenz_data.csv"
	timeseries.ToCSVwithNames(csvpath, []string{"x", "y", "z"})

	plotter := NewPlotter()
	lines := []string{
		fmt.Sprintf("plot.diagram3D('%s','x','y','z')", csvpath),
	}
	scriptpath := "out.lorenz_plot.py"
	err = plotter.Create(scriptpath, lines)
	if postpro {
		err = plotter.Execute(scriptpath)
	}

	return err
}
