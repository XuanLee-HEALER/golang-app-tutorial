package trace

/*
	go test -cover 可以得出有多少条语句在测试范围内
*/

import (
	"bytes"
	"testing"
)

func TestNew(t *testing.T) {
	var buf bytes.Buffer
	trace := New(&buf)
	// 红绿灯判断，有意义的测试
	if trace == nil {
		t.Error("return from New should not be null")
	} else {
		trace.Trace("Hello trace package")
		if buf.String() != "Hello trace package\n" {
			t.Errorf("trace should not write: %s", buf.String())
		}
	}
}

func TestOff(t *testing.T) {
	silentTracer := Off()
	silentTracer.Trace("something")
}
