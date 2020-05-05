package main

import (
	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/go-bip39"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync/atomic"
	"time"
)

const (
	ChainID             = "decimal-testnet"
	RootPath            = "$HOME/.decimal/cli"
	mnemonicEntropySize = 256

	DefaultGas    = uint64(50000)
	DefaultGasAdj = float64(1.1)

	RPCPrefix = "http://localhost:26657"
)

type Account struct {
	Address   sdk.AccAddress
	AccNumber uint64
	Name      string
	Sequence  *uint64
	Password  string
}

type Provider struct {
	cdc     *codec.Codec
	keybase keys.Keybase
}

func NewProvider() *Provider {
	cfg := sdk.GetConfig()
	cfg.SetCoinType(60)
	cfg.SetFullFundraiserPath("44'/60'/0'/0/0")
	cfg.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	cfg.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
	cfg.Seal()

	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	cdc.RegisterConcrete(coin.MsgSendCoin{}, "coin/SendCoin", nil)
	cdc.RegisterConcrete(coin.MsgBuyCoin{}, "coin/BuyCoin", nil)
	cdc.RegisterConcrete(coin.MsgSellCoin{}, "coin/SellCoin", nil)
	codec.RegisterCrypto(cdc)
	cdc.Seal()

	rootPath := os.ExpandEnv(RootPath)

	// Initialize and prepare keybase
	keybase, err := keys.NewKeyring(sdk.KeyringServiceName(), keys.BackendTest, rootPath, nil)
	if err != nil {
		log.Fatalf("ERROR: Unable to initialize keybase: %v", err)
	}

	return &Provider{
		cdc:     cdc,
		keybase: keybase,
	}
}

func (p *Provider) CreateAccount(name, password string) (Account, error) {
	entropySeed, err := bip39.NewEntropy(mnemonicEntropySize)
	if err != nil {
		return Account{}, err
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return Account{}, err
	}

	info, err := p.keybase.CreateAccount(name, mnemonic, password, password, "44'/60'/0'/0/0", keys.Secp256k1)
	if err != nil {
		return Account{}, err
	}

	return Account{
		Address:   info.GetAddress(),
		AccNumber: 6,
		Name:      name,
		Sequence:  new(uint64),
		Password:  password,
	}, nil
}

func main() {
	mainAddrRaw := flag.String("main-account", "", "Address of main account")
	configPath := flag.String("config", "config.json", "Path to config")

	flag.Parse()

	if *mainAddrRaw == "" {
		fmt.Println("error: you must specify the address of the spammer account")
		return
	}

	config, err := ImportConfig(*configPath)
	if err != nil {
		log.Println(err)
		return
	}

	provider := NewProvider()

	accounts := make([]Account, 1)

	accounts[0], err = provider.CreateAccount("spam30", "12345678")
	if err != nil {
		log.Println(err)
		return
	}

	mainAddr, err := sdk.AccAddressFromBech32(*mainAddrRaw)
	if err != nil {
		log.Println(err)
		return
	}

	mainSequence, mainAccNumber, err := GetSequenceAndAccNumber(mainAddr.String())
	if err != nil {
		log.Println(err)
		return
	}

	mainAccount := Account{
		Address:   mainAddr,
		AccNumber: mainAccNumber,
		Name:      "spammer",
		Sequence:  &mainSequence,
		Password:  "12345678",
	}

	err = provider.SendCoin(mainAccount, accounts[0], 1000000)
	if err != nil {
		log.Println(err)
		return
	}

	time.Sleep(time.Second * 5)

	*accounts[0].Sequence, accounts[0].AccNumber, err = GetSequenceAndAccNumber(accounts[0].Address.String())
	if err != nil {
		log.Println(err)
		return
	}

	for i := 0; i < len(accounts); i++ {
		for j := 0; j < 1; j++ {
			go func(account Account) {
				for {
					err = provider.SendCoin(account, account, 5)
					if err != nil {
						log.Println(err)
					}
					time.Sleep(config.TimeoutMs.Send)
				}
			}(accounts[i])

			go func(account Account) {
				for {
					err = provider.BuyCoin("KIR", "tDCL", sdk.NewInt(1), sdk.NewInt(0), account)
					if err != nil {
						log.Println(err)
					}
					time.Sleep(config.TimeoutMs.Buy)
				}
			}(accounts[i])
		}
	}

	select {}
}

