// Package collector provides an html parser that will extract all posible social information
// from a web page using the meta tags and know protocols as OGP, Twitter Cards, etc.
package collector

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"golang.org/x/net/html"
)

const (
	googleBotUA string = "Mozilla/5.0 (compatible; Googlebot/2.1; +http://www.google.com/bot.html)"
)

// Attributes holds all the information related to the page
type Attributes struct {
	SiteName    string      `json:"site_name,omitempty" meta:"og:site_name"`
	URL         string      `json:"url" meta:"og:url"`
	Description string      `json:"description,omitempty" meta:"og:description"`
	Title       string      `json:"title,omitempty" meta:"og:title"`
	Image       string      `json:"image,omitempty" meta:"og:image"`
	Type        string      `json:"type,omitempty" meta:"og:type"`
	FbApp       string      `json:"fbapp,omitempty" meta:"fb:app_id"`
	Robots      string      `json:"robots,omitempty" meta:"robots"`
	VideoTags   []*VideoTag `json:"video_tags,omitempty" meta:"og:video:tag"`
	Videos      []*Video    `json:"videos,omitempty"`
	Apps        []*App      `json:"apps,omitempty"`
	Twitter     *Twitter    `json:"twitter,omitempty"`
	Player      *Player     `json:"player,omitempty"`
	Icons       []*Icons    `json:"icons,omitempty"`
}

// VideoTag attribute holds the tags related to a Video.
type VideoTag string

// Video attribute holds OGP video information.
type Video struct {
	URL         string `json:"url" meta:"og:video:url,og:video"`
	SecureURL   string `json:"secure_url,omitempty" meta:"og:video:secure_url"`
	Type        string `json:"type,omitempty" meta:"og:video:type"`
	Width       string `json:"width,omitempty" meta:"og:video:width"`
	Height      string `json:"height,omitempty" meta:"og:video:height"`
	Duration    string `json:"duration,omitempty" meta:"og:video:duration"`
	ReleaseDate string `json:"release_date,omitempty" meta:"og:video:release_date"`
}

// Twitter attribute has twitter related data.
type Twitter struct {
	Username  string `json:"user_name,omitempty" meta:"twitter:site"`
	UserID    string `json:"user_id,omitempty" meta:"twitter:site:id"`
	Image     string `json:"image,omitempty" meta:"twitter:image"`
	ImageAlt  string `json:"image_alt,omitempty" meta:"twitter:image:alt"`
	Creator   string `json:"creator,omitempty" meta:"twitter:creator"`
	CreatorID string `json:"creator_id,omitempty" meta:"twitter:creator:id"`
}

// App attribute collect the info of an application related to the site.
type App struct {
	Type  string `json:"type,omitempty" meta:"type"`
	Name  string `json:"name,omitempty" meta:"al:*:app_name"`
	AppID string `json:"app_id,omitempty" meta:"al:*:app_id"`
	URL   string `json:"url,omitempty" meta:"al:*:url"`
	Class string `json:"class,omitempty" meta:"al:*:class"`
}

// Player attribute holds the player information if tha page has a video.
type Player struct {
	URL    string `json:"player,omitempty" meta:"twitter:player"`
	Width  string `json:"width,omitempty" meta:"twitter:player:width"`
	Height string `json:"height,omitempty" meta:"twitter:player:height"`
}

// Icons attribute holds all the icons the page define
type Icons struct {
	Size string `json:"size,omitempty" meta:"size"`
	URL  string `json:"url,omitempty" meta:"url"`
}

// MetaHandler is the function that will be called after each meta attribute
// is readed. This only need to be used if want to replace the default function
// used by ExtractMeta function
type MetaHandler func(attrs *Attributes, meta interface{})

// LinkAttrs attribute holds the Links attributes send to the MetaHandler function.
type LinkAttrs struct {
	Rel   string
	Sizes string
	ID    string
	Href  string
}

// MetaAttrs attribute holds the Meta attributes send to the MetaHandler function.
type MetaAttrs struct {
	Name    string
	Content string
}

// Used to parse the application attributes and extract the device
var appRegex = regexp.MustCompile(`^al:(\w*):.*$`)

