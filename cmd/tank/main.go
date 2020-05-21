package main

import (
	"encoding/hex"
	"encoding/json"
	"flag"
	"fmt"
	authexported "github.com/cosmos/cosmos-sdk/x/auth/exported"
	"io/ioutil"
	"log"
	"math/big"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync/atomic"
	"time"

	"github.com/cosmos/cosmos-sdk/codec"
	"github.com/cosmos/cosmos-sdk/crypto/keys"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/cosmos/cosmos-sdk/x/auth"
	"github.com/cosmos/go-bip39"

	"bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/utils/helpers"
	"bitbucket.org/decimalteam/go-node/x/coin"
)

const (
	ChainID             = config.ChainID
	RootPath            = "$HOME/.decimal/cli"
	TankKeyringBackend  = keys.BackendTest
	TankAccountName     = "tank"
	TankAccountPassword = "12345678"
	TankAddress         = "dx1esffyu0wxk6eez77fhzdxfgvjp4646hqm9sx6c"
	TankMnemonic        = "silver maximum item glass profit fragile require race decide sell gentle reflect success identify tray erosion gentle orchard wedding yard civil edge regret vote"
	TankBIP44Path       = "44'/60'/0'/0/0"

	DefaultGas    = uint64(20000000)
	DefaultGasAdj = float64(1.1)

	RPCPrefix  = "http://localhost:26657"
	RESTPrefix = "http://localhost:1317"
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
	cdc.RegisterInterface((*sdk.Msg)(nil), nil)
	cdc.RegisterConcrete(coin.MsgSendCoin{}, "coin/SendCoin", nil)
	cdc.RegisterConcrete(coin.MsgBuyCoin{}, "coin/BuyCoin", nil)
	cdc.RegisterConcrete(coin.MsgSellCoin{}, "coin/SellCoin", nil)
	cdc.RegisterConcrete(coin.MsgCreateCoin{}, "coin/CreateCoin", nil)
	cdc.RegisterConcrete(coin.MsgSellAllCoin{}, "coin/SellAllCoin", nil)

	cdc.RegisterInterface((*authexported.Account)(nil), nil)
	cdc.RegisterConcrete(&auth.BaseAccount{}, "test/validator/BaseAccount", nil)
	codec.RegisterCrypto(cdc)

	rootPath := os.ExpandEnv(RootPath)

	// Initialize and prepare keybase
	keybase, err := keys.NewKeyring(sdk.KeyringServiceName(), TankKeyringBackend, rootPath, nil)
	if err != nil {
		log.Fatalf("ERROR: Unable to initialize keybase: %v", err)
	}
	_, err = keybase.CreateAccount(TankAccountName, TankMnemonic, "", TankAccountPassword, TankBIP44Path, keys.Secp256k1)
	if err != nil {
		log.Println(err)
	}

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
		AccNumber: 0,
		Name:      name,
		Sequence:  new(uint64),
		Password:  password,
	}, nil
}

type Worker struct {
	account Account
	ch      chan func(Account) error
}

func (w *Worker) Run() {
	go func() {
		for {
			select {
			case fn := <-w.ch:
				err := fn(w.account)
				if err != nil {
					log.Println(err)
				}
			}
		}
	}()
}

type Timeout struct {
	LastTime time.Time
	Timeout  time.Duration
}

type Distributor struct {
	TxWeights       map[string]int
	TxTimeouts      map[string]Timeout
	Workers         []Worker
	CountTxPerBlock int
	provider        *Provider
	Coins           []string
}

func NewDistributor(txWeights map[string]int, txTimeouts map[string]Timeout, countTxPerBlock int, provider *Provider, coins []string) *Distributor {
	return &Distributor{
		TxWeights:       txWeights,
		TxTimeouts:      txTimeouts,
		Workers:         []Worker{},
		CountTxPerBlock: countTxPerBlock,
		provider:        provider,
		Coins:           coins,
	}
}

func (d *Distributor) AddWorker(account Account) {
	worker := Worker{
		ch:      make(chan func(Account) error),
		account: account,
	}
	d.Workers = append(d.Workers, worker)
}

