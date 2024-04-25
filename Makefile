all: server frontend

server:
	@(cd ./src/server/; go build -o ../../goChat)
	@echo server compiled

frontend:
	@(cd ./src/frontend; npm install && npx tsc)
	@echo frontend compiled
