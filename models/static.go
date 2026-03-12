package models

import (
	"fmt"
)

type ResponseStatic struct {
	Static Static `json:"static"`
}

func (rs ResponseStatic) GetMMPOrigin() string {
	return rs.Static.Features.MMPExternalSourceOrigin
}

func (rs ResponseStatic) GetMMPVerstion() string {
	return rs.Static.FeaturesV2.PlayerModuleExternalLoading.MmpVersion
}

func (rs ResponseStatic) GetBaseMMP() string {
	return fmt.Sprintf("%s/%s", rs.GetMMPOrigin(), rs.GetMMPVerstion())
}

type Static struct {
	AltcraftPixelBaseURL             string                  `json:"altcraftPixelBaseUrl"`
	AlternateHostRoot                any                     `json:"alternateHostRoot"`
	AmplitudeTrackingCountryCodes    []string                `json:"amplitudeTrackingCountryCodes"`
	AsianCountriesWordsBanCharacters string                  `json:"asianCountriesWordsBanCharacters"`
	Billing                          Billing                 `json:"billing"`
	BroadcastServerOverride          BroadcastServerOverride `json:"broadcastServerOverride"`
	Centrifugo                       []any                   `json:"centrifugo"`
	ComplianceEmailAddress           string                  `json:"complianceEmailAddress"`
	// Countries                        []Country                     `json:"countries"`
	Debug                         bool                          `json:"debug"`
	DMCA                          DMCA                          `json:"dmca"`
	Env                           string                        `json:"env"`
	ErrorReporter                 ErrorReporter                 `json:"errorReporter"`
	ExternalWidget                ExternalWidget                `json:"externalWidget"`
	FanClub                       FanClub                       `json:"fanClub"`
	Features                      Features                      `json:"features"`
	FeaturesV2                    FeaturesV2                    `json:"featuresV2"`
	FingerprintV2                 FingerprintV2                 `json:"fingerprintV2"`
	GiphyAPIKey                   string                        `json:"giphyApiKey"`
	GoogleAuthURL                 string                        `json:"googleAuthUrl"`
	GoogleClientID                string                        `json:"googleClientId"`
	GuestLimits                   GuestLimits                   `json:"guestLimits"`
	HasVipModels                  bool                          `json:"hasVipModels"`
	Hosts                         Hosts                         `json:"hosts"`
	Languages                     []Language                    `json:"languages"`
	Links                         Links                         `json:"links"`
	MainPersonConfig              MainPersonConfig              `json:"mainPersonConfig"`
	MessagesTranslationLanguages  []string                      `json:"messagesTranslationLanguages"`
	MinWatchTime                  int64                         `json:"minWatchTime"`
	MinimumVideoLengthForCutting  int64                         `json:"minimumVideoLengthForCutting"`
	MlAnalyticsHost               string                        `json:"mlAnalyticsHost"`
	ModelTop                      ModelTop                      `json:"modelTop"`
	Moengage                      Moengage                      `json:"moengage"`
	NotInterested                 NotInterested                 `json:"notInterested"`
	PlatformPlayerMap             PlatformPlayerMap             `json:"platformPlayerMap"`
	PlatformPlayerMapSpecialRules PlatformPlayerMapSpecialRules `json:"platformPlayerMapSpecialRules"`
	PlayerDebug                   bool                          `json:"playerDebug"`
	PornstarEmail                 string                        `json:"pornstarEmail"`
	PornstarLink                  string                        `json:"pornstarLink"`
	PressEmailAddress             string                        `json:"pressEmailAddress"`
	PrivateAnimationTypes         PrivateAnimationTypes         `json:"privateAnimationTypes"`
	PrivateMessages               PrivateMessages               `json:"privateMessages"`
	ProfileLanguagesMaxCount      int64                         `json:"profileLanguagesMaxCount"`
	PublicKey                     string                        `json:"publicKey"`
	RedeemGiftRedirectWebsite     string                        `json:"redeemGiftRedirectWebsite"`
	// Regions                          map[string][]Country          `json:"regions"`
	SnapshotTimeout                 int64                    `json:"snapshotTimeout"`
	SocialVerificationLinks         []SocialVerificationLink `json:"socialVerificationLinks"`
	StreamHostPreconnectURLTemplate string                   `json:"streamHostPreconnectUrlTemplate"`
	StudioDetachTime                int64                    `json:"studioDetachTime"`
	StudioDocsURL                   string                   `json:"studioDocsUrl"`
	SupportLinks                    SupportLinks             `json:"supportLinks"`
	TimeoutAddWatchModel            int64                    `json:"timeoutAddWatchModel"`
	TweetMyShow                     TweetMyShow              `json:"tweetMyShow"`
	UserAlerts                      []any                    `json:"userAlerts"`
	UserLevelsRanking               UserLevelsRanking        `json:"userLevelsRanking"`
	ViewersList                     ViewersList              `json:"viewersList"`
	WebrtcEnabled                   bool                     `json:"webrtcEnabled"`
	WebxrExternalWidgetURL          string                   `json:"webxrExternalWidgetUrl"`
	WhiteLabel                      WhiteLabel               `json:"whiteLabel"`
	WhiteLabelHost                  string                   `json:"whiteLabelHost"`
	WhitelistDomainsRegex           string                   `json:"whitelistDomainsRegex"`
	XhlConfig                       XhlConfig                `json:"xhlConfig"`
}

type Billing struct {
	ModelRefundAvailabilityMinutes int64 `json:"modelRefundAvailabilityMinutes"`
}

type BroadcastServerOverride struct {
	PeerConfig string `json:"peerConfig"`
	StreamURL  string `json:"streamUrl"`
}

type Country struct {
	Code  string `json:"code"`
	Title string `json:"title"`
}

type DMCA struct {
	UrlsLimit int64 `json:"urlsLimit"`
}

type ErrorReporter struct {
	IsClientErrorEnabled bool  `json:"isClientErrorEnabled"`
	IsServerErrorEnabled bool  `json:"isServerErrorEnabled"`
	Percent              int64 `json:"percent"`
}

type ExternalWidget struct {
	Enabled bool   `json:"enabled"`
	URL     string `json:"url"`
}

type FanClub struct {
	ExperiencePerMonthByTier ExperiencePerMonthByTier `json:"experiencePerMonthByTier"`
	Tokens                   Tokens                   `json:"tokens"`
}

type ExperiencePerMonthByTier struct {
	Tier1 int64 `json:"tier1"`
	Tier2 int64 `json:"tier2"`
	Tier3 int64 `json:"tier3"`
}

type Tokens struct {
	OneTime   OneTime `json:"oneTime"`
	Recurrent OneTime `json:"recurrent"`
}

type OneTime struct {
	Tier1 []Tier `json:"tier1"`
	Tier2 []Tier `json:"tier2"`
	Tier3 []Tier `json:"tier3"`
}

type Tier struct {
	Months int64 `json:"months"`
	Tokens int64 `json:"tokens"`
}

type Features struct {
	MMPExternalSourceOrigin                         string                         `json:"MMPExternalSourceOrigin"`
	MMPMetricsCollectorEndpoint                     string                         `json:"MMPMetricsCollectorEndpoint"`
	AbTesting                                       []any                          `json:"abTesting"`
	AddShowsToTopFreeForCountries                   []string                       `json:"addShowsToTopFreeForCountries"`
	AppBots                                         bool                           `json:"appBots"`
	AstroPayBanner                                  bool                           `json:"astroPayBanner"`
	AutoResolutionAvailable                         bool                           `json:"autoResolutionAvailable"`
	BlackFriday2020                                 bool                           `json:"blackFriday2020"`
	BlockedSocketEvents                             []any                          `json:"blockedSocketEvents"`
	BroadcastBalancing                              BroadcastBalancing             `json:"broadcastBalancing"`
	DoppioPlayerConfig                              DoppioPlayerConfig             `json:"doppioPlayerConfig"`
	DoppioPlayerExternalSourceOrigin                string                         `json:"doppioPlayerExternalSourceOrigin"`
	DoppioPlayerLoggingLevel                        string                         `json:"doppioPlayerLoggingLevel"`
	FallbackConfig                                  FallbackConfig                 `json:"fallbackConfig"`
	FrontLog                                        bool                           `json:"frontLog"`
	GeoViewCounter                                  bool                           `json:"geoViewCounter"`
	GeoViewWithoutFavoritesCounter                  bool                           `json:"geoViewWithoutFavoritesCounter"`
	HLSFallback                                     HLSFallback                    `json:"hlsFallback"`
	HLSPlayerBaseConfig                             HLSPlayerBaseConfig            `json:"hlsPlayerBaseConfig"`
	HLSPlayerLiveThumbConfig                        []any                          `json:"hlsPlayerLiveThumbConfig"`
	HLSPlayerLowLatencyConfig                       HLSPlayerLowLatencyConfig      `json:"hlsPlayerLowLatencyConfig"`
	IsAirPlayEnabled                                bool                           `json:"isAirPlayEnabled"`
	IsAutoResolutionEnabled                         bool                           `json:"isAutoResolutionEnabled"`
	IsBeta                                          bool                           `json:"isBeta"`
	IsBroadcastQualityEnabled                       bool                           `json:"isBroadcastQualityEnabled"`
	IsDirectMediaLinksEnabled                       bool                           `json:"isDirectMediaLinksEnabled"`
	IsExternalAutoIdleEnabled                       bool                           `json:"isExternalAutoIdleEnabled"`
	IsKycEnabled                                    bool                           `json:"isKycEnabled"`
	IsMobileSidebarEnabled                          bool                           `json:"isMobileSidebarEnabled"`
	IsModelBroadcastScheduleEnabled                 bool                           `json:"isModelBroadcastScheduleEnabled"`
	IsNewBubblesSliderEnabled                       bool                           `json:"isNewBubblesSliderEnabled"`
	IsPWAIOSEnabled                                 bool                           `json:"isPWAIOSEnabled"`
	IsPlayerErrorEnabled                            bool                           `json:"isPlayerErrorEnabled"`
	IsPlayerMeasurementEnabled                      bool                           `json:"isPlayerMeasurementEnabled"`
	IsSignUpWithoutEmailEnabled                     bool                           `json:"isSignUpWithoutEmailEnabled"`
	IsStudioTableEnabled                            bool                           `json:"isStudioTableEnabled"`
	IsStudioTokensHistoryEnable                     bool                           `json:"isStudioTokensHistoryEnable"`
	IsTelegramEnabled                               bool                           `json:"isTelegramEnabled"`
	IsUnThrottlePercentage                          int64                          `json:"isUnThrottlePercentage"`
	MassMessage                                     MassMessage                    `json:"massMessage"`
	MaxIndexPages                                   int64                          `json:"maxIndexPages"`
	MeasurementSocketEvents                         []string                       `json:"measurementSocketEvents"`
	MovingPersonsTablePageSize                      int64                          `json:"movingPersonsTablePageSize"`
	NewMLDB                                         NewMLDB                        `json:"newMLDb"`
	NewMessenger                                    bool                           `json:"newMessenger"`
	NewSessionEventProbabilityOfSendingAsPercentage int64                          `json:"newSessionEventProbabilityOfSendingAsPercentage"`
	NoPasswordMode                                  bool                           `json:"noPasswordMode"`
	Notifications                                   Notifications                  `json:"notifications"`
	Plasma                                          bool                           `json:"plasma"`
	PopularSnapshotInterval                         any                            `json:"popularSnapshotInterval"`
	PrivateLogRecorder                              PrivateLogRecorder             `json:"privateLogRecorder"`
	RecordPrivateSettingsEnabled                    bool                           `json:"recordPrivateSettingsEnabled"`
	RecordPublicSettingsEnabledForUser              bool                           `json:"recordPublicSettingsEnabledForUser"`
	RelatedModelsByTagsEnabled                      bool                           `json:"relatedModelsByTagsEnabled"`
	ShowGenderHallOfFame                            bool                           `json:"showGenderHallOfFame"`
	ShowHalloween2019BannerPeriod                   []string                       `json:"showHalloween2019BannerPeriod"`
	SiteMirrorURL                                   string                         `json:"siteMirrorUrl"`
	StickersSection1                                []string                       `json:"stickersSection1"`
	StickersSection2                                []string                       `json:"stickersSection2"`
	StreamControllerConfig                          StreamControllerConfig         `json:"streamControllerConfig"`
	StripscoreAndPrivateRatingEnabled               bool                           `json:"stripscoreAndPrivateRatingEnabled"`
	SyncUserContentMetricsByChunks                  SyncUserContentMetricsByChunks `json:"syncUserContentMetricsByChunks"`
	TagsAliases                                     map[string][]string            `json:"tagsAliases"`
	ThrottledPercentage                             int64                          `json:"throttledPercentage"`
	TipMenuTranslation                              TipMenuTranslation             `json:"tipMenuTranslation"`
	Top                                             Top                            `json:"top"`
	TwoFactorSettingsAuthEnabled                    bool                           `json:"twoFactorSettingsAuthEnabled"`
	UserEnjoyedModelsEnabled                        bool                           `json:"userEnjoyedModelsEnabled"`
	UsingUserFavoritesEnabled                       bool                           `json:"usingUserFavoritesEnabled"`
	ViewersRequestDelay                             ViewersRequestDelay            `json:"viewersRequestDelay"`
	WatchModelTimeout                               int64                          `json:"watchModelTimeout"`
	WebRTCBroadcastProtocol                         string                         `json:"webRTCBroadcastProtocol"`
	WebRTCCloudFlareEdges                           []string                       `json:"webRTCCloudFlareEdges"`
	WebRTCOriginTurnServers                         WebRTCTurnServers              `json:"webRTCOriginTurnServers"`
	WebRTCOriginTurnServersPortMap                  WebRTCOriginTurnServersPortMap `json:"webRTCOriginTurnServersPortMap"`
	WebRTCPlayForceTCPForHighResolution             string                         `json:"webRTCPlayForceTCPForHighResolution"`
	WebRTCTurnServers                               WebRTCTurnServers              `json:"webRTCTurnServers"`
	WebRTCTurnServersConfig                         WebRTCTurnServersConfig        `json:"webRTCTurnServersConfig"`
	WebrtcCountryWhiteList                          []any                          `json:"webrtcCountryWhiteList"`
	WebrtcCountryWhiteListForIOS                    []string                       `json:"webrtcCountryWhiteListForIOS"`
}

