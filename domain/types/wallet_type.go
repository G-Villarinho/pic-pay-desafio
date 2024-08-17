package types

type WalletType uint8

const (
	WalletTypeCOMMON WalletType = iota
	WalletTypeMERCHANT
)

func (wt WalletType) IsValid() bool {
	switch wt {
	case WalletTypeCOMMON, WalletTypeMERCHANT:
		return true
	}
	return false
}
