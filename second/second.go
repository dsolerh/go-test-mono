package second

import "first"

func Second() bool {
	return !first.First()
}
