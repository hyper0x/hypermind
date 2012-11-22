package utils

import (
	"log"
	"runtime"
	"fmt"
	"strconv"
	"strings"
)

func init() {
	log.SetFlags(log.LstdFlags)
}

func getInvokerLocation() string {
	pc, file, line, ok := runtime.Caller(2)
	if !ok {
		return ""
	}
	simpleFileName := ""
	if index := strings.LastIndex(file, "/"); index > 0 {
		simpleFileName = file[index + 1 : len(file)]
	}
	funcPath := ""
	funcPtr := runtime.FuncForPC(pc)
	if funcPtr != nil {
		funcPath = funcPtr.Name()
	}
	return fmt.Sprintf("%s : (%s:%s)", funcPath, simpleFileName, strconv.Itoa(line))
}

func LogError(v ...interface{}) {
	log.Printf("[ERROR] %s - %v", getInvokerLocation(), v)
}

func LogErrorf(format string, v ...interface{}) {
	log.Printf("[ERROR] %s - " + format, getInvokerLocation(), v)
}

func LogErrorln(v ...interface{}) {
	log.Printf("[ERROR] %s - %v\n", getInvokerLocation(), v)
}

func LogFatal(v ...interface{}) {
	log.Fatalf("[FATAL] %s - %v", getInvokerLocation(), v)
}

func LogFatalf(format string, v ...interface{}) {
	log.Fatalf("[FATAL] %s - " + format, getInvokerLocation(), v)
}

func LogFatalln(v ...interface{}) {
	log.Fatalf("[FATAL] %s - %v\n", getInvokerLocation(), v)
}

func LogInfo(v ...interface{}) {
	log.Printf("[INFO] %s - %v", getInvokerLocation(), v)
}

func LogInfof(format string, v ...interface{}) {
	log.Printf("[INFO] %s - " + format, getInvokerLocation(), v)
}

func LogInfoln(v ...interface{}) {
	log.Printf("[INFO] %s - %v\n", getInvokerLocation(), v)
}

func LogPanic(v ...interface{}) {
	log.Panicf("[PANIC] %s - %v", getInvokerLocation(), v)
}

func LogPanicf(format string, v ...interface{}) {
	log.Panicf("[PANIC] %s - " + format, getInvokerLocation(), v)
}

func LogPanicln(v ...interface{}) {
	log.Panicf("[PANIC] %s - %v\n", getInvokerLocation(), v)
}

func LogWarn(v ...interface{}) {
	log.Printf("[WARN] %s - %v", getInvokerLocation(), v)
}

func LogWarnf(format string, v ...interface{}) {
	log.Printf("[WARN] %s - " + format, getInvokerLocation(), v)
}

func LogWarnln(v ...interface{}) {
	log.Printf("[WARN] %s - %v\n", getInvokerLocation(), v)
}
