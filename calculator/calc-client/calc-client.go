package main

import (
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
    "math/rand"
    "strconv"
    "time"
)

func constructParams(cmd, num1, num2 string) string {
    return cmd + "?num1=" + num1 + "&num2=" + num2
}

func main() {
    if len(os.Args) != 2 && len(os.Args) != 5 {
        fmt.Println("Usage:", os.Args[0],
            "[<proxy-port>]|[<addr> <cmd> <num1> <num2>]")
        os.Exit(1)
    }

    if len(os.Args) == 5 {
        addr := os.Args[1]
        cmd := os.Args[2]
        num1 := os.Args[3]
        num2 := os.Args[4]

        resp, err := http.Get(
            "http://" + addr + "/" + constructParams(cmd, num1, num2))
        if err != nil {
           panic(err)
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }

        fmt.Println("Result:", string(body))
        return
    }

    proxyPort := os.Args[1]

    rand.Seed(time.Now().UnixNano())
    commands := []string{"add", "sub", "mul", "div"}

    // Do every 8 seconds
    pollInterval := 8
    timerCh := time.Tick(time.Duration(pollInterval) * time.Second)
    for range timerCh {
        // Generate random command
        cmd := commands[rand.Intn(len(commands))]
        num1 := strconv.Itoa(rand.Intn(10)+1)
        num2 := strconv.Itoa(rand.Intn(10)+1)
        fmt.Println("Command:", cmd, num1, num2)

        // Make request
        resp, err := http.Get("http://127.0.0.1:" + proxyPort +
            "/calc-server:1.0/" + constructParams(cmd, num1, num2))
        if err != nil {
            fmt.Println(err.Error())
            continue
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            panic(err)
        }

        fmt.Println("Response status:", resp.Status)
        fmt.Println("Result:", string(body))
    }
}
