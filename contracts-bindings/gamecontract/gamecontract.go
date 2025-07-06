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
	GameId *big.Int
	State  uint8
	Result uint8
}

// GamecontractMetaData contains all meta data concerning the Gamecontract contract.
var GamecontractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_factory\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"FACTORY\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"addVote\",\"inputs\":[{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"chainId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"team\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.Team\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endGame\",\"inputs\":[{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameResult\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"getGameId\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameInfo\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIGameContract.GameInfo\",\"components\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"state\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameState\"},{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameResult\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameResult\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameResult\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameStatus\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint8\",\"internalType\":\"enumIGameContract.GameState\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayerVoteCounts\",\"inputs\":[],\"outputs\":[{\"name\":\"players\",\"type\":\"address[]\",\"internalType\":\"address[]\"},{\"name\":\"voteCounts\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"},{\"name\":\"chainIds\",\"type\":\"uint32[]\",\"internalType\":\"uint32[]\"},{\"name\":\"teams\",\"type\":\"uint256[]\",\"internalType\":\"uint256[]\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"initialize\",\"inputs\":[{\"name\":\"_gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"_fixedStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"userExists\",\"inputs\":[{\"name\":\"userAddress\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"event\",\"name\":\"GameCreated\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"fixedStakeAmount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"createdAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"GameFinished\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"result\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIGameContract.GameResult\"},{\"name\":\"finishedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"Vote\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"team\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIGameContract.Team\"},{\"name\":\"chainId\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"}],\"anonymous\":false}]",
	Bin: "0x60a03461008d57601f61097b38819003918201601f19168301916001600160401b038311848410176100915780849260209460405283398101031261008d57516001600160a01b038116810361008d5760805261020061ffff1960025416176002556040516108d590816100a6823960805181818160ca0152818161019401528181610364015261061d0152f35b5f80fd5b634e487b7160e01b5f52604160045260245ffdfe6080806040526004361015610012575f80fd5b5f3560e01c9081630e666e49146106e8575080631746bd1b1461064c5780632dd3100014610608578063382396ee146105df5780636ec1482814610430578063863375ac146103415780638f66d4ec14610157578063b3fb14ad1461012d578063c0bd8351146101115763e4a301161461008a575f80fd5b3461010d57604036600319011261010d576004357fbd3f84a55fdc44f970adc6ccbaa93e78d7da5f3fc37c4b7da578d47c9ad2e46a60406024356100f8337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614610847565b835f55806001558151908152426020820152a2005b5f80fd5b3461010d575f36600319011261010d5760205f54604051908152f35b3461010d575f36600319011261010d57602060ff60025460081c166101556040518092610741565bf35b3461010d57606036600319011261010d57610170610721565b60243563ffffffff811680910361010d5760443591600283101561010d576101c2337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614610847565b6001600160a01b03165f818152600360205260409020600101549092906102e3576101eb610781565b8281526002602082019160018352604081019261020785610737565b848452865f52600360205263ffffffff60405f2092511663ffffffff1983541617825551600182015501905161023c81610737565b61024581610737565b60ff80198354169116179055600454680100000000000000008110156102cf577fca41c755a10d5981d66a7053c39564f2f68151dbd0ba5295fea79ed7ec1717869161029982600160409401600455610807565b81546001600160a01b0360039290921b91821b19169087901b1790555b5f54938251916102c581610737565b82526020820152a3005b634e487b7160e01b5f52604160045260245ffd5b825f526003602052600160405f2001908154905f19821461032d577fca41c755a10d5981d66a7053c39564f2f68151dbd0ba5295fea79ed7ec1717869260016040930190556102b6565b634e487b7160e01b5f52601160045260245ffd5b3461010d57602036600319011261010d57600435600381101561010d57610392337f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031614610847565b60025460ff81166103a281610737565b6103f65760019061ff008360081b169061ffff191617176002557e21feea4608e6756f982d5bd71c09645176b4865cf6bd600fcc4aa990f0240460405f54926103ed82518092610741565b426020820152a2005b60405162461bcd60e51b815260206004820152601260248201527147616d65206973206e6f742061637469766560701b6044820152606490fd5b3461010d575f36600319011261010d5760045461045461044f826107c7565b6107a1565b818152610460826107c7565b602082019290601f1901368437610476816107df565b9261048361044f836107c7565b90828252610490836107c7565b602083019590601f19013687376104a6846107df565b935f5b818110610555575050604051946080860190608087525180915260a0860192905f5b81811061053657505050816104e89186602094038488015261074e565b91848303604086015251918281520193905f5b81811061051a5784806105168887838203606085015261074e565b0390f35b825163ffffffff168652602095860195909201916001016104fb565b82516001600160a01b03168552602094850194909201916001016104cb565b80610561600192610807565b838060a01b0391549060031b1c168061057a838b610833565b52805f5260036020528260405f2001546105948387610833565b52805f52600360205263ffffffff60405f2054166105b28389610833565b525f52600360205260ff600260405f200154166105ce81610737565b6105d88289610833565b52016104a9565b3461010d575f36600319011261010d57602060ff600254166040519061060481610737565b8152f35b3461010d575f36600319011261010d576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b3461010d575f36600319011261010d575f6040610667610781565b82815282602082015201525f546002549060ff8083169260081c1661068a610781565b918252602082019261069b81610737565b83526040820160038210156106d45760609361015592825260405193518452516106c481610737565b6020840152516040830190610741565b634e487b7160e01b5f52602160045260245ffd5b3461010d57602036600319011261010d576020906001600160a01b0361070c610721565b165f5260038252600160405f20015415158152f35b600435906001600160a01b038216820361010d57565b600211156106d457565b9060038210156106d45752565b90602080835192838152019201905f5b81811061076b5750505090565b825184526020938401939092019160010161075e565b604051906060820182811067ffffffffffffffff8211176102cf57604052565b6040519190601f01601f1916820167ffffffffffffffff8111838210176102cf57604052565b67ffffffffffffffff81116102cf5760051b60200190565b906107ec61044f836107c7565b82815280926107fd601f19916107c7565b0190602036910137565b60045481101561081f5760045f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b805182101561081f5760209160051b010190565b1561084e57565b60405162461bcd60e51b815260206004820152602360248201527f4f6e6c7920666163746f72792063616e2063616c6c20746869732066756e637460448201526234b7b760e91b6064820152608490fdfea2646970667358221220ebc644537d00dfe39d8644bd39019f509d671b62b94410eaa417f197bda5268864736f6c634300081e0033",
}

