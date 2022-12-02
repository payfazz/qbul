package main

import (
	"github.com/payfazz/qbul/cmd/qbul-check/qbulcheck"

	"golang.org/x/tools/go/analysis/singlechecker"
)

func main() {
	singlechecker.Main(qbulcheck.Analizyer)
}
