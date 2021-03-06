package blkdb

import (
	"fmt"
	"io/ioutil"
	"reflect"
	"testing"

	"github.com/copernet/copernicus/log"
	"github.com/copernet/copernicus/model"
	"github.com/copernet/copernicus/model/block"
	"github.com/copernet/copernicus/model/blockindex"
	"github.com/copernet/copernicus/persist/db"
	"github.com/copernet/copernicus/util"
)

func initBlockDB() {
	path, err := ioutil.TempDir("/tmp", "blockIndex")
	if err != nil {
		panic(fmt.Sprintf("generate temp db path failed: %s\n", err))
	}
	bc := &BlockTreeDBConfig{
		Do: &db.DBOption{
			FilePath:  path,
			CacheSize: 1 << 20,
		},
	}

	InitBlockTreeDB(bc)
}

func TestWRTxIndex(t *testing.T) {
	initBlockDB()

	// test TxIndex && init block pos
	dbpos := block.NewDiskBlockPos(12, 12)
	//init tx pos
	dtp := block.NewDiskTxPos(dbpos, 1)
	wantVal := &block.DiskTxPos{
		BlockIn:    dbpos,
		TxOffsetIn: 1,
	}
	txindexs := make(map[util.Hash]block.DiskTxPos)
	h := util.HashFromString("000000002dd5588a74784eaa7ab0507a18ad16a236e7b1ce69f00d7ddfb5d011")
	txindexs[*h] = *dtp

	//test Write and Read TxIndex
	err := GetInstance().WriteTxIndex(txindexs)
	if err != nil {
		t.Errorf("write tx index failed: %v\n", err)
	}
	txpos, err := GetInstance().ReadTxIndex(h)
	if err != nil {
		t.Errorf("read tx index failed: %v\n", err)
	}

	if !reflect.DeepEqual(wantVal, txpos) {
		t.Errorf("the wantVal not equal except value: %v, %v\n", wantVal, txpos)
	}
}

func TestWriteFlag(t *testing.T) {
	initBlockDB()
	//test flag: value is false
	err := GetInstance().WriteFlag("b", false)
	if err != nil {
		t.Errorf("write flag failed: %v\n", err)
	}
	res := GetInstance().ReadFlag("b")
	if !res {
		t.Errorf("the flag should is true: %v\n", res)
	}

	//test flag: value is true
	err = GetInstance().WriteFlag("b", true)
	if err != nil {
		t.Errorf("write flag failed: %v\n", err)
	}
	res2 := GetInstance().ReadFlag("b")
	if res2 {
		t.Errorf("the flag should is false: %v\n", res2)
	}
}

func TestWriteReindexing(t *testing.T) {
	initBlockDB()
	//test write reindex: reindexing value is true
	err := GetInstance().WriteReindexing(true)
	if err != nil {
		t.Errorf("write the index failed: %v\n", err)
	}
	readReindex := GetInstance().ReadReindexing()
	if !readReindex {
		t.Errorf("the reindexing should is true: %v\n", readReindex)
	}

	//test write reindex: reindexing value is false
	err = GetInstance().WriteReindexing(false)
	if err != nil {
		t.Errorf("write the index failed: %v\n", err)
	}
	readReindexAgain := GetInstance().ReadReindexing()
	if readReindexAgain {
		t.Errorf("the reindexing should is false: %v\n", readReindexAgain)
	}
}

func TestReadBlockFileInfo(t *testing.T) {
	initBlockDB()
	//the block file info exist
	_, err := writeBlockFile()
	if err != nil {
		t.Errorf("write block file: %v\n ", err)
	}
	bfi, err := GetInstance().ReadBlockFileInfo(1)
	if err != nil {
		t.Errorf("read block file info<000> failed: %v\n", err)
	}
	log.Info("blockFileInfo value is:%v", bfi)
	if bfi == nil {
		t.Errorf("the block fileInfo not equal nil: %v\n", bfi)
	}

	//the block file info not exist
	bfi, err = GetInstance().ReadBlockFileInfo(1111)
	if err == nil {
		t.Errorf("read block file info<1111> failed: %v\n", err)
	}
	if bfi != nil {
		t.Errorf("the block file info error, please check: %v\n", bfi)
	}
}

