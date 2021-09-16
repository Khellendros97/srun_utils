package query

import (
	"github.com/sirupsen/logrus"
	"time"

	"github.com/jinzhu/gorm"
)

//TimeSection 时间区间， 时间格式为YYYY-MM-DD hh:mm:ss
type TimeSection struct {
	StartTime string `form:"startDate" json:"startDate"` //开始时间
	StopTime  string `form:"endDate" json:"endDate"`  //结束时间
}

// StartTimeStamp @title StartTimeStamp
// @description 获取起始时间的时间戳
// @auth 黄宇超 时间（2021/4/2）
func (t *TimeSection) StartTimeStamp() (int64, error) {
	// 解析时间
	if dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", t.StartTime, time.Local); err != nil {
		logrus.Error(err)
		return 0, err
	} else {
		return dateTime.Unix(), nil
	}
	//return strconv.ParseInt(t.StartTime, 10, 64)
}

// StopTimeStamp @title StopTimeStamp
// @description 获取结束时间的时间戳
// @auth 黄宇超 时间（2021/4/2）
func (t *TimeSection) StopTimeStamp() (int64, error) {
	if dateTime, err := time.ParseInLocation("2006-01-02 15:04:05", t.StopTime, time.Local); err != nil {
		logrus.Error(err)
		return 0, err
	} else {
		return dateTime.Unix(), nil
	}
	// return strconv.ParseInt(t.StopTime, 10, 64)
}

// SetTimeSectionForQuery @title SetTimeSectionForQuery
// @description 在查询中添加如下语句： WHERE <column> >= StartTimeStamp AND <column> <= StopTimeStamp
// @auth 黄宇超 时间（2021/4/2）
// @param column 时间字段的名称
func (t *TimeSection) SetTimeSectionForQuery(db *gorm.DB, column string) *gorm.DB {
	start, errStart := t.StartTimeStamp()
	stop, errStop := t.StopTimeStamp()

	if errStart == nil {
		db = db.Where(column+" >= ?", start)
	}
	if errStop == nil {
		db = db.Where(column+" <= ?", stop)
	}
	return db
}

func (t *TimeSection) SetTimeSectionForGorm(db *gorm.DB) *gorm.DB {
	if t.StartTime != "" {
		db = db.Where("updated_at >= ?", t.StartTime)
	}
	if t.StopTime != "" {
		db = db.Where("updated_at <= ?", t.StopTime)
	}
	
	return db
}