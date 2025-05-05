package log

import (
	"os"
	"time"
)

func Info(message string) {
	record(message, "info")
}

func Error(message string) {
	record(message, "error")
}

func Fatal(message string) {
	record(message, "fatal")
}

func record(message string, typeMessage string) {
	currentDate := time.Now()

	dirNameForLog := currentDate.Format("02_06_2006")
	fileNameForLog := currentDate.Format("02_06_2006_3_4")

	message = currentDate.String() + " " + message + "\n"

	dirs := map[int]map[string]string{
		0: {
			"file": "log",
			"type": "dir",
		},
		1: {
			"file": typeMessage,
			"type": "dir",
		},
		2: {
			"file": dirNameForLog,
			"type": "dir",
		},
		3: {
			"file": fileNameForLog,
			"type": "file",
		},
	}

	var pathToFile string

	for i := 0; i < len(dirs); i++ {
		fileName := dirs[i]["file"]
		fileType := dirs[i]["type"]

		file, err := os.Open(fileName)

		if !os.IsExist(err) {
			switch fileType {
			case "dir":
				err = os.Mkdir(pathToFile+fileName, 644)
				break
			case "file":
				file, err = os.OpenFile(pathToFile+fileName+".txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 644)
				break
			}
		}

		if fileType == "file" {
			file.WriteString(message)
		} else {
			pathToFile += fileName + "/"
		}
	}
}
