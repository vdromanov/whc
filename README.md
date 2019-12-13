# whc
Working Hours Counter

### Cross-compile requirements:
* libc6-dev-armhf-cross
* arm-linux-gnueabi-gcc
### Building:
* CGO_ENABLED=1 CC=arm-linux-gnueabi-gcc GOOS=linux GOARCH=arm GOARM=5 go build -o <binary_fname>