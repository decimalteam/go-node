package keeper

import (
	"encoding/binary"
	"fmt"
	"log"
	"net/http"
	"io"
	"os"
	"os/exec"
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

// ClearUpgradePlan clears any schedule upgrade
func (k Keeper) ClearUpgradePlan(ctx sdk.Context) {
	store := ctx.KVStore(k.storeKey)
	store.Delete(types.PlanKey())
}

// ApplyUpgrade will execute the handler associated with the Plan and mark the plan as done.
func (k Keeper) ApplyUpgrade(ctx sdk.Context, plan types.Plan) error {
	k.ClearUpgradePlan(ctx)
	/* workingDir, err := os.Getwd()
	if err != nil {
	 panic(err)
	}
   
	bin := workingDir + os.Args[0] */

	if k.EnsureBinary("update_decd") != nil {
		fmt.Println("Go download")
		/* err = k.DownloadBinary("decd2",k.GetUpdateUrl())
		if err != nil {
			panic(err)
		} */
	}

	bin := os.Args[0]

	syscall.Unlink(bin)
	err := os.Rename(filepath.Dir(bin)+"/update_decd", bin)
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

	
	return nil
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

// {
// 		"linux/amd64": "https://domain.com/bin"
//	}

func OSArch() string {
	return fmt.Sprintf("%s/%s", runtime.GOOS, runtime.GOARCH)
}

// EnsureBinary ensures the file exists and is executable, or returns an error
func (k *Keeper) EnsureBinary(path string) error {
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


//Function to generate the url to download the update
func (k *Keeper) GetUpdateUrl() string {
	return fmt.Sprintf("https://decimal/%s/update_decd",OSArch())
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
