package finalartwork

import (
    "archive/zip"
    "context"
    "crypto/md5"
    "errors"
    "fmt"
    "io"
    "io/ioutil"
    "os"
    "os/exec"
    "path/filepath"
    "strconv"
    "strings"
    "sync/atomic"
    "time"
)

type CompressionStatusCode int

const (
    COMPRESSERROR        CompressionStatusCode = 1
    COMPRESSING          CompressionStatusCode = 2
    COMPRESSIONCOMPLETED CompressionStatusCode = 3
)

type FinalArtworkValidation int

const (
    ZIPPABLE            FinalArtworkValidation = 1
    BUILDBODYONLY       FinalArtworkValidation = 2
    FINALARTWORKFAILURE FinalArtworkValidation = 3
)

const (
    dfpath = "/Volumes/datavolumn_bmkserver_Pub"
)

type JobInfo struct {
    JobCode                    string
    Token                      string
    Body                       string
    Jobpath                    string
    Error                      string
    FinalArtworkValidationCode FinalArtworkValidation
    FinalArtworkFilePath       [2]string
    CompressStatusCode         CompressionStatusCode
    SizeCounter                *SizeInfo
    BrokerStatusCode           int
}

func GetToken() string {
    crutime := time.Now().Unix()
    h := md5.New()
    io.WriteString(h, strconv.FormatInt(crutime, 10))
    token := fmt.Sprintf("%x", h.Sum(nil))
    return token
}

type PossibleJobpathRange struct {
    today     string
    yesterday string
}

type SizeInfo struct {
    Count            int64
    Accumulativation int64
    Totalfilesize    int64
}

func (p SizeInfo) GetReadedSize() int64 {
    return p.Accumulativation + p.Count
}

func (p SizeInfo) GetProcessPercentage() float64 {
    return float64(p.Accumulativation+p.Count) / float64(p.Totalfilesize)
}

func (p *JobInfo) GetTotalFileSize() error {
    var totalSize int64 = 0
    for _, v := range p.FinalArtworkFilePath {
        err := filepath.Walk(v, func(filePath string, info os.FileInfo, err error) error {

            if info.IsDir() {
                return nil
            }
            if info.Name()[0] == '.' {
                return nil
            }
            if err != nil {
                return err
            }
            totalSize += info.Size()
            return nil
        })

        if err != nil {
            return err
        }
    }
    p.SizeCounter = &SizeInfo{Count: 0, Totalfilesize: totalSize}
    return nil
}

type Reader struct {
    r *os.File
    n int64
}

func NewReader(r *os.File) *Reader {
    return &Reader{
        r: r,
    }
}

func (r *Reader) Read(p []byte) (n int, err error) {
    n, err = r.r.Read(p)
    atomic.AddInt64(&r.n, int64(n))
    return
}

func (r *Reader) N() int64 {
    return atomic.LoadInt64(&r.n)
}

func (p *PossibleJobpathRange) MakePath() {
    now := time.Now()
    today := now.Format("0102")
    yesterday := now.AddDate(0, 0, -1).Format("0102")
    month := now.Format("200601")
    p.today = filepath.Join(dfpath, month, today)
    p.yesterday = filepath.Join(dfpath, month, yesterday)
}

func GetToday() string {
    now := time.Now()
    today := now.Format("0102")
    return today
}

func GetZipFileName(jobpath string) string {
    baseName := filepath.Base(jobpath)
    baseName = strings.TrimSpace(baseName)
    time := GetToday()
    return baseName + "_" + time + ".zip"
}

func IfDirExist(path string) bool {
    _, err := os.Stat(path)
    if err != nil {
        if os.IsExist(err) {
            return true
        }
        return false
    }
    return true
}

func (p PossibleJobpathRange) SearchPath(job string) (string, error) {

    if IfDirExist(dfpath) == false {
        return "", errors.New("服务器是否连接")
    }

    plist := [2]string{p.today, p.yesterday}

    for _, path := range plist {

        if !IfDirExist(path) {
            continue
        }
        jobPath, err := SearchDir(path, job)
        if err == nil && jobPath != "" {
            return jobPath, nil
        }

    }

    return "", errors.New("Can't locate the Job today and yesterday.")
}

func SearchDir(path, job string) (string, error) {
    files, err := ioutil.ReadDir(path)
    if err != nil {
        return "", err
    }
    for _, file := range files {
        if file.IsDir() {
            if strings.Contains(file.Name(), job) {
                return filepath.Join(path, file.Name()), nil
            }
        }
    }
    return "", errors.New("no path of " + job + "on path: " + path)
}

func (p *JobInfo) GetJobPath() error {
    var pathRange PossibleJobpathRange
    pathRange.MakePath()
    jobPath, err := pathRange.SearchPath(p.JobCode)

    if err != nil {
        return err
    } else {
        p.Jobpath = jobPath
        return nil
    }
}

