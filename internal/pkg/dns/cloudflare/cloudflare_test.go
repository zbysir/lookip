package cloudflare

import (
	"context"
	"testing"
)

func TestDNS_List(t *testing.T) {
	d := CF{token: "x"}
	ls, err := d.List(context.Background(), "d2c71015294ef5631aa82abb126fb1eb", Params{Name: "*.bysir.top"})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", ls)
}

func TestUpdateRecordByName(t *testing.T) {
	d := CF{
		token: "x",
		zone:  "d2c71015294ef5631aa82abb126fb1eb",
		name:  "*.bysir.top",
	}

	err := d.UpdateRecord(context.Background(), "2.0.0.1", false)
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ok")
}
