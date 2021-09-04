package query

import (
	"fmt"
)

// Like @title Like
// @description 将str转化为%str%
// @auth 黄宇超 时间（2021/4/2）
func Like(s string) string {
	return fmt.Sprintf("%s%%", s)
}