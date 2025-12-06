package main

import (
	"strconv"
	"strings"
	"unicode"
)

// ProcessText takes the input text and applies all the required modifications.
func ProcessText(input string) string {
	// Pre-process: pad punctuation with spaces to ensure they are tokenized separately.
	// punctuations to handle: . , ! ? : ; '
	// Note: brackets () for markers are NOT included, so (hex) stays together.
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
			applyTransformation(processedWords, strings.ToUpper, 1)
			continue
		}

		if word == "(low)" {
			applyTransformation(processedWords, strings.ToLower, 1)
			continue
		}

		if word == "(cap)" {
			applyTransformation(processedWords, capitalize, 1)
			continue
		}

		// Handle (up, n), (low, n), (cap, n) with potential padded comma
		// Pattern: "(up", ",", "2)" due to pre-processing comma
		if strings.HasPrefix(word, "(") {
			cleanMarker := word[1:]
			// It might be just "(up" if comma was split, or "(up," if not?
			// Since we pad ",", it becomes " " + "," + " ". So "(up, 2)" -> "(up" " ," " 2)".

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
								applyTransformation(processedWords, strings.ToUpper, count)
							case "low":
								applyTransformation(processedWords, strings.ToLower, count)
							case "cap":
								applyTransformation(processedWords, capitalize, count)
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

func applyTransformation(words []string, transform func(string) string, count int) {
	startIndex := len(words) - count
	if startIndex < 0 {
		startIndex = 0
	}
	for j := startIndex; j < len(words); j++ {
		words[j] = transform(words[j])
	}
}

func capitalize(s string) string {
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
				// Open quote
				// Keeps space before (default true), no space after (handled by next word's logic not looking back usually, EXCEPT special check below)
				// Wait, "Space apart from next one" is NOT true for Open Quote. "Close to the next one".
				// My logic for "Next word" needs to handle this.
			} else {
				// Close quote: touches previous word.
				inQuote = false
				addSpaceBefore = false
			}

			if addSpaceBefore {
				sb.WriteString(" ")
			}
			sb.WriteString(word)

			// if inQuote {
			// 	inQuote = false
			// } else {
			// 	inQuote = true
			// }
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
			// note: '...' contains '.' which is in puncts.
			// But prompt groups example: "!?".
			return false
		}
	}
	return true
}
