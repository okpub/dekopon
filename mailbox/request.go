package mailbox

import (
	"context"
	"sync/atomic"

	"github.com/okpub/dekopon/utils"
)

func NewRequest(request interface{}) *mailboxRequest {
	return &mailboxRequest{
		request: request,
		exit:    make(chan struct{}),
	}
}

//class
type mailboxRequest struct {
	request interface{}
	/**todo*/
	lock int32
	exit chan struct{}
	resp interface{}
}

//支持actor
func (req *mailboxRequest) Done() { req.Respond(utils.TempErr) }

//request
func (req *mailboxRequest) Message() interface{} { return req.request }

func (req *mailboxRequest) Respond(message interface{}) (err error) {
	if atomic.CompareAndSwapInt32(&req.lock, 0, 1) {
		req.resp = message
		close(req.exit)
	} else {
		err = utils.EOF
	}
	return
}

//respond
func (req *mailboxRequest) Body(ctx context.Context) (resp interface{}, err error) {
	select {
	case <-ctx.Done():
		err = ctx.Err()
	case <-req.exit:
		switch data := req.resp.(type) {
		case error:
			err = data
		case nil:
			err = utils.NilErr
		default:
			resp = data
		}
	}
	return
}
