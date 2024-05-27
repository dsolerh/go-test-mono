package second

import (
	"fmt"

	"github.com/dsolerh/go-test-mono/first"
)

func Version() string {
	return fmt.Sprintf("second: v0.1.2 %s", first.Version())
}
