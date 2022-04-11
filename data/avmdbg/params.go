package avmdbg

import (
	"errors"
	"fmt"

	"github.com/algorand/go-algorand/config"
	algod "github.com/algorand/go-algorand/daemon/algod/api/client"
	v2 "github.com/algorand/go-algorand/daemon/algod/api/server/v2"
	"github.com/algorand/go-algorand/data/basics"
	"github.com/algorand/go-algorand/data/bookkeeping"
	"github.com/algorand/go-algorand/data/transactions"
	"github.com/algorand/go-algorand/ledger/ledgercore"
	"github.com/algorand/go-algorand/protocol"
	"github.com/algorand/go-algorand/rpcs"
)

type AssetParamsWithCreator struct {
	basics.AssetParams
	Creator basics.Address
}

type AppParamsWithCreator struct {
	basics.AppParams
	Creator basics.Address
}

type DebuggerParams struct {
	InputTxns []transactions.SignedTxn
	// CheckTxnSignatures bool

	AlgodClient algod.RestClient
}

type DebuggerContext struct {
	Params *DebuggerParams

	// on-chain state
	AccountTotals   ledgercore.AccountTotals
	LastBlockHeader bookkeeping.BlockHeader
	Accounts        map[basics.Address]ledgercore.AccountData
	Assets          map[basics.AssetIndex]AssetParamsWithCreator
	AssetHoldings   map[ledgercore.AccountAsset]basics.AssetHolding
	Apps            map[basics.AppIndex]AppParamsWithCreator
	AppLocalStates  map[ledgercore.AccountApp]basics.AppLocalState
}

func MakeContext(params *DebuggerParams) *DebuggerContext {
	return &DebuggerContext{Params: params}
}

