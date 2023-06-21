start:
	docker compose up --build

test:
	go test -v *.go

install:
	cd client && npm install

client_test:
	cd client && npm start