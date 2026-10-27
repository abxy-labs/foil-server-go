package foil

type ListResult[T any] struct {
	Items      []T
	Limit      int
	HasMore    bool
	NextCursor string
}

type FieldError struct {
	Name     string `json:"name"`
	Issue    string `json:"issue"`
	Expected string `json:"expected,omitempty"`
	Received any    `json:"received,omitempty"`
}

type ErrorDetails struct {
	Fields []FieldError `json:"fields,omitempty"`
}

type ApiErrorBody struct {
	Code      string       `json:"code"`
	Message   string       `json:"message"`
	Status    int          `json:"status"`
	Retryable bool         `json:"retryable"`
	RequestID string       `json:"request_id"`
	DocsURL   string       `json:"docs_url,omitempty"`
	Details   ErrorDetails `json:"details,omitempty"`
}

type apiErrorEnvelope struct {
	Error ApiErrorBody `json:"error"`
}

type meta struct {
	RequestID string `json:"request_id"`
}

type pagination struct {
	Limit      int    `json:"limit"`
	HasMore    bool   `json:"has_more"`
	NextCursor string `json:"next_cursor,omitempty"`
}

type resourceEnvelope[T any] struct {
	Data T    `json:"data"`
	Meta meta `json:"meta"`
}

type resourceListEnvelope[T any] struct {
	Data       []T        `json:"data"`
	Pagination pagination `json:"pagination"`
	Meta       meta       `json:"meta"`
}

type Manipulation struct {
	Score   *int    `json:"score,omitempty"`
	Verdict *string `json:"verdict,omitempty"`
}

type Decision struct {
	EventID              string        `json:"event_id"`
	Verdict              string        `json:"verdict"`
	RiskScore            int           `json:"risk_score"`
	Phase                *string       `json:"phase,omitempty"`
	IsProvisional        *bool         `json:"is_provisional,omitempty"`
	Manipulation         *Manipulation `json:"manipulation,omitempty"`
	EvaluationDurationMS *int          `json:"evaluation_duration_ms,omitempty"`
	EvaluatedAt          string        `json:"evaluated_at"`
}

type SessionDecision struct {
	EventID          string  `json:"event_id"`
	AutomationStatus string  `json:"automation_status"`
	RiskScore        int     `json:"risk_score"`
	EvaluationPhase  *string `json:"evaluation_phase,omitempty"`
	DecisionStatus   string  `json:"decision_status"`
	EvaluatedAt      string  `json:"evaluated_at"`
}

type SessionHighlightEvidence struct {
	Signal string `json:"signal"`
	Name   string `json:"name"`
}

type SessionHighlight struct {
	Key        string                     `json:"key"`
	Effect     string                     `json:"effect"`
	Importance string                     `json:"importance"`
	Summary    string                     `json:"summary"`
	Evidence   []SessionHighlightEvidence `json:"evidence,omitempty"`
}

type SessionAttributionLabel struct {
	Kind       string `json:"kind"`
	Value      string `json:"value"`
	Label      string `json:"label"`
	Confidence int    `json:"confidence"`
}

type SessionAttributionBehavior struct {
	Channel    string `json:"channel"`
	Value      string `json:"value"`
	Label      string `json:"label"`
	Confidence int    `json:"confidence"`
}

type SessionAttribution struct {
	Labels    []SessionAttributionLabel    `json:"labels"`
	Behaviors []SessionAttributionBehavior `json:"behaviors"`
}

type WebBotAuthStatus string

const (
	WebBotAuthStatusVerified WebBotAuthStatus = "verified"
	WebBotAuthStatusInvalid  WebBotAuthStatus = "invalid"
	WebBotAuthStatusError    WebBotAuthStatus = "error"
	WebBotAuthStatusAbsent   WebBotAuthStatus = "absent"
)

type SessionWebBotAuth struct {
	Status *WebBotAuthStatus `json:"status,omitempty"`
	Domain *string           `json:"domain,omitempty"`
}

