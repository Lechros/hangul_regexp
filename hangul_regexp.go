package hangul_regexp

import (
	"errors"
	"slices"
	"strings"
	"unicode/utf8"
)

func GetPattern(search string, ignoreSpace bool, fuzzy bool, matchChoseong bool) (string, error) {
	if ignoreSpace && fuzzy {
		return "", errors.New("ignoreSpace and fuzzy cannot be true at the same time")
	}

	connector := ""
	if ignoreSpace {
		connector = " *?"
	}
	if fuzzy {
		connector = ".*?"
	}

	builder := strings.Builder{}

	for i, ch := range search {
		if i+utf8.RuneLen(ch) == len(search) {
			if isHangul(ch) {
				builder.WriteString(getLastHangulPattern(ch, connector))
			} else if canBeChoseong(ch) {
				builder.WriteString(getChoseongPattern(ch))
			} else {
				if shouldEscapeRegexCharacter(ch) {
					builder.WriteRune('\\')
				}
				builder.WriteRune(ch)
			}
		} else {
			if canBeChoseong(ch) && matchChoseong {
				builder.WriteString(getChoseongPattern(ch))
			} else {
				if shouldEscapeRegexCharacter(ch) {
					builder.WriteRune('\\')
				}
				builder.WriteRune(ch)
			}
			builder.WriteString(connector)
		}
	}

	return builder.String(), nil
}

func shouldEscapeRegexCharacter(ch rune) bool {
	chars := [...]rune{'.', '^', '$', '*', '+', '?', '(', ')', '[', '{', '\\', '['}
	return slices.Contains(chars[:], ch)
}

func isHangul(ch rune) bool {
	return '가' <= ch && ch <= '힣'
}

func canBeChoseong(ch rune) bool {
	if 'ㄱ' <= ch && ch <= 'ㅎ' {
		_, found := slices.BinarySearch(choseongs[:], ch)
		return found
	}
	return false
}

func hasBatchim(hangul rune) bool {
	return (hangul-'가')%28 > 0
}

func getChoseongPattern(choseong rune) string {
	choOffset := getChoseongOffset(choseong)
	builder := strings.Builder{}
	builder.WriteString("(?:")
	builder.WriteRune(choseong)
	builder.WriteString("|[")
	builder.WriteRune(assemble(choOffset, 0, 0))
	builder.WriteRune('-')
	builder.WriteRune(assemble(choOffset, len(jungseongs)-1, len(jongseongs)-1))
	builder.WriteString("])")
	return builder.String()
}

func getLastHangulPattern(hangul rune, connector string) string {
	choOffset, jungOffset, jongOffset := disassemble(hangul)
	builder := strings.Builder{}
	if hasBatchim(hangul) {
		jongseong := jongseongs[jongOffset]
		if canBeChoseong(jongseong) {
			builder.WriteString("(?:")
			builder.WriteRune(hangul)
			builder.WriteRune('|')
			builder.WriteRune(assemble(choOffset, jungOffset, 0))
			builder.WriteString(connector)
			builder.WriteString(getChoseongPattern(jongseong))
			builder.WriteRune(')')
		} else {
			firstJong, secondJong := splitJongseong(jongseong)
			builder.WriteString("(?:")
			builder.WriteRune(hangul)
			builder.WriteRune('|')
			builder.WriteRune(assemble(choOffset, jungOffset, getJongseongOffset(firstJong)))
			builder.WriteString(connector)
			builder.WriteString(getChoseongPattern(secondJong))
			builder.WriteRune(')')
		}
	} else {
		builder.WriteString("(?:")
		builder.WriteRune(hangul)
		builder.WriteString("|[")
		builder.WriteRune(assemble(choOffset, jungOffset, 1))
		builder.WriteRune('-')
		builder.WriteRune(assemble(choOffset, jungOffset, len(jongseongs)-1))
		builder.WriteString("])")
	}
	return builder.String()
}

func getChoseongOffset(choseong rune) int {
	index, _ := slices.BinarySearch(choseongs[:], choseong)
	return index
}

func getJongseongOffset(jongseong rune) int {
	index, _ := slices.BinarySearch(jongseongs[:], jongseong)
	return index
}

func disassemble(hangul rune) (int, int, int) {
	hangul -= '가'
	return int(hangul / 28 / 21), int(hangul / 28 % 21), int(hangul % 28)
}

func assemble(choOffset, jungOffset, jongOffset int) rune {
	return rune('가' + (choOffset*21+jungOffset)*28 + jongOffset)
}

func splitJongseong(jongseong rune) (rune, rune) {
	switch jongseong {
	case 'ㄳ':
		return 'ㄱ', 'ㅅ'
	case 'ㄵ':
		return 'ㄴ', 'ㅈ'
	case 'ㄶ':
		return 'ㄴ', 'ㅎ'
	case 'ㄺ':
		return 'ㄹ', 'ㄱ'
	case 'ㄻ':
		return 'ㄹ', 'ㅁ'
	case 'ㄼ':
		return 'ㄹ', 'ㅂ'
	case 'ㄽ':
		return 'ㄹ', 'ㅅ'
	case 'ㄾ':
		return 'ㄹ', 'ㅌ'
	case 'ㄿ':
		return 'ㄹ', 'ㅍ'
	case 'ㅀ':
		return 'ㄹ', 'ㅎ'
	case 'ㅄ':
		return 'ㅂ', 'ㅅ'
	}
	panic(jongseong)
}

var choseongs = [...]rune{'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
var jungseongs = [...]rune{'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ', 'ㅙ', 'ㅚ', 'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ', 'ㅣ'}
var jongseongs = [...]rune{-1, 'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ', 'ㄹ', 'ㄺ', 'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ', 'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
