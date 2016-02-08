package memory

import (
	"bytes"
	"errors"
	"testing"
	"time"
)

func mockGetBitcoinPriceFunc(p int64, err error) getBitcoinPriceFunc {
	return func() (int64, error) {
		return p, err
	}
}

func TestMemory(t *testing.T) {
	c := New(mockGetBitcoinPriceFunc(8, nil), &bytes.Buffer{}, time.Second)

	p := c.GetBitcoinPrice()
	if p != 8 {
		t.Errorf("price should be 8.8 but get %v", p)
	}

	c.getBTCPrice = mockGetBitcoinPriceFunc(0, nil)
	c.setBitcoinPrice(false)

	c.getBTCPrice = mockGetBitcoinPriceFunc(0, errors.New(""))
	c.setBitcoinPrice(false)

	funcWithRecover(func() {
		c.getBTCPrice = mockGetBitcoinPriceFunc(0, nil)
		c.setBitcoinPrice(true)
	})

	funcWithRecover(func() {
		c.getBTCPrice = mockGetBitcoinPriceFunc(0, errors.New(""))
		c.setBitcoinPrice(true)
	})

	funcWithRecover(func() {
		c.getBTCPrice = mockGetBitcoinPriceFunc(0, errors.New("error"))
		c.backgroundJob(true, time.Second)
	})
}

func funcWithRecover(f func()) {
	defer func() {
		recover()
	}()
	f()
}
