package commands

import (
	"math/rand"

	"github.com/bwmarrin/discordgo"
)

// Gets a random color for a party role
func randomColor() *int {
	min := 0
	max := 16777215
	random := rand.Intn(max-min+1) + min
	return &random
}

// Taken from https://github.com/Necroforger/dgwidgets/blob/master/util.go#L16-L23
func nextMessageReactionAddC(s *discordgo.Session) chan *discordgo.MessageReactionAdd {
	out := make(chan *discordgo.MessageReactionAdd)
	s.AddHandlerOnce(func(_ *discordgo.Session, e *discordgo.MessageReactionAdd) {
		out <- e
	})
	return out
}
