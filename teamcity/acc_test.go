package teamcity_test

import (
	"strings"
)

func composeConfig(config ...string) string {
	var str strings.Builder

	for _, conf := range config {
		str.WriteString(conf)
	}

	return str.String()
}
