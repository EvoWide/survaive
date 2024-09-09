package sse

type Subscriber struct {
	uid      string
	dataChan chan string
}

func NewSubscriber(uid string, dataChan chan string) *Subscriber {
	return &Subscriber{
		uid,
		dataChan,
	}
}
