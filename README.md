# hypermind

It's the source code of [my personal website](http://hypermind.com.cn) written in Go.

## Usage

1.  Ensure the project in dir '$GOPATH/src'

2.  Install [Redis](http://redis.io/) database.

3.  Get and install the library dependencies:

```bash
# redis driver
go get github.com/garyburd/redigo/redis

# go_lib
cd <$GOPATH1/src> # $GOPATH1 is the first part of $GOPATH.
git clone https://github.com/hyper-carrot/go_lib.git
```
4.  Edit hypermind.config for your need.

5.  Change mode of script & run:

```bash
cd <project_path>
chmod +x deploy.sh
./deploy.sh
```

## License
 
Copyright (C) 2013

Distributed under the BSD-style license, the same as Go.

