package litcomments_test

import (
	"testing"

	"litcomments"

	"golang.org/x/tools/go/analysis/analysistest"
)

func Test(t *testing.T) {
	testdata := analysistest.TestData()
	analysistest.Run(t, testdata, litcomments.Analyzer, "a")
}
