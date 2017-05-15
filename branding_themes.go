package xero

// The BrandingTheme holds data regarding specific branding themes created in Xero
//   <BrandingThemes>
//      <BrandingTheme>
//         <BrandingThemeID>a94a78db-5cc6-4e26-a52b-045237e56e6e</BrandingThemeID>
//         <Name>Standard</Name>
//         <SortOrder>0</SortOrder>
//         <CreatedDateUTC>2010-06-29T18:16:36.27</CreatedDateUTC>
//      </BrandingTheme>
//      <BrandingTheme>
//         <BrandingThemeID>db5db9cd-b12e-4faf-8bdd-8eca8af46224</BrandingThemeID>
//         <Name>Special Projects</Name>
//         <SortOrder>1</SortOrder>
//         <CreatedDateUTC>2000-01-01T00:00:00</CreatedDateUTC>
//      </BrandingTheme>
//   </BrandingThemes>
type BrandingTheme struct {
	BrandingThemeID string  `xml:"BrandingThemeID,omitempty"`
	Name            string  `xml:"Name,omitempty"`
	SortOrder       string  `xml:"SortOrder,omitempty"`
	CreatedDateUTC  UTCDate `xml:"CreatedDateUTC,omitempty"`
}
