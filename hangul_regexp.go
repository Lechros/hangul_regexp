package hangul_regexp

import (
	"errors"
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
				writeLastHangulPattern(&builder, ch, connector)
			} else if canBeChoseong(ch) {
				writeChoseongPattern(&builder, ch)
			} else {
				writeEscaped(&builder, ch)
			}
		} else {
			if matchChoseong && canBeChoseong(ch) {
				writeChoseongPattern(&builder, ch)
			} else {
				writeEscaped(&builder, ch)
			}
			builder.WriteString(connector)
		}
	}

	return builder.String(), nil
}

func writeEscaped(builder *strings.Builder, ch rune) {
	switch ch {
	case '.', '^', '$', '*', '+', '?', '(', ')', '[', '{', '\\', '|':
		builder.WriteRune('\\')
	}
	builder.WriteRune(ch)
}

func writeChoseongPattern(builder *strings.Builder, choseong rune) {
	choOffset := getChoseongOffset(choseong)
	builder.WriteString("(?:")
	builder.WriteRune(choseong)
	builder.WriteString("|[")
	builder.WriteRune(assemble(choOffset, 0, 0))
	builder.WriteRune('-')
	builder.WriteRune(assemble(choOffset, len(jungseongs)-1, len(jongseongs)-1))
	builder.WriteString("])")
}

func writeLastHangulPattern(builder *strings.Builder, hangul rune, connector string) {
	choOffset, jungOffset, jongOffset := disassemble(hangul)
	if hasBatchim(hangul) {
		jongseong := jongseongs[jongOffset]
		if canBeChoseong(jongseong) {
			builder.WriteString("(?:")
			builder.WriteRune(hangul)
			builder.WriteRune('|')
			builder.WriteRune(assemble(choOffset, jungOffset, 0))
			builder.WriteString(connector)
			writeChoseongPattern(builder, jongseong)
			builder.WriteRune(')')
		} else {
			firstJong, secondJong := splitJongseong(jongseong)
			builder.WriteString("(?:")
			builder.WriteRune(hangul)
			builder.WriteRune('|')
			builder.WriteRune(assemble(choOffset, jungOffset, getJongseongOffset(firstJong)))
			builder.WriteString(connector)
			writeChoseongPattern(builder, secondJong)
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
}

func isHangul(ch rune) bool {
	return '가' <= ch && ch <= '힣'
}

func canBeChoseong(ch rune) bool {
	if 'ㄱ' <= ch && ch <= 'ㅎ' {
		return getChoseongOffset(ch) >= 0
	}
	return false
}

func hasBatchim(hangul rune) bool {
	return (hangul-'가')%28 > 0
}

func disassemble(hangul rune) (int, int, int) {
	hangul -= '가'
	return int(hangul / 28 / 21), int(hangul / 28 % 21), int(hangul % 28)
}

func assemble(choOffset, jungOffset, jongOffset int) rune {
	return rune('가' + (choOffset*21+jungOffset)*28 + jongOffset)
}

func getChoseongOffset(choseong rune) int {
	switch choseong {
	case 'ㄱ':
		return 0
	case 'ㄲ':
		return 1
	case 'ㄴ':
		return 2
	case 'ㄷ':
		return 3
	case 'ㄸ':
		return 4
	case 'ㄹ':
		return 5
	case 'ㅁ':
		return 6
	case 'ㅂ':
		return 7
	case 'ㅃ':
		return 8
	case 'ㅅ':
		return 9
	case 'ㅆ':
		return 10
	case 'ㅇ':
		return 11
	case 'ㅈ':
		return 12
	case 'ㅉ':
		return 13
	case 'ㅊ':
		return 14
	case 'ㅋ':
		return 15
	case 'ㅌ':
		return 16
	case 'ㅍ':
		return 17
	case 'ㅎ':
		return 18
	default:
		return -1
	}
}

func getJongseongOffset(jongseong rune) int {
	switch jongseong {
	case -1:
		return 0
	case 'ㄱ':
		return 1
	case 'ㄲ':
		return 2
	case 'ㄳ':
		return 3
	case 'ㄴ':
		return 4
	case 'ㄵ':
		return 5
	case 'ㄶ':
		return 6
	case 'ㄷ':
		return 7
	case 'ㄹ':
		return 8
	case 'ㄺ':
		return 9
	case 'ㄻ':
		return 10
	case 'ㄼ':
		return 11
	case 'ㄽ':
		return 12
	case 'ㄾ':
		return 13
	case 'ㄿ':
		return 14
	case 'ㅀ':
		return 15
	case 'ㅁ':
		return 16
	case 'ㅂ':
		return 17
	case 'ㅄ':
		return 18
	case 'ㅅ':
		return 19
	case 'ㅆ':
		return 20
	case 'ㅇ':
		return 21
	case 'ㅈ':
		return 22
	case 'ㅊ':
		return 23
	case 'ㅋ':
		return 24
	case 'ㅌ':
		return 25
	case 'ㅍ':
		return 26
	case 'ㅎ':
		return 27
	default:
		return -1
	}
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
