package main

import (
	"github.com/bwmarrin/discordgo"
	"os"
)

// Bot application API token.
var apiToken = os.Getenv("API_TOKEN")

// Presence online status (e.g. "online" or "idle".)
var onlineStatus = discordgo.StatusOnline

// Presence activity type (e.g. "Playing" or "Listening to".)
var activityType = discordgo.ActivityTypeListening

// Presence activity value (e.g. "Playing X" where X is this string.)
var activityString = "#just-monika"

// ID of the channel to listen for messages in.
var submissionsChannelId = "825666443675697172"

// ID of the channel to send posts to.
var galleryChannelId = "825667478519545906"

// Emoji name and ID (for custom emojis) or
// Unicode character for Unicode emojis:
// var emojiId = "💚"
var emojiId = "note:823616525156089926"

// Text to render in post footer.
// Substituted values:
// %[1]s: sender nickname
// %[2]s: time and date of message creation
// %[3]s: file size
// %[4]d: image width
// %[5]d: image height
var footerText = "Image sent by %[1]s at %[2]s with a file size of %[3]s and a size of %[4]d×%[5]d"