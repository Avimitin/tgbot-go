package currency

import (
	"testing"
	"time"
)

func TestListAllCurrenciesDesc(t *testing.T) {
	currencies, err := ListAllCurrenciesDescriptions()

	if err != nil {
		t.Fatal(err)
	}

	var testShortName = "cny"
	get := currencies[testShortName]
	want := "Chinese Yuan"

	if get != want {
		t.Errorf("get %s want %s", get, want)
	}
}

func TestListAvailableCurrencies(t *testing.T) {
	currencies, err := ListAllCurrencies()
	if err != nil {
		t.Fatal(err)
	}

	if len(currencies) == 0 {
		t.Errorf("no currencies available")
	}
}

func TestStoreBuffer(t *testing.T) {
	var testBuffer = map[string]string{"foo": "bar"}
	StoreBuffer(testBuffer)

	for key, value := range buffer.currencies {
		if testBuffer[key] != value {
			t.Errorf("got %s want %s", value, testBuffer[key])
		}
	}
}

func TestBufferOutDated(t *testing.T) {
	// let buffer time behind test time
	buffer.lastUpdateTime = time.Now().AddDate(0, 0, -2)
	if !BufferIsOutDated() {
		t.Errorf("buffer should out of dated")
	}
}

func TestExchange(t *testing.T) {
	var (
		currencyFixtures = []string{"cny", "usd"}
	)
	result, err := Exchange(currencyFixtures[0], currencyFixtures[1])
	if err != nil {
		t.Fatal(err)
	}

	date, ok := result["date"].(string)
	if !ok {
		t.Error("response is not expect")
	}

	data, ok := result[currencyFixtures[1]].(float64)
	if !ok {
		t.Error("response is not expect")
	}

	if date == "" || data == 0 {
		t.Errorf("unexpected date %s and data %f", date, data)
	}
}

func TestCalculate(t *testing.T) {
	result, err := CalculateExchange(100, "cny", "usd")
	if err != nil {
		t.Fatal(err)
	}

	if result == 0 {
		t.Fatal("no data calculated")
	}

	_, err = CalculateExchange(100, "aaa", "ccc")
	if err == nil {
		t.Error("expect error output")
	}
}
