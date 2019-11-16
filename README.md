# coolctl

[![Go Report Card](https://goreportcard.com/badge/github.com/arkste/coolctl)](https://goreportcard.com/report/github.com/arkste/coolctl)
[![License](https://img.shields.io/badge/license-MIT-blue.svg)](https://raw.githubusercontent.com/arkste/elsi/master/LICENSE)
[![Golang](https://img.shields.io/badge/Go-1.13-blue.svg)](https://golang.org)
![Linux](https://img.shields.io/badge/Supports-Linux-green.svg)

This is just a Playground for my gousb-Adventures (https://github.com/google/gousb).

Currently only reading the status of my NZXT Kraken X62 AIO:

```
$ sudo apt-get install libusb-1.0 pkg-config
$ git clone https://github.com/arkste/coolctl.git
$ cd coolctl
$ make dep
$ go run main.go
============================================
  Liquid temperature 32.7 °C
  Fan speed 527 rpm
  Pump speed 2040 rpm
  Firmware Version: 6.0.2
============================================
```
