package models

import (
	"regexp"
)

// domainFilterRegex отсеивает служебные CDN-домены с цифрами в имени.
var domainFilterRegex = regexp.MustCompile(`\d+`)

// InitialDynamicResponse — ответ /api/front/v3/config/initial-dynamic.
type InitialDynamicResponse struct {
	InitialDynamic InitialDynamic `json:"initialDynamic"`
}

// InitialDynamic — конфигурация клиента, выдаваемая сайтом при загрузке.
type InitialDynamic struct {
	Players Players `json:"players"`
}

// Players — настройки плееров и список CDN.
type Players struct {
	CDNConfig []CDNConfig `json:"cdnConfig"`
}

// CDNConfig — описание одного CDN-домена.
type CDNConfig struct {
	Domain string `json:"domain"`
}

// GetDomainsCDN возвращает CDN-домены без цифр в имени.
func (id InitialDynamic) GetDomainsCDN() []string {
	var cleaned = make([]string, 0, len(id.Players.CDNConfig))
	for i := range id.Players.CDNConfig {
		if !domainFilterRegex.MatchString(id.Players.CDNConfig[i].Domain) {
			cleaned = append(cleaned, id.Players.CDNConfig[i].Domain)
		}
	}
	return cleaned
}
