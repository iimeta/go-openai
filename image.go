package openai

import (
	"bytes"
	"context"
	"mime/multipart"
	"net/http"
	"os"
	"strconv"
)

// Image sizes defined by the OpenAI API.
const (
	CreateImageSize256x256   = "256x256"
	CreateImageSize512x512   = "512x512"
	CreateImageSize1024x1024 = "1024x1024"
	// dall-e-3 supported only.
	CreateImageSize1792x1024 = "1792x1024"
	CreateImageSize1024x1792 = "1024x1792"
)

const (
	CreateImageResponseFormatURL     = "url"
	CreateImageResponseFormatB64JSON = "b64_json"
)

const (
	CreateImageModelDallE2 = "dall-e-2"
	CreateImageModelDallE3 = "dall-e-3"
)

const (
	CreateImageQualityHD       = "hd"
	CreateImageQualityStandard = "standard"
)

const (
	CreateImageStyleVivid   = "vivid"
	CreateImageStyleNatural = "natural"
)

// ImageRequest represents the request structure for the image API.
type ImageRequest struct {
	Prompt            string `json:"prompt,omitempty"`
	Background        string `json:"background,omitempty"`
	Model             string `json:"model,omitempty"`
	Moderation        string `json:"moderation,omitempty"`
	N                 int    `json:"n,omitempty"`
	OutputCompression int    `json:"output_compression,omitempty"`
	OutputFormat      string `json:"output_format,omitempty"`
	Quality           string `json:"quality,omitempty"`
	ResponseFormat    string `json:"response_format,omitempty"`
	Size              string `json:"size,omitempty"`
	Style             string `json:"style,omitempty"`
	User              string `json:"user,omitempty"`
}

// ImageResponse represents a response structure for image API.
type ImageResponse struct {
	Created int64                    `json:"created,omitempty"`
	Data    []ImageResponseDataInner `json:"data,omitempty"`
	Usage   Usage                    `json:"usage,omitempty"`

	httpHeader
}

// ImageResponseDataInner represents a response data structure for image API.
type ImageResponseDataInner struct {
	URL           string `json:"url,omitempty"`
	B64JSON       string `json:"b64_json,omitempty"`
	RevisedPrompt string `json:"revised_prompt,omitempty"`
}

// CreateImage - API call to create an image. This is the main endpoint of the DALL-E API.
func (c *Client) CreateImage(ctx context.Context, request ImageRequest) (response ImageResponse, err error) {
	urlSuffix := "/images/generations"
	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL(urlSuffix, withModel(request.Model)),
		withBody(request),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// ImageEditRequest represents the request structure for the image API.
type ImageEditRequest struct {
	Image          []*multipart.FileHeader `json:"image,omitempty"`
	Prompt         string                  `json:"prompt,omitempty"`
	Background     string                  `json:"background,omitempty"`
	Mask           *multipart.FileHeader   `json:"mask,omitempty"`
	Model          string                  `json:"model,omitempty"`
	N              int                     `json:"n,omitempty"`
	Quality        string                  `json:"quality,omitempty"`
	ResponseFormat string                  `json:"response_format,omitempty"`
	Size           string                  `json:"size,omitempty"`
	User           string                  `json:"user,omitempty"`
}

// CreateEditImage - API call to create an image. This is the main endpoint of the DALL-E API.
func (c *Client) CreateEditImage(ctx context.Context, request ImageEditRequest) (response ImageResponse, err error) {

	body := &bytes.Buffer{}
	builder := c.createFormBuilder(body)

	if err = builder.WriteField("model", request.Model); err != nil {
		return
	}

	if len(request.Image) > 0 {
		if len(request.Image) == 1 {
			if err = builder.CreateFormFileHeader("image", request.Image[0]); err != nil {
				return
			}
		} else {
			for _, image := range request.Image {
				if err = builder.CreateFormFileHeader("image[]", image); err != nil {
					return
				}
			}
		}
	}

	if err = builder.WriteField("prompt", request.Prompt); err != nil {
		return
	}

	if request.Background != "" {
		if err = builder.WriteField("background", request.Background); err != nil {
			return
		}
	}

	if request.Mask != nil {
		if err = builder.CreateFormFileHeader("mask", request.Mask); err != nil {
			return
		}
	}

	if request.N != 0 {
		if err = builder.WriteField("n", strconv.Itoa(request.N)); err != nil {
			return
		}
	}

	if request.Quality != "" {
		if err = builder.WriteField("quality", request.Quality); err != nil {
			return
		}
	}

	if request.ResponseFormat != "" {
		if err = builder.WriteField("response_format", request.ResponseFormat); err != nil {
			return
		}
	}

	if request.Size != "" {
		if err = builder.WriteField("size", request.Size); err != nil {
			return
		}
	}

	if request.User != "" {
		if err = builder.WriteField("user", request.User); err != nil {
			return
		}
	}

	if err = builder.Close(); err != nil {
		return
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL("/images/edits"),
		withBody(body),
		withContentType(builder.FormDataContentType()),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}

// ImageVariRequest represents the request structure for the image API.
type ImageVariRequest struct {
	Image          *os.File `json:"image,omitempty"`
	Model          string   `json:"model,omitempty"`
	N              int      `json:"n,omitempty"`
	Size           string   `json:"size,omitempty"`
	ResponseFormat string   `json:"response_format,omitempty"`
}

// CreateVariImage - API call to create an image variation. This is the main endpoint of the DALL-E API.
// Use abbreviations(vari for variation) because ci-lint has a single-line length limit ...
func (c *Client) CreateVariImage(ctx context.Context, request ImageVariRequest) (response ImageResponse, err error) {
	body := &bytes.Buffer{}
	builder := c.createFormBuilder(body)

	// image
	err = builder.CreateFormFile("image", request.Image)
	if err != nil {
		return
	}

	err = builder.WriteField("n", strconv.Itoa(request.N))
	if err != nil {
		return
	}

	err = builder.WriteField("size", request.Size)
	if err != nil {
		return
	}

	err = builder.WriteField("response_format", request.ResponseFormat)
	if err != nil {
		return
	}

	err = builder.Close()
	if err != nil {
		return
	}

	req, err := c.newRequest(
		ctx,
		http.MethodPost,
		c.fullURL("/images/variations", withModel(request.Model)),
		withBody(body),
		withContentType(builder.FormDataContentType()),
	)
	if err != nil {
		return
	}

	err = c.sendRequest(req, &response)
	return
}
