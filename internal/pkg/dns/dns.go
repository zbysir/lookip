package dns

import (
	"context"
	"log"
)

type DNS interface {
	// UpdateRecord 更新记录, content: 域名解析的ip
	// 如果 force 为 true，则当 content 前后一致时也会更新，默认不更新
	UpdateRecord(ctx context.Context, content string, force bool) (bool, error)
}

type CombinedDNS struct {
	ds []DNS
}

func NewCombinedDNS(ds []DNS) *CombinedDNS {
	return &CombinedDNS{ds: ds}
}

func (c *CombinedDNS) UpdateRecord(ctx context.Context, content string, force bool) (bool, error) {
	var affected bool
	for _, d := range c.ds {
		aff, err := d.UpdateRecord(ctx, content, force)
		if err != nil {
			log.Printf("update record by name error: %s", err)
		} else {
			log.Printf("update record by name success")
		}
		if aff {
			affected = true
		}
	}
	return affected, nil
}
