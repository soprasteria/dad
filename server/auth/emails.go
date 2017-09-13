package auth

import (
	"net/mail"

	log "github.com/Sirupsen/logrus"
	"github.com/matcornic/hermes"
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

		Body: hermes.Email{
			Body: hermes.Body{
				Name: user.DisplayName,
				Intros: []string{
					"Welcome to D.A.D! We're very excited to have you on board. Your account has been created!",
				},
			},
		}})

	if err != nil {
		log.WithError(err).WithField("username", user.Username).Error("Failed to send welcome email")
	}
}
