package slug

import (
	"regexp"
	"strconv"
	"strings"
)

// Транслит как в Nest (slug.utils.ts) — поведение slug совпадает с бэкендом Node.
var translitMap = map[rune]string{
	'а': "a", 'б': "b", 'в': "v", 'г': "g", 'д': "d", 'е': "e", 'ё': "e",
	'ж': "zh", 'з': "z", 'и': "i", 'й': "y", 'к': "k", 'л': "l", 'м': "m",
	'н': "n", 'о': "o", 'п': "p", 'р': "r", 'с': "s", 'т': "t", 'у': "u",
	'ф': "f", 'х': "h", 'ц': "ts", 'ч': "ch", 'ш': "sh", 'щ': "sch",
	'ъ': "", 'ы': "y", 'ь': "", 'э': "e", 'ю': "yu", 'я': "ya",
}

var slugPattern = regexp.MustCompile(`^[a-z0-9]+(?:-[a-z0-9]+)*$`)

// GenerateSlug — kebab-case из названия (кириллица → латиница).
func GenerateSlug(text string) string {
	var b strings.Builder
	for _, r := range strings.ToLower(strings.TrimSpace(text)) {
		if tr, ok := translitMap[r]; ok {
			b.WriteString(tr)
			continue
		}
		b.WriteRune(r)
	}
	s := b.String()
	var out strings.Builder
	lastDash := false
	for _, r := range s {
		if (r >= 'a' && r <= 'z') || (r >= '0' && r <= '9') {
			out.WriteRune(r)
			lastDash = false
			continue
		}
		if !lastDash && out.Len() > 0 {
			out.WriteRune('-')
			lastDash = true
		}
	}
	res := strings.Trim(out.String(), "-")
	res = regexp.MustCompile(`-+`).ReplaceAllString(res, "-")
	return res
}

// MakeUniqueSlug — добавляет -2, -3, … при коллизии.
func MakeUniqueSlug(base string, existing []string) string {
	slug := base
	n := 2
	set := make(map[string]struct{}, len(existing))
	for _, e := range existing {
		set[e] = struct{}{}
	}
	for {
		if _, ok := set[slug]; !ok {
			return slug
		}
		slug = base + "-" + strconv.Itoa(n)
		n++
	}
}

// IsValidSlug — как regex в DTO Nest.
func IsValidSlug(s string) bool {
	return slugPattern.MatchString(s)
}

// IsBlankOrSpace — пустое имя после trim.
func IsBlankOrSpace(s string) bool {
	return strings.TrimSpace(s) == ""
}
