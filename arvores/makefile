export GOPATH=$(shell pwd)

default:
	go build expansao.go
	./expansao < map
	convert -delay 10 -loop 0 *.png animation.gif
	rm *.png

clear:
	rm *.png
