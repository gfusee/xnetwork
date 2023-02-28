package main

import (
	"io"
	"testing"

	"github.com/multiversx/mx-chain-tools-go/trieTools/trieToolsCommon"
	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/common"
	"github.com/multiversx/mx-chain-tools-go/trieTools/zeroBalanceSystemAccountChecker/mocks"
	"github.com/stretchr/testify/require"
)

func TestReadTokensWithNonce(t *testing.T) {
	workingDir := "working-dir"
	tokensDir := "tokens-dir"

	file1Name := "shard0"
	file2Name := "shard1"
	file3Name := "shard3"

	file1 := &mocks.FileStub{
		NameCalled: func() string {
			return file1Name
		},
	}
	file2 := &mocks.FileStub{
		NameCalled: func() string {
			return file2Name
		},
	}
	file3 := &mocks.FileStub{
		NameCalled: func() string {
			return file3Name
		},
		IsDirCalled: func() bool {
			return true
		},
	}

	adr1 := "adr1"
	sysAccAddr := "sysAccAddr"

	adr1Tokens := map[string]struct{}{
		"token1-r-0": {},
		"token2-r-0": {},
		"esdt1-rand": {},
	}
	sysAccTokensShard0 := map[string]struct{}{
		"token3-r-0": {},
		"token3-r-1": {},
		"esdt1-rand": {},
	}
	addressTokensMapShard0 := trieToolsCommon.NewAddressTokensMap()
	addressTokensMapShard0.Add(adr1, adr1Tokens)
	addressTokensMapShard0.Add(sysAccAddr, sysAccTokensShard0)

	adr2 := "adr2"
	adr2Tokens := map[string]struct{}{
		"token3-r-0": {},
		"token4-r-0": {},
	}
	sysAccTokensShard1 := map[string]struct{}{
		"token3-r-0": {},
		"token3-r-2": {},
		"token5-r-0": {},
		"esdt2-rand": {},
	}
	addressTokensMapShard1 := trieToolsCommon.NewAddressTokensMap()
	addressTokensMapShard1.Add(adr2, adr2Tokens)
	addressTokensMapShard1.Add(sysAccAddr, sysAccTokensShard1)

	openCt := 0
	fileHandlerStub := &mocks.FileHandlerStub{
		GetwdCalled: func() (dir string, err error) {
			return workingDir, nil
		},
		ReadDirCalled: func(dirname string) ([]common.FileInfo, error) {
			require.Equal(t, workingDir+"/"+tokensDir, dirname)
			return []common.FileInfo{file1, file2, file3}, nil
		},

		OpenCalled: func(name string) (io.Reader, error) {
			openCt++

			switch openCt {
			case 1:
				require.Equal(t, workingDir+"/"+tokensDir+"/"+file1Name, name)
			case 2:
				require.Equal(t, workingDir+"/"+tokensDir+"/"+file2Name, name)
			default:
				require.Fail(t, "should not have opened another file")
			}

			return nil, nil
		},
		ReadAllCalled: func(r io.Reader) ([]byte, error) {
			switch openCt {
			case 1:
				return jsonMarshaller.Marshal(addressTokensMapShard0.GetMapCopy())
			case 2:
				return jsonMarshaller.Marshal(addressTokensMapShard1.GetMapCopy())
			default:
				require.Fail(t, "should not have read another file")
			}

			return nil, nil
		},
	}

	reader, err := newAddressTokensMapFileReader(fileHandlerStub, jsonMarshaller)
	require.Nil(t, err)

	globalTokens, shardTokens, err := reader.readTokensWithNonce(tokensDir)
	require.Nil(t, err)

	expectedGlobalTokensMap := trieToolsCommon.NewAddressTokensMap()
	delete(adr1Tokens, "esdt1-rand")
	delete(sysAccTokensShard0, "esdt1-rand")
	delete(sysAccTokensShard1, "esdt2-rand")
	expectedGlobalTokensMap.Add(adr1, adr1Tokens)
	expectedGlobalTokensMap.Add(adr2, adr2Tokens)
	expectedGlobalTokensMap.Add(sysAccAddr, map[string]struct{}{
		"token3-r-0": {},
		"token3-r-1": {},
		"token3-r-2": {},
		"token5-r-0": {},
	})
	require.Equal(t, expectedGlobalTokensMap, globalTokens)

	expectedShardTokens := make(map[uint32]trieToolsCommon.AddressTokensMap)
	expectedAddressTokensMapShard0 := trieToolsCommon.NewAddressTokensMap()
	expectedAddressTokensMapShard0.Add(adr1, adr1Tokens)
	expectedAddressTokensMapShard0.Add(sysAccAddr, sysAccTokensShard0)
	expectedShardTokens[0] = expectedAddressTokensMapShard0

	expectedAddressTokensMapShard1 := trieToolsCommon.NewAddressTokensMap()
	expectedAddressTokensMapShard1.Add(adr2, adr2Tokens)
	expectedAddressTokensMapShard1.Add(sysAccAddr, sysAccTokensShard1)
	expectedShardTokens[1] = expectedAddressTokensMapShard1
	require.Equal(t, expectedShardTokens, shardTokens)
}
