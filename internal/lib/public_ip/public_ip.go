package public_ip

type IpGetter interface {
	Ip() (ip string, err error)
	Name() string
}

func Factory(t string) IpGetter {
	switch t {
	case "3322":
		return New3322()
	case "httpbin":
		fallthrough
	default:
		return NewHttpbinIp()
	}
}
