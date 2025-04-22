package publisher

import "github.com/Masterminds/semver/v3"

func SemverUpdater(vtype string) func(string) string {
	switch vtype {
	case "mayor":
		return func(s string) string {
			v, _ := semver.NewVersion(s)
			return "v" + v.IncMajor().String()
		}
	case "minor":
		return func(s string) string {
			v, _ := semver.NewVersion(s)
			return "v" + v.IncMinor().String()
		}
	case "patch":
		return func(s string) string {
			v, _ := semver.NewVersion(s)
			return "v" + v.IncPatch().String()
		}
	}
	panic("invalid version type")
}
