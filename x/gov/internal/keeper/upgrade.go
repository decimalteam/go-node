package keeper

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"path"
	"path/filepath"
	"runtime"
	"syscall"

	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"github.com/cosmos/cosmos-sdk/store/prefix"
	sdk "github.com/cosmos/cosmos-sdk/types"
	sdkerrors "github.com/cosmos/cosmos-sdk/types/errors"
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

	// plan.Name => url path to file
	myUrl, err := url.Parse(plan.Name)
	if err != nil {
		log.Fatal(err)
		return err
	}

	currbin := os.Args[0]
	baseFile := path.Base(myUrl.Path)
	nameFile := filepath.Join(filepath.Dir(currbin), baseFile)

	MarkExecutable(nameFile)

	syscall.Unlink(currbin)
	err = os.Rename(nameFile, currbin)
	if err != nil {
		panic(err)
	}

	k.setDone(ctx, plan.Name)
	k.ClearUpgradePlan(ctx)

	return nil
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
		fmt.Println("1", err)
		return err
	}

	defer resp.Body.Close()

	// Create the file
	out, err := os.Create(filepath)
	if err != nil {
		fmt.Println("2", err)
		return err
	}
	defer out.Close()

	// Write the body to file
	_, err = io.Copy(out, resp.Body)
	return err
}

func (k Keeper) OSArch() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
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