type BroadcastBalancing struct {
	Enabled     bool   `json:"enabled"`
	ServerModel string `json:"serverModel"`
	ServerUser  string `json:"serverUser"`
}

type DoppioPlayerConfig struct {
	PlaybackSpeedWatcher    PlaybackSpeedWatcher    `json:"playbackSpeedWatcher"`
	PlaybackStateController PlaybackStateController `json:"playbackStateController"`
	TargetBufferIncreaser   TargetBufferIncreaser   `json:"targetBufferIncreaser"`
}

type PlaybackSpeedWatcher struct {
	MaxTargetBufferMultiplier int64 `json:"maxTargetBufferMultiplier"`
}

type PlaybackStateController struct {
	AutoPlay       bool `json:"autoPlay"`
	SyncToLiveEdge bool `json:"syncToLiveEdge"`
}

type TargetBufferIncreaser struct {
	BufferAmountThreshold  int64 `json:"bufferAmountThreshold"`
	BufferAmountToIncrease int64 `json:"bufferAmountToIncrease"`
}

type FallbackConfig struct {
	ModelChat ModelChat `json:"modelChat"`
}

type ModelChat struct {
	Interval int64  `json:"interval"`
	Mode     string `json:"mode"`
	Name     string `json:"name"`
}

type HLSFallback struct {
	DefaultDomainCheckInterval           int64    `json:"defaultDomainCheckInterval"`
	DefaultDomainHealthCheckSuccessCount int64    `json:"defaultDomainHealthCheckSuccessCount"`
	FallbackDomains                      []string `json:"fallbackDomains"`
	HealthyCheckInterval                 int64    `json:"healthyCheckInterval"`
	IsEnabled                            bool     `json:"isEnabled"`
	IsReportingEnabled                   bool     `json:"isReportingEnabled"`
	LookupIntervalMinutes                int64    `json:"lookupIntervalMinutes"`
	MaxFailsCount                        int64    `json:"maxFailsCount"`
	NumberOfConnectionErrors             int64    `json:"numberOfConnectionErrors"`
	NumberOfConnectionSuccess            int64    `json:"numberOfConnectionSuccess"`
	TimeoutMS                            int64    `json:"timeoutMS"`
	UnhealthyCheckInterval               int64    `json:"unhealthyCheckInterval"`
}

type HLSPlayerBaseConfig struct {
	ABRBandWidthUpFactor       float64 `json:"abrBandWidthUpFactor"`
	ABREwmaDefaultEstimate     int64   `json:"abrEwmaDefaultEstimate"`
	AutoStartLoad              bool    `json:"autoStartLoad"`
	Debug                      bool    `json:"debug"`
	FragLoadingMaxRetry        int64   `json:"fragLoadingMaxRetry"`
	FragLoadingMaxRetryTimeout int64   `json:"fragLoadingMaxRetryTimeout"`
	FragLoadingMaxTimeOut      int64   `json:"fragLoadingMaxTimeOut"`
	FragLoadingRetryDelay      int64   `json:"fragLoadingRetryDelay"`
	FragLoadingTimeOut         int64   `json:"fragLoadingTimeOut"`
	LiveBackBufferLength       int64   `json:"liveBackBufferLength"`
	LowLatencyMode             bool    `json:"lowLatencyMode"`
	ManifestLoadingMaxRetry    int64   `json:"manifestLoadingMaxRetry"`
	ManifestLoadingRetryDelay  int64   `json:"manifestLoadingRetryDelay"`
	ManifestLoadingTimeOut     int64   `json:"manifestLoadingTimeOut"`
	MaxBufferLength            int64   `json:"maxBufferLength"`
}

type HLSPlayerLowLatencyConfig struct {
	LowLatencyMode          bool    `json:"lowLatencyMode"`
	MaxLiveSyncPlaybackRate float64 `json:"maxLiveSyncPlaybackRate"`
}

type MassMessage struct {
	CheckUnreadMessagesFromCamPage       bool  `json:"checkUnreadMessagesFromCamPage"`
	CheckUnreadMessagesFromNotifications bool  `json:"checkUnreadMessagesFromNotifications"`
	IsCometForUserEnabled                bool  `json:"isCometForUserEnabled"`
	IsEnabled                            bool  `json:"isEnabled"`
	IsPollingEnabled                     bool  `json:"isPollingEnabled"`
	IsSendingEnabled                     bool  `json:"isSendingEnabled"`
	PollingInterval                      int64 `json:"pollingInterval"`
}

type NewMLDB struct {
	ReadEnabled  bool `json:"readEnabled"`
	WriteEnabled bool `json:"writeEnabled"`
}

type Notifications struct {
	PromotedShow any `json:"promotedShow"`
}

type PrivateLogRecorder struct {
	IsRecordEnabled    bool `json:"isRecordEnabled"`
	IsRecordLogEnabled bool `json:"isRecordLogEnabled"`
}

type StreamControllerConfig struct {
	InitialUpdateNumber int64 `json:"initialUpdateNumber"`
	IsEnabled           bool  `json:"isEnabled"`
	NextUpdateNumber    int64 `json:"nextUpdateNumber"`
	TimeoutUpdateStats  int64 `json:"timeoutUpdateStats"`
}

type SyncUserContentMetricsByChunks struct {
	ChunkSize int64 `json:"chunkSize"`
	Enabled   bool  `json:"enabled"`
}

type TipMenuTranslation struct {
	LanguageDetectionCacheTTl int64 `json:"languageDetectionCacheTTl"`
	TranslationCacheTTL       int64 `json:"translationCacheTtl"`
}

type Top struct {
	Boost Boost `json:"boost"`
}

type Boost struct {
	End        string `json:"end"`
	Multiplier int64  `json:"multiplier"`
	Start      string `json:"start"`
}

type ViewersRequestDelay struct {
	Guest         int64 `json:"guest"`
	GuestInterval int64 `json:"guestInterval"`
	UserInterval  int64 `json:"userInterval"`
}

type WebRTCTurnServers struct {
	Credential         string   `json:"credential"`
	DefaultPort        *int64   `json:"defaultPort,omitempty"`
	DisableForFirefox  bool     `json:"disableForFirefox"`
	IceTransportPolicy string   `json:"iceTransportPolicy"`
	IsEnabled          bool     `json:"isEnabled"`
	Protocols          []string `json:"protocols"`
	Server             string   `json:"server"`
	Username           string   `json:"username"`
	IsEnabledInCam2Cam *bool    `json:"isEnabledInCam2Cam,omitempty"`
	Servers            *string  `json:"servers,omitempty"`
}

type WebRTCOriginTurnServersPortMap struct {
	Servers []any `json:"servers"`
}

type WebRTCTurnServersConfig struct {
	IceServersTemplate IceServersTemplate `json:"iceServersTemplate"`
	Servers            []string           `json:"servers"`
}

type IceServersTemplate struct {
	IceServers         []IceServer `json:"iceServers"`
	IceTransportPolicy string      `json:"iceTransportPolicy"`
}

type IceServer struct {
	Credential string `json:"credential"`
	URL        string `json:"url"`
	Username   string `json:"username"`
}

