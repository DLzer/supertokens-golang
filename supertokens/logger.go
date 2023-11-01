package supertokens

import (
	"fmt"
	"log"
	"os"
	"runtime"
	"time"
)

const supertokens_namespace = "com.supertokens"

/*
 The debug logger below can be used to log debug messages in the following format
    com.supertokens {t: "2022-03-21T17:10:42+05:30", message: "Test Message", file: "/home/supertokens-golang/supertokens/supertokens.go:51" sdkVer: "0.5.2"}
*/

var (
	Logger       = log.New(os.Stdout, supertokens_namespace, 0)
	debugEnabled = false
)

func formatMessage(message string) string {
	_, file, line, _ := runtime.Caller(2)
	return fmt.Sprintf(" {t: \"%s\", message: \"%s\", file: \"%s:%d\" sdkVer: \"%s\"}\n\n", time.Now().Format(time.RFC3339), message, file, line, VERSION)
}

func LogDebugMessage(message string) {
	_, exists := os.LookupEnv("SUPERTOKENS_DEBUG")
	if exists || debugEnabled == true {
		Logger.Printf(formatMessage(message))
	}
}