func (p *Provider) SendCoin(sender, receiver Account, amount int64) error {
	memo := "spam send"
	txEncoder := auth.DefaultTxEncoder(p.cdc)
	txBldr := auth.NewTxBuilder(
		txEncoder,
		sender.AccNumber, atomic.LoadUint64(sender.Sequence),
		DefaultGas, DefaultGasAdj,
		false, ChainID, memo, nil, nil,
	).WithKeybase(p.keybase)

	msgs := []sdk.Msg{coin.NewMsgSendCoin(sender.Address, "tDCL", sdk.NewInt(amount), receiver.Address)}

	tx, err := txBldr.BuildAndSign(sender.Name, sender.Password, msgs)
	if err != nil {
		return err
	}

	// TODO: Find the way to avoid this ugly hack!
	{
		hackPrefix, _ := hex.DecodeString("282816a9")
		hackLength := (int(tx[1])<<8 + int(tx[0])) + 4
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
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	atomic.AddUint64(sender.Sequence, 1)

	// Read broadcast response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Broadcast response: %s", string(respBody))
	return nil
}

func GetSequenceAndAccNumber(address string) (uint64, uint64, error) {
	resp, err := http.Get("http://localhost:1317/auth/accounts/" + address)
	if err != nil {
		return 0, 0, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	decoder := json.NewDecoder(resp.Body)

	var account struct {
		Result struct {
			Value struct {
				Sequence      uint64 `json:"sequence"`
				AccountNumber uint64 `json:"account_number"`
			} `json:"value"`
		} `json:"result"`
	}
	err = decoder.Decode(&account)
	if err != nil {
		return 0, 0, err
	}

	return account.Result.Value.Sequence, account.Result.Value.AccountNumber, nil
}

func (p *Provider) BuyCoin(coinToBuy, coinToSell string, amountToBuy, amountToSell sdk.Int, buyer Account) error {
	memo := "spam send"
	txEncoder := auth.DefaultTxEncoder(p.cdc)
	txBldr := auth.NewTxBuilder(
		txEncoder,
		buyer.AccNumber, atomic.LoadUint64(buyer.Sequence),
		DefaultGas, DefaultGasAdj,
		false, ChainID, memo, nil, nil,
	).WithKeybase(p.keybase)

	msgs := []sdk.Msg{coin.NewMsgBuyCoin(buyer.Address, coinToBuy, coinToSell, amountToBuy, amountToSell)}

	tx, err := txBldr.BuildAndSign(buyer.Name, buyer.Password, msgs)
	if err != nil {
		return err
	}

	// TODO: Find the way to avoid this ugly hack!
	{
		hackPrefix, _ := hex.DecodeString("282816a9")
		hackLength := (int(tx[1])<<8 + int(tx[0])) + 4
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
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	atomic.AddUint64(buyer.Sequence, 1)

	// Read broadcast response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Broadcast response: %s", string(respBody))
	return nil
}

func (p *Provider) SellCoin(coinToBuy, coinToSell string, amountToBuy, amountToSell sdk.Int, buyer Account) error {
	memo := "spam send"
	txEncoder := auth.DefaultTxEncoder(p.cdc)
	txBldr := auth.NewTxBuilder(
		txEncoder,
		buyer.AccNumber, atomic.LoadUint64(buyer.Sequence),
		DefaultGas, DefaultGasAdj,
		false, ChainID, memo, nil, nil,
	).WithKeybase(p.keybase)

	msgs := []sdk.Msg{coin.NewMsgSellCoin(buyer.Address, coinToBuy, coinToSell, amountToSell, amountToBuy)}

	tx, err := txBldr.BuildAndSign(buyer.Name, buyer.Password, msgs)
	if err != nil {
		return err
	}

	// TODO: Find the way to avoid this ugly hack!
	{
		hackPrefix, _ := hex.DecodeString("282816a9")
		hackLength := (int(tx[1])<<8 + int(tx[0])) + 4
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
		return err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()

	atomic.AddUint64(buyer.Sequence, 1)

	// Read broadcast response
	respBody, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	log.Printf("Broadcast response: %s", string(respBody))
	return nil
}
