.PHONY: build clean deploy gomodgen

build: gomodgen
	export GO111MODULE=on
	env GOOS=linux go build -ldflags="-s -w" -o bin/create_magicurl functions/create_magicurl/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/get_magicurl functions/get_magicurl/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/delete_magicurl functions/delete_magicurl/main.go

clean:
	rm -rf ./bin ./vendor Gopkg.lock

deploy: clean build
	sls deploy --verbose
	go run deployment_scripts/init_magicurl_db/main.go

gomodgen:
	chmod u+x gomod.sh
	./gomod.sh
