package keeper

import (
	"encoding/binary"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/url"
	"net/http"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"
	"syscall"
	"time"

	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/hashicorp/go-getter"
	
)

const (
	rootName        = "cosmovisor"
	genesisDir      = "genesis"
	upgradesDir     = "upgrades"
	currentLink     = "current"
	upgradeFilename = "upgrade-info.json"
)

// GetUpgradePlan returns the currently scheduled Plan if any, setting havePlan to true if there is a scheduled
// upgrade or false if there is none
func (k Keeper) GetUpgradePlan(ctx sdk.Context) (plan types.Plan, havePlan bool) {
	store := ctx.KVStore(k.storeKey)
	bz := store.Get(types.PlanKey())
	if bz == nil {
		return plan, false
	}

	k.cdc.MustUnmarshalBinaryBare(bz, &plan)
	return plan, true
}

// IsSkipHeight checks if the given height is part of skipUpgradeHeights
func (k Keeper) IsSkipHeight(height int64) bool {
	return k.skipUpgradeHeights[height]
}

// ClearUpgradePlan clears any schedule upgrade
func (k Keeper) ClearUpgradePlan(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PlanKey())
}

// ApplyUpgrade will execute the handler associated with the Plan and mark the plan as done.
func (k Keeper) ApplyUpgrade(ctx sdk.Context, plan types.Plan) error {
	k.ClearUpgradePlan(ctx)

	bin := os.Args[0]

	syscall.Unlink(bin)
	err := os.Rename(plan.Name, bin)
	if err != nil {
		panic(err)
	}

	cmd := exec.Command(bin, "start")

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin

	err = cmd.Start()
	if err != nil {
		log.Fatalf("cmd.Start() failed with %s\n", err)
	}

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGQUIT, syscall.SIGTERM)

	// go func() {
	// 	sig := <-sigs
	// 	if err := cmd.Process.Signal(sig); err != nil {
	// 		log.Fatal(bin, "terminated. Error:", err)
	// 	}
	// }()

	// sigs := make(chan os.Signal, 1)
	// signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	// go func() {
	// 	sig := <-sigs
	// 	if err := cmd.Process.Signal(sig); err != nil {
	// 		log.Fatal(bin, "terminated. Error:", err)
	// 		os.Exit(0)
	// 	}
	// }()

	// os.Exit(0)

	// cfg := Config{
	// 	Home: os.Getenv("HOME/.decimal"),
	// }

	// _, err := GetDownloadURL(plan.Info)
	// if err != nil {
	// 	panic(err)
	// }
	// if err := EnsureBinary(cfg.UpgradeBin(plan.Name)); err != nil {
	// 	return fmt.Errorf("downloaded binary doesn't check out: %w", err)
	// }

	//handler := k.upgradeHandlers[plan.Name]
	//if handler == nil {
	//	panic("ApplyUpgrade should never be called without first checking HasHandler")
	//}
	//
	//handler(ctx, plan)
	//
	//k.ClearUpgradePlan(ctx)
	//k.setDone(ctx, plan.Name)
	return nil
}

// Config is the information passed in to control the daemon
type Config struct {
	Home                  string
	Name                  string
	AllowDownloadBinaries bool
	RestartAfterUpgrade   bool
	PollInterval          time.Duration
	UnsafeSkipBackup      bool
}

// UpgradeBin is the path to the binary for the named upgrade
func (cfg *Config) UpgradeBin(upgradeName string) string {
	return filepath.Join(cfg.UpgradeDir(upgradeName), "bin")
}

// UpgradeDir is the directory named upgrade
func (cfg *Config) UpgradeDir(upgradeName string) string {
	safeName := url.PathEscape(upgradeName)
	return filepath.Join(cfg.Home, rootName, upgradesDir, safeName)
}

// UpgradeInfo is the update details created by `x/upgrade/keeper.DumpUpgradeInfoToDisk`.
type UpgradeInfo struct {
	Name   string `json:"name"`
	Info   string `json:"info"`
	Height uint   `json:"height"`
	ToDownload uint `json:"toDownload"`
}

// MarkExecutable will try to set the executable bits if not already set
// Fails if file doesn't exist or we cannot set those bits
func MarkExecutable(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("stating binary: %w", err)
	}
	// end early if world exec already set
	if info.Mode()&0001 == 1 {
		return nil
	}
	// now try to set all exec bits
	newMode := info.Mode().Perm() | 0111
	return os.Chmod(path, newMode)
}

