package apns

import (
	"io"
	"testing"
)

func TestNotificationError(t *testing.T) {
	{
		p := []byte{8, 1, 1, 2, 3, 4}
		e := NewNotificationError(p, nil)
		got := e.Error()
		expect := "Processing error(1): id(1020304)"
		if got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		p := []byte{1, 1, 1, 2, 3, 4}
		e := NewNotificationError(p, nil)
		got := e.Error()
		expect := "Unknow error, command(1), status(1), id(1020304)"
		if got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		p := []byte{1, 1, 1, 2, 3, 4}
		e := NewNotificationError(p, io.EOF)
		got := e.Error()
		expect := "Unknow error, command(1), status(1), id(1020304)"
		if got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		p := []byte{1, 2, 3, 4}
		e := NewNotificationError(p, nil)
		got := e.Error()
		expect := "Wrong data format, [01020304]"
		if got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		p := []byte{}
		e := NewNotificationError(p, io.EOF)
		got := e.Error()
		expect := "EOF"
		if got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		p := []byte{1, 2, 3, 4}
		e := NewNotificationError(p, io.EOF)
		got := e.Error()
		expect := "EOF"
		if got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}

	{
		e := NewNotificationError(nil, io.EOF)
		got := e.Error()
		expect := "EOF"
		if got != expect {
			t.Errorf("got: %s, expect: %s", got, expect)
		}
	}
}
