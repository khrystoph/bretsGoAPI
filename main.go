package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"io"
	"log"
	"net/http"
	"os"
	"sync"

	"golang.org/x/crypto/acme/autocert"
)

const (
	quoteBaseURL = "http://finance.google.com/finance/"
)

var (
	//Trace is log handling for Trace level messages
	Trace *log.Logger
	//Info is log handling for Info level messaging
	Info *log.Logger
	//Warning is log handling for Warning level messaging
	Warning *log.Logger
	//Error is log handling for Error level messaging
	Error         *log.Logger
	traceHandle   io.Writer
	infoHandle    io.Writer = os.Stdout
	warningHandle io.Writer = os.Stderr
	errorHandle   io.Writer = os.Stderr
	domain        string
	ticker        = "amzn"
	testing       bool
)

type quote struct {
	StockTicker string  `json:"ticker,omitempty"`
	AskPrice    float64 `json:"price"`
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
}

//getGoogleQuote will make a request out to the google finance apis and return the close price for the day
func getGoogleQuote(symbol string) (stockQuote quote, err error) {
	stockQuote = quote{StockTicker: symbol, AskPrice: 0.00}
	return stockQuote, nil
}

func main() {
	flag.Parse()

	var (
		bretsGoAPIServer *http.Server
		helloHandler     = func(w http.ResponseWriter, _ *http.Request) {
			io.WriteString(w, "Welcome to Bret's API!\n")
		}
		quoteHandler = func(w http.ResponseWriter, req *http.Request) {
			stockQuote, _ := getGoogleQuote(ticker)
			//stockQuote := quote{ticker, 0.0}
			stockQuoteString, err := json.Marshal(stockQuote)
			if err != nil {
				return
			}

			w.Header().Set("Content-Type", "application/json")
			w.Write(stockQuoteString)
			return
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
	http.HandleFunc("/quote", quoteHandler)

	if !testing {
		certManager := autocert.Manager{
			Prompt:     autocert.AcceptTOS,
			HostPolicy: autocert.HostWhitelist(domain),
		}

		bretsGoAPIServer = &http.Server{
			Addr: ":https",
			TLSConfig: &tls.Config{
				GetCertificate: certManager.GetCertificate,
			},
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
