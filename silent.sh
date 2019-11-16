#!/usr/bin/env bash

go run main.go color logo off
go run main.go color ring fading FF0000 00FF00 0000FF
go run main.go speed pump 20 60  35 60  55 100  60 100
go run main.go speed fan 20 25  35 25  50 55  60 100