package main

import (
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"net/http"
	"os"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/cosmos/cosmos-sdk/x/auth"

	"bitbucket.org/decimalteam/go-node/config"
)

// Token is Telegram Bot token received from @BotFather.
const Token = "1239473666:AAFJyytEi2_amYXCAZ3ADhSZaIqoW1jkeUQ"

// RootPath is path to CLI directory.
const RootPath = "$HOME/.decimal/cli"

// RPCPrefix is prefix of available Tendermint RPC endpoint.
const RPCPrefix = "http://139.59.133.148/rpc"

// Faucet settings.
const (
	FaucetChainID         = "decimal-testnet"
	FaucetGas             = uint64(200000)
	FaucetGasAdj          = float64(1.1)
	FaucetKeyringBackend  = keys.BackendTest
	FaucetAccountNumber   = 5
	FaucetAccountName     = "faucet"
	FaucetAccountPassword = "12345678"
	FaucetAddress         = "dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd"
	FaucetMnemonic        = "repair furnace west loud peasant false six hockey poem tube now alien service phone hazard winter favorite away sand fuel describe version tragic vendor"
)

var _ sdk.Msg = &MsgSendCoin{}
var cdc = codec.New()
var keybase keys.Keybase
var rootPath = os.ExpandEnv(RootPath)
var sequencePath = fmt.Sprintf("%s/%s.sequence", rootPath, FaucetAddress)
var sequence = uint64(0)

func init() {

	// Initialize cosmos-sdk configuration
	cfg := sdk.GetConfig()
	cfg.SetCoinType(60)
	cfg.SetFullFundraiserPath("44'/60'/0'/0/0")
	cfg.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	cfg.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
	cfg.Seal()

	// Initialize cosmos-sdk codec
	sdk.RegisterCodec(cdc)
	cdc.RegisterConcrete(MsgSendCoin{}, "coin/SendCoin", nil)
	codec.RegisterCrypto(cdc)
	cdc.Seal()

	// Initialize and prepare keybase
	var err error
	keybase, err = keys.NewKeyring(sdk.KeyringServiceName(), FaucetKeyringBackend, rootPath, nil)
	if err != nil {
		log.Fatalf("ERROR: Unable to initialize keybase: %v", err)
	}
	keybase.CreateAccount(FaucetAccountName, FaucetMnemonic, FaucetAccountPassword, FaucetAccountPassword, "44'/60'/0'/0/0", keys.Secp256k1)

	// Prepare file containing last used sequence for the faucet addres
	data, err := ioutil.ReadFile(sequencePath)
	if err != nil {
		err = ioutil.WriteFile(sequencePath, []byte(strconv.FormatUint(sequence, 10)), os.ModePerm)
		if err != nil {
			log.Printf("ERROR: Cannot write to file %s: %v", sequencePath, err)
		}
	}
	sequence, err = strconv.ParseUint(strings.Trim(string(data), " \n\r"), 10, 64)
	if err != nil {
		log.Printf("ERROR: Unable to parse sequence from file %s: %v", sequencePath, err)
	}
}

