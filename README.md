# alert-notifier
从 [Alertmanager](https://prometheus.io/docs/alerting/latest/alertmanager/) 接收原始告警数据，并通过短信方式将告警信息转发给相关联系人的中间件服务。它通过 Webhook 接收原始告警，并使用 worker pool 模型异步发送短信，支持失败重试和死信队列。

---

## ✨ 特性

- ������ 接收原始 Alertmanager 告警
- ������ 异步发送短信通知
- ������ 短信发送失败自动重试（默认最多 3 次）
- ������ Worker Pool 设计：自动伸缩，支持空闲自动退出
- ������ 支持死信队列（Deadletter）：重试失败的告警记录
- ������ 灵活配置通知接收人


---

## ⚙️ 配置说明

### 1. 联系人配置

```json
{
  "contacts": [
    {
      "name": "Alice",
      "phone": "13800138000"
    },
    {
      "name": "Bob",
      "phone": "13900139000"
    }
  ]
}
```


### 2. Webhook 绑定 Alertmanager
```yaml
global:
route:
  receiver: "noop"
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 2h

  routes:
    - match:
        severity: critical
      receiver: webhook

receivers:
  - name: noop
    webhook_configs: []
  - name: webhook
    webhook_configs:
      - send_resolved: true
        url: "http://127.0.0.1:5001/webhook"
```

### 3. 快速开始
```shell
git clone github.com/echoH00/alert-notifier && cd alert-notifier && make

# 启动服务 --config指定联系人
./alert-notifier --config=/etc/alert-notifier/config.json

部署deploy/目录
```

