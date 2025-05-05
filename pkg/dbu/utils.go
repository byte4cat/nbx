package dbu

import (
	"unicode"
	"unicode/utf8"
)

// toSnakeCase converts a camelCase or PascalCase string to snake_case.
// Handles potential acronyms (e.g., "ID", "URL") better than simple conversion.
func toSnakeCase(str string) string {
	var output []rune
	var previous rune // Track the previous character's case/type

	for i, r := range str {
		// If it's an uppercase letter
		if unicode.IsUpper(r) {
			// Add underscore if:
			// 1. It's not the first character (prevents _Foo).
			// 2. The previous character was a lowercase letter (handles fooBar -> foo_bar).
			// 3. The previous character was uppercase, AND the current char is followed by a lowercase (handles userID -> user_id, but not URL -> url).
			if i > 0 && (unicode.IsLower(previous) || (unicode.IsUpper(previous) && i+1 < len(str) && unicode.IsLower([]rune(str)[i+1]))) {
				output = append(output, '_')
			}
			output = append(output, unicode.ToLower(r))
		} else {
			// Append non-uppercase characters directly
			output = append(output, r)
		}
		previous = r // Update previous character *after* processing r
	}
	return string(output)
}

// firstCharToLower converts the first character of a string to lowercase.
// It handles empty strings and single-rune strings.
// This is used as the default key naming fallback for BuildMongoUpdateMap.
func firstCharToLower(s string) string {
	if s == "" {
		return ""
	}
	r, size := utf8.DecodeRuneInString(s)
	if r == utf8.RuneError && size <= 1 {
		return s // Handle errors or empty/invalid runes
	}
	lc := unicode.ToLower(r)
	if r == lc {
		return s // First character is already lowercase or not a letter
	}
	return string(lc) + s[size:]
}
