package main

import (
	"bufio"
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"strings"
)

func getInput(prompt string, reader *bufio.Reader) (string, error) {
	fmt.Print(prompt)
	input, err := reader.ReadString('\n')
	return strings.TrimSpace(input), err
}
func main() {
	fileName := flag.String("file", "problems.csv", "csv file in format of 'question,answer'")
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
	for _, record := range records {
		userAnswer, _ := getInput(record[0], reader)
		if strings.EqualFold(userAnswer, record[1]) {
			score++
		}
	}
	fmt.Printf("Congratulations your score is: %d/%d\n", score, maxScore)
}
