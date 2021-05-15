package mailbox

import "github.com/okpub/dekopon/utils"

func NewMailbox() Mailbox {
	return &defaultMailbox{
		taskMailbox: utils.MakeTest(),
	}
}

// Unbounded returns a producer which creates an unbounded mailbox
func Unbounded() Producer {
	return func() Mailbox {
		return NewMailbox()
	}
}
