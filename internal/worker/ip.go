package worker

import (
	"context"
	"github.com/zbysir/lookip/internal/lib/public_ip"
	"github.com/zbysir/lookip/internal/pkg/alidns"
	"log"
	"time"
)

type IpWorker struct {
	domain string
	rr     string

	aliCli *alidns.AliDns
}

func NewIpWorker(regionId, accessKey string, accessSecret string, domain string, rr string) *IpWorker {
	client := alidns.NewAliDns(regionId, accessKey, accessSecret)
	return &IpWorker{
		aliCli: client,
		domain: domain,
		rr:     rr,
	}
}

var nowIp = ""

// 将自己的公网ip上传到dns
func (i *IpWorker) LoopUpdateIp(ctx context.Context) {
	for {
		ip, err := public_ip.Get()
		if err != nil {
			select {
			case <-ctx.Done():
				goto Stop
			default:
			}

			log.Print(err)
			time.Sleep(3 * time.Second)
			continue
		}

		if ip != nowIp {
			_, err := i.aliCli.UpdateOrCreateDomainRecord(i.domain, i.rr, "A", ip)
			if err != nil {
				log.Print(err)
				time.Sleep(3 * time.Second)
				return
			}

			log.Printf("updateIp success: %s", ip)
			nowIp = ip
		}

		t := time.NewTimer(60 * time.Second)
		select {
		case <-ctx.Done():
			t.Stop()
			goto Stop
		case <-t.C:
		}
	}

Stop:
}
