package main

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Info("Initializing Discord bot...")
	bot, err := discordgo.New("Bot " + apiToken)
	if err != nil {
		log.Error("There was an error authenticating with token: " + err.Error())
		os.Exit(2)
	}

	bot.AddHandler(onMessage)

	err = bot.Open()
	if err != nil {
		log.Error("There was an error starting up the bot: " + err.Error())
		os.Exit(2)
	}

	err = bot.UpdateStatusComplex(discordgo.UpdateStatusData{
		Activities: []*discordgo.Activity{{
			Name: activityString,
			Type: activityType,
		}},
		Status: string(onlineStatus),
	})
	if err != nil {
		log.Error("There was an error updating bot presence: " + err.Error())
		os.Exit(2)
	}

	log.Info("Done. Press Ctrl+C to stop bot.")
	sc := make(chan os.Signal, 1) // Wait for Ctrl+C
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-sc

	log.Info("Stopping...")
	err = bot.Close()
	if err != nil {
		log.Error("There was an error shutting down the bot: " + err.Error())
		os.Exit(2)
	}

	log.Info("Bye!")
}

func onMessage(session *discordgo.Session, event *discordgo.MessageCreate) {
	if event.ChannelID != submissionsChannelId {
		return
	}

	err := NewPost(session, event.Message)
	if err != nil {
		if _, ok := err.(*NoImagesError); ok {
			// Ignored
		} else {
			log.Error("There was an error creating post (or multiple): " + err.Error())
		}

		return
	}

	log.Debug("Created a new post from message " + event.Message.ID)
}
