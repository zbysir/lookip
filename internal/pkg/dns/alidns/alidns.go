package alidns

import (
	"context"
	"errors"
	"github.com/aliyun/alibaba-cloud-sdk-go/sdk/requests"
	"github.com/aliyun/alibaba-cloud-sdk-go/services/alidns"
	"log"
	"strings"
)

const TTL = 600

type AliDns struct {
	cli    *alidns.Client
	domain string
	name   string // *.domain.com
}

func (a *AliDns) UpdateRecord(ctx context.Context, content string) error {
	domain := a.domain
	name := a.name
	// 只是完整域名的前段，如 *.domain.com 中的 *
	rr := strings.TrimSuffix(name, "."+domain)
	if rr == "" {
		rr = "@"
	}
	_, err := a.UpdateOrCreateDomainRecord(domain, rr, content)
	if err != nil {
		return err
	}
	return nil
}

func NewAliDns(regionId, key, secret string, domain, name string) *AliDns {
	if domain == "" {
		log.Panicf("domain not set")
	}
	if key == "" {
		log.Panicf("access-key-id not set")
	}
	if secret == "" {
		log.Panicf("access-key-secret not set")
	}
	client, err := alidns.NewClientWithAccessKey(regionId, key, secret)
	if err != nil {
		panic(err)
	}
	return &AliDns{
		cli:    client,
		domain: domain,
		name:   name,
	}
}

// 如果已经存在会报错: DomainRecordDuplicate
func (a *AliDns) AddDomainRecord(domainName, RR, typ, value string) (recordId string, err error) {
	q := alidns.CreateAddDomainRecordRequest()

	q.DomainName = domainName
	q.RR = RR
	q.Type = typ
	q.Value = value
	q.TTL = requests.NewInteger(TTL)

	rsp, err := a.cli.AddDomainRecord(q)
	if err != nil {
		err = errors.New(strings.Replace(err.Error(), "\n", " ", -1))

		if strings.Contains(err.Error(), "DomainRecordDuplicate") {
			return
		}
		return
	}

	recordId = rsp.RecordId

	return
}

// 根据domain和rr更新或者创建记录
func (a *AliDns) UpdateOrCreateDomainRecord(domain, rr, value string) (recordId string, err error) {
	rs, err := a.GetDomainRecordByRR(domain, rr)
	if err != nil {
		return
	}

	if len(rs) != 0 {
		// 更新第一个
		r := rs[0]

		recordId = r.RecordId

		if r.Value == value {
			return
		}
		err = a.UpdateDomainRecord(r.RecordId, rr, value)
		if err != nil {
			return
		}

		// 删除多余的
		for _, r := range rs[1:] {
			_ = a.DeleteDomainRecord(r.RecordId)
		}

	} else {
		recordId, err = a.AddDomainRecord(domain, rr, "A", value)
		if err != nil {
			return
		}
	}

	return
}

func (a *AliDns) UpdateDomainRecord(recordId string, RR, value string) (err error) {
	q := alidns.CreateUpdateDomainRecordRequest()

	q.RecordId = recordId
	q.RR = RR
	q.Value = value
	q.TTL = requests.NewInteger(TTL)

	_, err = a.cli.UpdateDomainRecord(q)
	if err != nil {
		err = errors.New(strings.Replace(err.Error(), "\n", " ", -1))
		return
	}

	return
}

func (a *AliDns) GetDomainRecordByRR(domain string, rr string) (rs []alidns.Record, err error) {
	q := alidns.CreateDescribeDomainRecordsRequest()
	q.RRKeyWord = rr
	q.DomainName = domain
	q.SearchMode = ""

	rsp, err := a.cli.DescribeDomainRecords(q)
	if err != nil {
		return
	}

	if len(rsp.DomainRecords.Record) == 0 {
		return
	}

	rs = rsp.DomainRecords.Record
	return
}

func (a *AliDns) DeleteDomainRecord(recordId string) (err error) {
	q := alidns.CreateDeleteDomainRecordRequest()

	q.RecordId = recordId

	_, err = a.cli.DeleteDomainRecord(q)
	if err != nil {
		err = errors.New(strings.Replace(err.Error(), "\n", " ", -1))
		return
	}

	return
}
