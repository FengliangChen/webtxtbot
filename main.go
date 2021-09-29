package main

import (
    "bytes"
    "errors"
    "finalartwork"
    "fmt"
    "html/template"
    "io"
    "net/http"
    "runtime"
    "searchcode"
    "txtbot"
)

// 端口
const (
    HTTP_PORT  string = "8989"
    HTTPS_PORT string = "8990"
)

// 目录
const (
    CSS_CLIENT_PATH   = "/css/"
    DART_CLIENT_PATH  = "/js/"
    IMAGE_CLIENT_PATH = "/image/"

    CSS_SVR_PATH   = "web"
    DART_SVR_PATH  = "web"
    IMAGE_SVR_PATH = "web"
)

func init() {
    runtime.GOMAXPROCS(runtime.NumCPU())
}

var (
    broker = finalartwork.NewBroker()
)

func main() {
    // 先把css和脚本服务上去
    http.Handle(CSS_CLIENT_PATH, http.FileServer(http.Dir(CSS_SVR_PATH)))
    http.Handle(DART_CLIENT_PATH, http.FileServer(http.Dir(DART_SVR_PATH)))
    http.Handle(IMAGE_CLIENT_PATH, http.FileServer(http.Dir(IMAGE_SVR_PATH)))

    // 网址与处理逻辑对应起来
    http.HandleFunc("/", HomePage)
    http.HandleFunc("/draftArtworkText", DraftArtworkText)

    http.HandleFunc("/final", FinalHomePage)
    http.HandleFunc("/query", Query)
    http.HandleFunc("/compress", Compress)
    http.HandleFunc("/track", Track)
    http.HandleFunc("/record", Record)
    http.HandleFunc("/autoemail", AutoEmail)

    http.HandleFunc("/searchcode", SearchHomePage)
    http.HandleFunc("/updatesupplier", UpdateSupplier)
    http.HandleFunc("/suppliersearch", SupplierSearch)

    // 开始服务
    err := http.ListenAndServe(":"+HTTP_PORT, nil)
    if err != nil {
        fmt.Println("服务失败 /// ", err)
    }
}

func WriteTemplateToHttpResponse(res http.ResponseWriter, t *template.Template) error {
    if t == nil || res == nil {
        return errors.New("WriteTemplateToHttpResponse: t must not be nil.")
    }
    var buf bytes.Buffer
    err := t.Execute(&buf, nil)
    if err != nil {
        return err
    }
    res.Header().Set("Content-Type", "text/html; charset=utf-8")
    _, err = res.Write(buf.Bytes())
    return err
}

func HomePage(res http.ResponseWriter, req *http.Request) {
    t, err := template.ParseFiles("web/txtb.html")
    if err != nil {
        fmt.Println(err)
        return
    }
    err = WriteTemplateToHttpResponse(res, t)
    if err != nil {
        fmt.Println(err)
        return
    }
    return
}

func DraftArtworkText(res http.ResponseWriter, req *http.Request) {
    keys, ok := req.URL.Query()["j"]
    if !ok || len(keys[0]) < 1 {
        fmt.Println("Url Param 'key' is missing")
        return
    }
    key := keys[0]

    rvst := false
    rvsts, ok := req.URL.Query()["rvst"]
    if !ok || len(rvsts[0]) < 1 {
        rvst = false
    } else {
        rvst = true
    }

    job := string(key)
    finaltxt := txtbot.Run(job, rvst)
    io.WriteString(res, finaltxt)
    return
}

func FinalHomePage(res http.ResponseWriter, req *http.Request) {
    t, err := template.ParseFiles("web/finalartwork.html")
    if err != nil {
        fmt.Println(err)
        return
    }
    err = WriteTemplateToHttpResponse(res, t)
    if err != nil {
        fmt.Println(err)
        return
    }
    return
}

