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

type DNS struct {
	token string
	zone  string
	name  string
}

func NewDNS(token string, zoneId string, name string) *DNS {
	if token == "" {
		log.Panicf("cf_tokne not set")
	}
	if zoneId == "" {
		log.Panicf("cf_zone_id not set")
	}
	return &DNS{
		token: token,
		zone:  zoneId,
		name:  name,
	}
}

func (d *DNS) Update(ctx context.Context, r Record) error {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r)
	if err != nil {
		return err
	}

	log.Printf("update record: %s", buf.Bytes())

	req, err := http.NewRequest("PUT", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records/%s", r.ZoneId, r.Id), &buf)
	if err != nil {
		return err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}
	defer rsp.Body.Close()

	rspBody, _ := ioutil.ReadAll(rsp.Body)

	if rsp.StatusCode != 200 {
		return errors.New(fmt.Sprintf("%s", rspBody))
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

type Parms struct {
	Name string `json:"name"`
}

func (d *DNS) List(ctx context.Context, zone string, p Parms) (rs []Record, err error) {
	query := ""
	if p.Name != "" {
		query = fmt.Sprintf("?name=%s", p.Name)
	}

	req, err := http.NewRequest("GET", fmt.Sprintf("https://api.cloudflare.com/client/v4/zones/%s/dns_records%s", zone, query), nil)
	if err != nil {
		return rs, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", d.token))
	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return rs, err
	}

	defer rsp.Body.Close()

	// {"result":[{"id":"34536642eee9956c6c73931bf1d6e238","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"editor.duckment.top","type":"A","content":"23.225.169.129","proxiable":true,"proxied":true,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2022-10-08T04:58:54.20771Z","modified_on":"2022-10-08T04:58:54.20771Z"},{"id":"0e5527ae305c131c8754a081b2d1c215","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"t.duckment.top","type":"A","content":"47.75.56.10","proxiable":true,"proxied":true,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2020-09-22T06:15:14.005604Z","modified_on":"2020-09-22T06:15:14.005604Z"},{"id":"97c90a56383bd7c376af1b6ec68f884f","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"dmtrace.duckment.top","type":"CNAME","content":"tracedm.aliyuncs.com","proxiable":true,"proxied":false,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2020-12-23T03:08:35.685329Z","modified_on":"2020-12-23T03:08:35.685329Z"},{"id":"80770bae578189ce6743735a0bc935cb","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"duckment.top","type":"CNAME","content":"x.bysir.top","proxiable":true,"proxied":true,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2021-03-29T14:25:30.384439Z","modified_on":"2021-03-29T14:25:30.384439Z"},{"id":"f313013917765e5896b844054c95d21e","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"preview.editor.duckment.top","type":"CNAME","content":"editor.duckment.top","proxiable":true,"proxied":false,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2022-10-08T06:16:09.283858Z","modified_on":"2022-10-08T06:16:09.283858Z"},{"id":"5550f20bc2fd82f98d6c8d54d9ff3650","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"zz.duckment.top","type":"CNAME","content":"bysir.e.cn.vc","proxiable":true,"proxied":true,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2020-09-22T06:04:31.327641Z","modified_on":"2020-09-22T06:04:31.327641Z"},{"id":"1fbeefd8d96908887a4836057c0f1d14","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"duckment.top","type":"MX","content":"mx01.dm.aliyun.com","priority":10,"proxiable":false,"proxied":false,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2020-12-23T03:08:03.755146Z","modified_on":"2020-12-23T03:08:03.755146Z"},{"id":"e0d264aef82f70e339d9c5dbba225f9d","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"aliyundm.duckment.top","type":"TXT","content":"a6a0e2477fca42deb04d","proxiable":false,"proxied":false,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2020-12-23T03:07:24.701024Z","modified_on":"2020-12-23T03:07:24.701024Z"},{"id":"3d97b71b0ed2a5087d2f029414551c76","zone_id":"cd3729b108a514b29c84e917c1d9a36f","zone_name":"duckment.top","name":"duckment.top","type":"TXT","content":"v=spf1 include:spf1.dm.aliyun.com -all","proxiable":false,"proxied":false,"ttl":1,"locked":false,"meta":{"auto_added":false,"managed_by_apps":false,"managed_by_argo_tunnel":false,"source":"primary"},"comment":null,"tags":[],"created_on":"2020-12-23T03:07:38.323484Z","modified_on":"2020-12-23T03:07:38.323484Z"}],"success":true,"errors":[],"messages":[],"result_info":{"page":1,"per_page":100,"count":9,"total_count":9,"total_pages":1}}

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
	rspBody, _ := ioutil.ReadAll(rsp.Body)
	if rsp.StatusCode != 200 {
		return rs, errors.New(fmt.Sprintf("%s", rspBody))
	}

	var r Rsp
	err = json.Unmarshal(rspBody, &r)
	if err != nil {
		return rs, err
	}

	return r.Result, nil
}

func (d *DNS) GetRecord(ctx context.Context) (id Record, exist bool, err error) {
	rs, err := d.List(ctx, d.zone, Parms{Name: d.name})
	if err != nil {
		return id, false, err
	}

	if len(rs) == 0 {
		return id, false, nil
	}

	return rs[0], true, nil
}

func (d *DNS) UpdateRecord(ctx context.Context, content string) (err error) {
	id, exist, err := d.GetRecord(ctx)
	if err != nil {
		return err
	}

	if !exist {
		return errors.New("record not exist")
	}

	proxyi := ctx.Value("proxy")
	proxy := false
	if proxyi != nil {
		proxy = proxyi.(bool)
	}

	id.Content = content
	id.Proxied = proxy
	return d.Update(ctx, id)
}
