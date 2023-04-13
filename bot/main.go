package main

import (
	"dndutils/bot/commands"
	"flag"
	"os"
	"os/signal"

	"github.com/bwmarrin/discordgo"
	"go.uber.org/zap"
)

var Token string
var GuildID string
var Endpoint string
var RemoveCommands bool

func init() {
	flag.StringVar(&Token, "t", "", "Bot Token")
	flag.StringVar(&GuildID, "g", "", "Guild ID")
	flag.StringVar(&Endpoint, "e", "", "API Endpoint")
	flag.BoolVar(&RemoveCommands, "R", true, "Remove commands after shutdown")
	flag.Parse()
}

func main() {
	discord, err := discordgo.New("Bot " + Token)

	logger, _ := zap.NewProduction()
	defer logger.Sync()
	sugar := logger.Sugar()

	if err != nil {
		sugar.Error("error creating Discord session,", err)
		return
	}

	discord.AddHandler(messageHandler)

	err = discord.Open()
	if err != nil {
		sugar.Error("error opening connection,", err)
	}

	defer discord.Close()

	c := commands.Commands{
		Sugar:          sugar,
		GuildID:        GuildID,
		DiscordSession: discord,
		Endpoint:       Endpoint,
	}

	c.AddCommandHandlers()
	c.RegisterCommands()

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	sugar.Info("Bot is now running. Press CTRL-C to exit.")
	<-sc

	err = c.RemoveCommands(RemoveCommands)

	if err != nil {
		sugar.Error("error removing commands", err)
	}
}

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}

	if m.Content == "ping" {
		s.ChannelMessageSend(m.ChannelID, "Pong!")
	}

	if m.Content == "pong" {
		s.ChannelMessageSend(m.ChannelID, "Ping!")
	}
}