type FeaturesV2 struct {
	Ab15_Vs30_DaysOldUsersPromo                  Ab                                  `json:"ab_15_vs_30_days_old_users_promo"`
	Ab25_TokensInstead20                         Ab                                  `json:"ab_25_tokens_instead_20"`
	Ab50TkGiveawayNewText                        Ab                                  `json:"ab_50tk_giveaway_new_text"`
	AbAbNavB                                     Ab                                  `json:"ab_ab_nav_b"`
	AbAbPlsB                                     Ab                                  `json:"ab_ab_pls_b"`
	AbAbcNewUndervc                              Ab                                  `json:"ab_abc_new_undervc"`
	AbAbcUltBannerFeatured                       Ab                                  `json:"ab_abc_ult_banner_featured"`
	AbAddFavoriteButtonMobile                    Ab                                  `json:"ab_add_favorite_button_mobile"`
	AbAlwaysStartAutoResolution                  Ab                                  `json:"ab_always_start_auto_resolution"`
	AbAvpColor                                   Ab                                  `json:"ab_avp_color"`
	AbAvpDocumentsOneProvider                    Ab                                  `json:"ab_avp_documents_one_provider"`
	AbAvpFranceXhl                               Ab                                  `json:"ab_avp_france_xhl"`
	AbAvpPixThumb                                Ab                                  `json:"ab_avp_pix_thumb"`
	AbAvpVendorSelector                          Ab                                  `json:"ab_avp_vendor_selector"`
	AbBlurLivestreamOnShutter                    Ab                                  `json:"ab_blur_livestream_on_shutter"`
	AbBuyTokensCountrySelector                   Ab                                  `json:"ab_buy_tokens_country_selector"`
	AbCcAVMethod1                                Ab                                  `json:"ab_cc_av_method_1"`
	AbCcAVMethod2                                Ab                                  `json:"ab_cc_av_method_2"`
	AbCFL                                        Ab                                  `json:"ab_cf_l"`
	AbCompletePurchasePopupV2                    Ab                                  `json:"ab_complete_purchase_popup_v2"`
	AbCurrencySelectorBuyTokensAll20             Ab                                  `json:"ab_currency_selector_buy_tokens_all_20"`
	AbDeduplication                              Ab                                  `json:"ab_deduplication"`
	AbDisableLl                                  Ab                                  `json:"ab_disable_ll"`
	AbDiscountOnPrivateFcBenefit                 Ab                                  `json:"ab_discount_on_private_fc_benefit"`
	AbDiscoveryNavigationBarOnMobileV3           Ab                                  `json:"ab_discovery_navigation_bar_on_mobile_v3"`
	AbEXTCountryWithProxy                        Ab                                  `json:"ab_ext_country_with_proxy"`
	AbFlirtingCategory                           Ab                                  `json:"ab_flirting_category"`
	AbGpayP1                                     Ab                                  `json:"ab_gpay_p1"`
	AbGroupTicketShutterUpd                      Ab                                  `json:"ab_group_ticket_shutter_upd"`
	AbGuestsRec                                  Ab                                  `json:"ab_guests_rec"`
	AbHLSLadderIncrease                          Ab                                  `json:"ab_hls_ladder_increase"`
	AbHLSOverrideDefaultConfig                   Ab                                  `json:"ab_hls_override_default_config"`
	AbHLSOverrideServerPartHoldback              Ab                                  `json:"ab_hls_override_server_part_holdback"`
	AbHLSSequentialMode                          Ab                                  `json:"ab_hls_sequential_mode"`
	AbHLSSmoothResolutionChange                  Ab                                  `json:"ab_hls_smooth_resolution_change"`
	AbImpEnjoyed                                 Ab                                  `json:"ab_imp_enjoyed"`
	AbImpRegStrip                                Ab                                  `json:"ab_imp_reg_strip"`
	AbIosHLSAsDefault                            Ab                                  `json:"ab_ios_hls_as_default"`
	AbIosLlHLS                                   Ab                                  `json:"ab_ios_ll_hls"`
	AbJoinScCommunityBannerDesktop               Ab                                  `json:"ab_join_sc_community_banner_desktop"`
	AbLocalCurrency2NdIteration                  Ab                                  `json:"ab_local_currency_2nd_iteration"`
	AbLocalCurrencySelectorBuyTokens             Ab                                  `json:"ab_local_currency_selector_buy_tokens"`
	AbLvl1PackagePixBrazil                       Ab                                  `json:"ab_lvl1_package_pix_brazil"`
	AbLvl5PackageBuyTokensTier1Countries         Ab                                  `json:"ab_lvl5_package_buy_tokens_tier1_countries"`
	AbMatchedSession                             Ab                                  `json:"ab_matched_session"`
	AbNewEXTCountry                              Ab                                  `json:"ab_new_ext_country"`
	AbNewEXTCountryPart3                         Ab                                  `json:"ab_new_ext_country_part3"`
	AbNewSidescrolling                           Ab                                  `json:"ab_new_sidescrolling"`
	AbOnboardRec                                 Ab                                  `json:"ab_onboard_rec"`
	AbOpenapiContractValidation                  Ab                                  `json:"ab_openapi_contract_validation"`
	AbPrivacyBankStatementUpdate                 Ab                                  `json:"ab_privacy_bank_statement_update"`
	AbPrivateModalActivities                     Ab                                  `json:"ab_private_modal_activities"`
	AbPromSECPosition                            Ab                                  `json:"ab_prom_sec_position"`
	AbRInCat                                     Ab                                  `json:"ab_r_in_cat"`
	AbRecInFeatured                              Ab                                  `json:"ab_rec_in_featured"`
	AbRecommended                                Ab                                  `json:"ab_recommended"`
	AbRedesignAVShutter                          Ab                                  `json:"ab_redesign_av_shutter"`
	AbRelModels                                  Ab                                  `json:"ab_rel_models"`
	AbRelatedV2                                  Ab                                  `json:"ab_related_v2"`
	AbSayHi5Tk                                   Ab                                  `json:"ab_say_hi_5tk"`
	AbSendTipAnimation                           Ab                                  `json:"ab_send_tip_animation"`
	AbSidebarUltimatePromotionNonPayingUsers     Ab                                  `json:"ab_sidebar_ultimate_promotion_non_paying_users"`
	AbSourceRendition                            Ab                                  `json:"ab_source_rendition"`
	AbStartingMessagesOptimization               Ab                                  `json:"ab_starting_messages_optimization"`
	AbStickAvpBanner                             Ab                                  `json:"ab_stick_avp_banner"`
	AbStreamingCDN                               Ab                                  `json:"ab_streaming_cdn"`
	AbStreamingCDN2                              Ab                                  `json:"ab_streaming_cdn_2"`
	AbStreamingCDN3                              Ab                                  `json:"ab_streaming_cdn_3"`
	AbStrips                                     Ab                                  `json:"ab_strips"`
	AbStripsBigPart                              Ab                                  `json:"ab_strips_big_part"`
	AbSyntAbTest1                                Ab                                  `json:"ab_synt_ab_test_1"`
	AbSyntAbTest2                                Ab                                  `json:"ab_synt_ab_test_2"`
	AbSyntAbTest3                                Ab                                  `json:"ab_synt_ab_test_3"`
	AbSyntAbTest4                                Ab                                  `json:"ab_synt_ab_test_4"`
	AbSyntAbTest5                                Ab                                  `json:"ab_synt_ab_test_5"`
	AbTestParallelTesting1                       Ab                                  `json:"ab_test_parallel_testing_1"`
	AbTestParallelTesting10                      Ab                                  `json:"ab_test_parallel_testing_10"`
	AbTestParallelTesting2                       Ab                                  `json:"ab_test_parallel_testing_2"`
	AbTestParallelTesting3                       Ab                                  `json:"ab_test_parallel_testing_3"`
	AbTestParallelTesting4                       Ab                                  `json:"ab_test_parallel_testing_4"`
	AbTestParallelTesting5                       Ab                                  `json:"ab_test_parallel_testing_5"`
	AbTestParallelTesting6                       Ab                                  `json:"ab_test_parallel_testing_6"`
	AbTestParallelTesting7                       Ab                                  `json:"ab_test_parallel_testing_7"`
	AbTestParallelTesting8                       Ab                                  `json:"ab_test_parallel_testing_8"`
	AbTestParallelTesting9                       Ab                                  `json:"ab_test_parallel_testing_9"`
	AbTipPromptInChat                            Ab                                  `json:"ab_tip_prompt_in_chat"`
	AbTranslateMessages                          Ab                                  `json:"ab_translate_messages"`
	AbTranslateMessagesPart2                     Ab                                  `json:"ab_translate_messages_part_2"`
	AbUpdTitleTextBuyTokens                      Ab                                  `json:"ab_upd_title_text_buy_tokens"`
	AbUpgateReskin                               Ab                                  `json:"ab_upgate_reskin"`
	AbVideoPoster                                Ab                                  `json:"ab_video_poster"`
	AbViewcamResolutionLimit                     Ab                                  `json:"ab_viewcam_resolution_limit"`
	AbVRGroupShow                                Ab                                  `json:"ab_vr_group_show"`
	AbWatchedModel                               Ab                                  `json:"ab_watched_model"`
	AbWebrtcFEC                                  Ab                                  `json:"ab_webrtc_fec"`
	AbWebrtcNewABR                               Ab                                  `json:"ab_webrtc_new_abr"`
	AbWebrtcPlayoutDelayHint                     Ab                                  `json:"ab_webrtc_playout_delay_hint"`
	AccessToStreamingApp                         AccessToStreamingApp                `json:"accessToStreamingApp"`
	AccountOwnerSelectorCryptoCosmoEpay          AccountOwnerSelectorCryptoCosmoEpay `json:"accountOwnerSelectorCryptoCosmoEpay"`
	ActivityUpdater                              ActivityUpdater                     `json:"activityUpdater"`
	AllModelsPage                                AccountOwnerSelectorCryptoCosmoEpay `json:"allModelsPage"`
	AllowBuyTokensWithoutEmail                   AccountOwnerSelectorCryptoCosmoEpay `json:"allowBuyTokensWithoutEmail"`
	AllowTubePromo                               AccountOwnerSelectorCryptoCosmoEpay `json:"allowTubePromo"`
	AllowUsersToClearTokenHistory                AccountOwnerSelectorCryptoCosmoEpay `json:"allowUsersToClearTokenHistory"`
	AllowUsersToDeleteChats                      AccountOwnerSelectorCryptoCosmoEpay `json:"allowUsersToDeleteChats"`
	AlphabetCategories                           AlphabetCategories                  `json:"alphabetCategories"`
	AlwaysSendModelStatusChangedEvent            AccountOwnerSelectorCryptoCosmoEpay `json:"alwaysSendModelStatusChangedEvent"`
	AmplitudeCompactAbTestProperties             AccountOwnerSelectorCryptoCosmoEpay `json:"amplitudeCompactAbTestProperties"`
	AmplitudeConfig                              AmplitudeConfig                     `json:"amplitudeConfig"`
	AmplitudeSideEffects                         AccountOwnerSelectorCryptoCosmoEpay `json:"amplitudeSideEffects"`
	AntispamRules                                AntispamRules                       `json:"antispamRules"`
	Apps                                         Apps                                `json:"apps"`
	AppsPerformanceThresholds                    AppsPerformanceThresholds           `json:"appsPerformanceThresholds"`
	AuthMiddlewareRefactoring                    AccountOwnerSelectorCryptoCosmoEpay `json:"authMiddlewareRefactoring"`
	AutoChangeInfoStudios                        AccountOwnerSelectorCryptoCosmoEpay `json:"autoChangeInfoStudios"`
	AutoDiscountPrivates                         AutoDiscountPrivates                `json:"autoDiscountPrivates"`
	AvpConfig                                    AccountOwnerSelectorCryptoCosmoEpay `json:"avpConfig"`
	AvpYotiStatus                                AccountOwnerSelectorCryptoCosmoEpay `json:"avpYotiStatus"`
	BanUserByFingerprintAndIP                    AccountOwnerSelectorCryptoCosmoEpay `json:"banUserByFingerprintAndIP"`
	BannerForJapan                               AccountOwnerSelectorCryptoCosmoEpay `json:"bannerForJapan"`
	BannerMyWebcamRoom                           BannerMyWebcamRoom                  `json:"bannerMyWebcamRoom"`
	BillingConfigsForUserWithoutEmail            AccountOwnerSelectorCryptoCosmoEpay `json:"billingConfigsForUserWithoutEmail"`
	BlurredLivestreamOnShutter                   BlurredLivestreamOnShutter          `json:"blurredLivestreamOnShutter"`
	BoostIncreaseModels                          AccountOwnerSelectorCryptoCosmoEpay `json:"boostIncreaseModels"`
	BroadcastBalancing                           BroadcastBalancing                  `json:"broadcastBalancing"`
	BroadcastDefaultQuality                      BroadcastDefaultQuality             `json:"broadcastDefaultQuality"`
	BroadcastQuality                             BroadcastQuality                    `json:"broadcastQuality"`
	BroadcastQualityLevels                       BroadcastQualityLevels              `json:"broadcastQualityLevels"`
	BroadcastService                             BroadcastService                    `json:"broadcastService"`
	BrowserCache                                 BrowserCache                        `json:"browserCache"`
	C2CLowLatencySwitch                          AccountOwnerSelectorCryptoCosmoEpay `json:"c2cLowLatencySwitch"`
	Centrifugo                                   Centrifugo                          `json:"centrifugo"`
	CentrifugoMetrics                            CentrifugoMetrics                   `json:"centrifugoMetrics"`
	CentrifugoRecovery                           AccountOwnerSelectorCryptoCosmoEpay `json:"centrifugoRecovery"`
	ChangeAPIMethod                              AccountOwnerSelectorCryptoCosmoEpay `json:"changeApiMethod"`
	ChangeCleanupViewersLogic                    AccountOwnerSelectorCryptoCosmoEpay `json:"changeCleanupViewersLogic"`
	CheckUserDescriptionForSpam                  AccountOwnerSelectorCryptoCosmoEpay `json:"checkUserDescriptionForSpam"`
	Comet                                        Comet                               `json:"comet"`
	CommunityReportCreatorV2                     CommunityReportCreatorV2            `json:"communityReportCreatorV2"`
	ConnectModelToStudio                         AccountOwnerSelectorCryptoCosmoEpay `json:"connectModelToStudio"`
	ContentMonitor                               ContentMonitor                      `json:"contentMonitor"`
	ContentMonitorQueueSizeCache                 ContentMonitorQueueSizeCache        `json:"contentMonitorQueueSizeCache"`
	ContentWatchAmplitudeEvent                   AccountOwnerSelectorCryptoCosmoEpay `json:"contentWatchAmplitudeEvent"`
	ContentsReports                              ContentsReports                     `json:"contentsReports"`
	ContinuousScrollInTipMenu                    AccountOwnerSelectorCryptoCosmoEpay `json:"continuousScrollInTipMenu"`
	CouplesGuysTransNewIndexCategoriesMar23      AccountOwnerSelectorCryptoCosmoEpay `json:"couplesGuysTransNewIndexCategoriesMar23"`
	CurrentReleaseVersion                        CurrentReleaseVersion               `json:"currentReleaseVersion"`
	CvCategoryIcon                               AccountOwnerSelectorCryptoCosmoEpay `json:"cvCategoryIcon"`
	CvMassMessageContentCheck                    AccountOwnerSelectorCryptoCosmoEpay `json:"cvMassMessageContentCheck"`
	CvPrivateContentCheck                        AccountOwnerSelectorCryptoCosmoEpay `json:"cvPrivateContentCheck"`
	CvReportHideEmptyStreams                     CvReportHideEmptyStreams            `json:"cvReportHideEmptyStreams"`
	CvReportModalProbability                     CvReportModalProbability            `json:"cvReportModalProbability"`
	CvSearch                                     AccountOwnerSelectorCryptoCosmoEpay `json:"cvSearch"`
	CvSearchConfig                               CvSearchConfig                      `json:"cvSearchConfig"`
	CvSearchExtension                            AccountOwnerSelectorCryptoCosmoEpay `json:"cvSearchExtension"`
	DatadomeBuyTokens                            AccountOwnerSelectorCryptoCosmoEpay `json:"datadomeBuyTokens"`
	DebugMeasurement                             DebugMeasurement                    `json:"debugMeasurement"`
	DeprecatedPages                              DeprecatedPages                     `json:"deprecatedPages"`
	DetectMediaAPITampering                      DetectMediaAPITampering             `json:"detectMediaApiTampering"`
	DiscountOnPrivateFcBenefit                   AccountOwnerSelectorCryptoCosmoEpay `json:"discountOnPrivateFcBenefit"`
	DocumentsMonitorAgeVerificationDocuments     AccountOwnerSelectorCryptoCosmoEpay `json:"documentsMonitorAgeVerificationDocuments"`
	DocumentsMonitorCriticalDuplicates           AccountOwnerSelectorCryptoCosmoEpay `json:"documentsMonitorCriticalDuplicates"`
	DownloadTranslatedAgreementAP                AccountOwnerSelectorCryptoCosmoEpay `json:"downloadTranslatedAgreementAP"`
	DownscaleBroadcastResolution                 DownscaleBroadcastResolution        `json:"downscaleBroadcastResolution"`
	DSAAPIV2                                     AccountOwnerSelectorCryptoCosmoEpay `json:"dsaApiV2"`
	DSARecSetting                                AccountOwnerSelectorCryptoCosmoEpay `json:"dsaRecSetting"`
	DSAReportAndAppeal                           AccountOwnerSelectorCryptoCosmoEpay `json:"dsaReportAndAppeal"`
	DSAReportUser                                DSAReportUser                       `json:"dsaReportUser"`
	DwhAdapter                                   DwhAdapter                          `json:"dwhAdapter"`
	DwhGuestFeatures                             DwhGuestFeatures                    `json:"dwhGuestFeatures"`
	DwhTabKeepAlive                              DwhTabKeepAlive                     `json:"dwhTabKeepAlive"`
	EnableTipMenuDiscount                        AccountOwnerSelectorCryptoCosmoEpay `json:"enableTipMenuDiscount"`
	EnableWebPushiOS                             EnableWebPushiOS                    `json:"enableWebPushiOS"`
	EnhancedMixedTags                            AccountOwnerSelectorCryptoCosmoEpay `json:"enhancedMixedTags"`
	EuDSAPage                                    EuDSAPage                           `json:"euDsaPage"`
	ExitPrivateSound                             AccountOwnerSelectorCryptoCosmoEpay `json:"exitPrivateSound"`
	ExpiredUltimateNotificationModal             AccountOwnerSelectorCryptoCosmoEpay `json:"expiredUltimateNotificationModal"`
	ExtensionsForTicketGroupShows                AccountOwnerSelectorCryptoCosmoEpay `json:"extensionsForTicketGroupShows"`
	FavoritesSectionPromoThumb                   AccountOwnerSelectorCryptoCosmoEpay `json:"favoritesSectionPromoThumb"`
	FingerprintCDN                               FingerprintCDN                      `json:"fingerprintCdn"`
	Fingerprintv2Giveaways                       AccountOwnerSelectorCryptoCosmoEpay `json:"fingerprintv2Giveaways"`
	FixDisclosureAPICamAccess                    AccountOwnerSelectorCryptoCosmoEpay `json:"fixDisclosureApiCamAccess"`
	FlirtingMode                                 FlirtingMode                        `json:"flirtingMode"`
	FlirtingModeCategorySorting                  AccountOwnerSelectorCryptoCosmoEpay `json:"flirtingModeCategorySorting"`
	FraudMonitor                                 Monitor                             `json:"fraudMonitor"`
	FriendsOnlyPMForModels                       AccountOwnerSelectorCryptoCosmoEpay `json:"friendsOnlyPMForModels"`
	FrontLogV2                                   FrontLogV2                          `json:"frontLogV2"`
	Glass2GlassLatencyMeasurement                Glass2GlassLatencyMeasurement       `json:"glass2GlassLatencyMeasurement"`
	GoogleAuthentication                         GoogleAuthentication                `json:"googleAuthentication"`
	GuestLandedAndIdentityTimeout                GuestLandedAndIdentityTimeout       `json:"guestLandedAndIdentityTimeout"`
	HalloweenReaction                            HalloweenReaction                   `json:"halloweenReaction"`
	HandyBuyButton                               AccountOwnerSelectorCryptoCosmoEpay `json:"handyBuyButton"`
	HideAgeForModels                             AccountOwnerSelectorCryptoCosmoEpay `json:"hideAgeForModels"`
	HideChangeUsernameXhl                        AccountOwnerSelectorCryptoCosmoEpay `json:"hideChangeUsernameXhl"`
	HideDeletedAccounts                          AccountOwnerSelectorCryptoCosmoEpay `json:"hideDeletedAccounts"`
	HideDuplicatesCountrySections                AccountOwnerSelectorCryptoCosmoEpay `json:"hideDuplicatesCountrySections"`
	HideFanCentro                                AccountOwnerSelectorCryptoCosmoEpay `json:"hideFanCentro"`
	HighlightInteractions                        AccountOwnerSelectorCryptoCosmoEpay `json:"highlightInteractions"`
	HLSABRConfig                                 HLSABRConfig                        `json:"hlsAbrConfig"`
	HLSEdgeBalancing                             AccountOwnerSelectorCryptoCosmoEpay `json:"hlsEdgeBalancing"`
	HLSPlayer                                    HLSPlayer                           `json:"hlsPlayer"`
	HotjarImprovements                           AccountOwnerSelectorCryptoCosmoEpay `json:"hotjarImprovements"`
	IndexCategoryLoadTime                        IndexCategoryLoadTime               `json:"indexCategoryLoadTime"`
	InpReporter                                  InpReporter                         `json:"inpReporter"`
	InvoiceV2                                    InvoiceV2                           `json:"invoiceV2"`
	Kyb                                          AccountOwnerSelectorCryptoCosmoEpay `json:"kyb"`
	Kyc                                          Kyc                                 `json:"kyc"`
	KycV2                                        AccountOwnerSelectorCryptoCosmoEpay `json:"kycV2"`
	LayersManagerForBuyTokens                    AccountOwnerSelectorCryptoCosmoEpay `json:"layersManagerForBuyTokens"`
	LayersManagerForStartPrivate                 AccountOwnerSelectorCryptoCosmoEpay `json:"layersManagerForStartPrivate"`
	LazyLoadModelThumbs                          LazyLoadModelThumbs                 `json:"lazyLoadModelThumbs"`
	LazyLoadingOnIndex                           LazyLoadingOnIndex                  `json:"lazyLoadingOnIndex"`
	LimitBroadcastHours                          LimitBroadcastHours                 `json:"limitBroadcastHours"`
	LiveThumbs                                   LiveThumbs                          `json:"liveThumbs"`
	LovenseAPIUrls                               LovenseAPIUrls                      `json:"lovenseApiUrls"`
	MainStudios                                  MainStudios                         `json:"mainStudios"`
	MaxAllowedUserDevices                        MaxAllowedUserDevices               `json:"maxAllowedUserDevices"`
	MessageRepositoryUnionQuery                  AccountOwnerSelectorCryptoCosmoEpay `json:"messageRepositoryUnionQuery"`
	MoEngageSupport                              MoEngageSupport                     `json:"moEngageSupport"`
	ModelActivitiesByGender                      AccountOwnerSelectorCryptoCosmoEpay `json:"modelActivitiesByGender"`
	ModelListUltimateBanner                      ModelListUltimateBanner             `json:"modelListUltimateBanner"`
	ModelRecordPublicShow                        ModelRecordPublicShow               `json:"modelRecordPublicShow"`
	ModelSearch                                  ModelSearch                         `json:"modelSearch"`
	ModelsBlackListCountries                     AccountOwnerSelectorCryptoCosmoEpay `json:"modelsBlackListCountries"`
	ModelsDeduplication                          ModelsDeduplication                 `json:"modelsDeduplication"`
	MonthlyTopModelsLimits                       MonthlyTopModelsLimits              `json:"monthlyTopModelsLimits"`
	MoreVisibleReportButton                      AccountOwnerSelectorCryptoCosmoEpay `json:"moreVisibleReportButton"`
	MouflonMediaPlayer                           MouflonMediaPlayer                  `json:"mouflonMediaPlayer"`
	MouflonStatsCollection                       MouflonStatsCollection              `json:"mouflonStatsCollection"`
	MoveAlbumsVideosAPIv2                        AccountOwnerSelectorCryptoCosmoEpay `json:"moveAlbumsVideosAPIv2"`
	NewBdsmCategory                              NewBdsmCategory                     `json:"newBdsmCategory"`
	NewFeedback                                  AccountOwnerSelectorCryptoCosmoEpay `json:"newFeedback"`
	NewImageViewer                               NewImageViewer                      `json:"newImageViewer"`
	NewIndianCategories                          AccountOwnerSelectorCryptoCosmoEpay `json:"newIndianCategories"`
	NewLocaleUpdateLogic                         AccountOwnerSelectorCryptoCosmoEpay `json:"newLocaleUpdateLogic"`
	NewModelTagsCalculation                      AccountOwnerSelectorCryptoCosmoEpay `json:"newModelTagsCalculation"`
	NewWatchTimeLogic                            NewWatchTimeLogic                   `json:"newWatchTimeLogic"`
	NoManualCountryChange                        AccountOwnerSelectorCryptoCosmoEpay `json:"noManualCountryChange"`
	NoReuploadForCritReasons                     AccountOwnerSelectorCryptoCosmoEpay `json:"noReuploadForCritReasons"`
	NonNude                                      NonNude                             `json:"nonNude"`
	NonNudeSnapshotsAPIV2                        AccountOwnerSelectorCryptoCosmoEpay `json:"nonNudeSnapshotsApiV2"`
	NotificationsV2                              AccountOwnerSelectorCryptoCosmoEpay `json:"notificationsV2"`
	ObsLinkToVRViewcam                           AccountOwnerSelectorCryptoCosmoEpay `json:"obsLinkToVrViewcam"`
	OneConfigEndpoint                            AccountOwnerSelectorCryptoCosmoEpay `json:"oneConfigEndpoint"`
	OneCountryTag                                AccountOwnerSelectorCryptoCosmoEpay `json:"oneCountryTag"`
	OnlineModelsDataExporter                     AccountOwnerSelectorCryptoCosmoEpay `json:"onlineModelsDataExporter"`
	Optimization                                 Optimization                        `json:"optimization"`
	PayoutsToCard                                AccountOwnerSelectorCryptoCosmoEpay `json:"payoutsToCard"`
	PerformancePages                             PerformancePages                    `json:"performancePages"`
	PermanentBlockDeletePublicMessages           AccountOwnerSelectorCryptoCosmoEpay `json:"permanentBlockDeletePublicMessages"`
	PersonalLinkDuplicateImages                  AccountOwnerSelectorCryptoCosmoEpay `json:"personalLinkDuplicateImages"`
	PipWindowPositioningAndDisablingState        AccountOwnerSelectorCryptoCosmoEpay `json:"pipWindowPositioningAndDisablingState"`
	PixelateShowsForJapan                        PixelateShowsForJapan               `json:"pixelateShowsForJapan"`
	PlasmaDeleteUIElements                       AccountOwnerSelectorCryptoCosmoEpay `json:"plasmaDeleteUIElements"`
	PlayerModuleExternalLoading                  PlayerModuleExternalLoading         `json:"playerModuleExternalLoading"`
	PolicySidebarAnalytics                       AccountOwnerSelectorCryptoCosmoEpay `json:"policySidebarAnalytics"`
	PreStreamVerification                        AccountOwnerSelectorCryptoCosmoEpay `json:"preStreamVerification"`
	PreviewsSectionInMyDetails                   PreviewsSectionInMyDetails          `json:"previewsSectionInMyDetails"`
	PrivateConversationsReadSocketEvent          AccountOwnerSelectorCryptoCosmoEpay `json:"privateConversationsReadSocketEvent"`
	PrivateToSpyModeTransition                   AccountOwnerSelectorCryptoCosmoEpay `json:"privateToSpyModeTransition"`
	ProfileVisibleToModelsOnly                   AccountOwnerSelectorCryptoCosmoEpay `json:"profileVisibleToModelsOnly"`
	ProhibitSubStudioChanges                     ProhibitSubStudioChanges            `json:"prohibitSubStudioChanges"`
	ProhibitedWords                              ProhibitedWords                     `json:"prohibitedWords"`
	ProhibitedWordsDuplicateRemoval              AccountOwnerSelectorCryptoCosmoEpay `json:"prohibitedWordsDuplicateRemoval"`
	ProhibitedWordsV2                            AccountOwnerSelectorCryptoCosmoEpay `json:"prohibitedWordsV2"`
	PromoThumbFewModelsCountry                   AccountOwnerSelectorCryptoCosmoEpay `json:"promoThumbFewModelsCountry"`
	PromoteHacksSection                          AccountOwnerSelectorCryptoCosmoEpay `json:"promoteHacksSection"`
	PulsePromotionBanners                        PulsePromotionBanners               `json:"pulsePromotionBanners"`
	ReRequestFPv2                                ReRequestFPv2                       `json:"reRequestFPv2"`
	RecordPlayerSwitches                         AccountOwnerSelectorCryptoCosmoEpay `json:"recordPlayerSwitches"`
	RecordStreams                                RecordStreams                       `json:"recordStreams"`
	RecordUserSessionAmplitude                   RecordUserSessionAmplitude          `json:"recordUserSessionAmplitude"`
	RecordVRStreams                              RecordVRStreams                     `json:"recordVrStreams"`
	RedeemDiscountPage                           AccountOwnerSelectorCryptoCosmoEpay `json:"redeemDiscountPage"`
	RefactoredTokensModalPromoBlock              AccountOwnerSelectorCryptoCosmoEpay `json:"refactoredTokensModalPromoBlock"`
	RefillBuyTokens                              InvoiceV2                           `json:"refillBuyTokens"`
	RefreshModelStatus                           RefreshModelStatus                  `json:"refreshModelStatus"`
	RegularAppealCountries                       RegularAppealCountries              `json:"regularAppealCountries"`
	ReloadChatAfterReconnectionForModel          AccountOwnerSelectorCryptoCosmoEpay `json:"reloadChatAfterReconnectionForModel"`
	RemoveMagicWandFromModelSignUp               AccountOwnerSelectorCryptoCosmoEpay `json:"removeMagicWandFromModelSignUp"`
	ReportAvatarsBackgroundPanels                AccountOwnerSelectorCryptoCosmoEpay `json:"reportAvatarsBackgroundPanels"`
	ReportMessagesMobile                         AccountOwnerSelectorCryptoCosmoEpay `json:"reportMessagesMobile"`
	ReportModelsContent                          AccountOwnerSelectorCryptoCosmoEpay `json:"reportModelsContent"`
	ReportMonitor                                Monitor                             `json:"reportMonitor"`
	ReportingForGuestsModelsStudios              AccountOwnerSelectorCryptoCosmoEpay `json:"reportingForGuestsModelsStudios"`
	ReportsInTextMonitor                         ReportsInTextMonitor                `json:"reportsInTextMonitor"`
	RewardLotteryAsync                           AccountOwnerSelectorCryptoCosmoEpay `json:"rewardLotteryAsync"`
	RulesModal                                   RulesModal                          `json:"rulesModal"`
	SafariCanvas                                 SafariCanvas                        `json:"safariCanvas"`
	SaveBackButtonBuyTokensAnalyticsToClickhouse AccountOwnerSelectorCryptoCosmoEpay `json:"saveBackButtonBuyTokensAnalyticsToClickhouse"`
	SayHi5Tk                                     SayHi5Tk                            `json:"sayHi5Tk"`
	ScWlLoginTooltipUpdate                       AccountOwnerSelectorCryptoCosmoEpay `json:"scWlLoginTooltipUpdate"`
	ScoreWithIP                                  AccountOwnerSelectorCryptoCosmoEpay `json:"scoreWithIp"`
	SecondFAXHL                                  AccountOwnerSelectorCryptoCosmoEpay `json:"secondFAXHL"`
	SecureImages                                 SecureImages                        `json:"secureImages"`
	SendChatInputEvent                           AccountOwnerSelectorCryptoCosmoEpay `json:"sendChatInputEvent"`
	SendModelPublisherReport                     AccountOwnerSelectorCryptoCosmoEpay `json:"sendModelPublisherReport"`
	SendTipInSpyMode                             AccountOwnerSelectorCryptoCosmoEpay `json:"sendTipInSpyMode"`
	SendUserDocumentToFraudSystem                AccountOwnerSelectorCryptoCosmoEpay `json:"sendUserDocumentToFraudSystem"`
	SeparateRedisForMessageValidator             SeparateRedisForMessageValidator    `json:"separateRedisForMessageValidator"`
	SeparateUpdateLastActivityAtWorker           SeparateUpdateLastActivityAtWorker  `json:"separateUpdateLastActivityAtWorker"`
	ServiceWorker                                ServiceWorker                       `json:"serviceWorker"`
	ServiceWorkerCentrifugo                      ServiceWorkerCentrifugo             `json:"serviceWorkerCentrifugo"`
	ShortMessageSpamCheckerExcludeRule           ShortMessageSpamCheckerExcludeRule  `json:"shortMessageSpamCheckerExcludeRule"`
	ShowAllAccessKeysUpdater                     AccountOwnerSelectorCryptoCosmoEpay `json:"showAllAccessKeysUpdater"`
	ShowDataCollector                            ShowDataCollector                   `json:"showDataCollector"`
	ShowFaphouseLink                             AccountOwnerSelectorCryptoCosmoEpay `json:"showFaphouseLink"`
	ShowHistory                                  AccountOwnerSelectorCryptoCosmoEpay `json:"showHistory"`
	ShowP2PUserPings                             ShowP2PUserPings                    `json:"showP2PUserPings"`
	ShowStreamBlockedBanner                      AccountOwnerSelectorCryptoCosmoEpay `json:"showStreamBlockedBanner"`
	SidebarAmplitude                             AccountOwnerSelectorCryptoCosmoEpay `json:"sidebarAmplitude"`
	SignUpCountrySpecific                        AccountOwnerSelectorCryptoCosmoEpay `json:"signUpCountrySpecific"`
	SignUpWithoutEmail                           AccountOwnerSelectorCryptoCosmoEpay `json:"signUpWithoutEmail"`
	SiteMirror                                   FingerprintCDN                      `json:"siteMirror"`
	SiteSettings                                 SiteSettings                        `json:"siteSettings"`
	SkipDeleteOnlineForBroadcast                 AccountOwnerSelectorCryptoCosmoEpay `json:"skipDeleteOnlineForBroadcast"`
	SkipSsr                                      SkipSsr                             `json:"skipSsr"`
	SlidableSegmentConfig                        SlidableSegmentConfig               `json:"slidableSegmentConfig"`
	SmoothAnimationThumb                         AccountOwnerSelectorCryptoCosmoEpay `json:"smoothAnimationThumb"`
	SnapshotsNewGateway                          SnapshotsNewGateway                 `json:"snapshotsNewGateway"`
	SnapshotsRetryInterval                       SnapshotsRetryInterval              `json:"snapshotsRetryInterval"`
	SpecialCategories                            SpecialCategories                   `json:"specialCategories"`
	SsrCaching                                   SsrCaching                          `json:"ssrCaching"`
	SsrInvalidCountryLogger                      AccountOwnerSelectorCryptoCosmoEpay `json:"ssrInvalidCountryLogger"`
	StreamIssuesModels                           AccountOwnerSelectorCryptoCosmoEpay `json:"streamIssuesModels"`
	StreamPixelization                           StreamPixelization                  `json:"streamPixelization"`
	StreamsTranscoding                           StreamsTranscoding                  `json:"streamsTranscoding"`
	StripScoreFactorWebRTC                       AccountOwnerSelectorCryptoCosmoEpay `json:"stripScoreFactorWebRTC"`
	StripsConfig                                 StripsConfig                        `json:"stripsConfig"`
	StudioAPIDocumentation                       AccountOwnerSelectorCryptoCosmoEpay `json:"studioApiDocumentation"`
	SuspiciousLoginCheckWithoutCity              AccountOwnerSelectorCryptoCosmoEpay `json:"suspiciousLoginCheckWithoutCity"`
	SwipeBackMessenger                           AccountOwnerSelectorCryptoCosmoEpay `json:"swipeBackMessenger"`
	SwitchToWebRTCForCam2Cam                     AccountOwnerSelectorCryptoCosmoEpay `json:"switchToWebRTCForCam2Cam"`
	TeasingVideoOnAvpShutter                     AccountOwnerSelectorCryptoCosmoEpay `json:"teasingVideoOnAvpShutter"`
	TemporaryPreloadMixedLivetagsByModelsAPI     AccountOwnerSelectorCryptoCosmoEpay `json:"temporaryPreloadMixedLivetagsByModelsApi"`
	TestConsumers                                TestConsumers                       `json:"testConsumers"`
	TestimonialBlackList                         AccountOwnerSelectorCryptoCosmoEpay `json:"testimonialBlackList"`
	TestimonialReports                           AccountOwnerSelectorCryptoCosmoEpay `json:"testimonialReports"`
	TextMonitor                                  TextMonitor                         `json:"textMonitor"`
	TextMonitorDoubleModeration                  TextMonitorDoubleModeration         `json:"textMonitorDoubleModeration"`
	TextMonitorMessagesPeriod                    TextMonitorMessagesPeriod           `json:"textMonitorMessagesPeriod"`
	ThornCsamCheckerLinks                        AccountOwnerSelectorCryptoCosmoEpay `json:"thornCsamCheckerLinks"`
	TicketShowEnterRefactor                      AccountOwnerSelectorCryptoCosmoEpay `json:"ticketShowEnterRefactor"`
	TipMenuImproveMobile                         TipMenuImproveMobile                `json:"tipMenuImproveMobile"`
	TopModelsAPIV2                               AccountOwnerSelectorCryptoCosmoEpay `json:"topModelsApiV2"`
	ToyMenuReskin                                InvoiceV2                           `json:"toyMenuReskin"`
	TranscodeVideo                               TranscodeVideo                      `json:"transcodeVideo"`
	TranslateChatsForModels                      AccountOwnerSelectorCryptoCosmoEpay `json:"translateChatsForModels"`
	TranslateMessagesNewLabel                    TranslateMessagesNewLabel           `json:"translateMessagesNewLabel"`
	TranslationsCaching                          TranslationsCaching                 `json:"translationsCaching"`
	TwitterAutopostSelectImage                   AccountOwnerSelectorCryptoCosmoEpay `json:"twitterAutopostSelectImage"`
	TwitterAutopostTimeLimit                     TwitterAutopostTimeLimit            `json:"twitterAutopostTimeLimit"`
	UkrainianInSpecials                          AccountOwnerSelectorCryptoCosmoEpay `json:"ukrainianInSpecials"`
	UnbanIPFingerprint                           AccountOwnerSelectorCryptoCosmoEpay `json:"unbanIpFingerprint"`
	UnicodeEmojies                               AccountOwnerSelectorCryptoCosmoEpay `json:"unicodeEmojies"`
	UnseenForm                                   AccountOwnerSelectorCryptoCosmoEpay `json:"unseenForm"`
	UseAbTestsConfig                             UseAbTestsConfig                    `json:"useAbTestsConfig"`
	UseMarkAfterPrivateShow                      AccountOwnerSelectorCryptoCosmoEpay `json:"useMarkAfterPrivateShow"`
	UserDeviceReport                             AccountOwnerSelectorCryptoCosmoEpay `json:"userDeviceReport"`
	UserStreamNotLoadingLogger                   RefreshModelStatus                  `json:"userStreamNotLoadingLogger"`
	VideoRecordingMigrationToDVR                 VideoRecordingMigrationToDVR        `json:"videoRecordingMigrationToDvr"`
	VideoReportSettings                          VideoReportSettings                 `json:"videoReportSettings"`
	ViewCamPerformanceMetrics                    ViewCamPerformanceMetrics           `json:"viewCamPerformanceMetrics"`
	ViewCamPrivateMessagesMultiline              AccountOwnerSelectorCryptoCosmoEpay `json:"viewCamPrivateMessagesMultiline"`
	ViewCamResolutionLimit                       ViewCamResolutionLimit              `json:"viewCamResolutionLimit"`
	Vr3DWebRTC                                   Vr3DWebRTC                          `json:"vr3DWebRTC"`
	VRABR                                        StreamsTranscoding                  `json:"vrABR"`
	VRAutoSwitch                                 Vr3DWebRTC                          `json:"vrAutoSwitch"`
	VRChatTip                                    Vr3DWebRTC                          `json:"vrChatTip"`
	VRFeeInEarningsByModels                      AccountOwnerSelectorCryptoCosmoEpay `json:"vrFeeInEarningsByModels"`
	VRHLSABR                                     Vr3DWebRTC                          `json:"vrHLSABR"`
	VRNewBuyTokens                               Vr3DWebRTC                          `json:"vrNewBuyTokens"`
	VRQuestModal                                 Vr3DWebRTC                          `json:"vrQuestModal"`
	VRRightPanel                                 Vr3DWebRTC                          `json:"vrRightPanel"`
	VRSocialNetworkLogin                         AccountOwnerSelectorCryptoCosmoEpay `json:"vrSocialNetworkLogin"`
	VRVideos                                     Vr3DWebRTC                          `json:"vrVideos"`
	WebRTCOriginTurnServers                      WebRTCTurnServers                   `json:"webRTCOriginTurnServers"`
	WebRTCTurnServers                            WebRTCTurnServers                   `json:"webRTCTurnServers"`
	WebrtcNewABR                                 map[string]float64                  `json:"webrtcNewAbr"`
	WebrtcNewABRMMP                              map[string]float64                  `json:"webrtcNewAbrMMP"`
	WebrtcPlayer                                 WebrtcPlayer                        `json:"webrtcPlayer"`
	WebrtcStatsCollection                        WebrtcStatsCollection               `json:"webrtcStatsCollection"`
	WhyYouUnsubscribingFromFanClubUltimate       AccountOwnerSelectorCryptoCosmoEpay `json:"whyYouUnsubscribingFromFanClubUltimate"`
	WizardDocsReviewModal                        WizardDocsReviewModal               `json:"wizardDocsReviewModal"`
	XConverter                                   SeparateUpdateLastActivityAtWorker  `json:"xConverter"`
	YouMightLikeSection                          YouMightLikeSection                 `json:"youMightLikeSection"`
	ZendeskExternalID                            AccountOwnerSelectorCryptoCosmoEpay `json:"zendeskExternalId"`
}

