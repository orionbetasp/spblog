package util

import (
	"os"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

//Logger 日志
var Logger *zap.Logger

func init() {
	Logger = NewLogger(1, 0, 0, true)
}

//NewLogger 获取日志
// filePath 日志文件路径
// level 日志级别
// maxSize 每个日志文件保存的最大尺寸 单位：M
// maxBackups 日志文件最多保存多少个备份
// maxAge 文件最多保存多少天
// compress 是否压缩
// serviceName 服务名
func NewLogger(maxSize int, maxBackups int, maxAge int, compress bool) *zap.Logger {
	core := zapcore.NewTee(
		newCore("./log/info.log", zapcore.InfoLevel, maxSize, maxBackups, maxAge, compress),
		newCore("./log/err.log", zapcore.ErrorLevel, maxSize, maxBackups, maxAge, compress),
	)
	//core := newCore(filePath, level, maxSize, maxBackups, maxAge, compress)
	return zap.New(core, zap.AddCaller(), zap.Development())
}

/**
 * zapcore构造
 */
func newCore(filePath string, level zapcore.Level, maxSize int, maxBackups int, maxAge int, compress bool) zapcore.Core {
	//日志文件路径配置2
	hook := lumberjack.Logger{
		Filename:   filePath,   // 日志文件路径
		MaxSize:    maxSize,    // 每个日志文件保存的最大尺寸 单位：M
		MaxBackups: maxBackups, // 日志文件最多保存多少个备份
		MaxAge:     maxAge,     // 文件最多保存多少天
		Compress:   compress,   // 是否压缩
	}
	// 设置日志级别
	atomicLevel := zap.NewAtomicLevel()
	atomicLevel.SetLevel(level)
	//公用编码器
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "linenum",
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.LowercaseLevelEncoder,  // 小写编码器
		EncodeTime:     timeEncoder,                    //zapcore.ISO8601TimeEncoder,  ISO8601 UTC 时间格式
		EncodeDuration: zapcore.SecondsDurationEncoder, //
		EncodeCaller:   zapcore.ShortCallerEncoder,     //zapcore.FullCallerEncoder,  全路径编码器
		EncodeName:     zapcore.FullNameEncoder,
	}
	return zapcore.NewCore(
		zapcore.NewConsoleEncoder(encoderConfig), // 编码器配置
		//json形式
		//zapcore.NewJSONEncoder(encoderConfig)
		//友好形式
		//zapcore.NewConsoleEncoder(encoderConfig)
		zapcore.NewMultiWriteSyncer(zapcore.AddSync(os.Stdout), zapcore.AddSync(&hook)), // 打印到控制台和文件
		atomicLevel, // 日志级别
	)
}

func timeEncoder(t time.Time, enc zapcore.PrimitiveArrayEncoder) {
	enc.AppendString(t.Format("2006-01-02 15:04:05.000"))
}
