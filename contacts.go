package xero

// The ContactTrackingCategory for SalesTrackingCategories and PurchasesTrackingCategories
type ContactTrackingCategory struct {
	TrackingCategoryName string `xml:"TrackingCategoryName,omitempty"`
	TrackingOptionName   string `xml:"TrackingOptionName,omitempty"`
}

// The ContactBatchPayments holds the batch payment details for a contact
type ContactBatchPayments struct {
	BankAccountNumber string `xml:"BankAccountNumber,omitempty"`
	BankAccountName   string `xml:"BankAccountName,omitempty"`
	Details           string `xml:"Details,omitempty"`
}

// The ContactBalance type holds the AccountsReceivable and AccountsPayable
// ContactBalances values
type ContactBalance struct {
	Outstanding string `xml:"Outstanding,omitempty"`
	Overdue     string `xml:"Overdue,omitempty"`
}

// The ContactBalances type is the raw AccountsReceivable(sales invoices) and
// AccountsPayable(bills) outstanding and overdue amounts,
// not converted to base currency
type ContactBalances struct {
	AccountsReceivable ContactBalance `xml:"AccountsReceivable,omitempty"`
	AccountsPayable    ContactBalance `xml:"AccountsPayable,omitempty"`
}

// A ContactPaymentTerm for bills or sales
type ContactPaymentTerm struct {
	Day  string      `xml:"Day,omitempty"`
	Type PaymentTerm `xml:"Type,omitempty"`
}

// ContactPaymentTerms is the default payment terms for the contact broken
// down into bills and sales
type ContactPaymentTerms struct {
	Bills ContactPaymentTerm `xml:"Bills,omitempty"`
	Sales ContactPaymentTerm `xml:"Sales,omitempty"`
}

// The ContactPerson allows a contact to hold multiple contact details
type ContactPerson struct {
	FirstName       string `xml:"FirstName,omitempty"`
	LastName        string `xml:"LastName,omitempty"`
	EmailAddress    string `xml:"EmailAddress,omitempty"`
	IncludeInEmails bool   `xml:"IncludeInEmails,omitempty"`
}

