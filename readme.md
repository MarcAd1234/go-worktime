# go-worktime
simple cli tool to track your time at work - written in Go Lang.

<img src="Images/go-worktime.jpg" alt="Logo" width="200"/>


# What it does
1. Save binary in a folder of your choice
2. 

# Normal Build

``` bash
go build -o worktime main.go
```


# Build for MacOS ARM64
``` bash
GOOS=darwin GOARCH=arm64 go build -o worktime_mac main.go
```
## Run on Mac
``` bash
chmod +x worktime_mac
sudo xattr -rd com.apple.quarantine ./worktime_mac
./worktime_mac
```


