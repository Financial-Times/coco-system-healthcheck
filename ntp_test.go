package main

import (
	"fmt"
	"testing"
)

func ntpCmdMockGood() string {
	out := []byte(`
	ntpd: reply from 212.45.32.250: offset:-0.001673 delay:0.011379 status:0x24 strat:2 refid:0xca4f43c1 rootdelay:0.002960 reach:0x01
	ntpd: reply from 212.45.32.250: offset:-0.001803 delay:0.011673 status:0x24 strat:2 refid:0xca4f43c1 rootdelay:0.002960 reach:0x03
	ntpd: reply from 212.45.32.250: offset:-0.001630 delay:0.010400 status:0x24 strat:2 refid:0xca4f43c1 rootdelay:0.002960 reach:0x07
	ntpd: reply from 212.45.32.250: offset:-0.000613 delay:0.012454 status:0x24 strat:2 refid:0xca4f43c1 rootdelay:0.002960 reach:0x0f
	ntpd: reply from 212.45.32.250: offset:-0.001514 delay:0.011780 status:0x24 strat:2 refid:0xca4f43c1 rootdelay:0.002960 reach:0x1f
	ntpd: reply from 212.45.32.250: offset:-0.001406 delay:0.011294 status:0x24 strat:2 refid:0xca4f43c1 rootdelay:0.002960 reach:0x3f
	Alarm clock
	`)
	return string(out)
}

func TestNTPCmd(t *testing.T) {
	offset := ntpOffset(ntpCmdMockGood)
	fmt.Printf("offset: %s\n", offset)
	if offset.err != nil {
		t.Fail()
	}
}

func TestNTPCheck(t *testing.T) {
	cases := []struct {
		offset offsetResult
		err    error
	}{
		{offsetResult{val: "0"}, nil},
		{offsetResult{val: "101"}, fmt.Errorf("offset is greater then limit of 100: 101")},
		{offsetResult{val: "-101"}, fmt.Errorf("offset is greater then limit of 100: -101")},
	}

	offsetCh = make(chan offsetResult)
	ntpChecker := ntpChecker{}

	for _, test := range cases {
		go func() {
			offsetCh <- test.offset
		}()

		val, err := ntpChecker.Check()
		if err != test.err && err.Error() == test.err.Error() {
			t.Fail()
		}
		if val != test.offset.val {
			t.Errorf("Offset is: '%v', wanted: '%v'", val, test.offset.val)

		}
	}
}
