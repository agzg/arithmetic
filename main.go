// https://github.com/jsm28/bsd-games/blob/master/arithmetic/arithmetic.c

package main

import (
	"fmt"
	"math"
	"math/rand"
	"strconv"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

const (
	maxN = 30
	minN = 15
)

var operands = []string{"+", "-", "x", "/"}

func main() {
	for {
		start := time.Now()
		n := rand.Intn(maxN-minN) + minN // Number of questions asked.
		var wrongs int

		for i := 0; i < n; i++ {
			a, b := rand.Intn(20), rand.Intn(20)

			if b > a {
				a, b = b, a
			}

			op := operands[rand.Intn(len(operands))]
			if b == 0 {
				op = operands[rand.Intn(len(operands)-1)]
			}
			var answer int

			switch op {
			case "+":
				answer = a + b
			case "-":
				answer = a - b
			case "×":
				answer = a * b
			case "/":
				answer = int(a / b)
			}

			var w bool
			fmt.Printf("%d %s %d =  ", a, op, b)
			for i := readInt(); i != answer; {
				fmt.Println("What?")
				i = readInt()
				w = true
			}
			fmt.Println("Correct!")

			if w {
				wrongs++
			}
		}
		rights := n - wrongs
		percent := int(math.Round(float64(rights)/float64(n)) * 100)
		fmt.Printf("\nRights %d; Wrongs %d; Score %d%%\n", rights, wrongs, percent)

		total := int(math.Round(time.Since(start).Seconds()))
		per := float64(total / n)
		fmt.Printf("Total time %d seconds; %.1f seconds per problem\n", total, per)

		fmt.Println("\nPress RETURN to continue")
		fmt.Scanln()
	}
}

func readInt() int {
	var a string

scan:
	_, err := fmt.Scan(&a)
	if err != nil {
		goto scan
	}

	i, err := strconv.Atoi(a)
	if err != nil {
		fmt.Println("Please enter a number.")
		goto scan
	}
	return i
}
