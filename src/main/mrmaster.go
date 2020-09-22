package main

//
// start the master process, which is implemented
// in ../mr/master.go
//
// go run mrmaster.go pg*.txt
//

import (
	"fmt"
	"mr"
	"os"
	"time"
)

func main() {
	if len(os.Args) < 2 {
		fmt.Fprintf(os.Stderr, "Usage: mrmaster inputfiles...\n")
		os.Exit(1)
	}

	mr.DPrintf("Starting the master...\n")
	m := mr.MakeMaster(os.Args[1:], 10)
	for m.Done() == false {
		mr.DPrintf("Sleeping...\n")
		time.Sleep(time.Second)
	}

	time.Sleep(time.Second)
}
