package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"strconv"
	"time"
)

var numCorrect int

var questionList = make(map[string]string)

//Iterate through map, asking each key as a question
//Compare the value to the answer read from the command line
func askQuestions(questionList map[string]string, quitChan chan bool) {
	for question, answer := range questionList {
		fmt.Print(question + ": ")
		var response int
		parsedAnswer, _ := strconv.Atoi(answer)
		if _, err := fmt.Scan(&response); err != nil {
			log.Fatal(err)
		}
		if response == parsedAnswer {
			numCorrect++
		}
	}
	quitChan <- true
}

func main() {
	//Arg parse
	problemFile := flag.String("csv", "problems.csv", "CSV file to load")
	timeLimit := flag.Int("limit", 30, "The time limit for the quiz in seconds")
	flag.Parse()

	//Open file
	questions, err := os.Open(*problemFile)
	if err != nil {
		log.Fatal(err)
	}
	defer questions.Close()

	//Could use bufio if file is large, but it shouldn't be.
	//reader := csv.NewReader(bufio.NewReader(questions))
	reader := csv.NewReader(questions)
	for {
		entry, err := reader.Read()
		if err == io.EOF {
			break
		}
		if err != nil {
			log.Fatal(err)
		}
		questionList[string(entry[0])] = string(entry[1])
	}

	timeout := time.NewTimer(time.Second * time.Duration(*timeLimit))
	defer timeout.Stop()
	quitChan := make(chan bool, 1)

	// Ask questions, but use go to allow for select to happen
	go askQuestions(questionList, quitChan)

	select {
	case <-timeout.C:
		fmt.Println("Timeout reached!")
		break
	case <-quitChan:
		break
	}
	fmt.Printf("You scored %d out of %d.", numCorrect, len(questionList))
}
