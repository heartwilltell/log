package log

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

//nolint:gocognit
func TestNew(t *testing.T) {
	t.Run("New()", func(t *testing.T) {
		got := New()
		if got.lvl != INF {
			t.Errorf("default log lvl should be INF but got %s", got.lvl.String())
		}

		if got.err.Writer() != os.Stderr || got.inf.Writer() != os.Stderr || got.dbg.Writer() != os.Stderr {
			t.Errorf("all default writers should be os.Stderr")
		}

		reflectStdLog(t, got)
	})

	t.Run("New(WithWriter)", func(t *testing.T) {
		tb := &testWriter{byf: make([]byte, 0, 20)}
		strToLog := "test string"
		got := New(WithWriter(tb))
		got.Info(strToLog)

		if got.err.Writer() != tb || got.inf.Writer() != tb || got.dbg.Writer() != tb {
			t.Errorf("all default writers should be os.Stderr")
		}

		if !strings.Contains(tb.String(), strToLog) {
			t.Errorf("expected %s but got %s", strToLog, tb.String())
		}

		reflectStdLog(t, got)
	})

	t.Run("New(WithLevel)", func(t *testing.T) {
		got := New(WithLevel(ERR))
		if got.lvl != ERR {
			t.Errorf("log lvl should be := ERR got := %s", got.lvl.String())
		}

		reflectStdLog(t, got)
	})

	t.Run("New(WithNoColor)", func(t *testing.T) {
		got := New(WithNoColor())
		if got.err.Prefix() != "ERR: " {
			t.Errorf("Unexpected ERR prefix := %s, expcted := 'ERR: '", got.err.Prefix())
		}

		if got.wrn.Prefix() != "WRN: " {
			t.Errorf("Unexpected WRN prefix := %s, expcted := 'WRN: '", got.wrn.Prefix())
		}

		if got.inf.Prefix() != "INF: " {
			t.Errorf("Unexpected INF prefix := %s, expcted := 'INF: '", got.inf.Prefix())
		}

		if got.dbg.Prefix() != "DBG: " {
			t.Errorf("Unexpected DBG prefix := %s, expcted := 'DBG: '", got.dbg.Prefix())
		}

		reflectStdLog(t, got)
	})

	t.Run("New(WithUTC)", func(t *testing.T) {
		defaultFlags := log.Ldate | log.Ltime
		got := New(WithUTC())

		// Magic! I know. I also hate bitwise operators.
		if (got.err.Flags() ^ log.LUTC) != defaultFlags {
			t.Errorf("UTC flag hs not been set")
		}

		if (got.inf.Flags() ^ log.LUTC) != defaultFlags {
			t.Errorf("UTC flag hs not been set")
		}

		if (got.dbg.Flags() ^ log.LUTC) != defaultFlags {
			t.Errorf("UTC flag hs not been set")
		}

		if (got.wrn.Flags() ^ log.LUTC) != defaultFlags {
			t.Errorf("UTC flag hs not been set")
		}

		reflectStdLog(t, got)
	})

	t.Run("New(WithLevelAtPrefixEnd)", func(t *testing.T) {
		defaultFlags := log.Ldate | log.Ltime
		got := New(WithLevelAtPrefixEnd())

		// Magic! I know. I also hate bitwise operators.
		if (got.err.Flags() ^ log.Lmsgprefix) != defaultFlags {
			t.Errorf("LevelAtPrefixEnd flag hs not been set")
		}

		if (got.inf.Flags() ^ log.Lmsgprefix) != defaultFlags {
			t.Errorf("LevelAtPrefixEnd flag hs not been set")
		}

		if (got.dbg.Flags() ^ log.Lmsgprefix) != defaultFlags {
			t.Errorf("LevelAtPrefixEnd flag hs not been set")
		}

		if (got.wrn.Flags() ^ log.Lmsgprefix) != defaultFlags {
			t.Errorf("LevelAtPrefixEnd flag hs not been set")
		}

		reflectStdLog(t, got)
	})

	t.Run("New(WithLineNum(ShortFmt))", func(t *testing.T) {
		defaultFlags := log.Ldate | log.Ltime
		got := New(WithLineNum(ShortFmt))

		// Magic! I know. I also hate bitwise operators.
		if (got.err.Flags() ^ log.Lshortfile) != defaultFlags {
			t.Errorf("LineNum flag hs not been set")
		}

		reflectStdLog(t, got)
	})

	t.Run("New(WithLineNum(LongFmt))", func(t *testing.T) {
		defaultFlags := log.Ldate | log.Ltime
		got := New(WithLineNum(LongFmt))

		// Magic! I know. I also hate bitwise operators.
		if (got.err.Flags() ^ log.Llongfile) != defaultFlags {
			t.Errorf("LineNum flag hs not been set")
		}

		reflectStdLog(t, got)
	})
}

