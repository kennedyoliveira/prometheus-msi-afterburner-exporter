#!/usr/bin/env bash

rm -rf build/

echo Building for MacOS x64...
env GOOS=darwin GOARCH=amd64 go build -o build/prometheus-msi-afterburner-exporter-darwin_x64 .

echo Building for Linux x64
env GOOS=linux GOARCH=amd64 go build -o build/prometheus-msi-afterburner-exporter-linux_x64 .

echo Building for Linux arm
env GOOS=linux GOARCH=arm go build -o build/prometheus-msi-afterburner-exporter-linux_arm .

echo BUilding for Windows x64
env GOOS=windows GOARCH=amd64 go build -o build/prometheus-msi-afterburner-exporter-windows_x64.exe .