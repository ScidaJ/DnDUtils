package commands

import (
	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

type (
	Commands struct {
		Sugar          *zap.SugaredLogger
		GuildID        string
		DiscordSession *discordgo.Session
		Commands       []*discordgo.ApplicationCommand
	}

	HandleFunc func(s *discordgo.Session, i *discordgo.InteractionCreate)
)

func (c *Commands) AddCommandHandlers() {
	c.DiscordSession.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		commandHandlers := getCommandsHandlers()
		if h, ok := commandHandlers[i.ApplicationCommandData().Name]; ok {
			h(c.DiscordSession, i)
		}
	})
}

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

func (c *Commands) RemoveCommands(r bool) error {
	if r {
		c.Sugar.Info("Removing commands...")
		// // We need to fetch the commands, since deleting requires the command ID.
		// // We are doing this from the returned commands on line 375, because using
		// // this will delete all the commands, which might not be desirable, so we
		// // are deleting only the commands that we added.
		// registeredCommands, err := s.ApplicationCommands(s.State.User.ID, *GuildID)
		// if err != nil {
		// 	log.Fatalf("Could not fetch registered commands: %v", err)
		// }

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

func getCommands() []SlashCommand {
	return []SlashCommand{
		MakeParty,
	}
}

func getCommandsHandlers() map[string]HandleFunc {
	return map[string]HandleFunc{
		"make-party": MakePartyHandler,
	}
}
