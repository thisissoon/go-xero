package xero

import (
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
)

// BankTransactions API Root
const apiBankTransactionsRoot = "/BankTransactions"

// BankTransactionsEndpoint defines the Xero bank transactions endpoint
var BankTransactionsEndpoint = Endpoint(apiBankTransactionsRoot)

// The BankTransaction type represents a single transaction within Xero.
//    <BankTransaction>
//      <Contact>...</Contact>
//      <Date>2010-07-30T00:00:00</Date>
//      <LineAmountTypes>Inclusive</LineAmountTypes>
//      <LineItems>
//        <LineItem>
//          <Description>Monthly account fee</Description>
//          <UnitAmount>15</UnitAmount>
//          <TaxType>NONE</TaxType>
//          <TaxAmount>0.00</TaxAmount>
//          <LineAmount>15.00</LineAmount>
//          <AccountCode>404</AccountCode>
//          <Quantity>1.0000</Quantity>
//          <LineItemID>52208ff9-528a-4985-a9ad-b2b1d4210e38</LineItemID>
//        </LineItem>
//      </LineItems>
//      <SubTotal>15.00</SubTotal>
//      <TotalTax>0.00</TotalTax>
//      <Total>15.00</Total>
//      <UpdatedDateUTC>2008-02-20T12:19:56.657</UpdatedDateUTC>
//      <FullyPaidOnDate>2010-07-30T00:00:00</FullyPaidOnDate>
//      <BankTransactionID>d20b6c54-7f5d-4ce6-ab83-55f609719126</BankTransactionID>
//      <BankAccount>
//        <AccountID>297c2dc5-cc47-4afd-8ec8-74990b8761e9</AccountID>
//        <Code>BANK</Code>
//      </BankAccount>
//      <Type>SPEND</Type>
//      <IsReconciled>true</IsReconciled>
//    </BankTransaction>
type BankTransaction struct {
	ValidationErrors // Used for validating POST/PUT requests

	Type              BankTransactionType   `xml:"Type,omitempty"`
	Contact           Contact               `xml:"Contact,omitempty"`
	LineItems         []LineItem            `xml:"Lineitems>Lineitem,omitempty"`
	BankAccount       BankAccount           `xml:"BankAccount,omitempty"`
	IsReconciled      bool                  `xml:"IsReconciled,omitempty"`
	Date              UTCDate               `xml:"Date,omitempty"`
	Reference         string                `xml:"Reference,omitempty"`
	CurrencyCode      string                `xml:"CurrencyCode,omitempty"`
	CurrencyRate      float32               `xml:"CurrencyRate,omitempty"`
	URL               string                `xml:"Url,omitempty"`
	Status            BankTransactionStatus `xml:"Status,omitempty"`
	LineAmountTypes   LineAmountType        `xml:"LineAmountTypes,omitempty"`
	SubTotal          float64               `xml:"SubTotal,omitempty"`
	TotalTax          float64               `xml:"TotalTax,omitempty"`
	Total             float64               `xml:"Total,omitempty"`
	BankTransactionID string                `xml:"BankTransactionID,omitempty"`
	PrepaymentID      string                `xml:"PrepaymentID,omitempty"`
	OverpaymentID     string                `xml:"OverpaymentID,omitempty"`
	UpdatedDateUTC    UTCDate               `xml:"UpdatedDateUTC,omitempty"`
	HasAttachments    bool                  `xml:"HasAttachments,omitempty"`
}

func (c BankTransaction) Encode(dst io.Writer) error {
	return encode(dst, &c)
}

type BankTransactions struct {
	BankTransactions []BankTransaction `xml:"BankTransactions>BankTransaction"`
}

func (c BankTransactions) Encode(dst io.Writer) error {
	return encode(dst, &c)
}

type BankTransactionsResponse struct {
	Response
	BankTransactions
}

// The BankTransactionIterator type allows for recursive paginated calls
// for n number of pages of contacts in 100 contact batches
type BankTransactionIterator struct {
	page   int
	getter getter
	root   *url.URL
}

// url constructs a url from the root url appending query params
func (c BankTransactionIterator) url() string {
	v := url.Values{}
	v.Set("page", fmt.Sprintf("%d", c.page))
	u := *c.root
	u.RawQuery = v.Encode()
	return u.String()
}

