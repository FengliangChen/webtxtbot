package txtbot

import (
	"errors"
	"github.com/sergeilem/xls"
	"github.com/tealeg/xlsx"
)

func ParseXls(path string) (string, error) {
	if xlFile, err := xls.Open(path, "utf-8"); err == nil {
		if sheet1 := xlFile.GetSheet(0); sheet1 != nil {
			brand = sheet1.Row(2).Col(1)
			xlsContent := ""
			for i := 4; i < 12; i++ {
				cell := sheet1.Row(i).Col(7)
				xlsContent = xlsContent + cell + "\n"
			}
			return xlsContent, nil
		}
	}
	return "", errors.New("Parse xls error.")
}

func ParseXlsx(path string) (string, error) {
	if xlFile, err := xlsx.OpenFile(path); err == nil {
		firstSheet := xlFile.Sheets[0]
		brandCell := firstSheet.Cell(2, 1)
		brand = brandCell.Value
		xlsContent := ""
		for i := 4; i < 12; i++ {
			tempCell := firstSheet.Cell(i, 7)
			xlsContent = xlsContent + tempCell.Value + "\n"
		}
		return xlsContent, nil
	}
	return "", errors.New("Parse xlsx error.")
}