type Ab struct {
	IsAbTest          bool         `json:"_isAbTest"`
	AbTestCountries   []string     `json:"abTestCountries"`
	IsAbTestEnabled   bool         `json:"isAbTestEnabled"`
	IsFeatureEnabled  bool         `json:"isFeatureEnabled"`
	Leagues           []any        `json:"leagues"`
	Name              string       `json:"name"`
	PlayerTypes       []PlayerType `json:"playerTypes"`
	Project           Project      `json:"project"`
	Regions           []string     `json:"regions"`
	RegisteredAfter   int64        `json:"registeredAfter"`
	RegisteredBefore  int64        `json:"registeredBefore"`
	SelectedUserTypes []string     `json:"selectedUserTypes"`
	UserPartFrom      int64        `json:"userPartFrom"`
	UserPartTo        int64        `json:"userPartTo"`
	Variants          []Variant    `json:"variants"`
}

type AccessToStreamingApp struct {
	IsEnabled      bool    `json:"isEnabled"`
	ModelsID       []int64 `json:"modelsId"`
	Percent        int64   `json:"percent"`
	PercentEnabled int64   `json:"percentEnabled"`
	StickyBanner   bool    `json:"stickyBanner"`
}

type AccountOwnerSelectorCryptoCosmoEpay struct {
	IsEnabled bool `json:"isEnabled"`
}

