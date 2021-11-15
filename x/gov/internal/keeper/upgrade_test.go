package keeper

/*
import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"syscall"
	"testing"
	"time"

	ncfg "bitbucket.org/decimalteam/go-node/config"
	"bitbucket.org/decimalteam/go-node/x/gov/internal/types"
	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
)

const (
	upgradeUrlPath = "https://testnet-repo.decimalchain.com/620004"
	upgradeJson    = `{"linux/centos/7":["bc60aa8ec4aa722200561f3b824944558ecebdef5f962f8029306a333babde8f","34a5c14288617ce8559f01548ab0c9e615db8159f0c4c77d036268d27e780345"],"linux/centos/8":["0dfd92ec87ee3159a6604c2b5719c0f00cf1788033778a580097472c6ac4be19","0c3077a5f6f44f4085dd2f11549d2d21bb9409b405e1c32bda66bcf61b5fbfac"],"linux/debian/10":["dc1490280f946beeeb9b592c34fc5e6aafd7b45e607330161b1ac156e62e4def","f6ade3a869449531b8759b1f30b43ce9169600c450f44f505a604315bb53cc5c"],"linux/ubuntu/16.04":["a0649a437fb07144649e2b3e344d671069c5bd9915b50701be04c5bc6be0244e","3230443d32e598416c9dc23e69d4319726334f665f78aaa1fe0f360266fdfb62"],"linux/ubuntu/18.04":["81b827174e1e0793ff0812a2c04fa1bc8006274c02ebafc81d67d8d77ef3f7d9","fe5faa26fb8f1eadda69023ec4d807ac050dd77b951b3d4f7b6d1ca81f79e8ab"],"linux/ubuntu/20.04":["6d04e93c3a90529df63b25e8ae7037d236915eab0c7b4de3bd83cf084f4df0eb","6a40906d8a55d714fca2c1cd64543f7e3077ea768d8b1a8c51980d1861321bbd"]}`
)

func createUpgradePlan(t *testing.T) (sdk.Context, Keeper, types.Plan) {
	ctx, _, keeper, _, _, _ := createTestInput(t, false, 100)
	newPlan := types.Plan{Name: upgradeUrlPath, Time: time.Time{}, Height: 10, Info: upgradeJson, ToDownload: 1}
	keeper.ScheduleUpgrade(ctx, newPlan)
	return ctx, keeper, newPlan
}

// create
func testCreateUpgradePlan(t *testing.T) (sdk.Context, Keeper, types.Plan) {
	ctx, keeper, newPlan := createUpgradePlan(t)
	plan, found := keeper.GetUpgradePlan(ctx)
	require.True(t, found)
	require.Equal(t, newPlan, plan)
	return ctx, keeper, plan
}

func TestCreateUpgradePlan(t *testing.T) {
	testCreateUpgradePlan(t)
}

func TestClearUpgradePlan(t *testing.T) {
	ctx, keeper, _ := testCreateUpgradePlan(t)

	keeper.ClearUpgradePlan(ctx)
	_, found := keeper.GetUpgradePlan(ctx)
	require.False(t, found)
}

// generate
func testGenerateUrl(t *testing.T) string {
	k := Keeper{}
	result := k.GenerateUrl(fmt.Sprintf("%s/decd", upgradeUrlPath))
	require.Equal(t, fmt.Sprintf("%s/%s/decd", upgradeUrlPath, k.OSArch()), result)
	return result
}

func TestGenerateUrl(t *testing.T) {
	testGenerateUrl(t)
}

// download
func testGetDownloadName(t *testing.T) string {
	k := Keeper{}
	downloadName := k.GetDownloadName("decd")
	require.Equal(t, filepath.Join(filepath.Dir(os.Args[0]), "decd.nv"), downloadName)
	return downloadName
}

func TestGetDownloadName(t *testing.T) {
	testGetDownloadName(t)
}

// (generate + download) -> binary
func testDownloadBinary(t *testing.T) (string, string) {
	newUrl := testGenerateUrl(t)
	downloadName := testGetDownloadName(t)

	err := Keeper{}.DownloadBinary(downloadName, newUrl)
	if err != nil {
		os.Remove(downloadName)
		require.NoError(t, err)
	}

	return newUrl, downloadName
}

func TestDownloadBinary(t *testing.T) {
	_, downloadName := testDownloadBinary(t)
	os.Remove(downloadName)
}

// (generate + download) -> binary -> mode
func testGetMode(t *testing.T) (string, os.FileMode) {
	_, downloadName := testDownloadBinary(t)

	mode, err := getMode(downloadName)
	require.NoError(t, err)

	require.Equal(t, os.FileMode(0644), mode)
	return downloadName, mode
}

func TestGetMode(t *testing.T) {
	downloadName, _ := testGetMode(t)
	os.Remove(downloadName)
}

func TestRunIsSuccess(t *testing.T) {
	downloadName, mode := testGetMode(t)
	require.NoError(t, MarkExecutableWithMode(downloadName, mode))

	mode, err := getMode(downloadName)
	require.NoError(t, err)

	require.Equal(t, os.FileMode(0775), mode)
	require.True(t, runIsSuccess(downloadName))
}

func TestUrlPageExist(t *testing.T) {
	newUrl := testGenerateUrl(t)
	require.True(t, Keeper{}.UrlPageExist(newUrl))
}

func TestFileHashEqual(t *testing.T) {
	_, k, plan := testCreateUpgradePlan(t)

	mapping := plan.Mapping()
	require.NotEqual(t, nil, mapping)

	hashes, ok := mapping[k.OSArch()]
	require.True(t, ok)

	_, downloadName := testDownloadBinary(t)
	require.True(t, fileHashEqual(downloadName, hashes[0]))
}

func TestApplyUpgrade(t *testing.T) {
	_, k, plan := testCreateUpgradePlan(t)

	for _, name := range ncfg.NameFiles {
		newUrl := k.GenerateUrl(fmt.Sprintf("%s/%s", plan.Name, name))
		require.NotEqual(t, "", newUrl)

		require.True(t, k.UrlPageExist(newUrl))

		downloadName := k.GetDownloadName(name)
		err := k.DownloadBinary(downloadName, newUrl)
		if err != nil {
			os.Remove(downloadName)
			require.NoError(t, err)
		}

		defer os.Remove(downloadName)
	}

	mapping := plan.Mapping()
	require.NotEqual(t, nil, mapping)

	currPath := filepath.Dir(os.Args[0])
	for i, name := range ncfg.NameFiles {
		downloadName := k.GetDownloadName(name)

		_, err := os.Stat(downloadName)
		require.False(t, os.IsNotExist(err))

		hashes, ok := mapping[k.OSArch()]
		require.True(t, ok)
		require.True(t, fileHashEqual(downloadName, hashes[i]))

		currBin := filepath.Join(currPath, name)
		err = copyFile(downloadName, currBin)
		if err != nil {
			os.Remove(currBin)
			require.NoError(t, err)
		}
		defer os.Remove(currBin)

		mode, err := getMode(currBin)
		require.NoError(t, err)

		err = MarkExecutableWithMode(downloadName, mode)
		require.NoError(t, err)

		ok = runIsSuccess(downloadName)
		require.True(t, ok)

		err = syscall.Unlink(currBin)
		require.NoError(t, err)

		err = os.Rename(downloadName, currBin)
		require.NoError(t, err)
	}
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	if err != nil {
		return err
	}
	return out.Close()
}
*/
