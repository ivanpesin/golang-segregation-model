package main

import (
	"fmt"
	"strconv"
)

var algdesc = []string{
	"Pick a random available site with higher utility if possible",
	"Pick a random available site",
	"Pick the first available site ordered by vacant time with higher utility",
	"Pick the first available site ordered by vacant time",
	"Pick the first vacant site with higher utility",
}

func cls() {
	fmt.Printf("\033[2J")
}

func pos(r, c int) {
	fmt.Printf("\033[" + strconv.Itoa(r) + ";" + strconv.Itoa(c) + "H")
}

func ind(r, c int) int {
	return r*(config.cols+2) + c
}

func coord(idx int) (int, int) {
	return idx / (config.cols + 2), idx % (config.cols + 2)
}

func debug(s string, args ...interface{}) {
	if config.debug {
		fmt.Printf(s+"\n", args)
	}
}
