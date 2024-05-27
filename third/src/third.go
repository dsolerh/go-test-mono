package third

import (
	"fmt"

	"github.com/dsolerh/go-test-mono/first"
)

func Version() string {
	return fmt.Sprintf("third: v0.1.2 | %s", first.Version())
}
