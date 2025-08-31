package main

import (
	"flag"
	"fmt"
	"log"

	"github.com/gboulant/dingo-ode/system"
)

type demofunc func(postpro bool) error

type demodef struct {
	label    string
	function demofunc
	comment  string
}

var demolist = []demodef{
	{"spring01", system.DemoSpring01, "damped spring simulation with basic implementation"},
	{"spring02", system.DemoSpring02, "damped spring simulation with structured implementation"},
	{"spring03", system.DemoSpring03, "damped spring simulation with comparrison to analytical solution"},
	{"lorenz", system.DemoLorenz, "demonstration of the Lorenz chaotic attractor"},
	{"laser01", system.DemoLaser, "demonstration of a chaotic laser dynamics"},
	{"laser02", system.DemoLaserFirstReturnMap, "Poincaré map of a chaotic laser dynamics"},
	{"watertank", system.DemoWaterTank, "water tank fill in/out"},
	{"cwatertank", system.DemoCascadingWaterTank, "cascading water tank fill in/out"},
	{"volterra", system.DemoVolterra, "model of preys/predators populations"},
}

func getDemoFunc(label string) (demofunc, error) {
	for index := range demolist {
		demo := demolist[index]
		if demo.label == label {
			return demo.function, nil
		}
	}
	return nil, fmt.Errorf("ERR: The demo %s does not exist", label)
}

func main() {

	listdemos := flag.Bool("l", false, "list the existing demos")
	demolabel := flag.String("d", "", "name of the demo to execute")
	postpro := flag.Bool("p", false, "plot results after the simulation process (requires matplotlib python library)")

	flag.Parse()

	if *listdemos {
		for index := range demolist {
			demo := demolist[index]
			fmt.Printf("%-10s: %s\n", demo.label, demo.comment)
		}
		return
	}

	if *demolabel != "" {
		function, err := getDemoFunc(*demolabel)
		if err != nil {
			panic(err)
		}
		log.Printf("Execution of the demo %s ...\n", *demolabel)
		err = function(*postpro)
		if err != nil {
			panic(err)
		}
	} else {
		flag.Usage()
	}
}
