package main

import (
	"crypto/sha1"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/garyburd/redigo/redis"
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/gorilla/mux"
	"github.com/ninetwentyfour/go-imago/Godeps/_workspace/src/github.com/ninetwentyfour/go-wkhtmltoimage"
	"net/http"
	"net/url"
	"regexp"
	"strconv"
)

// {"link":"https://d29sc4udwyhodq.cloudfront.net/309d0aa8f6edbe57b4d09630ec592f57.jpg","website":"http://www.reddit.com"}
type JsonResponse struct {
	Link    string `json:"link"`
	Website string `json:"website"`
}

type ImageParams struct {
	Width     int
	Height    int
	Url       string
	Name      string
	ParsedUrl *url.URL
}

type Page struct {
	Title string
	Body  []byte
}

func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html")
	w.Write(homeHtml())
}

// handle html requests
func htmlHandler(w http.ResponseWriter, r *http.Request) {
	link := ""
	imageParams, err := buildParams(r)
	if err != nil {
		link = ConNotFoundLink
	} else {
		link = getImageLink(r, imageParams)
	}

	w.Header().Set("Content-Type", "text/html")
	w.Write(pageHtml(link, imageParams.Url))
}

// handle image requests
func imageHandler(w http.ResponseWriter, r *http.Request) {
	link := ""
	imageParams, err := buildParams(r)
	if err != nil {
		link = ConNotFoundLink
	} else {
		link = imageParams.Name + ".png"
	}

	// need to make sure to get the bad image if there was an error getting image link
	response, err := GetFromS3(link)
	if err != nil {
		LogError(err.Error())
	}

	w.Header().Set("Content-Type", "image/png")
	w.Header().Set("Cache-Control", "max-age=2592000, no-transform, public")
	w.Header().Set("Expires", "Thu, 29 Sep 2022 01:22:54 GMT+00:00")

	w.Write(response)
}

// handle json requests
func jsonHandler(w http.ResponseWriter, r *http.Request) {
	link := ""
	imageParams, err := buildParams(r)
	if err != nil {
		link = ConNotFoundLink
	} else {
		link = getImageLink(r, imageParams)
	}

	responseStruct := JsonResponse{Link: link, Website: imageParams.Url}
	response, _ := json.Marshal(responseStruct)

	// handle jsonp callback
	callback := getCallBack(r)
	if callback != "" {
		w.Header().Set("Content-Type", "application/javascript")
		fmt.Fprintf(w, "%s(%s)", callback, response)
	} else {
		w.Header().Set("Content-Type", "application/json")
		w.Write(response)
	}
}

func oldHandler(w http.ResponseWriter, r *http.Request) {
	imageType := r.FormValue("format")
	switch imageType {
	case "image":
		imageHandler(w, r)
	case "html":
		htmlHandler(w, r)
	default:
		jsonHandler(w, r)
	}
}

func getImageLink(r *http.Request, imageParams *ImageParams) string {
	redisClient := pool.Get()
	defer redisClient.Close()

	if validateUrl(imageParams) == false {
		return ConNotFoundLink
	}

	// check for the url in redis
	link := ""
	rediscache, err := redis.String(redisClient.Do("GET", imageParams.ParsedUrl.Host))
	if err == nil {
		return rediscache
	}

	rediscache, err = redis.String(redisClient.Do("GET", imageParams.Name))
	if err != nil {
		c := wkhtmltoimage.ImageOptions{BinaryPath: ConWkhtmltoimageBinary, Input: imageParams.Url, Format: "png", Height: 720, Width: 1280, Quality: ConImageQuality}
		out, err := wkhtmltoimage.GenerateImage(&c)
		if err != nil {
			LogError(err.Error())
			return ConNotFoundLink
		}

		send_s3, err := resizeImage(out, imageParams)
		if err != nil {
			LogError(err.Error())
			return ConNotFoundLink
		}

		err = SaveToS3(send_s3, imageParams.Name)
		// err = SaveToS3(out, imageParams.Name)
		if err != nil {
			return ConNotFoundLink
		}

		link = ConImageUrl + imageParams.Name + ".png"

		// save the result to redis and set expiration
		cacheResult(redisClient, imageParams, link)
	} else {
		link = rediscache
	}

	return link
}

func buildParams(r *http.Request) (*ImageParams, error) {
	var uparam, wparam, hparam string
	params := mux.Vars(r)
	if len(params) > 0 {
		uparam, wparam, hparam = params["url"], params["width"], params["height"]
	} else {
		uparam, wparam, hparam = r.FormValue("website"), r.FormValue("width"), r.FormValue("height")
	}

	escaped, err := url.QueryUnescape(uparam)
	if err != nil {
		return new(ImageParams), err
	}
	u, err := url.Parse(escaped)
	if err != nil {
		// fmt.Println("ERROR PARSING URL")
		LogInfo("ERROR PARSING URL")
		// fmt.Println(escaped)
		LogInfo(escaped)
		return new(ImageParams), err
	}
	width, err := strconv.Atoi(wparam)
	if err != nil {
		return new(ImageParams), err
	}
	height, err := strconv.Atoi(hparam)
	if err != nil {
		return new(ImageParams), err
	}
	// hash := generateHash(u.Host, params["width"], params["height"])
	hash := generateHash(uparam, wparam, hparam)
	LogInfo(escaped)
	return &ImageParams{Width: width, Height: height, Url: escaped, Name: hash, ParsedUrl: u}, nil
}

func generateHash(name, width, height string) string {
	hasher := sha1.New()
	hasher.Write([]byte(name + "_" + width + "_" + height))
	return hex.EncodeToString(hasher.Sum(nil))
}

func cacheResult(redisClient redis.Conn, imageParams *ImageParams, response string) {
	redisClient.Do("SET", imageParams.ParsedUrl.Host, response)
	redisClient.Do("EXPIRE", imageParams.ParsedUrl.Host, 1)

	redisClient.Do("SET", imageParams.Name, response)
	redisClient.Do("EXPIRE", imageParams.Name, ConCacheLength)
}

func validateUrl(params *ImageParams) bool {
	if params.ParsedUrl.Scheme == "" {
		params.Url = "http://" + params.Url
	}

	// return govalidator.IsRequestURL(params.Url)

	reg, _ := regexp.Compile(`(http|https):\/\/([\w\-_]+(?:(?:\.[\w\-_]+)+))([\w\-\.,@?^=%&amp;:/~\+#]*[\w\-\@?^=%&amp;/~\+#])?`)
	return reg.MatchString(params.Url)
}

func getCallBack(r *http.Request) string {
	callback := ""
	cbs := []string{"callback", "jscallback", "jsonp", "jsoncallback", "cb"}
	for _, element := range cbs {
		cb := r.FormValue(element)
		if r.FormValue(element) != "" {
			callback = cb
			break
		}
	}
	return callback
}
