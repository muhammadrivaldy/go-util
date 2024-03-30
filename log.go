package goutil

import (
	"context"
	"os"
	"runtime"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type logs struct {
	logger   *zap.Logger
	undo     func()
	telegram TeleService
}

// Logs is a interface of method logging
type Logs interface {
	Config(osFile *os.File, createOutput bool)
	Info(ctx context.Context, msg string, zapFields ...zapcore.Field)
	Warning(ctx context.Context, err error)
	Error(ctx context.Context, err error)
	Undo()
	Sync()
}

func encodeConfig(osFile *os.File, telegram TeleService, createOutput bool) (logs, error) {

	encoder := zap.NewProductionEncoderConfig()
	encoder.EncodeTime.UnmarshalText([]byte("ISO8601"))

	config := zap.NewProductionConfig()
	config.EncoderConfig = encoder
	config.DisableStacktrace = true

	if createOutput {
		config.OutputPaths = []string{osFile.Name(), os.Stdout.Name()}
		config.ErrorOutputPaths = config.OutputPaths
	}

	logger, err := config.Build()
	if err != nil {
		return logs{}, err
	}

	return logs{
		logger:   logger.WithOptions(zap.AddCallerSkip(1)),
		undo:     zap.RedirectStdLog(logger),
		telegram: telegram}, nil
}

// NewLog is a function for set up log information in this service
func NewLog(osFile *os.File, telegram TeleService, createOutput bool) (Logs, error) {

	logs, err := encodeConfig(osFile, telegram, createOutput)
	if err != nil {
		return nil, err
	}

	return &logs, nil
}

func (l *logs) Config(osFile *os.File, createOutput bool) {

	logs, err := encodeConfig(osFile, l.telegram, createOutput)
	if err != nil {
		return
	}

	l.Sync()
	l.Undo()

	l.logger = logs.logger
	l.undo = logs.undo
}

func (l *logs) Info(ctx context.Context, msg string, zapFields ...zapcore.Field) {

	zapcore := zapcoreField(ctx)
	if zapFields != nil {
		zapcore = append(zapFields, zapcore...)
	}

	if ctx != nil {
		l.logger.Info(msg, zapcore...)
		return
	}

	l.logger.Info(msg)
}

func (l *logs) Warning(ctx context.Context, err error) {

	// convert error to string
	msgError := err.Error()

	// log with info context
	if ctx != nil {
		l.logger.Warn(msgError, zapcoreField(ctx)...)
		return
	}

	// log without info context
	l.logger.Warn(msgError)
}

func (l *logs) Error(ctx context.Context, err error) {

	// convert error to string
	msgError := err.Error()

	// get path & line
	_, path, line, _ := runtime.Caller(1)

	// send error to telegram
	go l.telegram.SendError(ctx, path, line, msgError)

	// log with info context
	if ctx != nil {
		l.logger.Error(msgError, zapcoreField(ctx)...)
		return
	}

	// log without info context
	l.logger.Error(msgError)
}

func (l *logs) Undo() {
	l.undo()
}

func (l *logs) Sync() {
	l.logger.Sync()
}

func zapcoreField(ctx context.Context) (zapFields []zapcore.Field) {

	res := GetContext(ctx)
	zapFields = append(zapFields, zap.String("request-id", res.RequestID))
	zapFields = append(zapFields, zap.String("method", res.Method))
	zapFields = append(zapFields, zap.String("endpoint", res.Endpoint))

	return
}

// LogMiddleware ..
func LogMiddleware(logger Logs) func(c *gin.Context) {

	// override handler
	return func(c *gin.Context) {

		// parse gin.Context to context.Context
		ctx := ParseContext(c)

		// log info request
		logger.Info(ctx, "request")

		// next handler
		c.Next()

		// log info response
		logger.Info(ctx, "response", zap.Int("status", c.Writer.Status()))
	}
}
