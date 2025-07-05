// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package gamecontract

import (
	"errors"
	"math/big"
	"strings"

	ethereum "github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/event"
)

// Reference imports to suppress errors if they are not otherwise used.
var (
	_ = errors.New
	_ = big.NewInt
	_ = strings.NewReader
	_ = ethereum.NotFound
	_ = bind.Bind
	_ = common.Big1
	_ = types.BloomLookup
	_ = event.NewSubscription
	_ = abi.ConvertType
)

// IGameContractGameInfo is an auto generated low-level Go binding around an user-defined struct.
type IGameContractGameInfo struct {
	GameId           *big.Int
	State            uint8
	Result           uint8
	FixedStakeAmount *big.Int
	CreatedAt        *big.Int
	EndedAt          *big.Int
	TotalWhiteStakes *big.Int
	TotalBlackStakes *big.Int
	WhitePlayerCount *big.Int
	BlackPlayerCount *big.Int
}

// IGameContractPlayerInfo is an auto generated low-level Go binding around an user-defined struct.
type IGameContractPlayerInfo struct {
	Team        uint8
	TotalStakes *big.Int
	MoveCount   *big.Int
	HasJoined   bool
}

// GamecontractMetaData contains all meta data concerning the Gamecontract contract.
var GamecontractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_authorizedBackend\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedBackend\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"calculateRewards\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"playerTotalStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"createGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fixedStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameResult\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"games\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameState\"},{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameResult\"},{\"name\":\"fixedStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalWhiteStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBlackStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"whitePlayerCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blackPlayerCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameInfo\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIGameContract.GameInfo\",\"components\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameState\"},{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameResult\"},{\"name\":\"fixedStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"endedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalWhiteStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalBlackStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"whitePlayerCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"blackPlayerCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameResult\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameResult\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayerInfo\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIGameContract.PlayerInfo\",\"components\":[{\"name\":\"team\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.Team\"},{\"name\":\"totalStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"moveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"hasJoined\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayerMoveCount\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayerTeam\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.Team\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayerTotalMoveCount\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"hasPlayerJoined\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isGameActive\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"joinTeam\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"team\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.Team\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"moves\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"chainId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"moveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"players\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"team\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.Team\"},{\"name\":\"totalStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"moveCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"hasJoined\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"recordMove\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint32\",\"internalType\":\"uint32\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"GameCreated\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"fixedStakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GameEnded\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"result\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIGameContract.GameResult\"},{\"name\":\"endedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"MoveRecorded\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"newMoveCount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"PlayerJoinedTeam\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"team\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIGameContract.Team\"}],\"anonymous\":false}]",
	Bin: "0x60a0346100dc57601f61126538819003918201601f19168301916001600160401b038311848410176100e0578084926020946040528339810103126100dc57516001600160a01b0381168082036100dc57156100855760805260405161117090816100f5823960805181818161010f015281816102ad0152818161062f015261081f0152f35b60405162461bcd60e51b815260206004820152602960248201527f417574686f72697a6564206261636b656e642063616e6e6f74206265207a65726044820152686f206164647265737360b81b6064820152608490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe60806040526004361015610011575f80fd5b5f3560e01c80630542172014610ccf57806307286e2f14610c83578063117a5b9014610bdf5780631a60a4e514610ba65780632b497f8514610b365780633305a7a914610aee5780633ccd10e914610a2f5780634294b984146109df57806345b7c6fc1461099857806347e1d5501461084e57806349430bcf1461080a57806360104cef14610612578063c0f809211461059b578063c3d4bd8d146103d7578063c544953114610287578063ed2df26d1461024e5763feee346f146100d4575f80fd5b3461024a577f2e5e5932470e4421499604d32a27d62dc99709b56faaea15de88d17fe1603aac604061010536610d43565b9391929061013d337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611054565b835f525f602052610152825f20541515610dcf565b835f525f60205261017560ff6001845f2001541661016f81610d2c565b156110d7565b835f526001602052815f2060018060a01b0382165f526020526101a060ff6003845f20015416610e1d565b835f526002602052815f2063ffffffff86165f52602052815f2060018060a01b0382165f52602052815f206101d58154611118565b9055835f526001602052815f2060018060a01b0382165f526020526002825f20016102008154611118565b9055835f526002602052815f2063ffffffff86165f52602052815f2060018060a01b0382165f52602052815f205463ffffffff835196168652602086015260018060a01b031693a3005b5f80fd5b3461024a57602036600319011261024a576004355f525f602052602060ff600160405f20015460081c166102856040518092610d36565bf35b3461024a57604036600319011261024a57600435602435600481101561024a576102db337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611054565b815f525f6020526102f160405f20541515610dcf565b815f525f60205261030f60ff600160405f2001541661016f81610d2c565b8015610386575f82815260208190526040908190206001908101805460ff191690911781557fad5c2fcf2f9497233e36f5fa5a21487acf87e422dd4d6535fd81cd01b37925d992906103629082906110b8565b835f525f602052426004835f20015561037d82518092610d36565b426020820152a2005b60405162461bcd60e51b815260206004820152602360248201527f43616e6e6f7420656e642067616d652077697468206f6e676f696e67207265736044820152621d5b1d60ea1b6064820152608490fd5b3461024a57604036600319011261024a57602435600435600282101561024a57805f525f60205261040d60405f20541515610dcf565b805f525f60205261042b60ff600160405f2001541661016f81610d2c565b5f81815260016020908152604080832033845290915290206003015460ff1661055657610456610d7a565b61045f83610d2c565b828152600360208201915f8352604081015f8152606082019360018552855f52600160205260405f2060018060a01b0333165f5260205260405f2092516104a581610d2c565b6104ae81610d2c565b60ff80198554169116178355516001830155516002820155019051151560ff801983541691161790556104e082610d2c565b8161053757805f525f602052600760405f20016104fd8154611118565b90555b6040519161050d81610d2c565b82527f345d86e28bee17ad4bc57f12601ea9a4817a1afaac78868cd4f7af041a0638fa60203393a3005b805f525f602052600860405f200161054f8154611118565b9055610500565b60405162461bcd60e51b815260206004820152601f60248201527f506c6179657220616c7265616479206a6f696e656420746869732067616d65006044820152606490fd5b3461024a57604036600319011261024a576105b4610d16565b6004355f52600160205260405f209060018060a01b03165f52602052608060405f2060ff8154169060018101549060ff600360028301549201541691604051936105fd81610d2c565b84526020840152604083015215156060820152f35b3461024a57604036600319011261024a5760043560243561065d337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614611054565b815f525f60205260405f20546107cf57801561077f5761067b610dae565b9082825260208201905f825260408301925f84526060810193828552608082019042825260a083015f815260c084015f815260e08501905f82526101008601925f84526101208701945f86528b5f525f60205260405f20975188556001880199516106e581610d2c565b6106ee81610d2c565b60ff80198c54169116178a555195600487101561076b577fbd3f84a55fdc44f970adc6ccbaa93e78d7da5f3fc37c4b7da578d47c9ad2e46a9a61073560089860409c6110b8565b516002890155516003880155516004870155516005860155516006850155516007840155519101558151908152426020820152a2005b634e487b7160e01b5f52602160045260245ffd5b60405162461bcd60e51b815260206004820152602260248201527f4669786564207374616b65206d7573742062652067726561746572207468616e604482015261020360f41b6064820152608490fd5b60405162461bcd60e51b815260206004820152601360248201527247616d6520616c72656164792065786973747360681b6044820152606490fd5b3461024a575f36600319011261024a576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b3461024a57602036600319011261024a575f61012061086b610dae565b8281528260208201528260408201528260608201528260808201528260a08201528260c08201528260e08201528261010082015201526004355f525f60205261014060405f206108b9610dae565b8154815260018201549091602083019160ff81166108d681610d2c565b83526108ec60ff604086019260081c1682610e11565b60028201546060850190815260038301546080860190815260048401549060a0870191825260058501549260c0880193845261096a60068701549560e08a0196875260086007890154986101008c01998a520154986101208b01998a526040519a518b525161095a81610d2c565b60208b01525160408a0190610d36565b5160608801525160808701525160a08601525160c08501525160e08401525161010083015251610120820152f35b3461024a57606036600319011261024a5760206109d76004356109b9610d16565b815f525f84526109ce60405f20541515610dcf565b60443591610e69565b604051908152f35b3461024a57604036600319011261024a576109f8610d16565b6004355f52600160205260405f209060018060a01b03165f52602052602060ff60405f20541660405190610a2b81610d2c565b8152f35b3461024a57604036600319011261024a57610a48610d16565b5f6060610a53610d7a565b82815282602082015282604082015201526004355f52600160205260405f209060018060a01b03165f52602052608060405f20610a8e610d7a565b60ff82541691610a9d83610d2c565b828252600181015460208301908152606060ff6003600285015494604087019586520154169301921515835260405193610ad681610d2c565b84525160208401525160408301525115156060820152f35b3461024a57610afc36610d43565b90915f52600260205263ffffffff60405f2091165f5260205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b3461024a57606036600319011261024a5760243563ffffffff8116810361024a57604435906001600160a01b038216820361024a576004355f52600260205263ffffffff60405f2091165f5260205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b3461024a57602036600319011261024a576004355f525f602052602060ff600160405f20015416610bd681610d2c565b60405190158152f35b3461024a57602036600319011261024a576004355f525f60205261014060405f2080549060018101549060ff82169160028201546003830154600484015490600585015492610c5c60068701549560086007890154980154986040519a8b52610c4781610d2c565b60208b015260ff60408b019160081c16610d36565b6060880152608087015260a086015260c085015260e0840152610100830152610120820152f35b3461024a57604036600319011261024a57610c9c610d16565b6004355f52600160205260405f209060018060a01b03165f52602052602060ff600360405f200154166040519015158152f35b3461024a57604036600319011261024a57610ce8610d16565b6004355f52600160205260405f209060018060a01b03165f526020526020600260405f200154604051908152f35b602435906001600160a01b038216820361024a57565b6002111561076b57565b90600482101561076b5752565b606090600319011261024a57600435906024356001600160a01b038116810361024a579060443563ffffffff8116810361024a5790565b604051906080820182811067ffffffffffffffff821117610d9a57604052565b634e487b7160e01b5f52604160045260245ffd5b60405190610140820182811067ffffffffffffffff821117610d9a57604052565b15610dd657565b60405162461bcd60e51b815260206004820152601360248201527211d85b5948191bd95cc81b9bdd08195e1a5cdd606a1b6044820152606490fd5b600482101561076b5752565b15610e2457565b60405162461bcd60e51b815260206004820152601f60248201527f506c6179657220686173206e6f74206a6f696e656420746869732067616d65006044820152606490fd5b90815f525f60205260405f20916001610e80610dae565b93805485528181015494610120600860ff8816936020840194610ea281610d2c565b8552610eb760ff604086019a841c168a610e11565b6002810154606085015260038101546080850152600481015460a0850152600581015460c0850152600681015460e08501526007810154610100850152015491015251610f0381610d2c565b610f0c81610d2c565b0361101a575f52600160205260405f209060018060a01b03165f5260205260405f20610f7160ff6003610f3d610d7a565b9382815416610f4b81610d2c565b855260018101546020860152600281015460408601520154161515806060840152610e1d565b8215611013578151600481101561076b57600303610f8e57505090565b8151600481101561076b576001149182610ff6575b8215610fb9575b505015610fb45790565b505f90565b90915051600481101561076b576002149081610fd8575b505f80610faa565b6001915051610fe681610d2c565b610fef81610d2c565b145f610fd0565b9150805161100381610d2c565b61100c81610d2c565b1591610fa3565b5050505f90565b60405162461bcd60e51b815260206004820152601260248201527111d85b59481a185cc81b9bdd08195b99195960721b6044820152606490fd5b1561105b57565b60405162461bcd60e51b815260206004820152602f60248201527f4f6e6c7920617574686f72697a6564206261636b656e642063616e207065726660448201526e37b936903a3434b99030b1ba34b7b760891b6064820152608490fd5b90600481101561076b5761ff0082549160081b169061ff001916179055565b156110de57565b60405162461bcd60e51b815260206004820152601260248201527147616d65206973206e6f742061637469766560701b6044820152606490fd5b5f1981146111265760010190565b634e487b7160e01b5f52601160045260245ffdfea2646970667358221220426343423c3f76dc9479b06af314e7dace94c36b563f5f2097a72b7051a2996c64736f6c634300081e0033",
}

