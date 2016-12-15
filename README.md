
# MetaCollector #

Is a library and command line tools to extract information from HTML documents.

## Usage ##

```go
  ...
  import mc "github.com/kascote/meta_collector"
  ...

  // From an URL
  attrs, err := mc.ParseHTML(urlString)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

  // From a File
  attrs, err := mc.ParseFile(filename)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

  // Or from an io.Reader
  attrs, err := mc.ExtractMeta(html.Body, mc.AttributesHandler)
	if err != nil {
		fmt.Printf("ERROR: %s\n", err.Error())
		os.Exit(1)
	}

```

The information extracted can be marshaled to this:
```json
{
	"site_name": "YouTube",
	"url": "https://www.youtube.com/watch?v=_70Yp8KPXH8",
	"title": "Pete Hunt - Component Based Style Reuse",
	"image": "https://i.ytimg.com/vi/_70Yp8KPXH8/maxresdefault.jpg",
	"type": "video",
	"fbapp": "87741124305",
	"video_tags": [
		"Open Web",
		"JavaScript",
		"Programming",
		"Open Source",
		"Bocoup",
		"ReactJS"
	],
	"videos": [
		{
			"url": "https://www.youtube.com/embed/_70Yp8KPXH8",
			"secure_url": "https://www.youtube.com/embed/_70Yp8KPXH8",
			"type": "text/html",
			"width": "1280",
			"height": "720"
		},
		{
			"url": "http://www.youtube.com/v/_70Yp8KPXH8?version=3\u0026autohide=1",
			"secure_url": "https://www.youtube.com/v/_70Yp8KPXH8?version=3\u0026autohide=1",
			"type": "application/x-shockwave-flash",
			"width": "1280",
			"height": "720"
		}
	],
	"apps": [
		{
			"type": "ios",
			"name": "YouTube",
			"app_id": "544007664",
			"url": "vnd.youtube://www.youtube.com/watch?v=_70Yp8KPXH8\u0026feature=applinks"
		},
		{
			"type": "android",
			"name": "YouTube",
			"app_id": "com.google.android.youtube",
			"url": "vnd.youtube://www.youtube.com/watch?v=_70Yp8KPXH8\u0026feature=applinks"
		},
		{
			"type": "web",
			"url": "https://www.youtube.com/watch?v=_70Yp8KPXH8\u0026feature=applinks"
		}
	],
	"twitter": {
		"user_name": "@youtube",
		"image": "https://i.ytimg.com/vi/_70Yp8KPXH8/maxresdefault.jpg"
	},
	"player": {
		"player": "https://www.youtube.com/embed/_70Yp8KPXH8",
		"width": "1280",
		"height": "720"
	},
	"icons": [
		{
			"size": "16x16",
			"url": "https://s.ytimg.com/yts/img/favicon-vflz7uhzw.ico"
		},
		{
			"size": "32x32",
			"url": "//s.ytimg.com/yts/img/favicon_32-vfl8NGn4k.png"
		},
		{
			"size": "48x48",
			"url": "//s.ytimg.com/yts/img/favicon_48-vfl1s0rGh.png"
		},
		{
			"size": "96x96",
			"url": "//s.ytimg.com/yts/img/favicon_96-vfldSA3ca.png"
		},
		{
			"size": "144x144",
			"url": "//s.ytimg.com/yts/img/favicon_144-vflWmzoXw.png"
		}
	]
}
```

## Licence ##

MIT
