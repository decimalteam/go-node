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
	"bitbucket.org/decimalteam/go-node/utils/formulas"
)

// Token is Telegram Bot token received from @BotFather.
const Token = "1239473666:AAFJyytEi2_amYXCAZ3ADhSZaIqoW1jkeUQ"

// RootPath is path to CLI directory.
const RootPath = "$HOME/.decimal/cli"

// RPCPrefix is prefix of available Tendermint RPC endpoint.
const RPCPrefix = "http://139.59.133.148/rpc"

// BaseCoin is base coin in original case.
const BaseCoin = "tDCL"

// BaseCoinLower is base coin in lower case.
const BaseCoinLower = "tdcl"

// Faucet settings.
const (
	FaucetChainID         = "decimal-testnet"
	FaucetGas             = uint64(200000)
	FaucetGasAdj          = float64(1.1)
	FaucetKeyringBackend  = keys.BackendTest
	FaucetAccountNumber   = 6
	FaucetAccountName     = "faucet"
	FaucetAccountPassword = "12345678"
	FaucetAddress         = "dx12k95ukkqzjhkm9d94866r4d9fwx7tsd82r8pjd"
	FaucetMnemonic        = "repair furnace west loud peasant false six hockey poem tube now alien service phone hazard winter favorite away sand fuel describe version tragic vendor"
	FaucetBIP44Path       = "44'/60'/0'/0/0"
)

// Coin is a type containing symbol, supply, reserve and CRR values of the coin.
type Coin struct {
	Symbol  string
	Supply  float64
	Reserve float64
	CRR     uint8
}

// Telegram related variables.
var bot *tgbotapi.BotAPI

// Cosmos SDK related variables.
var _ sdk.Msg = &MsgSendCoin{}
var cdc = codec.New()
var keybase keys.Keybase
var rootPath = os.ExpandEnv(RootPath)
var sequencePath = fmt.Sprintf("%s/%s.sequence", rootPath, FaucetAddress)
var sequence = uint64(0)

var err error

func init() {

	// Initialize cosmos-sdk configuration
	cfg := sdk.GetConfig()
	cfg.SetCoinType(60)
	cfg.SetFullFundraiserPath(FaucetBIP44Path)
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
	keybase, err = keys.NewKeyring(sdk.KeyringServiceName(), FaucetKeyringBackend, rootPath, nil)
	if err != nil {
		log.Fatalf("ERROR: Unable to initialize keybase: %v", err)
	}
	keybase.CreateAccount(FaucetAccountName, FaucetMnemonic, "", FaucetAccountPassword, FaucetBIP44Path, keys.Secp256k1)

	// Prepare file containing last used sequence for the faucet address
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
	bot, err = tgbotapi.NewBotAPI(Token)
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

		// Ignore any non-message updates and /start messages
		if update.Message == nil || update.Message.Text == "/start" {
			continue
		}

		log.Printf("[%s] %s", update.Message.From.UserName, update.Message.Text)

		// Try to handle request as trade calculation request
		if handleTradeCalculationRequest(update.Message) {
			continue
		}

		// Handle request as faucet request
		handleFaucetRequest(update.Message)
	}
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

