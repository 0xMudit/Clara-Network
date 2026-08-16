// Package instant implements the Clara Network instant-payment layer
// (docs/24, docs/25 §25.4 phase 10): ISO 20022 pacs.008 customer credit
// transfers settled in real time, 24/7/365, against fully prefunded member
// positions (the RTP model) with a 20-second scheme SLA. Before a payment is
// forwarded the engine verifies and reserves the sender's settlement
// capacity; insufficient prefunded position rejects the payment (overdrafts
// are not permitted), and a beneficiary that does not confirm within the SLA
// causes the reservation to be released.
//
// Money is always carried in minor units (cents) as int64 to avoid rounding.
package instant

import (
	"encoding/xml"
	"fmt"
	"strings"
)

// pacs008XMLNS is the ISO 20022 namespace for pacs.008.001.09 (CBPR+).
const pacs008XMLNS = "urn:iso:std:iso:20022:tech:xsd:pacs.008.001.09"

// Payment is a pacs.008 FIToFICustomerCreditTransfer as seen by the scheme.
type Payment struct {
	MsgID           string // GroupHeader.MsgId
	InstrID         string // PmtId.InstrId
	EndToEndID      string // PmtId.EndToEndId (UETR)
	TxID            string // PmtId.TxId
	Sender          string // DbtrAgt.FinInstnId.ClrSysMmbId.MmbId
	Beneficiary     string // CdtrAgt.FinInstnId.ClrSysMmbId.MmbId
	SenderIBAN      string // DbtrAcct.Id.IBAN
	BeneficiaryIBAN string // CdtrAcct.Id.IBAN
	AmountMinor     int64  // InstdAmt as minor units
	Currency        string // InstdAmt @Ccy
	Remittance      string // RmtInf.Ustrd
	CreatedAt       string // GroupHeader.CreDtTm (UTC, XML timestamp)
}

type pacs008Document struct {
	XMLName           xml.Name    `xml:"Document"`
	XMLNS             string      `xml:"xmlns,attr"`
	FIToFICstmrCdtTrf pacs008Body `xml:"FIToFICstmrCdtTrf"`
}

type pacs008Body struct {
	GrpHdr   pacs008GrpHdr `xml:"GrpHdr"`
	CdtTrfTx pacs008Tx     `xml:"CdtTrfTxInf"`
}

type pacs008GrpHdr struct {
	MsgID    string       `xml:"MsgId"`
	CreDtTm  string       `xml:"CreDtTm"`
	NbOfTxs  string       `xml:"NbOfTxs"`
	SttlmInf pacs008Sttlm `xml:"SttlmInf"`
}

type pacs008Sttlm struct {
	SttlmMtd string `xml:"SttlmMtd"`
}

type pacs008Tx struct {
	PmtId    pacs008PmtId  `xml:"PmtId"`
	Amt      pacs008Amt    `xml:"Amt"`
	DbtrAgt  pacs008Agent  `xml:"DbtrAgt"`
	CdtrAgt  pacs008Agent  `xml:"CdtrAgt"`
	DbtrAcct *pacs008Acct  `xml:"DbtrAcct,omitempty"`
	CdtrAcct pacs008Acct   `xml:"CdtrAcct"`
	RmtInf   pacs008RmtInf `xml:"RmtInf"`
}

type pacs008PmtId struct {
	InstrId    string `xml:"InstrId"`
	EndToEndId string `xml:"EndToEndId"`
	TxId       string `xml:"TxId"`
}

type pacs008Amt struct {
	InstdAmt pacs008Amount `xml:"InstdAmt"`
}

type pacs008Amount struct {
	Ccy  string `xml:"Ccy,attr"`
	Text string `xml:",chardata"`
}

type pacs008Agent struct {
	FinInstnId pacs008FinInstn `xml:"FinInstnId"`
}

type pacs008FinInstn struct {
	ClrSysMmbId pacs008Mmb `xml:"ClrSysMmbId"`
}

type pacs008Mmb struct {
	MmbId string `xml:"MmbId"`
}

type pacs008Acct struct {
	Id pacs008AcctId `xml:"Id"`
}

type pacs008AcctId struct {
	IBAN string `xml:"IBAN"`
}

type pacs008RmtInf struct {
	Ustrd string `xml:"Ustrd"`
}

