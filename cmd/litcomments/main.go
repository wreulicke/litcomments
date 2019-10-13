package main

import (
	"litcomments"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() { singlechecker.Main(litcomments.Analyzer) }
