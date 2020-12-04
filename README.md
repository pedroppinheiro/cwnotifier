# cwnotifier

![CircleCI](https://img.shields.io/circleci/build/github/pedroppinheiro/cwnotifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/pedroppinheiro/cwnotifier)](https://goreportcard.com/report/github.com/pedroppinheiro/cwnotifier)
[![GoDoc](https://godoc.org/github.com/pedroppinheiro/cwnotifier?status.svg)](https://godoc.org/github.com/pedroppinheiro/cwnotifier)

Notifier for cherwell's priority incidents. The purpose of this project is to never miss an SLA again.

## Build

To avoid opening a console at application startup, use these compile flags ([source](https://stackoverflow.com/a/36728885/1252947)):

```sh
go build -ldflags="-H=windowsgui -X main.version=$(git describe --tags --always)"
```

## Notes

- The .exe file must be placed alongside the "assets" folder, otherwise some images will not be displayed.

- The log file will be created in the same folder as the .exe file

- In order for the program to connect to the database and to perform other operations, there should be a "config.yaml" file in the same folder as the .exe file. Here is a basic template of the config.yaml:

```yaml
user:
  name: ""
  email: ""
  team: ""

notification:
   enableIncidentsWithoutOwnerNotification: true
   enableTasksWithoutOwnerNotification: true
   enableIncidentsWithClosedTasksNotification: true
   enableChangesThatNeedToBeValidatedNotification: true

job:
  start: "08:00"
  end: "17:59"
  sleepMinutes: 1

database:
  server: ""
  port: 1433
  user: ""
  password: ""
  databaseName: ""
```
