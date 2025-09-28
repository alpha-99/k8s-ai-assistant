```shell
go run main.go chat
```

```shell
# 创建pod
curl -X POST -H "Content-Type: application/json" -d '{"yaml": "{\"kind\":\"Pod\", \"apiVersion\":\"v1\",\"metadata\":{\"name\":\"foo-app\",\"labels\":{\"app\":\"foo\"}},\"spec\":{\"containers\":[{\"name\":\"foo-app\",\"image\":\"higress-registry.cn-hangzhou.cr.aliyuncs.com/higress/http-echo:0.2.4-alpine\",\"args\":[\"-text=foo\"]}]}}"}' http://localhost:8080/pods

# 查询pod
curl http://localhost:8080/pods?ns=infra

# 删除POD
curl -X DELETE 'http://localhost:8080/pods?ns=default&name=foo-app'

# 验证API
curl 'http://localhost:8080/get/gvr?resource=pod'

```