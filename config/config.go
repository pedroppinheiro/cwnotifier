package config

import (
	"fmt"
	"regexp"

	"gopkg.in/yaml.v2"
)

var timeRegex *regexp.Regexp = regexp.MustCompile(`(?m)\d\d:\d\d`)

// Configuration is the representation of the config.yaml file
type Configuration struct {
	User         User
	Notification Notification
	Job          Job
	Database     Database
}

// Validate validates configuration values
func (c Configuration) Validate() error {
	validationMessage := c.Job.Validate() + c.Database.Validate()
	if validationMessage != "" {
		return fmt.Errorf("Error in the config file.\n" + validationMessage)
	}

	return nil
}

// User holds the user's configuration
type User struct {
	Name  string
	Email string
	Team  string
}

// Notification holds the notification's configuration
type Notification struct {
	EnableIncidentsWithoutOwnerNotification        bool
	EnableTasksWithoutOwnerNotification            bool
	EnableIncidentsWithClosedTasksNotification     bool
	EnableChangesThatNeedToBeValidatedNotification bool
	gotMarshalled                                  bool
}

// IsNotificationsEnabled returns true if there is at least one notification enabled, otherwise returns false
func (notification Notification) IsNotificationsEnabled() bool {
	return notification.EnableIncidentsWithClosedTasksNotification ||
		notification.EnableIncidentsWithoutOwnerNotification ||
		notification.EnableTasksWithoutOwnerNotification ||
		notification.EnableChangesThatNeedToBeValidatedNotification
}

// UnmarshalYAML interface is implemented to give a custom behaviour when marshalling the yaml to the "Regex" field.
// See https://godoc.org/gopkg.in/yaml.v2#Unmarshaler for more details
func (notification *Notification) UnmarshalYAML(unmarshal func(interface{}) error) error {
	m := make(map[string]bool)
	var err error
	if err = unmarshal(&m); err != nil {
		return err
	}

	value, isPresent := m["enableIncidentsWithoutOwnerNotification"]
	if isPresent {
		notification.EnableIncidentsWithoutOwnerNotification = value
	} else {
		notification.EnableIncidentsWithoutOwnerNotification = true
	}

	value, isPresent = m["enableTasksWithoutOwnerNotification"]
	if isPresent {
		notification.EnableTasksWithoutOwnerNotification = value
	} else {
		notification.EnableTasksWithoutOwnerNotification = true
	}

	value, isPresent = m["enableIncidentsWithClosedTasksNotification"]
	if isPresent {
		notification.EnableIncidentsWithClosedTasksNotification = value
	} else {
		notification.EnableIncidentsWithClosedTasksNotification = true
	}

	value, isPresent = m["enableChangesThatNeedToBeValidatedNotification"]
	if isPresent {
		notification.EnableChangesThatNeedToBeValidatedNotification = value
	} else {
		notification.EnableChangesThatNeedToBeValidatedNotification = true
	}

	notification.gotMarshalled = true

	return nil
}

// Job holds the job's configuration
type Job struct {
	Start        string
	End          string
	SleepMinutes int `yaml:"sleepMinutes"`
}

// Validate validates database values
func (j Job) Validate() string {
	validationMessage := ""

	if !IsValidTime(j.Start) {
		validationMessage += fmt.Sprintf("job.start is invalid. Should be in the form of hh:mm, but got \"%v\"\n", j.Start)
	}

	if !IsValidTime(j.End) {
		validationMessage += fmt.Sprintf("job.end is invalid. Should be in the form of hh:mm, but got \"%v\"\n", j.End)
	}

	if j.SleepMinutes == 0 {
		validationMessage += fmt.Sprintln("job.sleepMinutes cannot be 0")
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

// Database holds the database's configuration
type Database struct {
	Server       string
	Port         int
	User         string
	Password     string
	DatabaseName string `yaml:"databaseName"`
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

	if !configuration.Notification.gotMarshalled {
		configuration.Notification.EnableIncidentsWithClosedTasksNotification = true
		configuration.Notification.EnableIncidentsWithoutOwnerNotification = true
		configuration.Notification.EnableTasksWithoutOwnerNotification = true
	}

	return configuration, err
}
