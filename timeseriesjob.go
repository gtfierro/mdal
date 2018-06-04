package main

type worker struct {
	localWork chan dataRequest
	quit      chan bool
	b         *btrdbClient
}

func newWorker(b *btrdbClient) worker {
	return worker{
		localWork: make(chan dataRequest),
		quit:      make(chan bool),
		b:         b,
	}
}

func (w worker) start() {
	go func() {
		for {
			w.b.workerpool <- w.localWork
			select {
			case req := <-w.localWork:
				if err := w.b.handleRequest(req); err != nil {
					log.Error(err, req)
					req.errs <- err
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w worker) stop() {
	w.quit <- true
}
