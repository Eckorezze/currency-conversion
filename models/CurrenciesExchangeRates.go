package models

import "time"

type CurrenciesExchangeRates struct {
	Id              int
	CurrencyId      int
	TargetCurencyId int
	ExchangeRate    float32 `json:"value"`
	RateSourceId    int
	CreatedAt       time.Time
	UpdatedAt       time.Time
}