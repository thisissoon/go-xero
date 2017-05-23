package xero

import (
	"encoding/xml"
	"fmt"
)

// Accounts API Root
const apiAccountsRoot = "/Accounts"

// The Xero Accounts endpoint
var AccountsEndpoint = Endpoint(apiAccountsRoot)

// The Account type represents a single account within Xero.
//    <Account>
//      <AccountID>297c2dc5-cc47-4afd-8ec8-74990b8761e9</AccountID>
//      <Code>200</Code>
//      <Name>BNZ Cheque Account</Name>
//      <Type>BANK</Type>
//      <BankAccountNumber>3809087654321500</BankAccountNumber>
//      <Description>Income from any normal business activity</Description>
//      <BankAccountType>BANK</BankAccountType>
//      <CurrencyCode>NZD</CurrencyCode>
//      <TaxType>NONE</TaxType>
//      <EnablePaymentsToAccount>false</EnablePaymentsToAccount>
//    </Account>
type Account struct {
	ValidationErrors // Used for validating POST/PUT requests

	// The following can be set on POST/PUT requests
	Code                    string          `xml:"Code,omitempty"`
	Name                    string          `xml:"Name,omitempty"`
	Type                    AccountType     `xml:"Type,omitempty"`
	BankAccountNumber       string          `xml:"BankAccountNumber,omitempty"`
	Status                  AccountStatus   `xml:"Status,omitempty"`
	Description             string          `xml:"Description,omitempty"`
	BankAccountType         BankAccountType `xml:"BankAccountType,omitempty"`
	CurrencyCode            string          `xml:"CurrencyCode,omitempty"`
	TaxType                 string          `xml:"TaxType,omitempty"` // TODO: implement tax types
	EnablePaymentsToAccount bool            `xml:"EnablePaymentsToAccount,omitempty"`
	ShowInExpenseClaims     bool            `xml:"ShowInExpenseClaims,omitempty"`
	// The following are only retrieved on GET requests
	AccountID         string       `xml:"AccountID,omitempty"`
	Class             AccountClass `xml:"Class,omitempty"`
	SystemAccount     string       `xml:"SystemAccount,omitempty"`
	ReportingCode     string       `xml:"ReportingCode,omitempty"`
	ReportingCodeName string       `xml:"ReportingCodeName,omitempty"`
	UpdatedDateUTC    UTCDate      `xml:"UpdatedDateUTC,omitempty"`
	HasAttachments    bool         `xml:"HasAttachments,omitempty"`
}

type AccountsResponse struct {
	Response
	Accounts []Account `xml:"Accounts>Account"`
}

// Account returns a specific singular account from the Xero API
// Identifier can be the Xero identifier for an account e.g. 297c2dc5-cc47-4afd-8ec8-74990b8761e9
func (c *Client) Account(identifier string) (Account, error) {
	var dst AccountsResponse
	var account Account
	urlStr := c.url(AccountsEndpoint, identifier).String()
	if err := c.get(urlStr, &dst); err != nil {
		return account, err
	}
	if len(dst.Accounts) == 0 {
		return account, fmt.Errorf("account %s not found", identifier)
	}
	account = dst.Accounts[0]
	return account, nil
}

// Accounts returns a list of Accounts from the /Accounts endpoint
func (c *Client) Accounts() ([]Account, error) {
	var dst AccountsResponse
	urlStr := c.url(AccountsEndpoint).String()
	if err := c.get(urlStr, &dst); err != nil {
		return []Account{}, err
	}
	return dst.Accounts, nil
}

// Account Class Type
// Predefined account class types from Xero
// https://developer.xero.com/documentation/api/types#AccountClassTypes
const (
	accountClassAsset   = "ASSET"
	accountClassEquity  = "EQUITY"
	accountClassExp     = "EXPENSE"
	accountClassLiab    = "LIABILITY"
	accountClassRevenue = "REVENUE"
)

// Xero Account class types
var (
	AccountClassAsset   = AccountClass{accountClassAsset}
	AccountClassEquity  = AccountClass{accountClassEquity}
	AccountClassExp     = AccountClass{accountClassExp}
	AccountClassLiab    = AccountClass{accountClassLiab}
	AccountClassRevenue = AccountClass{accountClassRevenue}
)

