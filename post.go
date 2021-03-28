package main

import (
	"errors"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/dustin/go-humanize"
	"net/http"
	"strconv"
)

type Image struct {
	URL    string
	Width  int
	Height int
}

type NoImagesError struct {
	Message *discordgo.Message
}

func (e *NoImagesError) Error() string {
	return "message has no image attachments or embeds"
}

func getImageURLs(message *discordgo.Message) []*Image {
	images := make([]*Image, 0, 1) // There generally would be just one
	for _, attachment := range message.Attachments {
		if attachment.Width*attachment.Height > 0 { // Images cannot have zero dimensions
			images = append(images, &Image{
				URL:    attachment.URL,
				Width:  attachment.Width,
				Height: attachment.Height,
			})
		}
	}

	for _, embed := range message.Embeds {
		if embed.Type == discordgo.EmbedTypeImage {
			images = append(images, &Image{
				URL:    embed.URL,
				Width:  embed.Thumbnail.Width,
				Height: embed.Thumbnail.Height,
			})
		} else if embed.Type == discordgo.EmbedTypeRich && embed.Image != nil {
			images = append(images, &Image{
				URL:    embed.Image.URL,
				Width:  embed.Image.Width,
				Height: embed.Image.Height,
			})
		}
	}

	return images
}

func NewPost(session *discordgo.Session, channelId string, message *discordgo.Message) error {
	images := getImageURLs(message)

	if len(images) == 0 {
		return &NoImagesError{Message: message}
	}

	for _, image := range images {
		var size string
		res, err := http.Head(image.URL)
		if err != nil {
			return errors.New("unable to send HTTP HEAD request: " + err.Error())
		}

		bytes, err := strconv.ParseUint(res.Header.Get("Content-Length"), 10, 64)
		if err != nil {
			return errors.New("unable to parse Content-Length header: " + err.Error())
		}
		size = humanize.IBytes(bytes)

		timestamp, err := message.Timestamp.Parse()
		if err != nil {
			return errors.New("unable to parse message timestamp: " + err.Error())
		}

		var nickname string // Retrieve user server-specific nickname
		member, err := session.GuildMember(message.GuildID, message.Author.ID)
		if err != nil || member == nil || member.Nick == "" {
			nickname = message.Author.Username
		} else {
			nickname = member.Nick
		}

		_, err = session.ChannelMessageSendEmbed(channelId, &discordgo.MessageEmbed{
			Type:  discordgo.EmbedTypeRich,
			Color: 44678,
			Footer: &discordgo.MessageEmbedFooter{
				Text: fmt.Sprintf(
					"Image sent by %s at %s with a file size of %s and a size of %d×%d",
					nickname, timestamp.Format("Mon Jan _2 2006 15:04:05 MST"), size,
					image.Width, image.Height,
				),
			},
			Image: &discordgo.MessageEmbedImage{
				URL: image.URL,
			},
		})

		if err != nil {
			return errors.New("could not send message embed: " + err.Error())
		}
	}

	return nil
}
