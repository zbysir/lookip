package public_ip

import (
	"encoding/json"
	"errors"
	"fmt"
	"strings"
)

type HttpbinIp struct {
	client HttpClient
}

func (h HttpbinIp) Name() string {
	return "http://httpbin.org/ip"
}

func (h HttpbinIp) Ip() (ip string, err error) {
	bs, err := h.client.Get("http://httpbin.org/ip")
	if err != nil {
		return
	}

	rsp := struct {
		Origin string `json:"origin"`
	}{}
	json.Unmarshal(bs, &rsp)

	if len(rsp.Origin) == 0 {
		err = errors.New(fmt.Sprintf("queryIp error:%s", bs))
		return
	}

	rsp.Origin = strings.TrimSpace(strings.Split(rsp.Origin, ",")[0])

	ip = rsp.Origin

	return
}

func NewHttpbinIp() IpGetter {
	return HttpbinIp{
		client: NewClient(),
	}
}