func (p *JobInfo) FinalArtworkValidation() error {
    // var AI_ThisFolderToPrinter, PDF_Locked_For_Visual_Ref string
    AI_ThisFolderToPrinter, err := SearchDir(p.Jobpath, "AI_ThisFolderToPrinter")
    if err != nil {
        p.FinalArtworkValidationCode = FINALARTWORKFAILURE
        return err
    }
    p.FinalArtworkFilePath[0] = AI_ThisFolderToPrinter

    PDF_Locked_For_Visual_Ref, err := SearchDir(p.Jobpath, "PDF_Locked_For_Visual_Ref")
    if err != nil {
        p.FinalArtworkValidationCode = BUILDBODYONLY
        return err
    }
    p.FinalArtworkFilePath[1] = PDF_Locked_For_Visual_Ref

    p.FinalArtworkValidationCode = ZIPPABLE

    return nil
}

func (p *JobInfo) BuildBody() error {
    // var body []string
    var body string
    // Walk function,walk by order of names?
    err := filepath.Walk(p.FinalArtworkFilePath[0], func(filePath string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }
        if info.IsDir() {
            return nil
        }
        if info.Name()[0] == '.' {
            return nil
        }

        if strings.HasSuffix(info.Name(), ".ai") {
            // body = append(body, strings.TrimSuffix(info.Name(), ".ai"))
            body = body + strings.TrimSuffix(info.Name(), ".ai") + "\n"
        }
        return nil

    })
    if err != nil {
        return err
    }
    p.Body = body
    return nil
}

func (p *JobInfo) RecursiveZip() error {
    p.CompressStatusCode = COMPRESSING

    firstLayer := p.Jobpath
    pathToZip := p.FinalArtworkFilePath
    destinationPath := filepath.Join(os.Getenv("HOME"), "Desktop", GetZipFileName(p.Jobpath))

    destinationFile, err := os.Create(destinationPath)
    if err != nil {
        return err
    }
    myZip := zip.NewWriter(destinationFile)

    for _, v := range pathToZip {
        err = filepath.Walk(v, func(filePath string, info os.FileInfo, err error) error {

            if info.IsDir() {
                return nil
            }
            if err != nil {
                return err
            }
            if info.Name()[0] == '.' {
                return nil
            }
            relPath := strings.TrimPrefix(filePath, filepath.Dir(firstLayer)+"/") // With "/" at front of path, will cause empty folder in Windows 10 built in compressor.
            header, err := zip.FileInfoHeader(info)
            if err != nil {
                return err
            }

            header.Name = relPath //keep the zip folder structure by relative path
            header.Method = zip.Deflate
            zipFile, err := myZip.CreateHeader(header)
            if err != nil {
                return err
            }

            fsFile, err := os.Open(filePath)
            if err != nil {
                p.CompressStatusCode = COMPRESSERROR
                return err
            }

            newreader := NewReader(fsFile)
            ctx, cancel := context.WithCancel(context.Background())
            go func(ctx context.Context) {
                for {
                    select {
                    case <-ctx.Done():
                        return
                    default:
                        p.SizeCounter.Count = newreader.N()
                        time.Sleep(100 * time.Millisecond)
                    }
                }
            }(ctx)

            _, err = io.Copy(zipFile, newreader)
            cancel()
            if err != nil {
                p.CompressStatusCode = COMPRESSERROR
                return err
            }
            p.SizeCounter.Accumulativation += newreader.N()
            p.SizeCounter.Count = 0
            return nil
        })
        if err != nil {
            p.CompressStatusCode = COMPRESSERROR
            return err
        }
    }

    err = myZip.Close()
    if err != nil {
        p.CompressStatusCode = COMPRESSERROR
        return err
    }
    p.CompressStatusCode = COMPRESSIONCOMPLETED
    return nil
}

func ProcessJob(queryCode string) *JobInfo {
    var job JobInfo

    job.JobCode = strings.ToUpper(queryCode)
    job.Token = GetToken()
    err := job.GetJobPath() // mostly can't get the job path.
    if err != nil {
        job.Error = err.Error()
        return &job
    }

    err = job.FinalArtworkValidation()
    if err != nil {
        job.Error = err.Error()
        return &job
    }

    if job.FinalArtworkValidationCode == BUILDBODYONLY {
        err = job.BuildBody()

        if err != nil {
            job.Error = err.Error()
            return &job
        }

    }

    if job.FinalArtworkValidationCode == ZIPPABLE {
        err = job.BuildBody()
        if err != nil {
            job.Error = err.Error()
            return &job
        }

        err = job.GetTotalFileSize()
        if err != nil {
            job.Error = err.Error()
            return &job
        }
    }
    return &job
}

func ProcessZip(job *JobInfo) {
    err := job.RecursiveZip()
    if err != nil {
        job.Error = err.Error()
    }

}

func (p JobInfo) OpenFolder() error {
    cmd := exec.Command("open", p.Jobpath)
    err := cmd.Run()
    if err != nil {
        return err
    }
    return nil
}
