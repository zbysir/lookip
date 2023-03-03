package cloudflare

import (
	"context"
	"testing"
)

func TestDNS_List(t *testing.T) {
	d := DNS{token: "wYktacHe18WHwd_9M5hw7vpT7dRIStpXfGdJQp2T"}
	ls, err := d.List(context.Background(), "d2c71015294ef5631aa82abb126fb1eb", Parms{Name: "*.bysir.top"})
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("%+v", ls)
}

func TestUpdateRecordByName(t *testing.T) {
	d := DNS{
		token: "wYktacHe18WHwd_9M5hw7vpT7dRIStpXfGdJQp2T",
		zone:  "d2c71015294ef5631aa82abb126fb1eb",
		name:  "*.bysir.top",
	}

	err := d.UpdateRecord(context.Background(), "2.0.0.1")
	if err != nil {
		t.Fatal(err)
	}

	t.Logf("ok")
}
