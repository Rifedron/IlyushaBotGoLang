package main

import (
	"awesomeProject/main/IlyushaBot"
	"awesomeProject/main/IlyushaBot/offers"
	"awesomeProject/main/IlyushaBot/privateVCs"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"os"
	"strings"
)

var (
	Bob, _   = discordgo.New("Bot " + IlyushaBot.Cfg.Token)
	shutdown = make(chan string)
)

func main() {
	fmt.Println("Starting bot")
	Bob.Identify.Intents = discordgo.IntentsAll
	Bob.AddHandler(onInteraction)
	for _, handler := range compiledEventHandlers(
		//Event handler slices
		offers.OfferEvents,
		privateVCs.PrivateVcEvents,
	) {
		Bob.AddHandler(handler)
	}
	err := Bob.Open()
	if err == nil {
		fmt.Println("Bot is working")
	} else {
		fmt.Println("Bot is ded")
		os.Exit(1)
	}

	addCommands(
		//Application command slices
		offers.OfferCommands,
		privateVCs.PrivatesCommands,
	)
	fmt.Println("Commands added")

	<-shutdown
}

func addCommands(commands ...[]*discordgo.ApplicationCommand) {
	guilds, _ := Bob.UserGuilds(100, "", "")
	for _, guild := range guilds {
		oldCommands, _ := Bob.ApplicationCommands(Bob.State.User.ID, guild.ID)
		for _, cmd := range oldCommands {
			_ = Bob.ApplicationCommandDelete(cmd.ApplicationID, cmd.GuildID, cmd.ID)
		}
		for _, commandCollection := range commands {
			for _, command := range commandCollection {
				go func(c *discordgo.ApplicationCommand) {
					_, err := Bob.ApplicationCommandCreate(Bob.State.User.ID, guild.ID, c)
					if err != nil {
						fmt.Println(err.Error())
					}
				}(command)
			}
		}
	}
}

func compiledEventHandlers(handlers ...[]interface{}) []interface{} {
	var allHandlers []interface{}
	for _, handler := range handlers {
		allHandlers = append(allHandlers, handler...)
	}
	return allHandlers
}

func onInteraction(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var (
		f func(s *discordgo.Session, i *discordgo.InteractionCreate)
		b bool
	)
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		f, b = mergedInteractionMap[i.ApplicationCommandData().Name]
		break
	case discordgo.InteractionMessageComponent:
		f, b = mergedInteractionMap[i.MessageComponentData().CustomID]
		break
	case discordgo.InteractionModalSubmit:
		for k, v := range mergedInteractionMap {
			_, found := strings.CutPrefix(i.ModalSubmitData().CustomID, k)
			if found {
				f = v
				b = true
				break
			}
		}
		break
	}
	if b {
		f(s, i)
	}
}

var mergedInteractionMap = makeMergedMap(
	offers.OfferInteractions,
	privateVCs.PrivatesInteractions,
)

func makeMergedMap(maps ...map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate)) map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	aMap := map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){}
	for _, m := range maps {
		for k, v := range m {
			aMap[k] = v
		}
	}
	return aMap
}
