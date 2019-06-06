package main

import (
    "flag"
    "fmt"
    "net"
    "strings"
)

func main() {
    port := flag.String("P", "4500", "Port to communicate with server.")
    flag.Parse()
    args := flag.Args()

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

    conn, err := net.Dial("tcp", host+":"+*port)
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

    fmt.Printf("%s\n", string(buf[0:bytesRead]))

}
