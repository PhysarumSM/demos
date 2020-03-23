package main

import (
    "fmt"
    "log"
    "net/http"
    "os"
)

func handler(w http.ResponseWriter, r *http.Request) {
    fmt.Println("Got request")
    fmt.Fprintf(w, "Hello, World\n")
}

func main() {
    if len(os.Args) != 2 {
        fmt.Println("Usage:", os.Args[0], "PORT")
        os.Exit(1)
    }

    port := os.Args[1]
    http.HandleFunc("/", handler)
    log.Fatal(http.ListenAndServe(":" + port, nil))
}
