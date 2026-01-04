package address_helper

import "testing"

func TestCheckAddress_ValidAddresses_ShouldReturnNil(t *testing.T) {
	addresses := []string{
		"0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
		"0x32be343b94f860124dc4fee278fdcbd38c102d88",
		"0xd8dA6BF26964aF9D7eEd9e03E53415D37aA96045",
		"0x000000000000000000000000000000000000dead",
		"0x5077d54024564758525049534575806950275845",
		"0xFE9986E75c407886F6927977D64843940c96D3C9",
		"0x1111111111111111111111111111111111111111",
		"0x90F8bf6A479f320ead074411a4B0e7944Ea8c9C1",
		"0xFFffffFFFFffffFFFFffffFFFFffffFFFFffffFF",
		"0xa5409ec958c83c3f309868babaca7c86dcb077c1",
	}

	for i := range addresses {
		err := CheckAddress(addresses[i])
		if err != nil {
			t.Errorf("%s: expected nil error, got %s", addresses[i], err)
		}
	}
}

func TestCheckAddress_InvalidAddresses_ShouldReturnError(t *testing.T) {
	addresses := []string{
		// Invalid characters (non-hex)
		"0x71C7656EC7ab88b098defB751B7401B5f6d8976G",
		"0x3c44cdddb6a900fa2b585dd299e03d12fa4293zh",
		"0x90f8bf6a479f320ead074411a4b0e7944ea8c9cR",
		"0xgggggggggggggggggggggggggggggggggggggggg",
		"0xd1220a0cf47c7b9be7a2e6ba89f429762e7b9a!!",
		"0xFE9986E75c407886F6927977D64843940c96D3-9",
		"0x71C7656EC7ab88b098defB751B7401B5f6d8976 ",

		// Invalid length
		"0x71C7656EC7ab88b098defB751B7401B5f6d897",
		"0x3c44cdddb6a900fa2b585dd299e03d12fa4293bc1",
		"0x123",
		"0x",
		"0xFE9986E75c407886F6927977D64843940c96D3C9abcde",

		// Invalid prefix
		"71C7656EC7ab88b098defB751B7401B5f6d8976F",
		"1x71C7656EC7ab88b098defB751B7401B5f6d8976F",
		"x071C7656EC7ab88b098defB751B7401B5f6d8976F",
		"#71C7656EC7ab88b098defB751B7401B5f6d8976F",

		// Other formatting errors
		" 0x71C7656EC7ab88b098defB751B7401B5f6d8976F",
		"0x71C7656EC 7ab88b098defB751B7401B5f6d8976F",
		"test_address_12345678901234567890123456",
		"0x........................................",
	}

	for i := range addresses {
		err := CheckAddress(addresses[i])
		if err == nil {
			t.Errorf("%s: expected error, got nil", addresses[i])
		}
	}
}
