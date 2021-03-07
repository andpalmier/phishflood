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
	phishingURL   	= "" // given URL
	numGoroutines 	= 10 // goroutines
	px            	= "" // proxies
	// iphone UA by default
	userAgent     	= "Mozilla/5.0 (iPhone; CPU iPhone OS 14_3 like Mac OS X) AppleWebKit/605.1.15 (KHTML, like " +
		"Gecko) CriOS/87.0.4280.77 Mobile/15E148 Safari/604.1"
	mindelay      	= 10
	maxdelay		= 3600
	seed int64 		= 0 // seed for random data generation
	asciiart string =
`
      _    _       _    ___  _              _ 
 ___ | |_ |_| ___ | |_ |  _|| | ___  ___  _| |
| . ||   || ||_ -||   ||  _|| || . || . || . |
|  _||_|_||_||___||_|_||_|  |_||___||___||___|
|_|                                           
by @andpalmier
`
)

func init() {
	flag.StringVar(&phishingURL, "url", phishingURL, "domain name or url, https is assumed.")
	flag.StringVar(&px, "proxies", px, "one or multiple proxies; specify the schema (http default) and port, and " +
		"use ',' as a separator.")
	flag.IntVar(&numGoroutines, "goroutines", numGoroutines, "number of goRoutines.")
	flag.StringVar(&userAgent, "ua", userAgent, "User Agent to be used, using Chrome on iPhone by default.")
	flag.IntVar(&mindelay, "dmin", mindelay, "minimun delay between consecutive requests, in seconds.")
	flag.IntVar(&maxdelay, "dmax", maxdelay, "maximum delay between consecutive requests, in seconds.")
	flag.Int64Var(&seed, "seed", seed, "seed for random data generation, random by default.")
}

// close execution with error
func die(format string, a ...interface{}) {
	fmt.Printf(format, a...)
	os.Exit(1)
}

func main() {

	flag.Parse()

	// check url was provided
	if phishingURL == "" {
		fmt.Printf("no -url specified.\n\n")
		flag.PrintDefaults()
		os.Exit(1)
	}

	// check delay provided
	if mindelay > maxdelay {
		die("minimum delay cannot be greater than maximum delay.\n")
	}

	// add schema
	if !strings.Contains(phishingURL, "://") {
		phishingURL = fmt.Sprintf("http://%s", phishingURL)
	}

	parsedProxies := []string{}
	// remove spaces from proxies list
	px = strings.ReplaceAll(px, " ", "")
	fmt.Println(asciiart)

	// check if proxies were provided
	if px == "" {
		fmt.Println("WARNING: no -proxies specified, this could expose your IP to the phishing kit. You have 5 " +
			 "seconds to stop the execution, or it will continue as it is.")
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
	postAction, inputNames, inputTypes := getPostData(phishingURL, parsedProxies, userAgent)
	if postAction == "" || len(inputNames) == 0 || len(inputTypes) == 0 {
		die("couldn't find a compatible form in the given page.\n")
	}

	fmt.Printf("\n[!] Found a form with action: %s \n[!] Input fields names found: %v\n[!] Input fields types " +
		"found: %v\n\n", postAction, inputNames, inputTypes)

	// set random seed if not provided by user
	rand.Seed(time.Now().UnixNano())
	if seed != 0 {
		rand.Seed(seed)
	}

	// create channel used for goroutines
	ch := make(chan string)

	// start routines
	for i := 0; i < numGoroutines; i++ {
		// waiting time
		var w int

		// user wants to wait exactly that time
		if maxdelay == mindelay {
			w = maxdelay
		} else {
			// wait for a random number of seconds between mindelay and maxdelay
			w = int(rand.Intn(maxdelay-mindelay) + mindelay)
		}
		time.Sleep(time.Duration(w) * time.Second)

		// if seed not declared -> set random
		seed = rand.Int63()

		// send requests with fake data
		go flood(i, postAction, inputNames, inputTypes, parsedProxies, ch, userAgent, seed)
	}

	// when POST request is completed, print the status code from the channel
	for i := 0; i < numGoroutines; i++ {
		fmt.Println(<-ch)
	}

}
