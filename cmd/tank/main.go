package main

import (
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"sync/atomic"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/go-bip39"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/coin"
)

const (
	ChainID             = "decimal-testnet"
	RootPath            = "$HOME/.decimal/cli"
	TankKeyringBackend  = keys.BackendTest
	TankAccountName     = "tank"
	TankAccountPassword = "12345678"
	TankAddress         = "dx1esffyu0wxk6eez77fhzdxfgvjp4646hqm9sx6c"
	TankMnemonic        = "silver maximum item glass profit fragile require race decide sell gentle reflect success identify tray erosion gentle orchard wedding yard civil edge regret vote"
	TankBIP44Path       = "44'/60'/0'/0/0"

	DefaultGas    = uint64(20000000)
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

	// Initialize cosmos-sdk configuration
	cfg := sdk.GetConfig()
	cfg.SetCoinType(60)
	cfg.SetFullFundraiserPath(TankBIP44Path)
	cfg.SetBech32PrefixForAccount(config.DecimalPrefixAccAddr, config.DecimalPrefixAccPub)
	cfg.SetBech32PrefixForValidator(config.DecimalPrefixValAddr, config.DecimalPrefixValPub)
	cfg.SetBech32PrefixForConsensusNode(config.DecimalPrefixConsAddr, config.DecimalPrefixConsPub)
	cfg.Seal()

	// Initialize cosmos-sdk codec
	cdc := codec.New()
	sdk.RegisterCodec(cdc)
	cdc.RegisterConcrete(coin.MsgSendCoin{}, "coin/SendCoin", nil)
	cdc.RegisterConcrete(coin.MsgBuyCoin{}, "coin/BuyCoin", nil)
	cdc.RegisterConcrete(coin.MsgSellCoin{}, "coin/SellCoin", nil)
	codec.RegisterCrypto(cdc)
	cdc.Seal()

	rootPath := os.ExpandEnv(RootPath)

	// Initialize and prepare keybase
	keybase, err := keys.NewKeyring(sdk.KeyringServiceName(), TankKeyringBackend, rootPath, nil)
	if err != nil {
		log.Fatalf("ERROR: Unable to initialize keybase: %v", err)
	}
	keybase.CreateAccount(TankAccountName, TankMnemonic, "", TankAccountPassword, TankBIP44Path, keys.Secp256k1)

	return &Provider{
		cdc:     cdc,
		keybase: keybase,
	}
}

