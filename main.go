package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/pedroppinheiro/cwnotifier/database"
	"github.com/pedroppinheiro/cwnotifier/notifier"

	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	validTimeStart string = "08:00"
	validTimeEnd   string = "18:00"
	jobSleepPeriod        = 5 * time.Minute
)

func init() {
	// configuring log to file and console
	logFile, err := os.OpenFile("log.txt", os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	//multiWriter := io.MultiWriter(os.Stdout, logFile)
	log.SetOutput(logFile)
}

func main() {
	log.Printf("CWNotifier program has started")
	notifier.NotifyProgramStart()

	defer func() {
		if r := recover(); r != nil {
			notifier.NotifyError()
			database.CloseConnection()
			log.Fatal("CWNotifier is closing")
		}
	}()

	database.Connect()
	defer database.CloseConnection()

	for true {

		shouldRun, err := shouldCheckDatabase(time.Now())

		if shouldRun {
			if database.GetNumberOfPriorityTasks() >= 1 {
				notifier.Notify()
			}
		} else {
			log.Println("Skipped checking cherwell. ", err)
		}

		time.Sleep(jobSleepPeriod)
	}
}

func shouldCheckDatabase(givenTime time.Time) (bool, error) {
	isSaturday := givenTime.Weekday() == time.Saturday
	isSunday := givenTime.Weekday() == time.Sunday

	if isSaturday || isSunday {
		return false, errors.New("Current date is a weekend")
	}

	newLayout := "15:04"
	start, _ := time.Parse(newLayout, validTimeStart)
	end, _ := time.Parse(newLayout, validTimeEnd)

	checkTimeString := fmt.Sprintf("%02d:%02d", givenTime.Hour(), givenTime.Minute()) // https://stackoverflow.com/a/51546906/1252947
	check, _ := time.Parse(newLayout, checkTimeString)

	isBetweenValidTime := inTimeSpan(start, end, check)

	if isBetweenValidTime {
		return true, nil
	}

	return false, errors.New("Current time is not between valid work time range")
}

// https://stackoverflow.com/a/55093788/1252947
func inTimeSpan(start, end, check time.Time) bool {
	if start.Before(end) {
		return !check.Before(start) && !check.After(end)
	}
	if start.Equal(end) {
		return check.Equal(start)
	}
	return !start.After(check) || !end.Before(check)
}
