// Code generated - DO NOT EDIT.
// This file is a generated binding and any manual changes will be lost.

package vaultcontract

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

// IVaultContractGameVaultInfo is an auto generated low-level Go binding around an user-defined struct.
type IVaultContractGameVaultInfo struct {
	GameId      *big.Int
	TotalStakes *big.Int
	PlayerCount *big.Int
	Result      uint8
	GameEnded   bool
	EndedAt     *big.Int
}

// IVaultContractPlayerStakeInfo is an auto generated low-level Go binding around an user-defined struct.
type IVaultContractPlayerStakeInfo struct {
	TotalStaked *big.Int
	StakeCount  *big.Int
	HasClaimed  bool
}

// VaultcontractMetaData contains all meta data concerning the Vaultcontract contract.
var VaultcontractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_authorizedBackend\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_gameContract\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"authorizedBackend\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"canClaimRewards\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"claimRewards\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"endGame\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIVaultContract.GameResult\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"gameContract\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gamePlayers\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"gameVaults\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"playerCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIVaultContract.GameResult\"},{\"name\":\"gameEnded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"endedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getGameVaultInfo\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIVaultContract.GameVaultInfo\",\"components\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"totalStakes\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"playerCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"result\",\"type\":\"uint8\",\"internalType\":\"enumIVaultContract.GameResult\"},{\"name\":\"gameEnded\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"endedAt\",\"type\":\"uint256\",\"internalType\":\"uint256\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayerStake\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getPlayerStakeInfo\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIVaultContract.PlayerStakeInfo\",\"components\":[{\"name\":\"totalStaked\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"hasClaimed\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalGameStakes\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isGameEnded\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"playerStakes\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"internalType\":\"address\"}],\"outputs\":[{\"name\":\"totalStaked\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"stakeCount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"hasClaimed\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"fixedStakeAmount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"payable\"},{\"type\":\"event\",\"name\":\"GameEndedInVault\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"result\",\"type\":\"uint8\",\"indexed\":false,\"internalType\":\"enumIVaultContract.GameResult\"},{\"name\":\"totalStakes\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"endedAt\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"RewardsClaimed\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeDeposited\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"player\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newTotal\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x60c03461014757601f610f5938819003918201601f19168301916001600160401b0383118484101761014b57808492604094855283398101031261014757610052602061004b8361015f565b920161015f565b906001600160a01b038116156100f0576001600160a01b0382161561009f5760805260a052604051610de5908161017482396080518181816102840152610981015260a051816101ae0152f35b60405162461bcd60e51b8152602060048201526024808201527f47616d6520636f6e74726163742063616e6e6f74206265207a65726f206164646044820152637265737360e01b6064820152608490fd5b60405162461bcd60e51b815260206004820152602960248201527f417574686f72697a6564206261636b656e642063616e6e6f74206265207a65726044820152686f206164647265737360b81b6064820152608490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b51906001600160a01b03821682036101475756fe6080806040526004361015610012575f80fd5b5f3560e01c9081630962ef7914610a3757508063182206dd146109e45780633c561555146109b057806349430bcf1461096c578063524828e9146109005780637b0472f014610651578063acbe73a81461057e578063b44f4a5514610547578063bb197c4f1461051b578063c544953114610262578063c748b362146101dd578063d3f3300914610199578063dfefef9f146101515763f328dc16146100b6575f80fd5b3461014d5760206100c636610c3c565b5f8281526001845260408082206001600160a01b039093168252918452209060ff60026100f1610cce565b9380548552600181015486860152015416906040830191151582525f525f835260ff600360405f20015460081c169182610142575b5081610138575b506040519015158152f35b905051155f61012d565b51151591505f610126565b5f80fd5b3461014d5761015f36610c6f565b905f52600260205260405f20805482101561014d5760209161018091610c85565b905460405160039290921b1c6001600160a01b03168152f35b3461014d575f36600319011261014d576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b3461014d576101eb36610c3c565b905f60406101f7610cce565b82815282602082015201525f52600160205260405f209060018060a01b03165f52602052606060405f20610229610cce565b815491828252604060ff600260018401549360208601948552015416920191151582526040519283525160208301525115156040820152f35b3461014d57604036600319011261014d57600435602435600481101561014d577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633036104c857801561047757815f525f60205260ff600360405f20015460081c1661043d57815f525f60205260405f2054155f146103d8576102ec610cae565b82815260208101915f835260408201915f8352606081019261030e8385610cee565b60808201906001825260a0830195428752875f525f60205260405f20935184555160018401555160028301556003820193519460048610156103c4576103958560049361037e7f8f2be90d4d3fe5486620a0848836ee043050dffa2ea6eb24d4449bd614823f9a99606099610cfa565b51815461ff00191690151560081b61ff0016179055565b519101555b835f525f602052600160405f2001546103b66040518093610c62565b6020820152426040820152a2005b634e487b7160e01b5f52602160045260245ffd5b60607f8f2be90d4d3fe5486620a0848836ee043050dffa2ea6eb24d4449bd614823f9a91835f525f60205261041381600360405f2001610cfa565b5f84815260208190526040902060038101805461ff0019166101001790554260049091015561039a565b60405162461bcd60e51b815260206004820152601260248201527111d85b5948185b1c9958591e48195b99195960721b6044820152606490fd5b60405162461bcd60e51b815260206004820152602360248201527f43616e6e6f7420656e642067616d652077697468206f6e676f696e67207265736044820152621d5b1d60ea1b6064820152608490fd5b60405162461bcd60e51b815260206004820152602560248201527f4f6e6c7920617574686f72697a6564206261636b656e642063616e20656e642060448201526467616d657360d81b6064820152608490fd5b3461014d57602036600319011261014d576004355f525f6020526020600160405f200154604051908152f35b3461014d5761055536610c3c565b905f52600160205260405f209060018060a01b03165f52602052602060405f2054604051908152f35b3461014d57602036600319011261014d575f60a061059a610cae565b82815282602082015282604082015282606082015282608082015201526004355f525f60205260c060405f206105ce610cae565b908054825260018101549060208301918252610641600282015491604085019283526003810154926004606087019261060a60ff871685610cee565b60ff608089019660081c161515865201549460a0870195865260405196518752516020870152516040860152516060850190610c62565b51151560808301525160a0820152f35b61065a36610c6f565b90805f525f60205260ff600360405f20015460081c166108bb5781156108645781340361080e57805f525f60205260405f205415610799575b5f81815260016020818152604080842033855290915290912090810180541561071f575b6106c2848354610d34565b82556106ce8154610d12565b9055815f525f602052600160405f20016106e9848254610d34565b90555460405192835260208301527f8700e7c955551ade34647b1ecf5ba06678ded0a43920ae3f84dc0941d3804f0060403393a3005b825f52600260205260405f208054680100000000000000008110156107855761074d91600182018155610c85565b81549060031b9033821b9160018060a01b03901b1916179055825f525f602052600260405f200161077e8154610d12565b90556106b7565b634e487b7160e01b5f52604160045260245ffd5b6107a1610cae565b818152602081015f815260408201905f8252606083015f815260808401925f845260a08501925f8452865f525f60205260405f20955186555160018601555160028501556003840190519060048210156103c45760049361037e6108059383610cfa565b51910155610693565b60405162461bcd60e51b815260206004820152602860248201527f53656e742076616c7565206d757374206d61746368206669786564207374616b6044820152671948185b5bdd5b9d60c21b6064820152608490fd5b60405162461bcd60e51b815260206004820152602960248201527f4669786564207374616b6520616d6f756e74206d75737420626520677265617460448201526806572207468616e20360bc1b6064820152608490fd5b60405162461bcd60e51b815260206004820152601c60248201527f47616d652068617320656e6465642c2063616e6e6f74207374616b65000000006044820152606490fd5b3461014d57602036600319011261014d576004355f525f60205260c060405f2080549060ff600182015491600281015460046003830154920154936040519586526020860152604085015261095a60608501838316610c62565b60081c161515608083015260a0820152f35b3461014d575f36600319011261014d576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b3461014d57602036600319011261014d576004355f525f602052602060ff600360405f20015460081c166040519015158152f35b3461014d576109f236610c3c565b905f52600160205260405f209060018060a01b03165f52602052606060405f2080549060ff600260018301549201541690604051928352602083015215156040820152f35b3461014d57602036600319011261014d5760043590815f525f60205260ff600360405f20015460081c1615610c0157505f8181526001602090815260408083203384529091529020805491908215610bc45760020180549260ff8416610b7f57610aa360019184610d41565b9360ff191617905581610adf575b6040519182527f3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb79696860203393a3005b5f80808085335af13d15610b7a573d67ffffffffffffffff81116107855760405190601f8101601f19908116603f0116820167ffffffffffffffff8111838210176107855760405281525f60203d92013e5b610ab15760405162461bcd60e51b815260206004820152601a60248201527f4661696c656420746f207472616e7366657220726577617264730000000000006044820152606490fd5b610b31565b60405162461bcd60e51b815260206004820152601760248201527f5265776172647320616c726561647920636c61696d65640000000000000000006044820152606490fd5b60405162461bcd60e51b81526020600482015260156024820152744e6f207374616b6520696e20746869732067616d6560581b6044820152606490fd5b62461bcd60e51b815260206004820152601660248201527511d85b59481a185cc81b9bdd08195b991959081e595d60521b6044820152606490fd5b604090600319011261014d57600435906024356001600160a01b038116810361014d5790565b9060048210156103c45752565b604090600319011261014d576004359060243590565b8054821015610c9a575f5260205f2001905f90565b634e487b7160e01b5f52603260045260245ffd5b6040519060c0820182811067ffffffffffffffff82111761078557604052565b604051906060820182811067ffffffffffffffff82111761078557604052565b60048210156103c45752565b9060048110156103c45760ff80198354169116179055565b5f198114610d205760010190565b634e487b7160e01b5f52601160045260245ffd5b91908201809211610d2057565b5f525f60205260405f20610d53610cae565b81548152600182015460208201526002820154604082015260a0600460038401549360ff6060850195610d8882821688610cee565b60081c161515608085015201549101525160048110156103c457600314610dac5790565b9056fea264697066735822122090edea765e07d193487489502e5312f69bd2f0a5050696184e55ff577295aed764736f6c634300081e0033",
}

