package main

import (
    "fmt"
    "io/ioutil"
    "log"
    "net/http"
    "os"
)

const homepage string =
`<!DOCTYPE html>
<html>
<body>

<h1>COVID-19 Tracker</h1>

<form action="/">
    <h2>Location</h2>
    <p>Leave blank to get worldwide data</p>

    <label for="country">Country/Region:</label><br>
    <input type="text" id="country" name="country" value=""><br>

    <label for="province">Province/State:</label><br>
    <input type="text" id="province" name="province" value=""><br>

    <label for="city">City:</label><br>
    <input type="text" id="city" name="city" value=""><br>
    <br>

    <h2>Date</h2>
    <p>Leave blank to get latest date</p>

    <label for="month">Month (1-12):</label><br>
    <input type="text" id="month" name="month" value=""><br>

    <label for="day">Day:</label><br>
    <input type="text" id="day" name="day" value=""><br>

    <label for="year">Year:</label><br>
    <input type="text" id="year" name="year" value=""><br>
    <br>

    <input type="submit" value="Submit">
</form>
<br>

<h2>Results:</h2>

<pre>
%s
</pre>

</body>
</html>
`

func rootHandler(proxyPort string) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Request from:", r.RemoteAddr, r.URL.RequestURI())

        // Query database
        resp, err := http.Get("http://127.0.0.1:" + proxyPort +
            "/covid19-db:1.0/query?" + r.URL.RawQuery)
        if err != nil {
            fmt.Println(err)
            fmt.Fprint(w, err)
            return
        }
        defer resp.Body.Close()

        body, err := ioutil.ReadAll(resp.Body)
        if err != nil {
            fmt.Println(err)
            fmt.Fprint(w, err)
            return
        }

        fmt.Println("Response status:", resp.Status)
        fmt.Println(string(body))

        fmt.Fprintf(w, homepage, string(body))
    }
}

func faviconHandler(w http.ResponseWriter, r *http.Request) {
    http.ServeFile(w, r, "favicon.ico")
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage:", os.Args[0], "<proxy-port> <listen-port>")
        os.Exit(1)
    }
    proxyPort := os.Args[1]
    listenPort := os.Args[2]

    http.HandleFunc("/", rootHandler(proxyPort))
    http.HandleFunc("/favicon.ico", faviconHandler)
    log.Fatal(http.ListenAndServe(":" + listenPort, nil))
}
