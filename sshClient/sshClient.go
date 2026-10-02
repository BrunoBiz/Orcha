package sshClient

import (
	"example/Orcha/util"
	"log/slog"
	"os"

	"golang.org/x/crypto/ssh"
)

type SshClient struct {
	sshClient *ssh.Client
	config    util.Config
}

func NewSshClient(config util.Config) (*SshClient, error) {
	// Reads private key
	privateBytes, err := os.ReadFile(config.SSHKeyFile)
	if err != nil {
		return nil, err
	}

	// Creates signer
	signer, err := ssh.ParsePrivateKeyWithPassphrase(privateBytes, []byte(config.SSHKeyPassphrase))
	if err != nil {
		return nil, err
	}

	// SSH client config
	configSSH := &ssh.ClientConfig{
		User:            "root", // TODO - maybe change this at some point
		Auth:            []ssh.AuthMethod{ssh.PublicKeys(signer)},
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), // TODO - Change HostKeyCallback
	}

	client, err := ssh.Dial("tcp", config.SSHPveIP+":"+config.SSHPvePort, configSSH)
	if err != nil {
		return nil, err
	}

	sshClient := &SshClient{
		sshClient: client,
		config:    config,
	}

	return sshClient, nil
}

func (ssh *SshClient) NewSession(context string) (string, error) {
	session, err := ssh.sshClient.NewSession()
	if err != nil {
		slog.Error("SSH - Failed to create session: " + err.Error())
		return "", err
	}

	defer session.Close()

	returnValue, err := session.CombinedOutput(context)

	if err != nil {
		return string(returnValue), err
	}

	return string(returnValue), nil
}

func (ssh *SshClient) CloseConnection() {
	ssh.sshClient.Close()
}