// The Contact type represnets a single contact within Xero.
//   <Contact>
//      <ContactID>bd2270c3-8706-4c11-9cfb-000b551c3f51</ContactID>
//      <ContactStatus>ACTIVE</ContactStatus>
//      <Name>ABC Limited</Name>
//      <FirstName>Andrea</FirstName>
//      <LastName>Dutchess</LastName>
//      <EmailAddress>a.dutchess@abclimited.com</EmailAddress>
//      <SkypeUserName>skype.dutchess@abclimited.com</SkypeUserName>
//      <BankAccountDetails>454611121</BankAccountDetails>
//      <TaxNumber>415465456454</TaxNumber>
//      <AccountsReceivableTaxType>INPUT2</AccountsReceivableTaxType>
//      <AccountsPayableTaxType>OUTPUT2</AccountsPayableTaxType>
//      <Addresses>
//         <Address>
//            <AddressType>POBOX</AddressType>
//            <AddressLine1>P O Box 123</AddressLine1>
//            <City>Wellington</City>
//            <PostalCode>6011</PostalCode>
//            <AttentionTo>Andrea</AttentionTo>
//         </Address>
//         <Address>
//            <AddressType>STREET</AddressType>
//         </Address>
//      </Addresses>
//      <Phones>
//         <Phone>
//            <PhoneType>DEFAULT</PhoneType>
//            <PhoneNumber>1111111</PhoneNumber>
//            <PhoneAreaCode>04</PhoneAreaCode>
//            <PhoneCountryCode>64</PhoneCountryCode>
//         </Phone>
//         <Phone>
//            <PhoneType>FAX</PhoneType>
//         </Phone>
//         <Phone>
//            <PhoneType>MOBILE</PhoneType>
//         </Phone>
//         <Phone>
//            <PhoneType>DDI</PhoneType>
//         </Phone>
//      </Phones>
//      <UpdatedDateUTC>2009-05-14T01:44:26.747</UpdatedDateUTC>
//      <IsSupplier>false</IsSupplier>
//      <IsCustomer>true</IsCustomer>
//      <DefaultCurrency>NZD</DefaultCurrency>
//      <PurchasesDefaultAccountCode>300</PurchasesDefaultAccountCode>
//      <SalesDefaultAccountCode>200</SalesDefaultAccountCode>
//      <SalesTrackingCategories>
//         <SalesTrackingCategory>
//            <TrackingCategoryName>Region</TrackingCategoryName>
//            <TrackingOptionName>Eastside</TrackingOptionName>
//         </SalesTrackingCategory>
//      </SalesTrackingCategories>
//      <PurchasesTrackingCategories>
//         <PurchasesTrackingCategory>
//            <TrackingCategoryName>Region</TrackingCategoryName>
//            <TrackingOptionName>North</TrackingOptionName>
//         </PurchasesTrackingCategory>
//      </PurchasesTrackingCategories>
//      <Balances>
//         <AccountsReceivable>
//            <Outstanding>849.50</Outstanding>
//            <Overdue>910.00</Overdue>
//         </AccountsReceivable>
//         <AccountsPayable>
//            <Outstanding>0.00</Outstanding>
//            <Overdue>0.00</Overdue>
//         </AccountsPayable>
//      </Balances>
//      <BatchPayments>
//         <BankAccountNumber>123456</BankAccountNumber>
//         <BankAccountName>bank acccount</BankAccountName>
//         <Details>details</Details>
//      </BatchPayments>
//      <PaymentTerms>
//         <Bills>
//            <Day>4</Day>
//            <Type>OFFOLLOWINGMONTH</Type>
//         </Bills>
//         <Sales>
//            <Day>2</Day>
//            <Type>OFFOLLOWINGMONTH</Type>
//         </Sales>
//      </PaymentTerms>
//   </Contact>
type Contact struct {
	ValidationErrors // Used for validating POST/PUT requests

	// The following can be set on POST/PUT requests
	ContactID                 string          `xml:"ContactID,omitempty`
	ContactNumber             string          `xml:"ContactNumber,omitempty`
	AccountNumber             string          `xml:"AccountNumber,omitempty"`
	ContactStatusa            string          `xml:"ContactStatus,omitempty"`
	Name                      string          `xml:"Name,omitempty"`
	FirstName                 string          `xml:"FirstName,omitempty"`
	LastName                  string          `xml:"LastName,omitempty"`
	EmailAddress              string          `xml:"EmailAddress,omitempty"`
	SkypeUserName             string          `xml:"SkypeUserName,omitempty"`
	ContactPersons            []ContactPerson `xml:"ContactPersons>ContactPerson,omitempty"`
	BankAccountDetails        string          `xml:"BankAccountDetails,omitempty"`
	TaxNumber                 string          `xml:"TaxNumber,omitempty"`
	AccountsReceivableTaxType string          `xml:"AccountsReceivableTaxType,omitempty"`
	AccountsPayableTaxType    string          `xml:"AccountsPayableTaxType,omitempty"`
	Addresses                 []Address       `xml:"Addresses>Address,omitempty"`
	Phones                    []Phone         `xml:"Phones>Phone,omitempty"`
	IsSupplier                bool            `xml:"IsSupplier,omitempty"`
	IsCustomer                bool            `xml:"IsCustomer,omitempty"`
	DefaultCurrency           string          `xml:"DefaultCurrency,omitempty"`
	UpdatedDateUTC            UTCDate         `xml:"UpdatedDateUTC,omitempty"`
	// The following are only retrieved on GET requests for a single contact or when pagination is used
	XeroNetworkKey              string                    `xml:"XeroNetworkKey,omitempty"`
	SalesDefaultAccountCode     string                    `xml:"SalesDefaultAccountCode,omitempty"`
	PurchasesDefaultAccountCode string                    `xml:"PurchasesDefaultAccountCode,omitempty"`
	SalesTrackingCategories     []ContactTrackingCategory `xml:"SalesTrackingCategories>SalesTrackingCategory,omitempty`
	PurchasesTrackingCategories []ContactTrackingCategory `xml:"PurchasesTrackingCategories>PurchasesTrackingCategory,omitempty`
	PaymentTerms                ContactPaymentTerms       `xml:"PaymentTerms,omitempty"`
	ContactGroups               []ContactGroup            `xml:"ContactGroups>ContactGroup,omitempty"`
	Website                     string                    `xml:"Website,omitempty"`
	BrandingTheme               BrandingTheme             `xml:"BrandingTheme,omitempty"`
	BatchPayments               ContactBatchPayments      `xml:"BatchPayments,omitempty"`
	Discount                    string                    `xml:"Discount,omitempty"`
	Balances                    ContactBalances           `xml:"Balances,omitempty"`
	HasAttachments              bool                      `xml:"HasAttachments,omitempty"`
}

// The ContactGroup type holds data regarding a users contact group(s) within Xero.
//  <ContactGroup>
//      <ContactGroupID>d0c68f1a-e5dd-4a45-aa02-27d8fdbfd562</ContactGroupID>
//      <Name>Preferred Suppliers</Name>
//      <Status>ACTIVE</Status>
//  </ContactGroup>
type ContactGroup struct {
	ContactGroupID string `xml:"ContactGroupID,omitempty"`
	Name           string `xml:"Name,omitempty"`
	Status         string `xml:"Status,omitempty"`
}