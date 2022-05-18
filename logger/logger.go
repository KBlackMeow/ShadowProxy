package logger

import (
	"fmt"
	"os"
	"shadowproxy/config"
	"sync"
	"time"
)

func TimeNow() string {
	return time.Now().Format("2006-01-02 15:04:05")
}

var Mutex = new(sync.Mutex)

func Log(info ...any) {
	if config.ShadowProxyConfig.LogLevel > 0 {
		return
	}
	Mutex.Lock()
	defer Mutex.Unlock()
	out := fmt.Sprint(TimeNow(), " [+] : ")
	for _, s := range info {
		out += fmt.Sprint(s, " ")
	}
	if config.ShadowProxyConfig.ConsoleOutput {
		fmt.Println(out)
	} else {
		WriteFileln(out)
	}
}

func Warn(info ...any) {
	if config.ShadowProxyConfig.LogLevel > 1 {
		return
	}
	Mutex.Lock()
	defer Mutex.Unlock()
	out := fmt.Sprint(TimeNow(), " [-] : ")
	for _, s := range info {
		out += fmt.Sprint(s, " ")
	}
	if config.ShadowProxyConfig.ConsoleOutput {
		fmt.Println(out)
	} else {
		WriteFileln(out)
	}
}

func Error(err ...any) {
	Mutex.Lock()
	defer Mutex.Unlock()
	out := fmt.Sprint(TimeNow(), " [*] : ")
	for _, s := range err {
		out += fmt.Sprint(s, " ")
	}
	if config.ShadowProxyConfig.ConsoleOutput {
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