type SessionRuntimeIntegrity struct {
	Tampering           string `json:"tampering"`
	DeveloperTools      string `json:"developer_tools"`
	Emulation           string `json:"emulation"`
	Virtualization      string `json:"virtualization"`
	PrivacyHardening    string `json:"privacy_hardening"`
	IdentitySpoofing    string `json:"identity_spoofing"`
	Replay              string `json:"replay"`
	OutdatedEnvironment string `json:"outdated_environment"`
}

type VisitorFingerprintLink struct {
	Object       string  `json:"object"`
	ID           string  `json:"id"`
	Confidence   *int    `json:"confidence,omitempty"`
	IdentifiedAt *string `json:"identified_at,omitempty"`
}

type RequestContext struct {
	UserAgent      string  `json:"user_agent"`
	URL            string  `json:"url"`
	ScreenSize     *string `json:"screen_size,omitempty"`
	IsTouchCapable *bool   `json:"is_touch_capable,omitempty"`
	IPAddress      string  `json:"ip_address"`
}

type SessionDetailRequest struct {
	URL       string  `json:"url"`
	Referrer  *string `json:"referrer,omitempty"`
	UserAgent string  `json:"user_agent"`
}

type BrowserEngine string

const (
	BrowserEngineBlink   BrowserEngine = "blink"
	BrowserEngineGecko   BrowserEngine = "gecko"
	BrowserEngineWebKit  BrowserEngine = "webkit"
	BrowserEngineUnknown BrowserEngine = "unknown"
)

type SessionBrowser struct {
	Name         *string       `json:"name,omitempty"`
	Version      *string       `json:"version,omitempty"`
	MajorVersion *string       `json:"major_version,omitempty"`
	Engine       BrowserEngine `json:"engine"`
}

type SessionDeviceOperatingSystem struct {
	Name    *string `json:"name,omitempty"`
	Version *string `json:"version,omitempty"`
}

type SessionDeviceScreen struct {
	Size            *string  `json:"size,omitempty"`
	ColorDepth      *int     `json:"color_depth,omitempty"`
	PixelRatio      *float64 `json:"pixel_ratio,omitempty"`
	OrientationType *string  `json:"orientation_type,omitempty"`
}

type SessionDeviceLocale struct {
	Timezone        *string  `json:"timezone,omitempty"`
	PrimaryLanguage *string  `json:"primary_language,omitempty"`
	Languages       []string `json:"languages"`
}

type SessionDeviceTouchCapabilities struct {
	Available      *bool `json:"available,omitempty"`
	MaxTouchPoints *int  `json:"max_touch_points,omitempty"`
}

type SessionDeviceStorageCapabilities struct {
	Cookies       *bool `json:"cookies,omitempty"`
	LocalStorage  *bool `json:"local_storage,omitempty"`
	IndexedDB     *bool `json:"indexed_db,omitempty"`
	ServiceWorker *bool `json:"service_worker,omitempty"`
	WindowName    *bool `json:"window_name,omitempty"`
}

type SessionDeviceAvailabilityCapability struct {
	Available *bool `json:"available,omitempty"`
}

type SessionDevicePlatformAuthenticatorCapabilities struct {
	Available            *bool `json:"available,omitempty"`
	ConditionalMediation *bool `json:"conditional_mediation,omitempty"`
}

type SessionDeviceCapabilities struct {
	Touch                 SessionDeviceTouchCapabilities                 `json:"touch"`
	Storage               SessionDeviceStorageCapabilities               `json:"storage"`
	WebGPU                SessionDeviceAvailabilityCapability            `json:"webgpu"`
	PlatformAuthenticator SessionDevicePlatformAuthenticatorCapabilities `json:"platform_authenticator"`
	MediaDevices          SessionDeviceAvailabilityCapability            `json:"media_devices"`
	SpeechSynthesis       SessionDeviceAvailabilityCapability            `json:"speech_synthesis"`
}

type DeviceFormFactor string

const (
	DeviceFormFactorDesktop DeviceFormFactor = "desktop"
	DeviceFormFactorPhone   DeviceFormFactor = "phone"
	DeviceFormFactorTablet  DeviceFormFactor = "tablet"
	DeviceFormFactorUnknown DeviceFormFactor = "unknown"
)

