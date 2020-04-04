package main

import (
    "errors"
    "fmt"
    "log"
    "net/http"
    "os"
    "strconv"
)

func getParam(r *http.Request, param string) (string, bool) {
    vals, ok := r.URL.Query()[param]
    if !ok || len(vals[0]) < 1 {
        return "", false
    }
    return vals[0], true
}

func getNums(r *http.Request) (int, int, error) {
    num1str, ok := getParam(r, "num1")
    if !ok {
        return 0, 0, errors.New("calc-server: Could not find param 'num1'")
    }
    num1, err := strconv.Atoi(num1str)
    if err != nil {
        return 0, 0, err
    }

    num2str, ok := getParam(r, "num2")
    if !ok {
        return 0, 0, errors.New("calc-server: Could not find param 'num2'")
    }
    num2, err := strconv.Atoi(num2str)
    if err != nil {
        return 0, 0, err
    }

    return num1, num2, nil
}

func addHandler(w http.ResponseWriter, r *http.Request) {
    num1, num2, err := getNums(r)
    if err != nil {
        fmt.Println(err)
        return
    }

    sum := num1 + num2

    fmt.Println("Request from:", r.RemoteAddr)
    fmt.Printf("Add: %d + %d = %d\n\n", num1, num2, sum)
    fmt.Fprintf(w, "%d", sum)
}

func subHandler(w http.ResponseWriter, r *http.Request) {
    num1, num2, err := getNums(r)
    if err != nil {
        fmt.Println(err)
        return
    }

    diff := num1 - num2

    fmt.Println("Request from:", r.RemoteAddr)
    fmt.Printf("Sub: %d - %d = %d\n\n", num1, num2, diff)
    fmt.Fprintf(w, "%d", diff)
}

func mulHandler(w http.ResponseWriter, r *http.Request) {
    num1, num2, err := getNums(r)
    if err != nil {
        fmt.Println(err)
        return
    }

    prod := num1 * num2

    fmt.Println("Request from:", r.RemoteAddr)
    fmt.Printf("Mul: %d * %d = %d\n\n", num1, num2, prod)
    fmt.Fprintf(w, "%d", prod)
}

func divHandler(w http.ResponseWriter, r *http.Request) {
    num1, num2, err := getNums(r)
    if err != nil {
        fmt.Println(err)
        return
    }

    quot := num1 / num2

    fmt.Println("Request from:", r.RemoteAddr)
    fmt.Printf("Div: %d / %d = %d\n\n", num1, num2, quot)
    fmt.Fprintf(w, "%d", quot)
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "<listen-port>")
        os.Exit(1)
    }
    listenPort := os.Args[1]

    http.HandleFunc("/add", addHandler)
    http.HandleFunc("/sub", subHandler)
    http.HandleFunc("/mul", mulHandler)
    http.HandleFunc("/div", divHandler)
    log.Fatal(http.ListenAndServe(":" + listenPort, nil))
}
