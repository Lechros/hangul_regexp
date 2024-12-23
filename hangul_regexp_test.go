package hangul_regexp

import (
	"fmt"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func TestGetPattern(t *testing.T) {
	type args struct {
		search      string
		ignoreSpace bool
		fuzzy       bool
		choseong    bool
		capturing   bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Alphabet", args{"a", false, false, false, false}, "a", false},
		{"Hangul", args{"ㄱ나다라123", false, false, false, false}, "ㄱ나다라123", false},
		{"Mixed", args{"Zx0ㅡㅡ", false, false, false, false}, "Zx0ㅡㅡ", false},
		{"Should escape", args{"[^가-힣]$", false, false, false, false}, "\\[\\^가-힣]\\$", false},

		{"Last char is choseong", args{"ㄱ", false, false, false, false}, "(?:ㄱ|[가-깋])", false},
		{"Last char without batchim", args{"가 나", false, false, false, false}, "가 (?:나|[낙-낳])", false},
		{"Last char with batchim", args{"가 안", false, false, false, false}, "가 (?:안|아(?:ㄴ|[나-닣]))", false},
		{"Last char with double batchim", args{"가 있", false, false, false, false}, "가 (?:있|이(?:ㅆ|[싸-앃]))", false},
		{"Last char with combined batchim", args{"가 얇", false, false, false, false}, "가 (?:얇|얄(?:ㅂ|[바-빟]))", false},
		{"Last char is combined choseong", args{"ㄻ", false, false, false, false}, "ㄻ", false},

		{"Mixed / ignoreSpace=true", args{"ㅁ가a항1", true, false, false, false}, "ㅁ *?가 *?a *?항 *?1", false},
		{"Last char with batchim / ignoreSpace=true / has space matcher between", args{"가 안", true, false, false, false}, "가 *?  *?(?:안|아 *?(?:ㄴ|[나-닣]))", false},

		{"Mixed / fuzzy=true", args{"ㅁ가a항1", false, true, false, false}, "ㅁ.*?가.*?a.*?항.*?1", false},
		{"Space / fuzzy=true / spaces are not concatenated", args{"가 s", false, true, false, false}, "가.*? .*?s", false},
		{"Last char with batchim / fuzzy=true / has any matcher between", args{"가 안", false, true, false, false}, "가.*? .*?(?:안|아.*?(?:ㄴ|[나-닣]))", false},

		{"Non-last char is choseong / choseong=true", args{"ㄱ1", false, false, true, false}, "(?:ㄱ|[가-깋])1", false},
		{"Multiple choseong chars / choseong=true", args{"ㄱ ㄴㄷ", false, false, true, false}, "(?:ㄱ|[가-깋]) (?:ㄴ|[나-닣])(?:ㄷ|[다-딯])", false},
		{"Mixed with choseong / choseong=true", args{"aㅎ1가ㄴ", false, false, true, false}, "a(?:ㅎ|[하-힣])1가(?:ㄴ|[나-닣])", false},
		{"Standalone batchim char / choseong=true", args{"ㄻㅄ", false, false, true, false}, "(?:ㄹ|[라-맇])(?:ㅁ|[마-밓])(?:ㅂ|[바-빟])(?:ㅅ|[사-싷])", false},

		{"Any / ignoreSpace=true, fuzzy=true / err", args{"", true, true, false, false}, "", true},

		{"Alphabet / capturing=true", args{"a", false, false, false, true}, "(a)", false},
		{"Hangul / capturing=true", args{"ㄱ나다라123", false, false, false, true}, "(ㄱ)(나)(다)(라)(1)(2)(3)", false},
		{"Mixed / capturing=true", args{"Zx0ㅡㅡ", false, false, false, true}, "(Z)(x)(0)(ㅡ)(ㅡ)", false},
		{"Special chars / capturing=true", args{"[^가-힣]$", false, false, false, true}, "(\\[)(\\^)(가)(-)(힣)(])(\\$)", false},

		{"Last char with batchim / capturing=true", args{"가 안", false, false, false, true}, "(가)( )(?:(안)|(아)(ㄴ|[나-닣]))", false},
		{"Last char with double batchim / capturing=true", args{"가 있", false, false, false, true}, "(가)( )(?:(있)|(이)(ㅆ|[싸-앃]))", false},
		{"Last char with combined batchim / capturing=true", args{"가 얇", false, false, false, true}, "(가)( )(?:(얇)|(얄)(ㅂ|[바-빟]))", false},
		{"Last char is combined choseong / capturing=true", args{"ㄻ", false, false, false, true}, "(ㄻ)", false},

		{"Mixed / fuzzy=true, capturing=true", args{"ㅁ가a항1", false, true, false, true}, "(ㅁ).*?(가).*?(a).*?(항).*?(1)", false},
		{"Space / fuzzy=true, capturing=true", args{"가 s", false, true, false, true}, "(가).*?( ).*?(s)", false},
		{"Last char with batchim / fuzzy=true, capturing=true", args{"가 안", false, true, false, true}, "(가).*?( ).*?(?:(안)|(아).*?(ㄴ|[나-닣]))", false},

		{"Standalone batchim char / choseong=true, capturing=true", args{"ㄻㅄ", false, false, true, true}, "(ㄹ|[라-맇])(ㅁ|[마-밓])(ㅂ|[바-빟])(ㅅ|[사-싷])", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPattern(tt.args.search, tt.args.ignoreSpace, tt.args.fuzzy, tt.args.choseong, tt.args.capturing)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetPattern() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetPattern() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestGetPatternCapturing(t *testing.T) {
	tests := []struct {
		search      string
		target      string
		matchConcat string
	}{
		{"이이저", "이글아이 레인저슈트", "이이저"},
		{"마깃아", "마력이 깃든 안대", "마깃안"},
		{"루컨ㅁ", "루즈 컨트롤 머신 마크", "루컨머"},
		{"ㅇㅋㅇ", "아케인셰이드 스태프", "아케인"},
		{"낢", "날아 먹", "날먹"},
		{"ㅄ", "보라색", "보색"},
	}
	for _, tt := range tests {
		t.Run(tt.search+" / "+tt.target, func(t *testing.T) {
			pattern, _ := GetPattern(tt.search, false, true, true, true)
			println(pattern)
			regex := regexp.MustCompile(pattern)
			match := regex.FindStringSubmatch(tt.target)
			for _, m := range match {
				println(m)
			}
			actual := strings.Join(match[1:], "")
			if actual != tt.matchConcat {
				t.Errorf("GetPattern() got = %v, want %v", actual, tt.matchConcat)
			}
		})
	}
}

func BenchmarkGetPattern_마깃안(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetPattern("마깃안", false, false, false, false)
	}
}

func BenchmarkGetPattern_마깃안_fuzzy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetPattern("마깃안", false, true, false, false)
	}
}

func BenchmarkGetPattern_아케인셰이드_에너지소드(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetPattern("아케인셰이드 에너지소드", false, false, false, false)
	}
}

