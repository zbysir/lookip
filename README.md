# Lookip

Save your public ip address to AliDns/Cloudflare, It is usually used for hosts whose IP addresses change frequently.

## Usage

### Docker

#### For ali DNS

```
docker run \
-e DOMAIN=youdomin.com \
-e NAME=*.youdomin.com \
-e ACCESS_KEY_ID=xxxx \
-e ACCESS_KEY_SECRET=xxxxx \
-e REGION_ID=zh-hangzhou \
-e IP_GETTER=httpbin \
-e DNS=ali
--restart=unless-stopped
bysir/lookip
```

#### For cloudflare DNS

```docker run \
-e NAME=*.youdomin.com \
-e CF_TOKEN=xxxxx \
-e CF_ZONE_ID=zone_id_xxxx \
-e IP_GETTER=httpbin \
-e DNS=cloudflare
--restart=unless-stopped
bysir/lookip
```

### Window

```
set DOMAIN=youdomin.com
set NAME=*.youdomin.com
set ACCESS_KEY_ID=xxxx
set ACCESS_KEY_SECRET=xxxxx
set REGION_ID=zh-hangzhou
set IP_GETTER=httpbin

lookip.exe
```

## Config (ENV)

| key               | tops                                                                          |
|-------------------|-------------------------------------------------------------------------------|
| IP_GETTER         | default is `httpbin`, but `3322` is recommended in China                      |
| DNS               | 'ali' or 'cloudflare', set to 'ali,cloudflare' if you want to update together |
| ACCESS_KEY_ID     | AliDns AccessKeyID                                                            |
| ACCESS_KEY_SECRET | aliDns AccessKeySecret                                                        |
| REGION_ID         | aliDns REGION_ID                                                              |
| CF_TOKEN          | cloudflare Global API Key                                                     |
| CF_ZONE_ID        | cloudflare ZoneID                                                             |
| DOMAIN            | domain name                                                                   |
| NAME              | subdomain name                                                                |
