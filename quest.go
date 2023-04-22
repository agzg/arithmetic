// Port to Golang by Ali Azam <azam.vw@gmail.com>.

package main

import (
	"bufio"
	"fmt"
	"math"
	"math/rand"
	"os"
	"strconv"
	"strings"
	"time"
)

// question represents a problem posed by arithmetic.
type question struct {
	prompt string
	answer int
}

var rights, wrongs int

// penalty is the number of times an incorrect problem is repeated when penalised.
const penalty = 5

// penalised contains questions that have been answered incorrectly.
var penalised = make(map[string][]question)

// ask poses a question to the player using the given operands and range.
//
// In accordance to the design of the original versions, the player must figure out the right
// answer themselves to progress.
func ask(ops []string, maxRange int) error {
	op := ops[rand.Intn(len(ops))]
	q := quest(op, maxRange)
	start := time.Now()
	fmt.Print(q.prompt)

	// Offsets the penalty each time. This ensures that problems gotten incorrect in a row aren't
	// too severely penalised.
	p := penalty
	offset := 1

	for {
		var c int

		bio := bufio.NewReader(os.Stdin)
		d, err := bio.ReadString('\n')
		d = strings.TrimSuffix(d, "\n")

		// Could be a float.
		if strings.Count(d, ".") == 1 {
			d = strings.Split(d, ".")[0]
		}
		c, err = strconv.Atoi(d)
		if err != nil {
			fmt.Println("Please enter a number.")
			continue
		}

		if c == q.answer {
			fmt.Println("Right!")
			rights++
			break
		}
		fmt.Println("What?")
		wrongs++

		// Add question to penalised questions for this op.
		for i := 0; i < p; i++ {
			penalised[op] = append(penalised[op], q)
		}
		p -= offset
		offset += 1
	}
	dur += time.Since(start)
	return nil
}

// quest constructs a question using the given operand and range.
//
// If you get a question wrong, it is penalised. Penalised questions are repeated a set number of
// times (penalty) before being forgiven.
func quest(op string, maxRange int) question {
	// TODO make incrementally more likely
	// one-third chance if few penalised questions
	// 50-50 chance if there are too many penalised questions
	l := len(penalised[op])
	if l > 0 && rand.Intn(maxRange) == 0 {
		var q question
		penalised[op], q = pop(penalised[op], rand.Intn(l))
		return q
	}

	var a, b, answer int

	switch op {
	case "+":
		answer = rand.Intn(maxRange) + 1
		a = rand.Intn(answer)
		b = answer - a
	case "-":
		a = rand.Intn(maxRange) + 1
		b = rand.Intn(a)
		answer = a - b
	case "x", "×":
		mr := int(math.Sqrt(float64(maxRange))) // modified range
		a = rand.Intn(mr)
		b = rand.Intn(mr)
		answer = a * b
	case "/", "÷":
		// TODO possibly add a "strict" mode which requires entering answer correct to 1 dp
		// at least if it is enabled
		mr := int(math.Sqrt(float64(maxRange)))
		b = rand.Intn(mr)
		answer = rand.Intn(mr)
		a = b*answer + rand.Intn(b-1)
	}
	return question{fmt.Sprintf("%d %s %d =  ", a, op, b), answer}
}

// stats prints statistics on the performance so far.
func stats() {
	quests := wrongs + rights
	if quests < 1 {
		return
	}

	percent := 100 * rights / quests
	fmt.Printf("\n\nRights %d; Wrongs %d; Score %d%%\n", rights, wrongs, percent)

	if rights == 0 {
		return
	}

	per := dur.Seconds() / float64(rights)
	fmt.Printf("Total time %.1f seconds; %.1f seconds per problem\n\n", dur.Seconds(), per)
}

func pop(slice []question, i int) ([]question, question) {
	return append(slice[:i], slice[i+1:]...), slice[i]
}
