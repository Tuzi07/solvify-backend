package api

import (
	"github.com/Tuzi07/solvify-backend/internal/util"
	"github.com/go-playground/validator/v10"
)

var validLanguage validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if language, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsSupportedLanguage(language)
	}
	return false
}

var validLanguages validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if languages, ok := fieldLevel.Field().Interface().([]string); ok {
		for _, language := range languages {
			if !util.IsSupportedLanguage(language) {
				return false
			}
		}
		return true
	}
	return false
}

var validFieldToOrderProblems validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if field, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsFieldToOrderProblems(field)
	}
	return false
}

var validLevelOfEducation validator.Func = func(fieldLevel validator.FieldLevel) bool {
	if level, ok := fieldLevel.Field().Interface().(string); ok {
		return util.IsLevelOfEducation(level)
	}
	return false
}
