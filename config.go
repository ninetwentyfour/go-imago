package main

import "os"

const (
	ConRootUrl = "/"
	// ConUrl     = "/{width:[0-9]+}/{height:[0-9]+}/{url}/{format:[a-z]+}"
	ConUrl         = "/{width:[0-9]+}/{height:[0-9]+}/{url}/"
	ConMethod      = "GET"
	ConCacheLength = 1209600
)

var (
	ConRedisUrl            = os.Getenv("IMAGO_GO_REDIS")
	ConMaxRedisActive      = 50
	ConMaxRedisIdle        = 10
	ConS3Key               = os.Getenv("S3_KEY")
	ConS3Secret            = os.Getenv("S3_SECRET")
	ConS3Bucket            = os.Getenv("IMAGO_S3_BUCKET")
	ConImageUrl            = os.Getenv("IMAGO_BASE_LINK_URL")
	ConPort                = ":6001"
	ConRateLimitLimit      = 10
	ConRateLimitTimeout    = 5 // how long before the limit attempts resets
	ConNotFoundLink        = os.Getenv("IMAGO_BASE_LINK_URL") + "not_found.jpg"
	ConImageQuality        = 90
	ConWkhtmltoimageBinary = os.Getenv("WKHTMLTOIMAGE_PATH")
	ConNewRelicKey         = os.Getenv("IMAGO_GO_NEW_RELIC")
)
