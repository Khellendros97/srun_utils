package excel

import (
	"fmt"
	"strings"
	"github.com/xuri/excelize/v2"
)

func ReadExcelF(fname, sheet string, cb func(int, map[string]interface{})) (err error) {
	cols := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	cellname := func(row int, col byte) string {
		return fmt.Sprintf("%s%d", string(col), row)
	}

	fp, err := excelize.OpenFile(fname)
	if err != nil {
		return
	}
	rows, err := fp.GetRows(sheet)
	if err != nil {
		return
	}
	if len(rows) < 1 {
		err = fmt.Errorf("empty sheet")
		return
	}

	header := rows[0]

	for i := range rows[1:] {
		data := make(map[string]interface{})
		var h string
		for j := 0; j < len(header); j++ {
			//fmt.Println(cellname(i + 1, cols[j]))
			h = strings.Trim(strings.TrimSpace(header[j]), "\n")
			data[h], err = fp.GetCellValue(sheet, cellname(i+2, cols[j]))
			if err != nil {
				data[h] = ""
			}
		}
		cb(i + 2, data)
	}
	return
}

func ReadExcel(fname string, sheet string) (data []map[string]interface{}, err error) {
	cols := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	cellname := func(row int, col byte) string {
		return fmt.Sprintf("%s%d", string(col), row)
	}

	fp, err := excelize.OpenFile(fname)
	if err != nil {
		return
	}
	rows, err := fp.GetRows(sheet)
	if err != nil {
		return
	}
	if len(rows) < 1 {
		err = fmt.Errorf("empty sheet")
		return
	}

	header := rows[0]
	data = make([]map[string]interface{}, len(rows)-1)

	for i := range rows[1:] {
		data[i] = make(map[string]interface{})
		var h string
		for j := 0; j < len(header); j++ {
			//fmt.Println(cellname(i + 1, cols[j]))
			h = strings.Trim(strings.TrimSpace(header[j]), "\n")
			data[i][h], err = fp.GetCellValue(sheet, cellname(i+2, cols[j]))
			if err != nil {
				data[i][h] = ""
			}
		}
	}
	return
}

type ExcelColumn struct {
	//Col string
	FieldName string
	FieldShow string
}

func WriteExcel(f *excelize.File, sheet string, columns []ExcelColumn, data []map[string]interface{}) *excelize.File {
	cols := []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	cellname := func(row int, col byte) string {
		return fmt.Sprintf("%s%d", string(col), row)
	}

	//write header
	irow := 1
	for i, col := range columns {
		f.SetCellValue(sheet, cellname(irow, cols[i]), col.FieldShow)
	}
	irow++
	for _, row := range data {
		for j, col := range columns {
			f.SetCellValue(sheet, cellname(irow, cols[j]), row[col.FieldName])
		}
		irow++
	}
	return f
}