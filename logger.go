package loggerLib

import (
	"context"
	"fmt"
	"log"
	"os"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	"go.opentelemetry.io/otel/trace"
)

type Logger struct {
	tracer      trace.Tracer
	infoLogger  *log.Logger
	errorLogger *log.Logger
}

func NewLogger(serviceName string) (*Logger, error) {
	ctx := context.Background()

	// OTLP exporter
	exp, err := otlptracegrpc.New(ctx)
	if err != nil {
		return nil, err
	}

	res, err := resource.New(
		ctx,
		resource.WithAttributes(
			attribute.String("service.name", serviceName),
		),
	)
	if err != nil {
		return nil, err
	}

	tp := sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(res),
	)
	otel.SetTracerProvider(tp)
	tracer := tp.Tracer(serviceName)

	// File loggers
	infoFile, err := os.OpenFile("info.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open info.log: %w", err)
	}

	errorFile, err := os.OpenFile("error.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		return nil, fmt.Errorf("failed to open error.log: %w", err)
	}

	return &Logger{
		tracer:      tracer,
		infoLogger:  log.New(infoFile, "INFO: ", log.Ldate|log.Ltime|log.Lshortfile),
		errorLogger: log.New(errorFile, "ERROR: ", log.Ldate|log.Ltime|log.Lshortfile),
	}, nil
}

func (l *Logger) Info(ctx context.Context, msg string) {
	l.infoLogger.Println(msg)
	_, span := l.tracer.Start(ctx, "Info")
	span.AddEvent(msg)
	span.End()
}

func (l *Logger) Error(ctx context.Context, msg string) {
	l.errorLogger.Println(msg)
	_, span := l.tracer.Start(ctx, "Error")
	span.AddEvent(msg)
	span.End()
}

func (l *Logger) Infof(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Info(ctx, msg)
}

func (l *Logger) Errorf(ctx context.Context, format string, v ...interface{}) {
	msg := fmt.Sprintf(format, v...)
	l.Error(ctx, msg)
}
