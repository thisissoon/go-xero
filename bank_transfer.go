package xero

import (
	"fmt"
	"io"
)

// BankTransfer API Root
const apiBankTransferRoot = "/BankTransfers"

// BankTransfersEndpoint defines the Xero bank transfers endpoint
var BankTransfersEndpoint = Endpoint(apiBankTransferRoot)

// The BankTransfer type represents a single bank transfer within Xero.
//    <BankTransfer>
//      <BankTransferID>d79f3e07-5f11-45e4-9d1a-30be536d0e13</BankTransferID>
//      <CreatedDateUTC>2014-02-25T19:27:15</CreatedDateUTC>
//      <Date>2014-02-26T00:00:00</Date>
//      <FromBankAccount>
//        <AccountID>ac993f75-035b-433c-82e0-7b7a2d40802c</AccountID>
//        <Name>Business Bank Account</Name>
//      </FromBankAccount>
//      <ToBankAccount>
//        <AccountID>ebd06280-af70-4bed-97c6-7451a454ad85</AccountID>
//        <Name>Business Savings Account</Name>
//      </ToBankAccount>
//      <Amount>20.00</Amount>
//      <FromBankTransactionID>b11794bc-775b-4f78-9b28-8f13240082ff</FromBankTransactionID>
//      <ToBankTransactionID>f589fb5e-34b3-4392-8207-4ba5a093eae</ToBankTransactionID>
//    </BankTransfer>
type BankTransfer struct {
	ValidationErrors // Used for validating POST/PUT requests

	// The following can be set on POST/PUT requests
	FromBankAccount BankAccount `xml:"FromBankAccount,omitempty"`
	ToBankAccount   BankAccount `xml:"ToBankAccount,omitempty"`
	Amount          float64     `xml:"Amount,omitempty"`
	Date            UTCDate     `xml:"Date,omitempty"`
	// The following are only retrieved on GET requests
	BankTransferID        string  `xml:"BankTransferID,omitempty"`
	CurrencyRate          float32 `xml:"CurrencyRate,omitempty"`
	FromBankTransactionID string  `xml:"FromBankTransactionID,omitempty"`
	ToBankTransactionID   string  `xml:"ToBankTransactionID,omitempty"`
	HasAttachments        bool    `xml:"HasAttachments,omitempty"`
	CreatedDateUTC        UTCDate `xml:"CreatedDateUTC,omitempty"`
}

func (c BankTransfer) Encode(dst io.Writer) error {
	return encode(dst, &c)
}

type BankTransfers struct {
	BankTransfers []BankTransfer `xml:"BankTransfers>BankTransfer"`
}

func (c BankTransfers) Encode(dst io.Writer) error {
	return encode(dst, &c)
}

type BankTransfersResponse struct {
	Response
	BankTransfers
}

// BankTransfer returns a specific singular banke transfer from the Xero API
// Identifier can be the Xero identifier for a transfer e.g. 297c2dc5-cc47-4afd-8ec8-74990b8761e9
func (c *Client) BankTransfer(identifier string) (BankTransfer, error) {
	var dst BankTransfersResponse
	var transfer BankTransfer
	urlStr := c.url(BankTransfersEndpoint, identifier).String()
	if err := c.get(urlStr, &dst); err != nil {
		return transfer, err
	}
	if len(dst.BankTransfers.BankTransfers) == 0 {
		return transfer, fmt.Errorf("transfer %s not found", identifier)
	}
	transfer = dst.BankTransfers.BankTransfers[0]
	return transfer, nil
}

// BankTransfers returns a list of BankTransfers from the /BankTransfers endpoint
func (c *Client) BankTransfers() ([]BankTransfer, error) {
	var dst BankTransfersResponse
	urlStr := c.url(AccountsEndpoint).String()
	if err := c.get(urlStr, &dst); err != nil {
		return []BankTransfer{}, err
	}
	return dst.BankTransfers.BankTransfers, nil
}
