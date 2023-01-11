package vmserver

// UpgradeRequest is a CLI / REST request message
type UpgradeRequest struct {
	DeployRequest
	ContractAddressHex string
	ContractAddress    []byte
}

func (request *UpgradeRequest) digest() error {
	err := request.DeployRequest.digest()
	if err != nil {
		return err
	}

	request.ContractAddress, err = fromHex(request.ContractAddressHex)
	if err != nil {
		return err
	}

	return nil
}

// UpgradeResponse is a CLI / REST response message
type UpgradeResponse struct {
	ContractResponseBase
}