type ActivityUpdater struct {
	IsEnabled          bool  `json:"isEnabled"`
	UpdateTimeout      int64 `json:"updateTimeout"`
	UpdateWebXRTimeout int64 `json:"updateWebXRTimeout"`
}

type AlphabetCategories struct {
	ExceptLocales []string `json:"exceptLocales"`
	IsEnabled     bool     `json:"isEnabled"`
}

type AmplitudeConfig struct {
	Events    []string `json:"events"`
	IsEnabled bool     `json:"isEnabled"`
}

type AntispamRules struct {
	CheckPeriods       []int64 `json:"checkPeriods"`
	DefaultCheckPeriod int64   `json:"defaultCheckPeriod"`
	DefaultPoints      int64   `json:"defaultPoints"`
}

type Apps struct {
	AvailableAppNames    []string `json:"availableAppNames"`
	UseAvailableAppNames bool     `json:"useAvailableAppNames"`
}

type AppsPerformanceThresholds struct {
	Battleships int64 `json:"battleships"`
	Reactions   int64 `json:"reactions"`
}

type AutoDiscountPrivates struct {
	DiscountDurationDays       int64 `json:"discountDurationDays"`
	DiscountPercent            int64 `json:"discountPercent"`
	EnabledForAllModels        bool  `json:"enabledForAllModels"`
	IsEnabled                  bool  `json:"isEnabled"`
	IsOfferProposalEnabled     bool  `json:"isOfferProposalEnabled"`
	TimeToActivateOfferMinutes int64 `json:"timeToActivateOfferMinutes"`
}

type BannerMyWebcamRoom struct {
	Countries []string `json:"countries"`
	IsEnabled bool     `json:"isEnabled"`
}

type BlurredLivestreamOnShutter struct {
	IsEnabled                bool  `json:"isEnabled"`
	MinFavoritesNumberToBlur int64 `json:"minFavoritesNumberToBlur"`
}

type BroadcastDefaultQuality struct {
	IsEnabled           bool     `json:"isEnabled"`
	ResolutionsPriority []string `json:"resolutionsPriority"`
	TargetResolution    string   `json:"targetResolution"`
}

type BroadcastQuality struct {
	AddProbingCacheKey      bool  `json:"addProbingCacheKey"`
	CollectInterval         int64 `json:"collectInterval"`
	IsEnabled               bool  `json:"isEnabled"`
	SendInterval            int64 `json:"sendInterval"`
	SetICEServersForFirefox bool  `json:"setICEServersForFirefox"`
}

type BroadcastQualityLevels struct {
	AudioLevelThreshold   float64 `json:"audioLevelThreshold"`
	BadThreshold          int64   `json:"badThreshold"`
	FPSBadThreshold       int64   `json:"fpsBadThreshold"`
	FPSPoorThreshold      int64   `json:"fpsPoorThreshold"`
	IsEnabled             bool    `json:"isEnabled"`
	NoConnectionThreshold int64   `json:"noConnectionThreshold"`
	PoorThreshold         int64   `json:"poorThreshold"`
}

type BroadcastService struct {
	BroadcastFinishTTL        int64 `json:"broadcastFinishTtl"`
	IsBroadcastServiceEnabled bool  `json:"isBroadcastServiceEnabled"`
	IsModelServiceEnabled     bool  `json:"isModelServiceEnabled"`
	IsViewServiceEnabled      bool  `json:"isViewServiceEnabled"`
}

type BrowserCache struct {
	IsLiveTagsCacheEnabled bool `json:"isLiveTagsCacheEnabled"`
}

type Centrifugo struct {
	CommonRecoverableChannels    []any    `json:"commonRecoverableChannels"`
	RecoverableChannelsForGuests []any    `json:"recoverableChannelsForGuests"`
	RecoverableChannelsForModels []string `json:"recoverableChannelsForModels"`
	RecoverableChannelsForUsers  []string `json:"recoverableChannelsForUsers"`
	RecoveryPercentageForGuests  int64    `json:"recoveryPercentageForGuests"`
	RecoveryPercentageForModels  int64    `json:"recoveryPercentageForModels"`
	RecoveryPercentageForUsers   int64    `json:"recoveryPercentageForUsers"`
}

type CentrifugoMetrics struct {
	IsCollectionEnabled bool     `json:"isCollectionEnabled"`
	IsWhitelistEnabled  bool     `json:"isWhitelistEnabled"`
	PercentageForGuests int64    `json:"percentageForGuests"`
	PercentageForModels int64    `json:"percentageForModels"`
	PercentageForUsers  int64    `json:"percentageForUsers"`
	Whitelist           []string `json:"whitelist"`
}

