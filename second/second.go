package second

import (
	"fmt"

	"github.com/dsolerh/go-test-mono/first"
)

func Version() string {
	return fmt.Sprintf("second: v1 %s", first.Version())
}
