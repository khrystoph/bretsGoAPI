package main

import (
	"flag"
	"io"
	"log"
	"net/http"
	"os"
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
	Error          *log.Logger
	traceHandle    io.Writer
	infoHandle     io.Writer = os.Stdout
	warningHandle  io.Writer = os.Stderr
	errorHandle    io.Writer = os.Stderr
	domain, ticker string
)

type quote struct {
	stockTicker string
	askPrice    float64
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
	flag.StringVar(&ticker, "ticker", "amzn", "enter the ticker symbol you want")
}

//getGoogleQuote will make a request out to the google finance apis and return the close price for the day
func getGoogleQuote(symbol string) (stockQuote quote, err error) {
	stockQuote = quote{stockTicker: ticker, askPrice: 0.00}
	return stockQuote, nil
}

func main() {
	flag.Parse()

	var (
		helloHandler = func(w http.ResponseWriter, _ *http.Request) {
			io.WriteString(w, "Welcome to Bret's API!\n")
		}
		quoteHandler = func(w http.ResponseWriter, req *http.Request) {
			getGoogleQuote(ticker)
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

	//certManager := autocert.Manager{
	//	Prompt:     autocert.AcceptTOS,
	//	HostPolicy: autocert.HostWhitelist(domain),
	//}

	//bretsGoAPIServer := &http.Server{
	//	Addr: ":https",
	//	TLSConfig: &tls.Config{
	//		GetCertificate: certManager.GetCertificate,
	//	},
	//}

	bretsGoAPIServer := &http.Server{
		Addr: ":8080",
	}

	//var wg sync.WaitGroup
	Info.Printf("Starting the letsencrypt server\n")
	//go func() {
	//	wg.Add(1)
	//	defer wg.Done()
	//	http.ListenAndServe(":http", certManager.HTTPHandler(nil))
	//}()

	Info.Printf("Starting the main TLS server.\n")

	Error.Fatal(bretsGoAPIServer.ListenAndServe())
	//Error.Fatal(bretsGoAPIServer.ListenAndServeTLS("", ""))

	return
}
