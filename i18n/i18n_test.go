package i18n_test

import (
	"testing"

	"github.com/Khellendros97/srun_utils/i18n"
)

/*
系统初始化时载入语言和翻译包
*/
func init() {
	// 设置语言为中文，也可以通过配置文件等载入语言配置
	i18n.Local = func()int { return i18n.LOCAL_ZN }
	i18n.Lang = i18n.LangT { // 载入翻译包
		//模块名
		// v
		"name": i18n.ModT {
			//键名            中文翻译       英文翻译
			// v                v             v
			"manager": { Zn: "管理员", En: "Manager" },
			// 注意，请不要在模块名和键名中使用特殊字符和空格
		},
		"action": i18n.ModT {
			"login": { Zn: "登录", En: "Login" },
		},
		"result": i18n.ModT {
			"succeed": { Zn: "成功", En: "Succeed"},
			"failed": { Zn: "失败", En: "Failed" },
		},
	}
}

func TestLangf(t *testing.T) {
	/*
	使用Langf函数配合特殊字符串语法载入翻译
	语法：
		插入翻译：$模块名.键名 如$name.manager就会从翻译包中载入name模块的manager翻译
		连接符：可以使用~在两个翻译之间加入连接符，在中文翻译中，连接符会被忽略；在英文翻译中，连接符会被替换为空格
	此外，Langf也能像Printf一样进行格式化
	*/
	// Langf根据i18n.Local载入翻译
	if i18n.Langf("$name.manager [%s] $action.login~$result.succeed", "hyc") != "管理员 [hyc] 登录成功" {
		t.Fail()
	}
	// 使用Znf载入中文
	if i18n.Znf("$name.manager [%s] $action.login~$result.succeed", "hyc") != "管理员 [hyc] 登录成功" {
		t.Fail()
	}
	// 使用Enf载入英文
	if i18n.Enf("$name.manager [%s] $action.login~$result.succeed", "hyc") != "Manager [hyc] Login Succeed" {
		t.Fail()
	}
}