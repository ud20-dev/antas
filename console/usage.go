package console

import (
	"slices"
	"strings"
	"maps"
)

func GetReportersUsage() string{
	keys := slices.Sorted(maps.Keys(REPORTERS))
	return strings.Join(keys, "|")
}