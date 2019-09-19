package types

import (
	"encoding/base64"
	"strings"
	"unicode/utf8"

	sdk "github.com/cosmos/cosmos-sdk/types"

	"github.com/coinexchain/shorthanzi"
)

// RouterKey is the name of the market module
const (
	TokenCommentKey = "token_comment"

	HanziLengthForLZ4 = 512
	MaxContentSize    = 16 * 1024
	MaxTitleSize      = 256
)

// /////////////////////////////////////////////////////////

var _ sdk.Msg = MsgCommentToken{}

const (
	IPFS          int8 = 0
	Magnet        int8 = 1
	HTTP          int8 = 2
	UTF8Text      int8 = 3
	ShortHanzi    int8 = 4
	ShortHanziLZ4 int8 = 5
	RawBytes      int8 = 6

	Like        int32 = 50
	Dislike     int32 = 51
	Laugh       int32 = 52
	Cry         int32 = 53
	Angry       int32 = 54
	Surprise    int32 = 55
	Heart       int32 = 56
	Sweat       int32 = 57
	Speechless  int32 = 58
	Favorite    int32 = 59
	Condolences int32 = 60
)

func ParseContentType(t string) int8 {
	switch strings.ToLower(t) {
	case "ipfs":
		return IPFS
	case "magnet":
		return Magnet
	case "http":
		return HTTP
	case "utf8text":
		return UTF8Text
	case "shorthanzi":
		return ShortHanzi
	case "rawbytes":
		return RawBytes
	case "":
		return UTF8Text
	default:
		return -1
	}
}

func ParseAttitude(a string) int32 {
	switch strings.ToLower(a) {
	case "like":
		return Like
	case "dislike":
		return Dislike
	case "laugh":
		return Laugh
	case "cry":
		return Cry
	case "angry":
		return Angry
	case "surprise":
		return Surprise
	case "heart":
		return Heart
	case "sweat":
		return Sweat
	case "speechless":
		return Speechless
	case "favorite":
		return Favorite
	case "condolences":
		return Condolences
	default:
		return -1
	}
}

type CommentRef struct {
	ID           uint64         `json:"id"`
	RewardTarget sdk.AccAddress `json:"reward_target"`
	RewardToken  string         `json:"reward_token"`
	RewardAmount int64          `json:"reward_amount"`
	Attitudes    []int32        `json:"attitudes"`
}

type MsgCommentToken struct {
	Sender      sdk.AccAddress `json:"sender"`
	Token       string         `json:"token"`
	Donation    int64          `json:"donation"`
	Title       string         `json:"title"`
	Content     []byte         `json:"content"`
	ContentType int8           `json:"content_type"`
	References  []CommentRef   `json:"references"`
}

type TokenComment struct {
	ID          uint64         `json:"id"`
	Height      int64          `json:"height"`
	Sender      sdk.AccAddress `json:"sender"`
	Token       string         `json:"token"`
	Donation    int64          `json:"donation"`
	Title       string         `json:"title"`
	Content     string         `json:"content"`
	ContentType int8           `json:"content_type"`
	References  []CommentRef   `json:"references"`
}

func NewMsgCommentToken(
	sender sdk.AccAddress,
	token string,
	donation int64,
	title string,
	contentStr string,
	contentType int8,
	references []CommentRef) *MsgCommentToken {

	content := []byte(contentStr)
	if contentType == ShortHanzi {
		content = []byte(shorthanzi.Transform(contentStr))
	} else if contentType == RawBytes {
		content, _ = base64.StdEncoding.DecodeString(contentStr)
	}

	if contentType == ShortHanzi && len(contentStr) > HanziLengthForLZ4 {
		contentLZ4, ok := shorthanzi.EncodeHanzi(contentStr)
		if ok {
			contentType = ShortHanziLZ4
			content = contentLZ4
		}
	}
	return &MsgCommentToken{
		Sender:      sender,
		Token:       token,
		Donation:    donation,
		Title:       title,
		Content:     content,
		ContentType: contentType,
		References:  references,
	}
}

func NewTokenComment(msg *MsgCommentToken, id uint64, height int64) *TokenComment {
	tokenComment := &TokenComment{
		ID:          id,
		Height:      height,
		Sender:      msg.Sender,
		Token:       msg.Token,
		Donation:    msg.Donation,
		Title:       msg.Title,
		ContentType: msg.ContentType,
		References:  msg.References,
	}

	if msg.ContentType == RawBytes {
		tokenComment.Content = base64.StdEncoding.EncodeToString(msg.Content)
	} else if msg.ContentType == ShortHanzi {
		tokenComment.ContentType = UTF8Text
		tokenComment.Content = shorthanzi.Transform(string(msg.Content))
	} else if msg.ContentType == ShortHanziLZ4 {
		tokenComment.ContentType = UTF8Text
		tokenComment.Content, _ = shorthanzi.DecodeHanzi(msg.Content)
	} else {
		tokenComment.Content = string(msg.Content)
	}
	return tokenComment
}

// --------------------------------------------------------
// sdk.Msg Implementation

func (msg MsgCommentToken) Route() string { return RouterKey }

func (msg MsgCommentToken) Type() string { return "comment_token" }

func (msg MsgCommentToken) ValidateBasic() sdk.Error {
	if len(msg.Sender) == 0 {
		return sdk.ErrInvalidAddress("missing sender address")
	}
	if len(msg.Token) == 0 {
		return ErrInvalidSymbol()
	}
	if msg.Donation < 0 {
		return ErrNegativeDonation()
	}
	if len(msg.Title) == 0 && len(msg.References) <= 1 {
		return ErrNoTitle()
	}
	if len(msg.Title) > MaxTitleSize {
		return ErrTitleTooLarge()
	}
	if msg.ContentType < IPFS || msg.ContentType > ShortHanziLZ4 {
		return ErrInvalidContentType(msg.ContentType)
	}
	if msg.ContentType == IPFS || msg.ContentType == Magnet || msg.ContentType == HTTP ||
		msg.ContentType == UTF8Text || msg.ContentType == ShortHanzi {
		if !utf8.Valid(msg.Content) {
			return ErrInvalidContent()
		}
	}
	if len(msg.Content) > MaxContentSize {
		return ErrContentTooLarge()
	}
	for _, ref := range msg.References {
		for _, a := range ref.Attitudes {
			if a < Like || a > Condolences {
				return ErrInvalidAttitude(a)
			}
		}
		if ref.RewardAmount < 0 {
			return ErrNegativeReward()
		}
	}
	return nil
}

func (msg MsgCommentToken) GetSignBytes() []byte {
	return sdk.MustSortJSON(ModuleCdc.MustMarshalJSON(msg))
}

func (msg MsgCommentToken) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{[]byte(msg.Sender)}
}