func main() {

	// Create and configure bot
	bot, err := tgbotapi.NewBotAPI(Token)
	if err != nil {
		log.Panic(err)
	}
	bot.Debug = true

	log.Printf("Authorized on account %s", bot.Self.UserName)

	// Create and configure bot updater
	updater := tgbotapi.NewUpdate(0)
	updater.Timeout = 60
	updates, err := bot.GetUpdatesChan(updater)
	if err != nil {
		log.Panic(err)
	}

	// Listen updates from the bot updater
	for update := range updates {

		// Ignore any non-message updates
		if update.Message == nil {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Parse input message as "address amount coin"
		strs := strings.Split(update.Message.Text, " ")
		if len(strs) != 3 {
			text := "Invalid faucet request: it should be in format \"address amount coin\""
			answerWithError(bot, update.Message, text)
			continue
		}

		// Validate address
		address := strs[0]
		if !strings.HasPrefix(address, "dx") {
			text := "Invalid address: it should be prefixed with \"dx\""
			answerWithError(bot, update.Message, text)
			continue
		}

		// Validate coin symbol
		amount, ok := big.NewInt(0).SetString(strs[1], 10)
		if !ok {
			text := "Invalid amount: it should be parseable to integer but it is not"
			answerWithError(bot, update.Message, text)
			continue
		}
		if amount.Sign() <= 0 {
			text := "Invalid amount: it should be greater than 0"
			answerWithError(bot, update.Message, text)
			continue
		}

		// Validate amount to transfer
		if coin := strings.ToLower(strs[2]); coin != "tdcl" {
			text := "Invalid coin symbol: only \"tDCL\" is allowed now"
			answerWithError(bot, update.Message, text)
			continue
		}

		// Send coins by preparing, signing and broadcasting a transaction to the Decimal blockchain
		response, txHash, err := sendCoins(address, amount)
		if err != nil {
			text := fmt.Sprintf("Unable to broadcast a transaction:\n%v", err)
			answerWithError(bot, update.Message, text)
			continue
		}

		// Respond with the broadcast response
		text := fmt.Sprintf("Response:\n```\n%s\n```\nTransaction: %s/tx?hash=0x%s", response, RPCPrefix, txHash)
		answerWithSuccess(bot, update.Message, text)
	}
}

func answerWithSuccess(bot *tgbotapi.BotAPI, m *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ReplyToMessageID = m.MessageID
	bot.Send(msg)
}

func answerWithError(bot *tgbotapi.BotAPI, m *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("ERROR: %s", text))
	msg.ReplyToMessageID = m.MessageID
	bot.Send(msg)
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

func sendCoins(address string, amount *big.Int) (response string, txHash string, err error) {

	memo := "faucet transfer"
	txEncoder := auth.DefaultTxEncoder(cdc)
	txBldr := auth.NewTxBuilder(
		txEncoder,
		FaucetAccountNumber, sequence,
		FaucetGas, FaucetGasAdj,
		false, FaucetChainID, memo, nil, nil,
	).WithKeybase(keybase)

	sender, err := sdk.AccAddressFromBech32(FaucetAddress)
	if err != nil {
		return
	}
	receiver, err := sdk.AccAddressFromBech32(address)
	if err != nil {
		return
	}
	msgs := []sdk.Msg{&MsgSendCoin{
		Sender:   sender,
		Coin:     "tDCL",
		Amount:   sdk.NewIntFromBigInt(amount),
		Receiver: receiver,
	}}

	tx, err := txBldr.BuildAndSign(FaucetAccountName, FaucetAccountPassword, msgs)
	if err != nil {
		return
	}

	// TODO: Find the way to avoid this ugly hack!
	{
		hackPrefix, _ := hex.DecodeString("282816a9")
		hackLength := (int(tx[1])<<8 + int(tx[0])) + 4
		fmt.Println(hackLength)
		hackTx := []byte{byte(hackLength & 0xFF), byte(hackLength >> 8)}
		hackTx = append(hackTx, hackPrefix...)
		hackTx = append(hackTx, tx[2:]...)
		tx = hackTx
	}

	// Broadcast signed transaction
	broadcastURL := fmt.Sprintf("%s/broadcast_tx_sync?tx=0x%x", RPCPrefix, tx)
	log.Printf("Broadcast request: %s", broadcastURL)
	resp, err := http.Get(broadcastURL)
	if err != nil {
		return
	}

	// Read broadcast response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	response = string(respBody)
	log.Printf("Broadcast response: %s", response)

	// Parse broadcast response
	var result map[string]interface{}
	err = json.NewDecoder(strings.NewReader(response)).Decode(&result)
	if err != nil {
		return
	}
	resultTx, ok := result["result"].(map[string]interface{})
	if ok {
		txHash = resultTx["hash"].(string)
	}

	// Update sequence in the sequence file
	sequence++
	err = ioutil.WriteFile(sequencePath, []byte(strconv.FormatUint(sequence, 10)), os.ModePerm)
	if err != nil {
		log.Printf("ERROR: Cannot write to file %s: %v", sequencePath, err)
	}

	return
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

type MsgSendCoin struct {
	Sender   sdk.AccAddress `json:"sender" yaml:"sender"`
	Coin     string         `json:"coin" yaml:"coin"`
	Amount   sdk.Int        `json:"amount" yaml:"amount"`
	Receiver sdk.AccAddress `json:"receiver" yaml:"receiver"`
}

func NewMsgSendCoin(sender sdk.AccAddress, coin string, amount sdk.Int, receiver sdk.AccAddress) MsgSendCoin {
	return MsgSendCoin{
		Sender:   sender,
		Coin:     coin,
		Amount:   amount,
		Receiver: receiver,
	}
}

func (msg MsgSendCoin) Route() string { return "coin" }
func (msg MsgSendCoin) Type() string  { return "SendCoin" }
func (msg MsgSendCoin) GetSigners() []sdk.AccAddress {
	return []sdk.AccAddress{msg.Sender}
}

func (msg MsgSendCoin) GetSignBytes() []byte {
	bz := cdc.MustMarshalJSON(msg)
	return sdk.MustSortJSON(bz)
}

func (msg MsgSendCoin) ValidateBasic() error {
	return ValidateSendCoin(msg)
}

func ValidateSendCoin(msg MsgSendCoin) error {
	if msg.Amount.LTE(sdk.NewInt(0)) {
		return sdkerrors.New("coin", 1, "Amount should be greater than 0")
	}
	return nil
}
