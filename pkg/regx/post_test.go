package regx

import (
	"fmt"
	"testing"
)

func TestA(t *testing.T) {
	var args = []string{
		"t2_11qnzrqv.",
		"t2_11qnzrqv",
		"t2_11qnzrqvqwdw",
		"t1_11qnzrqv",
		"11qnzrqv",
	}

	for _, each := range args {
		fmt.Printf("%v -> %s \n", ValidAuthorName.Match([]byte(each)), each)
	}
}
