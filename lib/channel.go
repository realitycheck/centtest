package centtest

import (
	"fmt"
)

type Channel struct {
	name string
	u    *User
}

func NewChannel(name string) *Channel {
	return &Channel{name: name}
}

func (ch *Channel) String() string {
	if ch.u != nil {
		return fmt.Sprintf("%s#%s", ch.name, ch.u.ID)
	}
	return fmt.Sprintf("%s", ch.name)
}

func (ch *Channel) Attach(u *User) *Channel {
	if ch.u != nil {
		panic(fmt.Sprintf("Channel %s is already attached to user %s", ch.name, ch.u))
	}
	ch.u = u
	return ch
}
