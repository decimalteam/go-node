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

	DefaultGas    = uint64(200000)
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
	keybase, err := keys.NewKeyring("Decimal", keys.BackendTest, rootPath, nil)
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
	configPath := flag.String("cfg", "cfg.json", "Path to cfg")

	flag.Parse()

	if *mainAddrRaw == "" {
		fmt.Println("error: you must specify the address of the spammer account")
		return
	}

	cfg, err := ImportConfig(*configPath)
	if err != nil {
		log.Println(err)
		return
	}

	provider := NewProvider()

	accounts := make([]Account, 2)

	accounts[0], err = provider.CreateAccount("spam30", "12345678")
	if err != nil {
		log.Println(err)
		return
	}

	accounts[1], err = provider.CreateAccount("spam40", "12345678")
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
		Name:      "tank",
		Sequence:  &mainSequence,
		Password:  "12345678",
	}

	for i := 0; i < len(accounts); i++ {
		err = provider.SendCoin(mainAccount, accounts[i], 1000000)
		if err != nil {
			log.Println("Init send", err)
			return
		}

		time.Sleep(time.Second * 5)

		var seq uint64
		seq, accounts[i].AccNumber, err = GetSequenceAndAccNumber(accounts[i].Address.String())
		if err != nil {
			log.Println(err)
			return
		}
		atomic.StoreUint64(accounts[i].Sequence, seq)
		for j := 0; j < 1; j++ {
			go func(account Account) {
				for {
					err = provider.SendCoin(accounts[0], accounts[1], 5)
					if err != nil {
						log.Println(err)
					}
					time.Sleep(cfg.TimeoutMs.Send)
					err = provider.SendCoin(accounts[1], accounts[0], 5)
					if err != nil {
						log.Println(err)
					}
					time.Sleep(cfg.TimeoutMs.Send)
				}
			}(accounts[i])

			//go func(account Account) {
			//	for {
			//		err = provider.BuyCoin("TEST1", "TEST2", sdk.NewInt(1), sdk.NewInt(1), account)
			//		if err != nil {
			//			log.Println(err)
			//		}
			//		time.Sleep(cfg.TimeoutMs.Buy)
			//	}
			//}(accounts[i])
			//
			//go func(account Account) {
			//	for {
			//		err = provider.SellCoin("TEST1", "TEST2", sdk.NewInt(1), sdk.NewInt(1), account)
			//		if err != nil {
			//			log.Println(err)
			//		}
			//		time.Sleep(cfg.TimeoutMs.Sell)
			//	}
			//}(accounts[i])
		}
	}

	select {}
}

type BroadcastResponse struct {
	Result struct {
		Code int    `json:"code"`
		Log  string `json:"log"`
		Hash string `json:"hash"`
	}
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
	//log.Printf("Broadcast request: %s", broadcastURL)
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
	//log.Printf("Broadcast response: %s", string(respBody))

	broadcastResp := BroadcastResponse{}
	err = json.Unmarshal(respBody, &broadcastResp)
	if err != nil {
		return err
	}
	if broadcastResp.Result.Code != 0 {
		log.Printf("Broadcast error: code: %d, log: %s", broadcastResp.Result.Code, broadcastResp.Result.Log)
	} else {
		log.Println("Broadcast hash: ", broadcastResp.Result.Hash)
	}

	return nil
}

func GetSequenceAndAccNumber(address string) (uint64, uint64, error) {
	resp, err := http.Get("http://139.59.133.148/rest/auth/accounts/" + address)
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
