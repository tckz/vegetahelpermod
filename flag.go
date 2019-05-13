/*
 * https://github.com/tsenart/vegeta/blob/master/flags.go
 */

package vegetahelper

import (
	"fmt"
	"strconv"
	"strings"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type RateFlag struct{ *vegeta.Rate }

func (f *RateFlag) Set(v string) (err error) {
	ps := strings.SplitN(v, "/", 2)
	switch len(ps) {
	case 1:
		ps = append(ps, "1s")
	case 0:
		return fmt.Errorf("-rate format %q doesn't match the \"freq/duration\" format (i.e. 50/1s)", v)
	}

	f.Freq, err = strconv.Atoi(ps[0])
	if err != nil {
		return err
	}

	switch ps[1] {
	case "ns", "us", "µs", "ms", "s", "m", "h":
		ps[1] = "1" + ps[1]
	}

	f.Per, err = time.ParseDuration(ps[1])
	return err
}

func (f *RateFlag) String() string {
	if f.Rate == nil {
		return ""
	}
	return fmt.Sprintf("%d/%s", f.Freq, f.Per)
}
