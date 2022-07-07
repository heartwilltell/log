package log

import (
	"fmt"
	"io"
	"log"
	"os"
	"strings"
)

// compilation time check for interface implementation.
var (
	_ Logger = (*StdLog)(nil)
	_ Logger = NopLog{}
)

const (
	// ERR represents error logging level.
	ERR Level = iota
	// WRN represents warning logging level.
	WRN
	// INF represents info logging level.
	INF
	// DBG represents debug logging level.
	DBG

	// Short represents format of the code line number.
	// Contains resulting file name and line number.
	Short LineNumFormat = iota
	// Long represents format of the code line number.
	// Contains full file path, name and line number.
	Long

	// ErrParseLevel indicates that string given to function ParseLevel can't be parsed to Level.
	ErrParseLevel Error = "string can't be parsed as Level, use: `error`, `warning`, `info`, `debug`"
)

// Logger formats the message according to standard format specifiers from the fmt package
// and writes the message to writer specified by the concrete interface implementation.
type Logger interface {
	// Error formats and writes the error level message.
	Error(format string, v ...any)
	// Warning formats and writes the warning level message.
	Warning(format string, v ...any)
	// Info formats and writes the information level message.
	Info(format string, v ...any)
	// Debug formats and writes the debug level message.
	Debug(format string, v ...any)
}

// Level represents an enumeration of logging levels.
type Level byte

func (l Level) String() string {
	return [4]string{
		"Error",
		"Warning",
		"Info",
		"Debug",
	}[l]
}

// LineNumFormat represents an enumeration of formats
// for printing the code line number in log output.
type LineNumFormat byte

// Error represents package level error related to logging work.
type Error string

func (e Error) Error() string { return string(e) }

// New returns a new instance of StdLog struct.
// Takes variadic options which will be applied to StdLog.
func New(options ...Option) *StdLog {
	l := &StdLog{
		err: log.New(os.Stderr, "\033[31mERR\033[0m: ", log.Ldate|log.Ltime),
		wrn: log.New(os.Stderr, "\033[33mWRN\033[0m: ", log.Ldate|log.Ltime),
		inf: log.New(os.Stderr, "\033[32mINF\033[0m: ", log.Ldate|log.Ltime),
		dbg: log.New(os.Stderr, "\033[35mDBG\033[0m: ", log.Ldate|log.Ltime),
		lvl: INF,
	}

	for _, option := range options {
		option(l)
	}

	return l
}

// NewStdLog returns a new instance of StdLog struct.
// Takes variadic options which will be applied to StdLog.
func NewStdLog(options ...Option) *StdLog { return New(options...) }

// StdLog represents wrapper around standard library logger
// which implements Logger interface.
type StdLog struct {
	err, wrn, inf, dbg *log.Logger
	lvl                Level
}

func (l *StdLog) Error(format string, v ...any) {
	if l.lvl < ERR {
		return
	}

	l.err.Printf(format, v...)
}

func (l *StdLog) Info(format string, v ...any) {
	if l.lvl < INF {
		return
	}

	l.inf.Printf(format, v...)
}

func (l *StdLog) Warning(format string, v ...any) {
	if l.lvl < WRN {
		return
	}

	l.wrn.Printf(format, v...)
}

func (l *StdLog) Debug(format string, v ...any) {
	if l.lvl < DBG {
		return
	}

	l.dbg.Printf(format, v...)
}

// Option represents a functional option type which can be
// passed to the NewStdLog function to change its underlying
// properties.
type Option func(l *StdLog)

// WithLevel changes the underlying logging level of StdLog to the given on.
func WithLevel(level Level) Option { return func(l *StdLog) { l.lvl = level } }

// WithWriter changes the writer for each leveled loggers of StdLog to the given on.
func WithWriter(w io.Writer) Option {
	return func(l *StdLog) {
		l.err.SetOutput(w)
		l.wrn.SetOutput(w)
		l.inf.SetOutput(w)
		l.dbg.SetOutput(w)
	}
}

// WithNoColor disables the color output of StdLog.
func WithNoColor() Option {
	return func(l *StdLog) {
		l.err.SetPrefix("ERR: ")
		l.wrn.SetPrefix("WRN: ")
		l.inf.SetPrefix("INF: ")
		l.dbg.SetPrefix("DBG: ")
	}
}

// WithLineNum enables the printing of code line number in log output.
func WithLineNum(format LineNumFormat) Option {
	return func(l *StdLog) {
		flags := l.inf.Flags() // because all flags are equal for each leveled logger

		switch format {
		case Short:
			flags |= log.Lshortfile
		case Long:
			flags |= log.Llongfile
		}

		l.err.SetFlags(flags)
		l.wrn.SetFlags(flags)
		l.inf.SetFlags(flags)
		l.dbg.SetFlags(flags)
	}
}

// WithUTC sets the log time format to UTC.
func WithUTC() Option {
	return func(l *StdLog) {
		flags := l.inf.Flags() | log.LUTC
		l.err.SetFlags(flags)
		l.wrn.SetFlags(flags)
		l.inf.SetFlags(flags)
		l.dbg.SetFlags(flags)
	}
}

// ParseLevel takes the string and tries to parse it to the Level.
func ParseLevel(lvl string) (Level, error) {
	if lvl == "" {
		return INF, ErrParseLevel
	}

	levels := map[string]Level{
		strings.ToLower(WRN.String()): WRN,
		strings.ToLower(ERR.String()): ERR,
		strings.ToLower(INF.String()): INF,
		strings.ToLower(DBG.String()): DBG,
	}

	level, ok := levels[strings.ToLower(lvl)]
	if !ok {
		return INF, fmt.Errorf("%s %w", lvl, ErrParseLevel)
	}

	return level, nil
}

// NopLog represents empty/disabled implementation of Logger interface.
type NopLog struct{}

// NewNopLog returns a new instance of NopLog.
func NewNopLog() NopLog { return NopLog{} }

func (l NopLog) Error(string, ...any)   {}
func (l NopLog) Warning(string, ...any) {}
func (l NopLog) Info(string, ...any)    {}
func (l NopLog) Debug(string, ...any)   {}
