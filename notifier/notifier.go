package notifier

import (
	"log"
	"unicode/utf8"

	"gopkg.in/toast.v1"

	"os"
	"path/filepath"
)

const (
	cherwellLogoName string = "assets\\cherwell.png"

	incidentsWithoutOwnerNotificationTitle   string = "Aviso de chamado prioritário sem responsável"
	incidentsWithoutOwnerNotificationMessage string = "Há chamados no backlog que demandam sua atenção urgente!"

	tasksWithoutOwnerNotificationTitle   string = "Aviso de tarefa prioritária sem responsável"
	tasksWithoutOwnerNotificationMessage string = "Há tarefas no backlog que demandam sua atenção urgente!"

	incidentsWithClosedTasksNotificationTitle   string = "Aviso de chamado prioritário apto a encerrar"
	incidentsWithClosedTasksNotificationMessage string = "Há chamados prioritários que já podem ser encerrados!"

	noNotificationsEnabledTitle   string = "Nenhum notificação habilitada"
	noNotificationsEnabledMessage string = "O programa está encerrando pois nenhuma notificação está habilitada. Por favor habilite no arquivo de configuração"

	errorNotificationTitle   string = "Erro!"
	errorNotificationMessage string = "Um erro ocorreu durante a execução e o programa foi encerrado. Verifique o arquivo de log."

	programStartNotificationTitle   string = "CWNotifier started!"
	programStartNotificationMessage string = "CWNotifier has started running."
)

var incidentsWithoutOwnerNotification toast.Notification = toast.Notification{
	AppID:    "CWNotifier",
	Title:    utf8toASCII(incidentsWithoutOwnerNotificationTitle),
	Message:  utf8toASCII(incidentsWithoutOwnerNotificationMessage),
	Duration: "short",
}

var tasksWithoutOwnerNotification toast.Notification = toast.Notification{
	AppID:    "CWNotifier",
	Title:    utf8toASCII(tasksWithoutOwnerNotificationTitle),
	Message:  utf8toASCII(tasksWithoutOwnerNotificationMessage),
	Duration: "short",
}

var incidentsWithClosedTasksNotification toast.Notification = toast.Notification{
	AppID:    "CWNotifier",
	Title:    utf8toASCII(incidentsWithClosedTasksNotificationTitle),
	Message:  utf8toASCII(incidentsWithClosedTasksNotificationMessage),
	Duration: "short",
}

var startNotification toast.Notification = toast.Notification{
	AppID:    "CWNotifier",
	Title:    utf8toASCII(programStartNotificationTitle),
	Message:  utf8toASCII(programStartNotificationMessage),
	Duration: "short",
}

var errorNotification toast.Notification = toast.Notification{
	AppID:    "CWNotifier",
	Title:    utf8toASCII(errorNotificationTitle),
	Message:  utf8toASCII(errorNotificationMessage),
	Duration: "short",
}

var noNotificationsEnabledNotification toast.Notification = toast.Notification{
	AppID:    "CWNotifier",
	Title:    utf8toASCII(noNotificationsEnabledTitle),
	Message:  utf8toASCII(noNotificationsEnabledMessage),
	Duration: "short",
}

func init() {
	// get the absolute path of the cherwell logo image to then present it in the notification
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}

	cherwellLogoLocation := currentDir + "\\" + cherwellLogoName

	if fileExists(cherwellLogoLocation) {
		incidentsWithoutOwnerNotification.Icon = cherwellLogoLocation
		tasksWithoutOwnerNotification.Icon = cherwellLogoLocation
		incidentsWithClosedTasksNotification.Icon = cherwellLogoLocation
	} else {
		log.Printf("File \"%v\" was not found.\n", cherwellLogoLocation)
	}
}

// NotifyIncidentsWithoutOwnerNotification emits the windows notification about a priority cherwell's incident
func NotifyIncidentsWithoutOwnerNotification() {
	err := incidentsWithoutOwnerNotification.Push()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Notification emitted.")
}

// NotifyTasksWithoutOwnerNotification emits the windows notification about a priority cherwell's incident
func NotifyTasksWithoutOwnerNotification() {
	err := tasksWithoutOwnerNotification.Push()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Notification emitted.")
}

// NotifyIncidentsWithClosedTasksNotification emits the windows notification about a priority cherwell's incident
func NotifyIncidentsWithClosedTasksNotification() {
	err := incidentsWithClosedTasksNotification.Push()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Notification emitted.")
}

// NotifyProgramStart emits the windows notification about the start of the program
func NotifyProgramStart() {
	err := startNotification.Push()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Program start notification emitted.")
}

// NotifyError emits the windows notification about an error that occurred in the program
func NotifyError() {
	err := errorNotification.Push()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Error notification emitted.")
}

// NotifyNoNotificationsEnabled emits the windows notification about being no notifications enabled
func NotifyNoNotificationsEnabled() {
	err := noNotificationsEnabledNotification.Push()
	if err != nil {
		log.Panic(err)
	}

	log.Printf("no notifications enabled notification emitted.")
}

// fileExists checks if a file exists and is not a directory before we
// try using it to prevent further errors. Code from: https://golangcode.com/check-if-a-file-exists/
func fileExists(filename string) bool {
	fileInfo, err := os.Stat(filename)
	if os.IsNotExist(err) {
		return false
	}

	fileIsDir := fileInfo.IsDir()
	if fileIsDir {
		log.Printf("\"%v\" found, but is a directory.\n", cherwellLogoName)
		return false
	}

	return true
}

// utf8toASCII converts a UTF-8 internal string representation to standard
// ASCII bytes. Code from: https://gist.github.com/jbuchbinder/5513891
// This function is needed because windows notifications do not deal with UTF-8
func utf8toASCII(s string) string {
	t := make([]byte, utf8.RuneCountInString(s))
	i := 0
	for _, r := range s {
		t[i] = byte(r)
		i++
	}
	return string(t)
}
