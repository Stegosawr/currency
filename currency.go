package currency

import (
	"crypto/tls"
	"errors"
	"io"
	"io/ioutil"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

const currAPIURL = "https://cdn.shopify.com/s/javascripts/currencies.js"

var reCurr = regexp.MustCompile(`(\w+):(\d*\.*\d+)`)

// GetCurrencies rates of the important currencies
func GetCurrencies() (map[string]float64, error) {

	data, err := get(currAPIURL)
	if err != nil {
		return nil, err
	}

	matchedCurr := reCurr.FindAllStringSubmatch(string(data), -1) //1=currencyKey 2=qouta
	if len(matchedCurr) < 1 {
		return nil, errors.New("no currencies found")
	}

	currencies := map[string]float64{}
	for _, match := range matchedCurr {
		currRate := parseCurrencyRate(match[2])
		if currRate == 0 {
			continue
		}

		currencies[match[1]] = currRate
	}

	return currencies, nil
}

func parseCurrencyRate(rate string) float64 {
	if strings.HasPrefix(rate, ".") {
		rate = "0" + rate
	}

	currRate, err := strconv.ParseFloat(rate, 64)
	if err != nil {
		return 0
	}

	return currRate
}

func get(URL string) ([]byte, error) {

	client := &http.Client{
		Transport: &http.Transport{
			DisableCompression:  true,
			TLSHandshakeTimeout: 10 * time.Second,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			IdleConnTimeout:     5 * time.Second,
		},
		Timeout: 5 * time.Minute,
	}

	req, err := http.NewRequest(http.MethodGet, URL, nil)
	if err != nil {
		return nil, err
	}

	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		if err != io.ErrUnexpectedEOF {
			return nil, err
		}
	}

	return body, nil

}
