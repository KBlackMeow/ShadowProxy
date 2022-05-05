package logger

import (
	"fmt"
	"time"
)

var LogLevel int

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

func Log(info ...any) {
	if LogLevel > 0 {
		return
	}
	fmt.Print(TimeNow(), " [+] : ")
	for _, s := range info {
		fmt.Print(s, " ")
	}
	fmt.Println()
}

func Warn(info ...any) {
	if LogLevel > 1 {
		return
	}
	fmt.Print(TimeNow(), " [-] : ")
	for _, s := range info {
		fmt.Print(s, " ")
	}
	fmt.Println()
}

func Error(err ...any) {
	fmt.Print(TimeNow(), " [*] : ")
	for _, s := range err {
		fmt.Print(s, " ")
	}
	fmt.Println()
}