type SessionDevice struct {
	FormFactor      DeviceFormFactor             `json:"form_factor"`
	OperatingSystem SessionDeviceOperatingSystem `json:"operating_system"`
	Architecture    *string                      `json:"architecture,omitempty"`
	Screen          SessionDeviceScreen          `json:"screen"`
	Locale          SessionDeviceLocale          `json:"locale"`
	Capabilities    SessionDeviceCapabilities    `json:"capabilities"`
}

type SessionRawDeviceNavigator struct {
	Platform            *string  `json:"platform,omitempty"`
	Vendor              *string  `json:"vendor,omitempty"`
	HardwareConcurrency *int     `json:"hardware_concurrency,omitempty"`
	DeviceMemory        *float64 `json:"device_memory,omitempty"`
	MaxTouchPoints      *int     `json:"max_touch_points,omitempty"`
	PDFViewerEnabled    *bool    `json:"pdf_viewer_enabled,omitempty"`
	CookieEnabled       *bool    `json:"cookie_enabled,omitempty"`
	ProductSub          *string  `json:"product_sub,omitempty"`
	PrimaryLanguage     *string  `json:"primary_language,omitempty"`
	Languages           []string `json:"languages"`
	MimeTypesCount      *int     `json:"mime_types_count,omitempty"`
	Plugins             []string `json:"plugins"`
}

type SessionRawDeviceStorage struct {
	Cookies        *bool `json:"cookies,omitempty"`
	LocalStorage   *bool `json:"local_storage,omitempty"`
	SessionStorage *bool `json:"session_storage,omitempty"`
	IndexedDB      *bool `json:"indexed_db,omitempty"`
	ServiceWorker  *bool `json:"service_worker,omitempty"`
	WindowName     *bool `json:"window_name,omitempty"`
}

type SessionRawDeviceCanvas struct {
	Hash any `json:"hash"`
	// Deprecated: Always null. Foil stopped measuring canvas geometry in July 2026.
	GeometryHash any `json:"geometry_hash"`
	TextHash     any `json:"text_hash"`
	// Deprecated: Always null. No Foil probe measures canvas winding.
	Winding             *bool `json:"winding,omitempty"`
	NoiseDetected       *bool `json:"noise_detected,omitempty"`
	OffscreenConsistent *bool `json:"offscreen_consistent,omitempty"`
}

type SessionRawDeviceGraphicsWebGL struct {
	Vendor                 *string `json:"vendor,omitempty"`
	Renderer               *string `json:"renderer,omitempty"`
	Version                *string `json:"version,omitempty"`
	ShadingLanguageVersion *string `json:"shading_language_version,omitempty"`
	ParametersHash         any     `json:"parameters_hash"`
	ExtensionsHash         any     `json:"extensions_hash"`
	// Deprecated: Use WebGL2ParametersHash, which carries the same value.
	ExtensionParametersHash any `json:"extension_parameters_hash"`
	WebGL2ParametersHash    any `json:"webgl2_parameters_hash"`
	ShaderPrecisionHash     any `json:"shader_precision_hash"`
}

type SessionRawDeviceGraphicsWebGPU struct {
	Available           *bool   `json:"available,omitempty"`
	AdapterVendor       *string `json:"adapter_vendor,omitempty"`
	AdapterArchitecture *string `json:"adapter_architecture,omitempty"`
	FallbackAdapter     *bool   `json:"fallback_adapter,omitempty"`
	FeaturesHash        any     `json:"features_hash"`
	LimitsHash          any     `json:"limits_hash"`
}

type SessionRawDeviceGraphics struct {
	WebGL  SessionRawDeviceGraphicsWebGL  `json:"webgl"`
	WebGPU SessionRawDeviceGraphicsWebGPU `json:"webgpu"`
}

type SessionRawDeviceAudio struct {
	Hash             any      `json:"hash"`
	SampleRate       *float64 `json:"sample_rate,omitempty"`
	ChannelCount     *int     `json:"channel_count,omitempty"`
	VoiceCount       *int     `json:"voice_count,omitempty"`
	LocalVoiceCount  *int     `json:"local_voice_count,omitempty"`
	DefaultVoiceLang *string  `json:"default_voice_lang,omitempty"`
	NoiseDetected    *bool    `json:"noise_detected,omitempty"`
}

