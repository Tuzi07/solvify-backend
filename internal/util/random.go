package util

import (
	"math/rand"
	"time"
	"fmt"
	"strings"
)

const alphanum = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

func RandomInt(min, max int) int {
	return min + rand.Intn(max-min)
}

func RandomString(length int) string {
	var sb strings.Builder
	sb.Grow(length)
	for i := 0; i < length; i++ {
		sb.WriteByte(alphanum[rand.Intn(len(alphanum))])
	}
	return sb.String()
}

func RandomStatement() string {
	return fmt.Sprintf("Statement: %s", RandomString(RandomInt(20, 100)))
}

func RandomFeedback() string {
	return fmt.Sprintf("Feedback: %s", RandomString(RandomInt(20, 100)))
}

func RandomSubjectId() string {
	return fmt.Sprintf("SubjectId: %s", RandomString(RandomInt(10, 20)))
}

func RandomTopicId() string {
	return fmt.Sprintf("TopicId: %s", RandomString(RandomInt(10, 20)))
}

func RandomSubtopicId() string {
	return fmt.Sprintf("SubtopicId: %s", RandomString(RandomInt(10, 20)))
}

func RandomLevelOfEducation() string {
	levels := []string{"fundamental", "fundamental 2", "ensino medio", "ensino superior"}
	return levels[rand.Intn(len(levels))]
}

func RandomLanguage() string {
	langs := []string{"pt", "en", "es", "jp", "fr", "de"}
	return langs[rand.Intn(len(langs))]
}

func RandomCreatorID() string {
	return fmt.Sprintf("CreatorId: %s", RandomString(RandomInt(10, 20)))
}

func RandomBoolAnswer() *bool {
	return &[]bool{true, false}[rand.Intn(2)]
}

func RandomBoolAnswers(n int) []bool {
	answers := make([]bool, n)
	for i := range answers {
		answers[i] = *RandomBoolAnswer()
	}
	return answers
}


func RandomItems() []string {
	items := make([]string, RandomInt(2, 6))
	for i := range items {
		items[i] = RandomString(RandomInt(3, 8))
	}
	return items
}

func RandomCorrectItem(n int) *int {
	x := RandomInt(0, n)
	return &x
}
