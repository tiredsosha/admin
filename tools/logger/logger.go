package logger

import (
	"io"
	"log"
	"os"
)

var (
	Warn  *log.Logger
	Debug *log.Logger
	Info  *log.Logger
	Error *log.Logger
)

func LogInit(debug bool) {
	var out interface{}
	out = io.Discard

	if debug {
		deleteLog := logCreation()
		if deleteLog {
			os.Remove("admin.log")
		}
		file, err := os.OpenFile("admin.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
		if err == nil {
			out = file
		}
	}

	Debug = log.New(out.(io.Writer), "DEBUG: ", log.Ldate|log.Ltime)
	Info = log.New(out.(io.Writer), "INFO:  ", log.Ldate|log.Ltime)
	Warn = log.New(out.(io.Writer), "WARN:  ", log.Ldate|log.Ltime)
	Error = log.New(out.(io.Writer), "ERROR: ", log.Ldate|log.Ltime)

	Debug.Println("")
	Debug.Println("")
	Info.Print("ADMIN STARTED")
}

func DebugLog(version string, debug bool, mqttOn, statusOn bool, hostname, broker, username, password string, port int) {
	Debug.Println("---------------------------")
	Debug.Println("common data:")
	Debug.Printf("\t\tversion    - %s\n", version)
	Debug.Println("- - - - - - - - - - - - - -")
	Debug.Println("http data:")
	Debug.Printf("\t\tweb port   - %d\n", port)
	Debug.Printf("\t\tstatus udp - %v\n", statusOn)
	Debug.Println("- - - - - - - - - - - - - -")
	Debug.Println("сonnection data:")
	Debug.Printf("\t\tmqtt on    - %v\n", mqttOn)
	Debug.Printf("\t\thostname   - %s\n", hostname)
	Debug.Printf("\t\tbroker     - %s\n", broker)
	Debug.Printf("\t\tusername   - %s\n", username)
	Debug.Printf("\t\tpassword   - %s\n", password)
	Debug.Println("---------------------------")
}

// старая проверка по времени создания, плохо работает, если файл первый раз создался больше года назад
// func logCreation() bool {
// 	var deleteLog bool = false

// 	if log, err := os.Stat("admin.log"); err == nil {
// 		createTime := time.Unix(0, log.Sys().(*syscall.Win32FileAttributeData).CreationTime.Nanoseconds())
// 		currTime := time.Now()
// 		diff := currTime.Sub(createTime).Milliseconds()
// 		// проверяем что лог не супер длинный и тяжелый.
// 		if diff > 604800000 {
// 			deleteLog = true
// 		}
// 	}
// 	return deleteLog
// }

func logCreation() bool {
	var deleteLog bool = false

	if log, err := os.Stat("admin.log"); err == nil {
		bytesSize := log.Size()
		// проверяем что лог не супер длинный и тяжелый.
		if bytesSize > 300000 {
			deleteLog = true
		}
	}
	return deleteLog
}
