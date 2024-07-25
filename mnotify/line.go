package mnotify

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"image"
	"io"
	"mime/multipart"
	"net/http"
	"time"
)

/*
Get access token
https://notify-bot.line.me/my/

Doc API
https://notify-bot.line.me/doc/en/

Sticker list
https://developers.line.biz/en/docs/messaging-api/sticker-list/
*/

var (
	ErrAccessToken           = "invalid access token"
	ErrImageDecode           = "decode image error"
	ErrMessageEmpty          = "message: must not be empty"
	ErrStickerPackageIdEmpty = "You should pass stickerPackageId parameter if you pass stickerId."
	ErrStickerIdEmpty        = "You should pass stickerId parameter if you pass stickerPackageId."
	ErrImageFullsizeEmpty    = "You should pass imageFullsize parameter if you pass imageThumbnail."
)

type LineMessageModel struct {
	ImageFile            []byte
	StickerPackageId     int
	StickerId            int
	NotificationDisabled bool
}

type LineResponse struct {
	Status  int    `json:"status"`
	Message string `json:"message"`
}

type Line struct {
	Authorization string
}

func NewLine(accesstoken string) Line {
	return Line{accesstoken}
}

func (l Line) Send(message string, options *LineMessageModel) (result LineResponse, err error) {
	// validate message
	if message == "" {
		return result, fmt.Errorf(ErrMessageEmpty)
	}

	// create payload
	payload := &bytes.Buffer{}
	writer := multipart.NewWriter(payload)
	_ = writer.WriteField("message", message)

	// validate sticker
	if options != nil {
		if options.StickerId != 0 {
			if options.StickerPackageId == 0 {
				return result, fmt.Errorf(ErrStickerPackageIdEmpty)
			}
			_ = writer.WriteField("stickerId", fmt.Sprint(options.StickerId))
		}
		if options.StickerPackageId != 0 {
			if options.StickerId == 0 {
				return result, fmt.Errorf(ErrStickerIdEmpty)
			}
			_ = writer.WriteField("stickerPackageId", fmt.Sprint(options.StickerPackageId))
		}

		// validate file image : Supported image format is png and jpeg.
		if options.ImageFile != nil && len(options.ImageFile) > 0 {
			_, typeImg, _ := image.Decode(bytes.NewBuffer(options.ImageFile))
			if err != nil {
				return result, fmt.Errorf(ErrImageDecode)
			}
			fw, _ := writer.CreateFormFile("imageFile", "image."+typeImg)
			fw.Write(options.ImageFile)
		}
	}

	writer.Close()

	client, _ := http.NewRequest("POST", "https://notify-api.line.me/api/notify", payload)
	client.Header.Set("Content-Type", writer.FormDataContentType())
	client.Header.Add("Authorization", "Bearer "+l.Authorization)

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client = client.WithContext(ctx)

	req := &http.Client{}
	res, err := req.Do(client)
	if err != nil {
		return LineResponse{500, ""}, err
	}
	defer res.Body.Close()

	resbody, _ := io.ReadAll(res.Body)
	json.Unmarshal(resbody, &result)
	return LineResponse{result.Status, result.Message}, nil
}

func (l Line) CheckStatus() (status bool, err error) {
	// validate an access token
	client, _ := http.NewRequest("GET", "https://notify-api.line.me/api/status", nil)
	client.Header.Add("Authorization", "Bearer "+l.Authorization)

	// set timeout
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	client = client.WithContext(ctx)

	req := &http.Client{}
	res, err := req.Do(client)
	if err != nil {
		return status, err
	}
	defer res.Body.Close()

	resbody, _ := io.ReadAll(res.Body)
	result := LineResponse{}
	json.Unmarshal(resbody, &result)

	// check status code of invalid access token
	if result.Status != 200 {
		return status, fmt.Errorf(ErrAccessToken)
	}

	return true, nil
}
