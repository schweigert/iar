export GOPATH=$(shell pwd)

default:
	go build expansao.go
	./expansao < map
	convert -delay 0.3 -loop 0 *.png animation.gif
	rm *.png

clear:
	rm *.png
