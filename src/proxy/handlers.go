package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type UrlGetter = func(string) string

func getCacheUrl(endpoint string) string {
    return fmt.Sprintf("http://%s%s", os.Getenv("CUSTOMER_CACHE_SERVICE"), endpoint)
}

func getDBUrl(endpoint string) string {
    return fmt.Sprintf("http://%s%s", os.Getenv("CUSTOMER_DB_SERVICE"), endpoint)
}

type Request = func(UrlGetter, string) (interface{}, error)

func request(urlGetter UrlGetter, endpoint string) (interface{}, error) {
	resp, err := http.Get(urlGetter(endpoint))
	if err != nil {
		return nil, err
	}
	if resp.Body == nil {
		return nil, errors.New("request failed")
	}

	buffer := new(bytes.Buffer)
	_, err = io.Copy(buffer, resp.Body)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(buffer.Bytes(), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

type Executor = func(Request)

var sourceUrlGetter map[string]UrlGetter = map[string]UrlGetter{
	// "cache":    getCacheUrl,
	"database": getDBUrl,
}

func customerTransactions(ctx *gin.Context) {
	if ctx.Param("id") == "" {
		ctx.JSON(http.StatusOK, New("").AddErrors(errors.New("id not provided")))
		return
	}

	var wg sync.WaitGroup
	responses := []Response{}
	for source, getter := range sourceUrlGetter {
		wg.Add(1)

		go func(r Request) {
			defer wg.Done()

			responses = append(responses, New(source).SetData(func() (interface{}, error) {
				return r(getter, fmt.Sprintf("/customers/%s/transactions", ctx.Param("id")))
			}))
		}(request)
	}

	wg.Wait()
	ctx.JSON(http.StatusOK, responses)
}

func customerTotal(ctx *gin.Context) {
	if ctx.Param("id") == "" {
		ctx.JSON(http.StatusOK, New("").AddErrors(errors.New("id not provided")))
		return
	}

	var wg sync.WaitGroup
	responses := []Response{}
	for source, getter := range sourceUrlGetter {
		wg.Add(1)

		go func(r Request) {
			defer wg.Done()

			responses = append(responses, New(source).SetData(func() (interface{}, error) {
				return r(getter, fmt.Sprintf("/customers/%s/transactions/total", ctx.Param("id")))
			}))
		}(request)
	}

	wg.Wait()
	ctx.JSON(http.StatusOK, responses)
}

func sortedCustomers(ctx *gin.Context) {
	if ctx.Param("id") == "" {
		ctx.JSON(http.StatusOK, New("").AddErrors(errors.New("id not provided")))
		return
	}

	var wg sync.WaitGroup
	responses := []Response{}
	for source, getter := range sourceUrlGetter {
		wg.Add(1)

		go func(r Request) {
			defer wg.Done()

			responses = append(responses, New(source).SetData(func() (interface{}, error) {
				return r(getter, "/customers/sorted")
			}))
		}(request)
	}

	wg.Wait()
	ctx.JSON(http.StatusOK, responses)
}

const (
	CUSTOMER_MAX_ROWS = 1_000
)

// parseRequests parses the query and returns a map of request names to number of request
// names are: customer_transactions, customer_total, customer_sorted
func parseRequests(query string) map[string]uint {
	pairs := strings.Split(query, ",")
	result := map[string]uint{}
	for _, pair := range pairs {
		parts := strings.Split(pair, ":")
		if len(parts) < 2 {
			continue
		}

		val, err := strconv.ParseUint(parts[1], 10, 64)
		if err != nil {
			log.Println(err)
			continue
		}
		if parts[0] == "customer_sorted" {
			if val > 40 {
				val = 40
			}
		}

		result[parts[0]] = uint(val)
	}

	return result
}

func genEndpoint(name string) string {
	switch name {
	case "customer_transactions":
		return fmt.Sprintf("/customers/%d/transactions", 1+rand.Intn(CUSTOMER_MAX_ROWS))
	case "customer_total":
		return fmt.Sprintf("/customers/%d/transactions/total", 1+rand.Intn(CUSTOMER_MAX_ROWS))
	default:
		return fmt.Sprintf("/customers/sorted?page=%d&pageSize=10", 1+rand.Intn(CUSTOMER_MAX_ROWS/10))
	}
}

var ch chan SSEResponse

func sse(path string, response Response) {
    ch <- NewSSE(response.Source).SetData(path).SetDuration(response.Duration)
}

func genExecuters(source, query string, aggregate func(string, time.Duration)) map[string][]func() {
	requests := parseRequests(query)
	executors := map[string][]func(){}
	for name, number := range requests {
		for i := uint(0); i < number; i++ {
			executors[name] = append(executors[name], func() {
				endpoint := genEndpoint(name)
				resp := New(source).SetData(func() (interface{}, error) {
					return request(sourceUrlGetter[source], endpoint)
				})
				sse(endpoint, resp)
				aggregate(name, resp.Duration)
			})
		}
	}

	return executors
}

func AddHeaders(ctx *gin.Context) {
	ctx.Writer.Header().Set("Connection", "keep-alive")
	ctx.Writer.Header().Set("Cache-Control", "no-cache")
	ctx.Writer.Header().Set("Content-Type", "text/event-stream")
	ctx.Writer.Header().Set("Access-Control-Expose-Headers", "Content-Type")
}

type Aggregate = func(string, time.Duration)

func bulk(ctx *gin.Context) {
	query := ctx.Query("requests")
	AddHeaders(ctx)
    ch = make(chan SSEResponse)
	if query == "" {
		ctx.JSON(http.StatusBadRequest, NewSSE("").Final())
		return
	}
	requests := parseRequests(query)

	totals := map[string]map[string]time.Duration{}
	aggregates := map[string]Aggregate{}
	executors := map[string]map[string][]func(){}
	for source := range sourceUrlGetter {
		totals[source] = map[string]time.Duration{}
		for name := range requests {
			totals[source][name] = 0
		}
		aggregates[source] = func(s string, d time.Duration) {
			totals[source][s] += d
		}
		executors[source] = genExecuters(source, query, aggregates[source])
	}

    streamCh := make(chan bool)
    go func (ctx *gin.Context)  {
        streamCh <- ctx.Stream(func(w io.Writer) bool {
            resp := <-ch
            if !resp.Done {
                ctx.SSEvent("message", resp)
                return true
            }
            ctx.SSEvent("message", resp)
            close(ch)
            return false
        })
    }(ctx)

	var wg sync.WaitGroup
	for _, executorMap := range executors {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for _, executorArr := range executorMap {
				for _, executor := range executorArr {
					executor()
				}
			}
		}()
	}

	wg.Wait()
	averages := map[string]map[string]time.Duration{}
	for source, results := range totals {
		averages[source] = map[string]time.Duration{}
		for name, total := range results {
			averages[source][name] = total / time.Duration(requests[name])
		}
	}
    totalsStr := map[string]map[string]string{}
    averagesStr := map[string]map[string]string{}
    for src, values := range totals {
        totalsStr[src] = map[string]string{}
        averagesStr[src] = map[string]string{}
        for name, value := range values {
            totalsStr[src][name] = value.String()
            averagesStr[src][name] = averages[src][name].String()
        }
    }

	ch <- NewSSE("").SetData(map[string]interface{}{
		"totals":   totalsStr,
		"averages": averagesStr,
	}).Final()
    <-streamCh
}
