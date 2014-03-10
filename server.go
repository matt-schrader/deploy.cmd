package main

// todos

import (
    "fmt"
    "net"
    "encoding/gob"
    "github.com/matt-schrader/deploy.cmd/model"
    "time"
)

var nextId int = 0

func handleConnection(node model.Node, wire net.Conn, queue chan *model.Work) {
    decoder := gob.NewDecoder(wire)
    encoder := gob.NewEncoder(wire)
    for work := range queue {
        err := encoder.Encode(work)
        if(err != nil) {
            fmt.Printf("Error occurred, will assume connectivity issue. %s", err.Error())
            queue <- work
            break
        }

        fmt.Printf("Sending the work to node %d\n", node.Id)
        for {
            workStatus := &model.WorkStatus{}
            fmt.Printf("Waiting for response\n")
            err := decoder.Decode(workStatus)
            if(err != nil) {
                fmt.Printf("Error occurred, will assume connectivity issue. %s", err.Error())
                queue <- work
                break
            }

            fmt.Printf("    Node %d responded with it's status\n", node.Id)
            if(workStatus.Done) {
                fmt.Printf("Node %d finished some work\n", node.Id)

                if(workStatus.Error != "") {
                    fmt.Printf("Error %s\n", workStatus.Error)
                }
                break
            }
        }
    }
    fmt.Printf("Connection closed %s\n", node.Id)
}

func assignWork(queue chan *model.Work, command string, args ...string) {
    queue <- &model.Work{Command: command, Args: args}
}


func main() {
    pending := make(chan *model.Work)

    go assignWork(pending, "cat", "/etc/hosts")
    go assignWork(pending, "grep", "-ir", "hello", "/home/matt/app/go/src/github.com/SchraderMJ11")

    go startServer(pending)
    
    var input string
    fmt.Scanln(&input)
    fmt.Printf("done\n")
}

func handshakeWithClient(node model.Node, wire net.Conn, queue chan *model.Work) {
    //encoder := gob.NewEncoder(wire)
    inTwoSeconds := time.Now().Add(2 * time.Second)
    fmt.Printf("Handshaking with node, must report for duty by %v\n", inTwoSeconds)
    wire.SetReadDeadline(inTwoSeconds)
    var tbuf [81920]byte
    n, err := wire.Read(tbuf[0:])
    fmt.Printf("Read bytes: %d\n", n)
    fmt.Printf("Read: %s\n", tbuf)
    if(err != nil) {
        fmt.Printf("Err: %s\n", err.Error())
        fmt.Printf("Node did not report for duty in time, terminating connection\n")
        return
    }
    var zero time.Time
    wire.SetReadDeadline(zero)
    go handleConnection(node, wire, queue)
}

func startServer(queue chan *model.Work) {
    port := "8080"
    fmt.Printf("start server %s\n", port);
    ln, err := net.Listen("tcp", ":" + port)
    if err != nil {
        // handle error
    }
    for {
        conn, err := ln.Accept() // this blocks until connection or error
        if err != nil {
            // handle error
            continue
        }
        newNode := model.Node{}
        nextId = nextId + 1
        newNode.Id = nextId
        newNode.Busy = false

        go handshakeWithClient(newNode, conn, queue)
    }
}

