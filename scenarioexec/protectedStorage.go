package scenarioexec

// ProtectedKeyPrefix prefixes all reserved storage. Only the protocol can write to keys starting with this.
const ProtectedKeyPrefix = "E" + "L" + "R" + "O" + "N" + "D"

// RewardKey is the storage key where the protocol writes when sending out rewards.
const RewardKey = ProtectedKeyPrefix + "reward"
