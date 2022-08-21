package debug

import (
	"fmt"
	"log"
	"os"
	"strings"
)

type Debug func(format string, v ...interface{})

var hookGetEnv = func() string {
	return os.Getenv("DEBUG")
}

var hookPrint = func(input string) {
	log.Println(input)
}

func Init(flag string) Debug {
	enable := false

	env := hookGetEnv()
	parts := strings.Split(env, ",")
	for _, part := range parts {
		if part == flag {
			enable = true
			break
		}
	}

	return func(format string, v ...interface{}) {
		if enable {
			hookPrint(fmt.Sprintf(format, v...))
		}
	}
}
