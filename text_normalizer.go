package main

import (
	"strconv"
	"strings"
	"unicode"
)

func NormalizeText(input string) string {
	puncts := ".,!?:;'"
	for _, p := range puncts {
		input = strings.ReplaceAll(input, string(p), " "+string(p)+" ")
	}

	words := strings.Fields(input)
	processedWords := []string{}

	// Pass 1: Handle markers (hex), (bin), (up), (low), (cap)
	for i := 0; i < len(words); i++ {
		word := words[i]

		if word == "(hex)" {
			if len(processedWords) > 0 {
				lastIdx := len(processedWords) - 1
				val, err := strconv.ParseInt(processedWords[lastIdx], 16, 64)
				if err == nil {
					processedWords[lastIdx] = strconv.Itoa(int(val))
				}
			}
			continue
		}

		if word == "(bin)" {
			if len(processedWords) > 0 {
				lastIdx := len(processedWords) - 1
				val, err := strconv.ParseInt(processedWords[lastIdx], 2, 64)
				if err == nil {
					processedWords[lastIdx] = strconv.Itoa(int(val))
				}
			}
			continue
		}

		if word == "(up)" {
			applyWordTransform(processedWords, strings.ToUpper, 1)
			continue
		}

		if word == "(low)" {
			applyWordTransform(processedWords, strings.ToLower, 1)
			continue
		}

		if word == "(cap)" {
			applyWordTransform(processedWords, capitalizeWord, 1)
			continue
		}

		if strings.HasPrefix(word, "(") {
			cleanMarker := word[1:]

			if (cleanMarker == "up" || cleanMarker == "low" || cleanMarker == "cap") && i+2 < len(words) {
				if words[i+1] == "," {
					countToken := words[i+2]
					if strings.HasSuffix(countToken, ")") {
						numStr := countToken[:len(countToken)-1]
						count, err := strconv.Atoi(numStr)
						if err == nil {
							i += 2 // Skip comma and number
							switch cleanMarker {
							case "up":
								applyWordTransform(processedWords, strings.ToUpper, count)
							case "low":
								applyWordTransform(processedWords, strings.ToLower, count)
							case "cap":
								applyWordTransform(processedWords, capitalizeWord, count)
							}
							continue
						}
					}
				}
			}
		}

		processedWords = append(processedWords, word)
	}

	// Pass 2: Vowels
	processedWords = handleVowels(processedWords)

	// Pass 3: Punctuation
	finalString := handlePunctuation(processedWords)

	return finalString
}

func applyWordTransform(words []string, transform func(string) string, count int) {
	startIndex := len(words) - count
	if startIndex < 0 {
		startIndex = 0
	}
	for j := startIndex; j < len(words); j++ {
		words[j] = transform(words[j])
	}
}

func capitalizeWord(s string) string {
	if len(s) == 0 {
		return s
	}
	runes := []rune(s)
	return string(unicode.ToUpper(runes[0])) + strings.ToLower(string(runes[1:]))
}

func handleVowels(words []string) []string {
	vowels := "aeiouhAEIOUH"
	for i := 0; i < len(words)-1; i++ {
		if words[i] == "a" || words[i] == "A" {
			if len(words[i+1]) == 0 {
				continue
			}
			firstChar := rune(words[i+1][0])
			if strings.ContainsRune(vowels, firstChar) {
				if words[i] == "a" {
					words[i] = "an"
				} else {
					words[i] = "An"
				}
			}
		}
	}
	return words
}

func handlePunctuation(words []string) string {
	var sb strings.Builder
	puncts := ".,!?:;"
	inQuote := false

	for i := 0; i < len(words); i++ {
		word := words[i]

		addSpaceBefore := true
		if i == 0 {
			addSpaceBefore = false
		}

		isPunct := false
		if len(word) > 0 && (strings.ContainsAny(string(word[0]), puncts) || isPunctuationGroup(word, puncts)) {
			isPunct = true
		}

		if word == "'" {
			if !inQuote {
				inQuote = true
			} else {
				inQuote = false
				addSpaceBefore = false
			}

			if addSpaceBefore {
				sb.WriteString(" ")
			}
			sb.WriteString(word)
			continue
		}

		if isPunct {
			addSpaceBefore = false
		}

		if i > 0 {
			prevWord := words[i-1]
			if prevWord == "'" && inQuote {
				// Previous was Open Quote. Current word should touch it.
				addSpaceBefore = false
			}
		}

		if addSpaceBefore {
			sb.WriteString(" ")
		}
		sb.WriteString(word)
	}
	return sb.String()
}

func isPunctuationGroup(s string, puncts string) bool {
	for _, r := range s {
		if !strings.ContainsRune(puncts, r) && r != '.' {
			// n '...' contains '.' which is in puncts.
			return false
		}
	}
	return true
}
