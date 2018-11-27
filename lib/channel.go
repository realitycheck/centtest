package centtest

import (
	"fmt"
)

type Channel string

func NewChannel(num int) Channel {
	return Channel(fmt.Sprintf("CH.%00000000d", num))
}

func (ch Channel) String() string {
	return string(ch)
}

func (ch Channel) Attach(u *User) Channel {
	return Channel(fmt.Sprintf("%s#%s", ch, u.ID))
}