func TestNewStdLog(t *testing.T) {
	if !reflect.DeepEqual(NewStdLog(), New()) {
		t.Errorf("Constructor values mismatch")
	}
}

func TestLevelPrinting(t *testing.T) {
	type tcase struct{ level Level }

	tests := map[string]tcase{
		"DBG":  {level: DBG},
		"INF":  {level: INF},
		"WRN":  {level: WRN},
		"ERR":  {level: ERR},
		"ERR1": {level: Level(0)},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var w bytes.Buffer
			logger := New(WithLevel(tc.level), WithWriter(&w))
			logger.Debug("debug")
			logger.Info("info")
			logger.Warning("warning")
			logger.Error("error")

			lines := strings.Split(w.String(), "\n")
			lines = lines[0 : len(lines)-1]

			switch tc.level {
			case ERR:
				checkErrorLevelPrinting(t, lines, w)
			case WRN:
				checkWarningLevelPrinting(t, lines, w)
			case INF:
				checkInfoLevelPrinting(t, lines, w)
			case DBG:
				checkDebugLevelPrinting(t, lines, w)
			}
		})
	}
}

func checkErrorLevelPrinting(t *testing.T, lines []string, w bytes.Buffer) {
	if len(lines)-int(ERR) != 1 {
		t.Errorf("Unexpected number of printed lines: want := %d got := %d", ERR+1, len(lines)-int(ERR))
	}

	for _, l := range lines {
		if strings.HasPrefix(l, "\u001B[31mERR\u001B[0m:") {
			continue
		}

		t.Errorf("Unexpected printed line! Only ERR logs are allowed in this case. Got := %s", w.String())
	}
}

func checkWarningLevelPrinting(t *testing.T, lines []string, w bytes.Buffer) {
	if len(lines)-int(WRN) != 1 {
		t.Errorf("Unexpected number of printed lines: want := %d got := %d", WRN+1, len(lines)-int(WRN))
	}

	for _, l := range lines {
		if strings.HasPrefix(l, "\u001B[31mERR\u001B[0m:") ||
			strings.HasPrefix(l, "\u001B[33mWRN\u001B[0m:") {
			continue
		}

		t.Errorf("Unexpected printed line! Only ERR and WRN logs are allowed in this case. Got := %s", w.String())
	}
}

func checkInfoLevelPrinting(t *testing.T, lines []string, w bytes.Buffer) {
	if len(lines)-int(INF) != 1 {
		t.Errorf("Unexpected number of printed lines: want := %d got := %d", INF+1, len(lines)-int(INF))
	}

	for _, l := range lines {
		if strings.HasPrefix(l, "\u001B[31mERR\u001B[0m:") ||
			strings.HasPrefix(l, "\u001B[33mWRN\u001B[0m:") ||
			strings.HasPrefix(l, "\u001B[32mINF\u001B[0m:") {
			continue
		}

		t.Errorf("Unexpected printed line! Only ERR WRN INF logs are allowed in this case. Got := %s", w.String())
	}
}

func checkDebugLevelPrinting(t *testing.T, lines []string, w bytes.Buffer) {
	if len(lines)-int(DBG) != 1 {
		t.Errorf("Unexpected number of printed lines: want := %d got := %d", DBG+1, len(lines)-int(DBG))
	}

	for _, l := range lines {
		if strings.HasPrefix(l, "\u001B[31mERR\u001B[0m:") ||
			strings.HasPrefix(l, "\u001B[33mWRN\u001B[0m:") ||
			strings.HasPrefix(l, "\u001B[32mINF\u001B[0m:") ||
			strings.HasPrefix(l, "\u001B[35mDBG\u001B[0m:") {
			continue
		}

		t.Errorf("Unexpected printed line! Only ERR WRN INF logs are allowed in this case. Got := %s", w.String())
	}
}