// GamecontractABI is the input ABI used to generate the binding from.
// Deprecated: Use GamecontractMetaData.ABI instead.
var GamecontractABI = GamecontractMetaData.ABI

// GamecontractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GamecontractMetaData.Bin instead.
var GamecontractBin = GamecontractMetaData.Bin

// DeployGamecontract deploys a new Ethereum contract, binding an instance of Gamecontract to it.
func DeployGamecontract(auth *bind.TransactOpts, backend bind.ContractBackend, _authorizedBackend common.Address) (common.Address, *types.Transaction, *Gamecontract, error) {
	parsed, err := GamecontractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GamecontractBin), backend, _authorizedBackend)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Gamecontract{GamecontractCaller: GamecontractCaller{contract: contract}, GamecontractTransactor: GamecontractTransactor{contract: contract}, GamecontractFilterer: GamecontractFilterer{contract: contract}}, nil
}

// Gamecontract is an auto generated Go binding around an Ethereum contract.
type Gamecontract struct {
	GamecontractCaller     // Read-only binding to the contract
	GamecontractTransactor // Write-only binding to the contract
	GamecontractFilterer   // Log filterer for contract events
}

// GamecontractCaller is an auto generated read-only Go binding around an Ethereum contract.
type GamecontractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GamecontractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type GamecontractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GamecontractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type GamecontractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// GamecontractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type GamecontractSession struct {
	Contract     *Gamecontract     // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// GamecontractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type GamecontractCallerSession struct {
	Contract *GamecontractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts       // Call options to use throughout this session
}

// GamecontractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type GamecontractTransactorSession struct {
	Contract     *GamecontractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts       // Transaction auth options to use throughout this session
}

// GamecontractRaw is an auto generated low-level Go binding around an Ethereum contract.
type GamecontractRaw struct {
	Contract *Gamecontract // Generic contract binding to access the raw methods on
}

// GamecontractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type GamecontractCallerRaw struct {
	Contract *GamecontractCaller // Generic read-only contract binding to access the raw methods on
}

// GamecontractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type GamecontractTransactorRaw struct {
	Contract *GamecontractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewGamecontract creates a new instance of Gamecontract, bound to a specific deployed contract.
func NewGamecontract(address common.Address, backend bind.ContractBackend) (*Gamecontract, error) {
	contract, err := bindGamecontract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Gamecontract{GamecontractCaller: GamecontractCaller{contract: contract}, GamecontractTransactor: GamecontractTransactor{contract: contract}, GamecontractFilterer: GamecontractFilterer{contract: contract}}, nil
}

// NewGamecontractCaller creates a new read-only instance of Gamecontract, bound to a specific deployed contract.
func NewGamecontractCaller(address common.Address, caller bind.ContractCaller) (*GamecontractCaller, error) {
	contract, err := bindGamecontract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &GamecontractCaller{contract: contract}, nil
}

// NewGamecontractTransactor creates a new write-only instance of Gamecontract, bound to a specific deployed contract.
func NewGamecontractTransactor(address common.Address, transactor bind.ContractTransactor) (*GamecontractTransactor, error) {
	contract, err := bindGamecontract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &GamecontractTransactor{contract: contract}, nil
}

// NewGamecontractFilterer creates a new log filterer instance of Gamecontract, bound to a specific deployed contract.
func NewGamecontractFilterer(address common.Address, filterer bind.ContractFilterer) (*GamecontractFilterer, error) {
	contract, err := bindGamecontract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &GamecontractFilterer{contract: contract}, nil
}

// bindGamecontract binds a generic wrapper to an already deployed contract.
func bindGamecontract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := GamecontractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gamecontract *GamecontractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gamecontract.Contract.GamecontractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gamecontract *GamecontractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gamecontract.Contract.GamecontractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gamecontract *GamecontractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gamecontract.Contract.GamecontractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Gamecontract *GamecontractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Gamecontract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Gamecontract *GamecontractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Gamecontract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Gamecontract *GamecontractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Gamecontract.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedBackend is a free data retrieval call binding the contract method 0x49430bcf.
//
// Solidity: function authorizedBackend() view returns(address)
func (_Gamecontract *GamecontractCaller) AuthorizedBackend(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "authorizedBackend")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AuthorizedBackend is a free data retrieval call binding the contract method 0x49430bcf.
//
// Solidity: function authorizedBackend() view returns(address)
func (_Gamecontract *GamecontractSession) AuthorizedBackend() (common.Address, error) {
	return _Gamecontract.Contract.AuthorizedBackend(&_Gamecontract.CallOpts)
}

// AuthorizedBackend is a free data retrieval call binding the contract method 0x49430bcf.
//
// Solidity: function authorizedBackend() view returns(address)
func (_Gamecontract *GamecontractCallerSession) AuthorizedBackend() (common.Address, error) {
	return _Gamecontract.Contract.AuthorizedBackend(&_Gamecontract.CallOpts)
}

// CalculateRewards is a free data retrieval call binding the contract method 0x45b7c6fc.
//
// Solidity: function calculateRewards(uint256 gameId, address player, uint256 playerTotalStakes) view returns(uint256)
func (_Gamecontract *GamecontractCaller) CalculateRewards(opts *bind.CallOpts, gameId *big.Int, player common.Address, playerTotalStakes *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "calculateRewards", gameId, player, playerTotalStakes)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// CalculateRewards is a free data retrieval call binding the contract method 0x45b7c6fc.
//
// Solidity: function calculateRewards(uint256 gameId, address player, uint256 playerTotalStakes) view returns(uint256)
func (_Gamecontract *GamecontractSession) CalculateRewards(gameId *big.Int, player common.Address, playerTotalStakes *big.Int) (*big.Int, error) {
	return _Gamecontract.Contract.CalculateRewards(&_Gamecontract.CallOpts, gameId, player, playerTotalStakes)
}

// CalculateRewards is a free data retrieval call binding the contract method 0x45b7c6fc.
//
// Solidity: function calculateRewards(uint256 gameId, address player, uint256 playerTotalStakes) view returns(uint256)
func (_Gamecontract *GamecontractCallerSession) CalculateRewards(gameId *big.Int, player common.Address, playerTotalStakes *big.Int) (*big.Int, error) {
	return _Gamecontract.Contract.CalculateRewards(&_Gamecontract.CallOpts, gameId, player, playerTotalStakes)
}

