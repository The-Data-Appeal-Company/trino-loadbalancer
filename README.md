# Trino Loadbalancer

![Docker](https://github.com/The-Data-Appeal-Company/trino-loadbalancer/workflows/Docker/badge.svg)
![Tests](https://github.com/The-Data-Appeal-Company/trino-loadbalancer/workflows/Tests/badge.svg)
[![Go Report Card](https://goreportcard.com/badge/github.com/The-Data-Appeal-Company/trino-loadbalancer)](https://goreportcard.com/report/github.com/The-Data-Appeal-Company/trino-loadbalancer)
[![Coverage Status](https://coveralls.io/repos/github/The-Data-Appeal-Company/trino-loadbalancer/badge.svg?branch=master)](https://coveralls.io/github/The-Data-Appeal-Company/trino-loadbalancer?branch=master)

Fast, high available load balancer for trino with smart routing rules


## Configuration 

#### Minimal configuration

```yaml
proxy:
  port: 8998

routing:
  rule: round-robin

persistence:
  postgres:
    db: 'postgres'
    host: '127.0.0.1'
    port: 5432
    username: 'postgres'
    password: 'test'
    ssl_mode: 'disable'

discovery:
  providers:
    - provider: static
      enabled: true
      static:
        clusters:
          - name: cluster-0
            enabled: true
            url: http://localhost:8080

session:
  store:
    redis:
      opts:
        prefix: 'trino::'
        max_ttl: 24h
```

## Deploy

Todo
