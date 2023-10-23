package apimodels

// BlockIssuerInfo is the response to the BlockIssuerAPIRouteInfo endpoint.
type BlockIssuerInfo struct {
	// The account address of the block issuer.
	BlockIssuerAddress string `serix:"0,mapKey=blockIssuerAddress"`
	// The number of trailing zeroes required for the proof of work to be valid.
	PowTargetTrailingZeros uint8 `serix:"1,mapKey=powTargetTrailingZeros"`
}