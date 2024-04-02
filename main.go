package main

import (
	"bytes"
	"encoding/xml"
	"fmt"
	"io"
	"log"
	"net/http"
)

type Envelope struct {
	Body Body `xml:"Body"`
}

type Body struct {
	Response ListOfCountryNamesByNameResponse `xml:"ListOfCountryNamesByNameResponse"`
}

type ListOfCountryNamesByNameResponse struct {
	Result ListOfCountryNamesByNameResult `xml:"ListOfCountryNamesByNameResult"`
}

type ListOfCountryNamesByNameResult struct {
	Countries []Country `xml:"tCountryCodeAndName"`
}

type Country struct {
	ISOCode string `xml:"sISOCode"`
	Name    string `xml:"sName"`
}

func main() {
	// define xml body
	xmlBody := `<?xml version="1.0" encoding="utf-8"?>
	<soap12:Envelope xmlns:soap12="http://www.w3.org/2003/05/soap-envelope">
	  <soap12:Body>
		<ListOfCountryNamesByName xmlns="http://www.oorsprong.org/websamples.countryinfo">
		</ListOfCountryNamesByName>
	  </soap12:Body>
	</soap12:Envelope>`

	// Create a new request with the XML body
	req, err := http.NewRequest("POST", "http://webservices.oorsprong.org/websamples.countryinfo/CountryInfoService.wso", bytes.NewBuffer([]byte(xmlBody)))
	if err != nil {
		log.Fatalf("failed to create request: %v", err)
	}

	// Set the Content-Type header
	req.Header.Set("Content-Type", "text/xml")

	// Make the request
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatalf("failed to make request: %v", err)
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Fatalf("failed to read response body: %v", err)
	}

	// Unmarshal the XML response into a struct
	var envelope Envelope
	if err := xml.Unmarshal(body, &envelope); err != nil {
		log.Fatalf("failed to unmarshal XML: %v", err)
	}

	// Extract country names
	var countryNames []string

	for _, country := range envelope.Body.Response.Result.Countries {
		countryNames = append(countryNames, country.Name)
	}

	// Print country names
	fmt.Println(countryNames)

	fmt.Println(envelope.Body.Response.Result.Countries)

}
