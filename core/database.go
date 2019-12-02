package core

import (
	"fmt"

	"github.com/go-pg/pg/v9"
)

func (b *Bot) connectToDB() error {
	if b.Config.Database.Type != "postgres" {
		return fmt.Errorf("database type '%s' is unsupported", b.Config.Database.Type)
	}

	b.DB = pg.Connect(&pg.Options{
		Network:         b.Config.Database.Protocol,
		Addr:            b.Config.Database.Address,
		User:            b.Config.Database.User,
		Password:        b.Config.Database.Password,
		Database:        b.Config.Database.Database,
		ApplicationName: "pyrowall",

		MaxRetries:            3,
		RetryStatementTimeout: true,
	})

	_, err := b.DB.Exec("SELECT 1")
	if err != nil {
		return err
	}

	return nil
}
