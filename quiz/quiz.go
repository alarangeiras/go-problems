package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"log"
	"os"
	"time"
)

const (
	DEFAULT_FILE_NAME = "problems.csv"
	DEFAULT_TIME = 30
)

var totalQuestions, correctQuestions int

func main() {
	var csvFileame string
	var quizTime int
	flag.StringVar(&csvFileame, "filename", DEFAULT_FILE_NAME, "full path of CSV filename")
	flag.IntVar(&quizTime, "time", DEFAULT_TIME, "time for answers in seconds")
	flag.Parse()
	finished := make(chan bool, 1)

	lines := readLines(csvFileame)
	totalQuestions = len(lines)

	fmt.Printf("\nWould you like to start, you'll have %d seconds?", quizTime)
	fmt.Scanln()
	go evaluateQuizTime(quizTime, finished)
	go presentQuiz(lines, &correctQuestions, finished)

	<-finished
	fmt.Println("total questions: ", totalQuestions)
	fmt.Println("correct questions: ", correctQuestions)

}

func presentQuiz(lines [][]string, totalCorrectQuestions *int, finished chan<- bool) {
	for pos, line := range lines {
		question := line[0]
		correctAnswer := line[1]

		var userAnswer string
		fmt.Printf("\nQuestion %d: %s >> ", pos + 1, question)
		fmt.Scan(&userAnswer)

		if userAnswer == correctAnswer {
			*totalCorrectQuestions++
		}
	}

	finished <- true
}

func evaluateQuizTime(quizTime int, finished chan<- bool) {
	timer := time.NewTimer(time.Duration(quizTime) * time.Second)
	<-timer.C
	fmt.Printf("\nThe time has been expired.\n")
	finished <- true
}

func extractQuestionsAndAnswers(lines [][]string) (questions []string, answers []string) {
	for _, line := range lines {
		questions = append(questions, line[0])
		answers = append(answers, line[1])
	}
	return
}

func evaluateResult(userAnswers, correctAnswers []string) (totalQuestions int, correctQuestions int) {
	for pos, userAnswer := range userAnswers {
		correctAnswer := correctAnswers[pos]

		totalQuestions++
		if correctAnswer == userAnswer {
			correctQuestions++
		}
	}	
	return
}

func readLines(csvFilename string) [][]string {
	file, err := os.Open(csvFilename)
	if err != nil {
		log.Fatalf("\ncould not open filename: %s", csvFilename)
	}
	defer file.Close()
	csvReader := csv.NewReader(file)
	records, err := csvReader.ReadAll()
	if err != nil {
		log.Fatalf("\ncould not read csv file")
	}
	return records
} 