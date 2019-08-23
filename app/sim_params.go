package app

const (
	//alias
	OpWeightMsgAliasUpdate = "op_weight_msg_alias_update"
	// asset
	OpWeightMsgIssueToken           = "op_weight_msg_issue_token"
	OpWeightMsgTransferOwnership    = "op_weight_msg_transfer_ownership"
	OpWeightMsgMintToken            = "op_weight_msg_mint_token"
	OpWeightMsgBurnToken            = "op_weight_msg_burn_token"
	OpWeightMsgForbidToken          = "op_weight_msg_forbid_token"
	OpWeightMsgUnForbidToken        = "op_weight_msg_unforbid_token"
	OpWeightMsgAddTokenWhitelist    = "op_weight_msg_add_token_whitelist"
	OpWeightMsgRemoveTokenWhitelist = "op_weight_msg_remove_token_whitelist"
	OpWeightMsgForbidAddr           = "op_weight_msg_forbid_addr"
	OpWeightMsgUnForbidAddr         = "op_weight_msg_unforbid_addr"
	OpWeightMsgModifyTokenInfo      = "op_weight_msg_modify_token_info"
	// bancorlite
	OpWeightMsgBancorInit   = "op_weight_msg_bancor_init"
	OpWeightMsgBancorTrade  = "op_weight_msg_bancor_trade"
	OpWeightMsgBancorCancel = "op_weight_msg_bancor_cancel"
	// bankx
	OpWeightMsgSetMemoRequired = "op_weight_msg_set_memo_required"
	//comment
	OpWeightCreateNewThread   = "op_weight_create_new_thread"
	OpWeightCreateCommentRefs = "op_weight_create_comment_refs"
	// distrx
	OpWeightMsgDonateToCommunityPool = "op_weight_msg_donate_to_community_pool"
	//market
	OpWeightMsgCreateTradingPair    = "op_weight_msg_create_trading_pair"
	OpWeightMsgCancelTradingPair    = "op_weight_msg_cancel_trading_pair"
	OpWeightMsgModifyPricePrecision = "op_weight_msg_modify_price_precision"
	OpWeightMsgCreateOrder          = "op_weight_msg_create_order"
	OpWeightMsgCancelOrder          = "op_weight_msg_cancel_order"
)
