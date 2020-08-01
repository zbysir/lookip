# lookip
Save your public ip address to AliDns

## Usage

#### docker

```
docker run \
-e DOMAIN=youdomin.com \
-e RR=* \
-e ACCESS_KEY_ID=xxxx \
-e ACCESS_KEY_SECRET=xxxxx \
-e REGION_ID=zh-hangzhou \
--restart=unless-stopped
bysir/lookip
```
#### window

```
set DOMAIN=youdomin.com
set RR=*
set ACCESS_KEY_ID=xxxx
set ACCESS_KEY_SECRET=xxxxx
set REGION_ID=zh-hangzhou

lookip.exe
```