// VaultcontractABI is the input ABI used to generate the binding from.
// Deprecated: Use VaultcontractMetaData.ABI instead.
var VaultcontractABI = VaultcontractMetaData.ABI

// VaultcontractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VaultcontractMetaData.Bin instead.
var VaultcontractBin = VaultcontractMetaData.Bin

// DeployVaultcontract deploys a new Ethereum contract, binding an instance of Vaultcontract to it.
func DeployVaultcontract(auth *bind.TransactOpts, backend bind.ContractBackend, _authorizedBackend common.Address, _gameContract common.Address) (common.Address, *types.Transaction, *Vaultcontract, error) {
	parsed, err := VaultcontractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VaultcontractBin), backend, _authorizedBackend, _gameContract)
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	return address, tx, &Vaultcontract{VaultcontractCaller: VaultcontractCaller{contract: contract}, VaultcontractTransactor: VaultcontractTransactor{contract: contract}, VaultcontractFilterer: VaultcontractFilterer{contract: contract}}, nil
}

// Vaultcontract is an auto generated Go binding around an Ethereum contract.
type Vaultcontract struct {
	VaultcontractCaller     // Read-only binding to the contract
	VaultcontractTransactor // Write-only binding to the contract
	VaultcontractFilterer   // Log filterer for contract events
}