type SessionRawDeviceFonts struct {
	DetectedCount   *int `json:"detected_count,omitempty"`
	TestedCount     *int `json:"tested_count,omitempty"`
	EnumerationHash any  `json:"enumeration_hash"`
	MetricsHash     any  `json:"metrics_hash"`
	PreferencesHash any  `json:"preferences_hash"`
	EmojiHash       any  `json:"emoji_hash"`
}

type SessionRawDeviceMedia struct {
	DeviceCount     *int           `json:"device_count,omitempty"`
	CountsByKind    map[string]int `json:"counts_by_kind"`
	BlankLabelCount *int           `json:"blank_label_count,omitempty"`
	TopologyHash    any            `json:"topology_hash"`
}

type SessionClientTelemetry struct {
	Navigator SessionRawDeviceNavigator `json:"navigator"`
	Storage   SessionRawDeviceStorage   `json:"storage"`
	Canvas    SessionRawDeviceCanvas    `json:"canvas"`
	Graphics  SessionRawDeviceGraphics  `json:"graphics"`
	Audio     SessionRawDeviceAudio     `json:"audio"`
	Fonts     SessionRawDeviceFonts     `json:"fonts"`
	Media     SessionRawDeviceMedia     `json:"media"`
}

type SessionSummary struct {
	Object              string                  `json:"object"`
	ID                  string                  `json:"id"`
	CreatedAt           *string                 `json:"created_at,omitempty"`
	ClientUserID        *string                 `json:"client_user_id,omitempty"`
	ClientApplicationID *string                 `json:"client_application_id,omitempty"`
	LatestDecision      Decision                `json:"latest_decision"`
	VisitorFingerprint  *VisitorFingerprintLink `json:"visitor_fingerprint,omitempty"`
}

type SessionNetworkLocation struct {
	City             *string  `json:"city,omitempty"`
	Region           *string  `json:"region,omitempty"`
	Country          *string  `json:"country,omitempty"`
	CountryCode      *string  `json:"country_code,omitempty"`
	Latitude         *float64 `json:"latitude,omitempty"`
	Longitude        *float64 `json:"longitude,omitempty"`
	Timezone         *string  `json:"timezone,omitempty"`
	PostalCode       *string  `json:"postal_code,omitempty"`
	AccuracyRadiusKm *float64 `json:"accuracy_radius_km,omitempty"`
}

type SessionNetworkRouting struct {
	ASN          *string `json:"asn,omitempty"`
	Organization *string `json:"organization,omitempty"`
}

type SessionNetworkAnonymity struct {
	VPN              bool `json:"vpn"`
	Proxy            bool `json:"proxy"`
	Tor              bool `json:"tor"`
	Relay            bool `json:"relay"`
	Hosting          bool `json:"hosting"`
	ResidentialProxy bool `json:"residential_proxy"`
	// Deprecated: Always false for sessions evaluated since July 30, 2026.
	CallbackProxy bool    `json:"callback_proxy"`
	Provider      *string `json:"provider,omitempty"`
}

type SessionNetworkReputation struct {
	Listed     bool     `json:"listed"`
	Categories []string `json:"categories"`
	// Deprecated: Always false for sessions evaluated since July 30, 2026.
	SuspiciousNetwork bool `json:"suspicious_network"`
}

type SessionNetworkEvidence struct {
	RiskSignals  []string `json:"risk_signals"`
	OperatorTags []string `json:"operator_tags"`
	ClientTypes  []string `json:"client_types"`
	// Deprecated: Always null for sessions evaluated since July 30, 2026.
	ClientCount *int `json:"client_count,omitempty"`
}

type SessionNetwork struct {
	IPAddress   *string                  `json:"ip_address,omitempty"`
	IPVersion   *string                  `json:"ip_version,omitempty"`
	Status      string                   `json:"status"`
	Summary     *string                  `json:"summary,omitempty"`
	Location    *SessionNetworkLocation  `json:"location,omitempty"`
	Routing     SessionNetworkRouting    `json:"routing"`
	Anonymity   SessionNetworkAnonymity  `json:"anonymity"`
	Reputation  SessionNetworkReputation `json:"reputation"`
	Evidence    SessionNetworkEvidence   `json:"evidence"`
	EvaluatedAt *string                  `json:"evaluated_at,omitempty"`
}

