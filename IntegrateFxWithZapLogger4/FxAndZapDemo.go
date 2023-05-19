package integratefxwithzaplogger4

import (
	"go.uber.org/fx"
	"go.uber.org/fx/fxevent"
	"go.uber.org/zap"
)

type (
	A struct{}
)

// see codes about how to use zap.Logger at https://github.com/hi-supergirl/go-micro-service-example/blob/master/dive-zapSugaredLogger/Demo.go
func createLogger() *zap.Logger {
	logger, err := zap.NewDevelopment()
	//logger, err := zap.NewProduction()
	if err != nil {
		logger = zap.NewNop()
	}
	return logger
}

func provider1(logger *zap.Logger) *A {
	logger.Sugar().Infoln("message", "this is a test for zap logger")
	return &A{}
}

func Main() {
	app := fx.New(
		fx.Supply(createLogger()),
		fx.WithLogger(func(logger *zap.Logger) fxevent.Logger {
			return &fxevent.ZapLogger{Logger: logger.Named("My Demo")}

		}),
		fx.Provide(
			provider1,
		),
		fx.Invoke(func(*A) {}),
	)
	app.Run()
}