// GamecontractABI is the input ABI used to generate the binding from.
// Deprecated: Use GamecontractMetaData.ABI instead.
var GamecontractABI = GamecontractMetaData.ABI

// GamecontractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use GamecontractMetaData.Bin instead.
var GamecontractBin = GamecontractMetaData.Bin

// DeployGamecontract deploys a new Ethereum contract, binding an instance of Gamecontract to it.
func DeployGamecontract(auth *bind.TransactOpts, backend bind.ContractBackend, _factory common.Address) (common.Address, *types.Transaction, *Gamecontract, error) {
	parsed, err := GamecontractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(GamecontractBin), backend, _factory)
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

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_Gamecontract *GamecontractCaller) FACTORY(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "FACTORY")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_Gamecontract *GamecontractSession) FACTORY() (common.Address, error) {
	return _Gamecontract.Contract.FACTORY(&_Gamecontract.CallOpts)
}

// FACTORY is a free data retrieval call binding the contract method 0x2dd31000.
//
// Solidity: function FACTORY() view returns(address)
func (_Gamecontract *GamecontractCallerSession) FACTORY() (common.Address, error) {
	return _Gamecontract.Contract.FACTORY(&_Gamecontract.CallOpts)
}

// GetGameId is a free data retrieval call binding the contract method 0xc0bd8351.
//
// Solidity: function getGameId() view returns(uint256)
func (_Gamecontract *GamecontractCaller) GetGameId(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getGameId")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetGameId is a free data retrieval call binding the contract method 0xc0bd8351.
//
// Solidity: function getGameId() view returns(uint256)
func (_Gamecontract *GamecontractSession) GetGameId() (*big.Int, error) {
	return _Gamecontract.Contract.GetGameId(&_Gamecontract.CallOpts)
}

// GetGameId is a free data retrieval call binding the contract method 0xc0bd8351.
//
// Solidity: function getGameId() view returns(uint256)
func (_Gamecontract *GamecontractCallerSession) GetGameId() (*big.Int, error) {
	return _Gamecontract.Contract.GetGameId(&_Gamecontract.CallOpts)
}

// GetGameInfo is a free data retrieval call binding the contract method 0x1746bd1b.
//
// Solidity: function getGameInfo() view returns((uint256,uint8,uint8))
func (_Gamecontract *GamecontractCaller) GetGameInfo(opts *bind.CallOpts) (IGameContractGameInfo, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getGameInfo")

	if err != nil {
		return *new(IGameContractGameInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IGameContractGameInfo)).(*IGameContractGameInfo)

	return out0, err

}

// GetGameInfo is a free data retrieval call binding the contract method 0x1746bd1b.
//
// Solidity: function getGameInfo() view returns((uint256,uint8,uint8))
func (_Gamecontract *GamecontractSession) GetGameInfo() (IGameContractGameInfo, error) {
	return _Gamecontract.Contract.GetGameInfo(&_Gamecontract.CallOpts)
}

// GetGameInfo is a free data retrieval call binding the contract method 0x1746bd1b.
//
// Solidity: function getGameInfo() view returns((uint256,uint8,uint8))
func (_Gamecontract *GamecontractCallerSession) GetGameInfo() (IGameContractGameInfo, error) {
	return _Gamecontract.Contract.GetGameInfo(&_Gamecontract.CallOpts)
}

// GetGameResult is a free data retrieval call binding the contract method 0xb3fb14ad.
//
// Solidity: function getGameResult() view returns(uint8)
func (_Gamecontract *GamecontractCaller) GetGameResult(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getGameResult")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetGameResult is a free data retrieval call binding the contract method 0xb3fb14ad.
//
// Solidity: function getGameResult() view returns(uint8)
func (_Gamecontract *GamecontractSession) GetGameResult() (uint8, error) {
	return _Gamecontract.Contract.GetGameResult(&_Gamecontract.CallOpts)
}

// GetGameResult is a free data retrieval call binding the contract method 0xb3fb14ad.
//
// Solidity: function getGameResult() view returns(uint8)
func (_Gamecontract *GamecontractCallerSession) GetGameResult() (uint8, error) {
	return _Gamecontract.Contract.GetGameResult(&_Gamecontract.CallOpts)
}

// GetGameStatus is a free data retrieval call binding the contract method 0x382396ee.
//
// Solidity: function getGameStatus() view returns(uint8)
func (_Gamecontract *GamecontractCaller) GetGameStatus(opts *bind.CallOpts) (uint8, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getGameStatus")

	if err != nil {
		return *new(uint8), err
	}

	out0 := *abi.ConvertType(out[0], new(uint8)).(*uint8)

	return out0, err

}

// GetGameStatus is a free data retrieval call binding the contract method 0x382396ee.
//
// Solidity: function getGameStatus() view returns(uint8)
func (_Gamecontract *GamecontractSession) GetGameStatus() (uint8, error) {
	return _Gamecontract.Contract.GetGameStatus(&_Gamecontract.CallOpts)
}

// GetGameStatus is a free data retrieval call binding the contract method 0x382396ee.
//
// Solidity: function getGameStatus() view returns(uint8)
func (_Gamecontract *GamecontractCallerSession) GetGameStatus() (uint8, error) {
	return _Gamecontract.Contract.GetGameStatus(&_Gamecontract.CallOpts)
}

// GetPlayerVoteCounts is a free data retrieval call binding the contract method 0x6ec14828.
//
// Solidity: function getPlayerVoteCounts() view returns(address[] players, uint256[] voteCounts, uint32[] chainIds, uint256[] teams)
func (_Gamecontract *GamecontractCaller) GetPlayerVoteCounts(opts *bind.CallOpts) (struct {
	Players    []common.Address
	VoteCounts []*big.Int
	ChainIds   []uint32
	Teams      []*big.Int
}, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "getPlayerVoteCounts")

	outstruct := new(struct {
		Players    []common.Address
		VoteCounts []*big.Int
		ChainIds   []uint32
		Teams      []*big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.Players = *abi.ConvertType(out[0], new([]common.Address)).(*[]common.Address)
	outstruct.VoteCounts = *abi.ConvertType(out[1], new([]*big.Int)).(*[]*big.Int)
	outstruct.ChainIds = *abi.ConvertType(out[2], new([]uint32)).(*[]uint32)
	outstruct.Teams = *abi.ConvertType(out[3], new([]*big.Int)).(*[]*big.Int)

	return *outstruct, err

}

// GetPlayerVoteCounts is a free data retrieval call binding the contract method 0x6ec14828.
//
// Solidity: function getPlayerVoteCounts() view returns(address[] players, uint256[] voteCounts, uint32[] chainIds, uint256[] teams)
func (_Gamecontract *GamecontractSession) GetPlayerVoteCounts() (struct {
	Players    []common.Address
	VoteCounts []*big.Int
	ChainIds   []uint32
	Teams      []*big.Int
}, error) {
	return _Gamecontract.Contract.GetPlayerVoteCounts(&_Gamecontract.CallOpts)
}

// GetPlayerVoteCounts is a free data retrieval call binding the contract method 0x6ec14828.
//
// Solidity: function getPlayerVoteCounts() view returns(address[] players, uint256[] voteCounts, uint32[] chainIds, uint256[] teams)
func (_Gamecontract *GamecontractCallerSession) GetPlayerVoteCounts() (struct {
	Players    []common.Address
	VoteCounts []*big.Int
	ChainIds   []uint32
	Teams      []*big.Int
}, error) {
	return _Gamecontract.Contract.GetPlayerVoteCounts(&_Gamecontract.CallOpts)
}

// UserExists is a free data retrieval call binding the contract method 0x0e666e49.
//
// Solidity: function userExists(address userAddress) view returns(bool)
func (_Gamecontract *GamecontractCaller) UserExists(opts *bind.CallOpts, userAddress common.Address) (bool, error) {
	var out []interface{}
	err := _Gamecontract.contract.Call(opts, &out, "userExists", userAddress)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// UserExists is a free data retrieval call binding the contract method 0x0e666e49.
//
// Solidity: function userExists(address userAddress) view returns(bool)
func (_Gamecontract *GamecontractSession) UserExists(userAddress common.Address) (bool, error) {
	return _Gamecontract.Contract.UserExists(&_Gamecontract.CallOpts, userAddress)
}

// UserExists is a free data retrieval call binding the contract method 0x0e666e49.
//
// Solidity: function userExists(address userAddress) view returns(bool)
func (_Gamecontract *GamecontractCallerSession) UserExists(userAddress common.Address) (bool, error) {
	return _Gamecontract.Contract.UserExists(&_Gamecontract.CallOpts, userAddress)
}

// AddVote is a paid mutator transaction binding the contract method 0x8f66d4ec.
//
// Solidity: function addVote(address player, uint32 chainId, uint8 team) returns()
func (_Gamecontract *GamecontractTransactor) AddVote(opts *bind.TransactOpts, player common.Address, chainId uint32, team uint8) (*types.Transaction, error) {
	return _Gamecontract.contract.Transact(opts, "addVote", player, chainId, team)
}

// AddVote is a paid mutator transaction binding the contract method 0x8f66d4ec.
//
// Solidity: function addVote(address player, uint32 chainId, uint8 team) returns()
func (_Gamecontract *GamecontractSession) AddVote(player common.Address, chainId uint32, team uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.AddVote(&_Gamecontract.TransactOpts, player, chainId, team)
}

// AddVote is a paid mutator transaction binding the contract method 0x8f66d4ec.
//
// Solidity: function addVote(address player, uint32 chainId, uint8 team) returns()
func (_Gamecontract *GamecontractTransactorSession) AddVote(player common.Address, chainId uint32, team uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.AddVote(&_Gamecontract.TransactOpts, player, chainId, team)
}

// EndGame is a paid mutator transaction binding the contract method 0x863375ac.
//
// Solidity: function endGame(uint8 result) returns()
func (_Gamecontract *GamecontractTransactor) EndGame(opts *bind.TransactOpts, result uint8) (*types.Transaction, error) {
	return _Gamecontract.contract.Transact(opts, "endGame", result)
}

// EndGame is a paid mutator transaction binding the contract method 0x863375ac.
//
// Solidity: function endGame(uint8 result) returns()
func (_Gamecontract *GamecontractSession) EndGame(result uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.EndGame(&_Gamecontract.TransactOpts, result)
}

// EndGame is a paid mutator transaction binding the contract method 0x863375ac.
//
// Solidity: function endGame(uint8 result) returns()
func (_Gamecontract *GamecontractTransactorSession) EndGame(result uint8) (*types.Transaction, error) {
	return _Gamecontract.Contract.EndGame(&_Gamecontract.TransactOpts, result)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4a30116.
//
// Solidity: function initialize(uint256 _gameId, uint256 _fixedStakeAmount) returns()
func (_Gamecontract *GamecontractTransactor) Initialize(opts *bind.TransactOpts, _gameId *big.Int, _fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Gamecontract.contract.Transact(opts, "initialize", _gameId, _fixedStakeAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4a30116.
//
// Solidity: function initialize(uint256 _gameId, uint256 _fixedStakeAmount) returns()
func (_Gamecontract *GamecontractSession) Initialize(_gameId *big.Int, _fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Gamecontract.Contract.Initialize(&_Gamecontract.TransactOpts, _gameId, _fixedStakeAmount)
}

// Initialize is a paid mutator transaction binding the contract method 0xe4a30116.
//
// Solidity: function initialize(uint256 _gameId, uint256 _fixedStakeAmount) returns()
func (_Gamecontract *GamecontractTransactorSession) Initialize(_gameId *big.Int, _fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Gamecontract.Contract.Initialize(&_Gamecontract.TransactOpts, _gameId, _fixedStakeAmount)
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

// GamecontractGameFinishedIterator is returned from FilterGameFinished and is used to iterate over the raw logs and unpacked data for GameFinished events raised by the Gamecontract contract.
type GamecontractGameFinishedIterator struct {
	Event *GamecontractGameFinished // Event containing the contract specifics and raw log

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
func (it *GamecontractGameFinishedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GamecontractGameFinished)
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
		it.Event = new(GamecontractGameFinished)
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
func (it *GamecontractGameFinishedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GamecontractGameFinishedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GamecontractGameFinished represents a GameFinished event raised by the Gamecontract contract.
type GamecontractGameFinished struct {
	GameId     *big.Int
	Result     uint8
	FinishedAt *big.Int
	Raw        types.Log // Blockchain specific contextual infos
}

// FilterGameFinished is a free log retrieval operation binding the contract event 0x0021feea4608e6756f982d5bd71c09645176b4865cf6bd600fcc4aa990f02404.
//
// Solidity: event GameFinished(uint256 indexed gameId, uint8 result, uint256 finishedAt)
func (_Gamecontract *GamecontractFilterer) FilterGameFinished(opts *bind.FilterOpts, gameId []*big.Int) (*GamecontractGameFinishedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Gamecontract.contract.FilterLogs(opts, "GameFinished", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &GamecontractGameFinishedIterator{contract: _Gamecontract.contract, event: "GameFinished", logs: logs, sub: sub}, nil
}

// WatchGameFinished is a free log subscription operation binding the contract event 0x0021feea4608e6756f982d5bd71c09645176b4865cf6bd600fcc4aa990f02404.
//
// Solidity: event GameFinished(uint256 indexed gameId, uint8 result, uint256 finishedAt)
func (_Gamecontract *GamecontractFilterer) WatchGameFinished(opts *bind.WatchOpts, sink chan<- *GamecontractGameFinished, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Gamecontract.contract.WatchLogs(opts, "GameFinished", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GamecontractGameFinished)
				if err := _Gamecontract.contract.UnpackLog(event, "GameFinished", log); err != nil {
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

// ParseGameFinished is a log parse operation binding the contract event 0x0021feea4608e6756f982d5bd71c09645176b4865cf6bd600fcc4aa990f02404.
//
// Solidity: event GameFinished(uint256 indexed gameId, uint8 result, uint256 finishedAt)
func (_Gamecontract *GamecontractFilterer) ParseGameFinished(log types.Log) (*GamecontractGameFinished, error) {
	event := new(GamecontractGameFinished)
	if err := _Gamecontract.contract.UnpackLog(event, "GameFinished", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// GamecontractVoteIterator is returned from FilterVote and is used to iterate over the raw logs and unpacked data for Vote events raised by the Gamecontract contract.
type GamecontractVoteIterator struct {
	Event *GamecontractVote // Event containing the contract specifics and raw log

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
func (it *GamecontractVoteIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(GamecontractVote)
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
		it.Event = new(GamecontractVote)
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
func (it *GamecontractVoteIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *GamecontractVoteIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// GamecontractVote represents a Vote event raised by the Gamecontract contract.
type GamecontractVote struct {
	GameId  *big.Int
	Player  common.Address
	Team    uint8
	ChainId uint32
	Raw     types.Log // Blockchain specific contextual infos
}

// FilterVote is a free log retrieval operation binding the contract event 0xca41c755a10d5981d66a7053c39564f2f68151dbd0ba5295fea79ed7ec171786.
//
// Solidity: event Vote(uint256 indexed gameId, address indexed player, uint8 team, uint32 chainId)
func (_Gamecontract *GamecontractFilterer) FilterVote(opts *bind.FilterOpts, gameId []*big.Int, player []common.Address) (*GamecontractVoteIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Gamecontract.contract.FilterLogs(opts, "Vote", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return &GamecontractVoteIterator{contract: _Gamecontract.contract, event: "Vote", logs: logs, sub: sub}, nil
}

// WatchVote is a free log subscription operation binding the contract event 0xca41c755a10d5981d66a7053c39564f2f68151dbd0ba5295fea79ed7ec171786.
//
// Solidity: event Vote(uint256 indexed gameId, address indexed player, uint8 team, uint32 chainId)
func (_Gamecontract *GamecontractFilterer) WatchVote(opts *bind.WatchOpts, sink chan<- *GamecontractVote, gameId []*big.Int, player []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Gamecontract.contract.WatchLogs(opts, "Vote", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(GamecontractVote)
				if err := _Gamecontract.contract.UnpackLog(event, "Vote", log); err != nil {
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

// ParseVote is a log parse operation binding the contract event 0xca41c755a10d5981d66a7053c39564f2f68151dbd0ba5295fea79ed7ec171786.
//
// Solidity: event Vote(uint256 indexed gameId, address indexed player, uint8 team, uint32 chainId)
func (_Gamecontract *GamecontractFilterer) ParseVote(log types.Log) (*GamecontractVote, error) {
	event := new(GamecontractVote)
	if err := _Gamecontract.contract.UnpackLog(event, "Vote", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
