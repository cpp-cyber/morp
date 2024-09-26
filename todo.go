package main

import (
	"fmt"
	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
	"regexp"
)

type todo struct {
	Id        int `sql:"AUTO_INCREMENT" gorm:"primaryKey"`
	User      string
	Task      string
	Completed bool
}

func getDiscordUser(userMention string) (*discordgo.Member, error) {
	// Extract the user ID from the mention using regex
	re := regexp.MustCompile(`<@!?(\d+)>`)
	matches := re.FindStringSubmatch(userMention)
	if len(matches) < 2 {
		return nil, fmt.Errorf("no user ID found in mention")
	}
	userID := matches[1]

	// Get the user by ID
	member, err := s.GuildMember(config.DiscordGuildID, userID)
	if err != nil {
		return nil, fmt.Errorf("error retrieving user from guild: %v", err)
	}

	return member, nil
}

func completeTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, todo *todo, err error) {
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to complete todo: " + err.Error(),
				Flags:   64,
			},
		})
		return
	}

	embed := embed.NewEmbed()
	embed.SetTitle("TODO Completed")
	embed.SetColor(0x33D6F5)
	embed.SetDescription(fmt.Sprintf("ㅤ\n%s\n```%s```", todo.User, todo.Task))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func addTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to add todo: " + err.Error(),
			},
		})
		return
	}
	user := i.ApplicationCommandData().Options[0].Options[0].StringValue()
	task := i.ApplicationCommandData().Options[0].Options[1].StringValue()

	embed := embed.NewEmbed()
	embed.SetTitle("TODO Added")
	embed.SetColor(0x33D6F5)
	embed.SetDescription(fmt.Sprintf("```%s```", task))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: user,
			Embeds:  []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func getTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, todos []todo, numCompleted int64, err error) {
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get todos: " + err.Error(),
				Flags:   64,
			},
		})
		return
	}

	member, err := getDiscordUser(todos[0].User)
	memberAvatar := member.User.AvatarURL("100")

	embed := embed.NewEmbed()
	embed.SetTitle("ㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤ")
	embed.SetColor(0x33D6F5)
	for _, t := range todos {
		embed.AddField(fmt.Sprintf("ID: %d", t.Id), fmt.Sprintf("```%s```ㅤ", t.Task))
	}
	embed.SetFooter(fmt.Sprintf("Uncompleted: %d ㅤ• ㅤCompleted: %d", len(todos), numCompleted))
	embed.SetAuthor(member.User.Username, memberAvatar)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func getAllTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, todos []todo, err error) {
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to get todos: " + err.Error(),
				Flags:   64,
			},
		})
		return
	}

	embed := embed.NewEmbed()
	embed.SetTitle("TODOsㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤㅤ\nㅤ")
	embed.SetColor(0x33D6F5)
	for _, t := range todos {
		embed.AddField(fmt.Sprintf("ID: %d", t.Id), fmt.Sprintf("%s\n```%s```ㅤ", t.User, t.Task))
	}
	embed.SetFooter(fmt.Sprintf("Uncompleted: %d", len(todos)))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}

func removeTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed := &discordgo.MessageEmbed{
		Title: "TODO(s) Removed",
		Color: 0x33D6F5,
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
		},
	})
}

func updateTodoResponse(s *discordgo.Session, i *discordgo.InteractionCreate, err error) {
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to update todo: " + err.Error(),
				Flags:   64,
			},
		})
		return
	}

	embed := embed.NewEmbed()
	embed.SetTitle("TODO Updated")
	embed.SetColor(0x33D6F5)
	embed.SetDescription(fmt.Sprintf("ㅤ\n**ID: %d**\n```%s```", i.ApplicationCommandData().Options[0].Options[0].IntValue(),
		i.ApplicationCommandData().Options[0].Options[1].StringValue()))

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
		},
	})
}
