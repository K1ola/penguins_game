package helpers

import (
	"fmt"
)

func RecoverPanic() {
		if err := recover(); err != nil {
			fmt.Println("Recovered panic: ", err)
		}
		//close all connections, do cleaning up
}
