package yoda_test

import (
	"os"
	"path"
	"strings"
	"testing"
	"github.com/vjmp/yoda"
)

type Nothing struct {}
type Example struct {Value int}

func panics(function yoda.Callable) (result bool) {
	result = false
	defer func() {
		if recover() != nil {
			result = true
		}
	}()
	function()
	return
}

func TestPanicOnWrongUse(be *testing.T) {
	yoda.True(true)
	if !panics(func() { yoda.True(true)}) {
		be.Error("Expecting panic, did not get it!")
	}
}

func TestCanUseTheTestingFramework(be *testing.T) {
	c := be
	c = nil
	a, b := 1, 2
	d, e := 3.14159, "hupsis"
	yoda.True(a < b).Must(be)
	yoda.True(b < a).Wont(be)
	yoda.Nil(a).Wont(be)
	yoda.Nil(d).Wont(be)
	yoda.Nil(e).Wont(be)
	yoda.Nil(be).Wont(be)
	yoda.Nil(c).Must(be)
	yoda.Nil(nil).Must(be)
	yoda.Equal(a, b).Wont(be)
	yoda.Equal(2, b).Must(be)
	yoda.Panic(func() {}).Wont(be)
	yoda.Panic(func() { panic("now") }).Must(be)
	yoda.Same(b, b).Must(be)
	yoda.Same(&Example{1}, &Example{2}).Wont(be)
	yoda.Same(&Example{1}, &Example{1}).Wont(be)
	yoda.Same(Example{1}, &Example{1}).Wont(be)
	yoda.Same(Example{1}, Example{1}).Must(be)
	yoda.Same(&Nothing{}, &Nothing{}).Must(be)
	yoda.Same(&Nothing{}, Nothing{}).Wont(be)
	yoda.Same(Nothing{}, Nothing{}).Must(be)
	yoda.Text("2", b).Must(be)
	yoda.Text("2", a).Wont(be)
	yoda.Text("<nil>", nil).Must(be)
	yoda.Text(nil, "<nil>").Must(be)
	yoda.Type("*testing.T", be).Must(be)
	yoda.All(positive_uints).Must(be)
	yoda.All(positive_ints).Wont(be)
}

func TestCanGetWorkingFolder(be *testing.T) {
	dir, err := os.Getwd()
	yoda.Nil(err).Must(be)
	yoda.Nil(dir).Wont(be)
	prefix, _ := path.Split(dir)
	yoda.Equal(prefix, dir).Wont(be)
	yoda.True(strings.HasPrefix(dir, prefix)).Must(be)
}

// some dummy propertis for yoda.All

func positive_uints(value uint64) bool {
	return value >= 0
}

func positive_ints(value int64) bool {
	return value >= 0
}