type SessionDetailVisitorFingerprintLifecycle struct {
	FirstSeenAt *string `json:"first_seen_at,omitempty"`
	LastSeenAt  *string `json:"last_seen_at,omitempty"`
	// Deprecated: Use ResolutionCount, which carries the same value. It does not match
	// SeenCount from Fingerprints.Get, which counts sessions.
	SeenCount       *int `json:"seen_count,omitempty"`
	ResolutionCount *int `json:"resolution_count,omitempty"`
}

type VisitorResolutionMethod string

const (
	VisitorResolutionMethodContinuity       VisitorResolutionMethod = "continuity"
	VisitorResolutionMethodFingerprint      VisitorResolutionMethod = "fingerprint"
	VisitorResolutionMethodNew              VisitorResolutionMethod = "new"
	VisitorResolutionMethodNativePersistent VisitorResolutionMethod = "native_persistent"
	VisitorResolutionMethodNativeNew        VisitorResolutionMethod = "native_new"
)

type VisitorResolutionBasis string

const (
	VisitorResolutionBasisDeterministic VisitorResolutionBasis = "deterministic"
	VisitorResolutionBasisProbabilistic VisitorResolutionBasis = "probabilistic"
)

type SessionDetailVisitorFingerprintResolution struct {
	Method *VisitorResolutionMethod `json:"method,omitempty"`
	Basis  *VisitorResolutionBasis  `json:"basis,omitempty"`
}

type SessionDetailVisitorFingerprint struct {
	Object       string                                    `json:"object"`
	ID           string                                    `json:"id"`
	Confidence   *int                                      `json:"confidence,omitempty"`
	IdentifiedAt *string                                   `json:"identified_at,omitempty"`
	Resolution   SessionDetailVisitorFingerprintResolution `json:"resolution"`
	Lifecycle    SessionDetailVisitorFingerprintLifecycle  `json:"lifecycle"`
}

type JA4Profile string

const (
	JA4ProfileChrome  JA4Profile = "chrome"
	JA4ProfileFirefox JA4Profile = "firefox"
	JA4ProfileSafari  JA4Profile = "safari"
	JA4ProfileNodeJS  JA4Profile = "nodejs"
	JA4ProfilePython  JA4Profile = "python"
	JA4ProfileGo      JA4Profile = "go"
	JA4ProfileCurl    JA4Profile = "curl"
	JA4ProfileUnknown JA4Profile = "unknown"
)

type JA4Family string

const (
	JA4FamilyChromium  JA4Family = "chromium"
	JA4FamilyFirefox   JA4Family = "firefox"
	JA4FamilyWebKit    JA4Family = "webkit"
	JA4FamilyNodeJS    JA4Family = "nodejs"
	JA4FamilyPython    JA4Family = "python"
	JA4FamilyGo        JA4Family = "go"
	JA4FamilyCurl      JA4Family = "curl"
	JA4FamilyMiddlebox JA4Family = "middlebox"
	JA4FamilyUnknown   JA4Family = "unknown"
)

type JA4Confidence string

const (
	JA4ConfidenceExact     JA4Confidence = "exact"
	JA4ConfidenceStrong    JA4Confidence = "strong"
	JA4ConfidenceHeuristic JA4Confidence = "heuristic"
	JA4ConfidenceUnknown   JA4Confidence = "unknown"
)

type JA4MatchSource string

const (
	JA4MatchSourceFullJA4    JA4MatchSource = "full_ja4"
	JA4MatchSourceBCHash     JA4MatchSource = "bc_hash"
	JA4MatchSourceCipherHash JA4MatchSource = "cipher_hash"
	JA4MatchSourcePrefix     JA4MatchSource = "prefix"
	JA4MatchSourceUnknown    JA4MatchSource = "unknown"
)

type JA4ReferenceSource string

const (
	JA4ReferenceSourceMasterDB  JA4ReferenceSource = "master_db"
	JA4ReferenceSourceHeuristic JA4ReferenceSource = "heuristic"
)

