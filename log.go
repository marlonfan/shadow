package shadow

import (
	"os"
	"log"
	"fmt"
)

var logger = log.New(os.Stderr, "", log.LstdFlags|log.Lshortfile|log.Lmsgprefix)

func logf(f string, v ...interface{}) {
	logger.Output(2, fmt.Sprintf(f, v...))
}

type logHelper struct {
	prefix string
}

func (l *logHelper) Write(p []byte) (n int, err error) {
	logger.Printf("%s%s\n", l.prefix, p)
	return len(p), nil
}

func newLogHelper(prefix string) *logHelper {
	return &logHelper{prefix}
}
