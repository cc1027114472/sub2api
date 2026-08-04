package service

import (
	"github.com/shopspring/decimal"
)

// Money helpers for hot-path balance/cost arithmetic.
//
// Boundary note: many public APIs and DB columns still use float64. Prefer these
// helpers for intermediate accumulation/subtraction, then convert at the edge
// with MoneyFloat. A full schema migration to numeric/decimal remains follow-up work.

const moneyScale = 8

// MoneyAdd returns a+b using fixed-scale decimal arithmetic.
func MoneyAdd(a, b float64) float64 {
	return MoneyFloat(decimal.NewFromFloat(a).Add(decimal.NewFromFloat(b)))
}

// MoneySub returns a-b using fixed-scale decimal arithmetic.
func MoneySub(a, b float64) float64 {
	return MoneyFloat(decimal.NewFromFloat(a).Sub(decimal.NewFromFloat(b)))
}

// MoneyMul returns a*b using fixed-scale decimal arithmetic.
func MoneyMul(a, b float64) float64 {
	return MoneyFloat(decimal.NewFromFloat(a).Mul(decimal.NewFromFloat(b)))
}

// MoneyFloat converts decimal to float64 at moneyScale.
func MoneyFloat(d decimal.Decimal) float64 {
	f, _ := d.Round(moneyScale).Float64()
	return f
}

// LiveSessionCostUSD computes Live wall-clock cost from duration and $/minute rate.
func LiveSessionCostUSD(durationMs int, costPerMinuteUSD float64) float64 {
	if durationMs <= 0 || costPerMinuteUSD <= 0 {
		return 0
	}
	minutes := decimal.NewFromInt(int64(durationMs)).Div(decimal.NewFromInt(60000))
	return MoneyFloat(minutes.Mul(decimal.NewFromFloat(costPerMinuteUSD)))
}
