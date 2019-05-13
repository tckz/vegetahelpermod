package vegetahelper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	vegeta "github.com/tsenart/vegeta/lib"
)

func TestRateFlagSet(t *testing.T) {
	assert := assert.New(t)

	t.Run("blank", func(t *testing.T) {
		rate := vegeta.Rate{}
		rf := &RateFlag{Rate: &rate}
		err := rf.Set("")
		assert.EqualError(err, `strconv.Atoi: parsing "": invalid syntax`)
	})

	t.Run("10", func(t *testing.T) {
		rate := vegeta.Rate{}
		rf := &RateFlag{Rate: &rate}
		err := rf.Set("10")
		assert.Nil(err)
		assert.Equal(10, rate.Freq)
		assert.Equal(time.Second, rate.Per)
	})

	t.Run("1000/m", func(t *testing.T) {
		rate := vegeta.Rate{}
		rf := &RateFlag{Rate: &rate}
		err := rf.Set("1000/m")
		assert.Nil(err)
		assert.Equal(1000, rate.Freq)
		assert.Equal(time.Minute, rate.Per)
	})

	t.Run("1000/2m", func(t *testing.T) {
		rate := vegeta.Rate{}
		rf := &RateFlag{Rate: &rate}
		err := rf.Set("1000/2m")
		assert.Nil(err)
		assert.Equal(1000, rate.Freq)
		assert.Equal(2*time.Minute, rate.Per)
	})

}

func TestRateFlagString(t *testing.T) {
	assert := assert.New(t)

	t.Run("100/s", func(t *testing.T) {
		rf := &RateFlag{Rate: &vegeta.Rate{Per: time.Second, Freq: 100}}
		assert.Equal("100/1s", rf.String())
	})

	t.Run("", func(t *testing.T) {
		rf := &RateFlag{}
		assert.Equal("", rf.String())
	})
}
