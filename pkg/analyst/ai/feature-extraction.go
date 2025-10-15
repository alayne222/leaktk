package ai

import (
	"math"
	"regexp"
	"strings"
)

type Features struct {
	Entropy                      float64
	LineHasKeyword               float64
	NumNumbers                   float64
	MatchHasKeyword              float64
	LineHasConsecutiveTrigrams   float64
	MatchHasConsecutiveTrigrams  float64
	SecretHasConsecutiveTrigrams float64
	NumSpecial                   float64
	SecretHasKeyword             float64
	LineHasRepeatingTrigrams     float64
	LineHasStopword              float64
	SecretLength                 float64
	SecretHasRepeatingTrigrams   float64
	MatchHasRepeatingTrigrams    float64
	SecretHasStopword            float64
	MatchHasStopword             float64
	SecretHasDictionaryWord      float64
}

func NewFeaturesPipeline(
	match string,
	secret string,
	keywords []string,
	stopwords []string,
	dictwords []string,
) *Features {
	f := &Features{}

	//line := input

	// Feature calculation
	f.Entropy = calculateEntropy(secret)
	//f.LineHasKeyword = hasAnyKeyword(line)
	f.NumNumbers = countNumbers(secret)
	f.MatchHasKeyword = hasAnyKeyword(match, keywords)
	//f.LineHasConsecutiveTrigrams = hasConsecutiveTrigrams(line)
	f.MatchHasConsecutiveTrigrams = hasConsecutiveTrigrams(match)
	f.SecretHasConsecutiveTrigrams = hasConsecutiveTrigrams(secret)
	f.NumSpecial = countSpecialChars(secret)
	f.SecretHasKeyword = hasAnyKeyword(secret, keywords)
	//f.LineHasRepeatingTrigrams = hasRepeatingTrigrams(line)
	//f.LineHasStopword = hasAnyStopword(line)
	f.SecretLength = float64(len(secret))
	f.SecretHasRepeatingTrigrams = hasRepeatingTrigrams(secret)
	f.MatchHasRepeatingTrigrams = hasRepeatingTrigrams(match)
	f.SecretHasStopword = hasAnyStopword(secret, stopwords)
	f.MatchHasStopword = hasAnyStopword(match, stopwords)
	//f.SecretHasDictionaryWord = hasDictionaryWord(secret)

	return f
}

func calculateEntropy(s string) float64 {
	if s == "" {
		return 0.0
	}
	freqs := make(map[rune]int)
	for _, r := range s {
		freqs[r]++
	}
	var entropy float64
	length := float64(len(s))
	for _, freq := range freqs {
		prob := float64(freq) / length
		entropy -= prob * math.Log2(prob)
	}
	return entropy
}

func countNumbers(s string) float64 {
	re := regexp.MustCompile(`[0-9]`)
	return float64(len(re.FindAllString(s, -1)))
}

// hasConsecutiveTrigrams checks for three identical consecutive characters.
func hasConsecutiveTrigrams(s string) float64 {
	for i := 0; i < len(s)-2; i++ {
		if s[i] == s[i+1] && s[i+1] == s[i+2] {
			return 1.0
		}
	}
	return 0.0
}

// countSpecialChars counts non-alphanumeric characters.
func countSpecialChars(s string) float64 {
	count := 0
	for _, r := range s {
		if !('a' <= r && r <= 'z' || 'A' <= r && r <= 'Z' || '0' <= r && r <= '9') {
			count++
		}
	}
	return float64(count)
}

// hasRepeatingTrigrams checks for any repeating trigram (3-character sequence).
func hasRepeatingTrigrams(s string) float64 {
	if len(s) < 6 {
		return 0.0
	}
	trigrams := make(map[string]bool)
	for i := 0; i < len(s)-2; i++ {
		trigram := s[i : i+3]
		if _, ok := trigrams[trigram]; ok {
			return 1.0
		}
		trigrams[trigram] = true
	}
	return 0.0
}

func hasAnyKeyword(s string, keywords []string) float64 {
	for _, keyword := range keywords {
		if strings.Contains(s, keyword) {
			return 1.0
		}
	}
	return 0.0
}

func hasAnyStopword(s string, stopwords []string) float64 {
	for _, stopword := range stopwords {
		if strings.Contains(s, stopword) {
			return 1.0
		}
	}
	return 0.0
}

// hasDictionaryWord checks if a string contains any word from a predefined dictionary.
// func hasDictionaryWord(s string) float64 {
// 	words := strings.Fields(strings.ToLower(s))
// 	for _, word := range words {
// 		if _, ok := dictionary[word]; ok {
// 			return 1.0
// 		}
// 	}
// 	return 0.0
// }
