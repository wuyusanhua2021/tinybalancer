package balancer

import (
	"hash/crc32"
	"math/rand"
	"sync"
	"time"
)

func init() {
	factories[P2CBalancer] = NewP2C
}

const Salt = "%#!"

type host struct {
	name string
	load uint64
}

type P2C struct {
	sync.RWMutex
	hosts   []*host
	rnd     rand.Rand
	loadMap map[string]*host
}

func NewP2C(hosts []string) Balancer {
	p := &P2C{
		hosts:   []*host{},
		loadMap: make(map[string]*host),
		rnd:     *rand.New(rand.NewSource(time.Now().UnixNano())),
	}

	for _, h := range hosts {
		p.Add(h)
	}
	return p
}

func (p *P2C) Add(hostName string) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.loadMap[hostName]; ok {
		return
	}

	h := &host{name: hostName, load: 0}
	p.hosts = append(p.hosts, h)
	p.loadMap[hostName] = h
}

func (p *P2C) Remove(host string) {
	p.Lock()
	defer p.Unlock()

	if _, ok := p.loadMap[host]; !ok {
		return
	}

	delete(p.loadMap, host)

	for i, h := range p.hosts {
		if h.name == host {
			p.hosts = append(p.hosts[:i], p.hosts[i+1:]...)
			return
		}
	}
}

func (p *P2C) Balance(key string) (string, error) {
	p.RLock()
	defer p.RUnlock()

	if len(p.hosts) == 0 {
		return "", ErrNoHost
	}

	n1, n2 := p.hash(key)
	host := n2
	if p.loadMap[n1].load <= p.loadMap[n2].load {
		host = n1
	}
	return host, nil
}

func (p *P2C) hash(key string) (n1 string, n2 string) {
	if len(key) > 0 {
		salyKey := key + Salt
		n1 = p.hosts[crc32.ChecksumIEEE([]byte(key))%uint32(len(p.hosts))].name
		n2 = p.hosts[crc32.ChecksumIEEE([]byte(salyKey))%uint32(len(p.hosts))].name
		return
	}

	n1 = p.hosts[p.rnd.Intn(len(p.hosts))].name
	n2 = p.hosts[p.rnd.Intn(len(p.hosts))].name
	return
}

func (p *P2C) Inc(host string) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return
	}
	h.load++
}

func (p *P2C) Done(host string) {
	p.Lock()
	defer p.Unlock()

	h, ok := p.loadMap[host]
	if !ok {
		return
	}
	if h.load > 0 {
		h.load--
	}
}
