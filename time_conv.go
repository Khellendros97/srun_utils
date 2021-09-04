package utils

import (
	"time"
)

// @title TimeFromInt64
// @description 时间戳 -> time.Time
// @auth 黄宇超 时间（2021/4/2）
func TimeFromInt64(n int64) time.Time {
	return time.Unix(n, 0)
}

// @title TimeToString
// @description time.Time -> string， 格式为 YYYY-MM-DD hh:mm:ss
// @auth 黄宇超 时间（2021/4/2）
func TimeToString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05")
}

// @title TimeFromString
// @description string -> time.Time
// @auth 黄宇超 时间（2021/4/2）
func TimeFromString(s string) (time.Time, error) {
	return time.ParseInLocation("2006-01-02 15:04:05", s, time.Local)
}

// @title TimeStampToTimeString
// @description 将时间戳转为可读的字符串
// @auth 黄宇超 时间（2021/4/2）
func TimeStampToTimeString(timestamp int64) string {
	return TimeToString(TimeFromInt64(timestamp))
}

// @title TimeStringToTimeStamp
// @description 将字符串转化为时间戳
// @auth 黄宇超 时间（2021/4/2）
func TimeStringToTimeStamp(timestr string) (timestamp int64, err error) {
	t, err := TimeFromString(timestr)
	if err != nil {
		return
	}
	timestamp = t.Unix()
	return
}

func NowString() string {
	return TimeToString(time.Now())
}

//@title TimeStringToInt
//@description  将时间戳转化为年月日 1625711519(2021-07-08 10:31:59) -> 1625673600(2021-07-08)
//@autor ZHANGYAN
//@time  2021/7/8 13:26
func TimeStringToInt(timeInt int64, timeFormat string) (timestamp int64) {
	timeDate := time.Unix(timeInt, 0).Local().Format(timeFormat)
	timeDateInt, _ := time.ParseInLocation(timeFormat, timeDate, time.Local)

	return timeDateInt.Unix()
}

//@title SegmentByDate
//@description 将一个时间段分割到数组中
//@autor ZHANGYAN
//@time  2021/7/8 15:08
func SegmentByDate(StartTime int64, EndTime int64, Unit string, Step int64) (date []string) {
	switch Unit {
	case "minute":
		for {
			if StartTime <= EndTime {
				row := time.Unix(StartTime, 0).Local().Format("2006-01-02 15:04:05")
				date = append(date, row)
				if StartTime == EndTime {
					break
				}
			}
			StartTime = StartTime + Step*60
		}
	}
	return
}
