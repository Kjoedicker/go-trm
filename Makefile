build:
	go build -o trm && ls -al | grep "trm"

install: build
	sudo install trm /usr/bin

push:
	git add .
	git commit
	git push origin
	git push github