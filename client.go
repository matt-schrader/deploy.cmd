package main

// todos

import (
    "fmt"
    "flag"
    "log"
    "net"
    "encoding/gob"
    "os/exec"
    "github.com/matt-schrader/deploy.cmd/model"
)

var host string
func init() {
    flag.StringVar(&host, "host", "localhost:8080", "Location of the server.  Format is <host>:<port>.")
}

func main() {
    flag.Parse()
    go startClient()

    var input string
    fmt.Scanln(&input)
    fmt.Printf("done\n")
}

func startClient() {
    fmt.Printf("start client connecting to %s\n", host);
    conn, err := net.Dial("tcp", host)
    if err != nil {
        log.Fatal("Connection error", err)
    }

    handleClient(conn)
}

func handleClient(conn net.Conn) {
    decoder := gob.NewDecoder(conn)
    encoder := gob.NewEncoder(conn)
    for {
        iwork := &model.Work{}
        err := decoder.Decode(iwork)
        if(err != nil) {
            fmt.Printf("Connection seems to have been closed.  %s", err.Error())
            break
        }

        cmd := exec.Command(iwork.Command, iwork.Args...)
        out, err := cmd.Output()

        workStatus := &model.WorkStatus{Done: true, Results: out}
        if(err != nil) {
            workStatus.Error = err.Error()
        }
        encoder.Encode(workStatus)
    }
}

