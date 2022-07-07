package log

import (
	"bytes"
	"fmt"
	"os"
	"reflect"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	t.Run("NewStdLog()", func(t *testing.T) {
		got := NewStdLog()
		if got.lvl != INF {
			t.Errorf("default log lvl should be INF but got %s", got.lvl.String())
		}

		if got.err.Writer() != os.Stderr || got.inf.Writer() != os.Stderr || got.dbg.Writer() != os.Stderr {
			t.Errorf("all default writers should be os.Stderr")
		}

		if !reflect.TypeOf(got).Implements(reflect.TypeOf((*Logger)(nil)).Elem()) {
			t.Errorf("type does't implement logger.Logger interface")
		}

		if reflect.TypeOf(*got).Name() != "StdLog" {
			t.Errorf("type name should be StdLog but got: %s", reflect.TypeOf(*got).Name())
		}

		if reflect.TypeOf(got).String() != "*log.StdLog" {
			t.Errorf("struct type returned by NewStdLog() should have name *log.StdLog but got: %s", reflect.TypeOf(got).String())
		}

		if reflect.TypeOf(*got).Kind() != reflect.Struct {
			t.Errorf("type kind returned by NewStdLog() should be struct but got: %s", reflect.TypeOf(got).Kind())
		}
	})

	t.Run("NewStdLog(WithWriter)", func(t *testing.T) {
		tb := &testWriter{byf: make([]byte, 0, 20)}
		strToLog := "test string"
		got := NewStdLog(WithWriter(tb))
		got.Info(strToLog)

		if !strings.Contains(tb.String(), strToLog) {
			t.Errorf("expected %s but got %s", strToLog, tb.String())
		}

		if got.lvl != INF {
			t.Errorf("default log lvl should be INF but got %s", got.lvl.String())
		}
		if got.err.Writer() != tb || got.inf.Writer() != tb || got.dbg.Writer() != tb {
			t.Errorf("all default writers should be os.Stderr")
		}
	})

	t.Run("NewStdLog(WithLevel)", func(t *testing.T) {
		got := NewStdLog(WithLevel(ERR))
		if got.lvl != ERR {
			t.Errorf("default log lvl should be INF but got %s", got.lvl.String())
		}

		if got.err.Writer() != os.Stderr || got.inf.Writer() != os.Stderr || got.dbg.Writer() != os.Stderr {
			t.Errorf("all default writers should be os.Stderr")
		}

		if !reflect.TypeOf(got).Implements(reflect.TypeOf((*Logger)(nil)).Elem()) {
			t.Errorf("type does't implement logger.Logger interface")
		}

		if reflect.TypeOf(*got).Name() != "StdLog" {
			t.Errorf("type name should be StdLog but got: %s", reflect.TypeOf(*got).Name())
		}

		if reflect.TypeOf(got).String() != "*log.StdLog" {
			t.Errorf("struct type returned by NewStdLog() should have name *log.StdLog but got: %s", reflect.TypeOf(got).String())
		}

		if reflect.TypeOf(*got).Kind() != reflect.Struct {
			t.Errorf("type kind returned by NewStdLog() should be struct but got: %s", reflect.TypeOf(got).Kind())
		}
	})
}

func TestLevelPrinting(t *testing.T) {
	type tcase struct {
		level Level
	}

	tests := map[string]tcase{
		"DBG": {level: DBG},
		"INF": {level: INF},
		"WRN": {level: WRN},
		"ERR": {level: ERR},
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			var w bytes.Buffer
			log := New(WithLevel(tc.level), WithWriter(&w))
			log.Debug("debug")
			log.Info("info")
			log.Warning("warning")
			log.Error("error")

			lines := strings.Split(w.String(), "\n")
			lines = lines[0 : len(lines)-1]

			t.Log(lines)

			switch tc.level {
			case ERR:
				if len(lines)-int(ERR) != 1 {
					t.Errorf("Unexpected number of printed lines: want := %d got := %d", ERR+1, len(lines)-int(ERR))
				}

				for _, l := range lines {
					if strings.HasPrefix(l, "\u001B[31mERR\u001B[0m:") {
						continue
					}

					t.Errorf("Unexpected printed line! Only ERR logs are alowed in this case. Got := %s", w.String())
				}

			case WRN:
				if len(lines)-int(WRN) != 1 {
					t.Errorf("Unexpected number of printed lines: want := %d got := %d", WRN+1, len(lines)-int(WRN))
				}

				for _, l := range lines {
					if strings.HasPrefix(l, "\u001B[31mERR\u001B[0m:") ||
						strings.HasPrefix(l, "\u001B[33mWRN\u001B[0m:") {
						continue
					}

					t.Errorf("Unexpected printed line! Only ERR and WRN logs are alowed in this case. Got := %s", w.String())
				}
			case INF:
				if len(lines)-int(INF) != 1 {
					t.Errorf("Unexpected number of printed lines: want := %d got := %d", INF+1, len(lines)-int(INF))
				}

				for _, l := range lines {
					if strings.HasPrefix(l, "\u001B[31mERR\u001B[0m:") ||
						strings.HasPrefix(l, "\u001B[33mWRN\u001B[0m:") ||
						strings.HasPrefix(l, "\u001B[32mINF\u001B[0m:") {
						continue
					}

					t.Errorf("Unexpected printed line! Only ERR WRN INF logs are alowed in this case. Got := %s", w.String())
				}
			case DBG:
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

					t.Errorf("Unexpected printed line! Only ERR WRN INF logs are alowed in this case. Got := %s", w.String())
				}
			}
		})
	}

}

func TestNewNopLog(t *testing.T) {
	t.Run("Reflect type", func(t *testing.T) {
		want := reflect.TypeOf(NopLog{})
		got := reflect.TypeOf(NewNopLog())

		if !reflect.DeepEqual(want, got) {
			t.Errorf("Type mismatch; got := %v; want := %v", got, want)
		}
	})
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
	}

	for name, tc := range tests {
		t.Run(name, func(t *testing.T) {
			if got := tc.l.String(); got != tc.want {
				t.Errorf("String() = %v, want %v", got, tc.want)
			}
		})
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
