package logger

import (
	"fmt"
	"log"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
}

func StringifyStruct(s any) string {
	return fmt.Sprintf("%#v", s)
}

func Info(msg string) {
	log.Println("[INFO] ", msg)
}

func Error(msg string, err error) {
	log.Println("[ERROR] ", msg)
	panic(err)
}
