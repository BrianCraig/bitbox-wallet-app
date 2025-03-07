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

package btc

import (
	"math/big"
	"strconv"

	"github.com/btcsuite/btcd/chaincfg/chainhash"
	"github.com/btcsuite/btcd/txscript"
	"github.com/btcsuite/btcd/wire"
	"github.com/btcsuite/btcutil"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/accounts"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/accounts/errors"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/addresses"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/blockchain"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/maketx"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/btc/transactions"
	"github.com/digitalbitbox/bitbox-wallet-app/backend/coins/coin"
	"github.com/digitalbitbox/bitbox-wallet-app/util/errp"
)

// unitSatoshi is 1 BTC (default unit) in Satoshi.
const unitSatoshi = 1e8

// getFeePerKb returns the fee rate to be used in a new transaction. It is deduced from the supplied
// fee target (priority) if one is given, or the provided args.FeePerKb if the fee taret is
// `FeeTargetCodeCustom`.
func (account *Account) getFeePerKb(args *accounts.TxProposalArgs) (btcutil.Amount, error) {
	if args.FeeTargetCode == accounts.FeeTargetCodeCustom {
		float, err := strconv.ParseFloat(args.CustomFee, 64)
		if err != nil {
			return 0, err
		}
		// Technically it is vKb (virtual Kb) since fees are computed from a transaction's weight
		// (measured in weight units or virtual bytes), but we keep the `Kb` unit to be consistent
		// with the rest of the codebase and Bitcoin Core.
		feePerKb := btcutil.Amount(float * 1000)
		if feePerKb < account.getMinRelayFeeRate() {
			return 0, errors.ErrFeeTooLow
		}
		return feePerKb, nil
	}
	var feeTarget *FeeTarget
	for _, target := range account.feeTargets {
		if target.code == args.FeeTargetCode {
			feeTarget = target
			break
		}
	}
	if feeTarget == nil || feeTarget.feeRatePerKb == nil {
		return 0, errp.New("Fee could not be estimated")
	}
	return *feeTarget.feeRatePerKb, nil
}

// newTx creates a new tx to the given recipient address. It also returns a set of used account
// outputs, which contains all outputs that spent in the tx. Those are needed to be able to sign the
// transaction. selectedUTXOs restricts the available coins; if empty, no restriction is applied and
// all unspent coins can be used.
func (account *Account) newTx(args *accounts.TxProposalArgs) (
	map[wire.OutPoint]*transactions.SpendableOutput, *maketx.TxProposal, error) {

	account.log.Debug("Prepare new transaction")

	address, err := account.coin.DecodeAddress(args.RecipientAddress)
	if err != nil {
		return nil, nil, err
	}
	pkScript, err := txscript.PayToAddrScript(address)
	if err != nil {
		return nil, nil, errp.WithStack(err)
	}
	utxo := account.transactions.SpendableOutputs()
	wireUTXO := make(map[wire.OutPoint]maketx.UTXO, len(utxo))
	for outPoint, txOut := range utxo {
		// Apply coin control.
		if len(args.SelectedUTXOs) != 0 {
			if _, ok := args.SelectedUTXOs[outPoint]; !ok {
				continue
			}
		}
		wireUTXO[outPoint] = maketx.UTXO{
			TxOut: txOut.TxOut,
			Configuration: account.getAddress(
				blockchain.NewScriptHashHex(txOut.TxOut.PkScript)).Configuration,
		}
	}
	feeRatePerKb, err := account.getFeePerKb(args)
	if err != nil {
		return nil, nil, err
	}

	var txProposal *maketx.TxProposal
	if args.Amount.SendAll() {
		txProposal, err = maketx.NewTxSpendAll(
			account.coin,
			wireUTXO,
			pkScript,
			feeRatePerKb,
			account.log,
		)
		if err != nil {
			return nil, nil, err
		}
	} else {
		allowZero := false
		parsedAmount, err := args.Amount.Amount(big.NewInt(unitSatoshi), allowZero)
		if err != nil {
			return nil, nil, err
		}
		parsedAmountInt64, err := parsedAmount.Int64()
		if err != nil {
			return nil, nil, errp.WithStack(errors.ErrInvalidAmount)
		}
		txProposal, err = maketx.NewTx(
			account.coin,
			wireUTXO,
			wire.NewTxOut(parsedAmountInt64, pkScript),
			feeRatePerKb,
			// Change address is of the first subaccount, always.
			account.subaccounts[0].changeAddresses.GetUnused()[0],
			account.log,
		)
		if err != nil {
			return nil, nil, err
		}
	}
	account.log.Debugf("creating tx with %d inputs, %d outputs",
		len(txProposal.Transaction.TxIn), len(txProposal.Transaction.TxOut))
	return utxo, txProposal, nil
}

func (account *Account) getAddress(scriptHashHex blockchain.ScriptHashHex) *addresses.AccountAddress {
	for _, subacc := range account.subaccounts {
		if address := subacc.receiveAddresses.LookupByScriptHashHex(scriptHashHex); address != nil {
			return address
		}
		if address := subacc.changeAddresses.LookupByScriptHashHex(scriptHashHex); address != nil {
			return address
		}
	}
	panic("address must be present")
}

// SendTx implements accounts.Interface.
func (account *Account) SendTx() error {
	unlock := account.activeTxProposalLock.RLock()
	txProposal := account.activeTxProposal
	unlock()
	if txProposal == nil {
		return errp.New("No active tx proposal")
	}

	note := account.BaseAccount.GetAndClearProposedTxNote()

	account.log.Info("Signing and sending transaction")
	utxos := account.transactions.SpendableOutputs()
	getPrevTx := func(txHash chainhash.Hash) *wire.MsgTx {
		txChan := make(chan *wire.MsgTx)
		account.coin.Blockchain().TransactionGet(txHash,
			func(tx *wire.MsgTx) {
				txChan <- tx
			},
			func(err error) {
				if err != nil {
					panic(err)
				}
			},
		)
		return <-txChan
	}
	if err := account.signTransaction(txProposal, utxos, getPrevTx); err != nil {
		return errp.WithMessage(err, "Failed to sign transaction")
	}

	account.log.Info("Signed transaction is broadcasted")
	if err := account.coin.Blockchain().TransactionBroadcast(txProposal.Transaction); err != nil {
		return err
	}
	if err := account.SetTxNote(txProposal.Transaction.TxHash().String(), note); err != nil {
		// Not critical.
		account.log.WithError(err).Error("Failed to save transaction note when sending a tx")
	}
	return nil
}

// TxProposal creates a tx from the relevant input and returns information about it for display in
// the UI (the output amount and the fee). At the same time, it validates the input. The proposal is
// stored internally and can be signed and sent with SendTx().
func (account *Account) TxProposal(
	args *accounts.TxProposalArgs,
) (
	coin.Amount, coin.Amount, coin.Amount, error) {
	defer account.activeTxProposalLock.Lock()()

	account.log.Debug("Proposing transaction")
	_, txProposal, err := account.newTx(args)
	if err != nil {
		return coin.Amount{}, coin.Amount{}, coin.Amount{}, err
	}

	account.activeTxProposal = txProposal

	account.log.WithField("fee", txProposal.Fee).Debug("Returning fee")
	return coin.NewAmountFromInt64(int64(txProposal.Amount)),
		coin.NewAmountFromInt64(int64(txProposal.Fee)),
		coin.NewAmountFromInt64(int64(txProposal.Total())), nil
}
