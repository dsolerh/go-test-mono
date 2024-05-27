package third

import (
	"fmt"

	"github.com/dsolerh/go-test-mono/first"
)

func Version() string {
	return fmt.Sprintf("third: v1 | %s", first.Version())
}
