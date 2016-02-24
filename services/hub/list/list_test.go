package list

import (
	"errors"
	"testing"
	"time"
)

type mockConn struct {
	err error
}

func (m mockConn) Write([]byte) error {
	return m.err
}

func TestHub(t *testing.T) {
	h := New(func(error) {})
	n := 10000

	go func(count int) {
		for i := 0; i < count; i++ {
			h.PutConn(mockConn{nil})
		}
	}(n)

	h.Broadcast([]byte(`i am message`))
	time.Sleep(time.Second) // wait long enough for broadcast done

	if l := h.Len(); l != n {
		t.Errorf("expected %v connection but get %v", n, l)
	}

	h.PutConn(mockConn{errors.New("hey write error")})
	h.Broadcast([]byte(`i am message`))
	time.Sleep(time.Second) // wait long enough for broadcast done

	if l := h.Len(); l != n {
		t.Errorf("expected %v connection but get %v", n, l)
	}
}
