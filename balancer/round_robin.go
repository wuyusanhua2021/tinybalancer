package balancer

import "sync/atomic"

func init() {
	factories[R2Balancer] = NewRoundRobin
}

type RoundRobin struct {
	BaseBalancer
	i atomic.Uint64
}

func NewRoundRobin(hosts []string) Balancer {
	return &RoundRobin{
		i: atomic.Uint64{},
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
	}
}

func (r *RoundRobin) Balance(_ string) (string, error) {
	r.RLock()
	defer r.RUnlock()

	if len(r.hosts) == 0 {
		return "", ErrNoHost
	}
	hosts := r.hosts[r.i.Add(1)%uint64(len(r.hosts))]
	return hosts, nil
}
