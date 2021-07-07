package util

import (
	"testing"
)

func TestFormatParam(t *testing.T) {
	s := FormatParam("/a/:b/c/:d")
	if s != `/a/-/c/-` {
		t.Error("s is not `/a/-/c/-` , but ", s)
		return
	}
	s = FormatParam("/")
	if s != `/` {
		t.Error("s is not `/` , but ", s)
		return
	}
	s = FormatParam("")
	if s != `` {
		t.Error("s is not `` , but ", s)
		return
	}
}

func TestMatchInParam(t *testing.T) {
	b := MatchInParam("/a/:b/c/:d", "/a/:f/c/:h")
	if b != true {
		t.Error("b is not true , but ", b)
		return
	}

	b = MatchInParam("/a/:b", "/e/:c/:d")
	if b != false {
		t.Error("b is not false , but ", b)
		return
	}
}

func TestMatchParam(t *testing.T) {
	b := MatchRoute("/a/:b/:c", "/a/asdih-1q/wkjl1-a")
	if b != true {
		t.Error("b is not true , but ", b)
		return
	}

	b = MatchRoute("/a/:b/:c", "/a/asdih-1q/")
	if b != false {
		t.Error("b is not false , but ", b)
		return
	}
}
