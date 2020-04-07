.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_magicurl create_magicurl/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_magicurl get_magicurl/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete_magicurl delete_magicurl/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/initialize_magicurl_id initialize_magicurl_id/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
