package auth

import (
	"net/mail"

	log "github.com/Sirupsen/logrus"
	"github.com/soprasteria/dad/server/email"
	"github.com/soprasteria/dad/server/types"
)

// SendWelcomeEmail sends a welcome email after a user's registration
func SendWelcomeEmail(user types.User) {
	err := email.Send(email.SendOptions{
		To: []mail.Address{
			{Name: user.DisplayName, Address: user.Email},
		},
		Subject: "Welcome to D.A.D",
		Body:    "Your account has been created !",
	})

	if err != nil {
		log.WithError(err).WithField("username", user.Username).Error("Failed to send welcome email")
	}
}
