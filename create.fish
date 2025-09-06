#!/usr/bin/fish

set flag false

echo "Input release dir name"
read input

echo "Put in existing dir? (y/n)"
read check

if $check == "y"
    ehco "Insert path"
    read path
    flag = true
else
    echo ""
end

GOARC=arm64 go build .
mv ./intext ./ITX-CLI-arm64

GOARC=amd64 go build .
mv ./intext ./ITX-CLI

GOOS=windows GOARC=arm64 go build .
mv ./intext.exe ./ITX-CLI-arm64-WinEd.exe

GOOS=windows GOARC=amd64 go build .
mv ./intext.exe ./ITX-CLI-WinEd.exe

if $flag -eq true
    mkdir -p ./$path/$input/
    mv ./ITX-* ./$path/$input/
else
    mkdir -p ./Releases/$input
    mv ./ITX-* ./Releases/$input
end