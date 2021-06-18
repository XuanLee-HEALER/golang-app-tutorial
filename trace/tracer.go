package trace

import (
	"fmt"
	"io"
)

type Tracer interface {
	Trace(...interface{})
}

type tracer struct {
	out io.Writer
}

func (t *tracer) Trace(a ...interface{}) {
	fmt.Fprint(t.out, a...)
	fmt.Fprintln(t.out)
}

// New 参数决定trace的内容输出到哪里，io.Writer是一个接口
func New(w io.Writer) Tracer {
	return &tracer{out: w}
}
