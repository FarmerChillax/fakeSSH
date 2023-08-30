package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/gliderlabs/ssh"
	stdssh "golang.org/x/crypto/ssh"
)

func CreateOrLoadKeySigner() (stdssh.Signer, error) {
	keyPath := filepath.Join(os.TempDir(), "fakessh.rsa")

	if _, err := os.Stat(keyPath); os.IsNotExist(err) {
		os.MkdirAll(filepath.Dir(keyPath), os.ModePerm)
		stderr, err := exec.Command("ssh-keygen", "-f", keyPath, "-t", "rsa", "-N", "").CombinedOutput()
		output := string(stderr)
		if err != nil {
			return nil, fmt.Errorf("fail to generate private key: %v - %s", err, output)

		}
	}

	privateBytes, err := os.ReadFile(keyPath)
	if err != nil {
		return nil, err
	}
	// 生成ssh.Signer
	return stdssh.ParsePrivateKey(privateBytes)
}

func Handler(s ssh.Session) {
	io.WriteString(s, "自由爱国公正平等公正友善公正公正诚信文明公正诚信平等文明富强法治法治公正诚信平等法治文明公正诚信文明公正自由\n")
}

func PasswordHandler(ctx ssh.Context, password string) bool {
	fmt.Printf("user: %s; pwd: %s\n", ctx.User(), password)
	return true
}

func main() {
	hostKeySigner, err := CreateOrLoadKeySigner()
	if err != nil {
		log.Fatal(err)
	}
	addr := ":22"
	s := &ssh.Server{
		Addr:    addr,
		Handler: Handler,
		// PublicKeyHandler: authHandler,
		PasswordHandler: PasswordHandler,
	}
	s.AddHostKey(hostKeySigner)

	log.Fatal(s.ListenAndServe())
}
