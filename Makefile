start:
	docker compose up --build

install:
	cd client && npm install

client_test:
	cd client && npm start