func (d *Distributor) Run() {
	for i := range d.Workers {
		d.Workers[i].Run()
	}
	totalWeight := 0
	for _, weight := range d.TxWeights {
		totalWeight += weight
	}
	go func() {
		for {
			log.Println(d.Coins)
			count := 0
			for tx, timeout := range d.TxTimeouts {
				if timeout.LastTime.Add(timeout.Timeout).Before(time.Now()) {
					d.createTx(tx, count)
					timeout.LastTime = time.Now()
					d.TxTimeouts[tx] = timeout
					count++
				}
			}
			for tx, weight := range d.TxWeights {
				countTx := int(float32(weight) / float32(totalWeight) * float32(d.CountTxPerBlock))
				for i := 0; i < countTx; i++ {
					d.createTx(tx, count)
					count++
				}
			}
			time.Sleep(time.Second * 7)
		}
	}()
}

func (d *Distributor) createTx(tx string, count int) {
	if count == len(d.Workers) {
		return
	}
	switch tx {
	case "send":
		d.Workers[count].ch <- func(account Account) error {
			return d.provider.SendCoin(account, d.Workers[rand.Intn(len(d.Workers))].account, sdk.NewInt(rand.Int63n(99)+1).Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(13), nil))))
		}
	case "buy":
		d.Workers[count].ch <- func(account Account) error {
			return d.provider.BuyCoin(d.Coins[rand.Intn(len(d.Coins))], "tDEL", sdk.NewInt(1000000000000000), Pow(sdk.NewInt(1), 25), account)
		}
	case "sell":
		d.Workers[count].ch <- func(account Account) error {
			return d.provider.SellCoin("tDEL", d.Coins[rand.Intn(len(d.Coins))], sdk.NewInt(1), sdk.NewInt(1000000000000000), account)
		}
	case "create_coin":
		d.Workers[count].ch <- func(account Account) error {
			d.Coins = append(d.Coins, "TEST"+strconv.Itoa(len(d.Coins)))
			return d.provider.CreateCoin("TEST"+strconv.Itoa(len(d.Coins)-1), "TEST"+strconv.Itoa(len(d.Coins)-1), 50, helpers.BipToPip(sdk.NewInt(100000)), helpers.BipToPip(sdk.NewInt(100000)), helpers.BipToPip(sdk.NewInt(100000000000000)), account)
		}
	case "sell_all":
		d.Workers[count].ch <- func(account Account) error {
			return d.provider.SellAllCoins(account, "tDEL", d.Coins[rand.Intn(len(d.Coins))], sdk.NewInt(1))
		}
	}
}

