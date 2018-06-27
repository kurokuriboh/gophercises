package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

func main() {
	ind := 0
	score := 0
	var q [][]string

	// Setting up flags
	var t = flag.Int("limit", 30, "the time limit for the quiz in seconds (default 30)")
	var c = flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
	flag.Parse()

	// Open CSV file
	f, err := os.Open(*c)
	if err != nil {
		panic(err)
	}

	// Read all lines from CSV file to array
	r := csv.NewReader(bufio.NewReader(f))
	for {
		line, err := r.Read()
		if err == io.EOF {
			break
		} else if err != nil {
			panic(err)
		}
		q = append(q, line)
	}

	// Shuffle the questions
	for i := range q {
		j := rand.Intn(i + 1)
		q[i], q[j] = q[j], q[i]
	}

	// Create a channel for receiving answers from Stdin
	ch := make(chan string)
	go func() {
		var answer string
		for {
			fmt.Scanf("%s\n", &answer)
			ch <- answer
		}
	}()

	// Start timer for the quiz
	timer := time.NewTimer(time.Duration(*t) * time.Second)

rangeLoop:
	for _, l := range q {
		ind++
		fmt.Printf("Problem #%d: %s = ", ind, l[0])

		select {
		case <-timer.C:
			fmt.Println()
			break rangeLoop
		case answer := <-ch:
			answer = strings.TrimRight(answer, "\n")
			answer = strings.TrimSpace(answer)
			if strings.ToLower(answer) == strings.ToLower(l[1]) {
				score++
			}
		}
	}

	fmt.Printf("You scored %d out of %d.\n", score, len(q))
}