func handleTradeCalculationRequest(m *tgbotapi.Message) (handled bool) {

	coinSpecification := `
> [COIN] supply=S reserve=R crr=CRR
where:
S - amount of COIN supplied at the moment (float)
R - amount of total tDCL reserved for the COIN at the moment (float)
CRR - amount of persentages from 10 to 100 (integer)`

	tradeCalculationRequests := `
> calc buy X COIN1 for COIN2
> calc buy COIN1 for Y COIN2
> calc sell Y COIN1 for COIN2
> calc sell COIN1 for X COIN2
where:
X - amoun of COIN want to receive (float)
Y - amoun of COIN want to spend (float)`

	e18 := big.NewFloat(1000000000000000000)
	floatToInt := func(amount float64) sdk.Int {
		bigFloat := big.NewFloat(0).Mul(big.NewFloat(amount), e18)
		bigInt := big.NewInt(0)
		bigInt, _ = bigFloat.Int(bigInt)
		return sdk.NewIntFromBigInt(bigInt)
	}
	floatFromInt := func(amount sdk.Int) float64 {
		bigFloat := big.NewFloat(0)
		bigFloat.SetInt(amount.BigInt())
		bigFloat = bigFloat.Quo(bigFloat, e18)
		float, _ := bigFloat.Float64()
		return float
	}

	// Parse input message as set of strings
	strs := strings.Split(m.Text, "\n\r")
	strsc := len(strs)
	if strsc <= 1 {
		text := fmt.Sprintf("Invalid trade calculation request: at least one coin should be specified. Usage:%s", coinSpecification)
		answerWithError(m, text)
		return
	}

	// Check if last string is trande calculation request
	calcStr := strings.TrimSpace(strs[strsc-1])
	if !strings.HasPrefix(calcStr, "calc") {
		return
	}

	handled = true

	// Parse trade calculation request string as set of strings and check everything
	calcStrs := strings.Split(calcStr, " ")
	calcStrsc := len(calcStrs)
	if calcStrsc != 6 {
		text := fmt.Sprintf("Invalid trade calculation request. Usage:%s", tradeCalculationRequests)
		answerWithError(m, text)
		return
	}
	if calcStrs[1] != "buy" && calcStrs[1] != "sell" {
		text := fmt.Sprintf("Invalid trade calculation request. Usage:%s", tradeCalculationRequests)
		answerWithError(m, text)
		return
	}
	if calcStrs[3] != "for" && calcStrs[4] != "for" {
		text := fmt.Sprintf("Invalid trade calculation request. Usage:%s", tradeCalculationRequests)
		answerWithError(m, text)
		return
	}

	// Parse coins specifications
	var coins = make(map[string]*Coin)
	for i, c := 0, strsc; i < c; i++ {
		coinStr := strings.TrimSpace(strs[i])
		coinStrs := strings.Split(coinStr, " ")
		coinStrsc := len(coinStrs)
		if coinStrsc != 4 {
			text := fmt.Sprintf("Invalid coin specification in trade calculation request. Usage:%s", coinSpecification)
			answerWithError(m, text)
			return
		}
		if !strings.HasPrefix(coinStrs[0], "[") || !strings.HasSuffix(coinStrs[0], "]") ||
			!strings.HasPrefix(coinStrs[1], "supply=") || !strings.HasPrefix(coinStrs[2], "reserve=") ||
			!strings.HasPrefix(coinStrs[3], "crr=") {
			text := fmt.Sprintf("Invalid coin specification in trade calculation request. Usage:%s", coinSpecification)
			answerWithError(m, text)
			return
		}
		symbol := coinStrs[0][1 : len(coinStrs[0])-1]
		supply, err := strconv.ParseFloat(strings.Trim(coinStrs[1][len("supply="):], " "), 64)
		if err != nil {
			text := fmt.Sprintf("Invalid coin specification in trade calculation request: supply should be parseable to float. Usage:%s", coinSpecification)
			answerWithError(m, text)
			return
		}
		reserve, err := strconv.ParseFloat(strings.Trim(coinStrs[2][len("reserve="):], " "), 64)
		if err != nil {
			text := fmt.Sprintf("Invalid coin specification in trade calculation request: reserve should be parseable to float. Usage:%s", coinSpecification)
			answerWithError(m, text)
			return
		}
		crr, err := strconv.ParseUint(strings.Trim(coinStrs[2][len("crr="):], " "), 10, 64)
		if err != nil {
			text := fmt.Sprintf("Invalid coin specification in trade calculation request: CRR should be parseable to integer. Usage:%s", coinSpecification)
			answerWithError(m, text)
			return
		}
		if crr < 10 || crr > 100 {
			text := fmt.Sprintf("Invalid coin specification in trade calculation request: CRR should be in range [10; 100]. Usage:%s", coinSpecification)
			answerWithError(m, text)
			return
		}
		coins[strings.ToLower(symbol)] = &Coin{
			Symbol:  symbol,
			Supply:  supply,
			Reserve: reserve,
			CRR:     uint8(crr),
		}
	}

	// Parse coins and amount
	var symbolA, symbolB string
	var amount float64
	var buy = (calcStrs[1] == "buy")
	var require = (calcStrs[3] == "for")
	if require {
		symbolA = strings.TrimSpace(calcStrs[2])
		symbolB = strings.TrimSpace(calcStrs[5])
		amount, err = strconv.ParseFloat(strings.TrimSpace(calcStrs[4]), 64)
	} else {
		symbolA = strings.TrimSpace(calcStrs[3])
		symbolB = strings.TrimSpace(calcStrs[5])
		amount, err = strconv.ParseFloat(strings.TrimSpace(calcStrs[2]), 64)
	}
	if err != nil {
		text := fmt.Sprintf("Invalid trade calculation request: amount should be parseable to float. Usage:%s", tradeCalculationRequests)
		answerWithError(m, text)
		return
	}

	// Check coins
	coinA := coins[strings.ToLower(symbolA)]
	coinB := coins[strings.ToLower(symbolB)]
	isBaseCoinA := strings.ToLower(symbolA) == BaseCoinLower
	isBaseCoinB := strings.ToLower(symbolB) == BaseCoinLower
	if isBaseCoinA {
		text := fmt.Sprintf("Invalid trade calculation request: trading coin cannot be %s. Usage:%s%s", BaseCoin, coinSpecification, tradeCalculationRequests)
		answerWithError(m, text)
		return
	}
	if coinA == nil {
		text := fmt.Sprintf("Invalid trade calculation request: unknown coin %s. Usage:%s%s", symbolA, coinSpecification, tradeCalculationRequests)
		answerWithError(m, text)
		return
	}
	if !isBaseCoinB && coinB == nil {
		text := fmt.Sprintf("Invalid trade calculation request: unknown coin %s. Usage:%s%s", symbolB, coinSpecification, tradeCalculationRequests)
		answerWithError(m, text)
		return
	}
	if isBaseCoinA && isBaseCoinB {
		text := fmt.Sprintf("Invalid trade calculation request: unknown coin %s. Usage:%s%s", symbolB, coinSpecification, tradeCalculationRequests)
		answerWithError(m, text)
		return
	}

	// Calculate
	s := floatToInt(coinA.Supply)
	r := floatToInt(coinA.Reserve)
	crr := uint(coinA.CRR)
	a := floatToInt(amount)
	if buy {
		if isBaseCoinB {
			if require {
				result := formulas.CalculatePurchaseAmount(s, r, crr, a)
				text := fmt.Sprintf("You will recieve %f %s by spending %f %s", floatFromInt(result), symbolA, amount, BaseCoin)
				answerWithSuccess(m, text)
			} else {
				result := formulas.CalculatePurchaseReturn(s, r, crr, a)
				text := fmt.Sprintf("To buy %f %s you want to spend %f %s", amount, symbolA, floatFromInt(result), BaseCoin)
				answerWithSuccess(m, text)
			}
		} else {
			// TODO
			text := "Trade calculation for custom coin is not yet supported!"
			answerWithError(m, text)
			return
		}
	} else {
		if isBaseCoinB {
			if require {
				result := formulas.CalculateSaleReturn(s, r, crr, a)
				text := fmt.Sprintf("To sell %f %s you want to spend %f %s", amount, symbolA, floatFromInt(result), BaseCoin)
				answerWithSuccess(m, text)
			} else {
				result := formulas.CalculateSaleAmount(s, r, crr, a)
				text := fmt.Sprintf("You will recieve %f %s by spending %f %s", floatFromInt(result), BaseCoin, amount, symbolA)
				answerWithSuccess(m, text)
			}
		} else {
			// TODO: Implement
			text := "Trade calculation for custom coin is not yet supported!"
			answerWithError(m, text)
			return
		}
	}

	return
}

