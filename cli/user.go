package cli

import (
	"fmt"
	"strings"

	"github.com/fabienbellanger/go-url-shortener/models"
	"github.com/fabienbellanger/go-url-shortener/repositories"
	"github.com/fabienbellanger/go-url-shortener/utils"
	"github.com/spf13/cobra"
)

var (
	userEmail     string
	userPassword  string
	userLastname  string
	userFirstname string
)

type userCreation struct {
	Lastname  string `validate:"required"`
	Firstname string `validate:"required"`
	Email     string `validate:"required,email"`
	Password  string `validate:"required,min=8"`
}

func init() {
	userCmd.Flags().StringVarP(&userLastname, "lastname", "l", "", "user lastname")
	userCmd.Flags().StringVarP(&userFirstname, "firstname", "f", "", "user firstname")
	userCmd.Flags().StringVarP(&userEmail, "email", "e", "", "user email")
	userCmd.Flags().StringVarP(&userPassword, "password", "p", "", "user password")

	userCmd.MarkFlagRequired("lastname")
	userCmd.MarkFlagRequired("firstname")
	userCmd.MarkFlagRequired("email")
	userCmd.MarkFlagRequired("password")

	rootCmd.AddCommand(userCmd)
}

var userCmd = &cobra.Command{
	Use:   "register",
	Short: "User creation",
	Long:  `User creation`,
	Run: func(cmd *cobra.Command, args []string) {
		user := userCreation{
			Lastname:  strings.TrimSpace(userLastname),
			Firstname: strings.TrimSpace(userFirstname),
			Password:  strings.TrimSpace(userPassword),
			Email:     strings.TrimSpace(userEmail),
		}

		// Validate data
		// -------------
		errs := utils.ValidateStruct(user)
		if len(errs) > 0 {
			fmt.Printf("\nError: invalid email or password (min 8 characters)\n")
			return
		}

		// Init config, logger and database
		// --------------------------------
		_, db, err := initConfigLoggerDatabase(false, true)
		if err != nil {
			fmt.Printf("\nError: %v\n", err)
			return
		}

		// User creation
		// -------------
		u := models.User{
			Lastname:  user.Lastname,
			Firstname: user.Firstname,
			Password:  user.Password,
			Username:  user.Email,
		}

		err = repositories.CreateUser(db, &u)
		if err != nil {
			fmt.Printf("\n%v\n", err)
			return
		}

		// Display result
		// --------------
		fmt.Printf(`
User successfully created:
    - Lastname:  %s
    - Firstname: %s
    - Email:     %s
    - Password:  %s
`,
			user.Lastname,
			user.Firstname,
			user.Email,
			user.Password,
		)
	},
}
