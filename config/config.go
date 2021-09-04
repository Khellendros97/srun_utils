package config

import (
	"fmt"
	"os"
	utils "github.com/srunsoft/srun_utils"
	"strconv"
	"strings"

	//"log"
	"errors"
)

type Config struct {
	cfg map[string]string
}

func New() *Config {
	return &Config{
		cfg: make(map[string]string),
	}
}

func (c *Config) Set(key string, value interface{}) (err error) {
	var val string
	switch value.(type) {
	case string:
		val = value.(string)
	case int:
		val = fmt.Sprintf("%d", value.(int))
	case bool:
		if value.(bool) {
			val = "1"
		} else {
			val = "0"
		}
	case []interface{}:
		strs := make([]string, len(value.([]interface{})))
		for i, x := range value.([]interface{}) {
			strs[i] = fmt.Sprintf("%v", x)
		}
		val = strings.Join(strs, ",")
	case []int:
		strs := make([]string, len(value.([]int)))
		for i, x := range value.([]int) {
			strs[i] = fmt.Sprintf("%v", x)
		}
		val = strings.Join(strs, ",")
	case []string:
		strs := make([]string, len(value.([]string)))
		for i, x := range value.([]string) {
			strs[i] = fmt.Sprintf("%v", x)
		}
		val = strings.Join(strs, ",")
	default:
		return errors.New("invalid value type ,only support string, int or bool")
	}
	c.cfg[key] = val
	return
}

func (c *Config) GetStringOrElse(key string, orelse string) (value string) {
	value, ok := c.cfg[key]
	if !ok {
		value = orelse
		return
	}
	return
}

func (c *Config) GetIntOrElse(key string, orelse int) (value int) {
	str, ok := c.cfg[key]
	if !ok {
		value = orelse
		return
	}
	value, err := strconv.Atoi(str)
	if err != nil {
		value = orelse
	}
	return
}

func (c *Config) GetBool(key string) (value bool) {
	n := c.GetIntOrElse(key, 0)
	if n == 0 {
		return false
	} else {
		return true
	}
}

func (c *Config) GetStrings(key string) (value []string) {
	str := c.GetStringOrElse(key, "")
	value = strings.Split(str, ",")
	return
}

func (c *Config) GetInts(key string, dft int) (value []int) {
	strs := c.GetStrings(key)
	value = make([]int, len(strs))
	for i, str := range strs {
		var err error
		value[i], err = strconv.Atoi(str)
		if err != nil {
			value[i] = dft
		}
	}
	return
}

func parseLine(line string) (key string, value string, ok bool) {
	strs := strings.Split(line, "=")
	if len(strs) != 2 {
		ok = false
		return
	}
	key = strings.Trim(strs[0], " \n\r\t")
	//value = strings.Trim(strings.ReplaceAll(strs[1], "\"", ""), " \r\n\t")
	value = strings.ReplaceAll(strs[1], "\"", "")
	ok = true
	return
}

func Load(fname string) (c *Config, err error) {
	lines, err := utils.ReadListFile(fname)
	if err != nil {
		return
	}
	c = New()
	for _, line := range lines {
		//log.Println(line)
		k, v, ok := parseLine(line)
		//log.Println(k, ": ", v)
		if ok {
			c.Set(k, v)
		}
	}
	return
}

func (c *Config) Save(dstFname string) (err error) {
	fp, err := os.Create(dstFname)
	if err != nil {
		return
	}
	defer fp.Close()
	for k, v := range c.cfg {
		line := fmt.Sprintf("%s=\"%s\"\n", k, v)
		fp.WriteString(line)
	}
	return
}