func handleFaucetRequest(m *tgbotapi.Message) {

	// Parse input message as "address amount coin"
	strs := strings.Split(m.Text, " ")
	if len(strs) != 3 {
		text := "Invalid faucet request: it should be in format \"address amount coin\""
		answerWithError(m, text)
		return
	}

	// Validate address
	address := strs[0]
	if !strings.HasPrefix(address, "dx") {
		text := "Invalid address: it should be prefixed with \"dx\""
		answerWithError(m, text)
		return
	}

	// Validate coin symbol
	amount, ok := big.NewInt(0).SetString(strs[1], 10)
	if !ok {
		text := "Invalid amount: it should be parseable to integer"
		answerWithError(m, text)
		return
	}
	if amount.Sign() <= 0 {
		text := "Invalid amount: it should be greater than 0"
		answerWithError(m, text)
		return
	}

	// Validate amount to transfer
	if coin := strings.ToLower(strs[2]); coin != BaseCoinLower {
		text := fmt.Sprintf("Invalid coin symbol: only %q is allowed now", BaseCoin)
		answerWithError(m, text)
		return
	}

	// Send coins by preparing, signing and broadcasting a transaction to the Decimal blockchain
	response, txHash, err := sendCoins(address, amount)
	if err != nil {
		text := fmt.Sprintf("Unable to broadcast a transaction:\n%v", err)
		answerWithError(m, text)
		return
	}

	// Respond with the broadcast response
	text := fmt.Sprintf("Response:\n%s\nTransaction: %s/tx?hash=0x%s", response, RPCPrefix, txHash)
	answerWithSuccess(m, text)

	return
}

func answerWithSuccess(m *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, text)
	msg.ReplyToMessageID = m.MessageID
	bot.Send(msg)
}

func answerWithError(m *tgbotapi.Message, text string) {
	msg := tgbotapi.NewMessage(m.Chat.ID, fmt.Sprintf("ERROR: %s", text))
	msg.ReplyToMessageID = m.MessageID
	bot.Send(msg)
}

////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////
////////////////////////////////////////////////////////////////

func sendCoins(address string, amount *big.Int) (response string, txHash string, err error) {

	// memo := "faucet transfer"
	memo := ""
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
		Coin:     BaseCoin,
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
		if resultTx["code"].(float64) == 0 {
			// Update sequence in the sequence file
			sequence++
			err = ioutil.WriteFile(sequencePath, []byte(strconv.FormatUint(sequence, 10)), os.ModePerm)
			if err != nil {
				log.Printf("ERROR: Cannot write to file %s: %v", sequencePath, err)
			}
		}
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
