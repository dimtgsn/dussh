package logger

import (
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"os"
)

func MustNew(level, encoding string) *zap.Logger {
	at, err := zap.ParseAtomicLevel(level)
	if err != nil {
		panic(err)
	}

	//return New(at.Level())
	config := zap.Config{
		Level:            zap.NewAtomicLevelAt(at.Level()),
		Development:      true,
		Encoding:         encoding,
		EncoderConfig:    zap.NewProductionEncoderConfig(),
		OutputPaths:      []string{"stdout"},
		ErrorOutputPaths: []string{"stderr"},
	}

	config.EncoderConfig.EncodeCaller = zapcore.FullCallerEncoder

	logger, err := config.Build()
	if err != nil {
		panic(err)
	}
	defer logger.Sync()

	return logger
}

func New(ll zapcore.Level) *zap.Logger {
	encoderCfg := zap.NewProductionEncoderConfig()

	encoderCfg.EncodeTime = zapcore.RFC3339NanoTimeEncoder
	encoderCfg.EncodeDuration = zapcore.SecondsDurationEncoder

	l := zap.New(zapcore.NewCore(
		zapcore.NewJSONEncoder(encoderCfg),
		zapcore.Lock(os.Stdout),
		zap.NewAtomicLevelAt(ll),
	), zap.AddStacktrace(zap.PanicLevel))
	zap.ReplaceGlobals(l)

	return l
}

const ctxLogger = "loggerKey"

// ContextWithLogger adds logger to context.
func ContextWithLogger(c *gin.Context, l *zap.Logger) {
	c.Set(ctxLogger, l)
}

// LoggerFromContext returns logger from context
func LoggerFromContext(c *gin.Context) *zap.Logger {
	if l, ok := c.Value(ctxLogger).(*zap.Logger); ok {
		return l
	}
	return zap.L()
}
