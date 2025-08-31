=========
Userguide
=========

Quick start
===========

Consider the differential equation that modelizes a damped spring:

.. math::

   m*\ddot{x} = -k*x - a*v

This equation can be rewriten to a 1-degree differential system:

.. math::

   \begin{array}{lcl}
   \dot{x} & = & v \\
   \dot{v} & = & -x*k/m - v*a/m
   \end{array}
   
The general form of the system is :math:`\dot{X} = f(X,t)` where
:math:`X` is a vector :math:`(x,v)` that represents the state of the
system and :math:`f(X,t)` a function that represents the evolution
rate of the system.

A numerical solution of this equation can be found by integrating the
system from an initial condition :math:`t0`, :math:`X_0=X(t_0)`, using
an explicit method that compute iteratively :math:`X(t_{n+1})` from
:math:`X(t_n)`. The result at time t is :math:`X(t) = X(t_n) =
X(t_0+n*h)` where :math:`h` is the integration time step.

The package ``solver`` of ``galuma.net/systemd/ode`` provides the core
function to execute this procedure. First of all, you need to import
the package in your go program:

.. code-block:: go

   import "galuma.net/systemd/ode/solver"

Then, you have to provide an implementation of the function
:math:`f(X,t)`. The only requirement is the interface:

.. code-block:: go

   k := 2.0
   m := 1.0
   a := 0.1
   
   f := func(t float64, X []float64) ([]float64, error) {
   		x := X[0]
   		v := X[1]
   		dx := v
   		dv := -x*k/m - v*a/m
   		return []float64{dx, dv}, nil
   }

Define the initial condition :math:`(t_0,X_0)` (variables ``t0`` and
``X0``, the integration step ``h`` and some additional parameters to
specify the stopping condition (a time limit ``tmax`` in this
example):

.. code-block:: go

   X0 := []float64{0.5, 0.0}
   t0 := 0.0
   h := 0.01
   tmax := 60.0

You then have to select a solving method and execute the solver with
all this parameters:

.. code-block:: go
   
   rk2 := solver.NewRK2Solver()
   n, err := rk2.Solve(f, t0, X0, h, solver.StopAtTime(tmax), nil)

The variables returned by the function ``Solve`` are ``n`` the number
of iterations and ``err`` the error of execution. If ``err`` is
``nil``, then no error occurs during execution and ``n`` should look
like ``(tmax-t0)/h``. You can finally retrieve the result, i.e. the
value of ``X(t=tmax)`` at ending time t=tmax:

.. code-block:: go

   t, X := rk2.Result()
   x := X[0]
   v := X[1]

Conversely, if ``err`` is not null, then an error occured during the
solving process and it probably stopped before tmax (n<(tmax-t0)/h and
t<tmax).
   
That's all you need to known to start with goode. The following
section give you some details concerning the options and good
practices.

Recording the timeseries
========================



Controlling the process
=======================

Define a stop condition.


Selecting the integration method
================================



Examples of dynamical systems
=============================


