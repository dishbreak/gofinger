package main

import (
    "fmt"
    "net"
    "os"
    "strings"
)

func main() {
    args := os.Args[1:]

    if len(args) != 1 {
        panic("Must provide username or username@host!")
    }

    parts := strings.Split(args[0], "@")

    var username string
    var host string

    switch len(parts) {
    case 1:
        username = parts[0]
        host = "127.0.0.1"
    case 2:
        username = parts[0]
        host = parts[1]
    default:
        panic("Must provide username or username@host!")

    }

    port := "4500"

    conn, err := net.Dial("tcp", host+":"+port)
    if err != nil {
        panic(err)
    }

    defer conn.Close()

    fmt.Fprintf(conn, username)
    buf := make([]byte, 1024)
    bytesRead, err := conn.Read(buf)
    if err != nil {
        panic(err)
    }

    fmt.Printf("%s", string(buf[0:bytesRead]))

}
