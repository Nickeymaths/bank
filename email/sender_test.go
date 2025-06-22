package email

import (
	"fmt"
	"testing"

	"github.com/Nickeymaths/bank/util"
	"github.com/stretchr/testify/require"
)

func TestEmailSender(t *testing.T) {
	config, err := util.LoadConfig("..")
	require.NoError(t, err)

	sender := NewVmailSender(config.EmailSenderName, config.EmailSenderAddress, config.EmailSenderPassword)
	subject := "Verify your email"
	verifyEmailId := util.RandomInt(0, 10)
	secret := util.RandomString(32)
	content := fmt.Sprintf(`<h1> Welcome to my bank application </h1>
		<p> Please follow <a href="http://localhost:4000/v1/verify_email?id=%d,secret=%s">this link</a> to verify your account</p>
		`, verifyEmailId, secret)
	to := []string{
		"vinhpt15@viettel.com.vn",
	}

	err = sender.Send(
		subject,
		content,
		to,
		nil,
		nil,
		nil,
	)
	require.NoError(t, err)
}
