

all:build

build:
	# build as linux amd64
	env GOOS=linux  GOARCH=amd64 go build  -o ./bin/jetbrains_license_servers

image:
	docker build -t jetbrains_license_servers .

clean:
	-rm -f ./bin/*