// VaultcontractCaller is an auto generated read-only Go binding around an Ethereum contract.
type VaultcontractCaller struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultcontractTransactor is an auto generated write-only Go binding around an Ethereum contract.
type VaultcontractTransactor struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultcontractFilterer is an auto generated log filtering Go binding around an Ethereum contract events.
type VaultcontractFilterer struct {
	contract *bind.BoundContract // Generic contract wrapper for the low level calls
}

// VaultcontractSession is an auto generated Go binding around an Ethereum contract,
// with pre-set call and transact options.
type VaultcontractSession struct {
	Contract     *Vaultcontract    // Generic contract binding to set the session for
	CallOpts     bind.CallOpts     // Call options to use throughout this session
	TransactOpts bind.TransactOpts // Transaction auth options to use throughout this session
}

// VaultcontractCallerSession is an auto generated read-only Go binding around an Ethereum contract,
// with pre-set call options.
type VaultcontractCallerSession struct {
	Contract *VaultcontractCaller // Generic contract caller binding to set the session for
	CallOpts bind.CallOpts        // Call options to use throughout this session
}

// VaultcontractTransactorSession is an auto generated write-only Go binding around an Ethereum contract,
// with pre-set transact options.
type VaultcontractTransactorSession struct {
	Contract     *VaultcontractTransactor // Generic contract transactor binding to set the session for
	TransactOpts bind.TransactOpts        // Transaction auth options to use throughout this session
}

// VaultcontractRaw is an auto generated low-level Go binding around an Ethereum contract.
type VaultcontractRaw struct {
	Contract *Vaultcontract // Generic contract binding to access the raw methods on
}

// VaultcontractCallerRaw is an auto generated low-level read-only Go binding around an Ethereum contract.
type VaultcontractCallerRaw struct {
	Contract *VaultcontractCaller // Generic read-only contract binding to access the raw methods on
}

// VaultcontractTransactorRaw is an auto generated low-level write-only Go binding around an Ethereum contract.
type VaultcontractTransactorRaw struct {
	Contract *VaultcontractTransactor // Generic write-only contract binding to access the raw methods on
}

// NewVaultcontract creates a new instance of Vaultcontract, bound to a specific deployed contract.
func NewVaultcontract(address common.Address, backend bind.ContractBackend) (*Vaultcontract, error) {
	contract, err := bindVaultcontract(address, backend, backend, backend)
	if err != nil {
		return nil, err
	}
	return &Vaultcontract{VaultcontractCaller: VaultcontractCaller{contract: contract}, VaultcontractTransactor: VaultcontractTransactor{contract: contract}, VaultcontractFilterer: VaultcontractFilterer{contract: contract}}, nil
}

// NewVaultcontractCaller creates a new read-only instance of Vaultcontract, bound to a specific deployed contract.
func NewVaultcontractCaller(address common.Address, caller bind.ContractCaller) (*VaultcontractCaller, error) {
	contract, err := bindVaultcontract(address, caller, nil, nil)
	if err != nil {
		return nil, err
	}
	return &VaultcontractCaller{contract: contract}, nil
}

// NewVaultcontractTransactor creates a new write-only instance of Vaultcontract, bound to a specific deployed contract.
func NewVaultcontractTransactor(address common.Address, transactor bind.ContractTransactor) (*VaultcontractTransactor, error) {
	contract, err := bindVaultcontract(address, nil, transactor, nil)
	if err != nil {
		return nil, err
	}
	return &VaultcontractTransactor{contract: contract}, nil
}

// NewVaultcontractFilterer creates a new log filterer instance of Vaultcontract, bound to a specific deployed contract.
func NewVaultcontractFilterer(address common.Address, filterer bind.ContractFilterer) (*VaultcontractFilterer, error) {
	contract, err := bindVaultcontract(address, nil, nil, filterer)
	if err != nil {
		return nil, err
	}
	return &VaultcontractFilterer{contract: contract}, nil
}

// bindVaultcontract binds a generic wrapper to an already deployed contract.
func bindVaultcontract(address common.Address, caller bind.ContractCaller, transactor bind.ContractTransactor, filterer bind.ContractFilterer) (*bind.BoundContract, error) {
	parsed, err := VaultcontractMetaData.GetAbi()
	if err != nil {
		return nil, err
	}
	return bind.NewBoundContract(address, *parsed, caller, transactor, filterer), nil
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vaultcontract *VaultcontractRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vaultcontract.Contract.VaultcontractCaller.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vaultcontract *VaultcontractRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vaultcontract.Contract.VaultcontractTransactor.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vaultcontract *VaultcontractRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vaultcontract.Contract.VaultcontractTransactor.contract.Transact(opts, method, params...)
}

// Call invokes the (constant) contract method with params as input values and
// sets the output to result. The result type might be a single field for simple
// returns, a slice of interfaces for anonymous returns and a struct for named
// returns.
func (_Vaultcontract *VaultcontractCallerRaw) Call(opts *bind.CallOpts, result *[]interface{}, method string, params ...interface{}) error {
	return _Vaultcontract.Contract.contract.Call(opts, result, method, params...)
}

