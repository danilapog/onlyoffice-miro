package common

import (
	"fmt"
	"strings"

	"golang.org/x/text/language"
)

func Concat(strs ...string) string { return strings.Join(strs, "") }

func ToTemplateLanguage(lang string) string {
	tag, err := language.Parse(lang)
	if err != nil {
		return "en-US"
	}

	base, _ := tag.Base()
	baseStr := base.String()

	region, _ := tag.Region()
	regionStr := region.String()

	if baseStr == "zh" {
		if regionStr == "TW" || regionStr == "HK" {
			return "zh-TW"
		}

		return "zh-CN"
	}

	if regionStr != "" {
		return fmt.Sprintf("%s-%s", baseStr, regionStr)
	}

	if likely, err := language.Parse(lang + "-US"); err == nil {
		if region, _ := likely.Region(); region.String() != "" {
			return fmt.Sprintf("%s-%s", baseStr, region.String())
		}
	}

	return "en-US"
}
