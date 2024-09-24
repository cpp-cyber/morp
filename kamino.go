package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/cookiejar"
	"time"

	"github.com/Clinet/discordgo-embed"
	"github.com/bwmarrin/discordgo"
)

func NewKaminoClient() {
    cl = Login()
    if cl == nil {
        fmt.Println("Login failed")
    }
}

func Login() *http.Client {
    jar, err := cookiejar.New(nil)
    if err != nil {
        fmt.Println(err)
        return nil
    }
    client := &http.Client{
        Jar: jar,
    }

    resp, err := doAPIRequest("POST", config.KaminoLoginEndpoint, map[string]any{
        "username": config.KaminoUser,
        "password": config.KaminoPass,
    })
    if err != nil {
        fmt.Println(err)
        return nil
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        return nil
    }

    client.Jar.SetCookies(resp.Request.URL, resp.Cookies())
    return client
}

func refreshLogin() {
    for {
        <-time.After(30 * time.Minute)
        NewKaminoClient()
    }
}

func getPods(s *discordgo.Session, i *discordgo.InteractionCreate) {
    resp, err := doAPIRequest("GET", config.KaminoGetPodsEndpoint, nil)
    if err != nil {
        fmt.Println(err)
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Failed to get pods",
            },
        })
        return
    }
    defer resp.Body.Close()

    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        fmt.Println(err)
    }

    type pods struct {
        Name string `json:"Name"`
        ResourceGroup string `json:"ResourceGroup"`
        ServerGUID string `json:"ServerGUID"`
    }

    var podList []pods
    err = json.Unmarshal(respBody, &podList)
    if err != nil {
        fmt.Println(err)
    }

    embed := embed.NewEmbed()
    embed.SetTitle("Pods")
    embed.SetColor(0xffab40)

    podString := ""
    for i, pod := range podList {
        podString += fmt.Sprintf("%d. %s\n", i+1, pod.Name)
    }

    embed.AddField("Pods", podString)

    s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseChannelMessageWithSource,
        Data: &discordgo.InteractionResponseData{
            Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
        },
    })
}

func deletePod(s *discordgo.Session, i *discordgo.InteractionCreate) {
    podId := i.ApplicationCommandData().Options[0].Options[0].StringValue()

    resp, err := doAPIRequest("DELETE", config.KaminoDeleteEndpoint + "/" + podId, nil)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        message := fmt.Sprintf("Failed to delete pod %s", podId)
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &message,
        })
        return
    }

    message := fmt.Sprintf("Pod %s deleted", podId)
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Content: &message,
    })
}

func bulkDeletePods(s *discordgo.Session, i *discordgo.InteractionCreate) {
    data := make(map[string]any)
    data["Filters"] = i.ApplicationCommandData().Options[0].Options[0].StringValue()

    resp, err := doAPIRequest("DELETE", config.KaminoBulkDeleteEndpoint, data)
    if err != nil {
        embed := embed.NewEmbed()
        embed.SetTitle("Failed to delete pods")
        embed.SetColor(0xff0000)
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: "Status",
            Value: resp.Status,
        })
        embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
            Name: "Error",
            Value: err.Error(),
        })

        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Embeds: []*discordgo.MessageEmbed{embed.MessageEmbed},
            },
        })
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        message := "Failed to delete pods. Check logs for more information"
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &message,
        })
    }

    message := "Pods deleted"
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Content: &message,
    })
}

func refreshTemplates(s *discordgo.Session, i *discordgo.InteractionCreate) {
    resp, err := doAPIRequest("POST", config.KaminoRefreshTemplatesEndpoint, nil)
    if err != nil {
        fmt.Println(err)
    }
    defer resp.Body.Close()

    if resp.StatusCode != http.StatusOK {
        message := "Failed to refresh templates. Check logs for more information"
        s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
            Content: &message,
        })
        return
    }

    message := "Templates refreshed"
    s.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
        Content: &message,
    })
}

func doAPIRequest(verb, endpoint string, data map[string]any) (*http.Response, error) {
    var body io.Reader
    if data != nil {
        jsonBody, err := json.Marshal(data)
        if err != nil {
            return nil, err
        }
        body = bytes.NewReader(jsonBody)
    }

    req, err := http.NewRequest(verb, config.KaminoURL + endpoint, body)
    if err != nil {
        return nil, err
    }

    if data != nil {
        req.Header.Set("Content-Type", "application/json")
    }

    resp, err := cl.Do(req)
    if err != nil {
        return nil, err
    }

    if resp.StatusCode != http.StatusOK {
        return nil, fmt.Errorf("API request failed with status code %d", resp.StatusCode)
    }

    return resp, nil
}
