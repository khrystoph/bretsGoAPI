package main

import (
	"encoding/json"
	"errors"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

const (
	quoteBaseURL = "https://finnhub.io/api/v1/"
)

var (
	//Trace is log handling for Trace level messages
	Trace *log.Logger
	//Info is log handling for Info level messaging
	Info *log.Logger
	//Warning is log handling for Warning level messaging
	Warning *log.Logger
	//Error is log handling for Error level messaging
	Error                 *log.Logger
	traceHandle           io.Writer
	infoHandle            io.Writer = os.Stdout
	warningHandle         io.Writer = os.Stderr
	errorHandle           io.Writer = os.Stderr
	domain, finnhubAPIKey string
	ticker                = "amzn"
	testing               bool
)

type error interface {
	Error() string
}

type quote struct {
	StockTicker string  `json:"ticker,omitempty"`
	AskPrice    float64 `json:"price,omitempty"`
}

type quoteResponse struct {
	Error         *string `jsoin:"error,omitempty"`
	Close         float64 `json:"c,omitempty"`
	High          float64 `json:"h,omitempty"`
	Low           float64 `json:"l,omitempty"`
	Open          float64 `json:"o,omitempty"`
	PreviousClose float64 `json:"pc,omitempty"`
	Timestamp     int64   `json:"t,omitempty"`
}

type requestString struct {
	JSONRequest string `json:"ticker"`
}

func init() {
	Trace = log.New(traceHandle,
		"TRACE: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Info = log.New(infoHandle,
		"INFO: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Warning = log.New(warningHandle,
		"WARNING: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	Error = log.New(errorHandle,
		"ERROR: ",
		log.Ldate|log.Ltime|log.Lshortfile)

	flag.StringVar(&domain, "d", "example.com", "enter your fully qualified domain name here. Default: example.com")
	flag.StringVar(&domain, "domain", "example.com", "enter your fully qualified domain name here. Default: example.com")
	flag.BoolVar(&testing, "testing", false, "set this flag if you want to disable running on SSL/TLS and run in unprotected mode")
	flag.StringVar(&finnhubAPIKey, "apikey", "", "enter your finnhub API key to use ")
}

//getStockQuote will make a request out to the finnhub apis and return the close price for the day
func getStockQuote(symbol string) (stockQuote quote, err error) {
	stockQuote = quote{StockTicker: symbol, AskPrice: 0.00}
	getURL := quoteBaseURL + "quote?symbol=" + symbol + "&token=" + finnhubAPIKey
	Info.Printf("URL Request: %s\n", getURL)
	resp, err := http.Get(getURL)
	if err != nil {
		Error.Printf("Unable to retrieve ticker: %s", symbol)
		return
	}
	defer resp.Body.Close()
	decoder := json.NewDecoder(resp.Body)
	var stockQuoteResponse quoteResponse
	err = decoder.Decode(&stockQuoteResponse)
	if err != nil {
		return
	}
	if stockQuoteResponse.Error != nil {
		Info.Printf("Error set:%v  Error Message: %v", stockQuoteResponse.Error != nil, *stockQuoteResponse.Error)
		stockQuote = quote{}
		return stockQuote, errors.New(*stockQuoteResponse.Error)
	}

	Info.Printf("%v\n", stockQuoteResponse.Close)
	stockQuote.AskPrice = stockQuoteResponse.Close

	return stockQuote, nil
}

func quoteHandler(w http.ResponseWriter, req *http.Request) {
	Info.Printf("%s\n", req.Method)
	decoder := json.NewDecoder(req.Body)
	var stockString requestString
	err := decoder.Decode(&stockString)
	if err != nil {
		Error.Printf("failed to decode JSON. Panic!")
		return
	}
	if req.Method == "POST" && stockString.JSONRequest != "" {
		ticker = stockString.JSONRequest
	} else if stockString.JSONRequest == "" {
		Error.Printf("empty input JSON. Using Default stock ticker: AMZN")
	}
	ticker = strings.ToUpper(ticker)
	Info.Printf(ticker)
	stockQuote, err := getStockQuote(ticker)
	Info.Printf("stock Quote: %v", stockQuote)
	if err != nil {
		Info.Printf("failed to get stock quote: %v", err)
		errJSON := map[string]string{"error": err.Error()}
		errJSONString, err := json.Marshal(errJSON)
		if err != nil {
			Info.Printf("Failed to marshal JSON in error handling %v", err)
			return
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		w.Write(errJSONString)
	}
	stockQuoteString, err := json.Marshal(stockQuote)
	if err != nil {
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(stockQuoteString)
	return
}

func main() {
	flag.Parse()

	var (
		bretsGoAPIServer *http.Server
		helloHandler     = func(w http.ResponseWriter, _ *http.Request) {
			io.WriteString(w, "Welcome to Bret's API!\n")
		}
	)

	if domain == "example.com" {
		if domainenv := os.Getenv("DOMAIN"); domainenv != "" {
			domain = domainenv
		} else {
			Error.Printf("Domain value: %s\n$DOMAIN = %s", domain, domainenv)
			Error.Fatal("Please set the domain via domain flag or set DOMAIN env var.")
		}
	}

	http.HandleFunc("/", helloHandler)
	http.HandleFunc("/api/v1/quote", quoteHandler)

	if !testing {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
			Cache:      autocert.DirCache("certs"),
		}

		bretsGoAPIServer = &http.Server{
			Addr:      ":https",
			TLSConfig: certManager.TLSConfig(),
		}

		var wg sync.WaitGroup
		Info.Printf("Starting the letsencrypt server\n")
		go func() {
			wg.Add(1)
			defer wg.Done()
			http.ListenAndServe(":http", certManager.HTTPHandler(nil))
		}()
	} else {
		bretsGoAPIServer = &http.Server{
			Addr: ":8080",
		}
	}

	if !testing {
		Info.Printf("Starting the main TLS server.\n")
		Error.Fatal(bretsGoAPIServer.ListenAndServeTLS("", ""))
	} else {
		Error.Fatal(bretsGoAPIServer.ListenAndServe())
	}

	return
}
