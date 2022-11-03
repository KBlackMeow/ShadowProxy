# ShadowProxy v0.1.0
## What is ShadowProxy.
> The SP (ShadowProxy) was designed to protect your services from Port Scaning Attacks and Replay Attacks. It will hide your real service in your servers and send attackers facker informations.
## 
> The project is still at an early stage, so there are some bugs in the stage and those bugs are not fixed in time.
## How to start 
> Run "go build" and you will get 'shadowproxy', and template and static must be together with binary file in a directory.

> When first run ./shadowproxy, it will generate a yaml config file named 'config.yaml', then the SP will read those configs to program.

## ShadowProxy Config Introduction

```
# 0 normal log; 1 warning log; 2 error log. SP will show log whose level is bigger than loglevel.
loglevel: 0 

# Password will be generated when config.yaml created
password: a9a139419a60a98db6ffde10b50e4a52

# Enable https in auth service. 'server.key' and 'server.crt' are files name of ssl key and ssl crt. 
authssl: false

# Enable fillter function
enablefillter: true

# Log will output to .log file if consoleoutput is false
consoleoutput: true

debug: false

# The shadow service will be visited if client ip is not in white list or in black list.
shadow: 127.0.0.1:57575

# Those services will be run background. You can ban some service by delete it from array
services:
- auth # 127.0.0.1:57575
- flag # 127.0.0.1:40000
- cmd

# Port proxy rules, tcp/udp://bind address->backend address
rules:
- tcp://0.0.0.0:30000->127.0.0.1:40000
- udp://0.0.0.0:30050->127.0.0.1:40000

whitelist:
- 127.0.0.1
blacklist: []

# Those commonds will be executed when SP begin to run
cmd:
- whoami
```