package foo

import (
	"fmt"

	"github.com/dsolerh/go-test-mono/utils"
)

func Version() string {
	return fmt.Sprintf("own: 0.1.4, utils: %s, utils.sub: %s", utils.Version(), utils.SubVersion())
}
