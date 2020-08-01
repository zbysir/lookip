package public_ip

import "testing"

func TestGet(t *testing.T) {
	id, err := Get()
	if err != nil {
		t.Fatal(err)
	}

	t.Log(id)
}
