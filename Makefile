build:win64
	@echo "only build in win64"
mac:
	@echo "build in macos"
	@CGO_ENABLED=0 GOOS=darwin GOARCH=arm64 go build -ldflags="-w -s" -o expatch-mac
	@upx expatch-mac
linux:
	@echo "build in linux"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-w -s" -o expatch-linux
	@upx expatch-linux
win64:
	@echo "build in windows64"
	@CGO_ENABLED=0 GOOS=windows GOARCH=amd64 go build -ldflags="-w -s" -o expatch-win64.exe
	@upx expatch-win64.exe
win32:
	@echo "build in windows32"
	@CGO_ENABLED=0 GOOS=windows GOARCH=386 go build -ldflags="-w -s" -o expatch-win32.exe
	@upx expatch-win32.exe
all:mac linux win64 win32
	@echo "build in mac/linux/windows"
# clean all build
clean:
	@echo "clean build"
	@rm -f expatch-mac
	@rm -f expatch-linux
	@rm -f expatch-win64.exe
	@rm -f expatch-win32.exe