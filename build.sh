GOOS=linux GOARCH=amd64 go build -o dist/linux.bin ./src;
GOOS=darwin GOARCH=arm64 go build -o dist/mac-arm.bin ./src;