// The AccountClass type defines the specific account class types within Xero:
type AccountClass struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the AccountClass
func (a AccountClass) String() string {
	return a.value
}

// MarshalXML marshals a AccountClass into valid XML for Xero
func (a *AccountClass) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero AccountClass XML data into valid AccountClass
func (a *AccountClass) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	switch value {
	case accountClassAsset:
		*a = AccountClassAsset
	case accountClassEquity:
		*a = AccountClassEquity
	case accountClassExp:
		*a = AccountClassExp
	case accountClassLiab:
		*a = AccountClassLiab
	case accountClassRevenue:
		*a = AccountClassRevenue
	default:
		return fmt.Errorf("unsupported account class: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero AccountClass XML data into valid AccountClass
func (a *AccountClass) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}

// Account Type
// Predefined account types from Xero
// https://developer.xero.com/documentation/api/types#AccountTypes
const (
	accountTypeBank       = "BANK"
	accountTypeCurrent    = "CURRENT"
	accountTypeCurrLiab   = "CURRLIAB"
	accountTypeDepreciatn = "DEPRECIATN"
	accountTypeDC         = "DIRECTCOSTS"
	accountTypeEquity     = "EQUITY"
	accountTypeExp        = "EXPENSE"
	accountTypeFixed      = "FIXED"
	accountTypeInventory  = "INVENTORY"
	accountTypeLiab       = "LIABILITY"
	accountTypeNonCurrent = "NONCURRENT"
	accountTypeOther      = "OTHERINCOME"
	accountTypeOverhead   = "OVERHEADS"
	accountTypePrepay     = "PREPAYMENT"
	accountTypeRevenue    = "REVENUE"
	accountTypeSale       = "SALES"
	accountTypeTermLiab   = "TERMLIAB"
	accountTypePAYG       = "PAYGLIABILITY"
	accountTypeSAExp      = "SUPERANNUATIONEXPENSE"
	accountTypeSALiab     = "SUPERANNUATIONLIABILITY"
	accountTypeWageExp    = "WAGESEXPENSE"
	accountTypeWageLiab   = "WAGESPAYABLELIABILITY"
)

// Xero Account types
var (
	AccountTypeBank       = AccountType{accountTypeBank}
	AccountTypeCurrent    = AccountType{accountTypeCurrent}
	AccountTypeCurrLiab   = AccountType{accountTypeCurrLiab}
	AccountTypeDepreciatn = AccountType{accountTypeDepreciatn}
	AccountTypeDC         = AccountType{accountTypeDC}
	AccountTypeEquity     = AccountType{accountTypeEquity}
	AccountTypeExp        = AccountType{accountTypeExp}
	AccountTypeFixed      = AccountType{accountTypeFixed}
	AccountTypeInventory  = AccountType{accountTypeInventory}
	AccountTypeLiab       = AccountType{accountTypeLiab}
	AccountTypeNonCurrent = AccountType{accountTypeNonCurrent}
	AccountTypeOther      = AccountType{accountTypeOther}
	AccountTypeOverhead   = AccountType{accountTypeOverhead}
	AccountTypePrepay     = AccountType{accountTypePrepay}
	AccountTypeRevenue    = AccountType{accountTypeRevenue}
	AccountTypeSale       = AccountType{accountTypeSale}
	AccountTypeTermLiab   = AccountType{accountTypeTermLiab}
	AccountTypePAYG       = AccountType{accountTypePAYG}
	AccountTypeSAExp      = AccountType{accountTypeSAExp}
	AccountTypeSALiab     = AccountType{accountTypeSALiab}
	AccountTypeWageExp    = AccountType{accountTypeWageExp}
	AccountTypeWageLiab   = AccountType{accountTypeWageLiab}
)

// AccountTypes is a slice of all account types
var AccountTypes = []AccountType{
	AccountTypeBank,
	AccountTypeCurrent,
	AccountTypeCurrLiab,
	AccountTypeDepreciatn,
	AccountTypeDC,
	AccountTypeEquity,
	AccountTypeExp,
	AccountTypeFixed,
	AccountTypeInventory,
	AccountTypeLiab,
	AccountTypeNonCurrent,
	AccountTypeOther,
	AccountTypeOverhead,
	AccountTypePrepay,
	AccountTypeRevenue,
	AccountTypeSale,
	AccountTypeTermLiab,
	AccountTypePAYG,
	AccountTypeSAExp,
	AccountTypeSALiab,
	AccountTypeWageExp,
	AccountTypeWageLiab,
}

// The AccountType type defines the specific account types within Xero:
type AccountType struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the AccountType
func (a AccountType) String() string {
	return a.value
}

// MarshalXML marshals a AccountType into valid XML for Xero
func (a *AccountType) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero AccountType XML data into valid AccountType
func (a *AccountType) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}

	err := fmt.Errorf("unsupported account type: %s", value)
	for i := 0; i < len(AccountTypes); i++ {
		if value == AccountTypes[i].value {
			*a = AccountTypes[i]
			err = nil
		}
	}

	if err != nil {
		return err
	}

	return nil
}

