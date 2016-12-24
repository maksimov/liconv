package main

// Component represents a dependency record
type Component struct {
	Name               string
	ArtifactID         string
	GroupID            string
	Version            string
	Description        string
	Use                string
	Type               string
	Language           string
	SiteURL            string
	Provider           string
	ContainsEncryption bool
	CopyrightNotice    string
	AllRightsReserved  bool
	PurchaseDate       string
	PurchaseBy         string

	LicenseName              string
	LicenseVersion           string
	LicenseType              string
	LicenseSiteURL           string
	LicenseText              string
	LicenseMisc              string
	LicenseIncludedOnWebsite bool
}
