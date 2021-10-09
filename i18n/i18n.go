package i18n

import (
	"fmt"
	"strings"
)

const (
	LOCAL_EN = 0
	LOCAL_ZN = 1
)

type Translate struct {
	Zn string
	En string
}

type I18N struct {
	Name   map[string]Translate
	Action map[string]Translate
	Status map[string]Translate
}

type LangT = map[string]ModT
type ModT = map[string]Translate

var Lang LangT

var Local func()int = func()int { return LOCAL_EN }

func getModule(mod string) ModT {
	m, ok := Lang[mod]
	if !ok {
		m = ModT{}
	}
	return m
}

func Znf(pattern string, args ...interface{}) string {
	var lexer Lexer
	tokens := lexer.Parse(pattern)
	zns := make([]string, len(tokens))
	for i, t := range tokens {
		if t.Kind == TOKEN_KIND_KEY {
			module := getModule(t.Module)
			v, ok := module[t.Key]
			if !ok {
				zns[i] = t.Key
			} else {
				zns[i] = v.Zn
			}
		} else if t.Kind == TOKEN_KIND_SPACER {
			//中文忽略分隔符
			zns[i] = ""
		} else {
			zns[i] = string(t.Raw)
		}
	}
	p := strings.Join(zns, "")
	return fmt.Sprintf(p, args...)
}

func Enf(pattern string, args ...interface{}) string {
	var lexer Lexer
	tokens := lexer.Parse(pattern)
	ens := make([]string, len(tokens))
	for i, t := range tokens {
		if t.Kind == TOKEN_KIND_KEY {
			module := getModule(t.Module)
			v, ok := module[t.Key]
			if !ok {
				ens[i] = t.Key
			} else {
				ens[i] = v.En
			}
		} else if t.Kind == TOKEN_KIND_SPACER {
			//英文插入分隔符
			ens[i] = " "
		} else {
			ens[i] = string(t.Raw)
		}
	}
	p := strings.Join(ens, "")
	return fmt.Sprintf(p, args...)
}

func Langf(pattern string, args ...interface{}) string {
	switch Local() {
	case LOCAL_EN:
		return Enf(pattern, args...)
	case LOCAL_ZN:
		return Znf(pattern, args...)
	}
	return Enf(pattern, args...)
}