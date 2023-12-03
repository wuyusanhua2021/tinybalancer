package balancer

import "errors"

var (
	ErrNoHost              = errors.New("no host")
	ErrAlgorithmNotSupport = errors.New("algorithm not support")
)

type Balancer interface {
	Add(string)
	Remove(string)
	Balance(string) (string, error)
	Inc(string)
	Done(string)
}

type Factory func([]string) Balancer

var factories = make(map[string]Factory)

func Build(algorithm string, hosts []string) (Balancer, error) {
	factory, ok := factories[algorithm]
	if !ok {
		return nil, ErrAlgorithmNotSupport
	}
	return factory(hosts), nil
}
