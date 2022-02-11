package keeper

import (
	"crypto/sha256"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"os/exec"
	"path"
	"path/filepath"
	"runtime"
	"syscall"

	ncfg "bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
	"github.com/go-ini/ini"
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

func (k Keeper) setDone(ctx sdk.Context, name string) {
	store := prefix.NewStore(ctx.KVStore(k.storeKey), types.DoneByte)
	bz := make([]byte, 8)
	binary.BigEndian.PutUint64(bz, uint64(ctx.BlockHeight()))
	store.Set([]byte(name), bz)
}

func (k Keeper) ClearUpgradePlan(ctx sdk.Context) {
	oldPlan, found := k.GetUpgradePlan(ctx)
	if found {
		k.ClearIBCState(ctx, oldPlan.Height)
	}
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PlanKey())
}

func (k Keeper) ClearIBCState(ctx sdk.Context, lastHeight int64) {
	// delete IBC client and consensus state from store if this is IBC plan
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.UpgradedClientKey(lastHeight))
	store.Delete(types.UpgradedConsStateKey(lastHeight))
}

// ApplyUpgrade will execute the handler associated with the Plan and mark the plan as done.
func (k Keeper) ApplyUpgrade(ctx sdk.Context, plan types.Plan) error {
	mapping := plan.Mapping()
	if mapping == nil {
		return fmt.Errorf("error: mapping decode")
	}

	currPath := filepath.Dir(os.Args[0])
	for i, name := range ncfg.NameFiles {
		downloadName := k.GetDownloadName(name)
		if _, err := os.Stat(downloadName); os.IsNotExist(err) {
			return err
		}

		hashes, ok := mapping[k.OSArch()]
		if !ok {
			return fmt.Errorf("error: mapping[os] undefined")
		}

		if !fileHashEqual(downloadName, hashes[i]) {
			os.Remove(downloadName)
			return fmt.Errorf("error: hash does not match")
		}

		currBin := filepath.Join(currPath, name)
		mode, err := getMode(currBin)
		if err != nil {
			os.Remove(downloadName)
			return err
		}

		err = MarkExecutableWithMode(downloadName, mode)
		if err != nil {
			os.Remove(downloadName)
			return err
		}

		ok = runIsSuccess(downloadName)
		if !ok {
			os.Remove(downloadName)
			return fmt.Errorf("error: file not running")
		}

		syscall.Unlink(currBin)
		err = os.Rename(downloadName, currBin)
		if err != nil {
			os.Remove(downloadName)
			return err
		}
	}

	return nil
}

// MarkExecutable will try to set the executable bits if not already set
// Fails if file doesn't exist or we cannot set those bits
func MarkExecutableWithMode(path string, mode os.FileMode) error {
	return os.Chmod(path, mode|0111)
}

func getMode(path string) (os.FileMode, error) {
	info, err := os.Stat(path)
	if err != nil {
		return 0, fmt.Errorf("stating binary: %w", err)
	}
	return info.Mode().Perm(), nil
}

func runIsSuccess(nameFile string) bool {
	cmd := exec.Command(nameFile, "version")
	err := cmd.Run()
	return err == nil
}

// Generate name of download file.
func (k Keeper) GetDownloadName(name string) string {
	baseFile := fmt.Sprintf("%s.nv", name)
	nameFile := filepath.Join(filepath.Dir(os.Args[0]), baseFile)
	return nameFile
}

func (k Keeper) GenerateUrl(urlName string) string {
	// example: "linux/ubuntu/20.04"
	u, err := url.Parse(k.OSArch())
	if err != nil {
		return ""
	}

	// example: "http://127.0.0.1/90500/decd"
	myUrl, err := url.Parse(urlName)
	if err != nil {
		return ""
	}

	// result: "http://127.0.0.1/90500/linux/ubuntu/20.04/decd"
	return fmt.Sprintf("%s/%s", myUrl.ResolveReference(u), path.Base(myUrl.Path))
}

// Check if page exists.
func (k Keeper) UrlPageExist(urlPage string) bool {
	resp, err := http.Head(urlPage)
	if err != nil {
		return false
	}
	return resp.StatusCode == 200
}

//Download file by url
func (k Keeper) DownloadBinary(ctx sdk.Context, filepath string, url string) error {
	// Get the data
	resp, err := http.Get(url)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("download binary from \"%s\" with '%s'", url, err.Error()))
		return err
	}
	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("create binary file \"%s\" for \"%s\" with '%s'", filepath, url, err.Error()))
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	if err != nil {
		ctx.Logger().Error(fmt.Sprintf("write binary file \"%s\" for \"%s\" with '%s'", filepath, url, err.Error()))
	}
	return err
}

// Detect OS to create a url
func (k Keeper) OSArch() string {
	switch runtime.GOOS {
	case "windows", "darwin":
		return runtime.GOOS
	case "linux":
		distr := readOSRelease("ID")
		if distr == "" {
			return ""
		}
		version := readOSRelease("VERSION_ID")
		if version == "" {
			return ""
		}
		return fmt.Sprintf("linux/%s/%s", distr, version)
	default:
		return ""
	}
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

//Get the hash of the download file, then check what was in the transaction
func fileHashEqual(nameFile, hash string) bool {
	f, err := os.Open(nameFile)
	if err != nil {
		return false
	}
	defer f.Close()

	h := sha256.New()
	if _, err := io.Copy(h, f); err != nil {
		return false
	}
	return hash == hex.EncodeToString(h.Sum(nil))
}

// Read the file under /etc/os-release to get the distribution name
func readOSRelease(key string) string {
	const cfgfile = "/etc/os-release"
	cfg, err := ini.Load(cfgfile)
	if err != nil {
		return ""
	}
	return cfg.Section("").Key(key).String()
}
