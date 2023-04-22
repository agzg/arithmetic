// Port to Golang by Ali Azam <azam.vw@gmail.com>.

package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"time"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// Questions asked in one go will be between minQuests and maxQuests inclusive.
const (
	minQuests = 15
	maxQuests = 30
)

var (
	ops        = []string{"+", "-", "x", "/"}
	opsList    = strings.Join(ops, "")
	opsDefault = []string{"+", "-"}
	opsUni     = []string{"×", "÷"}
	convertUni = map[string]string{"x": "×", "/": "÷"}

	flagRange = flag.Int("r", 10, "range")
	flagOps   = flag.String("o", strings.Join(opsDefault, ""), opsList)
	flagUni   = flag.Bool("u", false, "unicode")
)

func usage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [-o %s] [-r range] [-u unicode]\n", os.Args[0], opsList)
	os.Exit(2)
}

// Duration of time spent answering questions.
var dur time.Duration

func main() {
	flag.Usage = usage
	flag.Parse()

	var operands []string

	for _, o := range *flagOps {
		op := string(o)
		if strings.Count(strings.Join(opsUni, ""), op) > 0 {
			operands = append(operands, op)
			*flagUni = true // Use unicode ops if these are provided in -o.
			continue
		}
		if strings.Count(opsList, op) < 1 {
			fmt.Fprintf(os.Stderr, "arithmetic: unknown key.\n")
			usage()
		}
		operands = append(operands, op)
	}

	for _, op := range operands {
		penalised[op] = []question{}
	}

	// Convert operands to Unicode if required.
	if *flagUni {
		for i, op := range operands {
			if u, ok := convertUni[op]; ok {
				operands[i] = u
			}
		}
	}

	lr := int(^uint(0) >> 1) // Largest range we can accept.
	if *flagRange < 0 || *flagRange > lr {
		fmt.Fprintf(os.Stderr, "arithmetic: invalid range (0 < r < %d).\n", lr)
		usage()
	}
	maxRange := *flagRange
	if maxRange == 0 { // Someone's feeling adventurous.
		maxRange = lr
		fmt.Printf("Warning: numbers can get as high as %d!\n", lr)
		fmt.Println("Press ENTER to continue...")
		fmt.Scanln()
		fmt.Println()
	}

	// Handle keyboard interrupts in a separate goroutine.
	//
	// Much like the previous versions, once interrupted, arithmetic prints out performance
	// statistics before exiting.
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			stats()
			fmt.Println()
			os.Exit(0)
		}
	}()

	for {
		quests := rand.Intn(maxQuests-minQuests) + minQuests
		for i := 0; i < quests; i++ {
			if err := ask(operands, maxRange); err != nil { // Could be anything.
				fmt.Fprintf(os.Stderr, "arithmetic: error: %s\n", err)
				os.Exit(1)
			}
		}
		stats()
		fmt.Print("Press ENTER to continue...\n")
		fmt.Scanln()
		fmt.Println()
	}
}
