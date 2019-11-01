package types

const (
	AttributeValueCategory = ModuleName

	EventTypeIssueToken           = "issue_token"
	EventTypeTransferOwnership    = "transfer_ownership"
	EventTypeMintToken            = "mint_token"
	EventTypeBurnToken            = "burn_token"
	EventTypeForbidToken          = "forbid_token"
	EventTypeUnForbidToken        = "unforbid_token"
	EventTypeAddTokenWhitelist    = "add_token_whitelist"
	EventTypeRemoveTokenWhitelist = "remove_token_whitelist"
	EventTypeForbidAddr           = "forbid_addr"
	EventTypeUnForbidAddr         = "unforbid_addr"
	EventTypeModifyTokenInfo      = "modify_token_info"

	AttributeKeySymbol        = "symbol"
	AttributeKeyTokenOwner    = "owner"
	AttributeKeyOriginalOwner = "original_owner"
	AttributeKeyAmount        = "amount"
	AttributeKeyAddrList      = "address_list"
	AttributeKeyURL           = "url"
	AttributeKeyDescription   = "description"
	AttributeKeyIdentity      = "identity"
)