type SessionConnectionFingerprintJA4Range struct {
	Min int `json:"min"`
	Max int `json:"max"`
}

type SessionConnectionFingerprintJA4Reference struct {
	Label       string                               `json:"label"`
	Source      JA4ReferenceSource                   `json:"source"`
	Protocol    string                               `json:"protocol"`
	TLSVersion  string                               `json:"tls_version"`
	SNI         string                               `json:"sni"`
	ALPN        string                               `json:"alpn"`
	CipherCount SessionConnectionFingerprintJA4Range `json:"cipher_count"`
	ExtCount    SessionConnectionFingerprintJA4Range `json:"ext_count"`
	SampleSize  *int                                 `json:"sample_size,omitempty"`
	Verified    *bool                                `json:"verified,omitempty"`
}

type SessionConnectionFingerprintJA4 struct {
	Hash          *string                                   `json:"hash,omitempty"`
	Profile       *JA4Profile                               `json:"profile,omitempty"`
	Family        *JA4Family                                `json:"family,omitempty"`
	Product       *string                                   `json:"product,omitempty"`
	Confidence    *JA4Confidence                            `json:"confidence,omitempty"`
	MatchSource   *JA4MatchSource                           `json:"match_source,omitempty"`
	Deterministic *bool                                     `json:"deterministic,omitempty"`
	Reference     *SessionConnectionFingerprintJA4Reference `json:"reference,omitempty"`
}

type HTTP2Profile string

const (
	HTTP2ProfileChrome  HTTP2Profile = "chrome"
	HTTP2ProfileFirefox HTTP2Profile = "firefox"
	HTTP2ProfileSafari  HTTP2Profile = "safari"
	HTTP2ProfileEdge    HTTP2Profile = "edge"
	HTTP2ProfileNodeJS  HTTP2Profile = "nodejs"
	HTTP2ProfilePython  HTTP2Profile = "python"
	HTTP2ProfileGo      HTTP2Profile = "go"
	HTTP2ProfileCurl    HTTP2Profile = "curl"
	HTTP2ProfileUnknown HTTP2Profile = "unknown"
)

type SessionConnectionFingerprintHTTP2 struct {
	AkamaiFingerprint *string       `json:"akamai_fingerprint,omitempty"`
	Profile           *HTTP2Profile `json:"profile,omitempty"`
}

type SessionConnectionFingerprint struct {
	JA4                SessionConnectionFingerprintJA4   `json:"ja4"`
	HTTP2              SessionConnectionFingerprintHTTP2 `json:"http2"`
	UserAgentAlignment *string                           `json:"user_agent_alignment,omitempty"`
}

type SessionAnalysisCoverage struct {
	Browser         bool `json:"browser"`
	Device          bool `json:"device"`
	Network         bool `json:"network"`
	Runtime         bool `json:"runtime"`
	Behavioral      bool `json:"behavioral"`
	VisitorIdentity bool `json:"visitor_identity"`
}

type SessionSignalRole string

const (
	SessionSignalRoleRisk SessionSignalRole = "risk"
)

type SessionSignalCategory string

const (
	SessionSignalCategoryEnvironment SessionSignalCategory = "environment"
	SessionSignalCategoryEventTrust  SessionSignalCategory = "eventTrust"
	SessionSignalCategoryFingerprint SessionSignalCategory = "fingerprint"
	SessionSignalCategoryMouse       SessionSignalCategory = "mouse"
	SessionSignalCategoryKeyboard    SessionSignalCategory = "keyboard"
	SessionSignalCategoryTiming      SessionSignalCategory = "timing"
	SessionSignalCategoryTouch       SessionSignalCategory = "touch"
)

type SessionSignalStrength string

const (
	SessionSignalStrengthDefinitive SessionSignalStrength = "definitive"
	SessionSignalStrengthStrong     SessionSignalStrength = "strong"
	SessionSignalStrengthBehavioral SessionSignalStrength = "behavioral"
	SessionSignalStrengthWeak       SessionSignalStrength = "weak"
)

