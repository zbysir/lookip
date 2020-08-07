package public_ip

import "testing"

func TestGet3322(t *testing.T) {
	x := Factory("3322")
	ip, err := x.Ip()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(x.Name(), ip)
}

func TestGethttpbin(t *testing.T) {
	x := Factory("httpbin")
	ip, err := x.Ip()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(x.Name(), ip)
}
