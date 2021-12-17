package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	"github.com/brianvoe/gofakeit/v6"
	"math/rand"
	"net/http"
	"net/url"
	"path"
	"strconv"
	"strings"
	"time"
)

func makeClient(i int, parsedProxies []string) (*http.Client, error) {

	var myClient *http.Client

	// make post request using proxy if specified
	if len(parsedProxies) != 0 {
		proxyURL, err := url.Parse(parsedProxies[i%len(parsedProxies)])
		if err != nil {
			err := errors.New("error parsing the proxy address")
			return nil, err
		}
		myClient = &http.Client{Timeout: 15 * time.Second, Transport: &http.Transport{
			Proxy: http.ProxyURL(proxyURL)}}
	} else {
		myClient = &http.Client{Timeout: 15 * time.Second}
	}
	return myClient, nil
}

// make POST requests to the given url using the proxy
// generate data based on type or name -> to be improved
func flood(i int, postAction string, inputNames []string, inputTypes []string, parsedProxies []string,
	ch chan<- string, userAgent string, seed int64) {

	myClient, err := makeClient(i, parsedProxies)
	if err != nil {
		die("%v\n", err)
	}

	// generate fake data
	vals := url.Values{}

	// handle in future
	gofakeit.Seed(seed)
	var val string
	for _, valName := range inputNames {
		switch valName {
		// cellulare is mobile phone in italian
		case "cellulare":
			val = gofakeit.Phone()
		case "login":
			val = gofakeit.Email()
		case "password":
			val = gofakeit.Password(true, false, true, false, false, rand.Intn(16-8)+8)
		// by default use random number (to keep compatibility with phishing kit)
		default:
			val = strconv.Itoa(gofakeit.Number(1111111, 9999999))
		}
		vals.Set(valName, fmt.Sprintf("%s", val))
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
		ch <- fmt.Sprintf("Request #%d with these parameters: %s returned the following status "+
			"code: %d %s.\n", i+1, prettyVals, resp.StatusCode, http.StatusText(resp.StatusCode))
	}
}

// given phishingUrl, visit it using the first of the proxy list and the UserAgent
// get the action of the POST and the type and names of the attributes
func getPostData(phishingURL string, parsedProxies []string, userAgent string) (string, []string, []string) {

	postAction := ""
	var inputNames []string
	var inputTypes []string

	myClient, err := makeClient(0, parsedProxies)
	if err != nil {
		die("%v\n", err)
	}

	// prepare request
	req, err := http.NewRequest("GET", phishingURL, nil)
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
				if strings.ToLower(typeattr) == "submit" {
					typeOk = false
				}
				_, hiddenOk := input.Attr("hidden")

				// find input with name and attributes
				if nameOk && typeOk && typeattr != "hidden" && !hiddenOk {
					inputNames = append(inputNames, nameattr)
					inputTypes = append(inputTypes, typeattr)
					u, err := url.Parse(phishingURL)
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
