package xLog

import (
	"fmt"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
	"os"
	"strings"
	"time"
)

const (
	LevelDebug = iota //0
	LevelInfo
	LevelWarn
	LevelError

	DatetimeFormat             = "2006-01-02 15:04:05.000"
	DatetimeFormatWithTimezone = "2006-01-02 15:04:05.000 -0700"
)

var programInfo string

var loglevel int
var logger *zap.SugaredLogger

func MyCallerEncoder(caller zapcore.EntryCaller, enc zapcore.PrimitiveArrayEncoder) {
	path := caller.TrimmedPath()
	index := strings.LastIndexByte(path, '/')
	if index != -1 {
		startIndex := index + 1
		if len(path) > startIndex {
			path = path[startIndex:]
		}
	}

	str := programInfo + "[" + path + "]"
	enc.AppendString(str)
}

func Init(filePath string, maxSize, maxAge, maxBackups int, compress bool, level int, openTimeZone bool) {
	genProgramInfoStr()

	loglevel = level

	logFile := &lumberjack.Logger{
		Filename:   filePath,
		MaxSize:    maxSize,    // megabytes
		MaxAge:     maxAge,     // 0 表示不限制
		MaxBackups: maxBackups, // 0 表示不限制
		Compress:   compress,   // disabled by default
		LocalTime:  true,
	}
	w := zapcore.AddSync(logFile)

	// 自定义时间输出格式
	var customTimeEncoder zapcore.TimeEncoder
	if openTimeZone {
		customTimeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(DatetimeFormatWithTimezone))
		}
	} else {
		customTimeEncoder = func(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
			enc.AppendString(t.Format(DatetimeFormat))
		}
	}

	// 自定义日志级别显示
	customLevelEncoder := func(level zapcore.Level, enc zapcore.PrimitiveArrayEncoder) {
		levelStr := level.CapitalString()
		if len(levelStr) > 0 {
			levelStr = levelStr[0:1]
		}
		enc.AppendString("[" + levelStr + "]")
	}

	encoderConfig := zap.NewProductionEncoderConfig()
	encoderConfig.EncodeTime = customTimeEncoder
	encoderConfig.EncodeLevel = customLevelEncoder
	encoderConfig.EncodeCaller = MyCallerEncoder
	core := zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig),
		w,
		zapcore.DebugLevel,
	)
	caller := zap.AddCaller()
	callerSkip := zap.AddCallerSkip(1)
	zapLogger := zap.New(core, caller, callerSkip)
	logger = zapLogger.Sugar()
}

func genProgramInfoStr() {
	// procName := os.Args[0]
	// idx := strings.LastIndex(procName, "/")
	// programName := procName[idx+1:]

	programPid := os.Getpid()

	// programInfo = fmt.Sprintf("[%d][%s]", programPid, programName)
	programInfo = fmt.Sprintf("[%d]", programPid)
}

func SetLevel(level int) {
	loglevel = level
}

func GetLevel() int {
	return loglevel
}

func LogLevelToString(level int) string {
	switch level {
	case LevelError:
		{
			return "error"
		}
	case LevelWarn:
		{
			return "warn"
		}
	case LevelInfo:
		{
			return "info"
		}
	case LevelDebug:
		{
			return "debug"
		}
	default:
		return "unknown level"
	}
}

func ParseLogLevel(levelName string) int {
	switch strings.ToLower(levelName) {
	case "error":
		return LevelError
	case "warn":
		return LevelWarn
	case "info":
		return LevelInfo
	case "debug":
		return LevelDebug
	default:
		return LevelInfo
	}
}

func Debug(format string, args ...interface{}) {
	if loglevel > LevelDebug {
		return
	}

	if logger == nil {
		return
	}

	logger.Debugf(format, args...)
}

func Info(format string, args ...interface{}) {
	if loglevel > LevelInfo {
		return
	}

	if logger == nil {
		return
	}

	logger.Infof(format, args...)
}

func Warn(format string, args ...interface{}) {
	if loglevel > LevelWarn {
		return
	}

	if logger == nil {
		return
	}

	logger.Warnf(format, args...)
}

func Error(format string, args ...interface{}) {
	if loglevel > LevelError {
		return
	}

	if logger == nil {
		return
	}

	logger.Errorf(format, args...)
}
