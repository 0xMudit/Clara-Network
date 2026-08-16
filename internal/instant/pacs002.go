package instant

import (
	"encoding/xml"
	"fmt"
)

// pacs002XMLNS is the ISO 20022 namespace for pacs.002.001.10.
const pacs002XMLNS = "urn:iso:std:iso:20022:tech:xsd:pacs.002.001.10"

// Transaction statuses used in a pacs.002 status report.
const (
	// StatusACSC AcceptedSettlementCompleted: settled final in real time.
	StatusACSC = "ACSC"
	// StatusACCP AcceptedCreditTransfer: accepted but not yet settled.
	StatusACCP = "ACCP"
	// StatusRJCT Rejected: not accepted.
	StatusRJCT = "RJCT"
)

// Reason codes used when a payment is rejected.
const (
	ReasonInsufficientFunds = "AC04" // insufficient funds in the sender's position
	ReasonAccount           = "AC01" // unknown or incorrect beneficiary account
	ReasonForbidden         = "AG01" // transaction forbidden (unknown sender / self transfer)
	ReasonFormat            = "FF01" // invalid message format
	ReasonNoAnswer          = "NOAS" // beneficiary did not confirm within the SLA
)

// StatusReport is a pacs.002 FIToFIPaymentStatusReport for one payment.
type StatusReport struct {
	MsgID           string // status message id
	OriginalMsgID   string // GroupHeader.MsgId of the original pacs.008
	OriginalMsgName string // e.g. "pacs.008.001.09"
	EndToEndID      string
	TxID            string
	Status          string // ACSC / ACCP / RJCT
	Reason          string // reason code when RJCT
	CreatedAt       string // UTC XML timestamp
}

type pacs002Document struct {
	XMLName         xml.Name    `xml:"Document"`
	XMLNS           string      `xml:"xmlns,attr"`
	FIToFIPmtStsRpt pacs002Body `xml:"FIToFIPmtStsRpt"`
}

type pacs002Body struct {
	GrpHdr      pacs002GrpHdr   `xml:"GrpHdr"`
	OrgnlGrpInf pacs002OrgnlGrp `xml:"OrgnlGrpInf"`
	TxInfAndSts pacs002TxSts    `xml:"TxInfAndSts"`
}

type pacs002GrpHdr struct {
	MsgID   string `xml:"MsgId"`
	CreDtTm string `xml:"CreDtTm"`
}

type pacs002OrgnlGrp struct {
	OrgnlMsgID   string `xml:"OrgnlMsgId"`
	OrgnlMsgNmId string `xml:"OrgnlMsgNmId"`
}

type pacs002TxSts struct {
	OrgnlEndToEndId string        `xml:"OrgnlEndToEndId"`
	OrgnlTxId       string        `xml:"OrgnlTxId"`
	TxSts           string        `xml:"TxSts"`
	StsRsnInf       pacs002RsnInf `xml:"StsRsnInf"`
}

type pacs002RsnInf struct {
	Rsn pacs002Rsn `xml:"Rsn"`
}

type pacs002Rsn struct {
	Cd string `xml:"Cd"`
}

// BuildPacs002 renders a payment status report as ISO 20022 pacs.002.
func BuildPacs002(s StatusReport) ([]byte, error) {
	switch s.Status {
	case StatusACSC, StatusACCP, StatusRJCT:
	default:
		return nil, fmt.Errorf("pacs.002: unknown status %q", s.Status)
	}
	if s.OriginalMsgID == "" || s.TxID == "" || s.EndToEndID == "" {
		return nil, fmt.Errorf("pacs.002: missing original message reference")
	}
	if s.OriginalMsgName == "" {
		s.OriginalMsgName = "pacs.008.001.09"
	}
	doc := pacs002Document{
		XMLNS: pacs002XMLNS,
		FIToFIPmtStsRpt: pacs002Body{
			GrpHdr: pacs002GrpHdr{MsgID: s.MsgID, CreDtTm: s.CreatedAt},
			OrgnlGrpInf: pacs002OrgnlGrp{
				OrgnlMsgID:   s.OriginalMsgID,
				OrgnlMsgNmId: s.OriginalMsgName,
			},
			TxInfAndSts: pacs002TxSts{
				OrgnlEndToEndId: s.EndToEndID,
				OrgnlTxId:       s.TxID,
				TxSts:           s.Status,
			},
		},
	}
	if s.Reason != "" {
		doc.FIToFIPmtStsRpt.TxInfAndSts.StsRsnInf = pacs002RsnInf{Rsn: pacs002Rsn{Cd: s.Reason}}
	}
	return marshalDoc(doc)
}
