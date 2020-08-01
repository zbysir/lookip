package public_ip

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var client *http.Client

func init() {
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client = &http.Client{Transport: transport, Timeout: 30 * time.Second}
}

func Get() (ip string, err error) {
	r, err := client.Get("http://httpbin.org/ip")
	if err != nil {
		return
	}
	defer r.Body.Close()
	bs, err := ioutil.ReadAll(r.Body)
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
