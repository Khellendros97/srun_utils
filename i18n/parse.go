package i18n

import "unicode"

const (
	TOKEN_KIND_KEY = 0
	TOKEN_KIND_NORMAL = 1
	TOKEN_KIND_SPACER = 2
	TOKEN_KIND_END = 3

	LEXER_STATUS_INIT = 0
	LEXER_STATUS_NORMAL = 1
	LEXER_STATUS_MODULE = 2
	LEXER_STATUS_KEY = 3
)

type token struct {
	Kind int
	Raw []rune
	Module string
	Key string
}

type Lexer struct {
	Line []rune
	Pos int
}

func (lex *Lexer) SetLine(line string) {
	lex.Line = []rune(line)
	lex.Pos = 0
}

func (lex *Lexer) Parse(line string) (ts []token) {
	lex.SetLine(line)
	ts = make([]token, 0)
	for {
		t := lex.GetToken()
		if t.Kind == TOKEN_KIND_END {
			return
		}
		ts = append(ts, t)
	}
}

func (lex *Lexer) GetToken() (t token) {
	status := LEXER_STATUS_INIT
	t.Kind = TOKEN_KIND_END
	var cur rune

	for lex.Pos < len(lex.Line) {
		cur = lex.Line[lex.Pos]

		if status == LEXER_STATUS_MODULE {
			if cur == '.' {
				t.Module = string(t.Raw)
				t.Raw = []rune{}
				status = LEXER_STATUS_KEY
				lex.Pos++
				continue
			}
		}

		if status == LEXER_STATUS_KEY {
			if !unicode.IsLetter(cur) && cur != '_' {
				break
			}
		}

		if status == LEXER_STATUS_NORMAL {
			if cur == '$' || cur == '~' {
				break
			}
		}

		if status == LEXER_STATUS_INIT {
			if cur == '$' {
				t.Kind = TOKEN_KIND_KEY
				status = LEXER_STATUS_MODULE
				lex.Pos++
				continue
			}

			if cur == '~' {
				t.Kind = TOKEN_KIND_SPACER
				lex.Pos++
				break
			}

			
			t.Kind = TOKEN_KIND_NORMAL
			status = LEXER_STATUS_NORMAL
		}
		lex.Pos++
		t.Raw = append(t.Raw, cur)
	}

	if t.Kind == TOKEN_KIND_KEY {
		t.Key = string(t.Raw)
	}

	return
}