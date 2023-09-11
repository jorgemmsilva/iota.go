package iotago

import (
	"context"

	"github.com/iotaledger/hive.go/serializer/v2/serix"
)

func CommonSerixAPI() *serix.API {
	api := serix.NewAPI()

	{
		must(api.RegisterTypeSettings(Ed25519Address{},
			serix.TypeSettings{}.WithObjectType(uint8(AddressEd25519)).WithMapKey("pubKeyHash")),
		)
		must(api.RegisterTypeSettings(RestrictedEd25519Address{},
			serix.TypeSettings{}.WithObjectType(uint8(AddressRestrictedEd25519))),
		)
		must(api.RegisterTypeSettings(AccountAddress{},
			serix.TypeSettings{}.WithObjectType(uint8(AddressAccount)).WithMapKey("accountId")),
		)
		must(api.RegisterTypeSettings(RestrictedAccountAddress{},
			serix.TypeSettings{}.WithObjectType(uint8(AddressRestrictedAccount))),
		)
		must(api.RegisterTypeSettings(NFTAddress{},
			serix.TypeSettings{}.WithObjectType(uint8(AddressNFT)).WithMapKey("nftId")),
		)
		must(api.RegisterTypeSettings(RestrictedNFTAddress{},
			serix.TypeSettings{}.WithObjectType(uint8(AddressRestrictedNFT))),
		)
		must(api.RegisterTypeSettings(ImplicitAccountCreationAddress{},
			serix.TypeSettings{}.WithObjectType(uint8(AddressImplicitAccountCreation)).WithMapKey("pubKeyHash")),
		)
		must(api.RegisterInterfaceObjects((*Address)(nil), (*Ed25519Address)(nil)))
		must(api.RegisterInterfaceObjects((*Address)(nil), (*RestrictedEd25519Address)(nil)))
		must(api.RegisterInterfaceObjects((*Address)(nil), (*AccountAddress)(nil)))
		must(api.RegisterInterfaceObjects((*Address)(nil), (*RestrictedAccountAddress)(nil)))
		must(api.RegisterInterfaceObjects((*Address)(nil), (*NFTAddress)(nil)))
		must(api.RegisterInterfaceObjects((*Address)(nil), (*RestrictedNFTAddress)(nil)))
		must(api.RegisterInterfaceObjects((*Address)(nil), (*ImplicitAccountCreationAddress)(nil)))

		// All versions of the protocol need to be able to parse older protocol parameter versions.
		{
			must(api.RegisterTypeSettings(RentStructure{}, serix.TypeSettings{}))

			must(api.RegisterTypeSettings(V3ProtocolParameters{},
				serix.TypeSettings{}.WithObjectType(uint8(ProtocolParametersV3))),
			)
			must(api.RegisterInterfaceObjects((*ProtocolParameters)(nil), (*V3ProtocolParameters)(nil)))
		}
	}

	return api
}

func ProtocolParametersFromBytes(bytes []byte) (params ProtocolParameters, bytesRead int, err error) {
	var protocolParameters ProtocolParameters
	n, err := CommonSerixAPI().Decode(context.TODO(), bytes, &protocolParameters, serix.WithValidation())
	if err != nil {
		return nil, 0, err
	}

	return protocolParameters, n, nil
}
