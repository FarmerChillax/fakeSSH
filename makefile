build-win-x64:
	CGO_ENABLE=0 GOOS=windows go build -o bin/win_x64/fake_sshd.exe

build-darwin-amd64:
	CGO_ENABLE=0 GOOS=darwin go build -o bin/darwin/fake_sshd

build-linux-amd64:
	CGO_ENABLE=0 GOOS=linux go build -o bin/linux/fake_sshd

# GOARCH=arm

build-all: build-darwin build-linux build-win

release-linux:
	tar -c bin/linux -f release/fakesshd_0.2.0_linux.tar

release-darwin:
	tar -c bin/linux -f release/fakesshd_0.2.0_darwin.tar

release-win:
	zip -r release/fakesshd_0.2.0_win.zip  bin/win 

release-all: release-linux release-darwin release-win

test:
	go test ./...

export PATH := $(GOPATH)/bin:$(PATH)
export GO111MODULE=on
LDFLAGS := -s -w

os-archs=darwin:amd64 darwin:arm64 freebsd:386 freebsd:amd64 linux:386 linux:amd64 linux:arm linux:arm64 windows:386 windows:amd64 linux:mips64 linux:mips64le linux:mips:softfloat linux:mipsle:softfloat linux:riscv64

fake_sshd_vesion =  v0.0.1

all: build

build: app

app:
	@$(foreach n, $(os-archs),\
		os=$(shell echo "$(n)" | cut -d : -f 1);\
		arch=$(shell echo "$(n)" | cut -d : -f 2);\
		gomips=$(shell echo "$(n)" | cut -d : -f 3);\
		target_suffix=$${os}_$${arch};\
		echo "Build $${os}-$${arch}...";\
		env CGO_ENABLED=0 GOOS=$${os} GOARCH=$${arch} GOMIPS=$${gomips} go build -trimpath -ldflags "$(LDFLAGS)" -o ./release/"$(fake_sshd_vesion)"/fakesshd_"$(fake_sshd_vesion)"_$${target_suffix} ;\
		echo "Build $${os}-$${arch} done";\
	)
	@mv ./release/"$(fake_sshd_vesion)"/fakesshd_"$(fake_sshd_vesion)"_windows_386 ./release/"$(fake_sshd_vesion)"/fakesshd_"$(fake_sshd_vesion)"_windows_386.exe
	@mv ./release/"$(fake_sshd_vesion)"/fakesshd_"$(fake_sshd_vesion)"_windows_amd64 ./release/"$(fake_sshd_vesion)"/fakesshd_"$(fake_sshd_vesion)"_windows_amd64.exe