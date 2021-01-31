package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"time"
)

// make POST requests to the given url using the proxy
// generate data based on type or name -> to be improved
func flood(i int, postAction string, inputNames []string, inputTypes []string, parsedProxies []string, ch chan<- string, userAgent string) {

	var myClient *http.Client

	// make post request using proxy if specified
	if len(parsedProxies) != 0 {
		proxyURL, err := url.Parse(parsedProxies[i%len(parsedProxies)])
		if err != nil {
			die("error parsing the proxy address: %v\n", err)
		}
		myClient = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	} else {
		myClient = &http.Client{Timeout: 15 * time.Second}
	}

	// generate fake data
	vals := url.Values{}
	for _, valName := range inputNames {

		// "cellulare" stands for mobile phone, so we have a particular interval to make it realistic
		if valName == "cellulare" {
			val := rand.Intn(3499999999-3200000000) + 3200000000
			vals.Set(valName, fmt.Sprintf("%d", val))

			// these are generic numbers
		} else {
			val := rand.Intn(99999999-10000000) + 10000000
			vals.Set(valName, fmt.Sprintf("%d", val))
		}
	}

	// use bytes to encode the fake data for the post
	b := bytes.NewBufferString(vals.Encode())
	req, err := http.NewRequest("POST", postAction, b)
	if err != nil {
		die("could not make the POST request: %v \n", err)
	}

	// set User Agent
	req.Header.Set("User-Agent", userAgent)

	// make request
	resp, err := myClient.Do(req)
	if err != nil {
		die("could not make the POST request: %v \n", err)
	}

	// print error
	if err != nil {
		ch <- fmt.Sprintf("Request #%d terminated with error: %s", i+1, err)
	} else {
		// send to the channel the status code of the POST
		prettyVals, err := json.Marshal(vals)
		if err != nil {
			die("error on marshal for fake data: %v \n", err)
		}
		ch <- fmt.Sprintf("Request #%d with these parameters: %s returned the following status code: %d %s.\n", i+1, prettyVals, resp.StatusCode, http.StatusText(resp.StatusCode))
	}
}

// given phishingUrl, visit it using the first of the proxy list and the UserAgent
// get the action of the POST and the type and names of the attributes
func getPostData(phishingUrl string, parsedProxies []string, userAgent string) (string, []string, []string) {

	postAction := ""
	var inputNames []string
	var inputTypes []string
	var myClient *http.Client

	// make post request using proxy if specified
	if len(parsedProxies) != 0 {
		proxyURL, err := url.Parse(parsedProxies[0])
		if err != nil {
			die("error parsing the proxy address: %v\n", err)
		}
		myClient = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{Proxy: http.ProxyURL(proxyURL)}}
	} else {
		myClient = &http.Client{Timeout: 15 * time.Second}
	}

	// prepare request
	req, err := http.NewRequest("GET", phishingUrl, nil)
	if err != nil {
		die("could not make the GET request: %v \n", err)
	}
	// set User Agent
	req.Header.Set("User-Agent", userAgent)

	// make request
	resp, err := myClient.Do(req)
	if err != nil {
		die("could not make the GET request: %v \n", err)
	}

	// get response
	defer resp.Body.Close()
	if resp.StatusCode != 200 {
		die("status code error: %d %s \n", resp.StatusCode, resp.Status)
	}

	// Load the HTML document
	doc, err := goquery.NewDocumentFromReader(resp.Body)
	if err != nil {
		die("error on html page: %v \n", err)
	}

	// Find the form
	doc.Find("form").Each(func(i int, form *goquery.Selection) {
		action, actionOk := form.Attr("action")

		// found a form with action
		if actionOk {
			form.Find("input").Each(func(i int, input *goquery.Selection) {
				nameattr, nameOk := input.Attr("name")
				typeattr, typeOk := input.Attr("type")

				// find input with name and attributes
				if actionOk && nameOk && typeOk {
					inputNames = append(inputNames, nameattr)
					inputTypes = append(inputTypes, typeattr)
					u, err := url.Parse(phishingUrl)
					if err != nil {
						die("error parsing diven url: %v \n", err)
					}

					// create full url for path where to submit the form
					u.Path = path.Join(u.Path, action)
					postAction = u.String()
				}
			})
		}
	})

	return postAction, inputNames, inputTypes
}
