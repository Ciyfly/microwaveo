package logger

import (
	"fmt"
	"log"
	"os"
)

var mlog *log.Logger

func Init() {
	mlog = log.New(os.Stdout, "[mcrowaveo] ", log.LstdFlags)
}

func Printf(format string, args ...interface{}) {
	mlog.Println(fmt.Sprintf(format, args...))
}

func Print(p string) {
	mlog.Println(p)
}

func Fatalf(format string, args ...interface{}) {
	mlog.Fatal(fmt.Sprintf(format, args...))
}

func Fatal(f string) {
	mlog.Fatal(f)
}
