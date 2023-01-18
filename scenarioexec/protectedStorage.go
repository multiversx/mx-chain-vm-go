package scenarioexec

// ElrondProtectedKeyPrefix prefixes all Elrond reserved storage. Only the protocol can write to keys starting with this.
const ElrondProtectedKeyPrefix = "ELROND"

// ElrondRewardKey is the storage key where the protocol writes when sending out rewards.
const ElrondRewardKey = ElrondProtectedKeyPrefix + "reward"
