package main

import (
	"encoding/json"
	"errors"
	"fmt"
	//"io/ioutil"
	"log"
	"flag"
	"strings"

	"github.com/go-resty/resty/v2"
)

const ReconDevSearchEndpoint string = "https://recon.dev/api/search"

const JsonKeyword="\"rawDomains\":"

type ReconDevJSONOut []struct {
	Domains    []string `yaml:"domains"`
	IP         string   `yaml:"ip"`
	RawDomains []string `yaml:"rawDomains"`
	RawPort    string   `yaml:"rawPort"`
	RawIP      string   `yaml:"rawIp"`
}

// Function queries the domains/IPs parsed from the web output of Recon.dev API 
func queryReconDevAPI(apiKey, domain string) (int, string, error) {
	var e error
	
	url := fmt.Sprintf("%s?key=%s&domain=%s", ReconDevSearchEndpoint, apiKey, 
		domain)

	statusCode := 0
	respText := ""

	log.Printf("Making request to URL: %s", url)
	client := resty.New()
	resp, e := client.R().
		EnableTrace().
		Get(url)

	if e != nil {
		errStr := fmt.Sprintf("Error in request. Err: %s\n", e.Error())
		e = errors.New(errStr)
	} else {

		statusCode = resp.StatusCode()
		respText = resp.String()
		if statusCode == 200 {
			log.Printf("Response received: %d", statusCode)
			// Check if keyword in response as expected
			if !strings.Contains(respText, JsonKeyword) {
				errStr := fmt.Sprintf("Error: Output may not be JSON. Resp: %s", respText)
				e = errors.New(errStr)
			}

		} else {
			errStr := fmt.Sprintf("Error encountered in response. Status code: %d, %s\n", 
				statusCode, respText)
			e = errors.New(errStr)
		}
	}

	return statusCode, respText, e
}

// Function returns the domains/IPs parsed from the JSON output of Recon.dev API 
func parseReconDevOutput(respText string) ([]string, []string, error) {
	var e error
	var rd ReconDevJSONOut 

	// Create objects for storing maps,ips uniquely
	domainsMap := make(map[string]bool)
	ipsMap:= make(map[string] bool)
	var domains []string
	var ips []string

	// Parse the JSON response 
	//respTextBytes, _ := ioutil.ReadFile("/tmp/test.json")
	respTextBytes := []byte(respText)
	ej := json.Unmarshal(respTextBytes, &rd)
	if ej != nil {
		// If there were any errors, record them
		errStr := fmt.Sprintf("JSON Unmarshal error: %v", ej.Error())
		e = errors.New(errStr)

	} else {
		// Extract the domains/IPs from the response
		for _, r := range(rd) {
			for _, domain := range(r.RawDomains) { 
				domainsMap[domain] = true
			}
			ip := r.RawIP
			ipsMap[ip] = true
		}
	}
	
	// Extract IPs, domains
	for domain := range(domainsMap) {
		domains = append(domains, domain)
	}
	for ip := range(ipsMap) {
		ips = append(ips, ip)
	}
	return domains, ips, e
}

func main() {
	apiKey := ""
	domain := ""
	flag.StringVar(&apiKey, "apiKey", "", "API key for recon.dev")
	flag.StringVar(&domain, "domain", "", "Domain to query via recon.dev")
	flag.Parse()

	if apiKey == "" {
		log.Fatalf("API Key must be provided")
	}

	if domain == "" {
		log.Fatalf("Domain must be provided")
	}

	_, respText, e := queryReconDevAPI(apiKey, domain)
	if e != nil {
		log.Fatalf(e.Error())
	}
	//respText := ""

	ips, domains, e := parseReconDevOutput(respText)
	if e != nil {
		log.Fatalf(e.Error())
	}
	log.Printf("Number of domains found: %d\n", len(domains))
	log.Printf("Number of ips found: %d\n", len(ips))
	
	for _, ip := range(ips) {
		fmt.Println(ip)
	}
	for _, domain := range(domains) {
		fmt.Println(domain)
	}
}