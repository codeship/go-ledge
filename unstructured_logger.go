package ledge

import (
	"fmt"
	"sort"
	"strings"
)

type unstructuredLogger struct {
	logger Logger
	fields map[string]interface{}
}

func newUnstructuredLogger(
	logger Logger,
	fields map[string]interface{},
) *unstructuredLogger {
	return &unstructuredLogger{
		logger,
		fields,
	}
}

func (u *unstructuredLogger) WithField(key string, value interface{}) UnstructuredLogger {
	return u.WithFields(Fields{key: value})
}

func (u *unstructuredLogger) WithFields(fields Fields) UnstructuredLogger {
	newFields := make(map[string]interface{})
	for key, value := range u.fields {
		newFields[key] = value
	}
	for key, value := range fields {
		newFields[key] = value
	}
	return newUnstructuredLogger(
		u.logger,
		newFields,
	)
}

func (u *unstructuredLogger) Debug(args ...interface{}) {
	u.logger.Debug(Debug(u.value(fmt.Sprint(args...))))
}

func (u *unstructuredLogger) Debugf(format string, args ...interface{}) {
	u.logger.Debug(Debug(u.value(fmt.Sprintf(format, args...))))
}

func (u *unstructuredLogger) Debugln(args ...interface{}) {
	u.logger.Debug(Debug(u.value(fmt.Sprintln(args...))))
}

func (u *unstructuredLogger) Error(args ...interface{}) {
	u.logger.Error(Error(u.value(fmt.Sprint(args...))))
}

func (u *unstructuredLogger) Errorf(format string, args ...interface{}) {
	u.logger.Error(Error(u.value(fmt.Sprintf(format, args...))))
}

func (u *unstructuredLogger) Errorln(args ...interface{}) {
	u.logger.Error(Error(u.value(fmt.Sprintln(args...))))
}

func (u *unstructuredLogger) Fatal(args ...interface{}) {
	u.logger.Fatal(Fatal(u.value(fmt.Sprint(args...))))
}

func (u *unstructuredLogger) Fatalf(format string, args ...interface{}) {
	u.logger.Fatal(Fatal(u.value(fmt.Sprintf(format, args...))))
}

func (u *unstructuredLogger) Fatalln(args ...interface{}) {
	u.logger.Fatal(Fatal(u.value(fmt.Sprintln(args...))))
}

func (u *unstructuredLogger) Info(args ...interface{}) {
	u.logger.Info(Info(u.value(fmt.Sprint(args...))))
}

func (u *unstructuredLogger) Infof(format string, args ...interface{}) {
	u.logger.Info(Info(u.value(fmt.Sprintf(format, args...))))
}

func (u *unstructuredLogger) Infoln(args ...interface{}) {
	u.logger.Info(Info(u.value(fmt.Sprintln(args...))))
}

func (u *unstructuredLogger) Panic(args ...interface{}) {
	u.logger.Panic(Panic(u.value(fmt.Sprint(args...))))
}

func (u *unstructuredLogger) Panicf(format string, args ...interface{}) {
	u.logger.Panic(Panic(u.value(fmt.Sprintf(format, args...))))
}

func (u *unstructuredLogger) Panicln(args ...interface{}) {
	u.logger.Panic(Panic(u.value(fmt.Sprintln(args...))))
}

func (u *unstructuredLogger) Print(args ...interface{}) {
	u.logger.Info(Info(u.value(fmt.Sprint(args...))))
}

func (u *unstructuredLogger) Printf(format string, args ...interface{}) {
	u.logger.Info(Info(u.value(fmt.Sprintf(format, args...))))
}

func (u *unstructuredLogger) Println(args ...interface{}) {
	u.logger.Info(Info(u.value(fmt.Sprintln(args...))))
}

func (u *unstructuredLogger) Warn(args ...interface{}) {
	u.logger.Warn(Warn(u.value(fmt.Sprint(args...))))
}

func (u *unstructuredLogger) Warnf(format string, args ...interface{}) {
	u.logger.Warn(Warn(u.value(fmt.Sprintf(format, args...))))
}

func (u *unstructuredLogger) Warnln(args ...interface{}) {
	u.logger.Warn(Warn(u.value(fmt.Sprintln(args...))))
}

func (u *unstructuredLogger) value(s string) string {
	if len(u.fields) == 0 {
		return s
	}
	fieldKeys := make([]string, len(u.fields))
	i := 0
	for key, _ := range u.fields {
		fieldKeys[i] = key
		i++
	}
	sort.Sort(sort.StringSlice(fieldKeys))
	var slice []string
	for _, key := range fieldKeys {
		slice = append(slice, fmt.Sprintf("%s:%v", key, u.fields[key]))
	}
	return fmt.Sprintf("{%s} %s", strings.Join(slice, " "), s)
}
