package log

import (
	"os"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

func init() {
	ProvideLogger()
}

func ProvideLogger() {
	// ### Mostly copied from zap examples ###
	// The bundled Config struct only supports the most common configuration
	// options. More complex needs, like splitting logs between multiple files
	// or writing to non-file outputs, require use of the zapcore package.
	//
	// In this example, imagine we're both sending our logs to Kafka and writing
	// them to the console. We'd like to encode the console output and the Kafka
	// topics differently, and we'd also like special treatment for
	// high-priority logs.

	// First, define our level-handling logic.
	consoleLevel := zap.InfoLevel
	fileLevel := zap.DebugLevel

	file, _ := os.Create("arkwaifu-2x.log")

	// Assume that we have clients for two Kafka topics. The clients implement
	// zapcore.WriteSyncer and are safe for concurrent use. (If they only
	// implement io.Writer, we can use zapcore.AddSync to add a no-op Sync
	// method. If they're not safe for concurrent use, we can add a protecting
	// mutex with zapcore.Lock.)
	fileOutput := zapcore.Lock(file)

	// High-priority output should also go to standard error, and low-priority
	// output should also go to standard out.
	consoleOutput := zapcore.Lock(os.Stdout)

	// Optimize the Kafka output for machine consumption and the console output
	// for human operators.
	fileEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())
	consoleEncoder := zapcore.NewConsoleEncoder(zap.NewDevelopmentEncoderConfig())

	// Join the outputs, encoders, and level-handling functions into
	// zapcore.Cores, then tee the four cores together.
	core := zapcore.NewTee(
		zapcore.NewCore(fileEncoder, fileOutput, fileLevel),
		zapcore.NewCore(consoleEncoder, consoleOutput, consoleLevel),
	)

	// From a zapcore.Core, it's easy to construct a Logger.
	logger := zap.New(core)
	zap.ReplaceGlobals(logger)
}
