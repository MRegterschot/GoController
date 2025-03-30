package utils

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/MRegterschot/GoController/config"
	"github.com/MRegterschot/GoController/models"
)

var (
	delimiter string
	re        *regexp.Regexp
	theme     models.Theme
)

func SetTheme() {
	delimiter = config.AppEnv.Delimiter
	re = regexp.MustCompile(fmt.Sprintf(`%s([a-zA-Z0-9_]+)%s`, delimiter, delimiter))
	theme = config.Theme
}

// Processes a string and returns the modified string.
func ProcessString(str string) string {
	escapedInput := strings.ReplaceAll(str, `\#`, `\__ESCAPED_HASH__`)

	// Replace theme colors in the string
	modString := re.ReplaceAllStringFunc(escapedInput, func(match string) string {
		colorKey := strings.Trim(match, delimiter)
		for themeKey := range theme.Styling {
			if strings.EqualFold(colorKey, themeKey) {
				return "$" + theme.Styling[themeKey] // Replace with theme value
			}
		}
		return match
	})

	modString = strings.ReplaceAll(modString, `\__ESCAPED_HASH__`, `#`)

	return modString
}