func BenchmarkGetPattern_아케인셰이드_에너지소드_fuzzy(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetPattern("아케인셰이드 에너지소드", false, true, false, false)
	}
}

func BenchmarkGetPattern_아케인셰이드_에너지소드_fuzzy_matchChoseong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetPattern("아케인셰이드 에너지소드", false, true, true, false)
	}
}

func BenchmarkGetPattern_ㅇㅋㅇㅅㅇㄷ_ㅇㄴㅈㅅㄷ_fuzzy_matchChoseong(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_, _ = GetPattern("ㅇㅋㅇㅅㅇㄷ ㅇㄴㅈㅅㄷ", false, true, true, false)
	}
}

func BenchmarkLastHangulString_Sprint(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprint("(?:",
			'간',
			'|',
			'가',
			".*?",
			"(?:ㄴ|[나-닣])",
			')')
	}
}

func BenchmarkLastHangulString_Sprintf(b *testing.B) {
	for i := 0; i < b.N; i++ {
		_ = fmt.Sprintf("(?:%c|%c%s%s)",
			'간',
			'가',
			".*?",
			"(?:ㄴ|[나-닣])")
	}
}

func BenchmarkLastHangulString_Builder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := strings.Builder{}
		builder.WriteString("(?:")
		builder.WriteRune('간')
		builder.WriteRune('|')
		builder.WriteRune('가')
		builder.WriteString(".*?")
		builder.WriteString("(?:ㄴ|[나-닣])")
		builder.WriteRune(')')
		_ = builder.String()
	}
}

func BenchmarkLastHangulString_BuilderAllWriteRune(b *testing.B) {
	for i := 0; i < b.N; i++ {
		builder := strings.Builder{}
		builder.WriteRune('(')
		builder.WriteRune('?')
		builder.WriteRune(':')
		builder.WriteRune('간')
		builder.WriteRune('|')
		builder.WriteRune('가')
		builder.WriteString(".*?")
		builder.WriteString("(?:ㄴ|[나-닣])")
		builder.WriteRune(')')
		_ = builder.String()
	}
}

func BenchmarkBetweenFunction_ReturnString(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBetweenFunction_ReturnString_Outer()
	}
}

func benchmarkBetweenFunction_ReturnString_Outer() {
	builder := strings.Builder{}
	builder.WriteString(benchmarkBetweenFunction_ReturnString_Inner())
	_ = builder.String()
}

func benchmarkBetweenFunction_ReturnString_Inner() string {
	builder := strings.Builder{}
	builder.WriteString("(?:")
	builder.WriteRune('ㄱ')
	builder.WriteString("|[")
	builder.WriteRune('가')
	builder.WriteRune('-')
	builder.WriteRune('깋')
	builder.WriteRune(')')
	return builder.String()
}

func BenchmarkBetweenFunction_PassBuilder(b *testing.B) {
	for i := 0; i < b.N; i++ {
		benchmarkBetweenFunction_PassBuilder_Outer()
	}
}

func benchmarkBetweenFunction_PassBuilder_Outer() {
	builder := strings.Builder{}
	benchmarkBetweenFunction_PassBuilder_Inner(builder)
	_ = builder.String()
}