// Next calls the next page of the /BankTransactions endpoint returning the next
// page of transactions. If no transactions are returned we have reached the end
// and an io.EOF error is returned
func (c BankTransactionIterator) Next() (BankTransactionIterator, []BankTransaction, error) {
	var dst BankTransactionsResponse
	if err := c.getter.get(c.url(), &dst); err != nil {
		return c, nil, err
	}
	if len(dst.BankTransactions.BankTransactions) == 0 {
		return c, nil, io.EOF
	}
	c.page++
	return c, dst.BankTransactions.BankTransactions, nil
}

// BankTransaction returns a specific bank transaction from the Xero API
// Identifier can be the Xero identifier for a transaction e.g. 297c2dc5-cc47-4afd-8ec8-74990b8761e9
func (c *Client) BankTransaction(identifier string) (BankTransaction, error) {
	var dst BankTransactionsResponse
	var transaction BankTransaction
	urlStr := c.url(BankTransactionsEndpoint, identifier).String()
	if err := c.get(urlStr, &dst); err != nil {
		return transaction, err
	}
	if len(dst.BankTransactions.BankTransactions) == 0 {
		return transaction, fmt.Errorf("transaction %s not found", identifier)
	}
	transaction = dst.BankTransactions.BankTransactions[0]
	return transaction, nil
}

// The BankTransactions method returns a BankTransactionIterator and first batch of
// BankTransactions from the /BankTransactions endpoint. Call the iterator
// recursivly until the iterator errors with an io.EOF or the length of contacts is 0
func (c *Client) BankTransactions() (BankTransactionIterator, []BankTransaction, error) {
	return BankTransactionIterator{
		page:   1,
		getter: c,
		root:   c.url(BankTransactionsEndpoint), // https://api.xero.com/api.xro/2.0/BankTransactions
	}.Next()
}

// The BankAccount type represents a single bank account in Xero
//   <BankAccount>
//     <AccountID>297c2dc5-cc47-4afd-8ec8-74990b8761e9</AccountID>
//     <Code>BANK</Code>
//   </BankAccount>
type BankAccount struct {
	Code      string `xml:"Code,omitempty"`
	AccountID string `xml:"AccountID,omitempty"`
	Name      string `xml:"Name,omitempty"`
}

// The LineItem type represents a single line item in Xero
//    <LineItem>
//      <Description>Monthly account fee</Description>
//      <UnitAmount>15</UnitAmount>
//      <TaxType>NONE</TaxType>
//      <TaxAmount>0.00</TaxAmount>
//      <LineAmount>15.00</LineAmount>
//      <AccountCode>404</AccountCode>
//      <Quantity>1.0000</Quantity>
//      <LineItemID>52208ff9-528a-4985-a9ad-b2b1d4210e38</LineItemID>
//    </LineItem>
type LineItem struct {
	Description string  `xml:"Description,omitempty"`
	Quantity    float32 `xml:"Quantity,omitempty"`
	UnitAmount  float64 `xml:"UnitAmount,omitempty"`
	AccountCode string  `xml:"AccountCode,omitempty"`
	ItemCode    string  `xml:"ItemCode,omitempty"`
	LineItemID  string  `xml:"LineItemID,omitempty"`
	TaxType     string  `xml:"TaxType,omitempty"` // TODO implement tax types
	LineAmount  float64 `xml:"LineAmount,omitempty"`
	Tracking    string  `xml:"Tracking,omitempty"`
}

// Line Amount Types
// Predefined line amount types from Xero
// https://developer.xero.com/documentation/api/types#LineAmountTypes
const (
	lineAmountTypeExc   = "Exclusive"
	lineAmountTypeInc   = "Inclusive"
	lineAmountTypeNoTax = "NoTax"
)

// Xero Line amount types
var (
	LineAmountTypeExc   = LineAmountType{lineAmountTypeExc}
	LineAmountTypeInc   = LineAmountType{lineAmountTypeInc}
	LineAmountTypeNoTax = LineAmountType{lineAmountTypeNoTax}
)

// The LineAmountType defines specific line amount types within Xero:
type LineAmountType struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the LineAmountType
func (a LineAmountType) String() string {
	return a.value
}

