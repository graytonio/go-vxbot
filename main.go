package main

import (
	"fmt"
	"os"
	"os/signal"
	"regexp"
	"syscall"

	"github.com/bwmarrin/discordgo"
	"github.com/sirupsen/logrus"
)

var (
	BotToken string
	LogLevel logrus.Level
	err error
)

func parseEnv() {
	BotToken = os.Getenv("DISCORD_BOT_TOKEN")
	if BotToken == "" {
		logrus.Fatal("DISCORD_BOT_TOKEN is a required env")
	}

	logLevelSetting := os.Getenv("LOG_LEVEL")
	LogLevel, err = logrus.ParseLevel(logLevelSetting)
	if err != nil {
		logrus.Fatal("Invalid LOG_LEVEL")
	}
}

func main() {
	parseEnv()
	dg, err := discordgo.New("Bot " + BotToken)
	if err != nil {
		logrus.WithField("error", err).Fatal("error creating Discord session")
	}
	defer dg.Close()

	dg.AddHandler(handleNewMessage)
	dg.Identify.Intents |= discordgo.IntentMessageContent

	err = dg.Open()
	if err != nil {
		logrus.WithField("error", err).Fatal("error opening connection")
	}

	logrus.Info("Bot Online")
	waitForKillSignal()
}

func waitForKillSignal() {
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc
}

var twitterLinkRegex = regexp.MustCompile(`https:\/\/twitter\.com\/(.+)\/status/(.+)`)

func handleNewMessage(s *discordgo.Session, m *discordgo.MessageCreate) {
	// Do not modify messages sent from this bot
	if m.Author.ID == s.State.User.ID {
		return 
	}

	// Create logging context
	log := logrus.WithFields(logrus.Fields{
		"user": m.Author.Username,
		"guild": m.GuildID,
	})
	log.Debug("Got Message")
	
	// Check if message has twitter link in it
	changedText, changed := isTwitterLink(m.Content, log)
	if !changed {
		log.Debug("No twitter link ignoring")
		return
	}
	
	// Swap original message for updated link one
	_, err := s.ChannelMessageSend(m.ChannelID, fmt.Sprintf("From %s: %s", m.Author.Mention(), changedText))
	if err != nil {
		log.WithField("error", err).Error("Could not send new message")
		return
	}

	err = s.ChannelMessageDelete(m.ChannelID, m.ID)
	if err != nil {
		log.WithField("error", err).Error("Could not delete original message")
		return
	}
	
	log.Info("Replaced Twitter Link")
}

func isTwitterLink(content string, log *logrus.Entry) (string, bool) {
	hasTwitterLink := twitterLinkRegex.MatchString(content)
	if !hasTwitterLink {
		return "", false
	}

	replacedLink := twitterLinkRegex.ReplaceAllString(content, "https://vxtwitter.com/$1/status/$2")
	log.WithFields(logrus.Fields{
		"original": content,
		"replaced": replacedLink,
	}).Debug("Replaced Twitter Link")
	return replacedLink, true
}