package scenarioexec

// ProtectedKeyPrefix prefixes all Elrond reserved storage. Only the protocol can write to keys starting with this.
const ProtectedKeyPrefix = "E"+"L"+"R"+"O"+"N"+"D"

// ElrondRewardKey is the storage key where the protocol writes when sending out rewards.
const ElrondRewardKey = ProtectedKeyPrefix + "reward"
