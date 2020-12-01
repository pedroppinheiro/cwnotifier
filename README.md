# cwnotifier

![CircleCI](https://img.shields.io/circleci/build/github/pedroppinheiro/cwnotifier)
[![Go Report Card](https://goreportcard.com/badge/github.com/pedroppinheiro/cwnotifier)](https://goreportcard.com/report/github.com/pedroppinheiro/cwnotifier)
[![GoDoc](https://godoc.org/github.com/pedroppinheiro/cwnotifier?status.svg)](https://godoc.org/github.com/pedroppinheiro/cwnotifier)

Notifier for cherwell's tasks. The purpose of this project is to never miss an SLA again.

## Build

To avoid opening a console at application startup, use these compile flags ([source](https://stackoverflow.com/a/36728885/1252947)):

```sh
go build -ldflags="-H=windowsgui -X main.version=$(git describe --tags --always)"
```

## Notes

- The .exe file must be placed alongside the "assets" folder, otherwise some images will not be displayed.

- The log file will be created in the same folder as the .exe file

- In order for the program to connect to the database, there should be a "config.yaml" file in the same folder as the .exe file. Here is a basic template of config.yaml:

```yaml
job:
    start: "08:00" # A partir de qual horário o programa irá checar o cherwell
    end: "17:59" # Até qual horário o programa irá checar o cherwell
    sleepMinutes: 5 # De quanto em quanto tempo em minutos o programa deve checar o cherwell
database: # Dados para a conexão com o banco de dados
    server: ""
    port: 1433
    user: ""
    password: ""
    databaseName: "" # Nome do banco de dados do cherwell
```
