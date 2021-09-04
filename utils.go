package utils

import (
	"bytes"
	"encoding/csv"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/xuri/excelize/v2"
	"math/big"
	"math/rand"
	"net"
	"os"
	"reflect"
	"strings"
	"time"
)

// 获取秘钥
func CreateSecret(n int) string {
	var letters = []rune("abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789")
	b := make([]rune, n)
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

// 将结构体转换为map
func StructToMap(obj interface{}) map[string]string {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]string)

	switch obj1.Kind() {
	case reflect.Ptr:
		for i := 0; i < obj1.Elem().NumField(); i++ {
			if obj2.Elem().Field(i).Interface() != nil {
				data[obj1.Elem().Field(i).Tag.Get("json")] = fmt.Sprintf("%v", obj2.Elem().Field(i).Interface())
			}
		}
	default:
		for i := 0; i < obj1.NumField(); i++ {
			if obj2.Field(i).Interface() != nil {
				data[obj1.Field(i).Tag.Get("json")] = fmt.Sprintf("%v", obj2.Field(i).Interface())
			}
		}
	}
	return data
}

// 将结构体转换为map
func StructToMapValue(obj interface{}) map[string]string {
	obj1 := reflect.TypeOf(obj)
	obj2 := reflect.ValueOf(obj)

	var data = make(map[string]string)

	switch obj1.Kind() {
	case reflect.Ptr:
		for i := 0; i < obj1.Elem().NumField(); i++ {
			if obj2.Elem().Field(i).Interface() != nil && obj2.Elem().Field(i).Interface() != "" {
				data[obj1.Elem().Field(i).Tag.Get("json")] = fmt.Sprintf("%v", obj2.Elem().Field(i).Interface())
			}
		}
	default:
		for i := 0; i < obj1.NumField(); i++ {
			if obj2.Field(i).Interface() != nil && obj2.Field(i).Interface() != "" {
				data[obj1.Field(i).Tag.Get("json")] = fmt.Sprintf("%v", obj2.Field(i).Interface())
			}
		}
	}
	return data
}

func MapToString(m map[string]interface{}) string {
	var (
		b     strings.Builder
		count int
	)
	b.WriteString("{")
	for k, v := range m {
		switch reflect.ValueOf(v).Kind() {
		case reflect.Uint:
			b.WriteString(fmt.Sprintf("\"%s\":%d", k, v.(uint)))
		case reflect.Int:
			b.WriteString(fmt.Sprintf("\"%s\":%d", k, v.(int)))
		case reflect.Int64:
			b.WriteString(fmt.Sprintf("\"%s\":%d", k, v.(int64)))
		case reflect.String:
			b.WriteString(fmt.Sprintf("\"%s\":\"%s\"", k, v.(string)))
		case reflect.Bool:
			b.WriteString(fmt.Sprintf("\"%s\":%v", k, v.(bool)))
		}
		if count != len(m)-1 {
			b.WriteString(",")
		}
		count++
	}
	b.WriteString("}")
	return b.String()
}

func InArray(key string, arr []string) bool {
	for _, v := range arr {
		if key == v {
			return true
		}
	}

	return false
}

// 导出文件到csv中
func Download(filename string, data [][]string, headers []string, g *gin.Context) string {
	logrus.WithFields(logrus.Fields{
		"headers": headers,
		"data":    data,
	}).Info("download csv")
	// 定义文件名
	g.Writer.Header().Set("Content-Disposition", fmt.Sprintf("attachment;filename=%s", filename))
	g.Writer.Header().Add("Content-Type", "application/csv")

	// 将传进来的

	b := &bytes.Buffer{}
	b.WriteString("\xEF\xBB\xBF") // UTF-8 BOM 处理中文乱码
	wr := csv.NewWriter(b)
	wr.Write(headers) // 按行shu
	wr.WriteAll(data)
	/*for i := 0; i < len(data); i++ {
		wr.Write(data[i])
	}*/
	wr.Flush()
	return b.String()
}

// 往csv文件中写数据
func WriteCsv(filename string, data [][]string, headers []string) error {
	file, err := os.OpenFile(filename, os.O_CREATE|os.O_RDWR|os.O_TRUNC, 0766)
	if err != nil {
		return err
	}
	defer file.Close()
	wr := csv.NewWriter(file)
	wr.Write(headers)
	for i := 0; i < len(data); i++ {
		wr.Write(data[i])
	}
	wr.Flush()
	return nil
}

// 获取表名
func GetTableName(tableName string, params map[string]string) string {
	endTime, err := GetStopTime(params["start_time"], params["stop_time"])
	if err != nil {
		logrus.Error(err)
		return tableName
	}
	now := time.Now().Format("200601")
	if strings.EqualFold(now, endTime.Format("200601")) {
		return tableName
	}

	return fmt.Sprintf("%s_%s", tableName, endTime.Format("200601"))
}

