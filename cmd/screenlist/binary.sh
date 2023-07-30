GOOS=windows GOARCH=amd64 go build -o bin/screenlist.exe main.go
GOOS=darwin GOARCH=arm64 go build -o bin/darwin_arm64 main.go
GOOS=darwin GOARCH=amd64 go build -o bin/darwin_amd64 main.go
GOOS=linux GOARCH=arm64 go build -o bin/linux_arm64 main.go
GOOS=linux GOARCH=amd64 go build -o bin/linux_amd64 main.go

mv ./bin/darwin_arm64 ./bin/screenlist
zip darwin_arm64.zip ./bin/screenlist
mv darwin_arm64.zip ./bin/darwin_arm64.zip

mv ./bin/darwin_amd64 ./bin/screenlist
zip darwin_amd64.zip ./bin/screenlist
mv darwin_amd64.zip ./bin/darwin_amd64.zip

mv ./bin/linux_arm64 ./bin/screenlist
zip linux_arm64.zip ./bin/screenlist
mv linux_arm64.zip ./bin/linux_arm64.zip

mv ./bin/linux_amd64 ./bin/screenlist
zip linux_amd64.zip ./bin/screenlist
mv linux_amd64.zip ./bin/linux_amd64.zip

rm ./bin/screenlist