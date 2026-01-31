/*
 * 项目名称：blog-backend
 * 文件名称：logger.go
 * 创建时间：2026-01-31 16:12:38
 *
 * 系统用户：Administrator
 * 作　　者：無以菱
 * 联系邮箱：huangjing510@126.com
 * 功能描述：日志管理模块，基于zap实现结构化日志，支持开发环境和生产环境的不同输出格式
 */
package logger

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

// Logger 全局日志实例
var Logger *zap.Logger

// InitLogger 初始化日志系统
// 功能说明：
//  1. 根据配置的日志级别设置zap日志级别
//  2. 根据环境（开发/生产）选择不同的编码器
//  3. 开发环境使用彩色控制台输出，便于调试
//  4. 生产环境使用JSON格式输出，便于日志收集和分析
//
// 参数:
//   - level: 日志级别，可选值：debug、info、warn、error，默认为info
//   - isDev: true表示开发环境（彩色输出），false表示生产环境（JSON格式）
//
// 返回:
//   - error: 初始化失败时返回错误
func InitLogger(level string, isDev bool) error {
	// 设置日志级别
	var zapLevel zapcore.Level
	switch level {
	case "debug":
		zapLevel = zapcore.DebugLevel
	case "info":
		zapLevel = zapcore.InfoLevel
	case "warn":
		zapLevel = zapcore.WarnLevel
	case "error":
		zapLevel = zapcore.ErrorLevel
	default:
		zapLevel = zapcore.InfoLevel // 默认使用Info级别
	}

	var encoder zapcore.Encoder
	isDev = false // 强制使用生产环境格式（JSON），可根据需要调整
	if isDev {
		// 开发环境：彩色输出，易读格式
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "T",             // 时间字段键名
			LevelKey:       "L",             // 级别字段键名
			NameKey:        "N",             // 日志器名称字段键名
			CallerKey:      "C",             // 调用者字段键名
			FunctionKey:    zapcore.OmitKey, // 省略函数名
			MessageKey:     "M",             // 消息字段键名
			StacktraceKey:  "S",             // 堆栈跟踪字段键名
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalColorLevelEncoder,                   // 彩色级别编码
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // 简洁时间格式
			EncodeDuration: zapcore.StringDurationEncoder,                      // 持续时间编码
			EncodeCaller:   zapcore.ShortCallerEncoder,                         // 短调用者编码
		}
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	} else {
		// 生产环境：JSON格式，便于日志收集和分析
		encoderConfig := zapcore.EncoderConfig{
			TimeKey:        "ts",            // 时间字段键名
			LevelKey:       "level",         // 级别字段键名
			NameKey:        "logger",        // 日志器名称字段键名
			CallerKey:      "caller",        // 调用者字段键名
			FunctionKey:    zapcore.OmitKey, // 省略函数名
			MessageKey:     "msg",           // 消息字段键名
			StacktraceKey:  "stacktrace",    // 堆栈跟踪字段键名
			LineEnding:     zapcore.DefaultLineEnding,
			EncodeLevel:    zapcore.CapitalLevelEncoder,                        // 大写级别编码
			EncodeTime:     zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05"), // 时间格式
			EncodeDuration: zapcore.StringDurationEncoder,                      // 持续时间编码
			EncodeCaller:   zapcore.ShortCallerEncoder,                         // 短调用者编码
		}
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建日志核心，将编码器、输出目标和日志级别组合
	core := zapcore.NewCore(
		encoder,
		zapcore.AddSync(os.Stdout), // 输出到标准输出
		zapLevel,                   // 日志级别
	)

	// 创建Logger实例，添加调用者信息，跳过一层调用栈（避免显示logger包本身的调用）
	Logger = zap.New(core, zap.AddCaller(), zap.AddCallerSkip(1))

	return nil
}

// Info 记录Info级别日志
// 参数:
//   - msg: 日志消息
//   - fields: 可选的附加字段（键值对）
func Info(msg string, fields ...zap.Field) {
	Logger.Info(msg, fields...)
}

// Debug 记录Debug级别日志
// 参数:
//   - msg: 日志消息
//   - fields: 可选的附加字段（键值对）
func Debug(msg string, fields ...zap.Field) {
	Logger.Debug(msg, fields...)
}

// Warn 记录Warn级别日志
// 参数:
//   - msg: 日志消息
//   - fields: 可选的附加字段（键值对）
func Warn(msg string, fields ...zap.Field) {
	Logger.Warn(msg, fields...)
}

// Error 记录Error级别日志
// 参数:
//   - msg: 日志消息
//   - fields: 可选的附加字段（键值对）
func Error(msg string, fields ...zap.Field) {
	Logger.Error(msg, fields...)
}

// Fatal 记录Fatal级别日志并退出程序
// 注意：调用此函数会导致程序立即退出
// 参数:
//   - msg: 日志消息
//   - fields: 可选的附加字段（键值对）
func Fatal(msg string, fields ...zap.Field) {
	Logger.Fatal(msg, fields...)
}

// Sync 同步日志缓冲区，确保所有日志都已写入
// 建议在程序退出前调用此函数，确保日志不丢失
func Sync() {
	Logger.Sync()
}
