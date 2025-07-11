package config

import "web_utilidades/app/utils"

var MainLayoutFilePath string = "resources/layouts/app.html"
var appName string = utils.GetEnv("APP_NAME")
var appEnv string = utils.GetEnv("APP_ENV")
var appKey string = utils.GetEnv("APP_KEY")
var appDebug bool = utils.GetEnv("APP_DEBUG") == "true"
var appTimezone string = utils.GetEnv("APP_TIMEZONE")
var appUrl string = utils.GetEnv("APP_URL")
