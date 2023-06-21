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

func getInput(prompt string, reader *bufio.Reader, answerChannel chan string) {
	fmt.Print(prompt)
	input, _ := reader.ReadString('\n')
	answerChannel <- strings.TrimSpace(input)
}
func main() {
	fileName := flag.String("file", "problems.csv", "csv file in format of 'question,answer'")
	timerValue := flag.Int("limit", 30, "Time limit for each question")
	timerDuration := time.Duration(*timerValue) * time.Second
	timer := time.NewTimer(timerDuration)
	flag.Parse()

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
	for _, record := range records {
		go getInput(record[0], reader, answerChannel)
		select {
		case answer := <-answerChannel:
			if strings.EqualFold(answer, record[1]) {
				score++
			}
		case <-timer.C:
			fmt.Printf("\nCongratulations your score is: %d/%d\n", score, maxScore)
			return
		}

	}
	fmt.Printf("\nCongratulations your score is: %d/%d\n", score, maxScore)
}
