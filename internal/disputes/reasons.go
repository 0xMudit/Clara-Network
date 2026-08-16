package disputes

// ReasonCode is a documented, versioned dispute reason code that drives the
// evidence requirements, response timeline, and liability (docs/20 §20.2).
type ReasonCode struct {
	Code             string
	Category         string // fraud | authorization | processing | error
	Description      string
	Liability        string // who carries the loss if the chargeback stands
	RequiredEvidence []string
}

// Reason categories.
const (
	CategoryFraud        = "fraud"
	CategoryAuth         = "authorization"
	CategoryProcessing   = "processing"
	CategoryError        = "error"
)

// Liability targets.
const (
	LiabilityIssuer   = "issuer"
	LiabilityAcquirer = "acquirer"
)

// Response windows in days per reason category (VCR: fraud is streamlined).
const (
	DaysFraud       = 20
	DaysAuth        = 30
	DaysProcessing  = 45
	DaysStandard    = 46
)

// ReasonCodes is the scheme's dispute reason-code catalog.
var ReasonCodes = map[string]ReasonCode{
	"4837": {Code: "4837", Category: CategoryFraud, Description: "Cardholder does not recognize transaction",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"3ds", "chip"}},
	"4840": {Code: "4840", Category: CategoryFraud, Description: "Fraudulent transaction",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"chip", "cvv"}},
	"4849": {Code: "4849", Category: CategoryFraud, Description: "Cardholder dispute - suspected fraud",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"3ds", "device"}},
	"4870": {Code: "4870", Category: CategoryAuth, Description: "Non-receipt of merchandise",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"receipt", "delivery"}},
	"4871": {Code: "4871", Category: CategoryAuth, Description: "Merchandise not as described",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"receipt", "delivery"}},
	"4841": {Code: "4841", Category: CategoryAuth, Description: "Cancelled recurring transaction",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"terms", "receipt"}},
	"4831": {Code: "4831", Category: CategoryProcessing, Description: "Transaction amount differs",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"receipt", "avs"}},
	"4834": {Code: "4834", Category: CategoryProcessing, Description: "Point-of-interaction error",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"receipt"}},
	"13":  {Code: "13", Category: CategoryProcessing, Description: "Duplicate processing",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"refund", "receipt"}},
	"57":  {Code: "57", Category: CategoryError, Description: "Credit not processed",
		Liability: LiabilityAcquirer, RequiredEvidence: []string{"refund"}},
}

// Lookup returns the reason-code definition and the response window in days.
func Lookup(code string) (ReasonCode, bool) {
	rc, ok := ReasonCodes[code]
	return rc, ok
}

// ResponseDays returns the response window for a reason category.
func ResponseDays(category string) int {
	switch category {
	case CategoryFraud:
		return DaysFraud
	case CategoryAuth:
		return DaysAuth
	case CategoryProcessing:
		return DaysProcessing
	default:
		return DaysStandard
	}
}
