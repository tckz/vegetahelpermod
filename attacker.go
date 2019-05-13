/*
 * Based on
 * https://github.com/tsenart/vegeta/blob/master/lib/attack.go
 */

package vegetahelper

import (
	"context"
	"net/http"
	"sync"
	"time"

	vegeta "github.com/tsenart/vegeta/lib"
)

type HitResult struct {
	SentBytes uint64
	RecvBytes uint64
	Code      uint16
	Error     string
}

type HitFunc func(context.Context) (*HitResult, error)

type Attacker struct {
	hitter  HitFunc
	options *AttackerOptions
	seqmu   sync.Mutex
	seq     uint64
	began   time.Time
}

type AttackerOptions struct {
	numWorkers uint64
}

type AttackerOption func(*AttackerOptions)

func WithWorkers(n uint64) AttackerOption {
	return func(o *AttackerOptions) {
		o.numWorkers = n
	}
}

func NewAttacker(f HitFunc, opts ...AttackerOption) *Attacker {
	aopt := &AttackerOptions{
		numWorkers: vegeta.DefaultWorkers,
	}

	for _, o := range opts {
		o(aopt)
	}

	ret := &Attacker{
		hitter:  f,
		began:   time.Now(),
		options: aopt,
	}
	return ret
}

func (a *Attacker) Attack(ctx context.Context, r vegeta.Rate, du time.Duration, name string) <-chan *vegeta.Result {
	var wg sync.WaitGroup
	results := make(chan *vegeta.Result)
	ticks := make(chan uint64)
	for i := uint64(0); i < a.options.numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			a.attack(ctx, name, ticks, results)
		}()
	}

	go func() {
		defer close(results)
		defer wg.Wait()
		defer close(ticks)
		interval := uint64(r.Per.Nanoseconds() / int64(r.Freq))
		hits := uint64(du) / interval
		began, count := time.Now(), uint64(0)
		for {
			now, next := time.Now(), began.Add(time.Duration(count*interval))
			time.Sleep(next.Sub(now))
			select {
			case ticks <- count:
				if count++; count == hits {
					return
				}
			case <-ctx.Done():
				return
			default: // all workers are blocked. start one more and try again
				wg.Add(1)
				go func() {
					defer wg.Done()
					a.attack(ctx, name, ticks, results)
				}()
			}
		}
	}()

	return results
}

func (a *Attacker) attack(ctx context.Context, name string, ticks <-chan uint64, results chan<- *vegeta.Result) {
	for range ticks {
		results <- a.hit(ctx, name)
	}
}

func (a *Attacker) hit(ctx context.Context, name string) *vegeta.Result {
	res := vegeta.Result{Attack: name}

	a.seqmu.Lock()
	res.Timestamp = a.began.Add(time.Since(a.began))
	res.Seq = a.seq
	a.seq++
	a.seqmu.Unlock()

	result, err := a.hitter(ctx)
	res.Latency = time.Since(res.Timestamp)

	if result != nil {
		res.BytesIn = result.RecvBytes
		res.BytesOut = result.SentBytes
		res.Error = result.Error
		res.Code = result.Code
	} else {
		if err == nil {
			res.Code = http.StatusOK
		} else {
			res.Code = http.StatusInternalServerError
			res.Error = err.Error()
		}
	}

	return &res
}
