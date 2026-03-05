package main

import (
	"embed"
	"errors"
	"fmt"
	"net/http"
	"slices"
	"sync"
	"time"

	"github.com/a-h/templ"

	"experiment/htmlstream/templates"
)

//go:embed static
var staticFiles embed.FS

func main() {
	http.Handle("/static/", http.FileServer(http.FS(staticFiles)))

	http.HandleFunc("/async", concurrentStreamingPage)
	http.HandleFunc("/sync", basicStreamingPage)

	fmt.Println("Listening on :3000")
	http.ListenAndServe(":3000", nil)
}

// Values are computed concurrently, but are not rendered concurrently.
// This means that waiting for a value to be ready to render is also blocking subsequent renders;
// however, during this time other values can become ready for later immediate render.
func concurrentStreamingPage(w http.ResponseWriter, r *http.Request) {
	myData1 := templates.Concurrent(func() string {
		time.Sleep(3_500 * time.Millisecond)
		return "My data 1"
	})
	myData2 := templates.Concurrent(func() string {
		time.Sleep(3_500 * time.Millisecond)
		return "My data 2"
	})

	myData3 := templates.Concurrent(func() string {
		time.Sleep(1_500 * time.Millisecond)
		return "My data 3"
	})
	myError := templates.TryConcurrent(func() (string, error) {
		time.Sleep(1_500 * time.Millisecond)
		return "", errors.New("Deliberate failure")
	})

	mySeq := templates.ConcurrentSeq(
		func() string {
			time.Sleep(1_500 * time.Millisecond)
			return "My data seq 1"
		},
		func() string {
			time.Sleep(1_000 * time.Millisecond)
			return "My data seq 2"
		},
	)

	page := templates.StreamingPage(myData1, myData2, myData3, mySeq, myError)
	templ.Handler(page, templ.WithStreaming()).ServeHTTP(w, r)
}

// This is included purely for comparison purposes.
// It doesn't spawn any goroutines, it just makes evaluation lazy.
func basicStreamingPage(w http.ResponseWriter, r *http.Request) {
	myData1 := sync.OnceValue(func() string {
		time.Sleep(2_500 * time.Millisecond)
		return "My data 1"
	})
	myData2 := sync.OnceValue(func() string {
		time.Sleep(2_500 * time.Millisecond)
		return "My data 2"
	})

	myData3 := sync.OnceValue(func() string {
		time.Sleep(1_500 * time.Millisecond)
		return "My data 3"
	})
	myError := sync.OnceValues(func() (string, error) {
		time.Sleep(1_500 * time.Millisecond)
		return "", errors.New("Deliberate failure")
	})

	mySeq := slices.Values([]func() string{
		sync.OnceValue(func() string {
			time.Sleep(1_500 * time.Millisecond)
			return "My data seq 1"
		}),
		sync.OnceValue(func() string {
			time.Sleep(1_000 * time.Millisecond)
			return "My data seq 2"
		}),
	})

	page := templates.StreamingPage(myData1, myData2, myData3, mySeq, myError)
	templ.Handler(page, templ.WithStreaming()).ServeHTTP(w, r)
}
