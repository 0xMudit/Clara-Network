package acquiring

// MCC is a Merchant Category Code assigned at boarding (docs/23 §23.3). It
// drives interchange, risk tiering, monitoring, and compliance.
type MCC struct {
	Code       string
	Name       string
	Tier       string // low | medium | high
	EnhancedDD bool   // enhanced due diligence required (high-integrity-risk)
	RateBPS    int64  // default MDR rate in basis points
}

// Tier constants.
const (
	TierLow    = "low"
	TierMedium = "medium"
	TierHigh   = "high"
)

// MCCs is the scheme's Merchant Category Code catalog.
var MCCs = map[string]MCC{
	"5411": {Code: "5411", Name: "Grocery stores, supermarkets", Tier: TierLow, RateBPS: 150},
	"5812": {Code: "5812", Name: "Eating places, restaurants", Tier: TierLow, RateBPS: 220},
	"5814": {Code: "5814", Name: "Fast food restaurants", Tier: TierLow, RateBPS: 180},
	"5912": {Code: "5912", Name: "Drug stores, pharmacies", Tier: TierLow, RateBPS: 150},
	"5999": {Code: "5999", Name: "Miscellaneous retail", Tier: TierLow, RateBPS: 250},
	"7011": {Code: "7011", Name: "Lodging, hotels", Tier: TierLow, RateBPS: 250},
	"5734": {Code: "5734", Name: "Computer software stores", Tier: TierMedium, RateBPS: 280},
	"4814": {Code: "4814", Name: "Telecommunication services", Tier: TierMedium, RateBPS: 290},
	"5967": {Code: "5967", Name: "Direct marketing", Tier: TierMedium, RateBPS: 290},
	"6051": {Code: "6051", Name: "Non-financial money transfer", Tier: TierHigh, EnhancedDD: true, RateBPS: 350},
	"7995": {Code: "7995", Name: "Betting, lottery, gambling", Tier: TierHigh, EnhancedDD: true, RateBPS: 400},
	"7841": {Code: "7841", Name: "Adult entertainment", Tier: TierHigh, EnhancedDD: true, RateBPS: 400},
}

// LookupMCC returns the MCC definition for a code.
func LookupMCC(code string) (MCC, bool) {
	m, ok := MCCs[code]
	return m, ok
}
