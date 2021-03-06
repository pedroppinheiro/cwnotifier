package main

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/getlantern/systray"
	"github.com/pedroppinheiro/cwnotifier/config"
	"github.com/pedroppinheiro/cwnotifier/database"
	"github.com/pedroppinheiro/cwnotifier/notifier"

	"log"

	_ "github.com/denisenkom/go-mssqldb"
)

const (
	defaultYAMLName string = "config.yaml"
	defaultLogName  string = "log.txt"
)

// Version will be defined in compile time.
var version = "undefined"

func init() {
	// configuring log to file
	logFile, err := os.OpenFile(defaultLogName, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0666)
	if err != nil {
		panic(err)
	}
	log.SetOutput(logFile)
}

func main() {
	log.Printf("CWNotifier is starting. Program version: %v", version)
	systray.Run(onReady, nil)
	log.Printf("CWNotifier has finished")
}

func onReady() {
	defer recoverFromError()

	configureSystemtray()

	configuration, err := readConfiguration(defaultYAMLName)
	if err != nil {
		log.Panic(err)
	}

	// used to maintain compatibility with previous versions in which the default was "SUSIS - GERIN"
	if configuration.User.Team == "" {
		configuration.User.Team = "SUSIS - GERIN"
		log.Println("user.team is empty, using \"SUSIS - GERIN\" as fallback.")
	}

	if !configuration.Notification.IsNotificationsEnabled() {
		notifier.NotifyNoNotificationsEnabled()
		systray.Quit()
		return
	}

	database.Connect(configuration.Database)
	defer database.CloseConnection()

	notifier.NotifyProgramStart()
	for ; true; time.Sleep(time.Duration(configuration.Job.SleepMinutes) * time.Minute) {
		shouldNotify, err := shouldCheckDatabase(time.Now(), configuration.Job)
		if !shouldNotify || err != nil {
			log.Println("Skipped checking cherwell. ", err)
			continue
		}

		if configuration.Notification.EnableIncidentsWithoutOwnerNotification {
			notifyIncidentsWithoutOwnerNotification(configuration)
		}

		if configuration.Notification.EnableTasksWithoutOwnerNotification {
			notifyTasksWithoutOwnerNotification(configuration)
		}

		if configuration.Notification.EnableIncidentsWithClosedTasksNotification {
			notifyIncidentsWithClosedTasksNotification(configuration)
		}

		if configuration.Notification.EnableChangesThatNeedToBeValidatedNotification {
			notifyChangesThatNeedToBeValidated(configuration)
		}

		if configuration.Notification.EnableChangesThatRequireUpdateNotification {
			notifyChangesThatRequireUpdate(configuration)
		}
	}
}

func notifyIncidentsWithoutOwnerNotification(configuration config.Configuration) {
	incidents := database.GetIncidentsWithoutOwner(configuration.User.Team)
	if len(incidents) >= 1 {
		notifier.NotifyIncidentsWithoutOwner(strings.Join(incidents[:], ","))
	}
}

func notifyTasksWithoutOwnerNotification(configuration config.Configuration) {
	tasks := database.GetTasksWithoutOwner(configuration.User.Team, configuration.User.Email)
	if len(tasks) >= 1 {
		notifier.NotifyTasksWithoutOwner(strings.Join(tasks[:], ","))
	}
}

func notifyIncidentsWithClosedTasksNotification(configuration config.Configuration) {
	incidents := database.GetIncidentsWithClosedTasks(configuration.User.Team, configuration.User.Name)
	if len(incidents) >= 1 {
		notifier.NotifyIncidentsWithClosedTasks(strings.Join(incidents[:], ","))
	}
}

func notifyChangesThatNeedToBeValidated(configuration config.Configuration) {
	changes := database.GetChangesThatNeedToBeValidated(configuration.User.Name)
	if len(changes) >= 1 {
		notifier.NotifyChangesThatNeedToBeValidated(strings.Join(changes[:], ","))
	}
}

