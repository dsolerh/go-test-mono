package second

import "github.com/dsolerh/go-test-mono/first"

func Second() bool {
	return !first.First()
}
