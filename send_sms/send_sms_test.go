package send_sms

import (
	"fmt"
	"testing"
)

func TestSendSms(t *testing.T) {
	New("xxx", "yyy")
	fmt.Println(SendSms("18010489927", "1225", "verification code:5968"))
}
