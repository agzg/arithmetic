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

// penalised contains questions that have been answered incorrectly.
var penalised = make(map[string][]question)

const penalty = 5

func ask(ops []string, maxRange int) error {
	op := ops[rand.Intn(len(ops))]
	q := quest(op, maxRange)
	start := time.Now()
	fmt.Print(q.prompt)

	// Add an offset to the penalty per question.
	//
	// This ensures that answering a question incorrectly multiple times in a row doesn't
	// lead to TOO many penalties.
	//
	// But...maybe it should?
	offset := 1
	p := penalty

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

		// Add question to slice of penalised questions for this op.
		for i := 0; i < p; i++ {
			penalised[op] = append(penalised[op], q)
		}
		p -= offset
		offset++
	}
	dur += time.Since(start)
	return nil
}

func quest(op string, maxRange int) question {
	// TODO make incrementally more likely
	// one-third chance if few penalised questions
	// 50-50 chance if there are too many penalised questions
	l := len(penalised[op])
	r := l / maxRange
	if l > 0 && ((r < 1.0 && rand.Intn(3) == 0) || rand.Intn(2) == 0) {
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
		mr := int(math.Sqrt(float64(maxRange))) // modified range
		b = rand.Intn(mr)
		answer = rand.Intn(mr)
		a = b*answer + rand.Intn(b-1)
	}
	return question{fmt.Sprintf("%d %s %d =  ", a, op, b), answer}
}

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
