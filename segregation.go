package main

import (
	"flag"
	"fmt"
	"math/rand"
	"time"
)

var matrix [][]int
var vlist []int
var round, satisfied int

var config struct {
	rows, cols int
	similar    int
	red        int
	empty      int
	delay      int
	alg        int

	debug bool
}

const (
	empty = iota
	red
	blue
)

func init() {
	flag.IntVar(&config.rows, "rows", 24, "number of rows")
	flag.IntVar(&config.cols, "cols", 80, "number of columns")

	flag.IntVar(&config.similar, "similar", 30, "percent of similar neighbours to be satisfied")
	flag.IntVar(&config.red, "red", 50, "percent of red agents (blue = 100 - red)")
	flag.IntVar(&config.empty, "empty", 10, "percent of empty sites")
	flag.IntVar(&config.alg, "alg", 0, "relocation algorithm (0-4)")

	flag.IntVar(&config.delay, "delay", 100, "simulation step delay")
	flag.BoolVar(&config.debug, "debug", false, "debug output")

	flag.Parse()

	rand.Seed(time.Now().UnixNano())

	matrix = make([][]int, config.rows+2)
	for i := 0; i < len(matrix); i++ {
		matrix[i] = make([]int, config.cols+2)
	}
}

func initmatrix() {

	sites := config.rows * config.cols
	actorsnum := sites * (100 - config.empty) / 100
	rednum := actorsnum * config.red / 100

	// fmt.Printf("\n --- matrix configuration: \n")
	// fmt.Printf(" Total sites: %d\n", sites)
	// fmt.Printf(" Empty      : %d\n", sites*config.empty/100)
	// fmt.Printf(" Red        : %d\n", rednum)
	// fmt.Printf(" Blue       : %d\n", actorsnum-rednum)
	// fmt.Printf(" Checksum   : %d\n", sites-rednum-(actorsnum-rednum)-(sites*config.empty/100))
	// fmt.Printf(" --- \n")

	// set red sites
	for i := 0; i < rednum; i++ {
		for {
			rnd := rand.Intn(sites)
			r := rnd/config.cols + 1
			c := rnd%config.cols + 1
			if matrix[r][c] == empty {
				matrix[r][c] = red
				break
			}
		}
	}

	// set blue sites
	for i := 0; i < actorsnum-rednum; i++ {
		for {
			rnd := rand.Intn(sites)
			r := rnd/config.cols + 1
			c := rnd%config.cols + 1
			if matrix[r][c] == empty {
				matrix[r][c] = blue
				break
			}
		}
	}

	// create vacant sites list
	vlist = make([]int, 0)
	for r := 1; r < config.rows+1; r++ {
		for c := 1; c < config.cols+1; c++ {
			if matrix[r][c] == empty {
				vlist = append(vlist, r*(config.cols+2)+c)
			}
		}
	}
	//fmt.Printf("%v\n", vlist)
	draw()
	time.Sleep(time.Duration(config.delay) * time.Millisecond)
}

func draw() {
	pos(0, 0)
	for r := 0; r < config.rows+2; r++ {
		for c := 0; c < config.cols+2; c++ {
			switch matrix[r][c] {
			case empty:
				fmt.Printf("  ")
			case red:
				fmt.Printf("\033[31mX ")
			case blue:
				fmt.Printf("\033[36mX ")
			}
		}
		fmt.Printf("\033[0m\n")
	}
	fmt.Printf("\nRound %d | Satisfied %3d%% | Alg %d: %s",
		round,
		satisfied*100/(config.rows*config.cols),
		config.alg,
		algdesc[config.alg])
}

func findsite(kind, r, c int) (int, int) {
	switch config.alg {
	case 0:
		// find a random available sitew here agent is satisfied
		for j := 0; j < len(vlist); j++ { // limit number of attempts
			i := rand.Intn(len(vlist))
			if utility(kind, r, c) < utility(kind, vlist[i]/(config.cols+2), vlist[i]%(config.cols+2)) {
				v := vlist[i]
				vlist = append(vlist[:i], vlist[i+1:]...)
				return coord(v)
			}
		}
		fallthrough
	case 1:
		// find a random available site
		i := rand.Intn(len(vlist))
		v := vlist[i]
		vlist = append(vlist[:i], vlist[i+1:]...)
		return coord(v)
	case 2:
		// find first available vacant site ordered by vacant time where agent is satisfied
		for i := 0; i < len(vlist); i++ {
			if utility(kind, r, c) < utility(kind, vlist[i]/(config.cols+2), vlist[i]%(config.cols+2)) {
				v := vlist[i]
				vlist = append(vlist[:i], vlist[i+1:]...)
				return coord(v)
			}
		}
		return -1, -1
	case 3:
		// find first available vacant site ordered by vacant time
		v := vlist[0]
		vlist = vlist[1:]
		return coord(v)
	default:
		// find first available vacant site where agent is satisfied
		for row := 1; row < config.rows+1; row++ {
			for col := 1; col < config.cols+1; col++ {
				if matrix[row][col] == empty &&
					utility(kind, r, c) < utility(kind, row, col) {
					return row, col
				}
			}
		}
		return -1, -1
	}
}

func isSatisfied(kind, r, c int) bool {
	if kind == empty {
		return true
	}
	return utility(kind, r, c) >= config.similar
}

func utility(kind, r, c int) int {
	if kind == empty {
		return 0
	}

	redcount := 0
	bluecount := 0

	// count neighbours of each type
	for i := c - 1; i <= c+1; i++ {
		for j := r - 1; j <= r+1; j++ {
			if j != r || i != c {
				switch matrix[j][i] {
				case 1:
					redcount++
				case 2:
					bluecount++
				}
			}
		}
	}

	total := bluecount + redcount
	if total == 0 {
		return 0
	}

	if kind == red {
		return redcount * 100 / total
	}

	return bluecount * 100 / total
}

func move() {
	round++                 // increment round count
	satisfied = 0           // satisfied counter
	ulist := make([]int, 0) // list of unsatisfied agents

	// find all unsatisfied agents
	for r := 1; r < config.rows+1; r++ {
		for c := 1; c < config.cols+1; c++ {
			if isSatisfied(matrix[r][c], r, c) {
				satisfied++
			} else {
				ulist = append(ulist, ind(r, c))
			}
		}
	}

	// move unsatisfied agents if possible
	rand.Shuffle(len(ulist), func(i, j int) { ulist[i], ulist[j] = ulist[j], ulist[i] })
	for _, a := range ulist {
		r, c := coord(a)

		// relocate if possible
		nr, nc := findsite(matrix[r][c], r, c)
		if nr != -1 {
			matrix[r][c], matrix[nr][nc] = matrix[nr][nc], matrix[r][c]
			vlist = append(vlist, a)
			if config.alg == 2 || config.alg == 3 {
				vlist[0], vlist[len(vlist)-1] = vlist[len(vlist)-1], vlist[0]
			}
		}
	}

}

func main() {
	cls()
	initmatrix()
	for {
		move()
		draw()
		if satisfied == config.rows*config.cols {
			break
		}
		time.Sleep(time.Duration(config.delay) * time.Millisecond)
	}
	fmt.Println()
}