func main() {
	configPath := flag.String("cfg", "cfg.json", "Path to cfg")

	flag.Parse()

	cfg, err := ImportConfig(*configPath)
	if err != nil {
		log.Println(err)
		return
	}
	timeout := make(map[string]Timeout)
	for tx, duration := range cfg.Timeout {
		timeout[tx] = Timeout{Timeout: duration, LastTime: time.Unix(0, 0)}
	}

	provider := NewProvider()

	if cfg.CountAccounts == 0 {
		cfg.CountAccounts = 20
	}

	coins, err := GetCoins()
	if err != nil {
		log.Println(err)
		return
	}

	var testCoins []string
	for _, c := range coins {
		if strings.HasPrefix(c, "TEST") {
			testCoins = append(testCoins, c)
		}
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

	err = provider.SendAll(mainAccount, accounts, helpers.BipToPip(sdk.NewInt(100000)), "tDEL")
	if err != nil {
		log.Println("Init send", err)
		return
	}

	time.Sleep(time.Second * 10)

	if len(testCoins) == 0 {
		err = provider.CreateCoin("TEST0", "TEST0", 50, helpers.BipToPip(sdk.NewInt(100000)), helpers.BipToPip(sdk.NewInt(100000)), helpers.BipToPip(sdk.NewInt(100000000000000)), mainAccount)
		if err != nil {
			log.Println(err)
			return
		}
		testCoins = append(testCoins, "TEST0")
		time.Sleep(time.Second * 10)
	}

	err = provider.SendAll(mainAccount, accounts, helpers.BipToPip(sdk.NewInt(1)), "TEST0")
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

	distributor := NewDistributor(cfg.Weights, timeout, cfg.CountAccounts, provider, testCoins)

	for _, account := range accounts {
		distributor.AddWorker(account)
	}

	distributor.Run()

	select {}
}

type BroadcastResponse struct {
	Result struct {
		Code int    `json:"code"`
		Log  string `json:"log"`
		Hash string `json:"hash"`
	}
}

func GetSequenceAndAccNumber(address string) (uint64, uint64, error) {
	resp, err := http.Get(RESTPrefix + "/auth/accounts/" + address)
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

func (p *Provider) SendCoin(sender, receiver Account, amount sdk.Int) error {
	return p.SendTx([]sdk.Msg{coin.NewMsgSendCoin(sender.Address, "tDEL", amount, receiver.Address)}, sender)
}

func (p *Provider) BuyCoin(coinToBuy, coinToSell string, amountToBuy, amountToSell sdk.Int, buyer Account) error {
	return p.SendTx([]sdk.Msg{coin.NewMsgBuyCoin(buyer.Address, coinToBuy, coinToSell, amountToBuy, amountToSell)}, buyer)
}

func (p *Provider) SellCoin(coinToBuy, coinToSell string, amountToBuy, amountToSell sdk.Int, buyer Account) error {
	return p.SendTx([]sdk.Msg{coin.NewMsgSellCoin(buyer.Address, coinToBuy, coinToSell, amountToSell, amountToBuy)}, buyer)
}

func (p *Provider) CreateCoin(title string, symbol string, crr uint, initVolume sdk.Int, initReserve sdk.Int, limitVolume sdk.Int, sender Account) error {
	return p.SendTx([]sdk.Msg{coin.NewMsgCreateCoin(title, crr, symbol, initVolume, initReserve, limitVolume, sender.Address)}, sender)
}

func (p *Provider) SellAllCoins(seller Account, coinToBuy string, coinToSell string, amountToBuy sdk.Int) error {
	return p.SendTx([]sdk.Msg{coin.NewMsgSellAllCoin(seller.Address, coinToBuy, coinToSell, amountToBuy)}, seller)
}

func (p *Provider) SendAll(sender Account, accounts []Account, amount sdk.Int, token string) error {
	msgs := make([]sdk.Msg, len(accounts))
	for i, account := range accounts {
		msgs[i] = coin.NewMsgSendCoin(sender.Address, token, amount, account.Address)
	}
	return p.SendTx(msgs, sender)
}

func (p *Provider) SendTx(messages []sdk.Msg, sender Account) error {
	memo := "tank send"
	txEncoder := auth.DefaultTxEncoder(p.cdc)
	txBldr := auth.NewTxBuilder(
		txEncoder,
		sender.AccNumber, atomic.LoadUint64(sender.Sequence),
		DefaultGas, DefaultGasAdj,
		false, ChainID, memo, nil, nil,
	).WithKeybase(p.keybase)

	tx, err := txBldr.BuildAndSign(sender.Name, sender.Password, messages)
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
		log.Println("Sequence = ", atomic.LoadUint64(sender.Sequence))
		log.Printf("Broadcast error: code: %d, log: %s", broadcastResp.Result.Code, broadcastResp.Result.Log)
	} else {
		log.Println("Broadcast hash: ", broadcastResp.Result.Hash)
		atomic.AddUint64(sender.Sequence, 1)
	}
	return nil
}

func Pow(value sdk.Int, power int64) sdk.Int {
	return value.Mul(sdk.NewIntFromBigInt(new(big.Int).Exp(big.NewInt(10), big.NewInt(power), nil)))
}

func GetCoins() ([]string, error) {
	resp, err := http.Get(RESTPrefix + "/coins")
	if err != nil {
		return nil, err
	}
	defer func() {
		err = resp.Body.Close()
		if err != nil {
			log.Println(err)
		}
	}()
	decoder := json.NewDecoder(resp.Body)

	var account struct {
		Result []string `json:"result"`
	}
	err = decoder.Decode(&account)
	if err != nil {
		return nil, err
	}

	return account.Result, nil
}
