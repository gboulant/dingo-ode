package system

import "math"

/*

The equation of Volterra (or Lotka-Volterra) modelizes the evolution of
population of prey and predator:

dx/dt = x(t) * (a - b * y(t))
dy/dt = y(t) * (d * x(t) - g)

Where:

* t is the time
* x(t) is the population of preys (number of individuals)
* y(t) is the population of predators (number of individuals)
* dx/dt and dy/dt are the population rate in time

And the parameters:

* a reproduction rate of preys (independant of the number of predators)
* b death rate of preys du to the predators
* d reproduction rate of predators related to the captured preys
* g death rate of predators (independant of preys)

The equation has two fixed points:

* (x=0, y=0): total extinction of both species
* (x=g/d, y=a/b): equilibrium with no evolution of populations

A classical solution is an oscillating evolution around the non trivial fixed
point (x=g/d, y=a/b). Near the fixed point, the period is 2*PI/sqrt(ab).

see: https://fr.wikipedia.org/wiki/%C3%89quations_de_pr%C3%A9dation_de_Lotka-Volterra

*/

// VolterraSystem modelises a prey/predator system
type VolterraSystem struct {
	a float64 // reproduction rate of preys (independant of the number of predators)
	b float64 // death rate of preys du to the predators
	d float64 // reproduction rate of predators related to the captured preys
	g float64 // death rate of predators (independant of preys)
}

// F implements the function f of the volterra system
func (system VolterraSystem) F(t float64, X []float64) ([]float64, error) {
	x := X[0]
	y := X[1]
	dxdt := x * (system.a - system.b*y)
	dydt := y * (system.d*x - system.g)
	return []float64{dxdt, dydt}, nil
}

func (system VolterraSystem) GetDefaultInput() (t0 float64, X0 []float64, step float64, tmax float64) {
	xe := system.g / system.d
	ye := system.a / system.b
	T := 2. * math.Pi / math.Sqrt(system.a*system.b) // pseudo period

	t0 = 0.0
	X0 = []float64{0.8 * xe, 1.2 * ye}
	step = T / 40.
	tmax = 4. * T
	return
}

// DemoVolterra simulates the water tank fill in/out
func DemoVolterra(postpro bool) error {
	system := VolterraSystem{a: 2. / 3., b: 4. / 3., d: 1., g: 1.}
	syssolver := NewSystemSolver(system)
	t0, X0, step, tmax := system.GetDefaultInput()
	err := syssolver.Solve(t0, X0, step, tmax)
	if err != nil {
		return err
	}

	names := []string{"x", "y"}
	if postpro {
		err = syssolver.PlotTimeSeries(names, true)
	} else {
		err = syssolver.SaveTimeseries("out.volterra_data.csv", names)
	}
	return err
}
