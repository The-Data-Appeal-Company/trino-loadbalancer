# Presto Hub

# presto-loadbalancer 
[![Go Report Card](https://goreportcard.com/badge/github.com/The-Data-Appeal-Company/presto-loadbalancer)](https://goreportcard.com/report/github.com/The-Data-Appeal-Company/presto-loadbalancer)
![Docker](https://github.com/The-Data-Appeal-Company/presto-loadbalancer/workflows/Docker/badge.svg)
![Tests](https://github.com/The-Data-Appeal-Company/presto-loadbalancer/workflows/Tests/badge.svg)

Fast, high available load balancer for presto with smart routing rules

## Deploy

Todo

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
    username: 'prestohub'
    password: 'presto'
    ssl_mode: 'disable'

session:
  store:
      standalone:
        enabled: true
        host: '127.0.0.1:6379'
```
