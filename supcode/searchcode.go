package searchcode

import(
"fmt"
"strings"
"encoding/json"
"errors"
"os"
// "time"
)


var (
	DEFAULTSIZE = 60
	ClientMap = map[string]string{}
	JasonStat os.FileInfo
	jsonPath = "/Volumes/datavolumn_bmkserver_Pub/新做稿/已结束/NON-WMT/.database/clientcode.json"
)


func ReadJson(jpath string) (error) {
	Client := []map[string]string{}
	buf := make([]byte, DEFAULTSIZE*1024)
	file, err := os.Open(jpath)
	if err != nil {
		fmt.Println("open json err")
		return err
	}
	defer file.Close()
	n, _ := file.Read(buf)
	if n == DEFAULTSIZE*1024 {
		return errors.New("MaxSize")
	}
	err = json.Unmarshal(buf[0:n], &Client)
	if err != nil {
		return err
	}
	for _, clientdata := range Client {
		for code, fullname := range clientdata {
			ClientMap[code] = fullname
		}
	}
	return nil
}

func JsonSearch(name string)string{
	var resultText string
	for code,value := range ClientMap {
		if strings.Contains(code, name) || strings.Contains(value, name){
			result := "{"+"\""+ code + "\"" + ":" + "\"" + value + "\"" + "},"
			resultText += result
		}
	}

	if len(resultText) != 0 {
		return "[" + resultText[0:len(resultText)-1] + "]"
	}
	return resultText
}

func JsonSave()error {
	var mapList = []map[string]string{}

	for code, vaule := range ClientMap {
		tempMap := map[string]string{}
		tempMap[code] = vaule
		mapList = append(mapList, tempMap)
	}
	b, err := json.Marshal(&mapList)
	if err != nil {
		return err
	}

	file, err := os.Create(jsonPath)
	if err != nil {
		return err
	}
	_, err = file.Write(b)
	if err != nil {
		return err
	}
	defer file.Close()
	return nil

}

func JsonAdd(code, value string){
	ClientMap[code] = value
}

func JsonDel(code string){
	if _, ok := ClientMap[code]; ok{
		delete(ClientMap, code)
	}
}

func Mapcheck()error{
	if len(ClientMap) == 0 {
		var err error
		JasonStat, err = os.Stat(jsonPath)
	    if err != nil {
	        return err
	    }
		err = ReadJson(jsonPath)
		if err != nil {
			return err
		}
		return nil
	}

	stat, err := os.Stat(jsonPath)
    if err != nil {
        return err
    }
    if stat.Size() != JasonStat.Size() || stat.ModTime() != JasonStat.ModTime() {
    	err := ReadJson(jsonPath)
    	if err != nil {
        	return err
    	}
    	JasonStat = stat
    }
	return nil
}

// func watchFile(filePath string) error {
//     initialStat, err := os.Stat(filePath)
//     if err != nil {
//         return err
//     }

//     stat, err := os.Stat(filePath)
//     if err != nil {
//         return err
//     }

//     if stat.Size() != initialStat.Size() || stat.ModTime() != initialStat.ModTime() {

//     }

//     return nil
// }


// func main(){
// 	path := "/Users/bmk/Desktop/clientcode.json"
// 	ReadJson(path)

// 	JsonDel("SUM")

// 	// result := JsonSearch("h")

// 	// fmt.Println(result)
// 	err := JsonSave("/Users/bmk/Desktop/NEWclientcode.json")
// 	if err != nil {
// 		fmt.Println(err)
// 	}

// }
