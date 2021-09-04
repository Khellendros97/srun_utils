package utils

import "regexp"

func PipeValidators(validators ...func(string)bool) func(string)bool {
	return func (in string) bool {
		result := true
		for _, v := range validators {
			result = result && v(in)
			if !result {
				break
			}
		}
		return result
	}
}

func Test(re string) func(string)bool {
	m := regexp.MustCompile(re)
	return func (in string) bool {
		return m.MatchString(in)
	}
}