// ParseHTML parse an HTML page and try to extract all the information.
func ParseHTML(stringURL string) (attrs *Attributes, err error) {

	var req *http.Request
	req, err = http.NewRequest("GET", stringURL, nil)
	if err != nil {
		return
	}
	req.Header.Set("User-Agent", googleBotUA)

	client := &http.Client{}
	var resp *http.Response
	resp, err = client.Do(req)
	if err != nil {
		return
	}
	defer func() {
		if closeErr := resp.Body.Close(); closeErr != nil {
			err = closeErr
		}
	}()
	if resp.StatusCode != 200 {
		err = fmt.Errorf("error reading URL: %s ~ (%d) %s", stringURL, resp.StatusCode, resp.Status)
		return
	}

	attrs, err = ExtractMeta(resp.Body, AttrsHandler)
	return
}

// ParseFile parse and HTML file and try to extract all the information.
func ParseFile(filename string) (*Attributes, error) {

	data, err := ioutil.ReadFile(filename)
	if err != nil {
		return nil, err
	}

	r := bytes.NewReader(data)
	return ExtractMeta(r, AttrsHandler)
}

// ExtractMeta tokenize the imput stream has an HTML stream and try to parse it.
// For each meta token found will call a handler function to process the information.
// Base code shameless borrowed from : https://github.com/rojters/opengraph/blob/master/opengraph.go
func ExtractMeta(doc io.Reader, handler MetaHandler) (*Attributes, error) {

	attrs := &Attributes{}
	z := html.NewTokenizer(doc)

	for {
		tt := z.Next()
		if tt == html.ErrorToken {
			if z.Err() == io.EOF {
				return attrs, nil
			}
			return attrs, z.Err()
		}

		t := z.Token()

		if t.Data == "head" && t.Type == html.EndTagToken {
			return attrs, nil
		}

		if t.Data == "meta" {
			var prop, cont, name string
			for _, a := range t.Attr {
				switch a.Key {
				case "property":
					prop = strings.TrimSpace(a.Val)
				case "name":
					name = strings.TrimSpace(a.Val)
				case "content":
					cont = strings.TrimSpace(a.Val)
				}
			}
			if prop == "" {
				prop = name
			}

			if prop == "" || cont == "" {
				continue
			}

			handler(attrs, &MetaAttrs{
				Name:    prop,
				Content: cont,
			})
		}

		if t.Data == "link" {
			attr := &LinkAttrs{}

			for _, a := range t.Attr {
				switch a.Key {
				case "rel":
					attr.Rel = strings.TrimSpace(a.Val)
				case "sizes":
					attr.Sizes = strings.TrimSpace(a.Val)
				case "href":
					attr.Href = strings.TrimSpace(a.Val)
				case "id":
					attr.ID = strings.TrimSpace(a.Val)
				}
			}

			if strings.Contains(attr.Rel, "icon") || (strings.Contains(attr.Rel, "alternate") && !strings.Contains(attr.Rel, "stylesheet")) {
				handler(attrs, attr)
			}
		}

	}
}

// AttrsHandler is the default AttributesHandler implementation
func AttrsHandler(attrs *Attributes, meta interface{}) {
	switch meta.(type) {
	case *MetaAttrs:
		metaHandler(attrs, meta.(*MetaAttrs))
	case *LinkAttrs:
		linkHandler(attrs, meta.(*LinkAttrs))
	}
}

