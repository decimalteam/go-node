package keeper

import (
	"bytes"
	"crypto/sha256"
	"encoding/hex"
	"sync"
	"time"

	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"

	sdk "github.com/cosmos/cosmos-sdk/types"
	"github.com/stretchr/testify/require"
	"github.com/tendermint/tendermint/libs/log"
)

////////////////////////
// Temporary file system
type TmpFilesystem struct {
	basepath string
}

func NewTmpFilesystem() (tfs *TmpFilesystem, err error) {
	tfs = &TmpFilesystem{}
	tfs.basepath, err = os.MkdirTemp("", "gotest")
	if err != nil {
		return
	}
	return tfs, nil
}

func (tfs *TmpFilesystem) PathForFilename(filename string) string {
	return filepath.Join(tfs.basepath, filename)
}

func (tfs *TmpFilesystem) Close() {
	os.RemoveAll(tfs.basepath)
}

////////////////////////
// test Logger
// implement tendermint logger inteface

type TestLogger struct {
	InfoMsgs  []string
	ErrorMsgs []string
	mtx       sync.Mutex
}

func NewTestLogger() *TestLogger {
	return &TestLogger{}
}

func (tl *TestLogger) Debug(msg string, keyvals ...interface{}) {

}

func (tl *TestLogger) Info(msg string, keyvals ...interface{}) {
	tl.mtx.Lock()
	defer tl.mtx.Unlock()
	tl.InfoMsgs = append(tl.InfoMsgs, fmt.Sprintf(msg, keyvals...))
}

func (tl *TestLogger) Error(msg string, keyvals ...interface{}) {
	tl.mtx.Lock()
	defer tl.mtx.Unlock()
	tl.ErrorMsgs = append(tl.ErrorMsgs, fmt.Sprintf(msg, keyvals...))
}

func (tl *TestLogger) With(keyvals ...interface{}) log.Logger {
	return tl
}

////////////////////////

func genRandomBytes(n int) []byte {
	b := make([]byte, n)
	for i := range b {
		b[i] = byte(rand.Intn(256))
	}
	return b
}

func calcHash(source []byte) string {
	r := bytes.NewReader(source)
	h := sha256.New()
	io.Copy(h, r)
	return hex.EncodeToString(h.Sum(nil))

}

func TestDownloadCases(t *testing.T) {
	testBytes := genRandomBytes(1024 * 1024)
	normalHash := calcHash(testBytes)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/file":
			{
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Write(testBytes)
			}
		case "/500_error":
			{
				w.WriteHeader(http.StatusInternalServerError)
				w.Write([]byte("error"))
			}
		case "/301_redirect":
			{
				http.Redirect(w, r, "/file", http.StatusMovedPermanently)
			}
		}
	}))
	defer ts.Close()

	tfs, err := NewTmpFilesystem()
	if err != nil {
		t.Error("NewTmpFilesystem error ", err.Error())
	}
	defer tfs.Close()
	testFilepath := tfs.PathForFilename("testfile")

	k := Keeper{}
	testSuite := []struct {
		msg        string
		url        string
		hashValue  string
		infoCount  int
		errorCount int
	}{
		{"normal download", "/file", normalHash, 3, 0},
		{"wrong hash", "/file", normalHash + "0", 2, 1},
		{"http error download", "/500_error", normalHash, 1, 1},
		{"redirect", "/301_redirect", normalHash, 3, 0},
	}
	for _, suite := range testSuite {
		logger := NewTestLogger()
		ctx := sdk.Context{}.WithLogger(logger)
		k.DownloadAndCheckHash(ctx, testFilepath, ts.URL+suite.url, suite.hashValue)
		//fmt.Println(logger.InfoMsgs)
		//fmt.Println(logger.ErrorMsgs)
		require.Equal(t, suite.infoCount, len(logger.InfoMsgs), suite.msg)
		require.Equal(t, suite.errorCount, len(logger.ErrorMsgs), suite.msg)
	}
}

func TestFileWriteClash(t *testing.T) {
	testBytes := genRandomBytes(1024 * 1024)
	normalHash := calcHash(testBytes)
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.String() {
		case "/file":
			{
				time.Sleep(time.Duration((rand.Intn(5) + 1) * int(time.Microsecond)))
				w.Header().Set("Content-Type", "application/octet-stream")
				w.Write(testBytes)
			}
		}
	}))
	defer ts.Close()

	tfs, err := NewTmpFilesystem()
	if err != nil {
		t.Error("NewTmpFilesystem error ", err.Error())
	}
	defer tfs.Close()
	testFilepath := tfs.PathForFilename("testfile")

	k := Keeper{}
	logger := NewTestLogger()
	ctx := sdk.Context{}.WithLogger(logger)
	wg := sync.WaitGroup{}
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			if _, err := os.Stat(testFilepath); os.IsNotExist(err) { //!! without this we get errors
				k.DownloadAndCheckHash(ctx, testFilepath, ts.URL+"/file", normalHash)
			}
			wg.Done()
		}()
	}
	wg.Wait()
	//fmt.Println(logger.ErrorMsgs)
	require.Equal(t, 0, len(logger.ErrorMsgs), "file write clash")
}

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

	mode, err := getMode(downloadName)
	require.NoError(t, err)

	if mode.Perm() != os.FileMode(0644) {
		err = os.Chmod(downloadName, os.FileMode(0644))
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

	require.Equal(t, os.FileMode(0755), mode)
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
