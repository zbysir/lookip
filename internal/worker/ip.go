package worker

import (
	"context"
	"github.com/zbysir/lookip/internal/lib/public_ip"
	"github.com/zbysir/lookip/internal/pkg/dns"
	"log"
	"time"
)

type IpWorker struct {
	dns dns.DNS

	ipGetter public_ip.IpGetter
}

func NewIpWorker(dns dns.DNS, ipGetter public_ip.IpGetter) *IpWorker {
	return &IpWorker{
		dns:      dns,
		ipGetter: ipGetter,
	}
}

// LoopUpdateIp 将自己的公网 ip 上传到dns
func (i *IpWorker) LoopUpdateIp(ctx context.Context) {
	for {
		ip, err := i.ipGetter.Ip()
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

		aff, err := i.dns.UpdateRecord(context.Background(), ip, false)
		if err != nil {
			log.Print(err)
			time.Sleep(3 * time.Second)
			continue
		}

		if aff {
			log.Printf("updateIp success: %s", ip)
		}

		t := time.NewTimer(120 * time.Second)
		select {
		case <-ctx.Done():
			t.Stop()
			goto Stop
		case <-t.C:
		}
	}

Stop:
}
