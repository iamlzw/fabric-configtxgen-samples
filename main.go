package main

import (
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/hyperledger/fabric/common/tools/configtxgen/encoder"
	"github.com/hyperledger/fabric/common/tools/configtxgen/localconfig"
	cb "github.com/hyperledger/fabric/protos/common"
	"github.com/hyperledger/fabric/protos/utils"
)

func load(profile string) *localconfig.Profile{
	configPath := "/home/www/go/src/github.com/hyperledger/fabric-samples/first-network"
	return localconfig.Load(profile,configPath)
}

func outputBlock(){
	blockDest := filepath.Join("block/")

	config1 := load("TwoOrgsOrdererGenesis")

	doOutputBlock(config1, "byfn-sys-channel", blockDest)

	config2 := load("TwoOrgsChannel")
	doOutputChannelCreateTx(config2,nil,"mychannel",blockDest)
}

func doOutputBlock(config *localconfig.Profile, channelID string, outputBlock string){
	pgen := encoder.New(config)
	fmt.Println("Generating genesis block")
	if config.Orderer == nil {
		fmt.Println("refusing to generate block which is missing orderer section")
	}
	if config.Consortiums == nil {
		fmt.Println("Genesis block does not contain a consortiums group definition.  This block cannot be used for orderer bootstrap.")
	}
	genesisBlock := pgen.GenesisBlockForChannel(channelID)
	fmt.Println("Writing genesis block")
	err := ioutil.WriteFile("genesis.block", utils.MarshalOrPanic(genesisBlock), 0644)
	if err != nil {
		fmt.Println("Error writing genesis block: %s", err)
	}
}

func doOutputChannelCreateTx(conf, baseProfile *localconfig.Profile, channelID string, outputChannelCreateTx string) {
	fmt.Println("Generating new channel configtx")

	var configtx *cb.Envelope
	var err error
	if baseProfile == nil {
		configtx, err = encoder.MakeChannelCreationTransaction(channelID, nil, conf)
	} else {
		configtx, err = encoder.MakeChannelCreationTransactionWithSystemChannelContext(channelID, nil, conf, baseProfile)
	}
	if err != nil {
		fmt.Println(err)
	}

	fmt.Println("Writing new channel tx")
	err = ioutil.WriteFile("channel.tx", utils.MarshalOrPanic(configtx), 0644)
	if err != nil {
		fmt.Errorf("Error writing channel create tx: %s", err)
	}
}


func main(){
	outputBlock()
}
