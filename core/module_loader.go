package core

import (
	"fmt"
	"strings"

	"github.com/kdrag0n/pyrowall/commands"

	"github.com/rs/zerolog/log"
)

func (b *Bot) registerCommand(cmd commands.Command) error {
	log.Debug().
		Str("command", cmd.Name).
		Msg("Registering command")

	lName := strings.ToLower(cmd.Name)
	if _, ok := b.Commands[lName]; ok {
		return fmt.Errorf("register command '%s': name already used", cmd.Name)
	}
	b.Commands[lName] = cmd

	for _, alias := range cmd.Aliases {
		log.Debug().
			Str("command", cmd.Name).
			Str("alias", alias).
			Msg("Registering command alias")

		lAlias := strings.ToLower(alias)
		if _, ok := b.Commands[lAlias]; ok {
			return fmt.Errorf("register alias '%s' for command '%s': name already used", alias, cmd.Name)
		}
		b.Commands[lAlias] = cmd
	}

	return nil
}

func (b *Bot) registerCommands(mod Module) error {
	for _, cmd := range mod.Commands() {
		err := b.registerCommand(cmd)
		if err != nil {
			return err
		}
	}

	return nil
}

func (b *Bot) loadModule(name string, cstr ModuleConstructor) error {
	log.Info().Str("name", name).Msg("Loading module")
	mod, err := cstr(b)
	if err != nil {
		return err
	}

	b.Modules[name] = mod
	err = b.registerCommands(mod)
	if err != nil {
		return err
	}

	return nil
}

// LoadModules loads all of the bot's modules. Automatically called by Start.
func (b *Bot) LoadModules() error {
	log.Info().Msg("Loading modules...")

	for name, cstr := range Modules {
		err := b.loadModule(name, cstr)
		if err != nil {
			return fmt.Errorf("load module '%s': %w", name, err)
		}
	}

	log.Info().Msg("All modules loaded.")
	return nil
}