func (ctx *DebuggerContext) GatherResources() error {
	if (ctx.Params.AlgodClient == algod.RestClient{}) {
		return errors.New("An algod client must be provided")
	}

	status, err := ctx.Params.AlgodClient.Status()
	if err != nil {
		return err
	}

	rawBlock, err := ctx.Params.AlgodClient.RawBlock(status.LastRound)
	if err != nil {
		return err
	}

	var lastBlock rpcs.EncodedBlockCert
	err = protocol.DecodeMsgp(rawBlock, &lastBlock)
	if err != nil {
		return err
	}
	// free up payset for GC
	lastBlock.Block.Payset = nil

	proto := config.Consensus[lastBlock.Block.CurrentProtocol]

	supply, err := ctx.Params.AlgodClient.LedgerSupply()
	if err != nil {
		return err
	}

	if supply.Round != status.LastRound {
		return fmt.Errorf("Responses are from different rounds. Got round values %d and %d", status.LastRound, supply.Round)
	}

	// TODO: this is not completely accurate
	onlineMoney := basics.MicroAlgos{
		Raw: supply.OnlineMoney,
	}
	offlineMoney := basics.MicroAlgos{
		Raw: supply.TotalMoney - supply.OnlineMoney/2,
	}
	nonparticipatingMoney := basics.MicroAlgos{
		Raw: (supply.OnlineMoney + 1) / 2,
	}
	accountTotals := ledgercore.AccountTotals{
		Online: ledgercore.AlgoCount{
			Money:       onlineMoney,
			RewardUnits: onlineMoney.RewardUnits(proto),
		},
		Offline: ledgercore.AlgoCount{
			Money:       offlineMoney,
			RewardUnits: offlineMoney.RewardUnits(proto),
		},
		NotParticipating: ledgercore.AlgoCount{
			Money:       nonparticipatingMoney,
			RewardUnits: nonparticipatingMoney.RewardUnits(proto),
		},
		RewardsLevel: lastBlock.Block.RewardsLevel,
	}

	accounts := make(map[basics.Address]ledgercore.AccountData)
	assets := make(map[basics.AssetIndex]AssetParamsWithCreator)
	apps := make(map[basics.AppIndex]AppParamsWithCreator)

	for _, stxn := range ctx.Params.InputTxns {
		accountsToLoad := []basics.Address{
			lastBlock.Block.RewardsPool,
			lastBlock.Block.FeeSink,
		}
		accountsToLoad = append(accountsToLoad, referencedAccounts(stxn.Txn)...)
		for _, addr := range accountsToLoad {
			if _, seen := accounts[addr]; seen {
				continue
			}

			restAccount, err := ctx.Params.AlgodClient.AccountInformationV2(addr.String(), false)
			if err != nil {
				return fmt.Errorf("account %s: %w", addr, err)
			}

			if restAccount.Round != status.LastRound {
				return fmt.Errorf("Responses are from different rounds. Got round values %d and %d", status.LastRound, restAccount.Round)
			}

			accountData, err := v2.AccountToAccountData(&restAccount)
			if err != nil {
				return err
			}

			ledgerAccount := ledgercore.ToAccountData(accountData)
			ledgerAccount.TotalAssetParams = restAccount.TotalCreatedAssets
			ledgerAccount.TotalAssets = restAccount.TotalAppsOptedIn
			ledgerAccount.TotalAppParams = restAccount.TotalCreatedApps
			ledgerAccount.TotalAppLocalStates = restAccount.TotalAppsOptedIn

			accounts[addr] = ledgerAccount
		}

		for _, assetIndex := range referencedAssets(stxn.Txn) {
			if _, seen := assets[assetIndex]; seen {
				continue
			}

			restAsset, err := ctx.Params.AlgodClient.AssetInformationV2(uint64(assetIndex))
			if err != nil {
				return fmt.Errorf("asset %d: %w", assetIndex, err)
			}

			// this response type has no Round, can't validate it

			assetParams, err := v2.AssetToAssetParams(&restAsset)
			if err != nil {
				return err
			}

			creator, err := basics.UnmarshalChecksumAddress(restAsset.Params.Creator)
			if err != nil {
				return err
			}

			assets[assetIndex] = AssetParamsWithCreator{
				AssetParams: assetParams,
				Creator:     creator,
			}
		}

		for _, appIndex := range referencedApps(stxn.Txn) {
			if _, seen := apps[appIndex]; seen {
				continue
			}

			restApp, err := ctx.Params.AlgodClient.ApplicationInformation(uint64(appIndex))
			if err != nil {
				return fmt.Errorf("app %d: %w", appIndex, err)
			}

			// this response type has no Round, can't validate it

			appParams, err := v2.ApplicationParamsToAppParams(&restApp.Params)
			if err != nil {
				return err
			}

			creator, err := basics.UnmarshalChecksumAddress(restApp.Params.Creator)
			if err != nil {
				return err
			}

			apps[appIndex] = AppParamsWithCreator{
				AppParams: appParams,
				Creator:   creator,
			}
		}
	}

	assetHoldings := make(map[ledgercore.AccountAsset]basics.AssetHolding)
	appLocalStates := make(map[ledgercore.AccountApp]basics.AppLocalState)

	for addr := range accounts {
		// if addr == lastBlock.Block.RewardsPool || addr == lastBlock.Block.FeeSink {
		// 	continue
		// }

		addrStr := addr.String()

		for assetIndex := range assets {
			restAssetInfo, err := ctx.Params.AlgodClient.AccountAssetInformation(addrStr, uint64(assetIndex))
			if err != nil {
				var httpError algod.HTTPError
				if errors.As(err, &httpError) && httpError.StatusCode == 404 {
					// the account does not hold this asset
					continue
				}
				return fmt.Errorf("account asset info %s, %d: %w", addrStr, assetIndex, err)
			}

			if restAssetInfo.Round != status.LastRound {
				return fmt.Errorf("Responses are from different rounds. Got round values %d and %d", status.LastRound, restAssetInfo.Round)
			}

			if restAssetInfo.AssetHolding != nil {
				assetHoldings[ledgercore.AccountAsset{Address: addr, Asset: assetIndex}] = v2.GeneratedAssetHoldingToAssetHolding(restAssetInfo.AssetHolding)
			}
		}

		for appIndex := range apps {
			restAppInfo, err := ctx.Params.AlgodClient.AccountApplicationInformation(addrStr, uint64(appIndex))
			if err != nil {
				var httpError algod.HTTPError
				if errors.As(err, &httpError) && httpError.StatusCode == 404 {
					// the account does not hold this app
					continue
				}
				return fmt.Errorf("account app info %s, %d: %w", addrStr, appIndex, err)
			}

			if restAppInfo.Round != status.LastRound {
				return fmt.Errorf("Responses are from different rounds. Got round values %d and %d", status.LastRound, restAppInfo.Round)
			}

			if restAppInfo.AppLocalState != nil {
				appLocalState, err := v2.ApplicationLocalStateToAppLocalState(restAppInfo.AppLocalState)
				if err != nil {
					return err
				}
				appLocalStates[ledgercore.AccountApp{Address: addr, App: appIndex}] = appLocalState
			}
		}
	}

	ctx.AccountTotals = accountTotals
	ctx.LastBlockHeader = lastBlock.Block.BlockHeader
	ctx.Accounts = accounts
	ctx.Assets = assets
	ctx.AssetHoldings = assetHoldings
	ctx.Apps = apps
	ctx.AppLocalStates = appLocalStates

	return nil
}

