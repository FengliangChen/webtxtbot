package txtbot

import (
	// "bufio"
	"encoding/json"
	"errors"
	// "fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"time"
)

const (
	dfpath      = "/Volumes/datavolumn_bmkserver_Pub"
	wks         = "/Volumes/datavolumn_bmkserver_Pub/新做稿/未开始"
	jxz         = "/Volumes/datavolumn_bmkserver_Pub/新做稿/进行中"
	DEFAULTSIZE = 64 // Json, Estimated 2000 client if consider 30 bit per client.
)

var (
	now       = time.Now()
	today     = now.Format("0102")
	yesterday = now.AddDate(0, 0, -1).Format("0102")
	month     = now.Format("200601")
	re        *regexp.Regexp
	job       string
	brand     string
	emailSave = filepath.Join(os.Getenv("HOME"), "Desktop", "draftartwork.txt")
	jsonPath  = "/Volumes/datavolumn_bmkserver_Pub/新做稿/已结束/NON-WMT/.database/clientcode.json"
	// filepath.Join(os.Getenv("HOME"), "Documents", "txtbot", "clientcode.json")
	rvst = false
)


type TxtJson struct {
	PHQ string
	TxtBodies []TxtBody
}

type TxtBody struct {
	TxtCount int
	TxtBody string

}

func Run(queryJob string, rv bool) string{
	rvst = rv //cancle some process.

	if !TestConnect() {
		return "Connection Errors, Please check if the server is connected at: " + dfpath
	}
	
	if len(queryJob) != 6 {
		return "Please input 6 digits"
	}

	job = queryJob
	job = strings.ToUpper(job)

	re = regexp.MustCompile(job)

	DFjobpath, err := FetchJobPath()
	if err != nil {
		log.Println(err)
		return job + " ---> Please check today and yesterday, if job folder is existed!"
	}

	PFpath, err := FetchPFpath()
	if err != nil {
		return err.Error()
	}

	txtCount, txtFilePath, err := FetchTxtpath(DFjobpath)
	if err != nil {
		log.Println(err)
		if err.Error() == "nofile of .txt" {
			log.Println("No txt file in job folder, creating base on existing pdf files.")
		}
	}

	var allbody *AllBody
	switch len(txtFilePath) {
	case 0:
		allbody, err = ConstructPDFName(DFjobpath)
		if err != nil {
			return err.Error()
		}

	default:
		allbody, err = FetchBody(txtFilePath, txtCount)
		if err != nil {
			return err.Error()
		}
	}

	tail, err := FetchTail(PFpath)
	if err != nil {
		return err.Error()
	}

	emailTxtJson := CombineAll(allbody, &tail)

	PHQtitle, err := PHQtitle(DFjobpath)
	if err != nil {
		log.Println(err)
	}
	// *emailTxt = PHQtitle + "\n\n" + *emailTxt
	emailTxtJson.PHQ = PHQtitle

	b, err := json.Marshal(emailTxtJson)
	if err != nil {
        return err.Error()
    }
    emailTxt := string(b)

	if true {
		cmd := exec.Command("open", DFjobpath)
		err := cmd.Run()
		if err != nil {
			log.Println(err)
		}
	}
	return emailTxt
}

func TestConnect() bool {
	result := Exists(dfpath)
	return result
}

func Exists(path string) bool {
	_, err := os.Stat(path)
	if err != nil {
		if os.IsExist(err) {
			return true
		}
		return false
	}
	return true
}

func SearchJob(path string) (string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return "", err
	}
	for _, file := range files {
		if file.IsDir() {
			if re.MatchString(file.Name()) {
				return filepath.Join(path, file.Name()), nil
			}
		}
	}
	return "", errors.New("no job of " + job + "on path: " + path)

}

func SearchFile(path, suf string) (int, []string, error) {
	files, err := ioutil.ReadDir(path)
	if err != nil {
		log.Println(err)
		return 0, nil, err
	}
	var found []string
	for _, file := range files {
		if !file.IsDir() && file.Name()[0] != '.' {
			if strings.HasSuffix(file.Name(), suf) {
				found = append(found, filepath.Join(path, file.Name()))
				continue
			}
		}
	}
	length := len(found)
	if length > 0 {
		return length, found, nil
	}
	return 0, nil, errors.New("nofile of " + suf)
}

func FetchJobPath() (string, error) {
	tpath := filepath.Join(dfpath, month, today)
	ypath := filepath.Join(dfpath, month, yesterday)
	tStatus := Exists(tpath)
	yStatus := Exists(ypath)

	if tStatus {
		if jobpath, err := SearchJob(tpath); err == nil {
			return jobpath, nil
		}
	}

	if yStatus {
		if jobpath, err := SearchJob(ypath); err == nil {
			return jobpath, nil
		}
	}

	return "", errors.New("Today and Yesterday, can not be located the job: " + job)
}

func FetchPFpath() (string, error) {
	if rvst {return "",nil}
	jobpath, err := SearchJob(wks)
	if err != nil {
		jobpath, err = SearchJob(jxz)
	}
	if err != nil {
		return "", errors.New("PF job folder may not existed, please check!")
	}
	_, PFpath, err := SearchFile(jobpath, ".xls")
	if err != nil {
		_, PFpath, err = SearchFile(jobpath, ".xlsx")
	}
	if err != nil {
		return "", errors.New("PF sheet file is not located, please check!")
	}
	return PFpath[0], nil

}

func FetchTxtpath(jobpath string) (int, []string, error) {
	count, txtpath, err := SearchFile(jobpath, ".txt")
	if err != nil {
		return 0, nil, err
	}
	return count, txtpath, nil
}

func FetchTail(path string) (string, error) {
	if rvst {
		return "",nil
	}
	if strings.HasSuffix(path, ".xls") {
		return ParseXls(path)
	} else {
		return ParseXlsx(path)
	}
}

func CombineAll(allbody *AllBody, tail *string) *TxtJson {
	// emailtxt := ""
	// for _, v := range *allbody {
	// 	emailtxt = emailtxt + "\n" + Head(v.Count) + v.Content + "\n" + *tail + "\n"
	// }
	// return &emailtxt
	var txtjon TxtJson
	for _, v := range *allbody {
		txtjon.TxtBodies = append(txtjon.TxtBodies, TxtBody{TxtCount:v.Count, TxtBody: v.Content + "\n" + *tail })
	}
	return &txtjon
}


func PHQtitle(jobpath string) (string, error) {
	if rvst {
		return "", nil
	}
	buf := make([]byte, DEFAULTSIZE*1024)
	file, err := os.Open(jsonPath)
	if err != nil {
		log.Println("open json err")
		return "", err
	}
	defer file.Close()
	n, _ := file.Read(buf)
	if n == DEFAULTSIZE*1024 {
		return "", errors.New("MaxSize")
	}

	client := []map[string]string{}
	clientMap := map[string]string{}
	err = json.Unmarshal(buf[0:n], &client)
	if err != nil {
		return "", errors.New("jsonMarshalError")
	}

	for _, clientdata := range client {
		for code, fullname := range clientdata {
			clientMap[code] = fullname
		}
	}
	tbrand, tcode, tjob := TitleSplit(jobpath)

	return clientMap[tbrand] + " / " + clientMap[tcode] + " / " + tjob, nil
}
