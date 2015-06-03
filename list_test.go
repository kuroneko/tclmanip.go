package tclmanip

import (
	"testing"
)

func TestCurlySplit(t *testing.T) {
	tclList := TclList("ok {1 2 3} {4 5 6}")
	split := tclList.Split()

	if len(split) != 3 {
		t.Errorf("Split into wrong number of pieces")
	}
	if split[0] != "ok" {
		t.Errorf("Unexpected first section: got \"%s\"", split[0])
	}
	if split[1] != "1 2 3" {
		t.Errorf("Unexpected second section: got \"%s\"", split[1])
	}
}

func TestNestedCurlySplit(t *testing.T) {
	tclList := TclList("ok {1 {2 3}} {4 5 6}")
	split := tclList.Split()

	if len(split) != 3 {
		t.Errorf("Split into wrong number of pieces")
	}
	if split[0] != "ok" {
		t.Errorf("Unexpected first section: got \"%s\"", split[0])
	}
	if split[1] != "1 {2 3}" {
		t.Fatalf("Unexpected second section: got \"%s\"", split[1])
	}

	subsplit := split[1].Split()
	if len(subsplit) != 2 {
		t.Errorf("Subsplit split into wrong number of pieces")
	}
}

func TestSet(t *testing.T) {
	tclList := TclList("ok 1 2 3")
	if tclList.Index(3).String() != "3" {
		t.Fatalf("Index(3) returned unexpcted value (sanity check)")
	}
	tclList.Set(3, "Hello World")
	if tclList.String() != "ok 1 2 {Hello World}" {
		t.Errorf("String() gives unexpected value after set: %s", string(tclList))
	}
	if tclList.Index(3).String() != "Hello World" {
		t.Errorf("Index(3) gives unexpected value after set", tclList.Index(3))
	}
}