package config

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v2"
)

var (
	timeRegex *regexp.Regexp = regexp.MustCompile(`(?m)\d\d:\d\d`)

	defaultConfiguration Configuration = Configuration{
		Job: Job{
			Start:        "08:00",
			End:          "17:59",
			SleepMinutes: 1,
		},
		Database: Database{
			Server: "localhost",
			Port:   1433,
			User:   "sa",
		},
	}
)

// Configuration is the representation of the config.yaml file
type Configuration struct {
	Job      Job
	Database Database
}

// Validate validates configuration values
func (c Configuration) Validate() error {
	validationMessage := c.Job.Validate() + c.Database.Validate()
	if validationMessage != "" {
		return fmt.Errorf("Error in the config file.\n" + validationMessage)
	}

	return nil
}

// Job holds the job's configuration
type Job struct {
	Start        string
	End          string
	SleepMinutes int `yaml:"sleepMinutes"`
}

// Database holds the database's configuration
type Database struct {
	Server       string
	Port         int
	User         string
	Password     string
	DatabaseName string `yaml:"databaseName"`
}

// Validate validates database values
func (j Job) Validate() string {
	validationMessage := ""

	if !IsValidTime(j.Start) {
		validationMessage += fmt.Sprintf("Invalid time given: %v. Should be in the form of hh:mm\n", j.Start)
	}

	if !IsValidTime(j.End) {
		validationMessage += fmt.Sprintf("Invalid time given: %v. Should be in the form of hh:mm\n", j.End)
	}

	if j.SleepMinutes == 0 {
		validationMessage += fmt.Sprintln("Config value job.sleepMinutes cannot be 0")
	}

	return validationMessage
}

// IsValidTime checks if the given time matches the expected "hh:mm" format
func IsValidTime(time string) bool {
	if len(timeRegex.FindStringIndex(time)) > 0 {
		return true
	}
	return false
}

// Validate validates database values
func (d Database) Validate() string {
	validationMessage := ""

	if d.Server == "" {
		validationMessage += fmt.Sprintln("database.server cannot be empty")
	}

	if d.User == "" {
		validationMessage += fmt.Sprintln("database.user cannot be empty")
	}

	if d.Password == "" {
		validationMessage += fmt.Sprintln("database.password cannot be empty")
	}

	if d.DatabaseName == "" {
		validationMessage += fmt.Sprintln("database.databaseName cannot be empty")
	}

	return validationMessage
}

// ReadConfiguration reads a YAML content and returns the equivalent Configuration struct
func ReadConfiguration(yamlConfiguration []byte) (Configuration, error) {
	configuration := Configuration{}
	err := yaml.UnmarshalStrict(yamlConfiguration, &configuration)
	if err != nil {
		return Configuration{}, err
	}

	err = configuration.Validate()
	if err != nil {
		return Configuration{}, err
	}

	return configuration, err
}
