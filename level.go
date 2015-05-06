package ledge

import (
	"errors"
	"fmt"
)

var (
	// LevelDebug is used for debugging.
	LevelDebug Level
	// LevelInfo is used for general informational statements.
	LevelInfo Level = 1
	// LevelWarn is used for statements that require attention but are not errors.
	LevelWarn Level = 2
	// LevelError is used for non-critical errors.
	LevelError Level = 3
	// LevelFatal is used for fatal errors. After logging, os.Exit(1) will be called.
	// os.Exit(1) will be called even if the PanicFilter is set.
	LevelFatal Level = 4
	// LevelPanic will panic with the given statement.
	LevelPanic Level = 5

	levelToString = map[Level]string{
		LevelDebug: "debug",
		LevelInfo:  "info",
		LevelWarn:  "warn",
		LevelError: "error",
		LevelFatal: "fatal",
		LevelPanic: "panic",
	}
	lenLevelToString = len(levelToString)
	stringToLevel    = map[string]Level{
		"debug": LevelDebug,
		"info":  LevelInfo,
		"warn":  LevelWarn,
		"error": LevelError,
		"fatal": LevelFatal,
		"panic": LevelPanic,
	}
)

// Level represents the logging level.
type Level uint

// LevelOf returns the Level for the given string representation.
func LevelOf(s string) (Level, error) {
	level, ok := stringToLevel[s]
	if !ok {
		return 0, errors.New(unknownLevel(s))
	}
	return level, nil
}

// String returns the string representation of the Level.
func (l Level) String() string {
	if int(l) < lenLevelToString {
		return levelToString[l]
	}
	panic(unknownLevel(l))
}

func unknownLevel(unknownLevel interface{}) string {
	return fmt.Sprintf("ledge: unknown Level: %v", unknownLevel)
}
