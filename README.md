# Anti bruteforce service
[![Build Status](https://travis-ci.com/vitamin-nn/otus_anti_bruteforce.svg?branch=master)](https://travis-ci.com/vitamin-nn/otus_anti_bruteforce)
## Overview
The service was made to control the rate of authentication requests.
Based on the loginning user data (login, password, ip) and service settings it defines whether the request is a bruteforce attempt.
The service works with GRPC and Command line interfaces.

The service is developed as final project of "Golang developer" course (by Otus). For more information see: https://github.com/OtusGolang/final_project/blob/master/01-anti-bruteforce.md

## Main Service
### Running
1. Clone project
2. `make run`

### Finish work
`make down`

## Cli interface
1. Build: `make build`
2. Run command: `.bin/antibf cli [command] [params]`

Available Commands:
-  addblack: adds ip network to the black list ()
-  addwhite: adds ip network to the white list
-  clear: clears specified bucket
-  rmblack: remove ip network from the black list
-  rmwhite: remove ip network from the white list
  
Required params for command would be view by typing: `.bin/antibf cli [command] --help`