// BuildPacs008 renders a payment as an ISO 20022 pacs.008 credit transfer.
func BuildPacs008(p Payment) ([]byte, error) {
	if err := p.Validate(); err != nil {
		return nil, fmt.Errorf("pacs.008: %w", err)
	}
	doc := pacs008Document{
		XMLNS: pacs008XMLNS,
		FIToFICstmrCdtTrf: pacs008Body{
			GrpHdr: pacs008GrpHdr{
				MsgID:    p.MsgID,
				CreDtTm:  p.CreatedAt,
				NbOfTxs:  "1",
				SttlmInf: pacs008Sttlm{SttlmMtd: "CLRG"},
			},
			CdtTrfTx: pacs008Tx{
				PmtId: pacs008PmtId{
					InstrId:    p.InstrID,
					EndToEndId: p.EndToEndID,
					TxId:       p.TxID,
				},
				Amt:      pacs008Amt{InstdAmt: pacs008Amount{Ccy: p.Currency, Text: formatAmount(p.AmountMinor)}},
				DbtrAgt:  pacs008Agent{FinInstnId: pacs008FinInstn{ClrSysMmbId: pacs008Mmb{MmbId: p.Sender}}},
				CdtrAgt:  pacs008Agent{FinInstnId: pacs008FinInstn{ClrSysMmbId: pacs008Mmb{MmbId: p.Beneficiary}}},
				CdtrAcct: pacs008Acct{Id: pacs008AcctId{IBAN: p.BeneficiaryIBAN}},
				RmtInf:   pacs008RmtInf{Ustrd: p.Remittance},
			},
		},
	}
	if p.SenderIBAN != "" {
		doc.FIToFICstmrCdtTrf.CdtTrfTx.DbtrAcct = &pacs008Acct{Id: pacs008AcctId{IBAN: p.SenderIBAN}}
	}
	return marshalDoc(doc)
}

// ParsePacs008 decodes an ISO 20022 pacs.008 credit transfer.
func ParsePacs008(raw []byte) (Payment, error) {
	var doc pacs008Document
	if err := xml.Unmarshal(raw, &doc); err != nil {
		return Payment{}, fmt.Errorf("pacs.008: parse: %w", err)
	}
	if !strings.Contains(doc.XMLNS, "pacs.008") {
		return Payment{}, fmt.Errorf("pacs.008: unexpected namespace %q", doc.XMLNS)
	}
	p := Payment{
		MsgID:       doc.FIToFICstmrCdtTrf.GrpHdr.MsgID,
		CreatedAt:   doc.FIToFICstmrCdtTrf.GrpHdr.CreDtTm,
		InstrID:     doc.FIToFICstmrCdtTrf.CdtTrfTx.PmtId.InstrId,
		EndToEndID:  doc.FIToFICstmrCdtTrf.CdtTrfTx.PmtId.EndToEndId,
		TxID:        doc.FIToFICstmrCdtTrf.CdtTrfTx.PmtId.TxId,
		Sender:      doc.FIToFICstmrCdtTrf.CdtTrfTx.DbtrAgt.FinInstnId.ClrSysMmbId.MmbId,
		Beneficiary: doc.FIToFICstmrCdtTrf.CdtTrfTx.CdtrAgt.FinInstnId.ClrSysMmbId.MmbId,
		Currency:    doc.FIToFICstmrCdtTrf.CdtTrfTx.Amt.InstdAmt.Ccy,
		Remittance:  doc.FIToFICstmrCdtTrf.CdtTrfTx.RmtInf.Ustrd,
	}
	if a := doc.FIToFICstmrCdtTrf.CdtTrfTx.DbtrAcct; a != nil {
		p.SenderIBAN = a.Id.IBAN
	}
	p.BeneficiaryIBAN = doc.FIToFICstmrCdtTrf.CdtTrfTx.CdtrAcct.Id.IBAN
	minor, err := parseAmount(doc.FIToFICstmrCdtTrf.CdtTrfTx.Amt.InstdAmt.Text)
	if err != nil {
		return Payment{}, fmt.Errorf("pacs.008: amount: %w", err)
	}
	p.AmountMinor = minor
	return p, nil
}

// Validate checks a payment for mandatory fields and ranges.
func (p Payment) Validate() error {
	switch {
	case p.MsgID == "":
		return fmt.Errorf("missing MsgId")
	case p.TxID == "":
		return fmt.Errorf("missing TxId")
	case p.EndToEndID == "":
		return fmt.Errorf("missing EndToEndId")
	case p.Sender == "" || p.Beneficiary == "":
		return fmt.Errorf("missing sender or beneficiary PSP")
	case p.Sender == p.Beneficiary:
		return fmt.Errorf("sender %q cannot also be beneficiary", p.Sender)
	case p.AmountMinor <= 0:
		return fmt.Errorf("amount must be positive")
	case len(p.Currency) != 3 || !isUpper(p.Currency):
		return fmt.Errorf("invalid currency %q", p.Currency)
	}
	return nil
}

func isUpper(s string) bool {
	for i := 0; i < len(s); i++ {
		if s[i] < 'A' || s[i] > 'Z' {
			return false
		}
	}
	return true
}