type Comet struct {
	CentrifugoInWorkerPercent        int64 `json:"centrifugoInWorkerPercent"`
	GuestInitialReconnectionDelayMS  int64 `json:"guestInitialReconnectionDelayMs"`
	GuestMaxReconnectionAttempts     int64 `json:"guestMaxReconnectionAttempts"`
	ReconnectionClearOldSubscription bool  `json:"reconnectionClearOldSubscription"`
	UserInitialReconnectionDelayMS   int64 `json:"userInitialReconnectionDelayMs"`
	UserMaxReconnectionAttempts      int64 `json:"userMaxReconnectionAttempts"`
}

type CommunityReportCreatorV2 struct {
	ClickhouseRequestDelayMS int64 `json:"clickhouseRequestDelayMs"`
	IsEnabled                bool  `json:"isEnabled"`
	IsProcessingEnabled      bool  `json:"isProcessingEnabled"`
}

type ContentMonitor struct {
	ActionHistory        bool     `json:"actionHistory"`
	DisabledContentTypes []string `json:"disabledContentTypes"`
}

type ContentMonitorQueueSizeCache struct {
	EnabledForUsersOnly []any `json:"enabledForUsersOnly"`
	IsEnabled           bool  `json:"isEnabled"`
}

type ContentsReports struct {
	AmnestyInDays       int64 `json:"amnestyInDays"`
	CopyrightsThreshold int64 `json:"copyrightsThreshold"`
	FluidsThreshold     int64 `json:"fluidsThreshold"`
	IllegalThreshold    int64 `json:"illegalThreshold"`
	OtherThreshold      int64 `json:"otherThreshold"`
	SpeechThreshold     int64 `json:"speechThreshold"`
	UnderageThreshold   int64 `json:"underageThreshold"`
	ViolenceThreshold   int64 `json:"violenceThreshold"`
}

type CurrentReleaseVersion struct {
	CurrentVersion  string `json:"currentVersion"`
	IsEnabled       bool   `json:"isEnabled"`
	MaxGroupsNumber int64  `json:"maxGroupsNumber"`
}

type CvReportHideEmptyStreams struct {
	HideTTL   int64 `json:"hideTTL"`
	IsEnabled bool  `json:"isEnabled"`
}

type CvReportModalProbability struct {
	Probability int64 `json:"probability"`
}

type CvSearchConfig struct {
	AnimationTimeoutMS    int64    `json:"animationTimeoutMS"`
	Countries             []any    `json:"countries"`
	SuggestionKeysCouples []string `json:"suggestionKeysCouples"`
	SuggestionKeysGirls   []string `json:"suggestionKeysGirls"`
	SuggestionKeysMen     []string `json:"suggestionKeysMen"`
	SuggestionKeysTrans   []string `json:"suggestionKeysTrans"`
	SuggestionsCouples    []string `json:"suggestionsCouples"`
	SuggestionsGirls      []string `json:"suggestionsGirls"`
	SuggestionsMen        []string `json:"suggestionsMen"`
	SuggestionsTrans      []string `json:"suggestionsTrans"`
	WhiteLabelIDS         []any    `json:"whiteLabelIds"`
}

type DSAReportUser struct {
	HoursBetweenReportsOnUser int64 `json:"hoursBetweenReportsOnUser"`
	IsEnabled                 bool  `json:"isEnabled"`
	ReportsOnUsersInOneDay    int64 `json:"reportsOnUsersInOneDay"`
}

type DebugMeasurement struct {
	Events          []any `json:"events"`
	IsEnabled       bool  `json:"isEnabled"`
	MaxTimeDuration int64 `json:"maxTimeDuration"`
}

type DeprecatedPages struct {
	DeprecatedPagesV0 string `json:"deprecatedPagesV0"`
	DeprecatedPagesV1 string `json:"deprecatedPagesV1"`
}

type DetectMediaAPITampering struct {
	BreakEnabled  bool `json:"breakEnabled"`
	ReportEnabled bool `json:"reportEnabled"`
}

type DownscaleBroadcastResolution struct {
	IsEnabled               bool  `json:"isEnabled"`
	ResolutionChangeTimeout int64 `json:"resolutionChangeTimeout"`
}

type DwhAdapter struct {
	Beacon            []string `json:"beacon"`
	Countries         []any    `json:"countries"`
	EnabledEvents     []string `json:"enabledEvents"`
	IsEnabled         bool     `json:"isEnabled"`
	IsFallbackEnabled bool     `json:"isFallbackEnabled"`
	ThrottlePercent   int64    `json:"throttlePercent"`
}

type DwhGuestFeatures struct {
	DisabledFeatures []string `json:"disabledFeatures"`
	EnabledFeatures  []string `json:"enabledFeatures"`
}

type DwhTabKeepAlive struct {
	IntervalSeconds int64 `json:"intervalSeconds"`
	MaxTimes        int64 `json:"maxTimes"`
}

type EnableWebPushiOS struct {
	IsEnabled bool `json:"isEnabled"`
	IsNew     bool `json:"isNew"`
}

type EuDSAPage struct {
	AmarsNumberSC  string `json:"amarsNumberSC"`
	AmarsNumberXHL string `json:"amarsNumberXHL"`
	DateSC         string `json:"dateSC"`
	DateXHL        string `json:"dateXHL"`
	IsEnabled      bool   `json:"isEnabled"`
}

type FingerprintCDN struct {
	URL string `json:"url"`
}

type FlirtingMode struct {
	Countries []any `json:"countries"`
	IsEnabled bool  `json:"isEnabled"`
	ModelsID  []any `json:"modelsId"`
}

type Monitor struct {
	DisabledReviewTypes []any `json:"disabledReviewTypes"`
}

type FrontLogV2 struct {
	APILogDisableSamplingForModelsAndStudios bool     `json:"apiLogDisableSamplingForModelsAndStudios"`
	APILogSampleRate                         float64  `json:"apiLogSampleRate"`
	APILogSamplingPatterns                   []string `json:"apiLogSamplingPatterns"`
	IsAPILogEnabled                          bool     `json:"isApiLogEnabled"`
	LoggerDebugPatterns                      []string `json:"loggerDebugPatterns"`
	LoggerExcludePatterns                    []any    `json:"loggerExcludePatterns"`
	LoggerMoveToWarnPatterns                 []string `json:"loggerMoveToWarnPatterns"`
	SentrySampleRate                         float64  `json:"sentrySampleRate"`
}

type Glass2GlassLatencyMeasurement struct {
	Glass2GlassMeasurementEnableProbability int64  `json:"glass2GlassMeasurementEnableProbability"`
	IsEnabled                               bool   `json:"isEnabled"`
	TimeResyncInterval                      int64  `json:"timeResyncInterval"`
	TimeSyncURLTemplate                     string `json:"timeSyncURLTemplate"`
}

type GoogleAuthentication struct {
	IsEnabled         bool `json:"isEnabled"`
	IsEnabledForStage bool `json:"isEnabledForStage"`
}

type GuestLandedAndIdentityTimeout struct {
	IsEnabled          bool  `json:"isEnabled"`
	VisiblePageTimeout int64 `json:"visiblePageTimeout"`
}

type HLSABRConfig struct {
	ABRDebug              bool    `json:"abrDebug"`
	BandwidthSafetyFactor float64 `json:"bandwidthSafetyFactor"`
	CacheLoadThreshold    int64   `json:"cacheLoadThreshold"`
	UseDeadTimeLatency    bool    `json:"useDeadTimeLatency"`
}

type HLSPlayer struct {
	ABRDebug                              bool     `json:"abrDebug"`
	BandwidthSafetyFactor                 float64  `json:"bandwidthSafetyFactor"`
	CacheLoadThreshold                    int64    `json:"cacheLoadThreshold"`
	CollectDoppioMetrics                  bool     `json:"collectDoppioMetrics"`
	CollectPlaylistToFragmentRatio        bool     `json:"collectPlaylistToFragmentRatio"`
	CollectResolutionWatchedEvents        bool     `json:"collectResolutionWatchedEvents"`
	CollectResolutionWatchingDuration     bool     `json:"collectResolutionWatchingDuration"`
	DesktopResolutionType                 string   `json:"desktopResolutionType"`
	DoppioMetricsBaseURL                  string   `json:"doppioMetricsBaseUrl"`
	DoppioMetricsCollectionInterval       int64    `json:"doppioMetricsCollectionInterval"`
	FallbackToLLHLS                       bool     `json:"fallbackToLLHls"`
	FragLoadingTimeOut                    int64    `json:"fragLoadingTimeOut"`
	HLSBalancer                           string   `json:"hlsBalancer"`
	IsLoadBalancingEnabled                bool     `json:"isLoadBalancingEnabled"`
	LowLatencyMode                        bool     `json:"lowLatencyMode"`
	LowLatencyPlayerType                  []string `json:"lowLatencyPlayerType"`
	ManifestLoadingMaxRetry               int64    `json:"manifestLoadingMaxRetry"`
	ManifestLoadingRetryDelay             int64    `json:"manifestLoadingRetryDelay"`
	ManifestLoadingTimeOut                int64    `json:"manifestLoadingTimeOut"`
	MinimalResolutionWatchingDuration     int64    `json:"minimalResolutionWatchingDuration"`
	MobileResolutionType                  string   `json:"mobileResolutionType"`
	NumberOfStallsToFallback              int64    `json:"numberOfStallsToFallback"`
	PlayHLSUntilStreamEnd                 bool     `json:"playHLSUntilStreamEnd"`
	ResolutionWatchedEventTimeout         int64    `json:"resolutionWatchedEventTimeout"`
	StallsCollectionPeriod                int64    `json:"stallsCollectionPeriod"`
	StartDesktopResolution                string   `json:"startDesktopResolution"`
	StartMobileResolution                 string   `json:"startMobileResolution"`
	TimeoutCheckStreamURL                 int64    `json:"timeoutCheckStreamUrl"`
	TimeoutInitCheckStreamURL             int64    `json:"timeoutInitCheckStreamUrl"`
	TimeoutWaitingStream                  int64    `json:"timeoutWaitingStream"`
	Use404ErrorToDetectStreamEnding       bool     `json:"use404ErrorToDetectStreamEnding"`
	UseDeadTimeLatency                    bool     `json:"useDeadTimeLatency"`
	UseDoppioPlayerModuleInVODProbability int64    `json:"useDoppioPlayerModuleInVODProbability"`
	UseDoppioPlayerModuleProbability      int64    `json:"useDoppioPlayerModuleProbability"`
	UseFMP4SegmentsProbability            int64    `json:"useFMP4SegmentsProbability"`
	UseGlobalPlaylists                    bool     `json:"useGlobalPlaylists"`
	UseLowLatencyProbability              int64    `json:"useLowLatencyProbability"`
	UseNewHLSJSModule                     bool     `json:"useNewHlsJSModule"`
}

type HalloweenReaction struct {
	EndDate   string `json:"endDate"`
	IsEnabled bool   `json:"isEnabled"`
	StartDate string `json:"startDate"`
}

type IndexCategoryLoadTime struct {
	IsEnabled  bool  `json:"isEnabled"`
	SampleRate int64 `json:"sampleRate"`
}

type InpReporter struct {
	DurationThreshold int64 `json:"durationThreshold"`
	IsEnabled         bool  `json:"isEnabled"`
	SampleRate        int64 `json:"sampleRate"`
}

type InvoiceV2 struct {
	IsEnabled bool   `json:"isEnabled"`
	StartDate string `json:"startDate"`
}

type Kyc struct {
	AutoApproveCountries      []string `json:"autoApproveCountries"`
	CvCheckTimeoutMinutes     int64    `json:"cvCheckTimeoutMinutes"`
	IsAutoApproveForAll       bool     `json:"isAutoApproveForAll"`
	IsCVDuplicatesFlowEnabled bool     `json:"isCVDuplicatesFlowEnabled"`
	IsEnabled                 bool     `json:"isEnabled"`
	IsEnabledForStudios       bool     `json:"isEnabledForStudios"`
}

type LazyLoadModelThumbs struct {
	IsEnabled bool  `json:"isEnabled"`
	SkipCount int64 `json:"skipCount"`
}

type LazyLoadingOnIndex struct {
	InitialBlockLoadingCount int64 `json:"initialBlockLoadingCount"`
	IsEnabled                bool  `json:"isEnabled"`
	LazyLoadingBlockCount    int64 `json:"lazyLoadingBlockCount"`
}

type LimitBroadcastHours struct {
	IsEnabled                 bool `json:"isEnabled"`
	IsNewStreamStopAPIEnabled bool `json:"isNewStreamStopApiEnabled"`
}

type LiveThumbs struct {
	Duration                int64 `json:"duration"`
	RequestStreamTimeout    int64 `json:"requestStreamTimeout"`
	SchedulerUnblockTimeout int64 `json:"schedulerUnblockTimeout"`
	VisibilityTimeout       int64 `json:"visibilityTimeout"`
	WatchEventTimeout       int64 `json:"watchEventTimeout"`
}

type LovenseAPIUrls struct {
	AppStoreLink   string `json:"appStoreLink"`
	GooglePlayLink string `json:"googlePlayLink"`
	ModelAPI       string `json:"modelApi"`
	SettingsAPI    string `json:"settingsApi"`
	UserAPI        string `json:"userApi"`
}

type MainStudios struct {
	IsEnabledForAllStudios bool `json:"isEnabledForAllStudios"`
}

type MaxAllowedUserDevices struct {
	Count int64 `json:"count"`
}

type MoEngageSupport struct {
	EnabledOnStands    []string `json:"enabledOnStands"`
	IsEnabled          bool     `json:"isEnabled"`
	IsEnabledForGuests bool     `json:"isEnabledForGuests"`
	IsEventsEnabled    bool     `json:"isEventsEnabled"`
}

type ModelListUltimateBanner struct {
	DisplayAfterRows     int64 `json:"displayAfterRows"`
	KeepClosedForSeconds int64 `json:"keepClosedForSeconds"`
}

type ModelRecordPublicShow struct {
	DelayRecordButtonInSEC int64 `json:"delayRecordButtonInSec"`
	IsEnabled              bool  `json:"isEnabled"`
}

type ModelSearch struct {
	IsEnabledForGuests bool  `json:"isEnabledForGuests"`
	Percentage         int64 `json:"percentage"`
}

type ModelsDeduplication struct {
	MultipleCategoriesRequestLimit int64 `json:"multipleCategoriesRequestLimit"`
}

