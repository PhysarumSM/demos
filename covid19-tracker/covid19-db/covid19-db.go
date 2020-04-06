package main

// Based on data from https://github.com/CSSEGISandData/COVID-19

import (
    "encoding/csv"
    "encoding/json"
    "fmt"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "os"
    "strconv"
    "strings"
    "time"
)

type QueryResponse struct {
    Confirmed int
    Deaths int
    Recovered int
    Active int
}

const (
    cityCol int = 1
    provinceCol int = 2
    countryCol int = 3
    confirmedCol int = 7
    deathsCol int = 8
    recoveredCol int = 9
    activeCol int = 10
)

func datedReportName(month, day, year int) string {
    // MM-DD-YYYY.csv in UTC
    reportNameFormat := "%02d-%02d-%04d.csv"
    return fmt.Sprintf(reportNameFormat, month, day, year)
}

func latestReportName() string {
    year, month, day := time.Now().UTC().Date()
    return datedReportName(int(month), day-1, year)
}

func downloadReport(reportName string) error {
    urlBegin := "https://raw.githubusercontent.com/CSSEGISandData/COVID-19/master/csse_covid_19_data/csse_covid_19_daily_reports/"
    url := urlBegin + reportName
    
    resp, err := http.Get(url)
    if err != nil {
       return err
    }
    defer resp.Body.Close()

    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return err
    }

    ioutil.WriteFile(reportName, body, 0644)

    return nil
}

func getParam(r *http.Request, param string) (string, bool) {
    vals, ok := r.URL.Query()[param]
    if !ok || len(vals[0]) < 1 {
        return "", false
    }
    return vals[0], true
}

func getStrParam(r *http.Request, param string, defaultValue string) string {
    str, ok := getParam(r, param)
    if !ok {
        return defaultValue
    }
    return str
}

func getIntParam(r *http.Request, param string, defaultValue int) int {
    str, ok := getParam(r, param)
    if !ok {
        return defaultValue
    }
    integer, err := strconv.Atoi(str)
    if err != nil {
        return defaultValue
    }
    return integer
}

func queryParams(r *http.Request) (string, string, string, int, int, int) {
    city := getStrParam(r, "city", "")
    province := getStrParam(r, "province", "")
    country := getStrParam(r, "country", "")

    month := getIntParam(r, "month", 0)
    day := getIntParam(r, "day", 0)
    year := getIntParam(r, "year", 0)

    return city, province, country, month, day, year
}

func queryHandler(listenIp, listenPort string) func(http.ResponseWriter, *http.Request) {
    return func(w http.ResponseWriter, r *http.Request) {
        fmt.Println("Request from:", r.RemoteAddr, r.URL.RequestURI())

        city, province, country, month, day, year := queryParams(r)
        fmt.Println(city, province, country, month, day, year)

        var reportName string
        if month == 0 || day == 0 || year == 0 {
            reportName = latestReportName()
        } else {
            reportName = datedReportName(month, day, year)
        }

        _, err := os.Stat(reportName)
        if os.IsNotExist(err) {
            err2 := downloadReport(reportName)
            if err2 != nil {
                fmt.Println(err2)
                fmt.Fprint(w, err)
                return
            }
        } else if err != nil {
            fmt.Println(err)
            fmt.Fprint(w, err)
            return
        }

        reportFile, err := os.Open(reportName)
        if err != nil {
            fmt.Println(err)
            fmt.Fprint(w, err)
            return
        }

        var qr QueryResponse
        cr := csv.NewReader(reportFile)
        for {
            record, err := cr.Read()
            if err == io.EOF {
                break
            }
            if err != nil {
                fmt.Println(err)
                fmt.Fprint(w, err)
                return
            }

            if (strings.EqualFold(record[cityCol], city) || city == "") &&
                (strings.EqualFold(record[provinceCol], province) || province == "") &&
                (strings.EqualFold(record[countryCol], country) || country == "") {

                confirmed, err := strconv.Atoi(record[confirmedCol])
                if err != nil {
                    confirmed = 0
                }

                deaths, err := strconv.Atoi(record[deathsCol])
                if err != nil {
                    deaths = 0
                }

                recovered, err := strconv.Atoi(record[recoveredCol])
                if err != nil {
                    recovered = 0
                }

                active, err := strconv.Atoi(record[activeCol])
                if err != nil {
                    active = 0
                }

                qr.Confirmed += confirmed
                qr.Deaths += deaths
                qr.Recovered += recovered
                qr.Active += active
            }
        }

        reportFile.Close()

        data, err := json.Marshal(qr)
        if err != nil {
            fmt.Println(err)
            fmt.Fprint(w, err)
            return
        }

        dataStr := string(data)
        fmt.Println(dataStr)

        fmt.Fprintln(w, "Served by:", listenIp + ":" + listenPort)
        fmt.Fprintln(w, dataStr)
    }
}

func main() {
    if len(os.Args) != 3 {
        fmt.Println("Usage:", os.Args[0], "<listen-ip> <listen-port>")
        os.Exit(1)
    }
    listenIp := os.Args[1]
    listenPort := os.Args[2]

    http.HandleFunc("/query", queryHandler(listenIp, listenPort))
    log.Fatal(http.ListenAndServe(":" + listenPort, nil))
}
