package cloudflare

import (
	"context"
	"testing"
)

func TestDNS_List(t *testing.T) {
	d := DNS{token: "xx"}
	ls, err := d.List(context.Background(), "xx", Parms{Name: "*.bysir.top"})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", ls)
}

func TestUpdateRecordByName(t *testing.T) {
	d := DNS{token: "xx"}

	err := d.UpdateRecordByName(context.Background(), "xx", "*.bysir.top", "127.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ok")
}
