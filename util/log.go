package util

import (
	"io"
	"os"
	"path/filepath"
	"time"

	rotatelogs "github.com/lestrrat-go/file-rotatelogs"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

const (
	Minute = "minute"
	Hour   = "hour"
	Day    = "day"
	Month  = "month"
	Year   = "year"

	TimeDivision = "time"
	SizeDivision = "size"
)

var (
	Logger                    *zap.Logger
	_encoderNameToConstructor = map[string]func(zapcore.EncoderConfig) zapcore.Encoder{
		"console": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewConsoleEncoder(encoderConfig)
		},
		"json": func(encoderConfig zapcore.EncoderConfig) zapcore.Encoder {
			return zapcore.NewJSONEncoder(encoderConfig)
		},
	}
)

type timeUnit string

type logOptions struct {
	Encoding      string // 输出格式 "json" 或者 "console"
	InfoFilename  string // info级别日志文件名
	ErrorFilename string // warn级别日志文件名

	Division string //归档方式

	TimeUnit timeUnit // 时间归档 切割单位

	MaxSize    int  // 每个日志文件保存的最大尺寸 单位：M
	MaxBackups int  // 日志文件最多保存多少个备份
	MaxAge     int  // 文件最多保存多少天
	Compress   bool // 是否压缩

	LevelSeparate bool //是否日志分级
	stdoutDisplay bool //是否在控制台输出
	caller        bool //是否输出文件行号
}

func init() {
	c := newLog()

	//c.Division = "time"
	//c.TimeUnit = Day
	//c.MaxAge = 0

	c.Division = "size"
	c.MaxSize = 1
	c.MaxAge = 0
	c.Compress = true
	c.MaxBackups = 0

	c.Encoding = "console"
	c.InfoFilename = "./log/info.log"
	c.ErrorFilename = "./log/err.log"
	c.initLogger()
}

func newLog() *logOptions {
	return &logOptions{
		Encoding:      "console",
		InfoFilename:  "./log/log.log",
		ErrorFilename: "./log/err.log",
		Division:      "size",
		TimeUnit:      Day,
		MaxSize:       1,
		MaxBackups:    0,
		MaxAge:        0,
		Compress:      true,
		LevelSeparate: true,
		stdoutDisplay: true,
		caller:        true,
	}
}

func (c *logOptions) initLogger() {
	var (
		core               zapcore.Core
		infoHook, warnHook io.Writer
		wsInfo, wsWarn     []zapcore.WriteSyncer
	)

	if c.Encoding == "" {
		c.Encoding = "console"
	}
	encoder := _encoderNameToConstructor[c.Encoding]

	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "file", //"linenum"
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder, // 小写编码器
		EncodeTime:     timeEncoder,                   //zapcore.ISO8601TimeEncoder,  ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder, //zapcore.FullCallerEncoder,  全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}

	if c.stdoutDisplay {
		wsInfo = append(wsInfo, zapcore.AddSync(os.Stdout))
		wsWarn = append(wsWarn, zapcore.AddSync(os.Stdout))
	}

	// zapcore WriteSyncer setting
	if c.InfoFilename != "" {
		switch c.Division {
		case TimeDivision:
			err := os.MkdirAll(filepath.Dir(c.InfoFilename), 0744)
			if err != nil {
				panic("can't make directories for new logfile")
			}
			infoHook = c.timeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				err := os.MkdirAll(filepath.Dir(c.ErrorFilename), 0744)
				if err != nil {
					panic("can't make directories for new logfile")
				}
				warnHook = c.timeDivisionWriter(c.ErrorFilename)
			}
		case SizeDivision:
			infoHook = c.sizeDivisionWriter(c.InfoFilename)
			if c.LevelSeparate {
				warnHook = c.sizeDivisionWriter(c.ErrorFilename)
			}
		}
		wsInfo = append(wsInfo, zapcore.AddSync(infoHook))
	}

	if c.ErrorFilename != "" {
		wsWarn = append(wsWarn, zapcore.AddSync(warnHook))
	}

	// Separate info and warning log
	if c.LevelSeparate {
		core = zapcore.NewTee(
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), infoLevel()),
			zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsWarn...), warnLevel()),
		)
	} else {
		core = zapcore.NewCore(encoder(encoderConfig), zapcore.NewMultiWriteSyncer(wsInfo...), zap.InfoLevel)
	}

	// file line number display
	development := zap.Development()
	stackTrace := zap.AddStacktrace(zapcore.WarnLevel)
	// init default key
	//filed := zap.Fields(zap.String("serviceName", "serviceName"))
	if c.caller {
		Logger = zap.New(core, zap.AddCaller(), development, stackTrace)
	} else {
		Logger = zap.New(core, development, stackTrace)
	}
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}

func infoLevel() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl < zapcore.WarnLevel
	}
}

func warnLevel() zap.LevelEnablerFunc {
	return func(lvl zapcore.Level) bool {
		return lvl >= zapcore.WarnLevel
	}
}

func (c *logOptions) sizeDivisionWriter(filename string) io.Writer {
	hook := &lumberjack.Logger{
		Filename:   filename,
		MaxSize:    c.MaxSize,
		MaxBackups: c.MaxBackups,
		MaxAge:     c.MaxAge,
		Compress:   c.Compress,
	}
	return hook
}

func (c *logOptions) timeDivisionWriter(filename string) io.Writer {
	hook, err := rotatelogs.New(
		filename+c.TimeUnit.format(),
		rotatelogs.WithMaxAge(time.Duration(int64(24*time.Hour)*int64(c.MaxAge))),
		rotatelogs.WithRotationTime(c.TimeUnit.rotationGap()),
	)

	if err != nil {
		panic(err)
	}
	return hook
}

func (t timeUnit) format() string {
	switch t {
	case Minute:
		return ".%Y%m%d%H%M"
	case Hour:
		return ".%Y%m%d%H"
	case Day:
		return ".%Y%m%d"
	case Month:
		return ".%Y%m"
	case Year:
		return ".%Y"
	default:
		return ".%Y%m%d"
	}
}

func (t timeUnit) rotationGap() time.Duration {
	switch t {
	case Minute:
		return time.Minute
	case Hour:
		return time.Hour
	case Day:
		return time.Hour * 24
	case Month:
		return time.Hour * 24 * 30
	case Year:
		return time.Hour * 24 * 365
	default:
		return time.Hour * 24
	}
}
