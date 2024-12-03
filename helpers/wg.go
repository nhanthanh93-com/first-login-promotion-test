package helpers

import (
	"context"
	"sync"
)

type WgGroup struct {
	wg     *sync.WaitGroup
	err    error
	ctx    context.Context
	cancel context.CancelFunc
}

func NewWgGroup() *WgGroup {
	ctx, cancel := context.WithCancel(context.Background())
	return &WgGroup{
		wg:     new(sync.WaitGroup),
		ctx:    ctx,
		cancel: cancel,
	}
}

func (g *WgGroup) Wait() error {
	g.wg.Wait()
	return g.err
}

func (g *WgGroup) Go(f func() error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		g.err = f()
	}()
}

func (g *WgGroup) RunWithContext(f func(ctx context.Context) error) {
	g.wg.Add(1)
	go func() {
		defer g.wg.Done()
		g.err = f(g.ctx)
	}()
}

func (g *WgGroup) Cancel() {
	g.cancel()
}