func notifyChangesThatRequireUpdate(configuration config.Configuration) {
	changes := database.GetChangesThatRequireUpdate(configuration.User.Name)
	if len(changes) >= 1 {
		notifier.NotifyChangesThatRequireUpdate(strings.Join(changes[:], ","))
	}
}

func recoverFromError() {
	if r := recover(); r != nil {
		notifier.NotifyError()
		log.Fatal("CWNotifier is closing due to errors")
	}
}

// https://dev.to/osuka42/building-a-simple-system-tray-app-with-go-899
func configureSystemtray() {
	systray.SetIcon(readFileContent("assets\\app.ico"))
	systray.SetTitle("CWNotifier")
	systray.SetTooltip("CWNotifier")

	showLogMenuItem := systray.AddMenuItem("Show log", "Show the app's log")
	showLogMenuItem.SetIcon(readFileContent("assets\\log.ico"))
	go func() {
		for {
			<-showLogMenuItem.ClickedCh
			cmd := exec.Command("notepad", defaultLogName)
			if err := cmd.Run(); err != nil {
				log.Println("An error occurred during show log menu action. ", err)
			}
		}
	}()

	quitMenuItem := systray.AddMenuItem("Quit", "Quit the app")
	quitMenuItem.SetIcon(readFileContent("assets\\quit.ico"))
	go func() {
		<-quitMenuItem.ClickedCh
		log.Println("User requested to quit")
		systray.Quit()
	}()
}

func readConfiguration(yamlLocation string) (config.Configuration, error) {
	yamlContent := readFileContent(yamlLocation)
	return config.ReadConfiguration(yamlContent)
}

func readFileContent(filePath string) []byte {
	content, err := ioutil.ReadFile(filePath)
	if err != nil {
		log.Panic(err)
	}
	return content
}

func shouldCheckDatabase(givenTime time.Time, jobConfig config.Job) (bool, error) {

	if isWeekend(givenTime) {
		return false, errors.New("Current date is a weekend")
	}

	checkTimeString := fmt.Sprintf("%02d:%02d", givenTime.Hour(), givenTime.Minute()) // https://stackoverflow.com/a/51546906/1252947
	isBetweenValidTime, err := inTimeSpan(jobConfig.Start, jobConfig.End, checkTimeString)

	if err != nil {
		log.Panic("An unexpected error occurred during verification of the time to run the job. ", err)
	}

	if isBetweenValidTime {
		return true, nil
	}

	return false, errors.New("Current time is not between valid work time range")
}

// isWeekend returns true if the given time is a weekend, otherwise returns false
func isWeekend(givenTime time.Time) bool {
	isSaturday := givenTime.Weekday() == time.Saturday
	isSunday := givenTime.Weekday() == time.Sunday

	return isSaturday || isSunday
}

// inTimeSpan returns true if a "check" time is between a "start" and an "end" range.
// Parameters must be given in the form of "hh:mm"
// https://stackoverflow.com/a/55093788/1252947
func inTimeSpan(start, end, check string) (bool, error) {
	if !config.IsValidTime(start) {
		return false, fmt.Errorf("Invalid time given: \"%v\". Should be in the form of hh:mm", start)
	}

	if !config.IsValidTime(end) {
		return false, fmt.Errorf("Invalid time given: \"%v\". Should be in the form of hh:mm", end)
	}

	if !config.IsValidTime(check) {
		return false, fmt.Errorf("Invalid time given: \"%v\". Should be in the form of hh:mm", check)
	}

	newLayout := "15:04"
	startTime, err := time.Parse(newLayout, start)
	if err != nil {
		return false, err
	}
	endTime, err := time.Parse(newLayout, end)
	if err != nil {
		return false, err
	}
	checkTime, err := time.Parse(newLayout, check)
	if err != nil {
		return false, err
	}

	if startTime.Before(endTime) {
		return !checkTime.Before(startTime) && !checkTime.After(endTime), nil
	}
	if startTime.Equal(endTime) {
		return checkTime.Equal(startTime), nil
	}
	return !startTime.After(checkTime) || !endTime.Before(checkTime), nil
}
