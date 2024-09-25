@echo off
setlocal

rem Set the name for your binary
set BINARY_NAME=router_api

rem Build the Docker image
docker build -t router_api-aarch64-musl .

rem Create a temporary container
docker create --name temp router_api-aarch64-musl

rem Copy the binary from the container
docker cp temp:/router_api .\%BINARY_NAME%

rem Remove the temporary container
docker rm temp

echo Binary extracted as %BINARY_NAME%

scp router_api root@192.168.50.194:/secure_router/router_api/router_api

endlocal