package utils

import (
	"fmt"
	"os"
)

const DefaultErrorExitCode = 1

func CheckErr(err error) {
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(DefaultErrorExitCode)
	}
}
