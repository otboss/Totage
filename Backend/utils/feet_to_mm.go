package utils

func FeetToMM(feet float64) float64 {
	const conversionFactor = 304.8
	return feet * conversionFactor
}
