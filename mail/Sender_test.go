package mail

import (
	"testing"

	"github.com/314793615/simplebank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmail(t *testing.T) {
	config, err := util.NewConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "A test email"

	content := `
	<h1>HelloWorld</h1>
	<p>This is a test meassage from <a herf="http://baidu.com">Baidu.com</a></p>
	`
	attachFiles := []string{"../README.md"}
	to := []string{"314793615@qq.com"}
	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)
	require.NoError(t, err)

}