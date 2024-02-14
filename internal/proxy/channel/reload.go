package channel

var Reload = &reload{
	channel: make(chan bool, 1),
}

type reload struct {
	channel chan bool
}

func (r *reload) Send() {
	r.channel <- true
}

func (r *reload) Listen() <-chan bool {
	return r.channel
}
