package iso8583

// FieldType describes how a data element is encoded in an ISO 8583 message.
type FieldType uint8

const (
	// TypeN fixed-length numeric, right-justified, zero padded.
	TypeN FieldType = iota
	// TypeA fixed-length alpha, left-justified, space padded.
	TypeA
	// TypeAN fixed-length alphanumeric.
	TypeAN
	// TypeANS fixed-length alphanumeric and special characters.
	TypeANS
	// TypeNLLVAR numeric variable length with a 2-digit length prefix.
	TypeNLLVAR
	// TypeNLLLVAR numeric variable length with a 3-digit length prefix.
	TypeNLLLVAR
	// TypeANLLVAR alphanumeric variable length with a 2-digit length prefix.
	TypeANLLVAR
	// TypeANLLLVAR alphanumeric variable length with a 3-digit length prefix.
	TypeANLLLVAR
)

// FieldSpec describes the encoding of a single data element.
type FieldSpec struct {
	Type   FieldType
	MaxLen int
}

// DataElements holds the ISO 8583:1987 data element definitions supported by
// Clara Network. Only the elements used by the network are declared.
var DataElements = map[int]FieldSpec{
	2:   {TypeNLLVAR, 19},    // Primary account number (PAN)
	3:   {TypeN, 6},          // Processing code
	4:   {TypeN, 12},         // Amount, transaction
	7:   {TypeN, 10},         // Transmission date and time (MMDDhhmmss)
	11:  {TypeN, 6},          // Systems trace audit number (STAN)
	12:  {TypeN, 6},          // Local transaction time (hhmmss)
	13:  {TypeN, 4},          // Local transaction date (MMDD)
	22:  {TypeN, 3},          // POS entry mode
	24:  {TypeN, 3},          // Function code
	25:  {TypeN, 2},          // POS condition code
	32:  {TypeNLLVAR, 11},    // Acquiring institution ID code
	37:  {TypeAN, 12},        // Retrieval reference number
	39:  {TypeAN, 3},         // Response code
	41:  {TypeANS, 8},        // Card acceptor terminal identification
	42:  {TypeANS, 15},       // Card acceptor identification code
	49:  {TypeN, 3},          // Currency code
	62:  {TypeANLLLVAR, 999}, // Private use (stand-in marker)
	70:  {TypeN, 3},          // Network management information code
	100: {TypeNLLVAR, 11},    // Receiving institution ID code
}
