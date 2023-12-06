package balancer

import "github.com/lafikl/consistent"

func init() {
	factories[BoundedBalancer] = NewBounded
}

type Bounded struct {
	ch *consistent.Consistent
}

func NewBounded(hosts []string) Balancer {
	c := &Bounded{consistent.New()}
	for _, h := range hosts {
		c.ch.Add(h)
	}
	return c
}

func (b *Bounded) Add(host string) {
	b.ch.Add(host)
}

func (b *Bounded) Remove(host string) {
	b.ch.Remove(host)
}

func (b *Bounded) Balance(key string) (string, error) {
	if len(b.ch.Hosts()) == 0 {
		return "", ErrNoHost
	}
	return b.ch.GetLeast(key)
}

func (b *Bounded) Inc(host string) {
	b.ch.Inc(host)
}

func (b *Bounded) Done(host string) {
	b.ch.Done(host)
}