func (k *Keeper) DownloadBinary(filepath string, url string) error {

	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}


// DownloadBinary will grab the binary and place it in the proper directory
/* func DownloadBinary(cfg *Config, info UpgradeInfo) error {
	url, err := GetDownloadURL(info.Info)
	if err != nil {
		return err
	}

	// download into the bin dir (works for one file)
	binPath := cfg.UpgradeBin(info.Name)
	err = getter.GetFile(binPath, url)

	// if this fails, let's see if it is a zipped directory
	if err != nil {
		dirPath := cfg.UpgradeDir(info.Name)
		err = getter.Get(dirPath, url)
		if err != nil {
			return err
		}
		err = EnsureBinary(binPath)
		// copy binary to binPath from dirPath if zipped directory don't contain bin directory to wrap the binary
		if err != nil {
			err = copy.Copy(filepath.Join(dirPath, cfg.Name), binPath)
			if err != nil {
				return err
			}
		}
	}

	// if it is successful, let's ensure the binary is executable
	return MarkExecutable(binPath)
} */

// {
// 		"linux/amd64": "https://domain.com/bin"
//	}

func OSArch() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// GetDownloadURL will check if there is an arch-dependent binary specified in Info
func GetDownloadURL(info string) (string, error) {
	doc := strings.TrimSpace(info)
	// if this is a url, then we download that and try to get a new doc with the real info
	if _, err := url.Parse(doc); err == nil {
		tmpDir, err := ioutil.TempDir("", "upgrade-manager-reference")
		if err != nil {
			return "", fmt.Errorf("create tempdir for reference file: %w", err)
		}
		defer os.RemoveAll(tmpDir)

		refPath := filepath.Join(tmpDir, "ref")
		if err := getter.GetFile(refPath, doc); err != nil {
			return "", fmt.Errorf("downloading reference link %s: %w", doc, err)
		}

		refBytes, err := ioutil.ReadFile(refPath)
		if err != nil {
			return "", fmt.Errorf("reading downloaded reference: %w", err)
		}
		// if download worked properly, then we use this new file as the binary map to parse
		doc = string(refBytes)
	}

	// check if it is the upgrade config
	var config types.UpgradeConfig

	if err := json.Unmarshal([]byte(doc), &config); err == nil {
		url, ok := config.Binaries[OSArch()]
		if !ok {
			url, ok = config.Binaries["any"]
		}
		if !ok {
			return "", fmt.Errorf("cannot find binary for os/arch: neither %s, nor any", OSArch())
		}

		return url, nil
	}

	return "", errors.New("upgrade info doesn't contain binary map")
}

// EnsureBinary ensures the file exists and is executable, or returns an error
func EnsureBinary(path string) error {
	info, err := os.Stat(path)
	if err != nil {
		return fmt.Errorf("cannot stat dir %s: %w", path, err)
	}

	if !info.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", info.Name())
	}

	// this checks if the world-executable bit is set (we cannot check owner easily)
	exec := info.Mode().Perm() & 0001
	if exec == 0 {
		return fmt.Errorf("%s is not world executable", info.Name())
	}

	return nil
}

// ScheduleUpgrade schedules an upgrade based on the specified plan.
// If there is another Plan already scheduled, it will overwrite it
// (implicitly cancelling the current plan)
func (k Keeper) ScheduleUpgrade(ctx sdk.Context, plan types.Plan) error {
	if err := plan.ValidateBasic(); err != nil {
		return err
	}

	if !plan.Time.IsZero() {
		if !plan.Time.After(ctx.BlockHeader().Time) {
			return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "upgrade cannot be scheduled in the past")
		}
	} else if plan.Height <= ctx.BlockHeight() {
		return sdkerrors.Wrap(sdkerrors.ErrInvalidRequest, "upgrade cannot be scheduled in the past")
	}

	if k.GetDoneHeight(ctx, plan.Name) != 0 {
		return sdkerrors.Wrapf(sdkerrors.ErrInvalidRequest, "upgrade with name %s has already been completed", plan.Name)
	}

	bz := k.cdc.MustMarshalBinaryBare(plan)
	store := ctx.KVStore(k.storeKey)
	store.Set(types.PlanKey(), bz)

	return nil
}

// GetDoneHeight returns the height at which the given upgrade was executed
func (k Keeper) GetDoneHeight(ctx sdk.Context, name string) int64 {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DoneByte)
	bz := store.Get([]byte(name))
	if len(bz) == 0 {
		return 0
	}

	return int64(binary.BigEndian.Uint64(bz))
}
