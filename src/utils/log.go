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
	var params []interface {}
	params = append(params, "[ERROR]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Print(params...)
}

func LogErrorf(format string, v ...interface{}) {
	var params []interface {}
	params = append(params, getInvokerLocation())
	params = append(params, v...)
	log.Printf(("[ERROR] %s - " + format), params...)
}

func LogErrorln(v ...interface{}) {
	var params []interface {}
	params = append(params, "[ERROR]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Println(params...)
}

func LogFatal(v ...interface{}) {
	var params []interface {}
	params = append(params, "[FATAL]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Fatal(params...)
}

func LogFatalf(format string, v ...interface{}) {
	var params []interface {}
	params = append(params, getInvokerLocation())
	params = append(params, v...)
	log.Fatalf("[FATAL] %s - " + format, params...)
}

func LogFatalln(v ...interface{}) {
	var params []interface {}
	params = append(params, "[FATAL]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Fatalln(params...)
}

func LogInfo(v ...interface{}) {
	var params []interface {}
	params = append(params, "[INFO]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Print(params...)
}

func LogInfof(format string, v ...interface{}) {
	var params []interface {}
	params = append(params, getInvokerLocation())
	params = append(params, v...)
	log.Printf("[INFO] %s - " + format, params...)
}

func LogInfoln(v ...interface{}) {
	var params []interface {}
	params = append(params, "[INFO]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Println(params...)
}

func LogPanic(v ...interface{}) {
	var params []interface {}
	params = append(params, "[PANIC]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Panic(params...)
}

func LogPanicf(format string, v ...interface{}) {
	var params []interface {}
	params = append(params, getInvokerLocation())
	params = append(params, v...)
	log.Panicf("[PANIC] %s - " + format, params...)
}

func LogPanicln(v ...interface{}) {
	var params []interface {}
	params = append(params, "[PANIC]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Panicln(params...)
}

func LogWarn(v ...interface{}) {
	var params []interface {}
	params = append(params, "[WARN]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Print(params...)
}

func LogWarnf(format string, v ...interface{}) {
	var params []interface {}
	params = append(params, getInvokerLocation())
	params = append(params, v...)
	log.Printf("[WARN] %s - " + format, params...)
}

func LogWarnln(v ...interface{}) {
	var params []interface {}
	params = append(params, "[WARN]")
	params = append(params, getInvokerLocation())
	params = append(params, "-")
	params = append(params, v...)
	log.Println(params...)
}
