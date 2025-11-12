package logx

import (
	"context"
	"log/slog"
	"slices"
	"strings"
)

type PrefixOptions struct {
	// A list of keys used to fetch prefix values from the log record.
	PrefixKeys []string
	// PrefixFormatter is a function to format the prefix of the log record.
	// If it's not set, the DefaultPrefixFormatter with ColorizePrefix wrapper is used.
	PrefixFormatter func(prefixes []Value) string
}

// prefixHandler is a custom slog handler that wraps another slog.prefixHandler to prepend a prefix to the log messages.
// The prefix is sourced from the log record's attributes using the keys specified in PrefixKeys.
type prefixHandler struct {
	Next     slog.Handler  // The next log handler in the chain.
	opts     PrefixOptions // Configuration options for this handler.
	prefixes []Value       // Cached list of prefix values.
}

// Prefix creates a new prefix logging handler.
// The new handler will prepend a prefix sourced from the log record's attributes to each log
// message before passing the record to the next handler.
func Prefix(next slog.Handler, opts *PrefixOptions) slog.Handler {
	if opts == nil {
		opts = &PrefixOptions{}
	}
	return &prefixHandler{
		Next:     next,
		opts:     *opts,
		prefixes: make([]Value, len(opts.PrefixKeys)),
	}
}

func (h *prefixHandler) Enabled(ctx context.Context, level Level) bool {
	return h.Next.Enabled(ctx, level)
}

// Handle processes a log record, prepending a prefix to its message if needed, and then passes the
// record to the next handler.
func (h *prefixHandler) Handle(ctx context.Context, r Record) error {
	prefixes := h.prefixes

	if r.NumAttrs() > 0 {
		nr := slog.NewRecord(r.Time, r.Level, r.Message, r.PC)
		attrs := make([]slog.Attr, 0, r.NumAttrs())
		r.Attrs(func(a slog.Attr) bool {
			attrs = append(attrs, a)
			return true
		})
		if p, changed := h.extractPrefixes(attrs); changed {
			nr.AddAttrs(attrs...)
			r = nr
			prefixes = p
		}
	}

	f := h.opts.PrefixFormatter
	if f == nil {
		f = h.formatPrefix
	}

	r.Message = f(prefixes) + r.Message

	return h.Next.Handle(ctx, r)
}

func (h *prefixHandler) WithAttrs(attrs []slog.Attr) slog.Handler {
	p, _ := h.extractPrefixes(attrs)
	return &prefixHandler{
		Next:     h.Next.WithAttrs(attrs),
		opts:     h.opts,
		prefixes: p,
	}
}

func (h *prefixHandler) WithGroup(name string) slog.Handler {
	return &prefixHandler{
		Next:     h.Next.WithGroup(name),
		opts:     h.opts,
		prefixes: h.prefixes,
	}
}

// extractPrefixes scans the attributes for keys specified in PrefixKeys.
// If found, their values are saved in a new prefix list.
// The original attribute list will be modified to remove the extracted prefix attributes.
func (h *prefixHandler) extractPrefixes(attrs []slog.Attr) (prefixes []Value, changed bool) {
	prefixes = h.prefixes
	for i, attr := range attrs {
		idx := slices.IndexFunc(h.opts.PrefixKeys, func(s string) bool { return s == attr.Key })
		if idx >= 0 {
			if !changed {
				// make a copy of prefixes:
				prefixes = make([]Value, len(h.prefixes))
				copy(prefixes, h.prefixes)
			}
			prefixes[idx] = attr.Value
			attrs[i] = slog.Attr{} // remove the prefix attribute
			changed = true
		}
	}
	return
}

// func(prefixes Value) string
func (*prefixHandler) formatPrefix(prefixes []Value) string {
	p := make([]string, 0, len(prefixes))
	for _, prefix := range prefixes {
		if prefix.Any() == nil || prefix.String() == "" {
			continue // skip empty prefixes
		}
		p = append(p, prefix.String())
	}

	if len(p) == 0 {
		return ""
	}

	n := `[` + strings.Join(p, ":") + `]`
	if c := 6 - len(n); c > 0 {
		n += strings.Repeat(" ", c)
	}
	return n + " " //" > "
}
