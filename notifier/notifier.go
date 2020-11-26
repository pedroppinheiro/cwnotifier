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

	notificationTitle   string = "Aviso de chamado prioritário!"
	notificationMessage string = "Há chamados que demandam sua atenção urgente!"

	errorNotificationTitle   string = "Erro!"
	errorNotificationMessage string = "Um erro ocorreu durante a execução do programa. Verifique o arquivo de log."

	programStartNotificationTitle   string = "CWNotifier started!"
	programStartNotificationMessage string = "CWNotifier has started running."
)

var notification toast.Notification = toast.Notification{
	AppID:    "CWNotifier",
	Title:    utf8toASCII(notificationTitle),
	Message:  utf8toASCII(notificationMessage),
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

func init() {
	// get the absolute path of the cherwell logo image to then present it in the notification
	currentDir, err := filepath.Abs(filepath.Dir(os.Args[0]))
	if err != nil {
		log.Panic(err)
	}

	cherwellLogoLocation := currentDir + "\\" + cherwellLogoName

	if fileExists(cherwellLogoLocation) {
		notification.Icon = cherwellLogoLocation
	} else {
		log.Printf("File \"%v\" was not found.\n", cherwellLogoLocation)
	}
}

// Notify emits the windows notification about a priority cherwell's task
func Notify() {
	err := notification.Push()
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
