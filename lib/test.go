package centtest

import (
	"fmt"
	"log"

	centrifuge "github.com/centrifugal/centrifuge-go"
)

type Test struct {
	u  *User
	c  *Client
	ch *Channel

	debug bool
}

func NewTest(u *User, c *Client, ch *Channel, debug bool) *Test {
	return &Test{u, c, ch, debug}
}

func (t *Test) String() string {
	return fmt.Sprintf("user=%s, client=%s, channel=%s", t.u, t.c, t.ch)
}

func (t *Test) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	t.log("centtest: Published, %s, data=%s", t, e.Data)
}

func (t *Test) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	t.logdebug("centtest: Unsubscribed, %s", t)
}

func (t *Test) OnConnect(cli *centrifuge.Client, e centrifuge.ConnectEvent) {
	t.c.connected = true
	t.c.id = e.ClientID
	t.logdebug("centtest: Connected, %s", t)
}

func (t *Test) OnDisconnect(cli *centrifuge.Client, e centrifuge.DisconnectEvent) {
	if t.c.connected {
		t.logdebug("centtest: Disconnected, %s, reason=%s", t, e.Reason)
	} else {
		t.log("centtest: Can't connect, %s, reason=%s", t, e.Reason)
	}
	t.c.connected = false
}

func (t *Test) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	t.logdebug("centtest: Subscribed, %s", t)
}

func (t *Test) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	t.log("centtest: Can't subscribe, %s, reason=%s", t, e.Error)
}

func (t *Test) OnError(cli *centrifuge.Client, e centrifuge.ErrorEvent) {
	t.log("centest: Error, %s, reason=%s", t, e.Message)
}

func (t *Test) Run() {
	t.c.cli.OnConnect(t)
	t.c.cli.OnDisconnect(t)
	t.c.cli.OnError(t)
	t.c.cli.SetToken(t.u.token(0))

	err := t.c.cli.Connect()
	if err != nil {
		t.log("centtest: Connect error, %s, reason=%s", t, err)
	}

	sub, err := t.c.cli.NewSubscription(t.ch.String())
	sub.OnPublish(t)
	sub.OnSubscribeSuccess(t)
	sub.OnSubscribeError(t)
	sub.OnUnsubscribe(t)

	err = sub.Subscribe()
	if err != nil {
		t.log("centtest: Subscribe error, %s, reason=%s", t, err)
	}
}

func (t *Test) Close() {
	if t.c.connected {
		if err := t.c.disconnect(); err != nil {
			t.log("centtest: Disconnect error, %s, reason=%s", t, err)
		}
	}
}

func (t *Test) log(format string, v ...interface{}) {
	log.Printf(format, v...)
}

func (t *Test) logdebug(format string, v ...interface{}) {
	if t.debug {
		t.log(format, v...)
	}
}