func writeBlockFile() (*util.Hash, error) {

	blkHeader := block.NewBlockHeader()
	bfi1 := make(map[int32]*block.BlockFileInfo)
	fi := block.NewBlockFileInfo()
	fi.UndoSize = 0
	fi.Size = 1
	fi.Blocks = 1
	fi.HeightFirst = 1
	fi.HeightLast = 2
	bfi1[1] = fi

	//init block header
	blkHeader1 := block.NewBlockHeader()
	blkHeader.Time = uint32(1534822771)
	blkHeader.Version = 536870912
	blkHeader.Bits = 486604799
	preHash := util.HashFromString("00000000000001bcd6b635a1249dfbe76c0d001592a7219a36cd9bbd002c7238")
	merkleRoot := util.HashFromString("7e814211a7de289a490380c0c20353e0fd4e62bf55a05b38e1628e0ea0b4fd3d")
	blkHeader.HashPrevBlock = *preHash
	blkHeader.Nonce = 1391785674
	blkHeader.MerkleRoot = *merkleRoot
	//init block index
	blkidx := blockindex.NewBlockIndex(blkHeader1)
	blockHash := blkidx.GetBlockHash() //14508459b221041eab257d2baaa7459775ba748246c8403609eb708f0e57e74b
	blkidxs := make([]*blockindex.BlockIndex, 0, 10)
	blkidxs = append(blkidxs, blkidx)
	err := GetInstance().WriteBatchSync(bfi1, 1, blkidxs)
	return blockHash, err
}

func TestReadLastBlockFile(t *testing.T) {
	initBlockDB()
	_, err := writeBlockFile()
	if err != nil {
		t.Errorf("write blockFileInfo failed: %v\n", err)
	}
	lastFile, err := GetInstance().ReadLastBlockFile()
	if err != nil {
		t.Errorf("read last block file failed: %v\n", err)
	}

	bfi, err := GetInstance().ReadBlockFileInfo(lastFile)
	if err != nil {
		t.Errorf("read last block fileInfo failed: %v\n", err)
	}
	log.Info("last blockFileInfo value is: %v\n", bfi)
	if bfi == nil {
		t.Errorf("the last blockFileInfo not equal nil, the value is: %v\n", bfi)
	}
}

func TestWriteBatchSync(t *testing.T) {
	initBlockDB()
	blkHeader := block.NewBlockHeader()
	bfi1 := make(map[int32]*block.BlockFileInfo)
	fi := block.NewBlockFileInfo()
	fi.UndoSize = 0
	fi.Size = 1
	fi.Blocks = 1
	fi.HeightFirst = 1
	fi.HeightLast = 2
	bfi1[1] = fi

	//init block header
	blkHeader1 := block.NewBlockHeader()
	blkHeader.Time = uint32(1534822771)
	blkHeader.Version = 536870912
	blkHeader.Bits = 486604799
	preHash := util.HashFromString("00000000000001bcd6b635a1249dfbe76c0d001592a7219a36cd9bbd002c7238")
	merkleRoot := util.HashFromString("7e814211a7de289a490380c0c20353e0fd4e62bf55a05b38e1628e0ea0b4fd3d")
	blkHeader.HashPrevBlock = *preHash
	blkHeader.Nonce = 1391785674
	blkHeader.MerkleRoot = *merkleRoot
	//init block index
	blkidx := blockindex.NewBlockIndex(blkHeader1)
	hhash := blkidx.GetBlockHash() //14508459b221041eab257d2baaa7459775ba748246c8403609eb708f0e57e74b
	blkidxs := make([]*blockindex.BlockIndex, 0, 10)
	blkidxs = append(blkidxs, blkidx)
	err := GetInstance().WriteBatchSync(bfi1, 1, blkidxs)
	if err != nil {
		t.Errorf("write blockFileInfo failed:%v", err)
	}

	lastFile, err := GetInstance().ReadLastBlockFile()
	if lastFile != 1 {
		t.Error("the lastfile value is 1, please check.")
	}
	if err != nil {
		t.Errorf("read last block file failed: %v\n", err)
	}

	bfi, err := GetInstance().ReadBlockFileInfo(lastFile)
	if err != nil {
		t.Errorf("read last block fileInfo failed: %v\n", err)
	}
	log.Info("last blockFileInfo value is:%v", bfi)
	if !reflect.DeepEqual(bfi, bfi1[1]) {
		t.Errorf("the last blockFileInfo not equal nil, the value is: %v\n", bfi)
	}

	blkidxMap := make(map[util.Hash]*blockindex.BlockIndex)
	blkidxMap[*hhash] = blkidx
	ret := GetInstance().LoadBlockIndexGuts(blkidxMap, model.ActiveNetParams)
	if !ret {
		t.Error("load block index guts failed, please check.")
	}

	if !reflect.DeepEqual(blkidxMap[*hhash], blkidx) {
		t.Error("the blkidxMap should equal blkidx, please check")
	}
}
