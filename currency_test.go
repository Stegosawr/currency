package currency

import "testing"

func TestGetCurrencies(t *testing.T) {
	expectedCurrencies := []string{
		"USD",
		"EUR",
		"JPY",
	}

	t.Run("Default Request Currency Exchange Rate", func(t *testing.T) {
		currencies, err := GetCurrencyRates()
		if err != nil {
			t.Error(err)
		}

		for _, eCurr := range expectedCurrencies {
			if _, ok := currencies[eCurr]; !ok {
				t.Errorf("Expected currency rate: %s not found in map", eCurr)
			}
		}
	})
}

func TestParseCurrencyRate(t *testing.T) {
	tests := []struct {
		Name       string
		CurrString string
		Want       float64
	}{
		{
			Name:       "Default",
			CurrString: "1.15998",
			Want:       1.15998,
		}, {
			Name:       "Rate under 1",
			CurrString: ".808357",
			Want:       0.808357,
		}, {
			Name:       "Invalid input",
			CurrString: "F",
			Want:       0.0,
		},
	}
	for _, tt := range tests {
		t.Run(tt.Name, func(t *testing.T) {
			rate := parseCurrencyRate(tt.CurrString)

			if rate != tt.Want {
				t.Errorf("Got: %v - want: %v", rate, tt.Want)
			}
		})
	}
}
