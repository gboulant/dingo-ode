package solver

func eulerIteration(f Function, tn float64, Xn []float64, h float64) ([]float64, error) {
	Xs := make([]float64, len(Xn))
	slope, err := f(tn, Xn)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(Xn); i++ {
		Xs[i] = Xn[i] + h*slope[i]
	}
	return Xs, nil
}

// NewEulerSolver returns a Solver that implements the Euler algorithm
func NewEulerSolver() Solver {
	solver := StandardSolver{iteration: eulerIteration}
	return &solver
}
