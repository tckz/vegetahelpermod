package vegetahelper

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sort"
	"sync/atomic"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	vegeta "github.com/tsenart/vegeta/lib"
)

func TestAttack(t *testing.T) {
	assert := assert.New(t)

	called := uint64(0)
	atk := NewAttacker(func(context.Context) (*HitResult, error) {
		c := atomic.AddUint64(&called, 1)
		return &HitResult{
			RecvBytes: c,
			SentBytes: 456,
			Code:      http.StatusOK,
			Error:     "",
		}, nil
	}, WithWorkers(2))

	ctx := context.Background()

	ch := atk.Attack(ctx, vegeta.Rate{10, time.Second}, 1*time.Second, "attack!")
	res := []*vegeta.Result{}
	for r := range ch {
		res = append(res, r)
	}
	sort.SliceStable(res, func(i, j int) bool {
		return res[i].BytesIn < res[j].BytesIn
	})

	l := len(res)
	assert.True(9 <= l && l <= 11, fmt.Sprintf("len should arround 20(actual=%d)", l))

	for i, e := range res {
		if uint64(i+1) != e.BytesIn {
			t.Errorf("res does not ordered correctly at index=%d, BytesIn=%d", i, e.BytesIn)
			break
		}

		if e.Attack != "attack!" {
			t.Errorf("actual=%s, index=%d", e.Attack, i)
			break
		}
	}
}

func TestAttackCancel(t *testing.T) {
	assert := assert.New(t)

	called := uint64(0)
	atk := NewAttacker(func(context.Context) (*HitResult, error) {
		c := atomic.AddUint64(&called, 1)
		return &HitResult{
			RecvBytes: c,
			SentBytes: 456,
			Code:      http.StatusOK,
			Error:     "",
		}, nil
	}, WithWorkers(10))

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go func() {
		time.Sleep(500 * time.Millisecond)
		cancel()
	}()
	ch := atk.Attack(ctx, vegeta.Rate{10, time.Second}, 1*time.Second, "attack!")

	res := []*vegeta.Result{}
	for r := range ch {
		res = append(res, r)
	}

	l := len(res)
	assert.True(4 <= l && l <= 9, fmt.Sprintf("len should arround 5(actual=%d)", l))
}

func TestHit(t *testing.T) {
	assert := assert.New(t)

	atk := NewAttacker(func(context.Context) (*HitResult, error) {
		return &HitResult{
			RecvBytes: 123,
			SentBytes: 456,
			Code:      http.StatusInternalServerError,
			Error:     "wao",
		}, nil
	})

	ctx := context.Background()
	t.Run("Some error", func(t *testing.T) {
		res := atk.hit(ctx, "hittest")
		assert.NotNil(res)
		assert.Equal(uint16(500), res.Code)
		assert.Equal(uint64(456), res.BytesOut)
		assert.Equal(uint64(123), res.BytesIn)
		assert.Equal("wao", res.Error)
		assert.Equal(uint64(0), res.Seq)
		assert.Equal("hittest", res.Attack)
	})

	t.Run("2nd", func(t *testing.T) {
		res := atk.hit(ctx, "hittest2")
		assert.NotNil(res)
		assert.Equal(uint16(500), res.Code)
		assert.Equal(uint64(456), res.BytesOut)
		assert.Equal(uint64(123), res.BytesIn)
		assert.Equal("wao", res.Error)
		assert.Equal(uint64(1), res.Seq)
		assert.Equal("hittest2", res.Attack)
	})
}

func TestHit2(t *testing.T) {
	assert := assert.New(t)

	atk := NewAttacker(func(context.Context) (*HitResult, error) {
		// something went wrong
		return nil, errors.New("this is error")
	})

	ctx := context.Background()
	t.Run("Some error", func(t *testing.T) {
		res := atk.hit(ctx, "hittest")
		assert.NotNil(res)
		assert.Equal(uint16(500), res.Code)
		assert.Equal(uint64(0), res.BytesOut)
		assert.Equal(uint64(0), res.BytesIn)
		assert.Equal("this is error", res.Error)
		assert.Equal(uint64(0), res.Seq)
		assert.Equal("hittest", res.Attack)
	})
}

func TestHit3(t *testing.T) {
	assert := assert.New(t)

	atk := NewAttacker(func(context.Context) (*HitResult, error) {
		// Assume as succ.
		return nil, nil
	})

	ctx := context.Background()
	t.Run("Assume as succ", func(t *testing.T) {
		res := atk.hit(ctx, "hittest")
		assert.NotNil(res)
		assert.Equal(uint16(200), res.Code)
		assert.Equal(uint64(0), res.BytesOut)
		assert.Equal(uint64(0), res.BytesIn)
		assert.Equal("", res.Error)
		assert.Equal(uint64(0), res.Seq)
		assert.Equal("hittest", res.Attack)
	})
}
