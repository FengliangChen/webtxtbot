package finalartwork

import (
  "encoding/json"
  "errors"
  "fmt"
  "runtime"
  "time"
)

type Broker [5]*JobInfo

func NewBroker() *Broker {
  var broker Broker
  go broker.CleanBySeconds(2)
  return &broker
}

type TxtReturn struct {
  Jobcode       string
  Token         string
  Body          string
  Zippable      bool
  Error         string
  TotalFileSize int64
}

type Track struct {
  Jobcode        string
  Token          string
  Error          string
  CompressedSize int64
  TotalFileSize  int64
}

type TrackReturn []Track

func FirstStageResponse(p *JobInfo, broker *Broker) (string, error) {
  var Response TxtReturn

  if p.FinalArtworkValidationCode == ZIPPABLE {
    Response.Jobcode = p.JobCode
    Response.Token = p.Token
    Response.Body = p.Body
    Response.Zippable = true
    Response.Error = p.Error
    Response.TotalFileSize = p.SizeCounter.Totalfilesize
    broker.Push(p)
  } else {
    Response.Jobcode = p.JobCode
    Response.Token = ""
    Response.Body = p.Body
    Response.Zippable = false
    Response.Error = p.Error
  }
  b, err := json.Marshal(&Response)
  if err != nil {
    return "", err
  }

  stringResponse := string(b)

  return stringResponse, nil
}

func (p *Broker) Push(job *JobInfo) error {

  for i, v := range p {
    if v != nil {
      p[i].BrokerStatusCode += 1
    }
  } // all BrokerStatusCode in broker add 1 every round

  for i, v := range p {
    if v != nil {
      if v.Token == job.Token {
        return errors.New(" token Push error")
      }
      if v.JobCode == job.JobCode {
        p[i] = job
        return nil
      } //same inquery.
    }
  }

  for i, v := range p {
    if v == nil {
      p[i] = job
      return nil
    }

  }

  for i, j := 0, 0; i < len(p); i++ {
    if p[i] != nil {
      if p[i].BrokerStatusCode > p[j].BrokerStatusCode {
        j = i
      }
    }
    if i == len(p)-1 {
      if p[j].CompressStatusCode != COMPRESSING {
        p[j] = job
      }
    }
  }
  return nil
}

func (p *Broker) Pop(token string) (*JobInfo, error) {

  for i, v := range p {
    if v != nil {
      if v.Token == token {
        return p[i], nil
      }
    }

  }
  return nil, errors.New("无效token")

}

func (p *Broker) Clear(token string) {
  for i, v := range p {
    if v != nil {
      if v.Token == token {
        p[i] = nil
      }
    }

  }

}

func (p *Broker) CleanBySeconds(s int) {
  tick := time.Tick(1 * time.Second)
  for {
    select {
    case <-tick:
      p.IncrementTimeCount()
      p.CleanCompletedJob(s)
    }
  }

}

func (p *Broker) IncrementTimeCount() {
  for i, v := range p {
    if v != nil {
      if p[i].CompressStatusCode == COMPRESSIONCOMPLETED {
        p[i].Timeflag += 1
      }
    }

  }
}

func (p *Broker) CleanCompletedJob(s int) {
  for i, v := range p {
    if v != nil {
      if p[i].CompressStatusCode == COMPRESSIONCOMPLETED {
        if p[i].Timeflag > s {
          p[i] = nil
          runtime.GC() // compromise to release the resource(zip file) immediately so that the file folder can be moved to other directory
        }

      }

    }
  }

}

func (p *Broker) TrackResponse() (string, error) {
  var trackReturn TrackReturn
  for i, v := range p {
    if v != nil {
      var track Track
      if p[i].CompressStatusCode == COMPRESSING {
        track.Jobcode = GetZipFileName(p[i].Jobpath)
        track.Token = p[i].Token
        track.Error = p[i].Error
        track.CompressedSize = p[i].SizeCounter.GetReadedSize()
        track.TotalFileSize = p[i].SizeCounter.Totalfilesize
      }
      if p[i].CompressStatusCode == COMPRESSIONCOMPLETED {
        track.Jobcode = GetZipFileName(p[i].Jobpath)
        track.Token = p[i].Token
        track.Error = p[i].Error
        track.CompressedSize = p[i].SizeCounter.Totalfilesize
        track.TotalFileSize = p[i].SizeCounter.Totalfilesize
        // p[i] = nil
      }
      if p[i] != nil {
        if p[i].CompressStatusCode == COMPRESSERROR {
          track.Jobcode = "COMPRESSERROR"
          track.Token = "COMPRESSERROR"
          track.Error = "COMPRESSERROR"
          track.CompressedSize = 0
          track.TotalFileSize = 0
          p[i] = nil
          continue
        }
      }

      if p[i] != nil {
        if p[i].CompressStatusCode == 0 {
          continue
        }
      }
      trackReturn = append(trackReturn, track)
    }
  }

  b, err := json.Marshal(&trackReturn)

  if err != nil {
    return "", err
  }

  stringResponse := string(b) // will be null
  return stringResponse, nil
}

func Run() {
  var broker Broker
  job := ProcessJob("2004J3")
  job2 := ProcessJob("2004K2")
  broker.Push(job)
  broker.Push(job2)
  stop := make(chan bool)
  go func() {
    testcount := 0
    for {
      testcount++
      select {
      case <-stop:
        return
      default:

        jstring, err := broker.TrackResponse()
        fmt.Println(jstring)
        if err != nil {
          fmt.Println(err)
        }
        time.Sleep(1000 * time.Millisecond)
      }
    }
  }()

  p1, err := broker.Pop(job.Token)

  if err != nil {
    fmt.Println(err)
  }
  p2, err := broker.Pop(job2.Token)
  if err != nil {
    fmt.Println(err)
  }

  done := make(chan bool)
  go func() {
    ProcessZip(p1)
    done <- true
  }()
  go func() {
    ProcessZip(p2)
    done <- true
  }()
  <-done
  <-done
  stop <- true
}
