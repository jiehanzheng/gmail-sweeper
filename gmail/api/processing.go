package api

import (
	"encoding/base64"
	"fmt"
	"strings"

	"golang.org/x/net/html"
	"google.golang.org/api/gmail/v1"
)

// GetMessageHeader finds value of header name k in a list of headers hs
func GetMessageHeader(hs []*gmail.MessagePartHeader, k string) (string, error) {
	for _, h := range hs {
		if h.Name == k {
			return h.Value, nil
		}
	}
	return "", fmt.Errorf("GetMessageHeader: header name %v not found", k)
}

func ExtractMessageText(mp *gmail.MessagePart) (string, error) {
	base64ToString := func(b64 string) (string, error) {
		decoded, err := base64.URLEncoding.DecodeString(b64)
		if err != nil {
			return "", fmt.Errorf("ExtractMessageText base64ToString: cannot decode message body: %v", err)
		}
		return string(decoded), nil
	}

	var parsePart func(*gmail.MessagePart) (string, error)
	parsePart = func(mp *gmail.MessagePart) (string, error) {
		switch mp.MimeType {
		case "text/plain":
			return base64ToString(mp.Body.Data)
		case "text/html":
			var text string
			htmlSrc, err := base64ToString(mp.Body.Data)
			if err != nil {
				return "", fmt.Errorf("ExtractMessageText parsePart: %v", err)
			}
			doc, err := html.Parse(strings.NewReader(htmlSrc))
			if err != nil {
				return "", fmt.Errorf("ExtractMessageText parsePart: %v", err)
			}
			var f func(*html.Node)
			f = func(n *html.Node) {
				// <script> and <style> do not contain meaningful text for
				// non-programmer human beings, skip
				if n.Type == html.ElementNode {
					switch n.Data {
					case "style", "script":
						return
					}
				}
				if n.Type == html.TextNode {
					text += " " + n.Data
				}
				for c := n.FirstChild; c != nil; c = c.NextSibling {
					f(c)
				}
			}
			f(doc)
			return text, nil
		case "multipart/alternative", "multipart/mixed":
			for _, p := range mp.Parts {
				txt, err := parsePart(p)

				// Return first successfully-parsed part
				// But, per spec, sender prefers the last successfully-parsed part.
				// I will save that for later since AFAIK reversing a slice needs
				// another loop.
				if err == nil {
					return txt, err
				}
			}
			return "", fmt.Errorf("ExtractMessageText parsePart: unable to understand any part in multipart/*")
		default:
			return "", fmt.Errorf("ExtractMessageText parsePart: unknown MIME type")
		}
	}

	return parsePart(mp)
}
