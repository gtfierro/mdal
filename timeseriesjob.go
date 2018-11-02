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

type worker2 struct {
	localWork chan dataRequest
	quit      chan bool
	b         *btrdbClient
}

func newWorker2(b *btrdbClient) worker2 {
	return worker2{
		localWork: make(chan dataRequest),
		quit:      make(chan bool),
		b:         b,
	}
}

func (w worker2) start() {
	go func() {
		for {
			w.b.workerpool2 <- w.localWork
			select {
			case req := <-w.localWork:
				if err := w.b.handleRequest2(req); err != nil {
					log.Error(err, req)
					req.errs <- err
				}
			case <-w.quit:
				return
			}
		}
	}()
}

func (w worker2) stop() {
	w.quit <- true
}
