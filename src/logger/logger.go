package outlogger

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"log/slog"
	"strconv"
	"sync"
)

// Constants to make logger have color
const (
	reset     = "\033[0m"
	red       = 31
	green     = 32
	yellow    = 33
	blue      = 34
	lightGray = 37
	darkGray  = 90
)

const (
	timeFormat = "[15:00:00.000]"
)

// Color ouput
func colorize(colorCode int, msg string) string {
	return fmt.Sprintf("\033[%sm%s%s", strconv.Itoa(colorCode), msg, reset)
}

// Handler structure
type Handler struct {
	handle slog.Handler
	bytes  *bytes.Buffer
	mutex  *sync.Mutex
}

func (h *Handler) Enabled(ctx context.Context, level slog.Level) bool {
	return h.handle.Enabled(ctx, level)
}

func (h *Handler) WithAttrs(attrs []slog.Attr) slog.Handler {
	return &Handler{handle: h.handle.WithAttrs(attrs), bytes: h.bytes, mutex: h.mutex}
}

func (h *Handler) WithGroup(name string) slog.Handler {
	return &Handler{handle: h.handle.WithGroup(name), bytes: h.bytes, mutex: h.mutex}
}

func (h *Handler) Handle(ctx context.Context, record slog.Record) error {
	level := record.Level.String() + ":"

	var color int

	switch record.Level {
	case slog.LevelDebug:
		color = blue
		level = colorize(blue, level)
	case slog.LevelInfo:
		color = green
		level = colorize(green, level)
	case slog.LevelWarn:
		color = yellow
		level = colorize(yellow, level)
	case slog.LevelError:
		color = red
		level = colorize(red, level)
	}

	attrs, err := h.computeAttrs(ctx, record)
	if err != nil {
		return err
	}

	bytes, err := json.MarshalIndent(attrs, "", " ")
	if err != nil {
		return err
	}

	fmt.Println(
		colorize(lightGray, record.Time.Format(timeFormat)),
		level,
		colorize(color, record.Message),
		colorize(darkGray, string(bytes)),
	)
	return nil
}

func NewHandler(opts *slog.HandlerOptions) *Handler {
	if opts == nil {
		opts = &slog.HandlerOptions{}
	}
	b := &bytes.Buffer{}
	return &Handler{
		handle: slog.NewJSONHandler(b, &slog.HandlerOptions{
			Level:       opts.Level,
			AddSource:   opts.AddSource,
			ReplaceAttr: suppressDefaults(opts.ReplaceAttr),
		}),
		bytes: b,
		mutex: &sync.Mutex{},
	}
}

func (h *Handler) computeAttrs(ctx context.Context, record slog.Record) (map[string]any, error) {
	h.mutex.Lock()
	defer func() {
		h.bytes.Reset()
		h.mutex.Unlock()
	}()

	if err := h.handle.Handle(ctx, record); err != nil {
		return nil, fmt.Errorf("error when calling inner handler's Handle result %w", err)
	}

	var attrs map[string]any
	err := json.Unmarshal(h.bytes.Bytes(), &attrs)
	if err != nil {
		return nil, fmt.Errorf("error when unmarshaling inner handler's Handle result %w", err)
	}

	return attrs, nil
}

func suppressDefaults(next func([]string, slog.Attr) slog.Attr) func([]string, slog.Attr) slog.Attr {
	return func(groups []string, a slog.Attr) slog.Attr {
		if a.Key == slog.TimeKey ||
			a.Key == slog.LevelKey ||
			a.Key == slog.MessageKey {
			return slog.Attr{}
		}
		if next == nil {
			return a
		}
		return next(groups, a)
	}
}
