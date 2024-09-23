package main

import (
	"os"
)

var (
    configErrors []string
)

type Config struct {
    DiscordToken string
    DiscordGuildID string
    KaminoUser string
    KaminoPass string 
    KaminoURL string
    KaminoLoginEndpoint string
    KaminoGetPodsEndpoint string
    KaminoDeleteEndpoint string
    KaminoBulkDeleteEndpoint string
    KaminoRefreshTemplatesEndpoint string
}

func ReadConfig(conf *Config) {
    conf.KaminoUser = os.Getenv("KAMINO_USER")
    conf.KaminoPass = os.Getenv("KAMINO_PASS")
    conf.KaminoURL = os.Getenv("KAMINO_URL")
    conf.KaminoLoginEndpoint = os.Getenv("KAMINO_LOGIN_ENDPOINT")
    conf.KaminoGetPodsEndpoint = os.Getenv("KAMINO_GET_PODS_ENDPOINT")
    conf.KaminoDeleteEndpoint = os.Getenv("KAMINO_DELETE_ENDPOINT")
    conf.KaminoBulkDeleteEndpoint = os.Getenv("KAMINO_BULK_DELETE_ENDPOINT")
    conf.KaminoRefreshTemplatesEndpoint = os.Getenv("KAMINO_REFRESH_TEMPLATES_ENDPOINT")
}
