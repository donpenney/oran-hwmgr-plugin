/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package logging

import (
	"context"
	"log/slog"
)

//
// This module includes utilities to define a slog handler that includes attributes
// that have been added to the context in order to carry info through an execution
// flow without needing to explicitly include it in all logs.
//

type loggingContextKey string

const (
	slogFields loggingContextKey = "slog_fields"
)

type LoggingContextHandler struct {
	handler slog.Handler
}

// Handle adds attributes from the context to the log record
func (h LoggingContextHandler) Handle(ctx context.Context, record slog.Record) error {
	if attrs, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		for _, v := range attrs {
			record.AddAttrs(v)
		}
	}

	return h.handler.Handle(ctx, record) // nolint: wrapcheck
}

func (h LoggingContextHandler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handler.Enabled(ctx, level)
}

func (h LoggingContextHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	if len(attrs) == 0 {
		return h
	}
	return LoggingContextHandler{handler: h.handler.WithAttrs(attrs)}
}

func (h LoggingContextHandler) WithGroup(name string) slog.Handler {
	if name == "" {
		return h
	}
	return LoggingContextHandler{handler: h.handler.WithGroup(name)}
}

func NewLoggingContextHandler() *LoggingContextHandler {
	return &LoggingContextHandler{
		handler: slog.Default().Handler(),
	}
}

// AppendCtx adds an slog attribute to the provided context so that it will be
// included in any Record created with such context
func AppendCtx(ctx context.Context, attr slog.Attr) context.Context {
	if ctx == nil {
		ctx = context.Background()
	}

	if v, ok := ctx.Value(slogFields).([]slog.Attr); ok {
		v = append(v, attr)
		return context.WithValue(ctx, slogFields, v)
	}

	v := []slog.Attr{}
	v = append(v, attr)
	return context.WithValue(ctx, slogFields, v)
}