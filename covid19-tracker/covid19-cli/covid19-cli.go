package main

import (
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "net/url"
    "os"
)

func main() {
    cityFlag := flag.String("city", "", "City")
    provinceFlag := flag.String("province", "", "Province or state")
    countryFlag := flag.String("country", "", "Country or region")
    monthFlag := flag.String("month", "", "Month (1-12)")
    dayFlag := flag.String("day", "", "Day (1-31)")
    yearFlag := flag.String("year", "", "Year (2020)")

    usage := func() {
        fmt.Fprintln(os.Stderr, "Usage:", os.Args[0], "[<options>] <proxy-port>")
        fmt.Fprintln(os.Stderr,
`
<proxy-port>
        Port to connect to proxy

<options>`)
        flag.PrintDefaults()
    }
    flag.Usage = usage

    flag.Parse()
    if len(flag.Args()) < 1 {
        fmt.Fprintln(os.Stderr, "Error: missing required arguments")
        usage()
        return
    }

    if len(flag.Args()) > 1 {
        fmt.Fprintln(os.Stderr, "Error: too many arguments")
        usage()
        return
    }

    proxyPort := flag.Arg(0)

    proxyUrl, err := url.Parse("http://127.0.0.1:" + proxyPort + "/covid19-db:1.0/query")
    if err != nil {
        panic(err)
    }
    q := proxyUrl.Query()
    q.Set("city", *cityFlag)
    q.Set("province", *provinceFlag)
    q.Set("country", *countryFlag)
    q.Set("month", *monthFlag)
    q.Set("day", *dayFlag)
    q.Set("year", *yearFlag)
    proxyUrl.RawQuery = q.Encode()

    fmt.Println(proxyUrl.String())

    // Make request
    resp, err := http.Get(proxyUrl.String())
    if err != nil {
        panic(err)
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        panic(err)
    }

    fmt.Println("Response status:", resp.Status)
    fmt.Println(string(body))
}
