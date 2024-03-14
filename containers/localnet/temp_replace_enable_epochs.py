import os
import re
import subprocess
import sys

num_shards = sys.argv[1]
staking_v4_step3_max_num_nodes = 64 - (int(num_shards) + 1) * 2


enable_epochs_content = f"""
    [EnableEpochs]
    # SCDeployEnableEpoch represents the epoch when the deployment of smart contracts will be enabled
    SCDeployEnableEpoch = 1

    # BuiltInFunctionsEnableEpoch represents the epoch when the built in functions will be enabled
    BuiltInFunctionsEnableEpoch = 1

    # RelayedTransactionsEnableEpoch represents the epoch when the relayed transactions will be enabled
    RelayedTransactionsEnableEpoch = 1

    # PenalizedTooMuchGasEnableEpoch represents the epoch when the penalization for using too much gas will be enabled
    PenalizedTooMuchGasEnableEpoch = 0

    # SwitchJailWaitingEnableEpoch represents the epoch when the system smart contract processing at end of epoch is enabled
    SwitchJailWaitingEnableEpoch = 0

    # BelowSignedThresholdEnableEpoch represents the epoch when the change for computing rating for validators below signed rating is enabled
    BelowSignedThresholdEnableEpoch = 0

    # SwitchHysteresisForMinNodesEnableEpoch represents the epoch when the system smart contract changes its config to consider
    # also (minimum) hysteresis nodes for the minimum number of nodes
    SwitchHysteresisForMinNodesEnableEpoch = 1

    # TransactionSignedWithTxHashEnableEpoch represents the epoch when the node will also accept transactions that are
    # signed with the hash of transaction
    TransactionSignedWithTxHashEnableEpoch = 1

    # MetaProtectionEnableEpoch represents the epoch when the transactions to the metachain are checked to have enough gas
    MetaProtectionEnableEpoch = 1

    # AheadOfTimeGasUsageEnableEpoch represents the epoch when the cost of smart contract prepare changes from compiler per byte to ahead of time prepare per byte
    AheadOfTimeGasUsageEnableEpoch = 1

    # GasPriceModifierEnableEpoch represents the epoch when the gas price modifier in fee computation is enabled
    GasPriceModifierEnableEpoch = 1

    # RepairCallbackEnableEpoch represents the epoch when the callback repair is activated for scrs
    RepairCallbackEnableEpoch = 1

    # BlockGasAndFeesReCheckEnableEpoch represents the epoch when gas and fees used in each created or processed block are re-checked
    BlockGasAndFeesReCheckEnableEpoch = 1

    # BalanceWaitingListsEnableEpoch represents the epoch when the shard waiting lists are balanced at the start of an epoch
    BalanceWaitingListsEnableEpoch = 1

    # ReturnDataToLastTransferEnableEpoch represents the epoch when returned data is added to last output transfer for callbacks
    ReturnDataToLastTransferEnableEpoch = 1

    # SenderInOutTransferEnableEpoch represents the epoch when the feature of having different senders in output transfer is enabled
    SenderInOutTransferEnableEpoch = 1

    # StakeEnableEpoch represents the epoch when staking is enabled
    StakeEnableEpoch = 0

    # StakingV2EnableEpoch represents the epoch when staking v2 is enabled
    StakingV2EnableEpoch = 1

    # DoubleKeyProtectionEnableEpoch represents the epoch when the double key protection will be enabled
    DoubleKeyProtectionEnableEpoch = 1

    # ESDTEnableEpoch represents the epoch when ESDT is enabled
    ESDTEnableEpoch = 1

    # GovernanceEnableEpoch represents the epoch when governance is enabled
    GovernanceEnableEpoch = 1

    # DelegationManagerEnableEpoch represents the epoch when the delegation manager is enabled
    # epoch should not be 0
    DelegationManagerEnableEpoch = 1

    # DelegationSmartContractEnableEpoch represents the epoch when delegation smart contract is enabled
    # epoch should not be 0
    DelegationSmartContractEnableEpoch = 1

    # CorrectLastUnjailedEnableEpoch represents the epoch when the fix regaring the last unjailed node should apply
    CorrectLastUnjailedEnableEpoch = 1

    # RelayedTransactionsV2EnableEpoch represents the epoch when the relayed transactions V2 will be enabled
    RelayedTransactionsV2EnableEpoch = 1

    # UnbondTokensV2EnableEpoch represents the epoch when the new implementation of the unbond tokens function is available
    UnbondTokensV2EnableEpoch = 1

    # SaveJailedAlwaysEnableEpoch represents the epoch when saving jailed status at end of epoch will happen in all cases
    SaveJailedAlwaysEnableEpoch = 1

    # ReDelegateBelowMinCheckEnableEpoch represents the epoch when the check for the re-delegated value will be enabled
    ReDelegateBelowMinCheckEnableEpoch = 1

    # ValidatorToDelegationEnableEpoch represents the epoch when the validator-to-delegation feature will be enabled
    ValidatorToDelegationEnableEpoch = 1

    # IncrementSCRNonceInMultiTransferEnableEpoch represents the epoch when the fix for preventing the generation of the same SCRs
    # is enabled. The fix is done by adding an extra increment.
    IncrementSCRNonceInMultiTransferEnableEpoch = 1

    # ESDTMultiTransferEnableEpoch represents the epoch when esdt multitransfer built in function is enabled
    ESDTMultiTransferEnableEpoch = 1

    # GlobalMintBurnDisableEpoch represents the epoch when the global mint and burn functions are disabled
    GlobalMintBurnDisableEpoch = 1

    # ESDTTransferRoleEnableEpoch represents the epoch when esdt transfer role set is enabled
    ESDTTransferRoleEnableEpoch = 1

    # ComputeRewardCheckpointEnableEpoch represents the epoch when compute rewards checkpoint epoch is enabled
    ComputeRewardCheckpointEnableEpoch = 1

    # SCRSizeInvariantCheckEnableEpoch represents the epoch when the scr size invariant check is enabled
    SCRSizeInvariantCheckEnableEpoch = 1

    # BackwardCompSaveKeyValueEnableEpoch represents the epoch when the backward compatibility for save key value error is enabled
    BackwardCompSaveKeyValueEnableEpoch = 1

    # ESDTNFTCreateOnMultiShardEnableEpoch represents the epoch when esdt nft creation is enabled on multiple shards
    ESDTNFTCreateOnMultiShardEnableEpoch = 1

    # MetaESDTSetEnableEpoch represents the epoch when the backward compatibility for save key value error is enabled
    MetaESDTSetEnableEpoch = 1

    # AddTokensToDelegationEnableEpoch represents the epoch when adding tokens to delegation is enabled for whitelisted address
    AddTokensToDelegationEnableEpoch = 1

    # MultiESDTTransferFixOnCallBackOnEnableEpoch represents the epoch when multi esdt transfer on callback fix is enabled
    MultiESDTTransferFixOnCallBackOnEnableEpoch = 1

    # OptimizeGasUsedInCrossMiniBlocksEnableEpoch represents the epoch when gas used in cross shard mini blocks will be optimized
    OptimizeGasUsedInCrossMiniBlocksEnableEpoch = 1

    # CorrectFirstQueuedEpoch represents the epoch when the backward compatibility for setting the first queued node is enabled
    CorrectFirstQueuedEpoch = 1

    # DeleteDelegatorAfterClaimRewardsEnableEpoch represents the epoch when the delegators data is deleted for delegators that have to claim rewards after they withdraw all funds
    DeleteDelegatorAfterClaimRewardsEnableEpoch = 1

    # FixOOGReturnCodeEnableEpoch represents the epoch when the backward compatibility returning out of gas error is enabled
    FixOOGReturnCodeEnableEpoch = 1

    # RemoveNonUpdatedStorageEnableEpoch represents the epoch when the backward compatibility for removing non updated storage is enabled
    RemoveNonUpdatedStorageEnableEpoch = 1

    # OptimizeNFTStoreEnableEpoch represents the epoch when optimizations on NFT metadata store and send are enabled
    OptimizeNFTStoreEnableEpoch = 1

    # CreateNFTThroughExecByCallerEnableEpoch represents the epoch when nft creation through execution on destination by caller is enabled
    CreateNFTThroughExecByCallerEnableEpoch = 1

    # StopDecreasingValidatorRatingWhenStuckEnableEpoch represents the epoch when we should stop decreasing validator's rating if, for instance, a shard gets stuck
    StopDecreasingValidatorRatingWhenStuckEnableEpoch = 1

    # FrontRunningProtectionEnableEpoch represents the epoch when the first version of protection against front running is enabled
    FrontRunningProtectionEnableEpoch = 1

    # IsPayableBySCEnableEpoch represents the epoch when a new flag isPayable by SC is enabled
    IsPayableBySCEnableEpoch = 1

    # CleanUpInformativeSCRsEnableEpoch represents the epoch when the informative-only scrs are cleaned from miniblocks and logs are created from them
    CleanUpInformativeSCRsEnableEpoch = 1

    # StorageAPICostOptimizationEnableEpoch represents the epoch when new storage helper functions are enabled and cost is reduced in Wasm VM
    StorageAPICostOptimizationEnableEpoch = 1

    # TransformToMultiShardCreateEnableEpoch represents the epoch when the new function on esdt system sc is enabled to transfer create role into multishard
    TransformToMultiShardCreateEnableEpoch = 1

    # ESDTRegisterAndSetAllRolesEnableEpoch represents the epoch when new function to register tickerID and set all roles is enabled
    ESDTRegisterAndSetAllRolesEnableEpoch = 1

    # ScheduledMiniBlocksEnableEpoch represents the epoch when scheduled mini blocks would be created if needed
    ScheduledMiniBlocksEnableEpoch = 1

    # CorrectJailedNotUnstakedEpoch represents the epoch when the jailed validators will also be unstaked if the queue is empty
    CorrectJailedNotUnstakedEmptyQueueEpoch = 1

    # DoNotReturnOldBlockInBlockchainHookEnableEpoch represents the epoch when the fetch old block operation is
    # disabled in the blockchain hook component
    DoNotReturnOldBlockInBlockchainHookEnableEpoch = 1

    # AddFailedRelayedTxToInvalidMBsDisableEpoch represents the epoch when adding the failed relayed txs to invalid miniblocks is disabled
    AddFailedRelayedTxToInvalidMBsDisableEpoch = 1

    # SCRSizeInvariantOnBuiltInResultEnableEpoch represents the epoch when scr size invariant on built in result is enabled
    SCRSizeInvariantOnBuiltInResultEnableEpoch = 1

    # CheckCorrectTokenIDForTransferRoleEnableEpoch represents the epoch when the correct token ID check is applied for transfer role verification
    CheckCorrectTokenIDForTransferRoleEnableEpoch = 1

    # DisableExecByCallerEnableEpoch represents the epoch when the check on value is disabled on exec by caller
    DisableExecByCallerEnableEpoch = 1

    # RefactorContextEnableEpoch represents the epoch when refactoring/simplifying is enabled in contexts
    RefactorContextEnableEpoch = 1

    # FailExecutionOnEveryAPIErrorEnableEpoch represent the epoch when new protection in VM is enabled to fail all wrong API calls
    FailExecutionOnEveryAPIErrorEnableEpoch = 1

    # ManagedCryptoAPIsEnableEpoch represents the epoch when new managed crypto APIs are enabled in the wasm VM
    ManagedCryptoAPIsEnableEpoch = 1

    # CheckFunctionArgumentEnableEpoch represents the epoch when the extra argument check is enabled in vm-common
    CheckFunctionArgumentEnableEpoch = 1

    # CheckExecuteOnReadOnlyEnableEpoch represents the epoch when the extra checks are enabled for execution on read only
    CheckExecuteOnReadOnlyEnableEpoch = 1

    # ESDTMetadataContinuousCleanupEnableEpoch represents the epoch when esdt metadata is automatically deleted according to inshard liquidity
    ESDTMetadataContinuousCleanupEnableEpoch = 1

    # MiniBlockPartialExecutionEnableEpoch represents the epoch when mini block partial execution will be enabled
    MiniBlockPartialExecutionEnableEpoch = 1

    # FixAsyncCallBackArgsListEnableEpoch represents the epoch when the async callback arguments lists fix will be enabled
    FixAsyncCallBackArgsListEnableEpoch = 1

    # FixOldTokenLiquidityEnableEpoch represents the epoch when the fix for old token liquidity is enabled
    FixOldTokenLiquidityEnableEpoch = 1

    # RuntimeMemStoreLimitEnableEpoch represents the epoch when the condition for Runtime MemStore is enabled
    RuntimeMemStoreLimitEnableEpoch = 1

    # SetSenderInEeiOutputTransferEnableEpoch represents the epoch when setting the sender in eei output transfers will be enabled
    SetSenderInEeiOutputTransferEnableEpoch = 1

    # RefactorPeersMiniBlocksEnableEpoch represents the epoch when refactor of the peers mini blocks will be enabled
    RefactorPeersMiniBlocksEnableEpoch = 1

    # MaxBlockchainHookCountersEnableEpoch represents the epoch when the max blockchainhook counters are enabled
    MaxBlockchainHookCountersEnableEpoch = 1

    # WipeSingleNFTLiquidityDecreaseEnableEpoch represents the epoch when the system account liquidity is decreased for wipeSingleNFT as well
    WipeSingleNFTLiquidityDecreaseEnableEpoch = 1

    # AlwaysSaveTokenMetaDataEnableEpoch represents the epoch when the token metadata is always saved
    AlwaysSaveTokenMetaDataEnableEpoch = 1

    # RuntimeCodeSizeFixEnableEpoch represents the epoch when the code size fix in the VM is enabled
    RuntimeCodeSizeFixEnableEpoch = 1

    # RelayedNonceFixEnableEpoch represents the epoch when the nonce fix for relayed txs is enabled
    RelayedNonceFixEnableEpoch = 1

    # SetGuardianEnableEpoch represents the epoch when the guard account feature is enabled in the protocol
    SetGuardianEnableEpoch = 1

    # DeterministicSortOnValidatorsInfoEnableEpoch represents the epoch when the deterministic sorting on validators info is enabled
    DeterministicSortOnValidatorsInfoEnableEpoch = 1

    # SCProcessorV2EnableEpoch represents the epoch when SC processor V2 will be used
    SCProcessorV2EnableEpoch = 1

    # AutoBalanceDataTriesEnableEpoch represents the epoch when the data tries are automatically balanced by inserting at the hashed key instead of the normal key
    AutoBalanceDataTriesEnableEpoch = 1

    # MigrateDataTrieEnableEpoch represents the epoch when the data tries migration is enabled
    MigrateDataTrieEnableEpoch = 1

    # KeepExecOrderOnCreatedSCRsEnableEpoch represents the epoch when the execution order of created SCRs is ensured
    KeepExecOrderOnCreatedSCRsEnableEpoch = 1

    # MultiClaimOnDelegationEnableEpoch represents the epoch when the multi claim on delegation is enabled
    MultiClaimOnDelegationEnableEpoch = 1

    # ChangeUsernameEnableEpoch represents the epoch when changing username is enabled
    ChangeUsernameEnableEpoch = 4

    # ConsistentTokensValuesLengthCheckEnableEpoch represents the epoch when the consistent tokens values length check is enabled
    ConsistentTokensValuesLengthCheckEnableEpoch = 1

    # FixDelegationChangeOwnerOnAccountEnableEpoch represents the epoch when the fix for the delegation system smart contract is enabled
    FixDelegationChangeOwnerOnAccountEnableEpoch = 1

    # DynamicGasCostForDataTrieStorageLoadEnableEpoch represents the epoch when dynamic gas cost for data trie storage load will be enabled
    DynamicGasCostForDataTrieStorageLoadEnableEpoch = 1

    # ScToScLogEventEnableEpoch represents the epoch when the sc to sc log event feature is enabled
    ScToScLogEventEnableEpoch = 1

    # NFTStopCreateEnableEpoch represents the epoch when NFT stop create feature is enabled
    NFTStopCreateEnableEpoch = 1

    # ChangeOwnerAddressCrossShardThroughSCEnableEpoch represents the epoch when the change owner address built in function will work also through a smart contract call cross shard
    ChangeOwnerAddressCrossShardThroughSCEnableEpoch = 1

    # FixGasRemainingForSaveKeyValueBuiltinFunctionEnableEpoch represents the epoch when the fix for the remaining gas in the SaveKeyValue builtin function is enabled
    FixGasRemainingForSaveKeyValueBuiltinFunctionEnableEpoch = 1

    # CurrentRandomnessOnSortingEnableEpoch represents the epoch when the current randomness on sorting is enabled
    CurrentRandomnessOnSortingEnableEpoch = 4

    # BLSMultiSignerEnableEpoch represents the activation epoch for different types of BLS multi-signers
    BLSMultiSignerEnableEpoch = [
        {{ EnableEpoch = 0, Type = "no-KOSK" }},
        {{ EnableEpoch = 1, Type = "KOSK" }}
    ]

    # StakeLimitsEnableEpoch represents the epoch when stake limits on validators are enabled
    StakeLimitsEnableEpoch = 5

    # StakingV4Step1EnableEpoch represents the epoch when staking v4 is initialized. This is the epoch in which
    # all nodes from staking queue are moved in the auction list
    StakingV4Step1EnableEpoch = 4

    # StakingV4Step2EnableEpoch represents the epoch when staking v4 is enabled. Should have a greater value than StakingV4Step1EnableEpoch.
    # From this epoch, all shuffled out nodes are moved to auction nodes. No auction nodes selection is done yet.
    StakingV4Step2EnableEpoch = 5

    # StakingV4Step3EnableEpoch represents the epoch in which selected nodes from auction will be distributed to waiting list
    StakingV4Step3EnableEpoch = 6

    # MaxNodesChangeEnableEpoch holds configuration for changing the maximum number of nodes and the enabling epoch
    MaxNodesChangeEnableEpoch = [
        {{ EpochEnable = 0, MaxNumNodes = 48, NodesToShufflePerShard = 4 }},  # 4 shuffled out keys / shard will not be reached normally
        {{ EpochEnable = 1, MaxNumNodes = 64, NodesToShufflePerShard = 2 }},
        # Staking v4 configuration, where:
        # - Enable epoch = StakingV4Step3EnableEpoch
        # - NodesToShufflePerShard = same as previous entry in MaxNodesChangeEnableEpoch
        # - MaxNumNodes = (MaxNumNodesFromPreviousEpochEnable - (numOfShards+1)*NodesToShufflePerShard)
        {{ EpochEnable = 6, MaxNumNodes = {staking_v4_step3_max_num_nodes}, NodesToShufflePerShard = 2 }},
    ]

    [GasSchedule]
    # GasScheduleByEpochs holds the configuration for the gas schedule that will be applied from specific epochs
    GasScheduleByEpochs = [
        {{ StartEpoch = 0, FileName = "gasScheduleV7.toml" }},
    ]
"""

def temp_replace_enable_epochs():
    cwd = os.getcwd()

    validator_dirs = [dir_name for dir_name in os.listdir(cwd) if re.match(r'validator\d\d', dir_name)]
    for validator_dir in validator_dirs:
        enable_epochs_path = os.path.join(cwd, validator_dir, 'config', 'enableEpochs.toml')
        with open(enable_epochs_path, 'w') as file:
            file.write(enable_epochs_content)

temp_replace_enable_epochs()
