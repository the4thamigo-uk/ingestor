package log

// duplicated from github.com/sirupsen/logrus

// Logger is what your logrus-enabled library should take, that way
// it'll accept a stdlib logger and a logrus logger. There's no standard
// interface, this is the closest we get, unfortunately.
type Logger interface {
	Print(...interface{})
	Printf(string, ...interface{})
	Println(...interface{})

	Fatal(...interface{})
	Fatalf(string, ...interface{})
	Fatalln(...interface{})

	Panic(...interface{})
	Panicf(string, ...interface{})
	Panicln(...interface{})
}

// Disabled is a no-op logger
func Disabled() Logger {
	return &noOp{}
}

type noOp struct{}

func (*noOp) Print(...interface{})          {}
func (*noOp) Printf(string, ...interface{}) {}
func (*noOp) Println(...interface{})        {}

func (*noOp) Fatal(...interface{})          {}
func (*noOp) Fatalf(string, ...interface{}) {}
func (*noOp) Fatalln(...interface{})        {}

func (*noOp) Panic(...interface{})          {}
func (*noOp) Panicf(string, ...interface{}) {}
func (*noOp) Panicln(...interface{})        {}