// UnmarshalXML handles converting raw Xero AccountType XML data into valid AccountType
func (a *AccountType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}

// Account Status
// Predefined account statuses from Xero
// https://developer.xero.com/documentation/api/types#AccountStatusCodes
const (
	accountStatusActive  = "ACTIVE"
	accountStatusArchive = "ARCHIVED"
)

// Xero Account statuses
var (
	AccountStatusActive  = AccountStatus{accountStatusActive}
	AccountStatusArchive = AccountStatus{accountStatusArchive}
)

// The AccountStatus type defines the specific account statuses within Xero:
type AccountStatus struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the AccountStatus
func (a AccountStatus) String() string {
	return a.value
}

// MarshalXML marshals a AccountStatus into valid XML for Xero
func (a *AccountStatus) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero AccountStatus XML data into valid AccountStatus
func (a *AccountStatus) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	switch value {
	case accountStatusActive:
		*a = AccountStatusActive
	case accountStatusArchive:
		*a = AccountStatusArchive
	default:
		return fmt.Errorf("unsupported account status: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero AccountStatus XML data into valid AccountStatus
func (a *AccountStatus) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}

// ABank Account Types
// Predefined bank account types from Xero
// https://developer.xero.com/documentation/api/types#BankAccountTypes
const (
	bankAccountTypeBank   = "BANK"
	bankAccountTypeCC     = "CREDITCARD"
	bankAccountTypePaypal = "PAYPAL"
)

// Xero Bank account types
var (
	BankAccountTypeBank   = BankAccountType{bankAccountTypeBank}
	BankAccountTypeCC     = BankAccountType{bankAccountTypeCC}
	BankAccountTypePaypal = BankAccountType{bankAccountTypePaypal}
)

// The BankAccountType type defines the specific bank account types within Xero:
type BankAccountType struct {
	value string
}

// String implements the Stringer interface returning the string representation
// of the BankAccountType
func (a BankAccountType) String() string {
	return a.value
}

// MarshalXML marshals a BankAccountType into valid XML for Xero
func (a *BankAccountType) MarshalXML(encoder *xml.Encoder, start xml.StartElement) error {
	return encoder.EncodeElement(a.value, start)
}

// unmarshalXML handles converting raw Xero BankAccountType XML data into valid BankAccountType
func (a *BankAccountType) unmarshalXML(decoder elementDecoder, start xml.StartElement) error {
	var value string
	if err := decoder.DecodeElement(&value, &start); err != nil {
		return err
	}
	switch value {
	case bankAccountTypeBank:
		*a = BankAccountTypeBank
	case bankAccountTypeCC:
		*a = BankAccountTypeCC
	case bankAccountTypePaypal:
		*a = BankAccountTypePaypal
	default:
		return fmt.Errorf("unsupported bank account type: %s", value)
	}
	return nil
}

// UnmarshalXML handles converting raw Xero BankAccountType XML data into valid BankAccountType
func (a *BankAccountType) UnmarshalXML(decoder *xml.Decoder, start xml.StartElement) error {
	return a.unmarshalXML(decoder, start)
}