// Transfer initiates a plain transaction to move funds to the contract, calling
// its default method if one is available.
func (_Vaultcontract *VaultcontractTransactorRaw) Transfer(opts *bind.TransactOpts) (*types.Transaction, error) {
	return _Vaultcontract.Contract.contract.Transfer(opts)
}

// Transact invokes the (paid) contract method with params as input values.
func (_Vaultcontract *VaultcontractTransactorRaw) Transact(opts *bind.TransactOpts, method string, params ...interface{}) (*types.Transaction, error) {
	return _Vaultcontract.Contract.contract.Transact(opts, method, params...)
}

// AuthorizedBackend is a free data retrieval call binding the contract method 0x49430bcf.
//
// Solidity: function authorizedBackend() view returns(address)
func (_Vaultcontract *VaultcontractCaller) AuthorizedBackend(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "authorizedBackend")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AuthorizedBackend is a free data retrieval call binding the contract method 0x49430bcf.
//
// Solidity: function authorizedBackend() view returns(address)
func (_Vaultcontract *VaultcontractSession) AuthorizedBackend() (common.Address, error) {
	return _Vaultcontract.Contract.AuthorizedBackend(&_Vaultcontract.CallOpts)
}

// AuthorizedBackend is a free data retrieval call binding the contract method 0x49430bcf.
//
// Solidity: function authorizedBackend() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) AuthorizedBackend() (common.Address, error) {
	return _Vaultcontract.Contract.AuthorizedBackend(&_Vaultcontract.CallOpts)
}

// CanClaimRewards is a free data retrieval call binding the contract method 0xf328dc16.
//
// Solidity: function canClaimRewards(uint256 gameId, address player) view returns(bool)
func (_Vaultcontract *VaultcontractCaller) CanClaimRewards(opts *bind.CallOpts, gameId *big.Int, player common.Address) (bool, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "canClaimRewards", gameId, player)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// CanClaimRewards is a free data retrieval call binding the contract method 0xf328dc16.
//
// Solidity: function canClaimRewards(uint256 gameId, address player) view returns(bool)
func (_Vaultcontract *VaultcontractSession) CanClaimRewards(gameId *big.Int, player common.Address) (bool, error) {
	return _Vaultcontract.Contract.CanClaimRewards(&_Vaultcontract.CallOpts, gameId, player)
}

// CanClaimRewards is a free data retrieval call binding the contract method 0xf328dc16.
//
// Solidity: function canClaimRewards(uint256 gameId, address player) view returns(bool)
func (_Vaultcontract *VaultcontractCallerSession) CanClaimRewards(gameId *big.Int, player common.Address) (bool, error) {
	return _Vaultcontract.Contract.CanClaimRewards(&_Vaultcontract.CallOpts, gameId, player)
}

// GameContract is a free data retrieval call binding the contract method 0xd3f33009.
//
// Solidity: function gameContract() view returns(address)
func (_Vaultcontract *VaultcontractCaller) GameContract(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "gameContract")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GameContract is a free data retrieval call binding the contract method 0xd3f33009.
//
// Solidity: function gameContract() view returns(address)
func (_Vaultcontract *VaultcontractSession) GameContract() (common.Address, error) {
	return _Vaultcontract.Contract.GameContract(&_Vaultcontract.CallOpts)
}

// GameContract is a free data retrieval call binding the contract method 0xd3f33009.
//
// Solidity: function gameContract() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) GameContract() (common.Address, error) {
	return _Vaultcontract.Contract.GameContract(&_Vaultcontract.CallOpts)
}

// GamePlayers is a free data retrieval call binding the contract method 0xdfefef9f.
//
// Solidity: function gamePlayers(uint256 gameId, uint256 ) view returns(address)
func (_Vaultcontract *VaultcontractCaller) GamePlayers(opts *bind.CallOpts, gameId *big.Int, arg1 *big.Int) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "gamePlayers", gameId, arg1)

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GamePlayers is a free data retrieval call binding the contract method 0xdfefef9f.
//
// Solidity: function gamePlayers(uint256 gameId, uint256 ) view returns(address)
func (_Vaultcontract *VaultcontractSession) GamePlayers(gameId *big.Int, arg1 *big.Int) (common.Address, error) {
	return _Vaultcontract.Contract.GamePlayers(&_Vaultcontract.CallOpts, gameId, arg1)
}

// GamePlayers is a free data retrieval call binding the contract method 0xdfefef9f.
//
// Solidity: function gamePlayers(uint256 gameId, uint256 ) view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) GamePlayers(gameId *big.Int, arg1 *big.Int) (common.Address, error) {
	return _Vaultcontract.Contract.GamePlayers(&_Vaultcontract.CallOpts, gameId, arg1)
}

