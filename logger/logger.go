package logger

import (
	"fmt"
	"os"
	"time"
)

var LogLevel int
var ConsoleOutput bool

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Log(info ...any) {
	if LogLevel > 0 {
		return
	}
	out := fmt.Sprint(TimeNow(), " [+] : ")
	for _, s := range info {
		out += fmt.Sprint(s, " ")
	}
	if ConsoleOutput {
		fmt.Println(out)
	} else {
		WriteFileln(out)
	}
}

func Warn(info ...any) {
	if LogLevel > 1 {
		return
	}
	out := fmt.Sprint(TimeNow(), " [-] : ")
	for _, s := range info {
		out += fmt.Sprint(s, " ")
	}
	if ConsoleOutput {
		fmt.Println(out)
	} else {
		WriteFileln(out)
	}
}

func Error(err ...any) {
	out := fmt.Sprint(TimeNow(), " [*] : ")
	for _, s := range err {
		out += fmt.Sprint(s, " ")
	}
	if ConsoleOutput {
		fmt.Println(out)
	} else {
		WriteFileln(out)
	}

}

func WriteFileln(s string) {
	s = s + "\n"
	logFile, err := os.OpenFile("shadowproxy.log", os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)
	if err != nil {
		panic(err)
	}

	logFile.WriteString(s)
}
