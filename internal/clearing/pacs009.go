package clearing

import (
	"bytes"
	"encoding/xml"
	"fmt"
)

const pacs009XMLNS = "urn:iso:std:iso:20022:tech:xsd:pacs.009.001.08"

// The structs below mirror the pacs.009 (FI credit transfer) message shape.
// Real integration would use the CBPR+ schema (pacs.009.001.09); this subset
// carries the settlement-critical fields.

type pacsDocument struct {
	XMLName  xml.Name `xml:"Document"`
	XMLNS    string   `xml:"xmlns,attr"`
	FICdtTrf pacs009  `xml:"FICdtTrf"`
}

type pacs009 struct {
	GrpHdr   pacsGrpHdr   `xml:"GrpHdr"`
	CdtTrfTx pacsCdtTrfTx `xml:"CdtTrfTxInf"`
}

type pacsGrpHdr struct {
	MsgID      string       `xml:"MsgId"`
	CreDtTm    string       `xml:"CreDtTm"`
	NbOfTxs    string       `xml:"NbOfTxs"`
	SttlmInf   pacsSttlmInf `xml:"SttlmInf"`
}

type pacsSttlmInf struct {
	SttlmMtd string `xml:"SttlmMtd"`
	SttlmDt  string `xml:"SttlmDt"`
}

type pacsCdtTrfTx struct {
	PmtId      pacsPmtId    `xml:"PmtId"`
	Amt        pacsAmt      `xml:"Amt"`
	ForDbtAgt  *pacsAgent   `xml:"ForDbtAgt,omitempty"`
	CdtrAgt    pacsAgent    `xml:"CdtrAgt"`
	Ustrd      string       `xml:"Ustrd"`
}

type pacsPmtId struct {
	InstrId   string `xml:"InstrId"`
	EndToEndId string `xml:"EndToEndId"`
}

type pacsAmt struct {
	InstdAmt pacsAmount `xml:"InstdAmt"`
}

type pacsAmount struct {
	Ccy  string `xml:"Ccy,attr"`
	Text string `xml:",chardata"`
}

type pacsAgent struct {
	FinInstnId pacsFinInstnId `xml:"FinInstnId"`
}

type pacsFinInstnId struct {
	ClrSysMmbId pacsClrSysMmbId `xml:"ClrSysMmbId"`
}

type pacsClrSysMmbId struct {
	MmbId string `xml:"MmbId"`
}

// Pacs009XML renders a settlement instruction as an ISO 20022 pacs.009
// credit transfer: the scheme's settlement agent moves the member's net
// amount either from the member (DEBIT, into the scheme omnibus) or to the
// member (CREDIT).
func Pacs009XML(inst SettlementInstruction) ([]byte, error) {
	if inst.Amount <= 0 {
		return nil, fmt.Errorf("pacs.009: non-positive amount")
	}
	doc := pacsDocument{
		XMLNS: pacs009XMLNS,
		FICdtTrf: pacs009{
			GrpHdr: pacsGrpHdr{
				MsgID:   inst.MsgID,
				CreDtTm: inst.Instruction.UTC().Format("2006-01-02T15:04:05.000Z"),
				NbOfTxs: "1",
				SttlmInf: pacsSttlmInf{
					SttlmMtd: "CLRG",
					SttlmDt:  inst.Instruction.UTC().Format("2006-01-02"),
				},
			},
			CdtTrfTx: pacsCdtTrfTx{
				PmtId: pacsPmtId{
					InstrId:   inst.MsgID,
					EndToEndId: inst.CycleID + ":" + inst.Member,
				},
				Amt: pacsAmt{InstdAmt: pacsAmount{Ccy: inst.Currency, Text: FormatAmount(inst.Amount)}},
				Ustrd: fmt.Sprintf("%s net settlement cycle %s", inst.Direction, inst.CycleID),
			},
		},
	}

	switch inst.Direction {
	case DirCredit:
		doc.FICdtTrf.CdtTrfTx.CdtrAgt = pacsAgent{FinInstnId: pacsFinInstnId{ClrSysMmbId: pacsClrSysMmbId{MmbId: inst.Member}}}
	case DirDebit:
		// Funds flow from the member into the scheme omnibus account.
		doc.FICdtTrf.CdtTrfTx.CdtrAgt = pacsAgent{FinInstnId: pacsFinInstnId{ClrSysMmbId: pacsClrSysMmbId{MmbId: SchemeOperatorID}}}
		doc.FICdtTrf.CdtTrfTx.ForDbtAgt = &pacsAgent{FinInstnId: pacsFinInstnId{ClrSysMmbId: pacsClrSysMmbId{MmbId: inst.Member}}}
	default:
		return nil, fmt.Errorf("pacs.009: unknown direction %q", inst.Direction)
	}

	var buf bytes.Buffer
	buf.WriteString(xml.Header)
	enc, err := xml.MarshalIndent(doc, "", "  ")
	if err != nil {
		return nil, err
	}
	buf.Write(enc)
	buf.WriteString("\n")
	return buf.Bytes(), nil
}
