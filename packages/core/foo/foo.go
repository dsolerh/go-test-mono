package foo

import (
	"fmt"

	"github.com/dsolerh/go-test-mono/packages/utils"
)

func Version() string { return fmt.Sprintf("own: 0.0.9, utils: %s", utils.Version()) }
