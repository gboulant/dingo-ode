package solver

func rk4Iteration(f Function, tn float64, Xn []float64, h float64) ([]float64, error) {
	// Step 1: k1 = h*f(tn, Xn)
	slope, err := f(tn, Xn)
	if err != nil {
		return nil, err
	}
	k1 := make([]float64, len(Xn))
	for i := 0; i < len(Xn); i++ {
		k1[i] = h * slope[i]
	}

	// Step 2: k2 = h*f(tn+h/2, Xn+k1/2). We define Xm as the mediate point Xn+k1/2.
	Xm := make([]float64, len(Xn))
	for i := 0; i < len(Xn); i++ {
		Xm[i] = Xn[i] + k1[i]/2
	}
	slope, err = f(tn+h/2, Xm)
	if err != nil {
		return nil, err
	}
	k2 := make([]float64, len(Xn))
	for i := 0; i < len(Xn); i++ {
		k2[i] = h * slope[i]
	}

	// Step 3: k3 = h*f(tn+h/2, Xn+k2/2). We define Xm as the mediate point Xn+k2/2.
	for i := 0; i < len(Xn); i++ {
		Xm[i] = Xn[i] + k2[i]/2
	}
	slope, err = f(tn+h/2, Xm)
	if err != nil {
		return nil, err
	}
	k3 := make([]float64, len(Xn))
	for i := 0; i < len(Xn); i++ {
		k3[i] = h * slope[i]
	}

	// Step 4: k4 = h*f(tn+h, Xn+k3). We define Xm as the mediate point Xn+k3.
	for i := 0; i < len(Xn); i++ {
		Xm[i] = Xn[i] + k3[i]
	}
	slope, err = f(tn+h, Xm)
	if err != nil {
		return nil, err
	}
	k4 := make([]float64, len(Xn))
	for i := 0; i < len(Xn); i++ {
		k4[i] = h * slope[i]
	}

	// Computing the weigth average final value Xn+1 (denoted to as Xs below)
	Xs := make([]float64, len(Xn))
	for i := 0; i < len(Xn); i++ {
		Xs[i] = Xn[i] + (k1[i]+2*k2[i]+2*k3[i]+k4[i])/6
	}

	return Xs, nil
}

// NewRK4Solver returns a Solver that implements the Euler algorithm
func NewRK4Solver() Solver {
	solver := StandardSolver{iteration: rk4Iteration}
	return &solver
}