// 有开始时间和结束时间时，比较开始时间是否大于结束时间，大于，则提示开始查询时间不能大于结束时间
// 否则，比较开始时间和结束时间的月份是否相同，若不同，则提示开始查询时间和结束时间不在同一个月
// 否则，若有开始时间没有结束时间，则结束时间等于开始时间+1天
// 若有结束时间没有开始时间，则直接返回
// 若既没有开始时间也没有结束时间，结束时间为当前时间
func GetStopTime(startTimeStr, stopTimeStr string) (stopTime time.Time, err error) {
	var (
		startTimeOk, stopTimeOk bool
		startTime               time.Time
	)

	if len(startTimeStr) > 0 && startTimeStr != "" {
		startTime, err = time.ParseInLocation("2006-01-02", startTimeStr, time.Local)
		if err != nil {
			logrus.Error("开始查询时间格式有误")
		} else {
			startTimeOk = true
		}
	}

	if len(stopTimeStr) > 0 && stopTimeStr != "" {
		stopTime, err = time.ParseInLocation("2006-01-02", stopTimeStr, time.Local)
		if err != nil {
			logrus.Error("结束查询时间格式有误")
		} else {
			stopTimeOk = true
		}
	}

	// 当前查询开始时间大于结束时间
	if startTimeOk && stopTimeOk {
		if stopTime.Before(startTime) {
			return stopTime, errors.New("开始查询时间不能大于结束时间")
		}
		// 不在同一个月
		if stopTime.Format("2006-01") > startTime.Format("2006-01") {
			return stopTime, errors.New("开始查询时间和结束时间不在同一个月")
		}
		// 有开始时间没有结束时间
	} else if startTimeOk && !stopTimeOk {
		stopTime = startTime.Add(time.Hour * 24)
		// 有结束时间没有开始时间
	} else if !startTimeOk && !stopTimeOk {
		stopTime = time.Now()
	}

	return stopTime, nil
}

// 根据时间和单位转换成秒
func Time2Second(time int, unit, method string) int {
	var unitInt int
	switch unit {
	case "s":
		unitInt = 1
	case "m":
		unitInt = 60
	case "h":
		unitInt = 3600
	case "day":
		unitInt = 86400
	case "week":
		unitInt = 604800
	case "month":
		unitInt = 2592000
	case "year":
		unitInt = 31536000
	default:
		unitInt = 1
	}
	if method == "t2s" {
		return time * unitInt
	} else if method == "s2t" {
		return time / unitInt
	}
	return time
}

// 根据流量值和单位转换成相应的kb
func Flow2Kb(flow int64, unit, method string) int64 {
	var unitInt int64
	switch unit {
	case "b":
		unitInt = 1
	case "k":
		unitInt = 1024
	case "m":
		unitInt = 1048576
	case "g":
		unitInt = 1073741824
	case "t":
		unitInt = 1099511627776
	default:
		unitInt = 1
	}
	if method == "f2k" {
		return flow * unitInt
	} else if method == "k2f" {
		return flow / unitInt
	}
	return flow
}

// @title 下载多工作表模板文件
// @author zhanglianfeng 2021/7/16
func WriteMultipleSheetTemplate(filePath, fileName string, data map[string][]string) (string, error) {
	// Create a new sheet
	if !IsExist(filePath) {
		if err := os.MkdirAll(filePath, os.ModePerm); err != nil {
			return "", err
		}
	}
	url := fmt.Sprintf("%s/%s", filePath, fileName)
	excel := excelize.NewFile()
	excel.SetSheetName("Sheet1", "SRUNSOFTSheet1")
	for k, v := range data {
		excel.NewSheet(k)
		excel.DeleteSheet("SRUNSOFTSheet1")
		excel.SetSheetRow(k, "A1", &v)
	}
	if err := excel.SaveAs(url); err != nil {
		return "", err
	}
	return url, nil
}

// @title 数字类型转字符串类型ip地址
// @author zhanglianfeng 2021/7/20
func ParseInttoIpv4(ip int64) string {
	return fmt.Sprintf("%d.%d.%d.%d", byte(ip>>24), byte(ip>>16), byte(ip>>8), byte(ip))
}

// @title 字符串类型转数字类型ip地址
// @author zhanglianfeng 2021/7/20
func ParseIpv4toInt(ip string) int64 {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To4())
	return ret.Int64()
}

// @title 字符串类型转数字类型ip地址
// @author zhanglianfeng 2021/7/20
func ParseIpv6toInt(ip string) *big.Int {
	ret := big.NewInt(0)
	ret.SetBytes(net.ParseIP(ip).To16())
	return ret
}

// @title 111111~~999999六位随机数
// @author zhanglianfeng 2021/8/24
func GetSerialCode() int {
	rand.Seed(time.Now().UnixNano())
	var serialCode int
	for {
		serialCode = rand.Intn(1000000)
		if serialCode >= 111111 && serialCode <= 999999 {
			break
		}
	}
	return serialCode
}