type MonthlyTopModelsLimits struct {
	Couples   int64 `json:"couples"`
	Girls     int64 `json:"girls"`
	IsEnabled bool  `json:"isEnabled"`
	Men       int64 `json:"men"`
	Trans     int64 `json:"trans"`
}

type MouflonMediaPlayer struct {
	HLSPartHoldback           int64  `json:"hlsPartHoldback"`
	IsDetachedProbability     int64  `json:"isDetachedProbability"`
	LiveThumbProbability      int64  `json:"liveThumbProbability"`
	MetricsCollectorEndpoint  string `json:"metricsCollectorEndpoint"`
	MetricsProjectID          string `json:"metricsProjectId"`
	Probability               int64  `json:"probability"`
	ProbabilityBroadcast      int64  `json:"probabilityBroadcast"`
	SendCustomRefererHeader   bool   `json:"sendCustomRefererHeader"`
	TempUseSafariC2CFix       bool   `json:"tempUseSafariC2CFix"`
	UseBitrateBasedResolution bool   `json:"useBitrateBasedResolution"`
}

type MouflonStatsCollection struct {
	CollectStatsIntervalHLS                   int64   `json:"collectStatsIntervalHLS"`
	CollectStatsIntervalWebRTC                int64   `json:"collectStatsIntervalWebRTC"`
	CollectStatsProbabilityHLS                float64 `json:"collectStatsProbabilityHLS"`
	CollectStatsProbabilityWebRTCBroadcasting float64 `json:"collectStatsProbabilityWebRTCBroadcasting"`
	CollectStatsProbabilityWebRTCPlaying      float64 `json:"collectStatsProbabilityWebRTCPlaying"`
	MetricsBaseURL                            string  `json:"metricsBaseUrl"`
	ShouldCollectFirstFrameMetrics            bool    `json:"shouldCollectFirstFrameMetrics"`
}

type NewBdsmCategory struct {
	BlockBdsmInDays               int64 `json:"blockBdsmInDays"`
	DoNotCheckAfterApproveInHours int64 `json:"doNotCheckAfterApproveInHours"`
	IsEnabled                     bool  `json:"isEnabled"`
	NeedToBlock                   bool  `json:"needToBlock"`
}

type NewImageViewer struct {
	IsEnabledInNotes bool `json:"isEnabledInNotes"`
}

type NewWatchTimeLogic struct {
	PublicRecordLogicIsEnabled  bool  `json:"publicRecordLogicIsEnabled"`
	PublicRecordRequestInterval int64 `json:"publicRecordRequestInterval"`
	WatchTimeLogicIsEnabled     bool  `json:"watchTimeLogicIsEnabled"`
}

type NonNude struct {
	AdultThreshold        float64 `json:"adultThreshold"`
	AutoHide              bool    `json:"autoHide"`
	CheckInterval         int64   `json:"checkInterval"`
	HideOnError           bool    `json:"hideOnError"`
	IsEnabled             bool    `json:"isEnabled"`
	ReportsDuration       int64   `json:"reportsDuration"`
	UserReportCoefficient float64 `json:"userReportCoefficient"`
}

type Optimization struct {
	UpdateBroadcastSettings []string `json:"updateBroadcastSettings"`
}

type PerformancePages struct {
	CollectTimeout int64 `json:"collectTimeout"`
	Coverage       int64 `json:"coverage"`
	IsEnabled      bool  `json:"isEnabled"`
	SendInterval   int64 `json:"sendInterval"`
}

type PixelateShowsForJapan struct {
	Country         []string `json:"country"`
	IsEnabled       bool     `json:"isEnabled"`
	IsEnabledOnBeta bool     `json:"isEnabledOnBeta"`
	ModelIDS        []any    `json:"modelIds"`
	OnlyPublicShow  bool     `json:"onlyPublicShow"`
}

type PlayerModuleExternalLoading struct {
	MaxAttemptsToLoad            int64  `json:"maxAttemptsToLoad"`
	MmpProbability               int64  `json:"mmpProbability"`
	MmpVersion                   string `json:"mmpVersion"`
	PreviewMMPVersion            string `json:"previewMMPVersion"`
	PreviewMMPVersionProbability int64  `json:"previewMMPVersionProbability"`
	Probability                  int64  `json:"probability"`
	Version                      string `json:"version"`
}

type PreviewsSectionInMyDetails struct {
	BannerVideoTeaser bool `json:"bannerVideoTeaser"`
	IsEnabled         bool `json:"isEnabled"`
}

type ProhibitSubStudioChanges struct {
	AuthorityUserIDS []int64 `json:"authorityUserIds"`
	IsEnabled        bool    `json:"isEnabled"`
}

type ProhibitedWords struct {
	AccessSetupUserIDS []int64 `json:"accessSetupUserIds"`
	IsEnabled          bool    `json:"isEnabled"`
}

type PulsePromotionBanners struct {
	DiscordInvite  bool `json:"discordInvite"`
	IsEnabled      bool `json:"isEnabled"`
	TelegramLink   bool `json:"telegramLink"`
	YouTubeTrailer bool `json:"youTubeTrailer"`
}

type ReRequestFPv2 struct {
	IsEnabled    bool  `json:"isEnabled"`
	SamplingRate int64 `json:"samplingRate"`
}

type RecordStreams struct {
	Use720PPresetLimit bool `json:"use720pPresetLimit"`
}

type RecordUserSessionAmplitude struct {
	Countries  []string `json:"countries"`
	IsEnabled  bool     `json:"isEnabled"`
	SampleRate float64  `json:"sampleRate"`
}

type RecordVRStreams struct {
	Enabled bool `json:"enabled"`
}

type RefreshModelStatus struct {
	IsEnabled bool  `json:"isEnabled"`
	Timeout   int64 `json:"timeout"`
}

type RegularAppealCountries struct {
	Country     []string `json:"country"`
	IsWorldWide bool     `json:"isWorldWide"`
}

type ReportsInTextMonitor struct {
	SendChatMessageReportsToMonitor bool `json:"sendChatMessageReportsToMonitor"`
}

type RulesModal struct {
	Countries       []string `json:"countries"`
	IsEnabled       bool     `json:"isEnabled"`
	IsEnabledOnProd bool     `json:"isEnabledOnProd"`
}

type SafariCanvas struct {
	HideVideoWhileNaNDimensions bool  `json:"hideVideoWhileNaNDimensions"`
	ProbabilityHLS              int64 `json:"probabilityHLS"`
	ProbabilityWebRTC           int64 `json:"probabilityWebRTC"`
	UseCSSResizerForHLS         bool  `json:"useCSSResizerForHLS"`
	UseCSSResizerForWebRTC      bool  `json:"useCSSResizerForWebRTC"`
	UseResizerInAnyBrowser      bool  `json:"useResizerInAnyBrowser"`
}

type SayHi5Tk struct {
	IsEnabled       bool  `json:"isEnabled"`
	TimeoutDuration int64 `json:"timeoutDuration"`
}

type SecureImages struct {
	AdminPanelOpeningLimitDaily       int64   `json:"adminPanelOpeningLimitDaily"`
	AdminPanelOpeningLimitMonthly     int64   `json:"adminPanelOpeningLimitMonthly"`
	AllowedCreationLagSeconds         int64   `json:"allowedCreationLagSeconds"`
	AntiFraudOpeningLimitDaily        int64   `json:"antiFraudOpeningLimitDaily"`
	AntiFraudOpeningLimitMonthly      int64   `json:"antiFraudOpeningLimitMonthly"`
	EncryptImages                     bool    `json:"encryptImages"`
	ExternalOpeningLimitDaily         int64   `json:"externalOpeningLimitDaily"`
	ExternalOpeningLimitMonthly       int64   `json:"externalOpeningLimitMonthly"`
	IsEnabled                         bool    `json:"isEnabled"`
	IsLDAPUsernameCheckerEnabled      bool    `json:"isLdapUsernameCheckerEnabled"`
	IsLoggingEnabled                  bool    `json:"isLoggingEnabled"`
	MaxTTLSeconds                     int64   `json:"maxTtlSeconds"`
	OpeningTimeDurationSeconds        int64   `json:"openingTimeDurationSeconds"`
	ResizeEnabled                     bool    `json:"resizeEnabled"`
	RightMouseClickProhibit           bool    `json:"rightMouseClickProhibit"`
	RightMouseClickProhibitExclusions []int64 `json:"rightMouseClickProhibitExclusions"`
}

type SeparateRedisForMessageValidator struct {
	IsEnabled         bool `json:"isEnabled"`
	IsFallbackEnabled bool `json:"isFallbackEnabled"`
}

type SeparateUpdateLastActivityAtWorker struct {
	IsEnabled      bool  `json:"isEnabled"`
	PercentEnabled int64 `json:"percentEnabled"`
}

type ServiceWorker struct {
	PingInterval          int64 `json:"pingInterval"`
	PingTimeout           int64 `json:"pingTimeout"`
	RegisterTimeout       int64 `json:"registerTimeout"`
	WorkerFeaturesTimeout int64 `json:"workerFeaturesTimeout"`
}

type ServiceWorkerCentrifugo struct {
	MaxSubscriptions       int64 `json:"maxSubscriptions"`
	MaxSubscriptionsPerTab int64 `json:"maxSubscriptionsPerTab"`
}

type ShortMessageSpamCheckerExcludeRule struct {
	MinMessageLength  int64 `json:"minMessageLength"`
	MinUsernameLength int64 `json:"minUsernameLength"`
}

type ShowDataCollector struct {
	Delay     int64 `json:"delay"`
	IsEnabled bool  `json:"isEnabled"`
}

type ShowP2PUserPings struct {
	IsFrontendEnabled bool `json:"isFrontendEnabled"`
	IsStreamEnabled   bool `json:"isStreamEnabled"`
}

type SiteSettings struct {
	AbTestsPollingInterval int64 `json:"abTestsPollingInterval"`
}

type SkipSsr struct {
	Countries []string `json:"countries"`
}

type SlidableSegmentConfig struct {
	DelayMS        int64 `json:"delayMs"`
	ResetByHypotPx int64 `json:"resetByHypotPx"`
}

type SnapshotsNewGateway struct {
	IsAuthBearer bool `json:"isAuthBearer"`
	IsEnabled    bool `json:"isEnabled"`
}

type SnapshotsRetryInterval struct {
	RetryInterval int64 `json:"retryInterval"`
}

type SpecialCategories struct {
	DisplayTagEndDate   string   `json:"displayTagEndDate"`
	DisplayTagStartDate string   `json:"displayTagStartDate"`
	IsEnabled           bool     `json:"isEnabled"`
	OnboardingEndDate   string   `json:"onboardingEndDate"`
	OnboardingStartDate string   `json:"onboardingStartDate"`
	Tags                []string `json:"tags"`
}

type SsrCaching struct {
	CachingTimeInSecongs          int64 `json:"cachingTimeInSecongs"`
	DefaultCacheTime              int64 `json:"defaultCacheTime"`
	IndexCacheTime                int64 `json:"indexCacheTime"`
	IsEnabled                     bool  `json:"isEnabled"`
	IsViewcamMetadataCacheEnabled bool  `json:"isViewcamMetadataCacheEnabled"`
	MetadataTTLSeconds            int64 `json:"metadataTtlSeconds"`
	ViewCamCacheTime              int64 `json:"viewCamCacheTime"`
}

type StreamPixelization struct {
	FirstFeaturedGirls int64    `json:"firstFeaturedGirls"`
	IsEnabledOnBeta    bool     `json:"isEnabledOnBeta"`
	ModelCountries     []string `json:"modelCountries"`
	ModelIDS           []any    `json:"modelIds"`
}

type StreamsTranscoding struct {
	IsEnabled bool  `json:"isEnabled"`
	Percent   int64 `json:"percent"`
}

type StripsConfig struct {
	AddLinksForGuests              bool  `json:"addLinksForGuests"`
	MinUpscaleFactor               int64 `json:"minUpscaleFactor"`
	PlayerStartTimeout             int64 `json:"playerStartTimeout"`
	ResetInactivityDurationSeconds int64 `json:"resetInactivityDurationSeconds"`
	ResolutionLimitEnabled         bool  `json:"resolutionLimitEnabled"`
	UseQualitySwitcher             bool  `json:"useQualitySwitcher"`
}

type TestConsumers struct {
	ConsumeDelaySeconds    int64 `json:"consumeDelaySeconds"`
	ExceptionFactorPercent int64 `json:"exceptionFactorPercent"`
}

type TextMonitor struct {
	ChangedReportDisablePeriod                   int64 `json:"changedReportDisablePeriod"`
	DeletedUsersSpamAlertsTestimonialsDisplaying bool  `json:"deletedUsersSpamAlertsTestimonialsDisplaying"`
	DisabledReviewTypes                          []any `json:"disabledReviewTypes"`
}

type TextMonitorDoubleModeration struct {
	BanByModel      string `json:"banByModel"`
	ReportedMessage string `json:"reportedMessage"`
	ReportedUser    string `json:"reportedUser"`
	SpamAlert       string `json:"spamAlert"`
}

type TextMonitorMessagesPeriod struct {
	PeriodInDays int64 `json:"periodInDays"`
}

type TipMenuImproveMobile struct {
	IsEnabled       bool  `json:"isEnabled"`
	IsEnabledOnProd bool  `json:"isEnabledOnProd"`
	OverlapAmount   int64 `json:"overlapAmount"`
}

type TranscodeVideo struct {
	TranscodingEnabledUntil string `json:"transcodingEnabledUntil"`
}

type TranslateMessagesNewLabel struct {
	IsEnabled                  bool  `json:"isEnabled"`
	IsEnabledForBeta           bool  `json:"isEnabledForBeta"`
	MaxSymbolsForDetect        int64 `json:"maxSymbolsForDetect"`
	MinSymbolsForDetect        int64 `json:"minSymbolsForDetect"`
	UserTranslationDetectDelay int64 `json:"userTranslationDetectDelay"`
}

type TranslationsCaching struct {
	IsEnabled             bool  `json:"isEnabled"`
	MaxNumberOfCharacters int64 `json:"maxNumberOfCharacters"`
}

