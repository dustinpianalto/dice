package exts

import (
	"fmt"
	"log"
	"regexp"
	"strings"

	"github.com/bwmarrin/discordgo"
	"github.com/dustinpianalto/dice/pkg/roller"
)

func DiceCommand(session *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.Bot {
		return
	}

	content := strings.ToLower(m.Content)
	r := regexp.MustCompile(`[^ "]+|"([^"]*)"`)
	parts := r.FindAllString(content, -1)
	var outString string
	for _, part := range parts {
		var label string
		var die string
		if strings.Contains(part, ":") {
			ps := strings.Split(part, ":")
			if len(ps) > 1 {
				label = strings.Join(ps[1:], ":")
			}
			die = ps[0]
		} else {
			die = part
		}
		i, s, err := roller.ParseRollString(die)
		if label != "" {
			outString += label + ": "
		}
		if err != nil && i != -2 {
			return
		} else if err != nil {
			outString = err.Error()
			break
		}
		outString += fmt.Sprintf("`%d` %s\n", i, s)
	}

	channel, err := session.Channel(m.ChannelID)
	if err != nil {
		log.Printf("Could not find channel %s\n", m.ChannelID)
	}
	session.ChannelMessageSend(channel.ID, outString)
}
