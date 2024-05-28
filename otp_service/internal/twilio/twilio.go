package twilio

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/twilio/twilio-go"
	twilioApi "github.com/twilio/twilio-go/rest/verify/v2"
)

type Client struct {
	Client *twilio.RestClient
}

func NewClient(accountSID, authToken, from string) *Client {
	return &Client{Client: twilio.NewRestClientWithParams(twilio.ClientParams{
		Username: accountSID,
		Password: authToken,
	})}
}

func envSERVICESID() string {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}
	return os.Getenv("TWILIO_SERVICES_ID")
}

func (c *Client) SendSMS(to, message string) (string, error) {

	params := &twilioApi.CreateVerificationParams{}
	params.SetTo(to)
	params.SetChannel("sms")

	resp, err := c.Client.VerifyV2.CreateVerification(envSERVICESID(), params)
	if err != nil {
		return "", err
	}

	return *resp.AccountSid, nil
}

func (c *Client) OtpVerification(to, otp string) error {
	params := &twilioApi.CreateVerificationCheckParams{}
	params.SetTo(to)
	params.SetCode(otp)
	resp, err := c.Client.VerifyV2.CreateVerificationCheck(envSERVICESID(), params)
	if err != nil {
		return err
	} else if *resp.Status == "approved" {
		return nil
	}

	return nil
}
