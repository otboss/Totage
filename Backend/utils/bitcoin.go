package utils

import (
	"fmt"

	"github.com/btcsuite/btcd/btcutil"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/tyler-smith/go-bip39"
)

type Bitcoin struct {
}

// GenerateSeedPhrase generates a new 12-word BIP39 mnemonic seed phrase
// compatible with modern Bitcoin and Electrum wallets.
func (b Bitcoin) GenerateSeedPhrase() (string, error) {
	// 128 bits of entropy produces a 12-word seed phrase (use 256 for 24 words)
	entropy, err := bip39.NewEntropy(128)
	if err != nil {
		return "", fmt.Errorf("failed to generate entropy: %w", err)
	}

	mnemonic, err := bip39.NewMnemonic(entropy)
	if err != nil {
		return "", fmt.Errorf("failed to generate mnemonic: %w", err)
	}

	return mnemonic, nil
}

// DeriveKeysFromSeedPhrase accepts a seed phrase, validates it, and outputs
// the WIF private key and compressed hex public key.
func (b Bitcoin) DeriveKeysFromSeedPhrase(mnemonic string) (privateKeyWIF string, publicKeyHex string, err error) {
	// Validate the input seed phrase
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", "", fmt.Errorf("invalid seed phrase mnemonic")
	}

	// Generate a 512-bit seed from the mnemonic phrase
	seed := bip39.NewSeed(mnemonic, "")

	// Create a hierarchical deterministic (HD) master key using MainNet configuration
	// (Change to chaincfg.TestNet3Params if working with testnet)
	masterKey, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", fmt.Errorf("failed to create master key: %w", err)
	}

	// Derive a child key (e.g., m/0' - hardened derivation at index 0)
	childKey, err := masterKey.Derive(hdkeychain.HardenedKeyStart + 0)
	if err != nil {
		return "", "", fmt.Errorf("failed to derive child key: %w", err)
	}

	// Extract the private key
	privKey, err := childKey.ECPrivKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to extract private key: %w", err)
	}

	// Encode private key into Wallet Import Format (WIF)
	wif, err := btcutil.NewWIF(privKey, &chaincfg.MainNetParams, true)
	if err != nil {
		return "", "", fmt.Errorf("failed to create WIF: %w", err)
	}

	// Extract the public key and serialize it in compressed format (33 bytes)
	pubKey, err := childKey.ECPubKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to extract public key: %w", err)
	}
	publicKeyHex = fmt.Sprintf("%x", pubKey.SerializeCompressed())

	return wif.String(), publicKeyHex, nil
}
