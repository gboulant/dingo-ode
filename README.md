# Solving Dynamical Systems with go

**contact**: [Guillaume Boulant](mailto:gboulant@gmail.com?subject=dingo-ode)

The package `ode` can help you with solving Ordinary Differential
Equations (ODE) with explicit methods (Euler, RK2, RK4). It is
restricted to initial value problems, i.e. problems defined by a
differential equation completed with an initial condition:

* An **initial value problems** is defined by a differential equation
  dX/dt=f(X,t) completed with a particular initial condition t=t0,
  X0=X(t0).
* An **explicit method** is an iterative method where X at step i+1 is
  determined from X at step i.

The variable X is a vector of real numbers whose values evolve with a
real scalar parameter t. The parameter t represents the time in most
of cases, for example when the differential equation modelizes a
dynamical systems. In such a case, X(t) characterises the state of the
system at the time t.

The integration of the package `ode` in an application consists in (1)
implement the function f(X,t), (2) define the initial condition (t0,
X0=X(t0)), (3) select a solving method (Euler, RK2 or RK4) and (4)
execute the selected solver with all this parameters, completed with a
stopping conditions (e.g. a time limit tmax).

Have a look to the documenation:

* [Quick start guide](admin/doc/userguide.rst)

And the demonstrative examples in the folder [demos](demos).
