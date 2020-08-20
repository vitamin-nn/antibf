# Anti bruteforce Service
## Overview
The service was made to control the rate of authentication requests.
Based on the loginning user data (login, password, ip) and service settings it defines whether the request is a brute force attempt.
The service works with grpc and command line interfaces.

The service is developed as final project of "Golang developer" course (by Otus)

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
  addblack    Adds ip network to the black list ()
  addwhite    Adds ip network to the white list
  clear       Clears specified bucket
  rmblack     Remove ip network from the black list
  rmwhite     Remove ip network from the white list
  
Required params for command would be view by typing: `.bin/antibf cli [command] --help`