type SessionSignalFired struct {
	Signal      string                `json:"signal"`
	Role        SessionSignalRole     `json:"role"`
	Category    SessionSignalCategory `json:"category"`
	Strength    SessionSignalStrength `json:"strength"`
	SignalScore int                   `json:"signal_score"`
}

type SessionDetail struct {
	Object                 string                           `json:"object"`
	ID                     string                           `json:"id"`
	CreatedAt              *string                          `json:"created_at,omitempty"`
	ClientUserID           *string                          `json:"client_user_id,omitempty"`
	ClientApplicationID    *string                          `json:"client_application_id,omitempty"`
	Decision               SessionDecision                  `json:"decision"`
	Highlights             []SessionHighlight               `json:"highlights"`
	Attribution            *SessionAttribution              `json:"attribution,omitempty"`
	WebBotAuth             *SessionWebBotAuth               `json:"web_bot_auth,omitempty"`
	Network                SessionNetwork                   `json:"network"`
	RuntimeIntegrity       SessionRuntimeIntegrity          `json:"runtime_integrity"`
	NativeRuntimeIntegrity map[string]any                   `json:"native_runtime_integrity,omitempty"`
	NativeApp              map[string]any                   `json:"native_app,omitempty"`
	NativeCarrier          map[string]any                   `json:"native_carrier,omitempty"`
	NativeMotionPrint      map[string]any                   `json:"native_motion_print,omitempty"`
	DeviceIdentity         map[string]any                   `json:"device_identity,omitempty"`
	InstallID              *string                          `json:"install_id,omitempty"`
	VisitorFingerprint     *SessionDetailVisitorFingerprint `json:"visitor_fingerprint,omitempty"`
	ConnectionFingerprint  SessionConnectionFingerprint     `json:"connection_fingerprint"`
	PreviousDecisions      []SessionDecision                `json:"previous_decisions"`
	Request                SessionDetailRequest             `json:"request"`
	Browser                SessionBrowser                   `json:"browser"`
	Device                 SessionDevice                    `json:"device"`
	AnalysisCoverage       SessionAnalysisCoverage          `json:"analysis_coverage"`
	SignalsFired           []SessionSignalFired             `json:"signals_fired"`
	ClientTelemetry        SessionClientTelemetry           `json:"client_telemetry"`
}

type VisitorFingerprintLifecycle struct {
	FirstSeenAt string `json:"first_seen_at"`
	LastSeenAt  string `json:"last_seen_at"`
	SeenCount   int    `json:"seen_count"`
	// Deprecated: Always 9999-12-31T00:00:00.000Z. Visitor IDs do not expire.
	ExpiresAt string `json:"expires_at"`
}

type VisitorFingerprintLatestRequest struct {
	UserAgent string `json:"user_agent"`
	IPAddress string `json:"ip_address"`
}

// Deprecated: Every flag is false from Fingerprints.List and Fingerprints.Get.
// Treat the flags as unavailable, not as evidence that a store was absent.
type VisitorFingerprintStorage struct {
	Cookies       bool `json:"cookies"`
	LocalStorage  bool `json:"local_storage"`
	IndexedDB     bool `json:"indexed_db"`
	ServiceWorker bool `json:"service_worker"`
	WindowName    bool `json:"window_name"`
}

// Deprecated: Every hash is null from Fingerprints.List and Fingerprints.Get.
type VisitorFingerprintAnchors struct {
	WebGLHash      *string `json:"webgl_hash,omitempty"`
	ParametersHash *string `json:"parameters_hash,omitempty"`
	AudioHash      *string `json:"audio_hash,omitempty"`
}

type VisitorFingerprintSummary struct {
	Object        string                          `json:"object"`
	ID            string                          `json:"id"`
	Lifecycle     VisitorFingerprintLifecycle     `json:"lifecycle"`
	LatestRequest VisitorFingerprintLatestRequest `json:"latest_request"`
	// Deprecated: Every flag is false. Treat the flags as unavailable.
	Storage VisitorFingerprintStorage `json:"storage"`
	// Deprecated: Every hash is null.
	Anchors VisitorFingerprintAnchors `json:"anchors"`
}

type ScoreBreakdown struct {
	Categories map[string]int `json:"categories"`
}

