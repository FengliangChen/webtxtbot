/*
 @Author      : Simon Chen
 @Email       : bafelem@gmail.com
 @datetime    : 2021-09-17
 @Description : Description
 @FileName    : autofill.go
*/

package autofill

// package main

import (
	"encoding/json"
	"fmt"
	"github.com/otiai10/copy"
	"github.com/xuri/excelize/v2"
	"os"
	"path/filepath"
	"strings"
	"time"
)

const (
	seasonalFormPath = "codeTable.xlsx"
	tempaltePath     = "template.xlsx"
	newtemplatePath  = "Newtemplate.xlsx"
)

var destopPath = filepath.Join(os.Getenv("HOME"), "Desktop")

type Jobinfo struct {
	job         string
	seasonCode  string
	country     string
	createDay   string
	program     string
	supplier    string
	buyer       string
	dept        string
	dueDate     string
	packoutDate string
	shipDate    string
	instoreDate string
	contact     string
	saveTo      string
}

func NewJob(formData map[string]string) Jobinfo {
	job := Jobinfo{
		job:         formData["job"],
		seasonCode:  formData["season_code"],
		country:     formData["country"],
		createDay:   formData["create_date"],
		program:     formData["program"],
		supplier:    formData["supplier"],
		buyer:       formData["buyer"],
		dept:        formData["department"],
		dueDate:     formData["artdue_date"],
		packoutDate: formData["packout_date"],
		shipDate:    formData["ship_date"],
		instoreDate: formData["instore_date"],
		contact:     formData["contact"],
		saveTo:      "",
	}
	return job
}

func (a *Jobinfo) SetCreateDate() {
	if len(a.createDay) == 0 {
		a.createDay = GetToday("01/02/2006")
	}
}

func GetToday(format string) string { // yy=2006 mm=01 dd=02
	now := time.Now()
	now = now.Add(time.Hour * -8) //coodinate with US office.
	return now.Format(format)
}

func (a *Jobinfo) SaveForm(templatePath string) error {
	f, err := excelize.OpenFile(templatePath)
	if err != nil {
		return err
	}
	f.SetCellValue("Sheet1", "B2", "Job #:"+a.job)
	f.SetCellValue("Sheet1", "B3", a.seasonCode+"-"+a.country)
	f.SetCellValue("Sheet1", "A5", a.createDay)
	f.SetCellValue("Sheet1", "H4", "Job #: "+a.job)
	f.SetCellValue("Sheet1", "H5", "Program: "+a.program)
	f.SetCellValue("Sheet1", "H6", "Supplier: "+a.supplier)
	f.SetCellValue("Sheet1", "H7", "Buyer: "+a.buyer+"("+"D"+a.dept+")")
	f.SetCellValue("Sheet1", "H8", "Artwork due date: "+a.dueDate)
	f.SetCellValue("Sheet1", "H9", "Packout date: "+a.packoutDate)
	f.SetCellValue("Sheet1", "H10", "Shipdate: "+a.shipDate)
	f.SetCellValue("Sheet1", "H11", "In-store date: "+a.instoreDate)
	f.SetCellValue("Sheet1", "H13", "联系人: "+a.contact)
	f.Save()
	return nil
}

func (a *Jobinfo) MakeJobFolder() error {
	parentFolder := filepath.Join(a.saveTo, a.job+" 做稿")
	today := GetToday("0102")
	intakefolder := filepath.Join(parentFolder, "1 intake sheet & order", today)
	rawfolder := filepath.Join(parentFolder, "2 raw client files", today)

	if err := os.MkdirAll(intakefolder, os.ModePerm); err != nil {
		return err
	}
	if err := os.MkdirAll(rawfolder, os.ModePerm); err != nil {
		return err
	}
	sheetName := a.job + "_DetailList_W.xlsx"
	sheetPath := filepath.Join(parentFolder, sheetName)
	MoveFile(newtemplatePath, sheetPath)
	return nil
}

func ReadSeasonalMap(path string) (map[string]string, error) {
	ConvertionMap := map[string]string{}
	f, err := excelize.OpenFile(path)
	if err != nil {
		return ConvertionMap, err
	}
	rows, err := f.GetRows("Sheet1")
	if err != nil {
		return ConvertionMap, err
	}
	for _, row := range rows {
		if len(row) > 1 {
			indexValue := row[0]
			mapedValue := row[1]
			if len(indexValue) > 0 && len(mapedValue) > 0 {
				ConvertionMap[indexValue] = mapedValue
			}

		}

	}
	return ConvertionMap, nil
}

func GetSeasonalMap_json() (string, error) {
	var jstring string
	programMap, err := ReadSeasonalMap(seasonalFormPath)

	if err != nil {
		return jstring, err
	}

	j, err := json.Marshal(programMap)

	if err != nil {
		return jstring, err
	}

	jstring = string(j)

	return jstring, nil
}

func GetBuyerList_json() (string, error) {
	var jstring string
	buyerList, err := ReadBuyerList(seasonalFormPath)

	if err != nil {
		return jstring, err
	}

	j, err := json.Marshal(buyerList)

	if err != nil {
		return jstring, err
	}

	jstring = string(j)

	return jstring, nil
}

func ReadBuyerList(path string) ([]string, error) {
	var buyerList []string

	f, err := excelize.OpenFile(path)
	if err != nil {
		return buyerList, err
	}
	rows, err := f.GetRows("Sheet2")
	if err != nil {
		return buyerList, err
	}
	for _, row := range rows {
		if len(row) > 0 {
			buyerList = append(buyerList, row[0])
		}
	}
	return buyerList, nil
}

func (a *Jobinfo) SetProgram() {
	if len(a.program) == 0 {
		programMap, err := ReadSeasonalMap(seasonalFormPath)
		if err != nil {
			fmt.Println(err.Error())
		}
		if value, iscontain := programMap[a.seasonCode]; iscontain {
			a.program = value
		}

	}
}

func (a *Jobinfo) SetCountry() {
	if len(a.country) == 0 {
		a.country = "China"

	}
}

func (a *Jobinfo) SetSaveTo() {
	if len(a.saveTo) == 0 {
		a.saveTo = destopPath

	}
}

func (a *Jobinfo) SetUpperCase() {
	a.seasonCode = strings.ToUpper(a.seasonCode)
}

func (a *Jobinfo) SetTitleCase() {
	a.supplier = strings.Title(a.supplier)
	a.buyer = strings.Title(a.buyer)

}

func MoveFile(src string, dst string) {
	os.Rename(src, dst)
}

func (a *Jobinfo) Init() {
	// if value is null, set the default value.
	a.SetCreateDate()
	a.SetProgram()
	a.SetCountry()
	a.SetSaveTo()
	a.SetUpperCase()
	a.SetTitleCase()
}
func (a *Jobinfo) MakeJob() error {
	err := copy.Copy(tempaltePath, newtemplatePath)
	if err != nil {
		return err
	}
	a.SaveForm(newtemplatePath)
	err = a.MakeJobFolder()
	if err != nil {
		return err
	}
	return nil

}

// func main() {
// 	a := Jobinfo{
// 		job:        "C201118_HOL",
// 		seasonCode: "HOL",
// 		// saveTo:     destopPath,
// 	}
// 	a.Init()
// 	a.MakeJob()

// }