// GameVaults is a free data retrieval call binding the contract method 0x524828e9.
//
// Solidity: function gameVaults(uint256 gameId) view returns(uint256 gameId, uint256 totalStakes, uint256 playerCount, uint8 result, bool gameEnded, uint256 endedAt)
func (_Vaultcontract *VaultcontractCaller) GameVaults(opts *bind.CallOpts, gameId *big.Int) (struct {
	GameId      *big.Int
	TotalStakes *big.Int
	PlayerCount *big.Int
	Result      uint8
	GameEnded   bool
	EndedAt     *big.Int
}, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "gameVaults", gameId)

	outstruct := new(struct {
		GameId      *big.Int
		TotalStakes *big.Int
		PlayerCount *big.Int
		Result      uint8
		GameEnded   bool
		EndedAt     *big.Int
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.GameId = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.TotalStakes = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.PlayerCount = *abi.ConvertType(out[2], new(*big.Int)).(**big.Int)
	outstruct.Result = *abi.ConvertType(out[3], new(uint8)).(*uint8)
	outstruct.GameEnded = *abi.ConvertType(out[4], new(bool)).(*bool)
	outstruct.EndedAt = *abi.ConvertType(out[5], new(*big.Int)).(**big.Int)

	return *outstruct, err

}

// GameVaults is a free data retrieval call binding the contract method 0x524828e9.
//
// Solidity: function gameVaults(uint256 gameId) view returns(uint256 gameId, uint256 totalStakes, uint256 playerCount, uint8 result, bool gameEnded, uint256 endedAt)
func (_Vaultcontract *VaultcontractSession) GameVaults(gameId *big.Int) (struct {
	GameId      *big.Int
	TotalStakes *big.Int
	PlayerCount *big.Int
	Result      uint8
	GameEnded   bool
	EndedAt     *big.Int
}, error) {
	return _Vaultcontract.Contract.GameVaults(&_Vaultcontract.CallOpts, gameId)
}

// GameVaults is a free data retrieval call binding the contract method 0x524828e9.
//
// Solidity: function gameVaults(uint256 gameId) view returns(uint256 gameId, uint256 totalStakes, uint256 playerCount, uint8 result, bool gameEnded, uint256 endedAt)
func (_Vaultcontract *VaultcontractCallerSession) GameVaults(gameId *big.Int) (struct {
	GameId      *big.Int
	TotalStakes *big.Int
	PlayerCount *big.Int
	Result      uint8
	GameEnded   bool
	EndedAt     *big.Int
}, error) {
	return _Vaultcontract.Contract.GameVaults(&_Vaultcontract.CallOpts, gameId)
}

// GetGameVaultInfo is a free data retrieval call binding the contract method 0xacbe73a8.
//
// Solidity: function getGameVaultInfo(uint256 gameId) view returns((uint256,uint256,uint256,uint8,bool,uint256))
func (_Vaultcontract *VaultcontractCaller) GetGameVaultInfo(opts *bind.CallOpts, gameId *big.Int) (IVaultContractGameVaultInfo, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getGameVaultInfo", gameId)

	if err != nil {
		return *new(IVaultContractGameVaultInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IVaultContractGameVaultInfo)).(*IVaultContractGameVaultInfo)

	return out0, err

}

// GetGameVaultInfo is a free data retrieval call binding the contract method 0xacbe73a8.
//
// Solidity: function getGameVaultInfo(uint256 gameId) view returns((uint256,uint256,uint256,uint8,bool,uint256))
func (_Vaultcontract *VaultcontractSession) GetGameVaultInfo(gameId *big.Int) (IVaultContractGameVaultInfo, error) {
	return _Vaultcontract.Contract.GetGameVaultInfo(&_Vaultcontract.CallOpts, gameId)
}

// GetGameVaultInfo is a free data retrieval call binding the contract method 0xacbe73a8.
//
// Solidity: function getGameVaultInfo(uint256 gameId) view returns((uint256,uint256,uint256,uint8,bool,uint256))
func (_Vaultcontract *VaultcontractCallerSession) GetGameVaultInfo(gameId *big.Int) (IVaultContractGameVaultInfo, error) {
	return _Vaultcontract.Contract.GetGameVaultInfo(&_Vaultcontract.CallOpts, gameId)
}

// GetPlayerStake is a free data retrieval call binding the contract method 0xb44f4a55.
//
// Solidity: function getPlayerStake(uint256 gameId, address player) view returns(uint256)
func (_Vaultcontract *VaultcontractCaller) GetPlayerStake(opts *bind.CallOpts, gameId *big.Int, player common.Address) (*big.Int, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getPlayerStake", gameId, player)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetPlayerStake is a free data retrieval call binding the contract method 0xb44f4a55.
//
// Solidity: function getPlayerStake(uint256 gameId, address player) view returns(uint256)
func (_Vaultcontract *VaultcontractSession) GetPlayerStake(gameId *big.Int, player common.Address) (*big.Int, error) {
	return _Vaultcontract.Contract.GetPlayerStake(&_Vaultcontract.CallOpts, gameId, player)
}

// GetPlayerStake is a free data retrieval call binding the contract method 0xb44f4a55.
//
// Solidity: function getPlayerStake(uint256 gameId, address player) view returns(uint256)
func (_Vaultcontract *VaultcontractCallerSession) GetPlayerStake(gameId *big.Int, player common.Address) (*big.Int, error) {
	return _Vaultcontract.Contract.GetPlayerStake(&_Vaultcontract.CallOpts, gameId, player)
}

// GetPlayerStakeInfo is a free data retrieval call binding the contract method 0xc748b362.
//
// Solidity: function getPlayerStakeInfo(uint256 gameId, address player) view returns((uint256,uint256,bool))
func (_Vaultcontract *VaultcontractCaller) GetPlayerStakeInfo(opts *bind.CallOpts, gameId *big.Int, player common.Address) (IVaultContractPlayerStakeInfo, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getPlayerStakeInfo", gameId, player)

	if err != nil {
		return *new(IVaultContractPlayerStakeInfo), err
	}

	out0 := *abi.ConvertType(out[0], new(IVaultContractPlayerStakeInfo)).(*IVaultContractPlayerStakeInfo)

	return out0, err

}

// GetPlayerStakeInfo is a free data retrieval call binding the contract method 0xc748b362.
//
// Solidity: function getPlayerStakeInfo(uint256 gameId, address player) view returns((uint256,uint256,bool))
func (_Vaultcontract *VaultcontractSession) GetPlayerStakeInfo(gameId *big.Int, player common.Address) (IVaultContractPlayerStakeInfo, error) {
	return _Vaultcontract.Contract.GetPlayerStakeInfo(&_Vaultcontract.CallOpts, gameId, player)
}

// GetPlayerStakeInfo is a free data retrieval call binding the contract method 0xc748b362.
//
// Solidity: function getPlayerStakeInfo(uint256 gameId, address player) view returns((uint256,uint256,bool))
func (_Vaultcontract *VaultcontractCallerSession) GetPlayerStakeInfo(gameId *big.Int, player common.Address) (IVaultContractPlayerStakeInfo, error) {
	return _Vaultcontract.Contract.GetPlayerStakeInfo(&_Vaultcontract.CallOpts, gameId, player)
}

// GetTotalGameStakes is a free data retrieval call binding the contract method 0xbb197c4f.
//
// Solidity: function getTotalGameStakes(uint256 gameId) view returns(uint256)
func (_Vaultcontract *VaultcontractCaller) GetTotalGameStakes(opts *bind.CallOpts, gameId *big.Int) (*big.Int, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getTotalGameStakes", gameId)

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalGameStakes is a free data retrieval call binding the contract method 0xbb197c4f.
//
// Solidity: function getTotalGameStakes(uint256 gameId) view returns(uint256)
func (_Vaultcontract *VaultcontractSession) GetTotalGameStakes(gameId *big.Int) (*big.Int, error) {
	return _Vaultcontract.Contract.GetTotalGameStakes(&_Vaultcontract.CallOpts, gameId)
}

// GetTotalGameStakes is a free data retrieval call binding the contract method 0xbb197c4f.
//
// Solidity: function getTotalGameStakes(uint256 gameId) view returns(uint256)
func (_Vaultcontract *VaultcontractCallerSession) GetTotalGameStakes(gameId *big.Int) (*big.Int, error) {
	return _Vaultcontract.Contract.GetTotalGameStakes(&_Vaultcontract.CallOpts, gameId)
}

// IsGameEnded is a free data retrieval call binding the contract method 0x3c561555.
//
// Solidity: function isGameEnded(uint256 gameId) view returns(bool)
func (_Vaultcontract *VaultcontractCaller) IsGameEnded(opts *bind.CallOpts, gameId *big.Int) (bool, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "isGameEnded", gameId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsGameEnded is a free data retrieval call binding the contract method 0x3c561555.
//
// Solidity: function isGameEnded(uint256 gameId) view returns(bool)
func (_Vaultcontract *VaultcontractSession) IsGameEnded(gameId *big.Int) (bool, error) {
	return _Vaultcontract.Contract.IsGameEnded(&_Vaultcontract.CallOpts, gameId)
}

// IsGameEnded is a free data retrieval call binding the contract method 0x3c561555.
//
// Solidity: function isGameEnded(uint256 gameId) view returns(bool)
func (_Vaultcontract *VaultcontractCallerSession) IsGameEnded(gameId *big.Int) (bool, error) {
	return _Vaultcontract.Contract.IsGameEnded(&_Vaultcontract.CallOpts, gameId)
}

// PlayerStakes is a free data retrieval call binding the contract method 0x182206dd.
//
// Solidity: function playerStakes(uint256 gameId, address player) view returns(uint256 totalStaked, uint256 stakeCount, bool hasClaimed)
func (_Vaultcontract *VaultcontractCaller) PlayerStakes(opts *bind.CallOpts, gameId *big.Int, player common.Address) (struct {
	TotalStaked *big.Int
	StakeCount  *big.Int
	HasClaimed  bool
}, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "playerStakes", gameId, player)

	outstruct := new(struct {
		TotalStaked *big.Int
		StakeCount  *big.Int
		HasClaimed  bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.TotalStaked = *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)
	outstruct.StakeCount = *abi.ConvertType(out[1], new(*big.Int)).(**big.Int)
	outstruct.HasClaimed = *abi.ConvertType(out[2], new(bool)).(*bool)

	return *outstruct, err

}

// PlayerStakes is a free data retrieval call binding the contract method 0x182206dd.
//
// Solidity: function playerStakes(uint256 gameId, address player) view returns(uint256 totalStaked, uint256 stakeCount, bool hasClaimed)
func (_Vaultcontract *VaultcontractSession) PlayerStakes(gameId *big.Int, player common.Address) (struct {
	TotalStaked *big.Int
	StakeCount  *big.Int
	HasClaimed  bool
}, error) {
	return _Vaultcontract.Contract.PlayerStakes(&_Vaultcontract.CallOpts, gameId, player)
}

// PlayerStakes is a free data retrieval call binding the contract method 0x182206dd.
//
// Solidity: function playerStakes(uint256 gameId, address player) view returns(uint256 totalStaked, uint256 stakeCount, bool hasClaimed)
func (_Vaultcontract *VaultcontractCallerSession) PlayerStakes(gameId *big.Int, player common.Address) (struct {
	TotalStaked *big.Int
	StakeCount  *big.Int
	HasClaimed  bool
}, error) {
	return _Vaultcontract.Contract.PlayerStakes(&_Vaultcontract.CallOpts, gameId, player)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0962ef79.
//
// Solidity: function claimRewards(uint256 gameId) returns()
func (_Vaultcontract *VaultcontractTransactor) ClaimRewards(opts *bind.TransactOpts, gameId *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.contract.Transact(opts, "claimRewards", gameId)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0962ef79.
//
// Solidity: function claimRewards(uint256 gameId) returns()
func (_Vaultcontract *VaultcontractSession) ClaimRewards(gameId *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.ClaimRewards(&_Vaultcontract.TransactOpts, gameId)
}

// ClaimRewards is a paid mutator transaction binding the contract method 0x0962ef79.
//
// Solidity: function claimRewards(uint256 gameId) returns()
func (_Vaultcontract *VaultcontractTransactorSession) ClaimRewards(gameId *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.ClaimRewards(&_Vaultcontract.TransactOpts, gameId)
}

// EndGame is a paid mutator transaction binding the contract method 0xc5449531.
//
// Solidity: function endGame(uint256 gameId, uint8 result) returns()
func (_Vaultcontract *VaultcontractTransactor) EndGame(opts *bind.TransactOpts, gameId *big.Int, result uint8) (*types.Transaction, error) {
	return _Vaultcontract.contract.Transact(opts, "endGame", gameId, result)
}

// EndGame is a paid mutator transaction binding the contract method 0xc5449531.
//
// Solidity: function endGame(uint256 gameId, uint8 result) returns()
func (_Vaultcontract *VaultcontractSession) EndGame(gameId *big.Int, result uint8) (*types.Transaction, error) {
	return _Vaultcontract.Contract.EndGame(&_Vaultcontract.TransactOpts, gameId, result)
}

// EndGame is a paid mutator transaction binding the contract method 0xc5449531.
//
// Solidity: function endGame(uint256 gameId, uint8 result) returns()
func (_Vaultcontract *VaultcontractTransactorSession) EndGame(gameId *big.Int, result uint8) (*types.Transaction, error) {
	return _Vaultcontract.Contract.EndGame(&_Vaultcontract.TransactOpts, gameId, result)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 gameId, uint256 fixedStakeAmount) payable returns()
func (_Vaultcontract *VaultcontractTransactor) Stake(opts *bind.TransactOpts, gameId *big.Int, fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.contract.Transact(opts, "stake", gameId, fixedStakeAmount)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 gameId, uint256 fixedStakeAmount) payable returns()
func (_Vaultcontract *VaultcontractSession) Stake(gameId *big.Int, fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.Stake(&_Vaultcontract.TransactOpts, gameId, fixedStakeAmount)
}

// Stake is a paid mutator transaction binding the contract method 0x7b0472f0.
//
// Solidity: function stake(uint256 gameId, uint256 fixedStakeAmount) payable returns()
func (_Vaultcontract *VaultcontractTransactorSession) Stake(gameId *big.Int, fixedStakeAmount *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.Stake(&_Vaultcontract.TransactOpts, gameId, fixedStakeAmount)
}

// VaultcontractGameEndedInVaultIterator is returned from FilterGameEndedInVault and is used to iterate over the raw logs and unpacked data for GameEndedInVault events raised by the Vaultcontract contract.
type VaultcontractGameEndedInVaultIterator struct {
	Event *VaultcontractGameEndedInVault // Event containing the contract specifics and raw log

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
func (it *VaultcontractGameEndedInVaultIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultcontractGameEndedInVault)
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
		it.Event = new(VaultcontractGameEndedInVault)
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
func (it *VaultcontractGameEndedInVaultIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultcontractGameEndedInVaultIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultcontractGameEndedInVault represents a GameEndedInVault event raised by the Vaultcontract contract.
type VaultcontractGameEndedInVault struct {
	GameId      *big.Int
	Result      uint8
	TotalStakes *big.Int
	EndedAt     *big.Int
	Raw         types.Log // Blockchain specific contextual infos
}

// FilterGameEndedInVault is a free log retrieval operation binding the contract event 0x8f2be90d4d3fe5486620a0848836ee043050dffa2ea6eb24d4449bd614823f9a.
//
// Solidity: event GameEndedInVault(uint256 indexed gameId, uint8 result, uint256 totalStakes, uint256 endedAt)
func (_Vaultcontract *VaultcontractFilterer) FilterGameEndedInVault(opts *bind.FilterOpts, gameId []*big.Int) (*VaultcontractGameEndedInVaultIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Vaultcontract.contract.FilterLogs(opts, "GameEndedInVault", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &VaultcontractGameEndedInVaultIterator{contract: _Vaultcontract.contract, event: "GameEndedInVault", logs: logs, sub: sub}, nil
}

// WatchGameEndedInVault is a free log subscription operation binding the contract event 0x8f2be90d4d3fe5486620a0848836ee043050dffa2ea6eb24d4449bd614823f9a.
//
// Solidity: event GameEndedInVault(uint256 indexed gameId, uint8 result, uint256 totalStakes, uint256 endedAt)
func (_Vaultcontract *VaultcontractFilterer) WatchGameEndedInVault(opts *bind.WatchOpts, sink chan<- *VaultcontractGameEndedInVault, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Vaultcontract.contract.WatchLogs(opts, "GameEndedInVault", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultcontractGameEndedInVault)
				if err := _Vaultcontract.contract.UnpackLog(event, "GameEndedInVault", log); err != nil {
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

// ParseGameEndedInVault is a log parse operation binding the contract event 0x8f2be90d4d3fe5486620a0848836ee043050dffa2ea6eb24d4449bd614823f9a.
//
// Solidity: event GameEndedInVault(uint256 indexed gameId, uint8 result, uint256 totalStakes, uint256 endedAt)
func (_Vaultcontract *VaultcontractFilterer) ParseGameEndedInVault(log types.Log) (*VaultcontractGameEndedInVault, error) {
	event := new(VaultcontractGameEndedInVault)
	if err := _Vaultcontract.contract.UnpackLog(event, "GameEndedInVault", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultcontractRewardsClaimedIterator is returned from FilterRewardsClaimed and is used to iterate over the raw logs and unpacked data for RewardsClaimed events raised by the Vaultcontract contract.
type VaultcontractRewardsClaimedIterator struct {
	Event *VaultcontractRewardsClaimed // Event containing the contract specifics and raw log

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
func (it *VaultcontractRewardsClaimedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultcontractRewardsClaimed)
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
		it.Event = new(VaultcontractRewardsClaimed)
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
func (it *VaultcontractRewardsClaimedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultcontractRewardsClaimedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultcontractRewardsClaimed represents a RewardsClaimed event raised by the Vaultcontract contract.
type VaultcontractRewardsClaimed struct {
	GameId *big.Int
	Player common.Address
	Amount *big.Int
	Raw    types.Log // Blockchain specific contextual infos
}

// FilterRewardsClaimed is a free log retrieval operation binding the contract event 0x3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb796968.
//
// Solidity: event RewardsClaimed(uint256 indexed gameId, address indexed player, uint256 amount)
func (_Vaultcontract *VaultcontractFilterer) FilterRewardsClaimed(opts *bind.FilterOpts, gameId []*big.Int, player []common.Address) (*VaultcontractRewardsClaimedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Vaultcontract.contract.FilterLogs(opts, "RewardsClaimed", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return &VaultcontractRewardsClaimedIterator{contract: _Vaultcontract.contract, event: "RewardsClaimed", logs: logs, sub: sub}, nil
}

// WatchRewardsClaimed is a free log subscription operation binding the contract event 0x3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb796968.
//
// Solidity: event RewardsClaimed(uint256 indexed gameId, address indexed player, uint256 amount)
func (_Vaultcontract *VaultcontractFilterer) WatchRewardsClaimed(opts *bind.WatchOpts, sink chan<- *VaultcontractRewardsClaimed, gameId []*big.Int, player []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Vaultcontract.contract.WatchLogs(opts, "RewardsClaimed", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultcontractRewardsClaimed)
				if err := _Vaultcontract.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
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

// ParseRewardsClaimed is a log parse operation binding the contract event 0x3300bdb359cfb956935bca32e9db727413eab1ca84341f2e36caea85bb796968.
//
// Solidity: event RewardsClaimed(uint256 indexed gameId, address indexed player, uint256 amount)
func (_Vaultcontract *VaultcontractFilterer) ParseRewardsClaimed(log types.Log) (*VaultcontractRewardsClaimed, error) {
	event := new(VaultcontractRewardsClaimed)
	if err := _Vaultcontract.contract.UnpackLog(event, "RewardsClaimed", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}

// VaultcontractStakeDepositedIterator is returned from FilterStakeDeposited and is used to iterate over the raw logs and unpacked data for StakeDeposited events raised by the Vaultcontract contract.
type VaultcontractStakeDepositedIterator struct {
	Event *VaultcontractStakeDeposited // Event containing the contract specifics and raw log

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
func (it *VaultcontractStakeDepositedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultcontractStakeDeposited)
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
		it.Event = new(VaultcontractStakeDeposited)
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
func (it *VaultcontractStakeDepositedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultcontractStakeDepositedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultcontractStakeDeposited represents a StakeDeposited event raised by the Vaultcontract contract.
type VaultcontractStakeDeposited struct {
	GameId   *big.Int
	Player   common.Address
	Amount   *big.Int
	NewTotal *big.Int
	Raw      types.Log // Blockchain specific contextual infos
}

// FilterStakeDeposited is a free log retrieval operation binding the contract event 0x8700e7c955551ade34647b1ecf5ba06678ded0a43920ae3f84dc0941d3804f00.
//
// Solidity: event StakeDeposited(uint256 indexed gameId, address indexed player, uint256 amount, uint256 newTotal)
func (_Vaultcontract *VaultcontractFilterer) FilterStakeDeposited(opts *bind.FilterOpts, gameId []*big.Int, player []common.Address) (*VaultcontractStakeDepositedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Vaultcontract.contract.FilterLogs(opts, "StakeDeposited", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return &VaultcontractStakeDepositedIterator{contract: _Vaultcontract.contract, event: "StakeDeposited", logs: logs, sub: sub}, nil
}

// WatchStakeDeposited is a free log subscription operation binding the contract event 0x8700e7c955551ade34647b1ecf5ba06678ded0a43920ae3f84dc0941d3804f00.
//
// Solidity: event StakeDeposited(uint256 indexed gameId, address indexed player, uint256 amount, uint256 newTotal)
func (_Vaultcontract *VaultcontractFilterer) WatchStakeDeposited(opts *bind.WatchOpts, sink chan<- *VaultcontractStakeDeposited, gameId []*big.Int, player []common.Address) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}
	var playerRule []interface{}
	for _, playerItem := range player {
		playerRule = append(playerRule, playerItem)
	}

	logs, sub, err := _Vaultcontract.contract.WatchLogs(opts, "StakeDeposited", gameIdRule, playerRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultcontractStakeDeposited)
				if err := _Vaultcontract.contract.UnpackLog(event, "StakeDeposited", log); err != nil {
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

// ParseStakeDeposited is a log parse operation binding the contract event 0x8700e7c955551ade34647b1ecf5ba06678ded0a43920ae3f84dc0941d3804f00.
//
// Solidity: event StakeDeposited(uint256 indexed gameId, address indexed player, uint256 amount, uint256 newTotal)
func (_Vaultcontract *VaultcontractFilterer) ParseStakeDeposited(log types.Log) (*VaultcontractStakeDeposited, error) {
	event := new(VaultcontractStakeDeposited)
	if err := _Vaultcontract.contract.UnpackLog(event, "StakeDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
