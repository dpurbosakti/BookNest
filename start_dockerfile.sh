#!/bin/sh

appName=booknest
localPort=8080 # local port di server akan di bind ke port berapa. disarankan mengikuti exposedPort dan pastikan port nya tidak dipakai oleh service lain.
exposedPort=8080 # ganti dengan APP_PORT sesuai dengan config.yaml
imageName=$appName:latest
containerName=$appName

echo Stop old container if any...
docker stop $containerName || true

echo Delete old container if any...
docker rm -f $containerName || true

echo Delete old image if any...
docker rmi -f $imageName || true

echo Building new image...
docker build -t $imageName -f Dockerfile  .

echo Run new container...
docker run -d -p $localPort:$exposedPort --name $containerName $imageName
