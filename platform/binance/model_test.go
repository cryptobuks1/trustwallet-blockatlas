package binance

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

var newOrderDataDst = Data{OrderData: OrderData{
	Symbol:   "AWC-986_BNB",
	Base:     "AWC-986",
	Quote:    "BNB",
	Quantity: 2.0,
	Price:    0.00324939,
}}

var cancelOrderDataDst = Data{OrderData: OrderData{
	Symbol:   "GTO-908_BNB",
	Base:     "GTO-908",
	Quote:    "BNB",
	Quantity: 1.0,
	Price:    0.00104716,
}}

func TestTx_getData(t *testing.T) {
	tests := []struct {
		name string
		Data string
		want Data
	}{
		{
			"new order",
			"{\"orderData\":{\"symbol\":\"AWC-986_BNB\",\"orderType\":\"limit\",\"side\":\"buy\",\"price\":0.00324939,\"quantity\":2.00000000,\"timeInForce\":\"GTE\",\"orderId\":\"D13BAF4BD6638FA3AAD6EBCA0E4BEEA73DF4D519-30\"}}",
			newOrderDataDst,
		},
		{
			"cancel order",
			"{\"orderData\":{\"symbol\":\"GTO-908_BNB\",\"orderType\":\"limit\",\"side\":\"buy\",\"price\":0.00104716,\"quantity\":1.00000000,\"timeInForce\":\"GTE\",\"orderId\":\"D13BAF4BD6638FA3AAD6EBCA0E4BEEA73DF4D519-28\"}}",
			cancelOrderDataDst,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tx := &Tx{Data: tt.Data}
			got, _ := tx.getData()
			assert.Equal(t, tt.want, got)
		})
	}
}

func TestConvertValue(t *testing.T) {
	tests := []struct {
		name       string
		value      interface{}
		wantResult float64
		wantError  bool
	}{
		{"test string 1", "9", 9, false},
		{"test number 1", 9, 9, false},
		{"test string 2", "9380938973", 9380938973, false},
		{"test number 2", 9380938973, 9380938973, false},
		{"test string 3", "0.0000003", 0.0000003, false},
		{"test number 3", 0.0000003, 0.0000003, false},
		{"test string 4", "0.44", 0.44, false},
		{"test number 4", 0.44, 0.44, false},
		{"test string 5", "3334", 3334, false},
		{"test number 5", 3334, 3334, false},
		{"test error", time.Time{}, 3334, true},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, ok := convertValue(tt.value)
			if tt.wantError {
				assert.False(t, ok)
				return
			}
			assert.True(t, ok)
			assert.Equal(t, tt.wantResult, got)
		})
	}
}

func Test_removeFloatPoint(t *testing.T) {
	tests := []struct {
		name  string
		value float64
		want  int64
	}{
		{"test float 1", 0.0034, 340000},
		{"test float 2", 0.00000013, 13},
		{"test float 3", 0.938984, 93898400},
		{"test float 4", 0.1, 10000000},
		{"test int 1", 12, 1200000000},
		{"test int 2", 2333333333, 233333333300000000},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := removeFloatPoint(tt.value); got != tt.want {
				t.Errorf("removeFloatPoint() = %v, want %v", got, tt.want)
			}
		})
	}
}
