package util

const (
	en = "en"
	pt = "pt"
)

func IsSupportedLanguage(language string) bool {
	switch language {
	case en, pt:
		return true
	}
	return false
}
