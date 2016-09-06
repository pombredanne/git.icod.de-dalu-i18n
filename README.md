# github.com/dalu/i18n

Go middleware that utilizes https://github.com/nicksnyder/go-i18n

## Required

Go 1.7 for context-aware http Package

## License

modified BSD2

## Issues

Use the Github issue tracker

### Example

```go
package main

import (
	"net/http"

	"github.com/dalu/i18n"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
)

func main() {
	var hrHR = []byte(`[
  {
    "id": "hello",
    "translation": "Pozdrav na {{.Lang}}"
  },
  {
    "id": "coins",
    "translation":{
      "one": "Imas {{.Count}} kunu",
      "few": "Imas {{.Count}} kuna",
      "other": "Imas {{.Count}} kuna"
    }
  }
]`)

	fb := make(map[string][]byte)
	fb["hr-hr.all.json"] = hrHR

	imw := i18n.New(i18n.Config{
		DefaultLanguage: "en-us",
		Files:           []string{"files/en-us.all.json", "files/de-de.all.json"},
		FilesBytes:      fb,
		Debug:           true,
		URLParam:        "lang",
	})

	http.Handle("/", imw.Middleware(http.HandlerFunc(indexHandler)))
	http.ListenAndServe(":8080", nil)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
	lang := r.Context().Value("i18nlang").(string)
	T := r.Context().Value("i18nTfunc").(bundle.TranslateFunc)
	w.Write([]byte(T("hello", map[string]interface{}{"Lang":lang})))
	w.Write([]byte("\n"))
	w.Write([]byte(T("coins", 1)))
	w.Write([]byte("\n"))
	w.Write([]byte(T("coins", 2)))
	w.Write([]byte("\n"))
	w.Write([]byte(T("coins", 200)))
}
```


files/en-us.all.json
```json
[
  {
    "id": "coins",
    "translation": {
      "one": "You have {{.Count}} coin",
      "other": "You have {{.Count}} coins"
    }
  },
  {
    "id": "hello",
    "translation": "Hello in en-US"
  }
]
```

files/de-de.all.json
```json
[
  {
    "id": "coins",
    "translation": {
      "one": "Du hast {{.Count}} Münze",
      "other": "Du hast {{.Count}} Münzen"
    }
  },
  {
    "id": "hello",
    "translation": "Hallo in {{.Lang}}"
  }
]
```
