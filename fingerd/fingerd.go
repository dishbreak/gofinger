package main

import (
    "bytes"
    "fmt"
    "html/template"
    "io/ioutil"
    "log"
    "net"
    "os"
    "os/user"
    "path"
    "strconv"
    "strings"
)

var templateWithoutPlan string = `
username: {{.Username}}  name: {{.Name}}
directory: {{.HomeDir}}
No plan file.
`

var templateWithPlan string = `
username: {{.Username}}  name: {{.Name}}
directory: {{.HomeDir}}
plan file: {{.PlanFileContents}}
`

type fingerUserRecord struct {
    Username         string
    Name             string
    HomeDir          string
    PlanFileContents string
}

func getUserRecord(username string) (string, error) {
    log.Printf("Getting user record for '%q'", username)
    record, err := user.Lookup(username)
    if err != nil {
        return "", err
    }

    planFilePath := path.Join(record.HomeDir, ".plan")
    /* Path looks arbitrary but isn't. Base directory comes from Golang API.
       Besides, nobody says finger is secure by any means. */
    /* #nosec */
    buffer, err := ioutil.ReadFile(planFilePath)

    planFileContents := ""
    templateText := templateWithoutPlan

    if err == nil {
        planFileContents = string(buffer)
        templateText = templateWithPlan
    }

    fingerRecord := fingerUserRecord{
        record.Username,
        record.Name,
        record.HomeDir,
        planFileContents,
    }

    templ, err := template.New("finger").Parse(templateText)
    if err != nil {
        panic(err)
    }

    var responseBuffer bytes.Buffer
    err = templ.Execute(&responseBuffer, fingerRecord)
    if err != nil {
        panic(err)
    }

    return strings.TrimSpace(responseBuffer.String()), nil
}

func fingerUser(conn net.Conn) {

    defer conn.Close()
    defer func() {
        if err := recover(); err != nil {
            log.Printf("Unexpected system error: %s", err)
            conn.Write([]byte("System error!"))
        }
    }()

    /* Usernames can't be longer than 32 bytes. */
    buf := make([]byte, 32)
    bytesRead, err := conn.Read(buf)
    if err != nil {
        panic(err)
    }
    log.Printf("recv: %d bytes on %s from remote host %s", bytesRead, conn.LocalAddr(), conn.RemoteAddr())

    username := strings.TrimSpace(string(buf[0:bytesRead]))
    result, err := getUserRecord(username)
    if err != nil {
        result = fmt.Sprintf("error: %s\n", err)
    }

    log.Printf("responding to query for %q", username)
    log.Printf("result: %s", result)
    conn.Write([]byte(result))
}

func main() {
    args := os.Args[1:]

    var port string
    switch len(args) {
    case 0:
        port = "4500"
    case 1:
        if _, err := strconv.Atoi(args[0]); err != nil {
            panic(fmt.Sprintf("Didn't recognize '%s' as a valid port number!", args[1]))
        }
        port = args[0]
    }

    listeningInterface := "127.0.0.1:" + port

    ld, err := net.Listen("tcp", listeningInterface)
    if err != nil {
        panic(fmt.Sprintf("Error opening socket: %s", err))
    }

    defer ld.Close()

    log.Printf("listening on %s", listeningInterface)

    for {
        conn, err := ld.Accept()
        if err != nil {
            log.Fatalf("Error encountered: %s", err)
        } else {
            log.Printf("Incoming connection.")
            go fingerUser(conn)
        }
    }
}
