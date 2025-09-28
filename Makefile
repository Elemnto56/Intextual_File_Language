buildall: windows windows_arm linux linux_arm macos macos_arm
	mkdir -p ./Releases/changeName/
	mv ITX-CLI* ./Releases/changeName/
	
windows:
	GOOS=windows GOARCH=amd64 go build -o ITX-CLI-WindEd.exe .

windows_arm:
	GOOS=windows GOARCH=arm64 go build -o ITX-CLI-arm64-WindEd.exe .

linux:
	GOOS=linux GOARCH=amd64 go build -o ITX-CLI .

linux_arm:
	GOOS=linux GOARCH=arm64 go build -o ITX-CLI-arm64 .

macos:
	GOOS=darwin GOARCH=amd64 go build -o ITX-CLI-MacOS .

macos_arm:
	GOOS=darwin GOARCH=arm64 go build -o ITX-CLI-arm64-MacOS .