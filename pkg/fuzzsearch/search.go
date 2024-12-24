package fuzzsearch

import (
	"strings"
)

func Search(source []string, substring string) (string, bool) { 
	for _, item := range source {
		index := strings.Index(
			item, substring,
		)
		if (index != -1) {
			return item, true
		}
	}
	return "", false
}
