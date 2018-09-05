package typhoon

import (
	"image"
	"net/http"

	"github.com/otiai10/opengraph"
	"golang.org/x/net/html"
)

const (
	tenkijp = "http://www.tenki.jp/bousai/typhoon/japan_near"
)

// Entry ...
type Entry struct {
	NearJP    string
	satisfied bool
	Reference string
}

// GetEntry ...
func GetEntry(httpclient *http.Client) (*Entry, error) {
	res, err := httpclient.Get(tenkijp)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()

	node, err := html.Parse(res.Body)
	if err != nil {
		return nil, err
	}

	entry := &Entry{Reference: tenkijp}
	if err := entry.walk(node); err != nil {
		return nil, err
	}

	return entry, nil
}

// Image ...
func (entry *Entry) Image(httpclient *http.Client) (image.Image, error) {
	res, err := httpclient.Get(entry.NearJP)
	if err != nil {
		return nil, err
	}
	defer res.Body.Close()
	img, _, err := image.Decode(res.Body)
	return img, err
}

func (entry *Entry) walk(node *html.Node) error {

	if entry.satisfied {
		return nil
	}

	if node.Type == html.ElementNode {
		switch node.Data {
		case "meta":
			meta := opengraph.MetaTag(node)
			if meta.Property == "og:image" {
				entry.NearJP = meta.Content
				entry.satisfied = true
				return nil
			}
		default:
			// pass
		}
	}

	for child := node.FirstChild; child != nil; child = child.NextSibling {
		if err := entry.walk(child); err != nil {
			return err
		}
	}

	return nil
}
