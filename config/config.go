package config

import (
	"fmt"

	"gopkg.in/yaml.v2"
)

var defaultConfiguration Configuration = Configuration{
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

// Configuration is the representation of the config.yaml file
type Configuration struct {
	Job      Job
	Database Database
}

// Validate validates configuration values
func (c Configuration) Validate() error {
	return c.Database.Validate()
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
func (d Database) Validate() error {
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

	if validationMessage != "" {
		return fmt.Errorf("Error in the config file.\n" + validationMessage)
	}

	return nil
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
