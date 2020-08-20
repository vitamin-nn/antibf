# Anti bruteforce Service
## Overview
The service was made to control the rate of authentication requests.
Based on the loginning user data (login, password, ip) and service settings it defines whether the request is a brute force attempt.

## Install
1. Download
```
go get github.com/vitamin-nn/otus_anti_bruteforce
```
2. Set ENV-params (see configs/.local.env for example)
