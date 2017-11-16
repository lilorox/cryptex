package kraken

import (
	"crypto/hmac"
	"crypto/sha256"
	"crypto/sha512"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"mime"
	"net/http"
	"net/url"
	"strings"
	"time"
)

const (
	APIURL       = "https://api.kraken.com"
	APIVersion   = "0"
	APIUserAgent = "Cryptex Kraken Client (https://github.com/lilorox/cryptex)"
)

type KrakenClient struct {
	apiKey     string
	privKey    string
	httpClient *http.Client
}

func New(apiKey, privKey string) *KrakenClient {
	return &KrakenClient{
		apiKey:     apiKey,
		privKey:    privKey,
		httpClient: http.DefaultClient,
	}
}

func (kc *KrakenClient) GetName() string {
	return "Kraken"
}

func (kc *KrakenClient) Balance() (map[string]*big.Float, error) {
	resp, err := kc.queryPrivate("Balance", url.Values{}, &BalanceResponse{})
	if err != nil {
		return nil, err
	}

	balance := resp.(*BalanceResponse)
	return *balance, nil
}

/*
 * Private methods
 */

// Execute a public method query
func (kc *KrakenClient) queryPublic(method string, values url.Values, typ interface{}) (interface{}, error) {
	url := fmt.Sprintf("%s/%s/public/%s", APIURL, APIVersion, method)
	resp, err := kc.doRequest(url, values, nil, typ)

	return resp, err
}

// queryPrivate executes a private method query
func (kc *KrakenClient) queryPrivate(method string, values url.Values, typ interface{}) (interface{}, error) {
	urlPath := fmt.Sprintf("/%s/private/%s", APIVersion, method)
	reqURL := fmt.Sprintf("%s%s", APIURL, urlPath)
	secret, _ := base64.StdEncoding.DecodeString(kc.privKey)
	values.Set("nonce", fmt.Sprintf("%d", time.Now().UnixNano()))

	// Create signature
	signature := createSignature(urlPath, values, secret)

	// Add Key and signature to request headers
	headers := map[string]string{
		"API-Key":  kc.apiKey,
		"API-Sign": signature,
	}

	return kc.doRequest(reqURL, values, headers, typ)
}

// doRequest executes a HTTP Request to the Kraken API and returns the result
func (kc *KrakenClient) doRequest(reqURL string, values url.Values, headers map[string]string, typ interface{}) (interface{}, error) {
	// Create request
	req, err := http.NewRequest("POST", reqURL, strings.NewReader(values.Encode()))
	if err != nil {
		return nil, fmt.Errorf("Could not create a new POST request: %s", err.Error())
	}

	req.Header.Add("User-Agent", APIUserAgent)
	for key, value := range headers {
		req.Header.Add(key, value)
	}

	// Execute request
	resp, err := kc.httpClient.Do(req)
	if err != nil {
		return nil, fmt.Errorf("Could not execute request: %s", err.Error())
	}
	defer resp.Body.Close()

	// Read request
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Could not read the response: %s", err.Error())
	}
	fmt.Printf("DEBUG, body=%+v\n", string(body))

	// Check mime type of response
	mimeType, _, err := mime.ParseMediaType(resp.Header.Get("Content-Type"))
	if err != nil {
		return nil, fmt.Errorf("Could not parse response Content-Type: %s", err.Error())
	}
	if mimeType != "application/json" {
		return nil, fmt.Errorf("Response Content-Type is '%s' instead of 'application/json'.", mimeType)
	}

	// Parse request
	var jsonData GenericResponse

	// Set the GenericResponse.Result to typ so `json.Unmarshal` will
	// unmarshal it into given typ instead of `interface{}`.
	if typ != nil {
		jsonData.Result = typ
	}

	err = json.Unmarshal(body, &jsonData)
	if err != nil {
		return nil, fmt.Errorf("Could not parse the response body: %s", err.Error())
	}

	// Check for Kraken API error
	if len(jsonData.Error) > 0 {
		return nil, fmt.Errorf("API returned an error: %s", jsonData.Error)
	}

	return jsonData.Result, nil
}

// getSha256 creates a sha256 hash for given []byte
func getSha256(input []byte) []byte {
	sha := sha256.New()
	sha.Write(input)
	return sha.Sum(nil)
}

// getHMacSha512 creates a hmac hash with sha512
func getHMacSha512(message, secret []byte) []byte {
	mac := hmac.New(sha512.New, secret)
	mac.Write(message)
	return mac.Sum(nil)
}

func createSignature(urlPath string, values url.Values, secret []byte) string {
	// See https://www.kraken.com/help/api#general-usage for more information
	shaSum := getSha256([]byte(values.Get("nonce") + values.Encode()))
	macSum := getHMacSha512(append([]byte(urlPath), shaSum...), secret)
	return base64.StdEncoding.EncodeToString(macSum)
}
