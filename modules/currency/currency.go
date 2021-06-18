package currency

import (
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/Avimitin/go-bot/modules/net"
)

var (
	buffer = struct {
		currencies     map[string]string
		lastUpdateTime time.Time
		locker         sync.RWMutex
		exchangeRate   map[string]map[string]float64
	}{
		currencies:   make(map[string]string),
		exchangeRate: make(map[string]map[string]float64),
	}
)

func StoreRateInBuffer(from, to string, rate float64) {
	buffer.locker.Lock()
	defer buffer.locker.Unlock()

	buffer.exchangeRate[from] = map[string]float64{to: rate}
}

func BufferHasRate(from, to string) (float64, error) {
	buffer.locker.RLock()
	defer buffer.locker.RUnlock()

	if codeTo, ok := buffer.exchangeRate[from]; ok {
		if rate, ok := codeTo[to]; ok {
			return rate, nil
		}
	}

	return 0, fmt.Errorf("No data")
}

// StoreBuffer release buffer and replace the inner data with
// given data.
func StoreBuffer(currencies map[string]string) {
	buffer.locker.Lock()
	defer buffer.locker.Unlock()

	// clean old buffer
	if buffer.currencies != nil {
		buffer.currencies = nil
	}
	buffer.currencies = currencies

	// reset timer
	buffer.lastUpdateTime = time.Now()
}

// BufferIsOutDated return true if It has been more than 24 hours since
// the last cache being update.
func BufferIsOutDated() bool {
	buffer.locker.RLock()
	defer buffer.locker.RUnlock()

	if buffer.lastUpdateTime.IsZero() {
		return true
	}

	return time.Now().Sub(buffer.lastUpdateTime).Hours() > 24
}

// BufferIsNil return true if the buffer is empty
func BufferIsNil() bool {
	buffer.locker.RLock()
	defer buffer.locker.RUnlock()

	return len(buffer.currencies) == 0
}

// ListAllCurrenciesDescriptions return a map of the currency code and
// currency code descriptions. Return error if request failed.
func ListAllCurrenciesDescriptions() (map[string]string, error) {
	if !BufferIsOutDated() && !BufferIsNil() {
		return buffer.currencies, nil
	}

	response, err := net.Get(
		"https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies.min.json",
	)
	if err != nil {
		return nil, fmt.Errorf("request currency API: %w", err)
	}

	var currencies = make(map[string]string)
	err = json.Unmarshal(response, &currencies)
	if err != nil {
		return nil, fmt.Errorf("failed to unmarshall currencies result: %w", err)
	}

	StoreBuffer(currencies)

	return currencies, nil
}

// ListAllCurrencies return list of currency code that can be used for
// calculate rate. Return error if ListAllCurrenciesDescriptions failed.
func ListAllCurrencies() ([]string, error) {
	currencies, err := ListAllCurrenciesDescriptions()

	if err != nil {
		return nil, err
	}

	var currenciesKey = make([]string, 150)

	for k := range currencies {
		currenciesKey = append(currenciesKey, k)
	}

	return currenciesKey, nil
}

// Exchange fetch the rate about the two parameters. Valid parameters can be found
// by ListAllCurrencies. Return error if request failed.
func Exchange(from, to string) (map[string]interface{}, error) {
	response, err := net.Get(
		fmt.Sprintf(
			"https://cdn.jsdelivr.net/gh/fawazahmed0/currency-api@1/latest/currencies/%s/%s.json",
			from,
			to,
		),
	)
	if err != nil {
		return nil, fmt.Errorf("request exchange api: %w", err)
	}

	var result map[string]interface{}

	err = json.Unmarshal(response, &result)
	if err != nil {
		return nil, fmt.Errorf("decode exchange result: %w", err)
	}

	return result, nil
}

func FetchRate(from, to string) (float64, error) {
	if !BufferIsOutDated() {
		if rate, err := BufferHasRate(from, to); err == nil {
			return rate, nil
		}
	}

	result, err := Exchange(from, to)
	if err != nil {
		return 0, err
	}

	rateRaw, have := result[to]
	if !have {
		return 0, fmt.Errorf("unexpect currency code %s", to)
	}

	rate, ok := rateRaw.(float64)
	if !ok {
		return 0, fmt.Errorf("illegal rate value")
	}

	StoreRateInBuffer(from, to, rate)

	return rate, nil
}

// CalculateExchange convert the currency amount with given currency code.
func CalculateExchange(amount float64, from, to string) (float64, error) {
	if !CodeIsValid(from) {
		return 0, fmt.Errorf("%q is not valid code", from)
	}

	if !CodeIsValid(to) {
		return 0, fmt.Errorf("%q is not valid code", to)
	}

	rate, err := FetchRate(from, to)
	if err != nil {
		return 0, err
	}

	if rate <= 0 {
		return 0, fmt.Errorf("rate is unexpected")
	}

	return amount * rate, nil
}

// CodeIsValid test if the given code is supported
func CodeIsValid(code string) bool {
	result, err := ListAllCurrenciesDescriptions()
	if err != nil {
		return false
	}

	_, ok := result[code]

	return ok
}
