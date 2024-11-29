package main

import (
	"os"

	standaloneListener "github.com/jcastellanos/challenge_transactions/internal/challenge/ports/input"
)

func main() {

	runtime := os.Getenv("RUNTIME")

	if runtime == "lambda" {

	} else {
		folder := "./transactions"
		standaloneListener := standaloneListener.NewListener(folder)
		standaloneListener.Run()
	}
}
