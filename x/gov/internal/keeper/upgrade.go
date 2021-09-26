package keeper

import (
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	"encoding/json"
	"errors"
	"fmt"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"io/ioutil"
	"net/url"
	"os"
	"path/filepath"
	"runtime"
	"strings"
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
func (k Keeper) ApplyUpgrade(ctx sdk.Context, plan types.Plan) {
	//handler := k.upgradeHandlers[plan.Name]
	//if handler == nil {
	//	panic("ApplyUpgrade should never be called without first checking HasHandler")
	//}
	//
	//handler(ctx, plan)
	//
	//k.ClearUpgradePlan(ctx)
	//k.setDone(ctx, plan.Name)
}

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
		//if err := getter.GetFile(refPath, doc); err != nil {
		//	return "", fmt.Errorf("downloading reference link %s: %w", doc, err)
		//}

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
