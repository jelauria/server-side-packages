package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

const headerCORS = "Access-Control-Allow-Origin"
const corsAnyOrig = "*"

//PreviewImage represents a preview image for a page
type PreviewImage struct {
	URL       string `json:"url,omitempty"`
	SecureURL string `json:"secureURL,omitempty"`
	Type      string `json:"type,omitempty"`
	Width     int    `json:"width,omitempty"`
	Height    int    `json:"height,omitempty"`
	Alt       string `json:"alt,omitempty"`
}

//PageSummary represents summary properties for a web page
type PageSummary struct {
	Type        string          `json:"type,omitempty"`
	URL         string          `json:"url,omitempty"`
	Title       string          `json:"title,omitempty"`
	SiteName    string          `json:"siteName,omitempty"`
	Description string          `json:"description,omitempty"`
	Author      string          `json:"author,omitempty"`
	Keywords    []string        `json:"keywords,omitempty"`
	Icon        *PreviewImage   `json:"icon,omitempty"`
	Images      []*PreviewImage `json:"images,omitempty"`
}

//SummaryHandler handles requests for the page summary API.
//This API expects one query string parameter named `url`,
//which should contain a URL to a web page. It responds with
//a JSON-encoded PageSummary struct containing the page summary
//meta-data.
func SummaryHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Add(headerCORS, corsAnyOrig)

	rQuery := r.URL.Query()
	url := rQuery.Get("url")

	if len(url) == 0 {
		http.Error(w, "Url not supplied", http.StatusBadRequest)
		return
	}

	stream, err := fetchHTML(url)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	sumData, err := extractSummary(url, stream)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer stream.Close()

	finSum, err := json.Marshal(sumData)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Add("Content-Type", "application/json")
	w.Write(finSum)
}

//fetchHTML fetches `pageURL` and returns the body stream or an error.
//Errors are returned if the response status code is an error (>=400),
//or if the content type indicates the URL is not an HTML page.
func fetchHTML(pageURL string) (io.ReadCloser, error) {

	resp, err := http.Get(pageURL)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("Provided url was not found")
	}

	if !strings.HasPrefix(resp.Header.Get("Content-Type"), "text/html") {
		return nil, errors.New("Provided url is not a web page")
	}

	return resp.Body, nil
}

//extractSummary tokenizes the `htmlStream` and populates a PageSummary
//struct with the page's summary meta-data.
func extractSummary(pageURL string, htmlStream io.ReadCloser) (*PageSummary, error) {

	tokenizer := html.NewTokenizer(htmlStream)

	var result PageSummary
	imgs := []*PreviewImage{}

	for {
		tokenType := tokenizer.Next()

		if tokenType == html.ErrorToken {
			err := tokenizer.Err()
			if err == io.EOF {
				break
			}
			return nil, errors.New("Error encountered in processing the web page")
		}

		if tokenType == html.EndTagToken {
			token := tokenizer.Token()
			if token.Data == "head" {
				break
			}
		}

		if tokenType == html.StartTagToken || tokenType == html.SelfClosingTagToken {
			token := tokenizer.Token()

			tag := token.Data
			property, name, content := getBasicInfo(token)

			if "title" == tag && len(result.Title) == 0 {
				tokenType = tokenizer.Next()

				if tokenType == html.TextToken {
					result.Title = tokenizer.Token().Data
				}
			}

			if "og:title" == property {
				result.Title = content
			}

			if "og:type" == property {
				result.Type = content
			}

			if "og:url" == property {
				result.URL = content
			}

			if "og:site_name" == property {
				result.SiteName = content
			}

			if "description" == name && len(result.Description) == 0 {
				result.Description = content
			}

			if "og:description" == property {
				result.Description = content
			}

			if "author" == name {
				result.Author = content
			}

			if "keywords" == name {
				keyArray := strings.Split(content, ",")
				for i, s := range keyArray {
					keyArray[i] = strings.TrimSpace(s)
				}
				result.Keywords = keyArray
			}

			if "link" == tag {
				if getRel(token) == "icon" {
					result.Icon = generateIcon(token, pageURL)
				}
			}

			if "og:image" == property {
				var newImg PreviewImage
				imgs = append(imgs, &newImg)
				imgs[len(imgs)-1].URL = resolveLink(content, pageURL)
			}

			if "og:image:secure_url" == property {
				imgs[len(imgs)-1].SecureURL = content
			}

			if "og:image:type" == property {
				imgs[len(imgs)-1].Type = content
			}

			if "og:image:width" == property {
				newWidth, _ := strconv.Atoi(content)
				imgs[len(imgs)-1].Width = newWidth
			}

			if "og:image:height" == property {
				newHeight, _ := strconv.Atoi(content)
				imgs[len(imgs)-1].Height = newHeight
			}

			if "og:image:alt" == property {
				imgs[len(imgs)-1].Alt = content
			}
		}
	}

	if len(imgs) > 0 {
		result.Images = imgs
	}
	return &result, nil
}

func getBasicInfo(t html.Token) (property string, name string, content string) {
	for _, a := range t.Attr {
		if a.Key == "property" {
			property = a.Val
		}
		if a.Key == "name" {
			name = a.Val
		}
		if a.Key == "content" {
			content = a.Val
		}
	}
	return
}

func getRel(t html.Token) (name string) {
	for _, a := range t.Attr {
		if a.Key == "rel" {
			name = a.Val
			return
		}
	}
	return ""
}

func generateIcon(t html.Token, pageURL string) (image *PreviewImage) {
	var result PreviewImage
	for _, a := range t.Attr {
		if a.Key == "href" {
			result.URL = resolveLink(a.Val, pageURL)
		}
		if a.Key == "type" {
			result.Type = a.Val
		}
		if a.Key == "sizes" {
			if a.Val != "any" {
				heiWid := strings.Split(a.Val, "x")
				height, err := strconv.Atoi(heiWid[0])
				if err == nil {
					result.Height = height
				}
				width, err := strconv.Atoi(heiWid[1])
				if err == nil {
					result.Width = width
				}
			}
		}
	}
	return &result
}

func resolveLink(path string, pageURL string) (result string) {
	if !strings.HasPrefix(path, "/") {
		return path
	}
	base, _ := url.Parse(pageURL)
	urlPath, _ := url.Parse(path)
	return base.ResolveReference(urlPath).String()
}
