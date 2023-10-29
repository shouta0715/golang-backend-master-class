package mail

import (
	"testing"

	"github.com/shouta0715/simple-bank/util"
	"github.com/stretchr/testify/require"
)

func TestSendEmailWithGmail(t *testing.T) {

	if testing.Short() {
		t.Skip("email test skipped")
	}

	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewGmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)

	subject := "A test Email"

	content := `
	<h1>Verify Emails</h1>
	<p>Please click the link below to verify your email address.</p>
	`

	to := []string{"kshouta0715@gmail.com"}
	attachFiles := []string{"../README.md"}

	err = sender.SendEmail(subject, content, to, nil, nil, attachFiles)

	// テストは通過済みです。練習用なのでエラーが発生します。
	require.NoError(t, err)
}
