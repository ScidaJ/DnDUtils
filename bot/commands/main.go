package commands

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

// Types relating to commands.
type (
	Commands struct {
		Sugar          *zap.SugaredLogger
		GuildID        string
		DiscordSession *discordgo.Session
		Commands       []*discordgo.ApplicationCommand
	}

	HandleFunc func(s *discordgo.Session, i *discordgo.InteractionCreate, sugar *zap.SugaredLogger)
)

// Associates commands with their handlers.
func (c *Commands) AddCommandHandlers() {
	c.DiscordSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandlers := getCommandsHandlers()
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(c.DiscordSession, i, c.Sugar)
		}
	})
}

// Registers the commands with the bot.
// If they are not present within the slice returned by getComands()
// they will not be  registered
func (c *Commands) RegisterCommands() ([]*discordgo.ApplicationCommand, error) {
	commands := getCommands()

	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	c.Sugar.Info("Adding commands")

	for i, v := range commands {
		cmd := &discordgo.ApplicationCommand{
			Name:        v.Name,
			Description: v.Description,
			Options:     v.Options,
		}
		cmd, err := c.DiscordSession.ApplicationCommandCreate(c.DiscordSession.State.User.ID, c.GuildID, cmd)

		if err != nil {
			c.Sugar.Errorf("error adding command %v", v.Name, err)
			return registeredCommands, err
		}
		registeredCommands[i] = cmd
	}

	return registeredCommands, nil
}

// Removes the commands from the bot. This is done when the bot shuts down. Makes testing easier
// No need to fetch commands, though this may be changed in the future
func (c *Commands) RemoveCommands(r bool) error {
	if r {
		c.Sugar.Info("Removing commands...")

		for _, v := range c.Commands {
			err := c.DiscordSession.ApplicationCommandDelete(c.DiscordSession.State.User.ID, c.GuildID, v.ID)
			if err != nil {
				c.Sugar.Panicf("Cannot delete '%v' command: %v", v.Name, err)
				return err
			}
		}
	}

	c.Sugar.Info("Gracefully shutting down.")
	return nil
}

// Command structs located in commands.go must
// be in the returned slice or they will not be applied
func getCommands() []SlashCommand {
	return []SlashCommand{
		MakeParty,
	}
}

// Command handlers must be present in the returned
// map along with the command itself or they will
// not be registered.
func getCommandsHandlers() map[string]HandleFunc {
	return map[string]HandleFunc{
		"make-party": MakePartyHandler,
	}
}
