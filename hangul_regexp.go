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
	builder.Grow(preCalculateBytes(search, len(connector), matchChoseong))

	for i, ch := range search {
		if i+utf8.RuneLen(ch) == len(search) {
			if IsHangul(ch) {
				writeLastHangulPattern(&builder, ch, connector)
			} else if CanBeChoseong(ch) {
				writeChoseongPattern(&builder, ch)
			} else if matchChoseong && CanBeChoseongOrJongseong(ch) {
				writeCombinedChoseongPattern(&builder, ch)
			} else {
				writeEscaped(&builder, ch)
			}
		} else {
			if matchChoseong && CanBeChoseongOrJongseong(ch) {
				if CanBeChoseong(ch) {
					writeChoseongPattern(&builder, ch)
				} else {
					writeCombinedChoseongPattern(&builder, ch)
				}
			} else {
				writeEscaped(&builder, ch)
			}
			builder.WriteString(connector)
		}
	}

	return builder.String(), nil
}

func preCalculateBytes(str string, connectorLength int, matchChoseong bool) int {
	if matchChoseong {
		size := len(str)
		for _, ch := range str {
			if 'ㄱ' <= ch && ch <= 'ㅎ' {
				// Choseong pattern is 17 bytes and choseong character is 3 bytes
				size += 17 - 3
			}
			// Accommodate for connector and possible regex character escape
			size += connectorLength + 1
		}
		// Last hangul pattern is max. 28 bytes
		size += 28 - 3
		return size
	} else {
		return len(str) + utf8.RuneCountInString(str)*(connectorLength+1) + (28 - 3)
	}
}

func writeEscaped(builder *strings.Builder, ch rune) {
	switch ch {
	case '.', '^', '$', '*', '+', '?', '(', ')', '[', '{', '\\', '|':
		builder.WriteRune('\\')
	}
	builder.WriteRune(ch)
}

func writeChoseongPattern(builder *strings.Builder, choseong rune) {
	choOffset := GetChoseongOffset(choseong)
	builder.WriteString("(?:")
	builder.WriteRune(choseong)
	builder.WriteString("|[")
	builder.WriteRune(Assemble(choOffset, 0, 0))
	builder.WriteRune('-')
	builder.WriteRune(Assemble(choOffset, len(jungseongs)-1, len(jongseongs)-1))
	builder.WriteString("])")
}

func writeCombinedChoseongPattern(builder *strings.Builder, jongseong rune) {
	firstCho, secondCho := SplitJongseong(jongseong)
	writeChoseongPattern(builder, firstCho)
	writeChoseongPattern(builder, secondCho)
}

func writeLastHangulPattern(builder *strings.Builder, hangul rune, connector string) {
	choOffset, jungOffset, jongOffset := Disassemble(hangul)
	if HasBatchim(hangul) {
		jongseong := jongseongs[jongOffset]
		if CanBeChoseong(jongseong) {
			builder.WriteString("(?:")
			builder.WriteRune(hangul)
			builder.WriteRune('|')
			builder.WriteRune(Assemble(choOffset, jungOffset, 0))
			builder.WriteString(connector)
			writeChoseongPattern(builder, jongseong)
			builder.WriteRune(')')
		} else {
			firstJong, secondJong := SplitJongseong(jongseong)
			builder.WriteString("(?:")
			builder.WriteRune(hangul)
			builder.WriteRune('|')
			builder.WriteRune(Assemble(choOffset, jungOffset, GetJongseongOffset(firstJong)))
			builder.WriteString(connector)
			writeChoseongPattern(builder, secondJong)
			builder.WriteRune(')')
		}
	} else {
		builder.WriteString("(?:")
		builder.WriteRune(hangul)
		builder.WriteString("|[")
		builder.WriteRune(Assemble(choOffset, jungOffset, 1))
		builder.WriteRune('-')
		builder.WriteRune(Assemble(choOffset, jungOffset, len(jongseongs)-1))
		builder.WriteString("])")
	}
}
