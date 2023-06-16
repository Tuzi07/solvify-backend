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

func IsFieldToOrderProblems(field string) bool {
	switch field {
	case "created_at", "attempts", "accuracy", "upvotes":
		return true
	}
	return false
}
