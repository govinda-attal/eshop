
init:
	go mod tidy

test:
	go test ./... -race -count=1 -cover

build: test
	rm -rf ./dist
	mkdir dist/
	mkdir dist/configs
	mkdir dist/scripts
	go build -o ./dist/eshop .
	cp ./configs/app-cfg.yaml ./dist/configs/app-cfg.yaml
	cp -r ./scripts/db ./dist/scripts/db

migrate: build
	cd dist && ./eshop migrate up

serve: build
	cd dist && ./eshop serve

docker-serve:
	docker-compose up -d --build eshop swagger-ui

docker-serve-all:
	docker-compose up -d --build

docker-deps:
	docker-compose up -d crdb migrate

docker-clean:
	docker-compose down

clean:
	rm ./dist/ -rf

pack:
	docker build -t gattal/eshop:latest .

upload:
	docker push gattal/eshop:latest	

ship: init test pack upload clean


api-gen:
	mkdir -p pkg/eshop/api
	oapi-codegen -package api api/openapi.yaml  > pkg/eshop/api/eshop.gen.go