func (p *Provider) CreateAccount(name, password string) (Account, error) {
	entropySeed, err := bip39.NewEntropy(256)
	if err != nil {
		return Account{}, err
	}
	mnemonic, err := bip39.NewMnemonic(entropySeed)
	if err != nil {
		return Account{}, err
	}

	info, err := p.keybase.CreateAccount(name, mnemonic, password, password, TankBIP44Path, keys.Secp256k1)
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
	configPath := flag.String("cfg", "cfg.json", "Path to cfg")

	flag.Parse()

	cfg, err := ImportConfig(*configPath)
	if err != nil {
		log.Println(err)
		return
	}

	provider := NewProvider()

	if cfg.CountAccounts == 0 {
		cfg.CountAccounts = 20
	}

	accounts := make([]Account, cfg.CountAccounts)

	mainAddr, err := sdk.AccAddressFromBech32(TankAddress)
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
		Name:      TankAccountName,
		Sequence:  &mainSequence,
		Password:  TankAccountPassword,
	}

	for i := range accounts {
		accounts[i], err = provider.CreateAccount("tank"+strconv.Itoa(i), "12345678")
		if err != nil {
			log.Println(err)
			return
		}
	}

	err = provider.SendAll(mainAccount, accounts, helpers.BipToPip(sdk.NewInt(1)), "tDCL")
	if err != nil {
		log.Println("Init send", err)
		return
	}

	time.Sleep(time.Second * 10)

	err = provider.SendAll(mainAccount, accounts, helpers.BipToPip(sdk.NewInt(1)), "TEST4")
	if err != nil {
		log.Println("Init send", err)
		return
	}

	time.Sleep(time.Second * 10)

	for i := range accounts {
		_, accounts[i].AccNumber, err = GetSequenceAndAccNumber(accounts[i].Address.String())
		if err != nil {
			log.Println(err)
			return
		}
		atomic.StoreUint64(accounts[i].Sequence, 0)
	}

	for i := range accounts[:len(accounts)/2] {
		go func(accountNum int, accountNumNext int) {
			for {
				log.Println("Send ", accounts[accountNum].Address.String())
				err = provider.SendCoin(accounts[accountNum], accounts[accountNumNext], sdk.NewInt(rand.Int63n(99)+1).Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(13), nil))))
				if err != nil {
					log.Println(err)
				}
				time.Sleep(cfg.TimeoutMs.Send)
			}
		}(i, (i+1)%(len(accounts)/2))
	}

	for _, account := range accounts[len(accounts)/2:] {
		go func(account Account) {
			for {
				log.Println("Buy ", account.Address.String())
				err = provider.BuyCoin("TEST4", "tDCL", sdk.NewInt(1000000000000000), sdk.NewInt(1000000000000000000000000), account)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(cfg.TimeoutMs.Buy)
				err = provider.SellCoin("tDCL", "TEST4", sdk.NewInt(1), sdk.NewInt(1000000000000000), account)
				if err != nil {
					log.Println(err)
				}
				time.Sleep(cfg.TimeoutMs.Sell)
			}
		}(account)
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

func (p *Provider) SendCoin(sender, receiver Account, amount sdk.Int) error {
	memo := "tank send"
	txEncoder := auth.DefaultTxEncoder(p.cdc)
	txBldr := auth.NewTxBuilder(
		txEncoder,
		sender.AccNumber, atomic.LoadUint64(sender.Sequence),
		DefaultGas, DefaultGasAdj,
		false, ChainID, memo, nil, nil,
	).WithKeybase(p.keybase)

	msgs := []sdk.Msg{coin.NewMsgSendCoin(sender.Address, "tDCL", amount, receiver.Address)}

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
		log.Println("Sequence = ", atomic.LoadUint64(sender.Sequence))
		log.Printf("Broadcast error: code: %d, log: %s", broadcastResp.Result.Code, broadcastResp.Result.Log)
	} else {
		log.Println("Broadcast hash: ", broadcastResp.Result.Hash)
		atomic.AddUint64(sender.Sequence, 1)
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
	memo := "tank send"
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

	broadcastResp := BroadcastResponse{}
	err = json.Unmarshal(respBody, &broadcastResp)
	if err != nil {
		return err
	}
	if broadcastResp.Result.Code != 0 {
		log.Println("Sequence = ", atomic.LoadUint64(buyer.Sequence))
		log.Printf("Broadcast error: code: %d, log: %s", broadcastResp.Result.Code, broadcastResp.Result.Log)
	} else {
		log.Println("Broadcast hash: ", broadcastResp.Result.Hash)
		atomic.AddUint64(buyer.Sequence, 1)
	}
	return nil
}

func (p *Provider) SellCoin(coinToBuy, coinToSell string, amountToBuy, amountToSell sdk.Int, buyer Account) error {
	memo := "tank send"
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

	broadcastResp := BroadcastResponse{}
	err = json.Unmarshal(respBody, &broadcastResp)
	if err != nil {
		return err
	}
	if broadcastResp.Result.Code != 0 {
		log.Println("Sequence = ", atomic.LoadUint64(buyer.Sequence))
		log.Printf("Broadcast error: code: %d, log: %s", broadcastResp.Result.Code, broadcastResp.Result.Log)
	} else {
		log.Println("Broadcast hash: ", broadcastResp.Result.Hash)
		atomic.AddUint64(buyer.Sequence, 1)
	}
	return nil
}

func (p *Provider) SendAll(sender Account, accounts []Account, amount sdk.Int, token string) error {
	memo := "tank send"
	txEncoder := auth.DefaultTxEncoder(p.cdc)
	txBldr := auth.NewTxBuilder(
		txEncoder,
		sender.AccNumber, atomic.LoadUint64(sender.Sequence),
		DefaultGas, DefaultGasAdj,
		false, ChainID, memo, nil, nil,
	).WithKeybase(p.keybase)

	msgs := make([]sdk.Msg, len(accounts))
	for i, account := range accounts {
		msgs[i] = coin.NewMsgSendCoin(sender.Address, token, amount, account.Address)
	}

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
		log.Println("Sequence = ", atomic.LoadUint64(sender.Sequence))
		log.Printf("Broadcast error: code: %d, log: %s", broadcastResp.Result.Code, broadcastResp.Result.Log)
	} else {
		log.Println("Broadcast hash: ", broadcastResp.Result.Hash)
		atomic.AddUint64(sender.Sequence, 1)
	}

	return nil
}
