package foo

import (
	"fmt"

	"github.com/dsolerh/go-test-mono/packages/utils"
)

func Version() string {
	return fmt.Sprintf("own: 0.0.3, depends of utils: %s", utils.Version())
}
