package docker

import (
	"context"
	"errors"
	"fmt"
	"net"
	"strings"
	"time"

	"github.com/dbarzdys/docktest/config"
)

type Waiter interface {
	Wait(services map[string]config.Service, containers []Container) error
}

func NewWaiter() Waiter {
	return waiter{}
}

type waiter struct{}

func (w waiter) Wait(
	services map[string]config.Service,
	containers []Container,
) error {
	wl := makeWaitList(time.Second * 100)
	mp := make(map[string]Container)
	for _, c := range containers {
		mp[c.Name] = c
	}
	for name, svc := range services {
		c, ok := mp[name]
		if !ok {
			return errors.New("Could not find container")
		} else if svc.WaitOn != "" {
			wl.Add("tcp", fmt.Sprintf("%s:%s", c.IP, svc.WaitOn))
		}
	}
	return wl.Wait()
}

// WaitList - waits for open connections
type WaitList interface {
	Add(network, address string) WaitList
	Wait() error
}

type waitTarget struct {
	network, address string
}

type waitList struct {
	doneMap map[waitTarget]error
	done    int
	timeout time.Duration
}

// makeWaitList - creates a new wait list with given timeout
func makeWaitList(timeout time.Duration) WaitList {
	return &waitList{
		doneMap: make(map[waitTarget]error, 0),
		done:    0,
		timeout: timeout,
	}
}

func (wl *waitList) Add(network, address string) WaitList {
	target := waitTarget{network, address}
	wl.doneMap[target] = errors.New("Have not been dialed")
	return wl
}

func (wl *waitList) Wait() (err error) {
	if len(wl.doneMap) == 0 {
		return nil
	}
	ctx := context.Background()
	ctx, cancel := context.WithTimeout(ctx, wl.timeout)

	defer cancel()
	type doneRes struct {
		target waitTarget
		err    error
	}
	ch := make(chan doneRes)
	for target := range wl.doneMap {
		go func(target waitTarget, ch chan<- doneRes) {
			err := wl.loopDial(ctx, target)
			ch <- doneRes{target, err}
		}(target, ch)
	}
	length := len(wl.doneMap)
	count := 0
	for done := range ch {
		wl.doneMap[done.target] = done.err
		count++
		if count == length {
			close(ch)
		}
	}
	return wl.formatError()
}

func (wl *waitList) formatError() error {
	list := []string{}
	for t, err := range wl.doneMap {
		if err == nil {
			continue
		}
		str := fmt.Sprintf("error connecting to %s://%s:[%v]", t.network, t.address, err)
		list = append(list, str)
	}
	if len(list) > 0 {
		return errors.New(strings.Join(list, ", "))
	}
	return nil
}

func (wl *waitList) loopDial(ctx context.Context, target waitTarget) (err error) {
loop:
	for {
		select {
		case <-ctx.Done():
			break loop
		default:
			err = wl.dial(target)
			if err == nil {
				return nil
			}
			time.Sleep(time.Second)
		}
	}
	return
}

func (wl *waitList) dial(target waitTarget) error {
	_, err := net.DialTimeout(target.network, target.address, time.Second)
	return err
}
