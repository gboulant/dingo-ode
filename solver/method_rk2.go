package solver

func rk2Iteration(f Function, tn float64, Xn []float64, h float64) ([]float64, error) {
	// Step 1
	Xm := make([]float64, len(Xn))
	slope, err := f(tn, Xn)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(Xn); i++ {
		Xm[i] = Xn[i] + h*slope[i]/2
	}

	// Step 2
	Xs := make([]float64, len(Xn))
	slope, err = f(tn+h/2, Xm)
	if err != nil {
		return nil, err
	}
	for i := 0; i < len(Xn); i++ {
		Xs[i] = Xn[i] + h*slope[i]
	}

	return Xs, nil
}

// NewRK2Solver returns a Solver that implements the Euler algorithm
func NewRK2Solver() Solver {
	solver := StandardSolver{iteration: rk2Iteration}
	return &solver
}
