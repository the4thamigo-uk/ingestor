package interrupter

import (
	"os"
	"os/signal"
	"sync"
)

// Closer is the interface used to shutdown the interrupter instance
type Closer interface {
	Close()
}

type interrupter struct {
	c chan os.Signal
	w sync.WaitGroup
}

// New creates a new interrupter instance to handle process interrupts. On interruption the handler is executed.
func New(handler func()) Closer {
	i := interrupter{
		c: make(chan os.Signal),
	}
	signal.Notify(i.c, os.Interrupt)
	i.w.Add(1)
	go func() {
		defer i.w.Done()
		if sig := <-i.c; sig == os.Interrupt {
			handler()
		}
	}()
	return &i
}

// Close terminates interrupter instance cleanly. The handler is not executed in this case.
func (i *interrupter) Close() {
	close(i.c)
	i.w.Wait()
}
