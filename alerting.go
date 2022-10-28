package main

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"os"
)

func IsDown(t *target, ts *targetStatus) {
	if ts.Alerting == false {
		ts.Alerting = true
		go AlertDiscordEmbed(t, ts)
		go AlertPushover(t, ts)
		log.Print("monitoring /!\\ " + t.Name + " via " + t.Method + " is down.")
	}
}

func IsUp(t *target, ts *targetStatus) {
	if ts.Alerting == true {
		ts.Alerting = false
		go AlertDiscordEmbed(t, ts)
		go AlertPushover(t, ts)
		log.Print("monitoring /!\\ " + t.Name + " via " + t.Method + " is up.")
	}
}

func AlertPushover(t *target, ts *targetStatus) {
	var title string
	if ts.Up == false {
		title = t.Name + " is down"
	} else {
		title = t.Name + " is up"
	}
	pushOverNotif := pushNotif{
		Token:    os.Getenv("PUSHOVER_TOKEN"),
		User:     os.Getenv("PUSHOVER_USER"),
		Title:    title,
		Message:  t.Target + " via " + t.Method,
		Priority: 0,
		Url:      "https://ericz.dev",
	}
	str, err := json.Marshal(pushOverNotif)
	if err != nil {
		log.Print(err)
	}
	go SendPushOverNotif(str)
}

func SendPushOverNotif(str []byte) {
	http.Post("https://api.pushover.net/1/messages.json", "application/json", bytes.NewBuffer(str))
}

func AlertDiscordEmbed(t *target, ts *targetStatus) {
	var color int
	var title string
	if ts.Up == false {
		color = 16756141
		title = t.Name + " is down"
	} else {
		color = 13303743
		title = t.Name + " is up"
	}
	discordEmbed := embed{
		Content: nil,
		Embeds: []embedBody{
			embedBody{
				Title:       title,
				Description: t.Target + " via " + t.Method,
				Color:       color,
				Footer: embedFooter{
					Text:    "dev",
					IconURL: "https://ericz.dev/logo.png",
				},
			},
		},
	}
	str, err := json.Marshal(discordEmbed)
	if err != nil {
		log.Print(err)
	}
	go SendDiscordWebhook(str)

}

func SendDiscordWebhook(str []byte) {
	http.Post(os.Getenv("LG_DISCORD_HOOK"), "application/json", bytes.NewBuffer(str))
}
