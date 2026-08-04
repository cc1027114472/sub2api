package service

import "testing"

func TestMoneyAddAvoidsBinaryFloatDrift(t *testing.T) {
	got := MoneyAdd(0.1, 0.2)
	if got != 0.3 {
		t.Fatalf("MoneyAdd(0.1, 0.2)=%v, want 0.3", got)
	}
}

func TestMoneySubAndMul(t *testing.T) {
	if got := MoneySub(1.0, 0.3); got != 0.7 {
		t.Fatalf("MoneySub(1, 0.3)=%v, want 0.7", got)
	}
	if got := MoneyMul(0.1, 3); got != 0.3 {
		t.Fatalf("MoneyMul(0.1, 3)=%v, want 0.3", got)
	}
}

func TestLiveSessionCostUSD(t *testing.T) {
	// 60s at $0.06/min => $0.06
	if got := LiveSessionCostUSD(60_000, 0.06); got != 0.06 {
		t.Fatalf("LiveSessionCostUSD(60s, 0.06)=%v, want 0.06", got)
	}
	if got := LiveSessionCostUSD(0, 0.06); got != 0 {
		t.Fatalf("zero duration should be free, got %v", got)
	}
	if got := LiveSessionCostUSD(1000, 0); got != 0 {
		t.Fatalf("zero rate should be free, got %v", got)
	}
}