type TwitterAutopostTimeLimit struct {
	LimitHours int64 `json:"limitHours"`
}

type UseAbTestsConfig struct {
	IsEnabled                 bool  `json:"isEnabled"`
	PercentOfUsageOfNewConfig int64 `json:"percentOfUsageOfNewConfig"`
}

type Vr3DWebRTC struct {
	FeaturePercentage int64 `json:"featurePercentage"`
}

type VideoRecordingMigrationToDVR struct {
	AllowedModelsIDS []int64 `json:"allowedModelsIds"`
	IsEnabled        bool    `json:"isEnabled"`
	Percent          int64   `json:"percent"`
}

type VideoReportSettings struct {
	ComplaintDuplicationTimeout int64 `json:"complaintDuplicationTimeout"`
}

type ViewCamPerformanceMetrics struct {
	IsEnabled              bool  `json:"isEnabled"`
	StartPlayingSampleRate int64 `json:"startPlayingSampleRate"`
}

type ViewCamResolutionLimit struct {
	MinUpscaleFactor       float64 `json:"minUpscaleFactor"`
	MinUpscaleFactorGroupB float64 `json:"minUpscaleFactorGroupB"`
	MinUpscaleFactorGroupC float64 `json:"minUpscaleFactorGroupC"`
}

type WebrtcPlayer struct {
	AvailableABRTypes      []string `json:"availableABRTypes"`
	DesktopResolutionType  string   `json:"desktopResolutionType"`
	FastStart              bool     `json:"fastStart"`
	FECOverheadAmount      int64    `json:"fecOverheadAmount"`
	IsSafariCanvasEnabled  bool     `json:"isSafariCanvasEnabled"`
	IsStereoAudioEnabled   bool     `json:"isStereoAudioEnabled"`
	MobileResolutionType   string   `json:"mobileResolutionType"`
	NewABRProbability      int64    `json:"newABRProbability"`
	StartDesktopResolution string   `json:"startDesktopResolution"`
	StartMobileResolution  string   `json:"startMobileResolution"`
}

type WebrtcStatsCollection struct {
	CollectStatsInterval    int64  `json:"collectStatsInterval"`
	CollectStatsProbability int64  `json:"collectStatsProbability"`
	MetricsBaseURL          string `json:"metricsBaseUrl"`
}

type WizardDocsReviewModal struct {
	CustomCountryReviewTimeMap string `json:"customCountryReviewTimeMap"`
}

type YouMightLikeSection struct {
	IsEnabled  bool  `json:"isEnabled"`
	SectionIDS []any `json:"sectionIds"`
}

type FingerprintV2 struct {
	APIKey            string `json:"apiKey"`
	CacheLifetimeDays int64  `json:"cacheLifetimeDays"`
	Domain            string `json:"domain"`
}

type GuestLimits struct {
	Favorites           int64 `json:"favorites"`
	OnlineSubscriptions int64 `json:"onlineSubscriptions"`
	WatchTimeLimit      int64 `json:"watchTimeLimit"`
}

type Hosts struct {
	Stripchat    string `json:"stripchat"`
	Xhamsterlive string `json:"xhamsterlive"`
}

type Language struct {
	Code   string `json:"code"`
	Native string `json:"native"`
	Title  string `json:"title"`
}

type Links struct {
	BasicTourForModels                    string `json:"basicTourForModels"`
	CameraLightAction                     string `json:"cameraLightAction"`
	ContentReporting                      string `json:"contentReporting"`
	ContestForModels                      string `json:"contestForModels"`
	CookiePolicy                          string `json:"cookiePolicy"`
	DMCAProtection                        string `json:"dmcaProtection"`
	DSACompliance                         string `json:"dsaCompliance"`
	FanClubs                              string `json:"fanClubs"`
	FAQ                                   string `json:"faq"`
	FAQModel                              string `json:"faqModel"`
	Halloween                             string `json:"halloween"`
	HowToConnectStudio                    string `json:"howToConnectStudio"`
	HowToSetUpMyPaymentSettingsIndividual string `json:"howToSetUpMyPaymentSettingsIndividual"`
	HowToSetUpMyPaymentSettingsStudio     string `json:"howToSetUpMyPaymentSettingsStudio"`
	HowToUseOBS                           string `json:"howToUseOBS"`
	HowToUseVR                            string `json:"howToUseVR"`
	July4Th                               string `json:"july4th"`
	LearnMoreGroupShow                    string `json:"learnMoreGroupShow"`
	Levels                                string `json:"levels"`
	Oktoberfest                           string `json:"oktoberfest"`
	ReportingAndComplaint                 string `json:"reportingAndComplaint"`
	RulesForModels                        string `json:"rulesForModels"`
	RulesForStudios                       string `json:"rulesForStudios"`
	RulesPenalties                        string `json:"rulesPenalties"`
	StreamSpecifics                       string `json:"streamSpecifics"`
	StripScore                            string `json:"stripScore"`
	StudioRules                           string `json:"studioRules"`
	StudiosFAQ                            string `json:"studiosFAQ"`
	SupportVR                             string `json:"supportVr"`
	Tour                                  string `json:"tour"`
	TourModel                             string `json:"tourModel"`
	Valentines                            string `json:"valentines"`
	WaysOfMakingMoney                     string `json:"waysOfMakingMoney"`
}

type MainPersonConfig struct {
	DisplayBannerToNonMainPersonBeneficiary bool `json:"displayBannerToNonMainPersonBeneficiary"`
	IsEnabledForIndividualModel             bool `json:"isEnabledForIndividualModel"`
	IsEnabledForStudioModel                 bool `json:"isEnabledForStudioModel"`
	LimitBeneficiaryNameByMainPerson        bool `json:"limitBeneficiaryNameByMainPerson"`
}

type ModelTop struct {
	Continents      []string `json:"continents"`
	ViewerContinent string   `json:"viewerContinent"`
}

type Moengage struct {
	AppID     string `json:"appId"`
	DebugLogs int64  `json:"debugLogs"`
}

type NotInterested struct {
	Limit int64 `json:"limit"`
}

type PlatformPlayerMap struct {
	Flash  Flash                   `json:"flash"`
	Webrtc PlatformPlayerMapWebrtc `json:"webrtc"`
}

type Flash struct {
	Android WindowsClass `json:"android"`
	Ios     FlashIos     `json:"ios"`
	Linux   WindowsClass `json:"linux"`
	Macos   WindowsClass `json:"macos"`
	Other   WindowsClass `json:"other"`
	Windows WindowsClass `json:"windows"`
}

type WindowsClass struct {
	Any []PlatformName `json:"any"`
}

type FlashIos struct {
	Any          []PlatformName `json:"any"`
	Safari9Minus []PlatformName `json:"safari9Minus"`
}

type PlatformPlayerMapWebrtc struct {
	Android WebrtcAndroid `json:"android"`
	Ios     WebrtcIos     `json:"ios"`
	Linux   WebrtcAndroid `json:"linux"`
	Macos   WebrtcAndroid `json:"macos"`
	Other   WebrtcAndroid `json:"other"`
	Windows Windows       `json:"windows"`
}

type WebrtcAndroid struct {
	Aloha         []PlayerType `json:"aloha,omitempty"`
	Any           []PlayerType `json:"any"`
	Chrome51Minus []PlayerType `json:"chrome51Minus"`
	SmartTV       []PlayerType `json:"smartTV"`
	Safari11Minus []PlayerType `json:"safari11Minus,omitempty"`
}

type WebrtcIos struct {
	Aloha           []PlayerType   `json:"aloha"`
	Any             []PlayerType   `json:"any"`
	Safari9Minus    []PlatformName `json:"safari9Minus"`
	SnapchatWebView []PlayerType   `json:"snapchatWebView"`
}

type Windows struct {
	Any             []PlayerType `json:"any"`
	Chrome51Minus   []PlayerType `json:"chrome51Minus"`
	Edge15Minus     []PlayerType `json:"edge15Minus"`
	Ie              []PlayerType `json:"ie"`
	OperaOnWindows7 []PlayerType `json:"operaOnWindows7"`
	SmartTV         []PlayerType `json:"smartTV"`
}

type PlatformPlayerMapSpecialRules struct {
	Webrtc []WebrtcElement `json:"webrtc"`
}

type WebrtcElement struct {
	PlayerList []PlayerType `json:"playerList"`
	Rule       Rule         `json:"rule"`
}

type Rule struct {
	SamsungInternetForAndroid *string  `json:"Samsung Internet for Android,omitempty"`
	Linux                     *Linux   `json:"linux,omitempty"`
	Ios                       *RuleIos `json:"ios,omitempty"`
}

type RuleIos struct {
	Chrome string `json:"chrome"`
}

type Linux struct {
	Firefox string `json:"firefox"`
}

type PrivateAnimationTypes struct {
	Bouquet  int64 `json:"bouquet"`
	Diamond  int64 `json:"diamond"`
	Excited  int64 `json:"excited"`
	Hello    int64 `json:"hello"`
	JerkOff  int64 `json:"jerk_off"`
	LetsPlay int64 `json:"lets_play"`
	LoveYou  int64 `json:"love_you"`
	Rainbow  int64 `json:"rainbow"`
}

type PrivateMessages struct {
	MaxLength int64 `json:"maxLength"`
}

type SocialVerificationLink struct {
	BlackListedNames []string `json:"blackListedNames,omitempty"`
	Label            string   `json:"label"`
	Placeholder      string   `json:"placeholder"`
	Regexp           string   `json:"regexp"`
	Type             string   `json:"type"`
}

type SupportLinks struct {
	AutoRefill                      string `json:"autoRefill"`
	BdsmStreamBlockedViolationRules string `json:"bdsmStreamBlockedViolationRules"`
	CommunityGuidelines             string `json:"communityGuidelines"`
	CryptoPaymentsSettings          string `json:"cryptoPaymentsSettings"`
	FlirtingModeRules               string `json:"flirtingModeRules"`
	HelpCenter                      string `json:"helpCenter"`
	HowToBuyBiggerPackages          string `json:"howToBuyBiggerPackages"`
	HowToSetUpPlasmaStreaming       string `json:"howToSetUpPlasmaStreaming"`
	IDUploadGuidelines              string `json:"idUploadGuidelines"`
	JapanSafetySurvey               string `json:"japanSafetySurvey"`
	KybRegistrationDocuments        string `json:"kybRegistrationDocuments"`
	LovenseConnectHelp              string `json:"lovenseConnectHelp"`
	MakeRequest                     string `json:"makeRequest"`
	PaxumPaymentsSettings           string `json:"paxumPaymentsSettings"`
	PayoutFrequency                 string `json:"payoutFrequency"`
	PayoutTiming                    string `json:"payoutTiming"`
	PersonalOffer                   string `json:"personalOffer"`
	PromoPeriodForNewModel          string `json:"promoPeriodForNewModel"`
	ReferralProgram                 string `json:"referralProgram"`
	RulesForModels                  string `json:"rulesForModels"`
	RulesForStudios                 string `json:"rulesForStudios"`
	SendStudioTokens                string `json:"sendStudioTokens"`
	VideoTeaserRules                string `json:"videoTeaserRules"`
	VRFee                           string `json:"vrFee"`
	WhatAreFavorites                string `json:"whatAreFavorites"`
	WhatIsIncomePerHour             string `json:"whatIsIncomePerHour"`
	WhatIsPrivateShowRating         string `json:"whatIsPrivateShowRating"`
	WhatIsStripScore                string `json:"whatIsStripScore"`
}

type Handy struct {
	APIBaseURL     string `json:"apiBaseUrl"`
	AppStoreLink   string `json:"appStoreLink"`
	GooglePlayLink string `json:"googlePlayLink"`
	ToyStoreLink   string `json:"toyStoreLink"`
}

type TweetMyShow struct {
	AccountName    string `json:"accountName"`
	HashTags       string `json:"hashTags"`
	HashTagsTokens string `json:"hashTagsTokens"`
}

type UserLevelsRanking struct {
	LeagueLevels    LeagueLevels     `json:"leagueLevels"`
	LevelExperience map[string]int64 `json:"levelExperience"`
}

type LeagueLevels struct {
	Bronze  []int64 `json:"bronze"`
	Diamond []int64 `json:"diamond"`
	Gold    []int64 `json:"gold"`
	Grey    []int64 `json:"grey"`
	Legend  []int64 `json:"legend"`
	Royal   []int64 `json:"royal"`
	Silver  []int64 `json:"silver"`
}

type ViewersList struct {
	IsEnabled      bool  `json:"isEnabled"`
	UpdateInterval int64 `json:"updateInterval"`
}

type WhiteLabel struct {
	ModelsFilterHash     any   `json:"modelsFilterHash"`
	ModelsFilterPrepared []any `json:"modelsFilterPrepared"`
}

type XhlConfig struct {
	URL URL `json:"url"`
}

type URL struct {
	AccountSettings AccountSettings `json:"accountSettings"`
	CrossDomainAuth []string        `json:"crossDomainAuth"`
	Login           string          `json:"login"`
	Signup          string          `json:"signup"`
	SignupModel     string          `json:"signupModel"`
	SignupStudio    string          `json:"signupStudio"`
	Tracking        string          `json:"tracking"`
	XhlURL          string          `json:"xhlUrl"`
}

type AccountSettings struct {
	ChangeEmail    string `json:"changeEmail"`
	ChangePassword string `json:"changePassword"`
	RecoveryURL    string `json:"recoveryUrl"`
}

type PlayerType string

const (
	HLS                PlayerType = "hls"
	PlayerTypeFlash    PlayerType = "flash"
	PlayerTypeSnapshot PlayerType = "snapshot"
	Webrtc             PlayerType = "webrtc"
)

type Project string

const (
	Main  Project = "MAIN"
	Webxr Project = "WEBXR"
)

type Variant string

const (
	A Variant = "A"
	B Variant = "B"
	C Variant = "C"
	E Variant = "E"
)

type PlatformName string

const (
	Html5                PlatformName = "html5"
	PlatformNameFlash    PlatformName = "flash"
	PlatformNameSnapshot PlatformName = "snapshot"
)
