# carbon-go-relay
A tool that work with graphite. Receive traffic from brubeck, then transfer to carbon-c-relay, follow predefined patterns.

## Backgroup of this project
- We want to use aggregation, one of carbon-c-relay's outstanding features, so metrics matched the same pattern should be transfered by brubeck to one carbon-c-relay instance, and carbon-c-relay must relay metrics and aggregation metrics in time. But when we confront with large traffic, at least 1Mpps,  one carbon-c-relay instance start to get hard, then metrics can't be relayed to carbon-cache timely, finally we can't get real-time and precise metrics. Obviously, carbon-c-relay can't be scaled horizontally if we use aggregation. So I do this project to split big traffic vertically by patterns, and metrics matched one pattern group will be transfered to the same carbon-c-relay instance, and finally be pushed into whisper timely.

## Features
- receive metrics from brubeck then transfer them to carbon-c-relay
- receive metrics from your own customized data source then transfered them to carbon-c-relay, with graphite metric foramt: `metric value timestamp`
- config statsGroup, then give some statistic

## Tutorial
### Run
- `cp cfg.example.josn your_cfg.json`, edit your own configure
- `go run main.go -c your_cfg.json`
