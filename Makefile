run:
	@ go run cmd/main.go

image:
	@ docker build . -t echoserver

docker: image
	@ docker run --rm -p 8080:8080 echoserver

up:
	@ docker-compose up -d

down:
	@ docker-compose down --volumes
