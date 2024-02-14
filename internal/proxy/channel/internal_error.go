package channel

var InternalError = &internalError{
	channel: make(chan error, 1),
}

type internalError struct {
	channel chan error
}

func (ie *internalError) Send(err error) {
	ie.channel <- err
}

func (ie *internalError) Listen() <-chan error {
	return ie.channel
}
