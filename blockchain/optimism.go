package blockchain

import (
	"context"
	"fmt"
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/rs/zerolog/log"

	"github.com/smartcontractkit/chainlink-testing-framework/utils"
)

// OptimismMultinodeClient represents a multi-node, EVM compatible client for the Optimism network
type OptimismMultinodeClient struct {
	*EthereumMultinodeClient
}

// OptimismClient represents a single node, EVM compatible client for the Optimism network
type OptimismClient struct {
	*EthereumClient
}

// Fund sends some OP to an address using the default wallet
func (m *OptimismClient) Fund(toAddress string, amount *big.Float) error {
	privateKey, err := crypto.HexToECDSA(m.DefaultWallet.PrivateKey())
	to := common.HexToAddress(toAddress)
	if err != nil {
		return fmt.Errorf("invalid private key: %v", err)
	}
	// Optimism uses legacy transactions and gas estimations
	suggestedGasPrice, err := m.Client.SuggestGasPrice(context.Background())
	if err != nil {
		return err
	}

	nonce, err := m.GetNonce(context.Background(), common.HexToAddress(m.DefaultWallet.Address()))
	if err != nil {
		return err
	}

	tx, err := types.SignNewTx(privateKey, types.LatestSignerForChainID(m.GetChainID()), &types.LegacyTx{
		Nonce:    nonce,
		To:       &to,
		Value:    utils.EtherToWei(amount),
		GasPrice: suggestedGasPrice,
		Gas:      21000,
	})
	if err != nil {
		return err
	}

	log.Info().
		Str("Token", "OP").
		Str("From", m.DefaultWallet.Address()).
		Str("To", toAddress).
		Str("Amount", amount.String()).
		Msg("Funding Address")
	if err := m.Client.SendTransaction(context.Background(), tx); err != nil {
		return err
	}

	return m.ProcessTransaction(tx)
}
