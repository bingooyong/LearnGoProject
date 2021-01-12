package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"
)

func main() {
	http.HandleFunc("/", handler)
	if err := http.ListenAndServe(":18090", nil); err != nil {
		panic(err)
	}
}

func handler(w http.ResponseWriter, r *http.Request) {
	requestMap := createRequestMap(r)
	marshal, _ := json.Marshal(requestMap)
	w.Write(marshal)
}

func createRequestMap(r *http.Request) map[string]interface{} {
	m := map[string]interface{}{
		"timeGo":     TimeFmt(time.Now()),
		"proto":      r.Proto,
		"host":       r.Host,
		"requestUri": r.RequestURI,
		"remoteAddr": r.RemoteAddr,
		"method":     r.Method,
		"url":        r.URL.String(),
		"headers":    ConvertHeader(r.Header),
	}
	m["timeTo"] = TimeFmt(time.Now())
	return m
}

// TimeFmt format time.
func TimeFmt(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.0000")
}

// ConvertHeader convert s head map[string][]string to map[string]string.
func ConvertHeader(query map[string][]string) map[string]string {
	q := make(map[string]string)
	for k, v := range query {
		q[k] = strings.Join(v, " ")
	}

	return q
}
