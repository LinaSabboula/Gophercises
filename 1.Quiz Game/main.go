package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
	"time"
)

const (
	HARD   = "hard"
	NORMAL = "normal"
)

func getInput(prompt string, reader *bufio.Reader, answerChannel chan string) {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	answerChannel <- strings.TrimSpace(input)
}

// difficulty to switch between ticker and ticker
func main() {
	fileName := flag.String("file", "problems.csv", "csv file in format of 'question,answer'")
	difficulty := flag.String("difficulty", NORMAL, "Difficulty setting for the quiz")
	var defaultTimer int
	if *difficulty == HARD {
		defaultTimer = 5
	} else {
		defaultTimer = 90
	}
	timerValue := flag.Int("limit", defaultTimer, "Time limit for each question")

	timerDuration := time.Duration(*timerValue) * time.Second
	flag.Parse()
	fmt.Println(*fileName, *difficulty, *timerValue)

	fmt.Println("Welcome to the super smart quiz game!")

	file, err := os.Open(*fileName)
	if err != nil {
		log.Fatal("Unable to read input file "+*fileName, err)
	}
	records, _ := csv.NewReader(file).ReadAll()
	maxScore := len(records)
	score := 0
	reader := bufio.NewReader(os.Stdin)

	answerChannel := make(chan string)

	timer := time.NewTimer(timerDuration)

	for _, record := range records {
		go getInput(record[0], reader, answerChannel)
		select {
		case answer := <-answerChannel:
			if strings.EqualFold(answer, record[1]) {
				score++
			}
		case <-timer.C:
			fmt.Printf("\nTimer's up! Next Question\n")
			if *difficulty == HARD {
				timer.Reset(timerDuration)
			}
			break
		}

	}
	fmt.Printf("\nCongratulations your score is: %d/%d\n", score, maxScore)
}
