package public_ip

import (
	"crypto/tls"
	"io/ioutil"
	"net/http"
	"time"
)

type HttpClient struct {
	client *http.Client
}

func NewClient() HttpClient {
	transport := &http.Transport{
		TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
		DisableCompression: true,
	}
	client := &http.Client{Transport: transport, Timeout: 30 * time.Second}

	return HttpClient{
		client: client,
	}
}

func (c HttpClient) Get(url string) (bs []byte, err error) {
	r, err := c.client.Get(url)
	if err != nil {
		return
	}
	defer r.Body.Close()
	bs, err = ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	return
}
