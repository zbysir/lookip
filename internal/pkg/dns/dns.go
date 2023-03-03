package dns

import (
	"context"
	"log"
)

type DNS interface {
	// UpdateRecord 更新记录, content: 域名解析的ip, proxied: 是否开启代理
	UpdateRecord(ctx context.Context, content string) error
}

type CombinedDNS struct {
	ds []DNS
}

func NewCombinedDNS(ds []DNS) *CombinedDNS {
	return &CombinedDNS{ds: ds}
}

func (c *CombinedDNS) UpdateRecord(ctx context.Context, content string) error {
	for _, d := range c.ds {
		err := d.UpdateRecord(ctx, content)
		if err != nil {
			log.Printf("update record by name error: %s", err)
		} else {
			log.Printf("update record by name success")
		}
	}
	return nil
}