func referencedAccounts(txn transactions.Transaction) []basics.Address {
	referenced := []basics.Address{txn.Sender}

	switch txn.Type {
	case protocol.PaymentTx:
		referenced = append(referenced, txn.Receiver)
		if !txn.CloseRemainderTo.IsZero() {
			referenced = append(referenced, txn.CloseRemainderTo)
		}
	case protocol.AssetTransferTx:
		referenced = append(referenced, txn.AssetReceiver)
		if !txn.AssetCloseTo.IsZero() {
			referenced = append(referenced, txn.AssetCloseTo)
		}
		if !txn.AssetSender.IsZero() {
			referenced = append(referenced, txn.AssetSender)
		}
	case protocol.ApplicationCallTx:
		referenced = append(referenced, txn.Accounts...)

		if txn.ApplicationID != 0 {
			referenced = append(referenced, txn.ApplicationID.Address())
		}

		for _, app := range txn.ForeignApps {
			referenced = append(referenced, app.Address())
		}
	}

	return referenced
}

func referencedAssets(txn transactions.Transaction) []basics.AssetIndex {
	switch txn.Type {
	case protocol.AssetConfigTx:
		if txn.ConfigAsset != 0 {
			return []basics.AssetIndex{txn.ConfigAsset}
		}
	case protocol.AssetFreezeTx:
		return []basics.AssetIndex{txn.FreezeAsset}
	case protocol.AssetTransferTx:
		return []basics.AssetIndex{txn.XferAsset}
	case protocol.ApplicationCallTx:
		return txn.ForeignAssets
	}

	return nil
}

func referencedApps(txn transactions.Transaction) []basics.AppIndex {
	if txn.Type != protocol.ApplicationCallTx {
		return nil
	}

	referenced := []basics.AppIndex{}

	if txn.ApplicationID != 0 {
		referenced = append(referenced, txn.ApplicationID)
	}

	referenced = append(referenced, txn.ForeignApps...)

	return referenced
}

func (ctx *DebuggerContext) LatestBlockHdr() bookkeeping.BlockHeader {
	return ctx.LastBlockHeader
}

func (ctx *DebuggerContext) GetAccount(addr basics.Address) (ledgercore.AccountData, bool) {
	account, exists := ctx.Accounts[addr]
	return account, exists
}

func (ctx *DebuggerContext) GetAsset(assetIndex basics.AssetIndex) (basics.AssetParams, bool) {
	asset, exists := ctx.Assets[assetIndex]
	return asset.AssetParams, exists
}

func (ctx *DebuggerContext) GetAssetCreator(assetIndex basics.AssetIndex) (basics.Address, bool) {
	asset, exists := ctx.Assets[assetIndex]
	return asset.Creator, exists
}

func (ctx *DebuggerContext) GetAssetHolding(assetIndex basics.AssetIndex, addr basics.Address) (basics.AssetHolding, bool) {
	key := ledgercore.AccountAsset{Asset: assetIndex, Address: addr}
	assetHolding, exists := ctx.AssetHoldings[key]
	return assetHolding, exists
}

func (ctx *DebuggerContext) GetApp(appIndex basics.AppIndex) (basics.AppParams, bool) {
	app, exists := ctx.Apps[appIndex]
	return app.AppParams, exists
}

func (ctx *DebuggerContext) GetAppCreator(appIndex basics.AppIndex) (basics.Address, bool) {
	app, exists := ctx.Apps[appIndex]
	return app.Creator, exists
}

func (ctx *DebuggerContext) GetAppLocalState(appIndex basics.AppIndex, addr basics.Address) (basics.AppLocalState, bool) {
	key := ledgercore.AccountApp{App: appIndex, Address: addr}
	appLocalState, exists := ctx.AppLocalStates[key]
	return appLocalState, exists
}

func (ctx *DebuggerContext) LatestTotals() ledgercore.AccountTotals {
	return ctx.AccountTotals
}
