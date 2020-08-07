package public_ip

import (
	"strings"
)

type x3322 struct {
	client HttpClient
}

func (h x3322) Name() string {
	return "http://ip.3322.net"
}

func (h x3322) Ip() (ip string, err error) {
	bs, err := h.client.Get("http://ip.3322.net")
	if err != nil {
		return
	}

	ip = strings.TrimSpace(string(bs))

	return
}

func New3322() IpGetter {
	return x3322{
		client: NewClient(),
	}
}