type VisitorFingerprintSessionSummary struct {
	SessionID      string         `json:"session_id"`
	Decision       Decision       `json:"decision"`
	Request        RequestContext `json:"request"`
	ScoreBreakdown ScoreBreakdown `json:"score_breakdown"`
}

type VisitorFingerprintComponents struct {
	// Deprecated: Always empty from Fingerprints.Get.
	Vector []int `json:"vector"`
}

type VisitorFingerprintActivity struct {
	Sessions []VisitorFingerprintSessionSummary `json:"sessions"`
}

type VisitorFingerprintDetail struct {
	VisitorFingerprintSummary
	Components VisitorFingerprintComponents `json:"components"`
	Activity   VisitorFingerprintActivity   `json:"activity"`
}

type Organization struct {
	Object    string  `json:"object"`
	ID        string  `json:"id"`
	Name      string  `json:"name"`
	Slug      string  `json:"slug"`
	Status    string  `json:"status"`
	CreatedAt string  `json:"created_at"`
	UpdatedAt *string `json:"updated_at,omitempty"`
}

type APIKey struct {
	Object         string   `json:"object"`
	ID             string   `json:"id"`
	Type           string   `json:"type"`
	Name           string   `json:"name"`
	Environment    string   `json:"environment"`
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
	RateLimit      *int     `json:"rate_limit,omitempty"`
	Status         string   `json:"status"`
	KeyPreview     string   `json:"key_preview"`
	DisplayKey     *string  `json:"display_key,omitempty"`
	LastUsedAt     *string  `json:"last_used_at,omitempty"`
	CreatedAt      string   `json:"created_at"`
	RotatedAt      *string  `json:"rotated_at,omitempty"`
	RevokedAt      *string  `json:"revoked_at,omitempty"`
	GraceExpiresAt *string  `json:"grace_expires_at,omitempty"`
}

type IssuedAPIKey struct {
	APIKey
	RevealedKey string `json:"revealed_key"`
}

type VerifiedFoilSignal struct {
	ID         string         `json:"id"`
	Category   string         `json:"category"`
	Confidence string         `json:"confidence"`
	Score      int            `json:"score"`
	Raw        map[string]any `json:"raw,omitempty"`
}

type Attribution struct {
	Bot map[string]any `json:"bot,omitempty"`
	Raw map[string]any `json:"raw,omitempty"`
}

type VerifiedFoilToken struct {
	Object             string                  `json:"object"`
	SessionID          string                  `json:"session_id"`
	Decision           Decision                `json:"decision"`
	Request            RequestContext          `json:"request"`
	VisitorFingerprint *VisitorFingerprintLink `json:"visitor_fingerprint,omitempty"`
	Signals            []VerifiedFoilSignal    `json:"signals"`
	ScoreBreakdown     ScoreBreakdown          `json:"score_breakdown"`
	Attribution        Attribution             `json:"attribution"`
	Embed              map[string]any          `json:"embed,omitempty"`
	Raw                map[string]any          `json:"raw,omitempty"`
}

type VerificationResult struct {
	OK    bool
	Data  *VerifiedFoilToken
	Error error
}

type SessionListParams struct {
	Limit               int
	Cursor              string
	Verdict             string
	Search              string
	ClientUserID        string
	ClientApplicationID string
}

type FingerprintListParams struct {
	Limit  int
	Cursor string
	Search string
	Sort   string
}

type APIKeyListParams struct {
	Limit  int
	Cursor string
}

type CreateOrganizationParams struct {
	Name string `json:"name"`
	Slug string `json:"slug"`
}

type UpdateOrganizationParams struct {
	Name   string `json:"name,omitempty"`
	Status string `json:"status,omitempty"`
}

type CreateAPIKeyParams struct {
	Name           string   `json:"name,omitempty"`
	Type           string   `json:"type,omitempty"`
	Environment    string   `json:"environment,omitempty"`
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}

type UpdateAPIKeyParams struct {
	Name           string   `json:"name,omitempty"`
	AllowedOrigins []string `json:"allowed_origins,omitempty"`
	Scopes         []string `json:"scopes,omitempty"`
}