// Games is a free data retrieval call binding the contract method 0x117a5b90.
//
// Solidity: function games(uint256 gameId) view returns(uint256 gameId, uint8 state, uint8 result, uint256 fixedStakeAmount, uint256 createdAt, uint256 endedAt, uint256 totalWhiteStakes, uint256 totalBlackStakes, uint256 whitePlayerCount, uint256 blackPlayerCount)
func (_Gamecontract *GamecontractCaller) Games(opts *bind.CallOpts, gameId *big.Int) (struct {
	GameId           *big.Int
	State            uint8
	Result           uint8
	FixedStakeAmount *big.Int
	CreatedAt        *big.Int
	EndedAt          *big.Int
	TotalWhiteStakes *big.Int
	TotalBlackStakes *big.Int
	WhitePlayerCount *big.Int
	BlackPlayerCount *big.Int
}, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "games", gameId)

	outstruct := new(struct {
		GameId           *big.Int
		State            uint8
		Result           uint8
		FixedStakeAmount *big.Int
		CreatedAt        *big.Int
		EndedAt          *big.Int
		TotalWhiteStakes *big.Int
		TotalBlackStakes *big.Int
		WhitePlayerCount *big.Int
		BlackPlayerCount *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.State = *abi.ConvertType(out[1], new(uint8)).(*uint8)
	outstruct.Result = *abi.ConvertType(out[2], new(uint8)).(*uint8)
	outstruct.FixedStakeAmount = *abi.ConvertType(out[3], new(*big.Int)).(**big.Int)
	outstruct.CreatedAt = *abi.ConvertType(out[4], new(*big.Int)).(**big.Int)
	outstruct.EndedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)
	outstruct.TotalWhiteStakes = *abi.ConvertType(out[6], new(*big.Int)).(**big.Int)
	outstruct.TotalBlackStakes = *abi.ConvertType(out[7], new(*big.Int)).(**big.Int)
	outstruct.WhitePlayerCount = *abi.ConvertType(out[8], new(*big.Int)).(**big.Int)
	outstruct.BlackPlayerCount = *abi.ConvertType(out[9], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// Games is a free data retrieval call binding the contract method 0x117a5b90.
//
// Solidity: function games(uint256 gameId) view returns(uint256 gameId, uint8 state, uint8 result, uint256 fixedStakeAmount, uint256 createdAt, uint256 endedAt, uint256 totalWhiteStakes, uint256 totalBlackStakes, uint256 whitePlayerCount, uint256 blackPlayerCount)
func (_Gamecontract *GamecontractSession) Games(gameId *big.Int) (struct {
	GameId           *big.Int
	State            uint8
	Result           uint8
	FixedStakeAmount *big.Int
	CreatedAt        *big.Int
	EndedAt          *big.Int
	TotalWhiteStakes *big.Int
	TotalBlackStakes *big.Int
	WhitePlayerCount *big.Int
	BlackPlayerCount *big.Int
}, error) {
	return _Gamecontract.Contract.Games(&_Gamecontract.CallOpts, gameId)
}

// Games is a free data retrieval call binding the contract method 0x117a5b90.
//
// Solidity: function games(uint256 gameId) view returns(uint256 gameId, uint8 state, uint8 result, uint256 fixedStakeAmount, uint256 createdAt, uint256 endedAt, uint256 totalWhiteStakes, uint256 totalBlackStakes, uint256 whitePlayerCount, uint256 blackPlayerCount)
func (_Gamecontract *GamecontractCallerSession) Games(gameId *big.Int) (struct {
	GameId           *big.Int
	State            uint8
	Result           uint8
	FixedStakeAmount *big.Int
	CreatedAt        *big.Int
	EndedAt          *big.Int
	TotalWhiteStakes *big.Int
	TotalBlackStakes *big.Int
	WhitePlayerCount *big.Int
	BlackPlayerCount *big.Int
}, error) {
	return _Gamecontract.Contract.Games(&_Gamecontract.CallOpts, gameId)
}

// GetGameInfo is a free data retrieval call binding the contract method 0x47e1d550.
//
// Solidity: function getGameInfo(uint256 gameId) view returns((uint256,uint8,uint8,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_Gamecontract *GamecontractCaller) GetGameInfo(opts *bind.CallOpts, gameId *big.Int) (IGameContractGameInfo, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getGameInfo", gameId)

	if err != nil {
		return *new(IGameContractGameInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IGameContractGameInfo)).(*IGameContractGameInfo)

	return out0, err

}

// GetGameInfo is a free data retrieval call binding the contract method 0x47e1d550.
//
// Solidity: function getGameInfo(uint256 gameId) view returns((uint256,uint8,uint8,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_Gamecontract *GamecontractSession) GetGameInfo(gameId *big.Int) (IGameContractGameInfo, error) {
	return _Gamecontract.Contract.GetGameInfo(&_Gamecontract.CallOpts, gameId)
}

// GetGameInfo is a free data retrieval call binding the contract method 0x47e1d550.
//
// Solidity: function getGameInfo(uint256 gameId) view returns((uint256,uint8,uint8,uint256,uint256,uint256,uint256,uint256,uint256,uint256))
func (_Gamecontract *GamecontractCallerSession) GetGameInfo(gameId *big.Int) (IGameContractGameInfo, error) {
	return _Gamecontract.Contract.GetGameInfo(&_Gamecontract.CallOpts, gameId)
}

