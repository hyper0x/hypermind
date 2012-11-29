# go-web-demo

A demo website written in Go.

## Usage

1.  Ensure the project path in $GOPATH

2.  Install [Redis](http://redis.io/) database.

3.  Get and install the library dependencies:

```bash
# session manager 
go get github.com/astaxie/session
go get github.com/astaxie/session/providers/memory

# redis driver
go get github.com/garyburd/redigo/redis
```
4.  Edit go-web-demo.config for your need.

5.  Change mode of script & run:

```bash
cd <project_path>
chmod +x deploy.sh
./deploy.sh
```

## License
 
Copyright (C) 2012

Distributed under the BSD-style license, the same as Go.

