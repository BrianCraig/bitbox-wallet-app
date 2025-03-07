// Copyright 2018 Shift Devices AG
// Copyright 2020 Shift Crypto AG
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package keystore

import (
	"errors"

	"github.com/btcsuite/btcutil/hdkeychain"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/coin"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/signing"
)

// Type denotes the type of a keystore.
type Type string

const (
	// TypeHardware means the keystore is provided by a hardware wallet.
	TypeHardware Type = "hardware"
	// TypeSoftware mans the keystore is provided by a software (hot) wallet. Currently only used in
	// devmode for testing.
	TypeSoftware Type = "software"
)

// ErrSigningAborted is used when the user aborts a signing in process (e.g. abort on HW wallet).
var ErrSigningAborted = errors.New("signing aborted by user")

// Keystore supports hardened key derivation according to BIP32 and signing of transactions.
//go:generate moq -pkg mocks -out mocks/keystore.go . Keystore
type Keystore interface {
	// Type denotes the type of the keystore.
	Type() Type

	// RootFingerprint returns the keystore's root fingerprint, which is the first 32 bits of the
	// hash160 of the pubkey at the keypath m/.
	//
	// https://github.com/bitcoin/bips/blob/master/bip-0032.mediawiki#key-identifiers
	RootFingerprint() ([]byte, error)

	// SupportsCoin returns true if the keystore supports at least one account type for this coin.
	SupportsCoin(coinInstance coin.Coin) bool

	// SupportsAccount returns true if they keystore supports the given coin/account.
	// meta is a coin-specific metadata related to the account type.
	SupportsAccount(coinInstance coin.Coin, meta interface{}) bool

	// SupportsUnifiedAccounts returns true if the keystore supports signing transactions with mixed
	// input script types in BTC/LTC, for single-sig accounts.
	// If false, the backend will add one account per supported script type.
	SupportsUnifiedAccounts() bool

	// SupportsMultipleAccounts returns true if the keystore can handle more than one account per
	// coin.
	SupportsMultipleAccounts() bool

	// CanVerifyAddress returns whether the keystore supports to output an address securely.
	// This is typically done through a screen on the device or through a paired mobile phone.
	// optional is true if the user can skip verification, and false if they should be forced to
	// verify.
	CanVerifyAddress(coin.Coin) (secureOutput bool, optional bool, err error)

	// VerifyAddress outputs the public key at the given configuration for the given coin.
	// Please note that this is only supported if the keystore has a secure output channel.
	VerifyAddress(*signing.Configuration, coin.Coin) error

	// CanVerifyExtendedPublicKey returns whether the keystore supports to output an xpub/zpub/tbup/ypub securely.
	CanVerifyExtendedPublicKey() bool

	// VerifyExtendedPublicKey displays the public key on the device for verification
	VerifyExtendedPublicKey(coin.Coin, *signing.Configuration) error

	// ExtendedPublicKey returns the extended public key at the given absolute keypath.
	ExtendedPublicKey(coin.Coin, signing.AbsoluteKeypath) (*hdkeychain.ExtendedKey, error)

	// CanSignMessage returns true if the keystore can sign a message for a coin.
	CanSignMessage(coin.Code) bool
	// SignBTCMessage signs the message using the private key at the keypath. The scriptType is
	// required to compute and verify the address. The returned signature is a 65 byte signature in
	// Electrum format.
	SignBTCMessage(message []byte, keypath signing.AbsoluteKeypath, scriptType signing.ScriptType) ([]byte, error)
	// SignETHMessage signs the message using thep rivate key at the keypath. The result contains a
	// 65 byte signature. The first 64 bytes are the secp256k1 signature in / compact format (R and
	// S values), and the last byte is the recoverable id (recid).
	SignETHMessage(message []byte, keypath signing.AbsoluteKeypath) ([]byte, error)

	// SignTransaction signs the given transaction proposal. Returns ErrSigningAborted if the user
	// aborts.
	SignTransaction(interface{}) error
}
