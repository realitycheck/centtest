package centtest

import (
	"log"
	"sync"

	"github.com/centrifugal/centrifuge-go"

	"github.com/dgrijalva/jwt-go"
)

var (
	jwtAlgorithm = jwt.SigningMethodHS256
	jwtKey       = []byte("secret")
)

type User struct {
	ID string
}

func NewUser(g IDGenerator) *User {
	u := &User{
		ID: g.ID(),
	}
	return u
}

func (u *User) String() string {
	return u.ID
}

func (u *User) Token(exp int64) string {
	claims := jwt.MapClaims{
		"sub": u.ID,
		"info": map[string]string{
			"name": u.ID,
		},
	}
	if exp > 0 {
		claims["exp"] = exp
	}
	t, err := jwt.NewWithClaims(jwtAlgorithm, claims).SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	return t
}

type receiver struct {
	wg *sync.WaitGroup
	u  *User
	c  *Client
	ch Channel
}

func (r *receiver) OnPublish(sub *centrifuge.Subscription, e centrifuge.PublishEvent) {
	r.wg.Done()
}

func (r *receiver) OnUnsubscribe(sub *centrifuge.Subscription, e centrifuge.UnsubscribeEvent) {
	r.wg.Done()
}

func (r *receiver) OnConnect(cli *centrifuge.Client, e centrifuge.ConnectEvent) {
	r.c.connected = true
	r.c.id = e.ClientID
	log.Printf("centtest: Connected, user=%s, client=%s", r.u, r.c)
}

func (r *receiver) OnDisconnect(cli *centrifuge.Client, e centrifuge.DisconnectEvent) {
	if !e.Reconnect {
		defer r.wg.Done()
	}

	if r.c.connected {
		log.Printf("centtest: Disconnected, user=%s, client=%s, reason=%s", r.u, r.c, e.Reason)
	} else {
		log.Printf("centtest: Can't connect, user=%s, reason=%s", r.u, e.Reason)
	}
	r.c.connected = false
}

func (r *receiver) OnSubscribeSuccess(sub *centrifuge.Subscription, e centrifuge.SubscribeSuccessEvent) {
	log.Printf("centtest: Subscribed, channel=%s, user=%s, client=%s", r.ch, r.u, r.c)
}

func (r *receiver) OnSubscribeError(sub *centrifuge.Subscription, e centrifuge.SubscribeErrorEvent) {
	log.Printf("centtest: Can't subscribe, channel=%s, user=%s, client=%s, reason=%s", r.ch, r.u, r.c, e.Error)
}

func (r *receiver) run() {
	defer r.wg.Done()

	r.c.cli.OnConnect(r)
	r.c.cli.OnDisconnect(r)
	r.c.cli.SetToken(r.u.Token(0))

	err := r.c.cli.Connect()
	if err != nil {
		log.Printf("centtest: Can't connect, user=%s, reason=%s", r.u.ID, err)
	}

	sub, err := r.c.cli.NewSubscription(r.ch.String())
	sub.OnPublish(r)
	sub.OnSubscribeSuccess(r)
	sub.OnSubscribeError(r)
	sub.OnUnsubscribe(r)

	err = sub.Subscribe()
	if err != nil {
		log.Printf("centtest: Can't subscribe, user=%s, client=%s, channel=%s, reason=%s",
			r.u.ID, r.c.id, r.ch, err)
	}

	r.wg.Add(1)
}

func Run(u *User, c *Client, ch Channel, wg *sync.WaitGroup) {
	(&receiver{wg, u, c, ch.Attach(u)}).run()
}
