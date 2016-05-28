package rpcsv

import (
	"fmt"
	"time"
)

func TimeoutCoder(f func(interface{}) error, e interface{}, msg string) error {
	echan := make(chan error, 1)
	go func() { echan <- f(e) }()
	select {
	case e := <-echan:
		return e
	case <-time.After(time.Second * 5):
		return fmt.Errorf("Timeout %s", msg)
	}
}
