package main

import (
	"flag"
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/FarmerChillax/fakeSSH/config"
	"github.com/FarmerChillax/fakeSSH/model"
	"github.com/FarmerChillax/fakeSSH/repository"
	"github.com/FarmerChillax/fakeSSH/vars"
	"github.com/gliderlabs/ssh"
	"github.com/sirupsen/logrus"
	stdssh "golang.org/x/crypto/ssh"
	"gorm.io/gorm"
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
	log.Printf("user: %s; pwd: %s\n", ctx.User(), password)
	db := vars.GetDB()
	err := db.Transaction(func(tx *gorm.DB) error {
		dataRepo := repository.NewDataRepository(db)
		data, err := dataRepo.FirstOrCreate(&model.Data{
			Username: ctx.User(),
			Password: password,
		})
		if err != nil {
			logrus.Errorf("PasswordHandler.dataRepo.FirstOrCreate err: %v; username: %s; password: %s",
				err, ctx.User(), password)
			return err
		}
		data.Count++
		if err := dataRepo.UpdateById(data, int64(data.Id)); err != nil {
			logrus.Errorf("PasswordHandler.dataRepo.UpdateById(%#v) err: %v", data, err)
			return err
		}
		return nil
	})

	if err != nil {
		logrus.Errorf("PasswordHandler.Transaction err: %v", err)
		return false
	}

	return false
}

var port int

func main() {
	if err := config.Load(); err != nil {
		log.Fatalf("config.Load err: %v", err)
	}
	hostKeySigner, err := CreateOrLoadKeySigner()
	if err != nil {
		log.Fatal(err)
	}
	flag.IntVar(&port, "port", 22, "SSH server port")
	flag.Parse()

	addr := fmt.Sprintf(":%d", port)
	s := &ssh.Server{
		Addr:    addr,
		Handler: Handler,
		// PublicKeyHandler: authHandler,
		PasswordHandler: PasswordHandler,
	}
	s.AddHostKey(hostKeySigner)
	logrus.Infof("ssh server is running on: %s", addr)
	log.Fatal(s.ListenAndServe())
}
