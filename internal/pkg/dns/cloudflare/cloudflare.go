package cloudflare

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"time"
)

type CF struct {
	token string
	zone  string
	name  string
}

func NewCF(token string, zoneId string, name string) *CF {
	if token == "" {
		log.Panicf("cf_tokne not set")
	}
	if zoneId == "" {
		log.Panicf("cf_zone_id not set")
	}
	return &CF{
		token: token,
		zone:  zoneId,
		name:  name,
	}
}

func (d *CF) Update(ctx context.Context, r Record) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r)
	if err != nil {
		return err
	}

	log.Printf("update record: %s", buf.Bytes())

	err = d.req(ctx, "PUT", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", r.ZoneId, r.Id), r, nil)
	if err != nil {
		return err
	}

	return nil
}

type Record struct {
	Id        string `json:"id"`
	ZoneId    string `json:"zone_id"`
	ZoneName  string `json:"zone_name"`
	Name      string `json:"name"`
	Type      string `json:"type"`
	Content   string `json:"content"`
	Proxiable bool   `json:"proxiable"`
	Proxied   bool   `json:"proxied"`
	Ttl       int    `json:"ttl"`
	Locked    bool   `json:"locked"`
	Meta      struct {
		AutoAdded           bool   `json:"auto_added"`
		ManagedByApps       bool   `json:"managed_by_apps"`
		ManagedByArgoTunnel bool   `json:"managed_by_argo_tunnel"`
		Source              string `json:"source"`
	} `json:"meta"`
	Comment    string        `json:"comment"`
	Tags       []interface{} `json:"tags"`
	CreatedOn  time.Time     `json:"created_on"`
	ModifiedOn time.Time     `json:"modified_on"`
	Priority   int           `json:"priority,omitempty"`
}

type Params struct {
	Name string `json:"name"`
}

func (d *CF) req(ctx context.Context, method string, url string, body interface{}, rsp interface{}) (err error) {
	var buf *bytes.Buffer
	if body != nil {
		buf = &bytes.Buffer{}
		err := json.NewEncoder(buf).Encode(body)
		if err != nil {
			return err
		}
	}

	req, err := http.NewRequestWithContext(ctx, method, url, buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))
	rspr, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rspr.Body.Close()
	rspBody, _ := ioutil.ReadAll(rspr.Body)
	if rspr.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%s", rspBody))
	}

	if rsp != nil {
		err = json.Unmarshal(rspBody, &rsp)
		if err != nil {
			return err
		}
	}
	return nil
}

func (d *CF) List(ctx context.Context, zone string, p Params) (rs []Record, err error) {
	query := ""
	if p.Name != "" {
		query = fmt.Sprintf("?name=%s", p.Name)
	}

	type Rsp struct {
		Result     []Record      `json:"result"`
		Success    bool          `json:"success"`
		Errors     []interface{} `json:"errors"`
		Messages   []interface{} `json:"messages"`
		ResultInfo struct {
			Page       int `json:"page"`
			PerPage    int `json:"per_page"`
			Count      int `json:"count"`
			TotalCount int `json:"total_count"`
			TotalPages int `json:"total_pages"`
		} `json:"result_info"`
	}

	var r Rsp
	err = d.req(ctx, "GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records%s", zone, query), nil, &r)
	if err != nil {
		return rs, err
	}

	return r.Result, nil
}

func (d *CF) GetRecord(ctx context.Context) (id Record, exist bool, err error) {
	log.Printf("cloudflare GetRecord")

	rs, err := d.List(ctx, d.zone, Params{Name: d.name})
	if err != nil {
		return id, false, err
	}

	if len(rs) == 0 {
		return id, false, nil
	}

	return rs[0], true, nil
}

func (d *CF) UpdateRecord(ctx context.Context, content string, force bool) (aff bool, err error) {
	log.Printf("cloudflare UpdateRecord")
	id, exist, err := d.GetRecord(ctx)
	if err != nil {
		return false, err
	}

	if !exist {
		return false, errors.New("record not exist")
	}

	if !force && id.Content == content {
		log.Printf("record not change")
		return
	}

	proxy := id.Proxied

	proxyi := ctx.Value("proxy")
	if proxyi != nil {
		proxy = proxyi.(bool)
	}

	id.Content = content
	id.Proxied = proxy
	return true, d.Update(ctx, id)
}
