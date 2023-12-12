podman run -d --name postgresfinalproject -p 5432:5432 -e POSTGRES_USER=postgres -e POSTGRES_PASSWORD=password -e POSTGRES_DB=concurrency -v postgres:14.2

podman run -d --name redis -p 6379:6379 redis:alpine

docker run -d --name mailhog -p 1025:1025 -p 8025:8025 mailhog/mailhog:latest
