package balancer

import "hash/crc32"

func init() {
	factories[IPHashBalancer] = NewIPHash
}

type IPHash struct {
	BaseBalancer
}

func NewIPHash(hosts []string) Balancer {
	return &IPHash{
		BaseBalancer: BaseBalancer{
			hosts: hosts,
		},
	}
}

func (ih *IPHash) Balance(key string) (string, error) {
	ih.RLock()
	defer ih.RUnlock()

	if len(ih.hosts) == 0 {
		return "", ErrNoHost
	}
	value := crc32.ChecksumIEEE([]byte(key)) % uint32(len(ih.hosts))
	return ih.hosts[value], nil
}
