package discover

var discoverMap map[string]Discover

func GetDiscover(name string) Discover {
	if discover, ok := discoverMap[name]; ok {
		return discover
	}
	return Consul{}
}

func DiscoverRegister(name string, discover Discover) {
	discoverMap[name] = discover
}

type Discover interface {
}

type Consul struct {
}
