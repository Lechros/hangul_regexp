package hangul_regexp

import (
	"errors"
	"strings"
	"unicode/utf8"
)

func GetPattern(search string, ignoreSpace bool, fuzzy bool, matchChoseong bool, capturing bool) (string, error) {
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
	builder.Grow(preCalculateBytes(search, len(connector), matchChoseong, capturing))

	for i, ch := range search {
		if i+utf8.RuneLen(ch) == len(search) {
			if IsHangul(ch) {
				writeLastHangulPattern(&builder, ch, connector, capturing)
			} else if CanBeChoseong(ch) {
				writeChoseongPattern(&builder, ch, capturing)
			} else if matchChoseong && CanBeChoseongOrJongseong(ch) {
				writeCombinedChoseongPattern(&builder, ch, connector, capturing)
			} else {
				writeEscaped(&builder, ch, capturing)
			}
		} else {
			if matchChoseong && CanBeChoseongOrJongseong(ch) {
				if CanBeChoseong(ch) {
					writeChoseongPattern(&builder, ch, capturing)
				} else {
					writeCombinedChoseongPattern(&builder, ch, connector, capturing)
				}
			} else {
				writeEscaped(&builder, ch, capturing)
			}
			builder.WriteString(connector)
		}
	}

	return builder.String(), nil
}

func preCalculateBytes(str string, connectorLength int, matchChoseong bool, capturing bool) int {
	if matchChoseong {
		size := len(str)
		for _, ch := range str {
			if 'ㄱ' <= ch && ch <= 'ㅎ' {
				// Choseong pattern is 17 bytes and choseong character is 3 bytes
				size += 17 - 3
			}
			// Accommodate for connector and possible regex character escape
			size += connectorLength + 1
			if capturing {
				size += 2
			}
		}
		// Last hangul pattern is max. 28 bytes
		size += 28 - 3
		if capturing {
			size += 4
		}
		return size
	} else if capturing {
		return len(str) + utf8.RuneCountInString(str)*(connectorLength+1+2) + (28 - 3 + 4)
	} else {
		return len(str) + utf8.RuneCountInString(str)*(connectorLength+1) + (28 - 3)
	}
}

func writeEscaped(builder *strings.Builder, ch rune, capturing bool) {
	if capturing {
		builder.WriteRune('(')
	}
	switch ch {
	case '.', '^', '$', '*', '+', '?', '(', ')', '[', '{', '\\', '|':
		builder.WriteRune('\\')
	}
	builder.WriteRune(ch)
	if capturing {
		builder.WriteRune(')')
	}
}

func writeChoseongPattern(builder *strings.Builder, choseong rune, capturing bool) {
	choOffset := GetChoseongOffset(choseong)
	if capturing {
		builder.WriteRune('(')
		builder.WriteRune(choseong)
		builder.WriteString("|[")
		builder.WriteRune(Assemble(choOffset, 0, 0))
		builder.WriteRune('-')
		builder.WriteRune(Assemble(choOffset, len(jungseongs)-1, len(jongseongs)-1))
		builder.WriteString("])")
	} else {
		builder.WriteString("(?:")
		builder.WriteRune(choseong)
		builder.WriteString("|[")
		builder.WriteRune(Assemble(choOffset, 0, 0))
		builder.WriteRune('-')
		builder.WriteRune(Assemble(choOffset, len(jungseongs)-1, len(jongseongs)-1))
		builder.WriteString("])")
	}
}

func writeCombinedChoseongPattern(builder *strings.Builder, jongseong rune, connector string, capturing bool) {
	firstCho, secondCho := SplitJongseong(jongseong)
	writeChoseongPattern(builder, firstCho, capturing)
	builder.WriteString(connector)
	writeChoseongPattern(builder, secondCho, capturing)
}

func writeLastHangulPattern(builder *strings.Builder, hangul rune, connector string, capturing bool) {
	choOffset, jungOffset, jongOffset := Disassemble(hangul)
	if HasBatchim(hangul) {
		jongseong := jongseongs[jongOffset]
		if CanBeChoseong(jongseong) {
			builder.WriteString("(?:")
			writeRune(builder, hangul, capturing)
			builder.WriteRune('|')
			writeRune(builder, Assemble(choOffset, jungOffset, 0), capturing)
			builder.WriteString(connector)
			writeChoseongPattern(builder, jongseong, capturing)
			builder.WriteRune(')')
		} else {
			firstJong, secondJong := SplitJongseong(jongseong)
			builder.WriteString("(?:")
			writeRune(builder, hangul, capturing)
			builder.WriteRune('|')
			writeRune(builder, Assemble(choOffset, jungOffset, GetJongseongOffset(firstJong)), capturing)
			builder.WriteString(connector)
			writeChoseongPattern(builder, secondJong, capturing)
			builder.WriteRune(')')
		}
	} else {
		if capturing {
			builder.WriteRune('(')
			builder.WriteRune(hangul)
			builder.WriteString("|[")
			builder.WriteRune(Assemble(choOffset, jungOffset, 1))
			builder.WriteRune('-')
			builder.WriteRune(Assemble(choOffset, jungOffset, len(jongseongs)-1))
			builder.WriteString("])")
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
}

func writeRune(builder *strings.Builder, ch rune, capturing bool) {
	if capturing {
		builder.WriteRune('(')
		builder.WriteRune(ch)
		builder.WriteRune(')')
	} else {
		builder.WriteRune(ch)
	}
}
