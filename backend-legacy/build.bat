@echo off

REM Check if build folder exists, if not create it
if not exist ".\builds" (
    echo Build folder does not exist. Creating it...
    mkdir ".\builds"
)

REM Build the Docker image
docker build -t backend-aarch64-musl .

REM Create a temporary container
docker create --name temp backend-aarch64-musl

REM Copy the binary from the container
docker cp temp:/backend .\builds\backend

REM Remove the temporary container
docker rm temp

echo Binary extracted as backend