package system

import (
	"fmt"
	"math/rand"
)

/*

The water tank system modelizes a tank filled in with a constant input flow and
getting empty because of a small hole on the bottom surface. The out flow
is proportional to the water pressure on the hole, i.e. a term proportional to
the water height, say a*h. The input flow is constant and equal to d. The
equation that governs the time evolution of the water height h is then something
like:

	dh/dt = d - a*h

The analytical solution for an initial condition like (t0=0,h=h0) is:

	h = d/a + (h0 - d/a)*exp(-a*t)

(find a solution like h=A+B*exp(r*t), i.e. the sum of a particular solution of
the general equation and the general solution of the equation without the second
term d)

The evolution fit an exponential curve starting from h0 and decreasing
(increasing) towards the equilibrium height he=d/a if h0>d/a (h0<d/a) with a
typical time rate of T=1/a.

For a relevant plot, we could choose the parameters:

h0 = 1 + d/a
t0 = 0
tmax = 8 * T (where T = 1/a)
step = T/40

*/

// WaterTankSystem modelises a single water tank.
type WaterTankSystem struct {
	d, a float64
}

// F implements the function f of the watertank system (in dX/dt = f(X,t))
func (system WaterTankSystem) F(t float64, X []float64) ([]float64, error) {
	h := X[0]
	dhdt := system.d - system.a*h
	return []float64{dhdt}, nil
}

func (system WaterTankSystem) GetDefaultInput() (t0 float64, X0 []float64, step float64, tmax float64) {
	t0 = 0.0
	h0 := 1 + system.d/system.a
	X0 = []float64{h0}

	T := 1 / system.a
	step = T / 40
	tmax = 8 * T
	return
}

// DemoWaterTank simulates the water tank fill in/out
func DemoWaterTank(postpro bool) error {
	watertank := WaterTankSystem{d: 2, a: 1}
	syssolver := NewSystemSolver(watertank)
	t0, X0, step, tmax := watertank.GetDefaultInput()
	err := syssolver.Solve(t0, X0, step, tmax)
	if err != nil {
		return err
	}

	if postpro {
		err = syssolver.PlotTimeSeries([]string{"h"}, false)
	} else {
		err = syssolver.SaveTimeseries("out.watertank_data.csv", []string{"h"})
	}
	return err
}

// ----------------------------------------------------------------------------

// CascadingWaterTankSystem modelises a cascade of water tank, where the tank0
// is filled in by a constant input flow, the tank 01 is filled in by the output
// flow of the tank 0, the tank 02 is filled in by the output flow of the tank
// 1, etc.
type CascadingWaterTankSystem struct {
	d, a float64
	n    int
}

// F implements the function f of the chained watertank system
func (system CascadingWaterTankSystem) F(t float64, X []float64) ([]float64, error) {
	dXdt := make([]float64, len(X))
	dXdt[0] = system.d - system.a*X[0]
	for i := 1; i < len(X); i++ {
		dXdt[i] = system.a * (X[i-1] - X[i])
	}
	return dXdt, nil
}

func (system CascadingWaterTankSystem) GetDefaultInput() (t0 float64, X0 []float64, step float64, tmax float64) {
	t0 = 0.0
	//X0 := NewConstantX(watertank.n, 1+watertank.d/watertank.a)
	X0 = NewRandomX(system.n, system.d/system.a, 0.2*system.d/system.a)

	T := 1 / system.a
	step = T / 40
	tmax = T * 2. * float64(system.n)
	return
}

// NewConstantX creates a default X value of size n and with all component
// values equal to h.
func NewConstantX(n int, h float64) []float64 {
	X := make([]float64, n)
	for i := 0; i < n; i++ {
		X[i] = h
	}
	return X
}

func NewRandomX(n int, h float64, dh float64) []float64 {
	X := make([]float64, n)
	for i := 0; i < n; i++ {
		X[i] = h + dh*(rand.Float64()-0.5)
	}
	return X
}

// DemoCascadingWaterTank simulates the cascading water tank fill in/out
func DemoCascadingWaterTank(postpro bool) error {
	watertank := CascadingWaterTankSystem{d: 2, a: 1, n: 8}
	syssolver := NewSystemSolver(watertank)
	t0, X0, step, tmax := watertank.GetDefaultInput()
	err := syssolver.Solve(t0, X0, step, tmax)
	if err != nil {
		return err
	}

	names := make([]string, len(X0))
	for i := 0; i < len(names); i++ {
		names[i] = fmt.Sprintf("h%d", i)
	}
	if postpro {
		err = syssolver.PlotTimeSeries(names, false)
	} else {
		err = syssolver.SaveTimeseries("out.watertank_data.csv", names)
	}
	return err
}