func benchmarkBetweenFunction_PassBuilder_Inner(builder strings.Builder) {
	builder.WriteString("(?:")
	builder.WriteRune('ㄱ')
	builder.WriteString("|[")
	builder.WriteRune('가')
	builder.WriteRune('-')
	builder.WriteRune('깋')
	builder.WriteRune(')')
}

func BenchmarkGetChoseongOffset_BinarySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range choseongs {
			slices.BinarySearch(choseongs[:], choseongs[ci])
		}
	}
}

func BenchmarkGetChoseongOffset_Linear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range choseongs {
			slices.Index(choseongs[:], choseongs[ci])
		}
	}
}

func BenchmarkGetChoseongOffset_Map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range choseongs {
			_ = choseongMap[choseongs[ci]]
		}
	}
}

func BenchmarkGetChoseongOffset_Switch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range choseongs {
			switch choseongs[ci] {
			case 'ㄱ':
				_ = 0
			case 'ㄲ':
				_ = 1
			case 'ㄴ':
				_ = 2
			case 'ㄷ':
				_ = 3
			case 'ㄸ':
				_ = 4
			case 'ㄹ':
				_ = 5
			case 'ㅁ':
				_ = 6
			case 'ㅂ':
				_ = 7
			case 'ㅃ':
				_ = 8
			case 'ㅅ':
				_ = 9
			case 'ㅆ':
				_ = 10
			case 'ㅇ':
				_ = 11
			case 'ㅈ':
				_ = 12
			case 'ㅉ':
				_ = 13
			case 'ㅊ':
				_ = 14
			case 'ㅋ':
				_ = 15
			case 'ㅌ':
				_ = 16
			case 'ㅍ':
				_ = 17
			case 'ㅎ':
				_ = 18
			}
		}
	}
}

func BenchmarkGetJongseongOffset_BinarySearch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range jongseongs {
			slices.BinarySearch(jongseongs[:], jongseongs[ci])
		}
	}
}

func BenchmarkGetJongseongOffset_Linear(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range jongseongs {
			slices.Index(jongseongs[:], jongseongs[ci])
		}
	}
}

func BenchmarkGetJongseongOffset_Map(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range jongseongs {
			_ = jongseongMap[jongseongs[ci]]
		}
	}
}

func BenchmarkGetJongseongOffset_Switch(b *testing.B) {
	for i := 0; i < b.N; i++ {
		for ci := range jongseongs {
			switch jongseongs[ci] {
			case -1:
				_ = 0
			case 'ㄱ':
				_ = 1
			case 'ㄲ':
				_ = 2
			case 'ㄳ':
				_ = 3
			case 'ㄴ':
				_ = 4
			case 'ㄵ':
				_ = 5
			case 'ㄶ':
				_ = 6
			case 'ㄷ':
				_ = 7
			case 'ㄹ':
				_ = 8
			case 'ㄺ':
				_ = 9
			case 'ㄻ':
				_ = 10
			case 'ㄼ':
				_ = 11
			case 'ㄽ':
				_ = 12
			case 'ㄾ':
				_ = 13
			case 'ㄿ':
				_ = 14
			case 'ㅀ':
				_ = 15
			case 'ㅁ':
				_ = 16
			case 'ㅂ':
				_ = 17
			case 'ㅄ':
				_ = 18
			case 'ㅅ':
				_ = 19
			case 'ㅆ':
				_ = 20
			case 'ㅇ':
				_ = 21
			case 'ㅈ':
				_ = 22
			case 'ㅊ':
				_ = 23
			case 'ㅋ':
				_ = 24
			case 'ㅌ':
				_ = 25
			case 'ㅍ':
				_ = 26
			case 'ㅎ':
				_ = 27
			}
		}
	}
}

var choseongMap = map[rune]int{
	'ㄱ': 0,
	'ㄲ': 1,
	'ㄴ': 2,
	'ㄷ': 3,
	'ㄸ': 4,
	'ㄹ': 5,
	'ㅁ': 6,
	'ㅂ': 7,
	'ㅃ': 8,
	'ㅅ': 9,
	'ㅆ': 10,
	'ㅇ': 11,
	'ㅈ': 12,
	'ㅉ': 13,
	'ㅊ': 14,
	'ㅋ': 15,
	'ㅌ': 16,
	'ㅍ': 17,
	'ㅎ': 18,
}

var jongseongMap = map[rune]int{
	-1:  0,
	'ㄱ': 1,
	'ㄲ': 2,
	'ㄳ': 3,
	'ㄴ': 4,
	'ㄵ': 5,
	'ㄶ': 6,
	'ㄷ': 7,
	'ㄹ': 8,
	'ㄺ': 9,
	'ㄻ': 10,
	'ㄼ': 11,
	'ㄽ': 12,
	'ㄾ': 13,
	'ㄿ': 14,
	'ㅀ': 15,
	'ㅁ': 16,
	'ㅂ': 17,
	'ㅄ': 18,
	'ㅅ': 19,
	'ㅆ': 20,
	'ㅇ': 21,
	'ㅈ': 22,
	'ㅊ': 23,
	'ㅋ': 24,
	'ㅌ': 25,
	'ㅍ': 26,
	'ㅎ': 27,
}
