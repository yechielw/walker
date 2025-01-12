package translation

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/abenz1267/walker/internal/config"
	"github.com/abenz1267/walker/internal/util"
)

var providers = map[string]Provider{
	"googlefree": &GoogleFree{},
}

type Provider interface {
	Name() string
	Translate(text, src, dest string) string
}

type Translation struct {
	config     config.Translation
	providers  []Provider
	systemLang string
}

func (translation *Translation) Cleanup() {
}

func (translation *Translation) Entries(term string) []util.Entry {
	entries := []util.Entry{}

	src, dest := "auto", translation.systemLang

	splits := strings.Split(term, ">")

	if len(splits) == 2 {
		if len(splits[0]) == 2 {
			src = splits[0]
			term = splits[1]
		} else {
			dest = splits[1]
			term = splits[0]
		}
	}

	if len(splits) == 3 {
		src = splits[0]
		dest = splits[2]
		term = splits[1]
	}

	for _, v := range translation.providers {
		res := v.Translate(term, src, dest)

		if res == "" {
			continue
		}

		entries = append(entries, util.Entry{
			Label:            strings.TrimSpace(res),
			Sub:              "Translation",
			Exec:             "",
			Class:            "translation",
			Matching:         util.AlwaysTop,
			RecalculateScore: true,
			SpecialFunc:      translation.SpecialFunc,
		})
	}

	return entries
}

func (translation *Translation) General() *config.GeneralModule {
	return &translation.config.GeneralModule
}

func (translation *Translation) Refresh() {
	translation.config.IsSetup = !translation.config.Refresh
}

func (translation *Translation) Setup() bool {
	translation.config = config.Cfg.Builtins.Translation
	translation.config.IsSetup = true

	for _, v := range translation.config.Providers {
		if provider, ok := providers[v]; ok {
			translation.providers = append(translation.providers, provider)
		}
	}

	langFull := config.Cfg.Locale

	if langFull == "" {
		langFull = os.Getenv("LANG")

		lang_messages := os.Getenv("LC_MESSAGES")
		if lang_messages != "" {
			langFull = lang_messages
		}

		lang_all := os.Getenv("LC_ALL")
		if lang_all != "" {
			langFull = lang_all
		}

		langFull = strings.Split(langFull, ".")[0]
	}

	translation.systemLang = strings.Split(langFull, "_")[0]

	return true
}

func (translation *Translation) SetupData() {
}

func (translation *Translation) SpecialFunc(args ...interface{}) {
	cmd := exec.Command("wl-copy")
	cmd.Stdin = strings.NewReader(args[0].(string))

	err := cmd.Start()
	if err != nil {
		log.Println(err)
		return
	}
}
