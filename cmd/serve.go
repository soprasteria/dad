// Copyright © 2017 NAME HERE <EMAIL ADDRESS>
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"github.com/soprasteria/dad/server"

	"github.com/soprasteria/dad/server/email"
	"github.com/soprasteria/dad/server/mongo"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

// serveCmd represents the serve command
var serveCmd = &cobra.Command{
	Use:   "serve",
	Short: "Launch D.A.D server",
	Long:  `D.A.D server will listen on 0.0.0.0:8080`,
	Run: func(cmd *cobra.Command, args []string) {
		email.InitSMTPConfiguration(viper.GetString("smtp.server"), viper.GetString("admin.name"), viper.GetString("smtp.user"), viper.GetString("smtp.identity"), viper.GetString("smtp.password"), viper.GetString("smtp.logo"))
		mongo.Connect()
		server.New(Version)
	},
}

func init() {
	// Get configuration from command line flags
	serveCmd.Flags().StringP("mongo-addr", "m", "localhost:27017", "URL to access MongoDB")
	serveCmd.Flags().StringP("mongo-username", "", "", "A user which has access to MongoDB")
	serveCmd.Flags().StringP("mongo-password", "", "", "Password of the mongo user")
	serveCmd.Flags().StringP("jwt-secret", "j", "dev-dad-secret", "Secret key used for JWT token authentication. Change it in your instance")
	serveCmd.Flags().StringP("reset-pwd-secret", "", "dev-dad-reset-pwd-to-change", "Secret key used when resetting the password. Change it in your instance")
	serveCmd.Flags().StringP("bcrypt-pepper", "p", "dev-dad-bcrypt", "Pepper used in password generation. Change it in your instance")
	serveCmd.Flags().BoolP("ldap-enable", "", true, "Enable LDAP")
	serveCmd.Flags().String("ldap-address", "", "LDAP full address like : ldap.server:389. Optional")
	serveCmd.Flags().String("ldap-baseDN", "", "BaseDN. Optional")
	serveCmd.Flags().String("ldap-domain", "", "Domain of the user. Optional")
	serveCmd.Flags().String("ldap-bindDN", "", "DN of system account. Optional")
	serveCmd.Flags().String("ldap-bindPassword", "", "Password of system account. Optional")
	serveCmd.Flags().String("ldap-searchFilter", "", "LDAP request to find users. Optional")
	serveCmd.Flags().String("ldap-attr-username", "cn", "LDAP attribute for username of users.")
	serveCmd.Flags().String("ldap-attr-firstname", "givenName", "LDAP attribute for firstname of users.")
	serveCmd.Flags().String("ldap-attr-lastname", "sn", "LDAP attribute for lastname of users.")
	serveCmd.Flags().String("ldap-attr-realname", "cn", "LDAP attribute for firstname of users.")
	serveCmd.Flags().String("ldap-attr-email", "mail", "LDAP attribute for lastname of users.")
	serveCmd.Flags().String("smtp-server", "", "SMTP server with its port.")
	serveCmd.Flags().String("smtp-user", "", "SMTP user for authentication.")
	serveCmd.Flags().String("smtp-password", "", "SMTP password for authentication.")
	serveCmd.Flags().String("smtp-logo", "", "Logo image that will be used in header of emails sent by DAD.")
	serveCmd.Flags().String("admin-name", "", "Email used as sender of emails")
	serveCmd.Flags().String("admin-email", "", "Email used as receiver of emails")
	serveCmd.Flags().String("name-receiver", "", "Email receiver's name")
	serveCmd.Flags().String("smtp-identity", "", "Identity of the sender")
	serveCmd.Flags().String("docktor-addr", "http://localhost:3000", "Docktor HTTP address. Format http://host:port")
	serveCmd.Flags().String("docktor-user", "user", "Docktor user to connect with")
	serveCmd.Flags().String("docktor-password", "password", "Docktor password to connect with")
	serveCmd.Flags().StringP("tasks-recurrence", "", "@every 20m", "Recurrence of back-end update tasks, like updating the deployment indicator (see https://godoc.org/github.com/robfig/cron)")

	// Bind env variables.
	_ = viper.BindPFlag("server.mongo.addr", serveCmd.Flags().Lookup("mongo-addr"))
	_ = viper.BindPFlag("server.mongo.username", serveCmd.Flags().Lookup("mongo-username"))
	_ = viper.BindPFlag("server.mongo.password", serveCmd.Flags().Lookup("mongo-password"))
	_ = viper.BindPFlag("auth.jwt-secret", serveCmd.Flags().Lookup("jwt-secret"))
	_ = viper.BindPFlag("auth.reset-pwd-secret", serveCmd.Flags().Lookup("reset-pwd-secret"))
	_ = viper.BindPFlag("auth.bcrypt-pepper", serveCmd.Flags().Lookup("bcrypt-pepper"))
	_ = viper.BindPFlag("ldap.enable", serveCmd.Flags().Lookup("ldap-enable"))
	_ = viper.BindPFlag("ldap.address", serveCmd.Flags().Lookup("ldap-address"))
	_ = viper.BindPFlag("ldap.baseDN", serveCmd.Flags().Lookup("ldap-baseDN"))
	_ = viper.BindPFlag("ldap.domain", serveCmd.Flags().Lookup("ldap-domain"))
	_ = viper.BindPFlag("ldap.bindDN", serveCmd.Flags().Lookup("ldap-bindDN"))
	_ = viper.BindPFlag("ldap.bindPassword", serveCmd.Flags().Lookup("ldap-bindPassword"))
	_ = viper.BindPFlag("ldap.searchFilter", serveCmd.Flags().Lookup("ldap-searchFilter"))
	_ = viper.BindPFlag("ldap.attr.username", serveCmd.Flags().Lookup("ldap-attr-username"))
	_ = viper.BindPFlag("ldap.attr.firstname", serveCmd.Flags().Lookup("ldap-attr-firstname"))
	_ = viper.BindPFlag("ldap.attr.lastname", serveCmd.Flags().Lookup("ldap-attr-lastname"))
	_ = viper.BindPFlag("ldap.attr.realname", serveCmd.Flags().Lookup("ldap-attr-realname"))
	_ = viper.BindPFlag("ldap.attr.email", serveCmd.Flags().Lookup("ldap-attr-email"))
	_ = viper.BindPFlag("smtp.server", serveCmd.Flags().Lookup("smtp-server"))
	_ = viper.BindPFlag("smtp.user", serveCmd.Flags().Lookup("smtp-user"))
	_ = viper.BindPFlag("smtp.password", serveCmd.Flags().Lookup("smtp-password"))
	_ = viper.BindPFlag("smtp.logo", serveCmd.Flags().Lookup("smtp-logo"))
	_ = viper.BindPFlag("admin.name", serveCmd.Flags().Lookup("admin-name"))
	_ = viper.BindPFlag("admin.email", serveCmd.Flags().Lookup("admin-email"))
	_ = viper.BindPFlag("name.receiver", serveCmd.Flags().Lookup("name-receiver"))
	_ = viper.BindPFlag("smtp.identity", serveCmd.Flags().Lookup("smtp-identity"))
	_ = viper.BindPFlag("docktor.addr", serveCmd.Flags().Lookup("docktor-addr"))
	_ = viper.BindPFlag("docktor.user", serveCmd.Flags().Lookup("docktor-user"))
	_ = viper.BindPFlag("docktor.password", serveCmd.Flags().Lookup("docktor-password"))
	_ = viper.BindPFlag("tasks.recurrence", serveCmd.Flags().Lookup("tasks-recurrence"))
	RootCmd.AddCommand(serveCmd)

}
