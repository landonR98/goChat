all: server frontend

server:
	@(cd ./src/golang/; go build -o ../../goChat)
	@echo server compiled

frontend:
	@(cd ./src/typescript; npm install && npx tsc)
	@echo frontend compiled