// GetGameResult is a free data retrieval call binding the contract method 0xed2df26d.
//
// Solidity: function getGameResult(uint256 gameId) view returns(uint8)
func (_Gamecontract *GamecontractCaller) GetGameResult(opts *bind.CallOpts, gameId *big.Int) (uint8, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getGameResult", gameId)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetGameResult is a free data retrieval call binding the contract method 0xed2df26d.
//
// Solidity: function getGameResult(uint256 gameId) view returns(uint8)
func (_Gamecontract *GamecontractSession) GetGameResult(gameId *big.Int) (uint8, error) {
	return _Gamecontract.Contract.GetGameResult(&_Gamecontract.CallOpts, gameId)
}

// GetGameResult is a free data retrieval call binding the contract method 0xed2df26d.
//
// Solidity: function getGameResult(uint256 gameId) view returns(uint8)
func (_Gamecontract *GamecontractCallerSession) GetGameResult(gameId *big.Int) (uint8, error) {
	return _Gamecontract.Contract.GetGameResult(&_Gamecontract.CallOpts, gameId)
}

// GetPlayerInfo is a free data retrieval call binding the contract method 0x3ccd10e9.
//
// Solidity: function getPlayerInfo(uint256 gameId, address player) view returns((uint8,uint256,uint256,bool))
func (_Gamecontract *GamecontractCaller) GetPlayerInfo(opts *bind.CallOpts, gameId *big.Int, player common.Address) (IGameContractPlayerInfo, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getPlayerInfo", gameId, player)

	if err != nil {
		return *new(IGameContractPlayerInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IGameContractPlayerInfo)).(*IGameContractPlayerInfo)

	return out0, err

}

// GetPlayerInfo is a free data retrieval call binding the contract method 0x3ccd10e9.
//
// Solidity: function getPlayerInfo(uint256 gameId, address player) view returns((uint8,uint256,uint256,bool))
func (_Gamecontract *GamecontractSession) GetPlayerInfo(gameId *big.Int, player common.Address) (IGameContractPlayerInfo, error) {
	return _Gamecontract.Contract.GetPlayerInfo(&_Gamecontract.CallOpts, gameId, player)
}

// GetPlayerInfo is a free data retrieval call binding the contract method 0x3ccd10e9.
//
// Solidity: function getPlayerInfo(uint256 gameId, address player) view returns((uint8,uint256,uint256,bool))
func (_Gamecontract *GamecontractCallerSession) GetPlayerInfo(gameId *big.Int, player common.Address) (IGameContractPlayerInfo, error) {
	return _Gamecontract.Contract.GetPlayerInfo(&_Gamecontract.CallOpts, gameId, player)
}

// GetPlayerMoveCount is a free data retrieval call binding the contract method 0x3305a7a9.
//
// Solidity: function getPlayerMoveCount(uint256 gameId, address player, uint32 chainId) view returns(uint256)
func (_Gamecontract *GamecontractCaller) GetPlayerMoveCount(opts *bind.CallOpts, gameId *big.Int, player common.Address, chainId uint32) (*big.Int, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getPlayerMoveCount", gameId, player, chainId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlayerMoveCount is a free data retrieval call binding the contract method 0x3305a7a9.
//
// Solidity: function getPlayerMoveCount(uint256 gameId, address player, uint32 chainId) view returns(uint256)
func (_Gamecontract *GamecontractSession) GetPlayerMoveCount(gameId *big.Int, player common.Address, chainId uint32) (*big.Int, error) {
	return _Gamecontract.Contract.GetPlayerMoveCount(&_Gamecontract.CallOpts, gameId, player, chainId)
}

// GetPlayerMoveCount is a free data retrieval call binding the contract method 0x3305a7a9.
//
// Solidity: function getPlayerMoveCount(uint256 gameId, address player, uint32 chainId) view returns(uint256)
func (_Gamecontract *GamecontractCallerSession) GetPlayerMoveCount(gameId *big.Int, player common.Address, chainId uint32) (*big.Int, error) {
	return _Gamecontract.Contract.GetPlayerMoveCount(&_Gamecontract.CallOpts, gameId, player, chainId)
}

// GetPlayerTeam is a free data retrieval call binding the contract method 0x4294b984.
//
// Solidity: function getPlayerTeam(uint256 gameId, address player) view returns(uint8)
func (_Gamecontract *GamecontractCaller) GetPlayerTeam(opts *bind.CallOpts, gameId *big.Int, player common.Address) (uint8, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getPlayerTeam", gameId, player)

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetPlayerTeam is a free data retrieval call binding the contract method 0x4294b984.
//
// Solidity: function getPlayerTeam(uint256 gameId, address player) view returns(uint8)
func (_Gamecontract *GamecontractSession) GetPlayerTeam(gameId *big.Int, player common.Address) (uint8, error) {
	return _Gamecontract.Contract.GetPlayerTeam(&_Gamecontract.CallOpts, gameId, player)
}

// GetPlayerTeam is a free data retrieval call binding the contract method 0x4294b984.
//
// Solidity: function getPlayerTeam(uint256 gameId, address player) view returns(uint8)
func (_Gamecontract *GamecontractCallerSession) GetPlayerTeam(gameId *big.Int, player common.Address) (uint8, error) {
	return _Gamecontract.Contract.GetPlayerTeam(&_Gamecontract.CallOpts, gameId, player)
}

// GetPlayerTotalMoveCount is a free data retrieval call binding the contract method 0x05421720.
//
// Solidity: function getPlayerTotalMoveCount(uint256 gameId, address player) view returns(uint256)
func (_Gamecontract *GamecontractCaller) GetPlayerTotalMoveCount(opts *bind.CallOpts, gameId *big.Int, player common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getPlayerTotalMoveCount", gameId, player)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlayerTotalMoveCount is a free data retrieval call binding the contract method 0x05421720.
//
// Solidity: function getPlayerTotalMoveCount(uint256 gameId, address player) view returns(uint256)
func (_Gamecontract *GamecontractSession) GetPlayerTotalMoveCount(gameId *big.Int, player common.Address) (*big.Int, error) {
	return _Gamecontract.Contract.GetPlayerTotalMoveCount(&_Gamecontract.CallOpts, gameId, player)
}

// GetPlayerTotalMoveCount is a free data retrieval call binding the contract method 0x05421720.
//
// Solidity: function getPlayerTotalMoveCount(uint256 gameId, address player) view returns(uint256)
func (_Gamecontract *GamecontractCallerSession) GetPlayerTotalMoveCount(gameId *big.Int, player common.Address) (*big.Int, error) {
	return _Gamecontract.Contract.GetPlayerTotalMoveCount(&_Gamecontract.CallOpts, gameId, player)
}

// HasPlayerJoined is a free data retrieval call binding the contract method 0x07286e2f.
//
// Solidity: function hasPlayerJoined(uint256 gameId, address player) view returns(bool)
func (_Gamecontract *GamecontractCaller) HasPlayerJoined(opts *bind.CallOpts, gameId *big.Int, player common.Address) (bool, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "hasPlayerJoined", gameId, player)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// HasPlayerJoined is a free data retrieval call binding the contract method 0x07286e2f.
//
// Solidity: function hasPlayerJoined(uint256 gameId, address player) view returns(bool)
func (_Gamecontract *GamecontractSession) HasPlayerJoined(gameId *big.Int, player common.Address) (bool, error) {
	return _Gamecontract.Contract.HasPlayerJoined(&_Gamecontract.CallOpts, gameId, player)
}

// HasPlayerJoined is a free data retrieval call binding the contract method 0x07286e2f.
//
// Solidity: function hasPlayerJoined(uint256 gameId, address player) view returns(bool)
func (_Gamecontract *GamecontractCallerSession) HasPlayerJoined(gameId *big.Int, player common.Address) (bool, error) {
	return _Gamecontract.Contract.HasPlayerJoined(&_Gamecontract.CallOpts, gameId, player)
}

// IsGameActive is a free data retrieval call binding the contract method 0x1a60a4e5.
//
// Solidity: function isGameActive(uint256 gameId) view returns(bool)
func (_Gamecontract *GamecontractCaller) IsGameActive(opts *bind.CallOpts, gameId *big.Int) (bool, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "isGameActive", gameId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsGameActive is a free data retrieval call binding the contract method 0x1a60a4e5.
//
// Solidity: function isGameActive(uint256 gameId) view returns(bool)
func (_Gamecontract *GamecontractSession) IsGameActive(gameId *big.Int) (bool, error) {
	return _Gamecontract.Contract.IsGameActive(&_Gamecontract.CallOpts, gameId)
}

// IsGameActive is a free data retrieval call binding the contract method 0x1a60a4e5.
//
// Solidity: function isGameActive(uint256 gameId) view returns(bool)
func (_Gamecontract *GamecontractCallerSession) IsGameActive(gameId *big.Int) (bool, error) {
	return _Gamecontract.Contract.IsGameActive(&_Gamecontract.CallOpts, gameId)
}

// Moves is a free data retrieval call binding the contract method 0x2b497f85.
//
// Solidity: function moves(uint256 gameId, uint32 chainId, address player) view returns(uint256 moveCount)
func (_Gamecontract *GamecontractCaller) Moves(opts *bind.CallOpts, gameId *big.Int, chainId uint32, player common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "moves", gameId, chainId, player)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// Moves is a free data retrieval call binding the contract method 0x2b497f85.
//
// Solidity: function moves(uint256 gameId, uint32 chainId, address player) view returns(uint256 moveCount)
func (_Gamecontract *GamecontractSession) Moves(gameId *big.Int, chainId uint32, player common.Address) (*big.Int, error) {
	return _Gamecontract.Contract.Moves(&_Gamecontract.CallOpts, gameId, chainId, player)
}

// Moves is a free data retrieval call binding the contract method 0x2b497f85.
//
// Solidity: function moves(uint256 gameId, uint32 chainId, address player) view returns(uint256 moveCount)
func (_Gamecontract *GamecontractCallerSession) Moves(gameId *big.Int, chainId uint32, player common.Address) (*big.Int, error) {
	return _Gamecontract.Contract.Moves(&_Gamecontract.CallOpts, gameId, chainId, player)
}

// Players is a free data retrieval call binding the contract method 0xc0f80921.
//
// Solidity: function players(uint256 gameId, address player) view returns(uint8 team, uint256 totalStakes, uint256 moveCount, bool hasJoined)
func (_Gamecontract *GamecontractCaller) Players(opts *bind.CallOpts, gameId *big.Int, player common.Address) (struct {
	Team        uint8
	TotalStakes *big.Int
	MoveCount   *big.Int
	HasJoined   bool
}, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "players", gameId, player)

	outstruct := new(struct {
		Team        uint8
		TotalStakes *big.Int
		MoveCount   *big.Int
		HasJoined   bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Team = *abi.ConvertType(out[0], new(uint8)).(*uint8)
	outstruct.TotalStakes = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.MoveCount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.HasJoined = *abi.ConvertType(out[3], new(bool)).(*bool)

	return *outstruct, err

}

// Players is a free data retrieval call binding the contract method 0xc0f80921.
//
// Solidity: function players(uint256 gameId, address player) view returns(uint8 team, uint256 totalStakes, uint256 moveCount, bool hasJoined)
func (_Gamecontract *GamecontractSession) Players(gameId *big.Int, player common.Address) (struct {
	Team        uint8
	TotalStakes *big.Int
	MoveCount   *big.Int
	HasJoined   bool
}, error) {
	return _Gamecontract.Contract.Players(&_Gamecontract.CallOpts, gameId, player)
}

// Players is a free data retrieval call binding the contract method 0xc0f80921.
//
// Solidity: function players(uint256 gameId, address player) view returns(uint8 team, uint256 totalStakes, uint256 moveCount, bool hasJoined)
func (_Gamecontract *GamecontractCallerSession) Players(gameId *big.Int, player common.Address) (struct {
	Team        uint8
	TotalStakes *big.Int
	MoveCount   *big.Int
	HasJoined   bool
}, error) {
	return _Gamecontract.Contract.Players(&_Gamecontract.CallOpts, gameId, player)
}

// CreateGame is a paid mutator transaction binding the contract method 0x60104cef.
//
// Solidity: function createGame(uint256 gameId, uint256 fixedStakeAmount) returns()
func (_Gamecontract *GamecontractTransactor) CreateGame(opts *bind.TransactOpts, gameId *big.Int, fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Gamecontract.contract.Transact(opts, "createGame", gameId, fixedStakeAmount)
}

// CreateGame is a paid mutator transaction binding the contract method 0x60104cef.
//
// Solidity: function createGame(uint256 gameId, uint256 fixedStakeAmount) returns()
func (_Gamecontract *GamecontractSession) CreateGame(gameId *big.Int, fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Gamecontract.Contract.CreateGame(&_Gamecontract.TransactOpts, gameId, fixedStakeAmount)
}

// CreateGame is a paid mutator transaction binding the contract method 0x60104cef.
//
// Solidity: function createGame(uint256 gameId, uint256 fixedStakeAmount) returns()
func (_Gamecontract *GamecontractTransactorSession) CreateGame(gameId *big.Int, fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Gamecontract.Contract.CreateGame(&_Gamecontract.TransactOpts, gameId, fixedStakeAmount)
}

// EndGame is a paid mutator transaction binding the contract method 0xc5449531.
//
// Solidity: function endGame(uint256 gameId, uint8 result) returns()
func (_Gamecontract *GamecontractTransactor) EndGame(opts *bind.TransactOpts, gameId *big.Int, result uint8) (*types.Transaction, error) {
	return _Gamecontract.contract.Transact(opts, "endGame", gameId, result)
}

// EndGame is a paid mutator transaction binding the contract method 0xc5449531.
//
// Solidity: function endGame(uint256 gameId, uint8 result) returns()
func (_Gamecontract *GamecontractSession) EndGame(gameId *big.Int, result uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.EndGame(&_Gamecontract.TransactOpts, gameId, result)
}

// EndGame is a paid mutator transaction binding the contract method 0xc5449531.
//
// Solidity: function endGame(uint256 gameId, uint8 result) returns()
func (_Gamecontract *GamecontractTransactorSession) EndGame(gameId *big.Int, result uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.EndGame(&_Gamecontract.TransactOpts, gameId, result)
}

// JoinTeam is a paid mutator transaction binding the contract method 0xc3d4bd8d.
//
// Solidity: function joinTeam(uint256 gameId, uint8 team) returns()
func (_Gamecontract *GamecontractTransactor) JoinTeam(opts *bind.TransactOpts, gameId *big.Int, team uint8) (*types.Transaction, error) {
	return _Gamecontract.contract.Transact(opts, "joinTeam", gameId, team)
}

// JoinTeam is a paid mutator transaction binding the contract method 0xc3d4bd8d.
//
// Solidity: function joinTeam(uint256 gameId, uint8 team) returns()
func (_Gamecontract *GamecontractSession) JoinTeam(gameId *big.Int, team uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.JoinTeam(&_Gamecontract.TransactOpts, gameId, team)
}

// JoinTeam is a paid mutator transaction binding the contract method 0xc3d4bd8d.
//
// Solidity: function joinTeam(uint256 gameId, uint8 team) returns()
func (_Gamecontract *GamecontractTransactorSession) JoinTeam(gameId *big.Int, team uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.JoinTeam(&_Gamecontract.TransactOpts, gameId, team)
}

// RecordMove is a paid mutator transaction binding the contract method 0xfeee346f.
//
// Solidity: function recordMove(uint256 gameId, address player, uint32 chainId) returns()
func (_Gamecontract *GamecontractTransactor) RecordMove(opts *bind.TransactOpts, gameId *big.Int, player common.Address, chainId uint32) (*types.Transaction, error) {
	return _Gamecontract.contract.Transact(opts, "recordMove", gameId, player, chainId)
}

// RecordMove is a paid mutator transaction binding the contract method 0xfeee346f.
//
// Solidity: function recordMove(uint256 gameId, address player, uint32 chainId) returns()
func (_Gamecontract *GamecontractSession) RecordMove(gameId *big.Int, player common.Address, chainId uint32) (*types.Transaction, error) {
	return _Gamecontract.Contract.RecordMove(&_Gamecontract.TransactOpts, gameId, player, chainId)
}

// RecordMove is a paid mutator transaction binding the contract method 0xfeee346f.
//
// Solidity: function recordMove(uint256 gameId, address player, uint32 chainId) returns()
func (_Gamecontract *GamecontractTransactorSession) RecordMove(gameId *big.Int, player common.Address, chainId uint32) (*types.Transaction, error) {
	return _Gamecontract.Contract.RecordMove(&_Gamecontract.TransactOpts, gameId, player, chainId)
}

// GamecontractGameCreatedIterator is returned from FilterGameCreated and is used to iterate over the raw logs and unpacked data for GameCreated events raised by the Gamecontract contract.
type GamecontractGameCreatedIterator struct {
	Event *GamecontractGameCreated // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GamecontractGameCreatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GamecontractGameCreated)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GamecontractGameCreated)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GamecontractGameCreatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GamecontractGameCreatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GamecontractGameCreated represents a GameCreated event raised by the Gamecontract contract.
type GamecontractGameCreated struct {
	GameId           *big.Int
	FixedStakeAmount *big.Int
	CreatedAt        *big.Int
	Raw              types.Log // Blockchain specific contextual infos
}

// FilterGameCreated is a free log retrieval operation binding the contract event 0xbd3f84a55fdc44f970adc6ccbaa93e78d7da5f3fc37c4b7da578d47c9ad2e46a.
//
// Solidity: event GameCreated(uint256 indexed gameId, uint256 fixedStakeAmount, uint256 createdAt)
func (_Gamecontract *GamecontractFilterer) FilterGameCreated(opts *bind.FilterOpts, gameId []*big.Int) (*GamecontractGameCreatedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Gamecontract.contract.FilterLogs(opts, "GameCreated", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &GamecontractGameCreatedIterator{contract: _Gamecontract.contract, event: "GameCreated", logs: logs, sub: sub}, nil
}

// WatchGameCreated is a free log subscription operation binding the contract event 0xbd3f84a55fdc44f970adc6ccbaa93e78d7da5f3fc37c4b7da578d47c9ad2e46a.
//
// Solidity: event GameCreated(uint256 indexed gameId, uint256 fixedStakeAmount, uint256 createdAt)
func (_Gamecontract *GamecontractFilterer) WatchGameCreated(opts *bind.WatchOpts, sink chan<- *GamecontractGameCreated, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Gamecontract.contract.WatchLogs(opts, "GameCreated", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GamecontractGameCreated)
				if err := _Gamecontract.contract.UnpackLog(event, "GameCreated", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameCreated is a log parse operation binding the contract event 0xbd3f84a55fdc44f970adc6ccbaa93e78d7da5f3fc37c4b7da578d47c9ad2e46a.
//
// Solidity: event GameCreated(uint256 indexed gameId, uint256 fixedStakeAmount, uint256 createdAt)
func (_Gamecontract *GamecontractFilterer) ParseGameCreated(log types.Log) (*GamecontractGameCreated, error) {
	event := new(GamecontractGameCreated)
	if err := _Gamecontract.contract.UnpackLog(event, "GameCreated", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GamecontractGameEndedIterator is returned from FilterGameEnded and is used to iterate over the raw logs and unpacked data for GameEnded events raised by the Gamecontract contract.
type GamecontractGameEndedIterator struct {
	Event *GamecontractGameEnded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GamecontractGameEndedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GamecontractGameEnded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GamecontractGameEnded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GamecontractGameEndedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GamecontractGameEndedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GamecontractGameEnded represents a GameEnded event raised by the Gamecontract contract.
type GamecontractGameEnded struct {
	GameId  *big.Int
	Result  uint8
	EndedAt *big.Int
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterGameEnded is a free log retrieval operation binding the contract event 0xad5c2fcf2f9497233e36f5fa5a21487acf87e422dd4d6535fd81cd01b37925d9.
//
// Solidity: event GameEnded(uint256 indexed gameId, uint8 result, uint256 endedAt)
func (_Gamecontract *GamecontractFilterer) FilterGameEnded(opts *bind.FilterOpts, gameId []*big.Int) (*GamecontractGameEndedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Gamecontract.contract.FilterLogs(opts, "GameEnded", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &GamecontractGameEndedIterator{contract: _Gamecontract.contract, event: "GameEnded", logs: logs, sub: sub}, nil
}

// WatchGameEnded is a free log subscription operation binding the contract event 0xad5c2fcf2f9497233e36f5fa5a21487acf87e422dd4d6535fd81cd01b37925d9.
//
// Solidity: event GameEnded(uint256 indexed gameId, uint8 result, uint256 endedAt)
func (_Gamecontract *GamecontractFilterer) WatchGameEnded(opts *bind.WatchOpts, sink chan<- *GamecontractGameEnded, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Gamecontract.contract.WatchLogs(opts, "GameEnded", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GamecontractGameEnded)
				if err := _Gamecontract.contract.UnpackLog(event, "GameEnded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseGameEnded is a log parse operation binding the contract event 0xad5c2fcf2f9497233e36f5fa5a21487acf87e422dd4d6535fd81cd01b37925d9.
//
// Solidity: event GameEnded(uint256 indexed gameId, uint8 result, uint256 endedAt)
func (_Gamecontract *GamecontractFilterer) ParseGameEnded(log types.Log) (*GamecontractGameEnded, error) {
	event := new(GamecontractGameEnded)
	if err := _Gamecontract.contract.UnpackLog(event, "GameEnded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GamecontractMoveRecordedIterator is returned from FilterMoveRecorded and is used to iterate over the raw logs and unpacked data for MoveRecorded events raised by the Gamecontract contract.
type GamecontractMoveRecordedIterator struct {
	Event *GamecontractMoveRecorded // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GamecontractMoveRecordedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GamecontractMoveRecorded)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GamecontractMoveRecorded)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GamecontractMoveRecordedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GamecontractMoveRecordedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GamecontractMoveRecorded represents a MoveRecorded event raised by the Gamecontract contract.
type GamecontractMoveRecorded struct {
	GameId       *big.Int
	Player       common.Address
	ChainId      uint32
	NewMoveCount *big.Int
	Raw          types.Log // Blockchain specific contextual infos
}

// FilterMoveRecorded is a free log retrieval operation binding the contract event 0x2e5e5932470e4421499604d32a27d62dc99709b56faaea15de88d17fe1603aac.
//
// Solidity: event MoveRecorded(uint256 indexed gameId, address indexed player, uint32 chainId, uint256 newMoveCount)
func (_Gamecontract *GamecontractFilterer) FilterMoveRecorded(opts *bind.FilterOpts, gameId []*big.Int, player []common.Address) (*GamecontractMoveRecordedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Gamecontract.contract.FilterLogs(opts, "MoveRecorded", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return &GamecontractMoveRecordedIterator{contract: _Gamecontract.contract, event: "MoveRecorded", logs: logs, sub: sub}, nil
}

// WatchMoveRecorded is a free log subscription operation binding the contract event 0x2e5e5932470e4421499604d32a27d62dc99709b56faaea15de88d17fe1603aac.
//
// Solidity: event MoveRecorded(uint256 indexed gameId, address indexed player, uint32 chainId, uint256 newMoveCount)
func (_Gamecontract *GamecontractFilterer) WatchMoveRecorded(opts *bind.WatchOpts, sink chan<- *GamecontractMoveRecorded, gameId []*big.Int, player []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Gamecontract.contract.WatchLogs(opts, "MoveRecorded", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GamecontractMoveRecorded)
				if err := _Gamecontract.contract.UnpackLog(event, "MoveRecorded", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParseMoveRecorded is a log parse operation binding the contract event 0x2e5e5932470e4421499604d32a27d62dc99709b56faaea15de88d17fe1603aac.
//
// Solidity: event MoveRecorded(uint256 indexed gameId, address indexed player, uint32 chainId, uint256 newMoveCount)
func (_Gamecontract *GamecontractFilterer) ParseMoveRecorded(log types.Log) (*GamecontractMoveRecorded, error) {
	event := new(GamecontractMoveRecorded)
	if err := _Gamecontract.contract.UnpackLog(event, "MoveRecorded", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GamecontractPlayerJoinedTeamIterator is returned from FilterPlayerJoinedTeam and is used to iterate over the raw logs and unpacked data for PlayerJoinedTeam events raised by the Gamecontract contract.
type GamecontractPlayerJoinedTeamIterator struct {
	Event *GamecontractPlayerJoinedTeam // Event containing the contract specifics and raw log

	contract *bind.BoundContract // Generic contract to use for unpacking event data
	event    string              // Event name to use for unpacking event data

	logs chan types.Log        // Log channel receiving the found contract events
	sub  ethereum.Subscription // Subscription for errors, completion and termination
	done bool                  // Whether the subscription completed delivering logs
	fail error                 // Occurred error to stop iteration
}

// Next advances the iterator to the subsequent event, returning whether there
// are any more events found. In case of a retrieval or parsing error, false is
// returned and Error() can be queried for the exact failure.
func (it *GamecontractPlayerJoinedTeamIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GamecontractPlayerJoinedTeam)
			if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
				it.fail = err
				return false
			}
			it.Event.Raw = log
			return true

		default:
			return false
		}
	}
	// Iterator still in progress, wait for either a data or an error event
	select {
	case log := <-it.logs:
		it.Event = new(GamecontractPlayerJoinedTeam)
		if err := it.contract.UnpackLog(it.Event, it.event, log); err != nil {
			it.fail = err
			return false
		}
		it.Event.Raw = log
		return true

	case err := <-it.sub.Err():
		it.done = true
		it.fail = err
		return it.Next()
	}
}

// Error returns any retrieval or parsing error occurred during filtering.
func (it *GamecontractPlayerJoinedTeamIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GamecontractPlayerJoinedTeamIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GamecontractPlayerJoinedTeam represents a PlayerJoinedTeam event raised by the Gamecontract contract.
type GamecontractPlayerJoinedTeam struct {
	GameId *big.Int
	Player common.Address
	Team   uint8
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterPlayerJoinedTeam is a free log retrieval operation binding the contract event 0x345d86e28bee17ad4bc57f12601ea9a4817a1afaac78868cd4f7af041a0638fa.
//
// Solidity: event PlayerJoinedTeam(uint256 indexed gameId, address indexed player, uint8 team)
func (_Gamecontract *GamecontractFilterer) FilterPlayerJoinedTeam(opts *bind.FilterOpts, gameId []*big.Int, player []common.Address) (*GamecontractPlayerJoinedTeamIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Gamecontract.contract.FilterLogs(opts, "PlayerJoinedTeam", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return &GamecontractPlayerJoinedTeamIterator{contract: _Gamecontract.contract, event: "PlayerJoinedTeam", logs: logs, sub: sub}, nil
}

// WatchPlayerJoinedTeam is a free log subscription operation binding the contract event 0x345d86e28bee17ad4bc57f12601ea9a4817a1afaac78868cd4f7af041a0638fa.
//
// Solidity: event PlayerJoinedTeam(uint256 indexed gameId, address indexed player, uint8 team)
func (_Gamecontract *GamecontractFilterer) WatchPlayerJoinedTeam(opts *bind.WatchOpts, sink chan<- *GamecontractPlayerJoinedTeam, gameId []*big.Int, player []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Gamecontract.contract.WatchLogs(opts, "PlayerJoinedTeam", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GamecontractPlayerJoinedTeam)
				if err := _Gamecontract.contract.UnpackLog(event, "PlayerJoinedTeam", log); err != nil {
					return err
				}
				event.Raw = log

				select {
				case sink <- event:
				case err := <-sub.Err():
					return err
				case <-quit:
					return nil
				}
			case err := <-sub.Err():
				return err
			case <-quit:
				return nil
			}
		}
	}), nil
}

// ParsePlayerJoinedTeam is a log parse operation binding the contract event 0x345d86e28bee17ad4bc57f12601ea9a4817a1afaac78868cd4f7af041a0638fa.
//
// Solidity: event PlayerJoinedTeam(uint256 indexed gameId, address indexed player, uint8 team)
func (_Gamecontract *GamecontractFilterer) ParsePlayerJoinedTeam(log types.Log) (*GamecontractPlayerJoinedTeam, error) {
	event := new(GamecontractPlayerJoinedTeam)
	if err := _Gamecontract.contract.UnpackLog(event, "PlayerJoinedTeam", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
