package kit

import (
	"context"
	"fmt"
	"time"
)

// Await allows awaiting some state by periodically hitting fn unless either it returns true or error or timeout
// It returns nil when fn results true
func Await(fn func() (bool, error), tick, timeout time.Duration) chan error {
	c := make(chan error)
	go func() {
		// first try without ticker
		res, err := fn()
		if err != nil {
			c <- err
			return
		}
		if res {
			c <- nil
			return
		}
		// if first try fails, run ticker
		ticker := time.NewTicker(tick)
		defer ticker.Stop()
		ctx, cancel := context.WithTimeout(context.Background(), timeout)
		defer cancel()
		for {
			select {
			case <-ticker.C:
				res, err := fn()
				if err != nil {
					c <- err
					return
				}
				if res {
					c <- nil
					return
				}
			case <-ctx.Done():
				c <- fmt.Errorf("timeout")
				return
			}
		}
	}()
	return c
}