func TestParseLevel(t *testing.T) {
	type tcase struct {
		str     string
		want    Level
		wantErr error
	}

	tests := map[string]tcase{
		"OK info level":       {str: "info", want: INF, wantErr: nil},
		"OK debug level":      {str: "debug", want: DBG, wantErr: nil},
		"OK warning level":    {str: "warning", want: WRN, wantErr: nil},
		"OK error level":      {str: "error", want: ERR, wantErr: nil},
		"OK uppercase":        {str: "INFO", want: INF, wantErr: nil},
		"Error empty string":  {str: "", want: INF, wantErr: ErrParseLevel},
		"Error invalid level": {str: "invalid-level", want: 0, wantErr: fmt.Errorf("invalid-level %w", ErrParseLevel)},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			got, err := ParseLevel(tc.str)
			if tc.wantErr != nil {
				if !reflect.DeepEqual(err, tc.wantErr) {
					t.Errorf("Error mismatch; got := %v; want := %v", err, tc.wantErr)
				}
			} else {
				if !reflect.DeepEqual(got, tc.want) {
					t.Errorf("Value mismatch; got := %v; want := %v", got, tc.want)
				}
			}
		})
	}
}

func TestLevel_String(t *testing.T) {
	type tcase struct {
		l    Level
		want string
	}

	tests := map[string]tcase{
		"Error":   {ERR, "Error"},
		"Info":    {INF, "Info"},
		"Warning": {WRN, "Warning"},
		"Debug":   {DBG, "Debug"},
		"Unknown": {Level(5), ""},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.l.String(); got != tc.want {
				t.Errorf("String() = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestNewNopLog(t *testing.T) {
	logger := NewNopLog()

	if !reflect.DeepEqual(logger, NopLog{}) {
		t.Errorf("Type mismatch; got := %v; want := %v", logger, NopLog{})
	}

	if !reflect.TypeOf(logger).Implements(reflect.TypeOf((*Logger)(nil)).Elem()) {
		t.Errorf("type does't implement logger.Logger interface")
	}

	if reflect.TypeOf(logger).Name() != "NopLog" {
		t.Errorf("type name should be NopLog but got: %s", reflect.TypeOf(logger).Name())
	}

	if reflect.TypeOf(logger).String() != "log.NopLog" {
		t.Errorf("struct type returned by () should have name *log.StdLog but got: %s", reflect.TypeOf(logger).String())
	}

	if reflect.TypeOf(logger).Kind() != reflect.Struct {
		t.Errorf("type kind returned by () should be struct but got: %s", reflect.TypeOf(logger).Kind())
	}
}

func TestNopLogPrinting(t *testing.T) {
	type tcase struct{ method func(string, ...any) }

	tests := map[string]tcase{
		"Info":    {method: NewNopLog().Info},
		"Error":   {method: NewNopLog().Error},
		"Warning": {method: NewNopLog().Warning},
		"Debug":   {method: NewNopLog().Debug},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) { tc.method("") })
	}
}

func reflectStdLog(t *testing.T, logger *StdLog) {
	t.Helper()

	if !reflect.TypeOf(logger).Implements(reflect.TypeOf((*Logger)(nil)).Elem()) {
		t.Errorf("type does't implement logger.Logger interface")
	}

	if reflect.TypeOf(*logger).Name() != "StdLog" {
		t.Errorf("type name should be StdLog but got: %s", reflect.TypeOf(logger).Name())
	}

	if reflect.TypeOf(logger).String() != "*log.StdLog" {
		t.Errorf("struct type returned by () should have name *log.StdLog but got: %s", reflect.TypeOf(logger).String())
	}

	if reflect.TypeOf(*logger).Kind() != reflect.Struct {
		t.Errorf("type kind returned by () should be struct but got: %s", reflect.TypeOf(logger).Kind())
	}
}

type testWriter struct {
	n   int
	byf []byte
}

func (w *testWriter) String() string {
	return string(w.byf)
}

func (w *testWriter) Write(p []byte) (n int, err error) {
	for i, b := range p {
		w.byf = append(w.byf, b)
		w.n = i
	}

	return w.n, nil
}
