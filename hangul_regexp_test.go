package hangul_regexp

import "testing"

func TestGetPattern(t *testing.T) {
	type args struct {
		search      string
		ignoreSpace bool
		fuzzy       bool
		choseong    bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{"Alphabet", args{"a", false, false, false}, "a", false},
		{"Hangul", args{"ㄱ나다라123", false, false, false}, "ㄱ나다라123", false},
		{"Mixed", args{"Zx0ㅡㅡ", false, false, false}, "Zx0ㅡㅡ", false},
		{"Should escape", args{"[^가-힣]$", false, false, false}, "\\[\\^가-힣]\\$", false},

		{"Last char is choseong", args{"ㄱ", false, false, false}, "(?:ㄱ|[가-깋])", false},
		{"Last char without batchim", args{"가 나", false, false, false}, "가 (?:나|[낙-낳])", false},
		{"Last char with batchim", args{"가 안", false, false, false}, "가 (?:안|아(?:ㄴ|[나-닣]))", false},
		{"Last char with double batchim", args{"가 있", false, false, false}, "가 (?:있|이(?:ㅆ|[싸-앃]))", false},
		{"Last char with combined batchim", args{"가 얇", false, false, false}, "가 (?:얇|얄(?:ㅂ|[바-빟]))", false},

		{"Mixed / ignoreSpace=true", args{"ㅁ가a항1", true, false, false}, "ㅁ *?가 *?a *?항 *?1", false},
		{"Last char with batchim / ignoreSpace=true / has space matcher between", args{"가 안", true, false, false}, "가 *?  *?(?:안|아 *?(?:ㄴ|[나-닣]))", false},

		{"Mixed / fuzzy=true", args{"ㅁ가a항1", false, true, false}, "ㅁ.*?가.*?a.*?항.*?1", false},
		{"Space / fuzzy=true / spaces are not concatenated", args{"가 s", false, true, false}, "가.*? .*?s", false},
		{"Last char with batchim / fuzzy=true / has any matcher between", args{"가 안", false, true, false}, "가.*? .*?(?:안|아.*?(?:ㄴ|[나-닣]))", false},

		{"Non-last char is choseong / choseong=true", args{"ㄱ1", false, false, true}, "(?:ㄱ|[가-깋])1", false},
		{"Multiple choseong chars / choseong=true", args{"ㄱ ㄴㄷ", false, false, true}, "(?:ㄱ|[가-깋]) (?:ㄴ|[나-닣])(?:ㄷ|[다-딯])", false},
		{"Mixed with choseong / choseong=true", args{"aㅎ1가ㄴ", false, false, true}, "a(?:ㅎ|[하-힣])1가(?:ㄴ|[나-닣])", false},

		{"Any / ignoreSpace=true, fuzzy=true / err", args{"", true, true, false}, "", true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := GetPattern(tt.args.search, tt.args.ignoreSpace, tt.args.fuzzy, tt.args.choseong)
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
