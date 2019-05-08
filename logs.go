package main

import "log"

func LogMsg(msg ...interface{}) {
	log.Println(msg...)
}
