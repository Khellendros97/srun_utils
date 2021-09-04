package utils

import (
	"bufio"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"sync"

	log "github.com/sirupsen/logrus"
)

/**
操作文件相关
*/
var (
	bt = sync.Pool{
		New: func() interface{} {
			b := make([]byte, 1024)
			return &b
		},
	}
)

/**
读取配置文件
*/
func ReadFile(fileName string) (m map[string]interface{}, err error) {
	var (
		b *[]byte
		f *os.File
		r *bufio.Reader
	)
	if _, err = os.Stat(fileName); err != nil && os.IsNotExist(err) {
		err = nil
		return
	}
	m = make(map[string]interface{})
	f, err = os.OpenFile(fileName, os.O_RDONLY, 0666)
	if err != nil {
		log.Error(err)
		return
	}
	defer f.Close()
	r = bufio.NewReader(f)
	for {
		b = bt.Get().(*[]byte)
		*b, err = r.ReadBytes('\n')
		if r := doStringToArray(string(*b)); len(r) > 0 {
			m[r[0]] = r[1]
		}
		*b = (*b)[:0]
		bt.Put(b)
		if err != nil {
			if err == io.EOF {
				return m, nil
			}
			log.Error(err)
		}

	}
}

//ReadListFile 读取列表文件
func ReadListFile(fname string) (strs []string, err error) {
	fp, err := os.Open(fname)
	if err != nil {
		return
	}
	defer fp.Close()

	reader := bufio.NewReader(fp)
	for {
		line, _, err := reader.ReadLine()
		if err == nil {
			strs = append(strs, strings.TrimSpace(string(line)))
		} else if err == io.EOF {
			break
		} else {
			return nil, err
		}
	}
	if err == io.EOF {
		err = nil
	}
	return
}

/**
string reserve to Array struct
*/
func doStringToArray(r string) (s []string) {
	//先判断字符串是否以#开头，注释
	if b := strings.HasPrefix(r, "#"); b {
		return nil
	}
	// 分号也是注释，用于企业微信配置文件
	if b := strings.HasPrefix(r, ";"); b {
		return nil
	}
	result := strings.Split(r, "=")
	if len(result) < 2 {
		return nil
	}
	ss := make([]string, 2)
	re3, _ := regexp.Compile("\"|(\r\n)|\n")

	// 避免出现值中出现=的情况
	var key, value string
	for k, v := range result {
		if k == 0 {
			key = re3.ReplaceAllString(v, "")
			key = strings.TrimSpace(key)
		} else if value == "" {
			value = re3.ReplaceAllString(v, "")
			value = strings.TrimSpace(value)
		} else {
			valueNew := re3.ReplaceAllString(v, "")
			valueNew = strings.TrimSpace(valueNew)
			value = value + "=" + valueNew
		}
	}

	ss[0] = key
	ss[1] = value
	return ss
}

// @title IsExist
// @description checks whether a file or directory exists
// @author DM
// @time 2021/4/9 13:36
// @param f
// @return bool It returns false when the file or directory does not exist
func IsExist(f string) bool {
	_, err := os.Stat(f)
	return err == nil || os.IsExist(err)
}

// @title IsFile
// @description checks whether the path is a file
// @author DM
// @time 2021/4/9 13:37
// @param f
// @return bool it returns false when it's a directory or does not exist
func IsFile(f string) bool {
	fi, e := os.Stat(f)
	if e != nil {
		return false
	}
	return !fi.IsDir()
}

// @title IsDir
// @description 判断所给路径是否为文件夹
// @author DM
// @time 2021/4/9 13:50
// @param path
// @return bool
func IsDir(path string) bool {
	s, err := os.Stat(path)
	if err != nil {
		return false
	}
	return s.IsDir()
}

//ListFiles 列出目录下所有文件（非递归）
func ListFiles(path string, pred ...func(os.FileInfo)bool) (infos [] os.FileInfo) {
	infos = make([]os.FileInfo, 0)
	finfos, err := ioutil.ReadDir(path)
	if err != nil {
		fmt.Println(err)
		return
	}
	if len(pred) < 1 {
		infos = finfos
		return
	}

	p := pred[0]
	for _, info := range finfos {
		if p(info) {
			infos = append(infos, info)
		}
	}
	return
}