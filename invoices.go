package xero

const (
	invoiceTypeACCPAY = "ACCPAY"
	invoiceTypeACCREC = "ACCREC"
)

var (
	ACCPAY = InvoiceType{invoiceTypeACCPAY}
	ACCREC = InvoiceType{invoiceTypeACCREC}
)

// TODO: XML encoding / decoding support
type InvoiceType struct {
	value string
}

// TODO: Payment Dayte Type
// TODO: Line Items
// TODO: Payments
type Invoice struct {
	Type                InvoiceType `xml:"Type,omitempty"`
	Contact             Contact     `xml:"Contact,omitempty"`
	Date                string      `xml:"Date,omitempty"`
	DueDate             string      `xml:"DueDate,omitempty"`
	Status              string      `xml:"Status,omitempty"`
	LineAmountTypes     string      `xml:"LineAmountTypes,omitempty"`
	LineItems           []LineItem  `xml:"LineItems>LineItem,omitempty"`
	SubTotal            string      `xml:"SubTotal,omitempty"`
	TotalTax            string      `xml:"TotalTax,omitempty"`
	Total               string      `xml:"Total,omitempty"`
	TotalDiscount       string      `xml:"TotalDiscount,omitempty"`
	UpdateDateUTC       UTCDate     `xml:"UpdateDateUTC,omitempty"`
	CurrencyCode        string      `xml:"CurrencyCode,omitempty"`
	CurrencyRate        string      `xml:"CurrencyRate,omitempty"`
	InvoiceID           string      `xml:"InvoiceID,omitempty"`
	InvoiceNumber       string      `xml:"InvoiceNumber,omitempty"`
	Reference           string      `xml:"Reference,omitempty"`
	BrandingThemeID     string      `xml:"BrandingThemeID,omitempty"`
	URL                 string      `xml:"Url,omitempty"`
	SendToContact       bool        `xml:"SendToContact,omitempty"`
	ExpectedPaymentDate string      `xml:"ExpectedPaymentDate,omitempty"`
	PlannedPaymentDate  string      `xml:"PlannedPaymentDate,omitempty"`
	HasAttachments      bool        `xml:"HasAttachments,omitempty"`
	Payments            string      `xml:"Payments,omitempty"`
	CreditNotes         string      `xml:"CreditNotes,omitempty"`
	PrePayments         string      `xml:"PrePayments,omitempty"`
	OverPayments        string      `xml:"OverPayments,omitempty"`
	AmountDue           string      `xml:"AmountDue,omitempty"`
	AmountPaid          string      `xml:"AmountPaid,omitempty"`
	FullyPaidDate       string      `xml:"FullyPaidDate,omitempty"`
	AmountCredited      string      `xml:"AmountCredited,omitempty"`
}

type LineItem struct {
	Description  string `xml:"Description,omitempty"`
	Quantity     string `xml:Quantity,omitempty"`
	UnitAmount   string `xml:UnitAmount,omitempty"`
	ItemCode     string `xml:"ItemCode,omitempty"`
	AccountCode  string `xml:"AccountCode,omitempty"`
	LimeItemID   string `xml:"LineItemID,omitempty"`
	TaxType      string `xml:"TaxType,omitempty"`
	LimeAmount   string `xml:"LineAmount,omitempty"`
	Tracking     string `xml:"Tracking,omitempty"`
	DiscountRate string `xml:"DiscountRate,omitempty"`
}
