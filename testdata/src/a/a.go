package a

import (
	"context"
	"net/http"
)

func f(name []int) {}

func x(x bool) {}

func main() {
	f([]int{1, 2, 3})                                          // want "Composite literal without comments is found"
	r, _ := http.NewRequest("method" /* method */, "url", nil) // want "Basic literal without comments \"url\" is found" "Nil literal without comments is found."

	// false positive
	f([]int{1, 2, // want "Composite literal without comments is found"
		3} /* int */)

	r.Clone(nil) // want "Nil literal without comments is found."

	x(true)  // want "true literal without comments is found."
	x(false) // want "false literal without comments is found."

	if true == true {
		nil := context.TODO()
		r.Clone(nil)
		true := context.TODO()
		r.Clone(true)
		false := context.TODO()
		r.Clone(false)
	}
}
