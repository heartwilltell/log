package log

import (
	"os"
	"reflect"
	"strings"
	"testing"
)

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
