package config

import (
    "bytes"
    "encoding/json"
    "io"
    "os"
)

type AuthInfoStruct struct {
    Use      bool
    Username string
    Password string
}

type ConfStruct struct {
    Host         string
    Port         string
    Queue_length int
    Middlewares  []string
    Auth         AuthInfoStruct
}

// 该辅助函数来自golang标准库io/ioutil/ioutil.go
func readAll(r io.Reader, capacity int64) (b []byte, err error) {
    buf := bytes.NewBuffer(make([]byte, 0, capacity))

    defer func() {
        e := recover()
        if e == nil {
            return
        }
        if panicErr, ok := e.(error); ok && panicErr == bytes.ErrTooLarge {
            err = panicErr
        } else {
            panic(e)
        }
    }()
    _, err = buf.ReadFrom(r)
    return buf.Bytes(), err
}

func readJson(f *os.File) (conf ConfStruct, err error) {
    var n int64
    if fi, err := f.Stat(); err == nil {
        if size := fi.Size(); size < 1e9 {
            n = size
        }
    }
    content, err := readAll(f, n+bytes.MinRead)
    err = json.Unmarshal(content, &conf)
    return
}

func ParseConf() (conf ConfStruct, err error) {
    f, err := os.OpenFile("./app.json", os.O_RDONLY, 0666)
    if err != nil {
        return conf, err
    }
    conf, err = readJson(f)
    return
}
