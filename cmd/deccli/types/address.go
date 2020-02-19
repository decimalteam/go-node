package types

const (
	DecimalMainPrefix = "dx"

	// PrefixValidator is the prefix for validator keys
	PrefixValidator = "val"
	// PrefixConsensus is the prefix for consensus keys
	PrefixConsensus = "cons"
	// PrefixPublic is the prefix for public keys
	PrefixPublic = "pub"
	// PrefixOperator is the prefix for operator keys
	PrefixOperator = "oper"

	// DecimalPrefixAccAddr defines the Decimal prefix of an account's address
	DecimalPrefixAccAddr = DecimalMainPrefix
	// DecimalPrefixAccPub defines the Decimal prefix of an account's public key
	DecimalPrefixAccPub = DecimalMainPrefix + PrefixPublic
	// DecimalPrefixValAddr defines the Decimal prefix of a validator's operator address
	DecimalPrefixValAddr = DecimalMainPrefix + PrefixValidator + PrefixOperator
	// DecimalPrefixValPub defines the Decimal prefix of a validator's operator public key
	DecimalPrefixValPub = DecimalMainPrefix + PrefixValidator + PrefixOperator + PrefixPublic
	// DecimalPrefixConsAddr defines the Decimal prefix of a consensus node address
	DecimalPrefixConsAddr = DecimalMainPrefix + PrefixValidator + PrefixConsensus
	// DecimalPrefixConsPub defines the Decimal prefix of a consensus node public key
	DecimalPrefixConsPub = DecimalMainPrefix + PrefixValidator + PrefixConsensus + PrefixPublic
)
