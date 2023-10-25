package main

import (
	"encoding/csv"
	"flag"
	"fmt"
	"os"
	"time"
)

func exit(msg string) {
	fmt.Println(msg)
	os.Exit(1)
}

func parse_csv(fileName string) [][]string {
	f, err := os.Open(fileName)
	if err != nil {
		exit(fmt.Sprintf("Failed to open the CSV file: %s\n", fileName))
	}

	r := csv.NewReader(f)

	records, err := r.ReadAll()
	if err != nil {
		exit("Failed to parse the provided CSV file.")
	}

	f.Close()
	return records
}

var fileName = flag.String("csv", "problems.csv", "a csv file in the format of 'question,answer'")
var limit = flag.Int("limit", 30, "the time limit for the quiz in seconds")

func quiz(records [][]string, right *int, stop chan bool) {
	for number, question := range records {
		var awnser string

		fmt.Printf("Problem #%d: %s = ", number+1, question[0])

		fmt.Scanf("%s\n", &awnser)

		if awnser == question[1] {
			*right++
		}
	}

	stop <- true
}

func exitAndShowScore(right int, tam int) {
	fmt.Printf("You Scored %d out of %d", right, tam)
	os.Exit(0)
}

func timerRun(c chan bool) {
	timer := time.NewTimer(time.Duration(*limit) * time.Second)
	<-timer.C
	c <- true
}

func main() {
	flag.Parse()
	records := parse_csv(*fileName)
	right := 0

	done := make(chan bool, 1)

	go quiz(records, &right, done)
	go timerRun(done)

	<-done
	exitAndShowScore(right, len(records))
}
