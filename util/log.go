package util

import (
	"fmt"
	"time"
)

func init() {
	//file := "./httpcase" + ".log"
	//logFile, err := os.OpenFile(file, os.O_RDWR|os.O_CREATE, 0766)
	//if err != nil {
	//	panic(err)
	//}
	//log.SetOutput(logFile)
}

func Log(v ...interface{}) {
	vs := make([]interface{}, 0)
	vs = append(vs, time.Now().Format("2006/01/02 15:04:05.000"))
	vs = append(vs, v...)
	fmt.Println(vs...)
}

func Print(v ...interface{}) {
	fmt.Print(v...)
}

func Println(v ...interface{}) {
	fmt.Println(v...)
}
