package api

import (
	"bytes"
	"encoding/json"
	"fmt"
	"go-tether/crypto"
	"go-tether/utils"
	"io"
	"net/http"
	"net/url"
	"time"
)

// doRequest performs the HTTP request with appropriate headers and signature
func doRequest(verb, uri string, options map[string]interface{}) (map[string]interface{}, error) {
	baseURI := utils.GetEnv("BASE_URI", "https://app.tether.to/api/v1")
	path := baseURI + uri
	var requestBody []byte
	var contentMD5 string
	var err error

	if verb == "GET" || verb == "DELETE" {
		if len(options) > 0 {
			query := url.Values{}
			for key, value := range options {
				query.Add(key, fmt.Sprintf("%v", value))
			}
			path += "?" + query.Encode()
		}
		contentMD5 = ""
	} else {
		requestBody, err = json.Marshal(options)
		if err != nil {
			return nil, err
		}
		contentMD5 = crypto.Md5Base64Digest(string(requestBody))
	}

	headers := map[string]string{
		"Content-MD5":  contentMD5,
		"Date":         time.Now().UTC().Format(http.TimeFormat),
		"Content-Type": "application/json",
	}

	canonicalString := fmt.Sprintf("%s,%s,%s,%s,%s",
		verb,
		headers["Content-Type"],
		headers["Content-MD5"],
		uri,
		headers["Date"])

	apiKey := utils.GetEnv("API_KEY", "")
	secret := utils.GetEnv("API_SECRET", "")

	signature := crypto.HmacSignature(secret, canonicalString)
	headers["Authorization"] = "APIAuth " + apiKey + ":" + signature

	client := &http.Client{}
	req, err := http.NewRequest(verb, path, bytes.NewBuffer(requestBody))
	if err != nil {
		return nil, err
	}

	for key, value := range headers {
		req.Header.Set(key, value)
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	var parsedResponse map[string]interface{}
	if err := json.Unmarshal(body, &parsedResponse); err != nil {
		return nil, err
	}
	return parsedResponse, nil
}

// get performs a GET request
func get(path string, options map[string]interface{}) (map[string]interface{}, error) {
	return doRequest("GET", path, options)
}

// balances returns current account balances
func Balances() (map[string]interface{}, error) {
	return get("/balances.json", nil)
}

// transactions returns list of most recent transactions
func Transactions() (map[string]interface{}, error) {
	return get("/transactions.json", nil)
}

// transactionPage returns a specific page of transactions
func TransactionPage(page int) (map[string]interface{}, error) {
	return get(fmt.Sprintf("/transactions/page/%d", page), nil)
}

// getTransaction returns details of a specific transaction
func GetTransaction(id int) (map[string]interface{}, error) {
	return get(fmt.Sprintf("/transactions/%d.json", id), nil)
}