func Query(res http.ResponseWriter, req *http.Request) {
    keys, ok := req.URL.Query()["j"]
    if !ok || len(keys[0]) < 1 {
        fmt.Println("Url Param 'key' is missing")
        return
    }
    key := keys[0]

    job := finalartwork.ProcessJob(string(key))
    jstring, err := finalartwork.FirstStageResponse(job, &broker)
    if err != nil {
        fmt.Println(err)
    }
    io.WriteString(res, jstring)
    if len(job.Jobpath) != 0 {
        job.OpenFolder()
    }
    return
}

func Compress(res http.ResponseWriter, req *http.Request) {

    keys, ok := req.URL.Query()["j"]
    if !ok || len(keys[0]) < 1 {
        fmt.Println("Url Param 'key' is missing")
        return
    }
    key := keys[0]

    job, err := broker.Pop(string(key))
    if err != nil {
        fmt.Println(err)
        io.WriteString(res, err.Error())
        return
    }

    if job.CompressStatusCode != 0 {
        io.WriteString(res, "已经提交压缩过了。")
        return
    }
    go func() {
        finalartwork.ProcessZip(job)
    }()
    io.WriteString(res, "压缩中。。。")
    return
}

func Track(res http.ResponseWriter, req *http.Request) {
    jstring, err := broker.TrackResponse()
    if err != nil {
        fmt.Println(err)
    }
    io.WriteString(res, jstring)
    return
}

func SearchHomePage(res http.ResponseWriter, req *http.Request) {
    t, err := template.ParseFiles("web/searchcode.html")
    if err != nil {
        fmt.Println(err)
        return
    }
    err = WriteTemplateToHttpResponse(res, t)
    if err != nil {
        fmt.Println(err)
        return
    }
    return
}

func UpdateSupplier(res http.ResponseWriter, req *http.Request) {
    var sup, supCode, responseTxt string

    key1, ok1 := req.URL.Query()["c"]
    key2, ok2 := req.URL.Query()["v"]
    if !ok1 || len(key1[0]) < 1 {
        return
    }
    supCode = key1[0]

    if !ok2 || len(key2[0]) < 1 {
        searchcode.JsonDel(supCode)

        err := searchcode.JsonSave()
        if err != nil {
            responseTxt = err.Error()
        } else {
            responseTxt = "已删除" + supCode
        }

    } else {
        sup = key2[0]

        searchcode.JsonAdd(supCode, sup)
        err := searchcode.JsonSave()
        if err != nil {
            responseTxt = err.Error()
        } else {
            responseTxt = supCode + "已经更新为" + sup
        }

    }
    io.WriteString(res, responseTxt)
    return
}

func SupplierSearch(res http.ResponseWriter, req *http.Request) {
    keys, ok := req.URL.Query()["c"]
    if !ok || len(keys[0]) < 1 {
        return
    }
    key := keys[0]

    var result string

    if finalartwork.IfDirExist("/Volumes/datavolumn_bmkserver_Pub") == false {
        io.WriteString(res, "服务器是否连接。")
        return
    }

    err := searchcode.Mapcheck()
    if err != nil {
        io.WriteString(res, err.Error())
        return
    }
    result = searchcode.JsonSearch(key)
    io.WriteString(res, result)
    return
}

func Record(res http.ResponseWriter, req *http.Request) {
    keys, ok := req.URL.Query()["j"]
    if !ok || len(keys[0]) < 1 {
        fmt.Println("Url Param 'key' is missing")
        return
    }
    key := keys[0]

    jstring := txtbot.ToFaris(string(key))
    // if err != nil {
    //     fmt.Println(err)
    // }
    io.WriteString(res, jstring)
    return
}

func AutoEmail(res http.ResponseWriter, req *http.Request) {
    if err := req.ParseForm(); err != nil {
        fmt.Fprintf(res, "ParseForm err: %v", err)
        return
    }

    title := req.FormValue("title")
    content := req.FormValue("content")
    finalartwork.MakeEmail(title, content)
    return

}
