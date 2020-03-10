# go-relogger
<p align="left">
<a href="https://hits.seeyoufarm.com"/><img src="https://hits.seeyoufarm.com/api/count/incr/badge.svg?url=https%3A%2F%2Fgithub.com%2Fgjbae1212%2Fgo-relogger"/></a>
<a href="https://goreportcard.com/report/github.com/gjbae1212/go-relogger"><img src="https://goreportcard.com/badge/github.com/gjbae1212/go-relogger" alt="Go Report Card" /></a>
<a href="https://godoc.org/github.com/gjbae1212/go-relogger"><img src="https://img.shields.io/badge/godoc-reference-5272B4"/></a>
<a href="/LICENSE"><img src="https://img.shields.io/badge/license-MIT-GREEN.svg" alt="license" /></a>
</p>

## Overview
**go-relogger** is a logger that supports to reopen logging file when traps signal or is passed a time interval.
It can use to the logrotate(current logging file moves to an other logging file with suffix).
> Because a logger prints previous logging file until logger reopens.    

## Getting Started
```go
package main
import (   
   "os"
   "time"
   "github.com/gjbae1212/go-relogger"
)

func main() {
   logger, err := NewReLogger("your-log-path", relogger.WithSignals([]os.Signal{
   	                                           syscall.SIGHUP}), 
                                               relogger.WithRefreshDuration(time.Hour))   
}
```


## LICENSE
This project is following The MIT.
