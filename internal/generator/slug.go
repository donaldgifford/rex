package generator

import (
	"regexp"
	"strings"
	"unicode"
)

// Slugify converts a title into a URL-friendly slug
// Example: "Use SQLite as Cache Layer" -> "use-sqlite-as-cache-layer"
func Slugify(title string) string {
	// Convert to lowercase
	slug := strings.ToLower(title)

	// Replace spaces with hyphens
	slug = strings.ReplaceAll(slug, " ", "-")

	// Remove special characters, keep only alphanumeric and hyphens
	slug = removeNonAlphanumeric(slug)

	// Replace multiple consecutive hyphens with a single hyphen
	re := regexp.MustCompile(`-+`)
	slug = re.ReplaceAllString(slug, "-")

	// Trim hyphens from start and end
	slug = strings.Trim(slug, "-")

	return slug
}

// removeNonAlphanumeric removes all characters except alphanumeric and hyphens
func removeNonAlphanumeric(s string) string {
	var result strings.Builder
	for _, r := range s {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || r == '-' {
			result.WriteRune(r)
		}
	}
	return result.String()
}

// GenerateADRFilename generates a filename for an ADR
// Format: {number:04d}-{slug}.md
// Example: 0001-use-sqlite-as-cache-layer.md
func GenerateADRFilename(number int, title string) string {
	slug := Slugify(title)
	return slugWithNumber(number, slug)
}

// GenerateRFCFilename generates a filename for an RFC
// Format: {number:04d}-{slug}.md
// Example: 0001-rex-documentation-management-system.md
func GenerateRFCFilename(number int, title string) string {
	slug := Slugify(title)
	return slugWithNumber(number, slug)
}

// GenerateTaskFilename generates a filename for a Task
// Format: {id}-{slug}.md
// Example: TASK-001-implement-sqlite-cache.md
func GenerateTaskFilename(id, title string) string {
	slug := Slugify(title)
	return id + "-" + slug + ".md"
}

// GeneratePlanFilename generates a filename for a Plan
// Format: {number:04d}-{slug}.md
// Example: 0001-phase-1-implementation.md
func GeneratePlanFilename(number int, title string) string {
	slug := Slugify(title)
	return slugWithNumber(number, slug)
}

// slugWithNumber creates a filename with zero-padded number prefix
func slugWithNumber(number int, slug string) string {
	// Zero-pad to 4 digits
	var numberStr string
	if number < 10 {
		numberStr = "000" + string(rune('0'+number))
	} else if number < 100 {
		numberStr = "00" + string(rune('0'+(number/10))) + string(rune('0'+(number%10)))
	} else if number < 1000 {
		numberStr = "0" + string(rune('0'+(number/100))) + string(rune('0'+((number/10)%10))) + string(rune('0'+(number%10)))
	} else {
		// For numbers >= 1000, use the full number
		digits := []rune{}
		n := number
		for n > 0 {
			digits = append([]rune{rune('0' + (n % 10))}, digits...)
			n /= 10
		}
		numberStr = string(digits)
	}

	return numberStr + "-" + slug + ".md"
}
