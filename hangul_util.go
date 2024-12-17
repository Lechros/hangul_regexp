package hangul_regexp

func IsHangul(ch rune) bool {
	return '가' <= ch && ch <= '힣'
}

func CanBeChoseongOrJongseong(ch rune) bool {
	return 'ㄱ' <= ch && ch <= 'ㅎ'
}

func CanBeChoseong(ch rune) bool {
	return GetChoseongOffset(ch) >= 0
}

func HasBatchim(hangul rune) bool {
	return (hangul-'가')%28 > 0
}

func Disassemble(hangul rune) (int, int, int) {
	hangul -= '가'
	return int(hangul / 28 / 21), int(hangul / 28 % 21), int(hangul % 28)
}

func Assemble(choOffset, jungOffset, jongOffset int) rune {
	return rune('가' + (choOffset*21+jungOffset)*28 + jongOffset)
}

func GetChoseongOffset(choseong rune) int {
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

func GetJongseongOffset(jongseong rune) int {
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

func SplitJongseong(jongseong rune) (rune, rune) {
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
