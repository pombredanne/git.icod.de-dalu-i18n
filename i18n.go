package i18n

import (
	"log"
	"github.com/nicksnyder/go-i18n/i18n/bundle"
	"net/http"
	"context"
)

type Config struct {
	DefaultLanguage string
	Files           []string          // files to load
	FilesBytes      map[string][]byte // or slices of []bytes with the embedded file data
	Debug           bool
	URLParam        string
	bundle          *bundle.Bundle
}

type I18nMiddleware struct {
	config Config
}

func New(c Config) *I18nMiddleware {
	if c.DefaultLanguage == "" {
		log.Fatal("i18n: No default language set")
	}
	if len(c.Files) == 0 && len(c.FilesBytes) == 0 {
		log.Fatal("i18n: You need to supply either Config.Files and|or Config.FileBytes for language files to be loaded|parsed")
	}

	b := bundle.New()
	c.bundle = b

	for _, file := range c.Files {
		if e := b.LoadTranslationFile(file); e != nil {
			log.Fatal("i18n:", e.Error())
		}
	}
	for s, by := range c.FilesBytes {
		if e := b.ParseTranslationFileBytes(s, by); e != nil {
			log.Fatal("i18n:", e.Error())
		}
	}

	if c.Debug {
		log.Println("i18n: Loaded languages")
		log.Println(b.Translations())
	}
	return &I18nMiddleware{config: c}
}

func (i *I18nMiddleware) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		bycookie := false
		lang := r.URL.Query().Get(i.config.URLParam)
		rlang := r.Header.Get("Accept-Language")
		if lang == "" {
			lc ,e := r.Cookie("lang")
			if e != nil {
				lang = ""
			} else {
				lang = lc.Value
				bycookie = true
			}
		}
		if !bycookie {
			http.SetCookie(w, &http.Cookie{HttpOnly:true, Name: "lang", Value:lang})
		}
		ctx0 := context.WithValue(r.Context(), "i18nlang", lang)
		ctx1 := context.WithValue(ctx0, "i18nrlang", rlang)
		ctx2 := context.WithValue(ctx1, "i18ndlang", i.config.DefaultLanguage)
		ctx3 := context.WithValue(ctx2, "i18nTfunc", i.config.bundle.MustTfunc(lang, rlang, i.config.DefaultLanguage))
		next.ServeHTTP(w, r.WithContext(ctx3))
	})
}