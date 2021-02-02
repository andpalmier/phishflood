package main

import (
	"flag"
	"fmt"
	"math/rand"
	"os"
	"strings"
	"time"
)

var (
	phishingUrl   = ""
	numGoroutines = 10
	px            = ""
	userAgent     = "Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like Gecko) CriOS/87.0.4280.77 Mobile/15E148 Safari/604.1"
	mindelay      = 10
	maxdelay      = 3600
)

func init() {
	flag.StringVar(&phishingUrl, "url", phishingUrl, "domain name or url, if schema is not specified, https is assumed.")
	flag.StringVar(&px, "proxies", px, "one or multiple proxies; specify the schema (http default) and port, and use ',' as a separator.")
	flag.IntVar(&numGoroutines, "goroutines", numGoroutines, "number of goRoutines.")
	flag.StringVar(&userAgent, "ua", userAgent, "User Agent to be used, using Chrome on iPhone by default.")
	flag.IntVar(&mindelay, "dmin", mindelay, "minimun delay between consecutive requests, in seconds.")
	flag.IntVar(&maxdelay, "dmax", maxdelay, "maximum delay between consecutive requests, in seconds.")
}

// close execution with error
func die(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func main() {

	flag.Parse()

	// check url was provided
	if phishingUrl == "" {
		fmt.Printf("no -url specified.\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// check delay provided
	if mindelay > maxdelay {
		die("minimum delay cannot be greater than maximum delay.\n")
	}

	// add schema
	if !strings.Contains(phishingUrl, "://") {
		phishingUrl = fmt.Sprintf("http://%s", phishingUrl)
	}

	parsedProxies := []string{}
	// remove spaces from proxies list
	px = strings.ReplaceAll(px, " ", "")

	// check if proxies were provided
	if px == "" {
		fmt.Printf("WARNING: no -proxies specified, this could expose your IP to the phishing kit. I'll give you 5 seconds to stop the execution, then it will continue at your own risk.\n\n")
		time.Sleep(time.Duration(5) * time.Second)
	} else {
		// split proxies in a list and add schema if not provided
		parsedProxies = strings.Split(px, ",")
		for i, proxy := range parsedProxies {
			if !strings.Contains(proxy, "://") {
				parsedProxies[i] = fmt.Sprintf("http://%s", proxy)
			}
		}
	}

	// go to the url and print findings
	postAction, inputNames, inputTypes := getPostData(phishingUrl, parsedProxies, userAgent)
	if postAction == "" || len(inputNames) == 0 || len(inputTypes) == 0 {
		die("couldn't find a compatible form in the given page.\n")
	}

	fmt.Printf("[!] Found a form with action: %s \n[!] Input fields names found: %v\n[!] Input fields types found: %v\n\n", postAction, inputNames, inputTypes)

	// set random seed
	rand.Seed(time.Now().UnixNano())

	// create channel used for goroutines
	ch := make(chan string)

	// start routines
	for i := 0; i < numGoroutines; i++ {

		// wait for a random number of seconds between mindelay and maxdelay
		w := int(rand.Intn(maxdelay - mindelay) + mindelay)
		time.Sleep(time.Duration(w) * time.Second)

		// send requests with fake data
		go flood(i, postAction, inputNames, inputTypes, parsedProxies, ch, userAgent)
	}

	// when POST request is completed, print the status code from the channel
	for i := 0; i < numGoroutines; i++ {
		fmt.Println(<-ch)
	}

}
