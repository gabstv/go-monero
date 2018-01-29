package unit

// MonetaryUnitType represents the monetary unit type of Monero
type MonetaryUnitType uint64

// Stolen from https://getmonero.org/resources/moneropedia/denominations.html
const (
	Piconero  MonetaryUnitType = 1
	Nanonero  MonetaryUnitType = 1e3 * Piconero  // 1000
	Micronero MonetaryUnitType = 1e3 * Nanonero  // 1000
	Millinero MonetaryUnitType = 1e3 * Micronero // 1000
	Centinero MonetaryUnitType = 1e1 * Millinero // 10
	Decinero  MonetaryUnitType = 1e1 * Centinero // 10
	Monero    MonetaryUnitType = 1e1 * Decinero  // 10
	Decanero  MonetaryUnitType = 1e1 * Monero    // 10
	Hectonero MonetaryUnitType = 1e1 * Decanero  // 10
	Kilonero  MonetaryUnitType = 1e1 * Hectonero // 10
	Meganero  MonetaryUnitType = 1e3 * Kilonero  // 1000
)