// MarshalXML marshals a LineAmountType into valid XML for Xero
func (a *LineAmountType) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero LineAmountType XML data into valid LineAmountType
func (a *LineAmountType) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}

	switch value {
	case lineAmountTypeExc:
		*a = LineAmountTypeExc
	case lineAmountTypeInc:
		*a = LineAmountTypeInc
	case lineAmountTypeNoTax:
		*a = LineAmountTypeNoTax
	default:
		return fmt.Errorf("unsupported line amount type: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero LineAmountType XML data into valid LineAmountType
func (a *LineAmountType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}

// Bank Transaction Status
// Predefined bank transaction statuses from Xero
// https://developer.xero.com/documentation/api/types#BankTransactionStatuses
const (
	bankTransStatusAuth = "AUTHORISED"
	bankTransStatusDel  = "DELETED"
)

// Xero Bank transaction statuses
var (
	BankTransStatusAuth = BankTransactionStatus{bankTransStatusAuth}
	BankTransStatusDel  = BankTransactionStatus{bankTransStatusDel}
)

// The BankTransactionStatus defines bank transaction statuses within Xero:
type BankTransactionStatus struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the BankTransactionStatus
func (a BankTransactionStatus) String() string {
	return a.value
}

// MarshalXML marshals a BankTransactionStatus into valid XML for Xero
func (a *BankTransactionStatus) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero BankTransactionStatus XML data into valid BankTransactionStatus
func (a *BankTransactionStatus) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}

	switch value {
	case bankTransStatusAuth:
		*a = BankTransStatusAuth
	case bankTransStatusDel:
		*a = BankTransStatusDel
	default:
		return fmt.Errorf("unsupported bank transaction status: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero BankTransactionStatus XML data into valid BankTransactionStatus
func (a *BankTransactionStatus) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}

// Bank Transaction Type
// Predefined bank transaction types from Xero
// https://developer.xero.com/documentation/api/types#BankTransactionTypes
const (
	bankTransTypeReceive = "RECEIVE"
	bankTransTypeROver   = "RECEIVE-OVERPAYMENT"
	bankTransTypeRPrepay = "RECEIVE-PREPAYMENT"
	bankTransTypeSpend   = "SPEND"
	bankTransTypeSOver   = "SPEND-OVERPAYMENT"
	bankTransTypeSPrepay = "SPEND-PREPAYMENT"
	// The following values are only supported via the GET method at the moment
	bankTransTypeRTransfer = "RECEIVE-TRANSFER"
	bankTransTypeSTransfer = "SPEND-TRANSFER"
)

// Xero Bank transaction types
var (
	BankTransTypeReceive   = BankTransactionType{bankTransTypeReceive}
	BankTransTypeROver     = BankTransactionType{bankTransTypeROver}
	BankTransTypeRPrepay   = BankTransactionType{bankTransTypeRPrepay}
	BankTransTypeSpend     = BankTransactionType{bankTransTypeSpend}
	BankTransTypeSOver     = BankTransactionType{bankTransTypeSOver}
	BankTransTypeSPrepay   = BankTransactionType{bankTransTypeSPrepay}
	BankTransTypeRTransfer = BankTransactionType{bankTransTypeRTransfer}
	BankTransTypeSTransfer = BankTransactionType{bankTransTypeSTransfer}
)

// BankTransactionTypes is a slice of Xero bank transaction types
var BankTransactionTypes = []BankTransactionType{
	BankTransTypeReceive,
	BankTransTypeROver,
	BankTransTypeRPrepay,
	BankTransTypeSpend,
	BankTransTypeSOver,
	BankTransTypeSPrepay,
	BankTransTypeRTransfer,
	BankTransTypeSTransfer,
}

// The BankTransactionType defines the specific bank transaction types within Xero:
type BankTransactionType struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the BankTransactionType
func (a BankTransactionType) String() string {
	return a.value
}

// MarshalXML marshals a BankTransactionType into valid XML for Xero
func (a *BankTransactionType) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero BankTransactionType XML data into valid BankTransactionType
func (a *BankTransactionType) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}

	err := fmt.Errorf("unsupported bank transaction type: %s", value)
	for i := 0; i < len(BankTransactionTypes); i++ {
		if value == BankTransactionTypes[i].value {
			*a = BankTransactionTypes[i]
			err = nil
		}
	}

	if err != nil {
		return err
	}
	return nil
}

// UnmarshalXML handles converting raw Xero BankTransactionType XML data into valid BankTransactionType
func (a *BankTransactionType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}
