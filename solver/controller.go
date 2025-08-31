package solver

import (
	"math"
)

// Controller defines a function that handles the execution of the solving
// process. A controller is a function that decides when the process should
// stop. A solving process can stop in normal condition (e.g. when time reach a
// given threshold) or because an abnormal behavior occurs (e.g. divergence of
// the solution outside of an predefined X domain). The Controller is called by the solver
// at each iteration to know if the solving process should stop. To indicate a
// normal stop, the Controller should return (bool,error)=(true,nil), bool=true
// meaning "yes, stop the process". To indicate an abnormal condition, the
// controller should return (bool,error)=(_,error), error specifying the details
// concerning the abnormal condition. In this case, the solving process stops
// returning this error whatever the bool value is.
type Controller func(t float64, X []float64) (bool, error)

// StopAtTime implements a default stop Controller that stops the solving process
// when the time exceed the maximum value tmax. In this implementation, the time
// tmax is included and the process stops, i.e. the time tmax is considered as a
// good value, and then it will be the result time of the Solve function.
func StopAtTime(tmax float64) Controller {
	controller := func(t float64, X []float64) (bool, error) {
		delta := math.Abs(t - tmax)
		if delta < 1e-8 || t <= tmax {
			return false, nil
		}
		return true, nil
	}
	return controller
}

// MultiController creates a controller that aggregates a list of controllers.
// The created controller executes the input controllers in the specified order
// and return true if a controller return true (meaning that the solving process
// should be stopped).
func MultiController(controllers ...Controller) Controller {
	return func(t float64, X []float64) (bool, error) {
		for _, c := range controllers {
			ok, err := c(t, X)
			if err != nil {
				return true, err
			}
			if ok {
				return true, nil
			}
		}
		return false, nil
	}
}