func metaHandler(attrs *Attributes, meta *MetaAttrs) {

	if strings.HasPrefix(meta.Name, "og:") || strings.HasPrefix(meta.Name, "fb:") {
		if strings.HasPrefix(meta.Name, "og:video") {
			if meta.Name == "og:video:tag" {
				updateAttributes(attrs, meta)
			}
			if (meta.Name == "og:video") || (meta.Name == "og:video:url") || (meta.Name == "og:video:secure_url") {
				// Check if the last Videos entry already has an url
				// if already has one, this is a new entry
				nvideos := len(attrs.Videos) - 1
				if nvideos >= 0 {
					if (attrs.Videos[nvideos].URL != "") && (attrs.Videos[nvideos].SecureURL != "") {
						attrs.Videos = append(attrs.Videos, &Video{})
					}
				} else {
					// Videos is empty and this is the first tag we see
					attrs.Videos = append(attrs.Videos, &Video{})
				}
			}
			updateAttributes(attrs.Videos, meta)

		} else {
			updateAttributes(attrs, meta)
		}
	} else if strings.HasPrefix(meta.Name, "al:") {
		rgx := appRegex.FindSubmatch([]byte(meta.Name))
		if len(rgx) > 0 {

			napps := len(attrs.Apps) - 1
			appType := string(rgx[1])
			// Chekc if the las Apps entry is from the same kind
			// if not, create a new entry
			if napps >= 0 {
				if attrs.Apps[napps].Type != appType {
					attrs.Apps = append(attrs.Apps, &App{Type: appType})
				}
			} else {
				// if is empty, add a new entry
				attrs.Apps = append(attrs.Apps, &App{Type: appType})
			}

			// Normalize store ids
			appParts := strings.Split(meta.Name, ":")
			if (appParts[2] == "app_store_id") || (appParts[2] == "package") || (appParts[2] == "app_id") {
				meta.Name = "al:*:app_id"
			} else {
				meta.Name = "al:*:" + appParts[2]
			}

			updateAttributes(attrs.Apps, meta)
		}

	} else if strings.HasPrefix(meta.Name, "twitter:") {
		if attrs.Twitter == nil {
			attrs.Twitter = &Twitter{}
		}
		if strings.HasPrefix(meta.Name, "twitter:player") {
			if attrs.Player == nil {
				attrs.Player = &Player{}
			}
			updateAttributes(attrs.Player, meta)
		} else {
			updateAttributes(attrs.Twitter, meta)
		}
	} else {
		updateAttributes(attrs, meta)
	}
}

func updateAttributes(obj interface{}, meta *MetaAttrs) {

	var st reflect.Value

	ar := reflect.TypeOf(obj)
	if ar.Kind() == reflect.Ptr {
		// deference if is a pointer
		ar = ar.Elem()
	}

	if ar.Kind() == reflect.Slice {
		slice := reflect.ValueOf(obj)
		if slice.Kind() == reflect.Ptr {
			// deference the pointer if exists
			slice = slice.Elem()
		}
		last := slice.Index(slice.Len() - 1)
		st = last.Elem()

	} else if ar.Kind() == reflect.Struct {

		st = reflect.ValueOf(obj)
		if st.Kind() == reflect.Ptr {
			// deference the pointer
			st = st.Elem()
		}

	} else {
		fmt.Printf("!>!>!> %v\n", ar.Kind().String())
	}
	typeOf := st.Type()

	// go throw the structure
	for i := 0; i < st.NumField(); i++ {
		field := typeOf.Field(i)
		if tag, ok := field.Tag.Lookup("meta"); ok {
			metaVals := strings.Split(tag, ",")
			for _, mval := range metaVals {
				if mval == meta.Name {
					valField := st.Field(i)
					if valField.Kind() == reflect.String {
						valField.SetString(meta.Content)
					} else if valField.Kind() == reflect.Slice {
						// get the type of the slice
						vt := valField.Type().Elem()
						// if is a pointer, deference it
						if vt.Kind() == reflect.Ptr {
							vt = vt.Elem()
						}
						newVal := reflect.New(vt).Elem()
						cast := reflect.ValueOf(meta.Content).Convert(vt)
						newVal.Set(cast)
						jose := reflect.Append(valField, newVal.Addr())
						valField.Set(jose)

						//valField.Set(reflect.Append(valField, reflect.ValueOf(meta.Content)))
					}
				}
			}
		}
	}
}

func linkHandler(c *Attributes, link *LinkAttrs) {

	if strings.Contains(link.Rel, "icon") {
		if link.Sizes == "" {
			if link.Rel == "mask-icon" {
				link.Sizes = "mask-icon"
			} else if strings.Contains(link.Rel, "shortcut") {
				link.Sizes = "16x16"
			} else {
				link.Sizes = "unknown"
			}
		}

		icn := Icons{
			Size: link.Sizes,
			URL:  link.Href,
		}

		c.Icons = append(c.Icons, &icn)
	}

	// if strings.Contains(link.Rel, "alternate") {
	// 	fmt.Printf("%#v\n", link)
	// }
}
