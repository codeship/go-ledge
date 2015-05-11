package ledge

import (
	"fmt"
	"io"
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

func (u *unstructuredLogger) DebugWriter() io.Writer {
	return u.logger.DebugWriter(Debug(u.value("")))
}

func (u *unstructuredLogger) ErrorWriter() io.Writer {
	return u.logger.ErrorWriter(Error(u.value("")))
}

func (u *unstructuredLogger) InfoWriter() io.Writer {
	return u.logger.InfoWriter(Info(u.value("")))
}

func (u *unstructuredLogger) WarnWriter() io.Writer {
	return u.logger.WarnWriter(Warn(u.value("")))
}

func (u *unstructuredLogger) value(s string) string {
	if len(u.fields) == 0 {
		return s
	}
	fieldKeys := make([]string, len(u.fields))
	i := 0
	for key := range u.fields {
		fieldKeys[i] = key
		i++
	}
	sort.Sort(sort.StringSlice(fieldKeys))
	var slice []string
	for _, key := range fieldKeys {
		slice = append(slice, fmt.Sprintf("%s:%v", key, u.fields[key]))
	}
	fieldsString := strings.Join(slice, " ")
	if fieldsString != "" {
		fieldsString = fmt.Sprintf("{%s}", fieldsString)
	}
	if s == "" {
		return fmt.Sprintf("%s", fieldsString)
	}
	return fmt.Sprintf("%s %s", fieldsString, s)
}
