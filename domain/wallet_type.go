package domain

type WalletType uint8

const (
	WalletTypeCOMMON   WalletType = 1
	WalletTypeMERCHANT WalletType = 2
)

func (wt WalletType) IsValid() bool {
	switch wt {
	case WalletTypeCOMMON, WalletTypeMERCHANT:
		return true
	}
	return false
}
