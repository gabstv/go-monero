package unit

// MonetaryUnitType represents the monetary unit type of Monero
type MonetaryUnitType uint64

// Stolen from https://getmonero.org/resources/moneropedia/denominations.html
const (
	Piconero  MonetaryUnitType = 1
	Nanonero                   = 1e3 * Piconero  // 1000
	Micronero                  = 1e3 * Nanonero  // 1000
	Millinero                  = 1e3 * Micronero // 1000
	Centinero                  = 1e1 * Millinero // 10
	Decinero                   = 1e1 * Centinero // 10
	Monero                     = 1e1 * Decinero  // 10
	Decanero                   = 1e1 * Monero    // 10
	Hectonero                  = 1e1 * Decanero  // 10
	Kilonero                   = 1e1 * Hectonero // 10
	Meganero                   = 1e3 * Kilonero  // 1000
)
