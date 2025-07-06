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

// IVaultContractChainConfig is an auto generated low-level Go binding around an user-defined struct.
type IVaultContractChainConfig struct {
	DomainId           uint32
	UsdcAddress        common.Address
	TokenMessenger     common.Address
	MessageTransmitter common.Address
	IsSupported        bool
}

// VaultcontractMetaData contains all meta data concerning the Vaultcontract contract.
var VaultcontractMetaData = &bind.MetaData{
	ABI: "[{\"type\":\"constructor\",\"inputs\":[{\"name\":\"_authorizedBackend\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_usdcContractAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_tokenMessengerV2\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"_messageTransmitterV2\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"AUTHORIZED_BACKEND\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"MESSAGE_TRANSMITTER_V2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"TOKEN_MESSENGER_V2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"USDC_CONTRACT_ADDRESS\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"chainConfigs\",\"inputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"domainId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"usdcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenMessenger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"messageTransmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isSupported\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getAuthorizedBackend\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getChainConfig\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"tuple\",\"internalType\":\"structIVaultContract.ChainConfig\",\"components\":[{\"name\":\"domainId\",\"type\":\"uint32\",\"internalType\":\"uint32\"},{\"name\":\"usdcAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"tokenMessenger\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"messageTransmitter\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"isSupported\",\"type\":\"bool\",\"internalType\":\"bool\"}]}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getMessageTransmitterV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTokenMessengerV2\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getTotalStakes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"getUsdcContractAddress\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"address\",\"internalType\":\"address\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"isChainSupported\",\"inputs\":[{\"name\":\"chainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[{\"name\":\"\",\"type\":\"bool\",\"internalType\":\"bool\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"stake\",\"inputs\":[{\"name\":\"playerAddress\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"function\",\"name\":\"totalStakes\",\"inputs\":[],\"outputs\":[{\"name\":\"\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"stateMutability\":\"view\"},{\"type\":\"function\",\"name\":\"transferRewardsCrossChain\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"destinationChainId\",\"type\":\"uint256\",\"internalType\":\"uint256\"},{\"name\":\"recipient\",\"type\":\"address\",\"internalType\":\"address\"},{\"name\":\"useFastTransfer\",\"type\":\"bool\",\"internalType\":\"bool\"},{\"name\":\"maxFee\",\"type\":\"uint256\",\"internalType\":\"uint256\"}],\"outputs\":[],\"stateMutability\":\"nonpayable\"},{\"type\":\"event\",\"name\":\"CrossChainTransferInitiated\",\"inputs\":[{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"destinationDomain\",\"type\":\"uint32\",\"indexed\":false,\"internalType\":\"uint32\"},{\"name\":\"nonce\",\"type\":\"uint64\",\"indexed\":false,\"internalType\":\"uint64\"}],\"anonymous\":false},{\"type\":\"event\",\"name\":\"StakeDeposited\",\"inputs\":[{\"name\":\"playerAddress\",\"type\":\"address\",\"indexed\":true,\"internalType\":\"address\"},{\"name\":\"gameId\",\"type\":\"uint256\",\"indexed\":true,\"internalType\":\"uint256\"},{\"name\":\"amount\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"},{\"name\":\"newTotal\",\"type\":\"uint256\",\"indexed\":false,\"internalType\":\"uint256\"}],\"anonymous\":false}]",
	Bin: "0x6101006040523461103257604051601f611b3338819003918201601f19168301916001600160401b03831184841017611036578084926080946040528339810103126110325761004e81611069565b9061005b60208201611069565b610073606061006c60408501611069565b9301611069565b926001600160a01b03811615610fdb576001600160a01b03821615610f81576001600160a01b03831615610f2b576001600160a01b03841615610ed15760805260a05260c05260e0526100c461104a565b5f808252731c7d4b196cb0c7b01d743fbc6116a902379c723860208084019182525f516020611b135f395f51905f52604085019081525f516020611af35f395f51905f526060860190815260016080870181815262aa36a7909652835294517fd0ca2b9e48613413ba8bd488b2bc9928edc96f8f17ec2af94a885da5292afebb805494516001600160c01b031990951663ffffffff92909216919091179390921b600160201b600160c01b0316929092179055517fd0ca2b9e48613413ba8bd488b2bc9928edc96f8f17ec2af94a885da5292afebc80546001600160a01b0319166001600160a01b0392831617905591517fd0ca2b9e48613413ba8bd488b2bc9928edc96f8f17ec2af94a885da5292afebd805492516001600160a81b0319909316919093161790151560a01b60ff60a01b1617905561020261104a565b6001808252735425890298aed601595a70ab815c96711a31bc6560208084019182525f516020611b135f395f51905f52604085019081525f516020611af35f395f51905f52606086019081526080860185815261a8695f5294835294517f801d15fd97e52318f13586549c4f4b4ff53f255093eff97aaa42bd8b6e920db2805494516001600160c01b031990951663ffffffff92909216919091179390921b600160201b600160c01b0316929092179055517f801d15fd97e52318f13586549c4f4b4ff53f255093eff97aaa42bd8b6e920db380546001600160a01b0319166001600160a01b0392831617905591517f801d15fd97e52318f13586549c4f4b4ff53f255093eff97aaa42bd8b6e920db4805492516001600160a81b0319909316919093161790151560a01b60ff60a01b1617905561033e61104a565b60028152735fd84259d66cd46123540766be93dfe6d43130d760208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f526060850190815260016080860181815262aa37dc5f5290845294517f6d09c67303814678cfc65ff246f203a38d994ed1aa4deb4318c71dffd348fb98805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517f6d09c67303814678cfc65ff246f203a38d994ed1aa4deb4318c71dffd348fb9980546001600160a01b0319166001600160a01b0392831617905590517f6d09c67303814678cfc65ff246f203a38d994ed1aa4deb4318c71dffd348fb9a805493516001600160a81b0319909416919092161791151560a01b60ff60a01b1691909117905561048061104a565b600381527375faf114eafb1bdbe2f0316df893fd58ce46aa4d60208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f526060850190815260016080860181815262066eee5f5290845294517fc9f67f25c901b9944d4f155866de98249444082c5aef4a85958cd69c7d860ab1805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517fc9f67f25c901b9944d4f155866de98249444082c5aef4a85958cd69c7d860ab280546001600160a01b0319166001600160a01b0392831617905590517fc9f67f25c901b9944d4f155866de98249444082c5aef4a85958cd69c7d860ab3805493516001600160a81b0319909416919092161791151560a01b60ff60a01b169190911790556105c261104a565b6006815273036cbd53842c5426634e7929541ec2318f3dcf7e60208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f526060850190815260016080860181815262014a345f5290845294517ffc8758130b1edbec4992818204a8d94723427ce838142fe6d080d638bc744b5f805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517ffc8758130b1edbec4992818204a8d94723427ce838142fe6d080d638bc744b6080546001600160a01b0319166001600160a01b0392831617905590517ffc8758130b1edbec4992818204a8d94723427ce838142fe6d080d638bc744b61805493516001600160a81b0319909416919092161791151560a01b60ff60a01b1691909117905561070461104a565b600781527341e94eb019c0762f9bfcf9fb1e58725bfb0e758260208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f5260608501908152600160808601818152620138825f5290845294517feb83f3a4ee9fa292d004e4b2a68e5d43a38a279c3d324bf454d05a21a08b3272805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517feb83f3a4ee9fa292d004e4b2a68e5d43a38a279c3d324bf454d05a21a08b327380546001600160a01b0319166001600160a01b0392831617905590517feb83f3a4ee9fa292d004e4b2a68e5d43a38a279c3d324bf454d05a21a08b3274805493516001600160a81b0319909416919092161791151560a01b60ff60a01b1691909117905561084661104a565b600a81527331d0220469e10c4e71834a79b1f276d740d3768f60208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f52606085019081526001608086018181526105155f5290845294517f8ca952c4fcae17da8cd26e08b2d8c72a97eeebcc56c47768d7cb95d93aa5ffce805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517f8ca952c4fcae17da8cd26e08b2d8c72a97eeebcc56c47768d7cb95d93aa5ffcf80546001600160a01b0319166001600160a01b0392831617905590517f8ca952c4fcae17da8cd26e08b2d8c72a97eeebcc56c47768d7cb95d93aa5ffd0805493516001600160a81b0319909416919092161791151560a01b60ff60a01b1691909117905561098761104a565b600b815273fece4462d57bd51a6a552365a011b95f0e16d9b760208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f526060850190815260016080860181815261e7055f5290845294517f4288a2b46aaf833d0f54834fc6329a9bb31b2f0873fd1753b9434b5050bd21fc805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517f4288a2b46aaf833d0f54834fc6329a9bb31b2f0873fd1753b9434b5050bd21fd80546001600160a01b0319166001600160a01b0392831617905590517f4288a2b46aaf833d0f54834fc6329a9bb31b2f0873fd1753b9434b5050bd21fe805493516001600160a81b0319909416919092161791151560a01b60ff60a01b16919091179055610ac861104a565b600c8152736d7f141b6819c2c9cc2f818e6ad549e7ca090f8f60208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f52606085019081526001608086018181526204f5885f5290845294517fa1678fa44bfd55960b9dd4ee162c093cce67a601b142eda4734339890e44cc8a805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517fa1678fa44bfd55960b9dd4ee162c093cce67a601b142eda4734339890e44cc8b80546001600160a01b0319166001600160a01b0392831617905590517fa1678fa44bfd55960b9dd4ee162c093cce67a601b142eda4734339890e44cc8c805493516001600160a81b0319909416919092161791151560a01b60ff60a01b16919091179055610c0a61104a565b600d815273a4879fed32ecbef99399e5cbc247e533421c4ec660208083019182525f516020611b135f395f51905f52604084019081525f516020611af35f395f51905f526060850190815260016080860181815261faa55f5290845294517f36988d627d7a3d641e03a2a19859f5472b09421687db4b53a29b33e8d252508d805495516001600160c01b031990961663ffffffff92909216919091179490931b600160201b600160c01b031693909317909155517f36988d627d7a3d641e03a2a19859f5472b09421687db4b53a29b33e8d252508e80546001600160a01b0319166001600160a01b0392831617905590517f36988d627d7a3d641e03a2a19859f5472b09421687db4b53a29b33e8d252508f805493516001600160a81b0319909416919092161791151560a01b60ff60a01b16919091179055610d4b61104a565b600e81527366145f38cbac35ca6f1dfb4914df98f1614aea8860208083019182525f516020611b135f395f51905f5260408085019182525f516020611af35f395f51905f52606086019081526001608087018181526112c15f5290855295517f204118ecda1e33bd18f6e7950d6206e9bba11f1a703f99bdee5cb15a3f2f813a805496516001600160c01b031990971663ffffffff92909216919091179590941b600160201b600160c01b031694909417909255517f204118ecda1e33bd18f6e7950d6206e9bba11f1a703f99bdee5cb15a3f2f813b80546001600160a01b0319166001600160a01b0392831617905591517f204118ecda1e33bd18f6e7950d6206e9bba11f1a703f99bdee5cb15a3f2f813c805494516001600160a81b0319909516919093161792151560a01b60ff60a01b1692909217905551610a75908161107e823960805181818160fb015261091e015260a0518181816101ed0152818161075201526109a6015260c0518181816101bc01526108da015260e051816109620152f35b60405162461bcd60e51b815260206004820152602c60248201527f4d6573736167655472616e736d69747465722056322063616e6e6f742062652060448201526b7a65726f206164647265737360a01b6064820152608490fd5b60405162461bcd60e51b815260206004820152602860248201527f546f6b656e4d657373656e6765722056322063616e6e6f74206265207a65726f604482015267206164647265737360c01b6064820152608490fd5b60405162461bcd60e51b815260206004820152602c60248201527f5553444320636f6e747261637420616464726573732063616e6e6f742062652060448201526b7a65726f206164647265737360a01b6064820152608490fd5b60405162461bcd60e51b815260206004820152602960248201527f417574686f72697a6564206261636b656e642063616e6e6f74206265207a65726044820152686f206164647265737360b81b6064820152608490fd5b5f80fd5b634e487b7160e01b5f52604160045260245ffd5b6040519060a082016001600160401b0381118382101761103657604052565b51906001600160a01b03821682036110325756fe60806040526004361015610011575f80fd5b5f3560e01c8063085eac90146108c55780630c51b88f146106f2578063160c2b4b1461061357806319ed16dc1461061d5780631c4cc6511461061857806329513693146105665780633d932ea11461061357806351bfbc32146105a05780635221c1f01461056b57806368c336271461054a578063780e956214610566578063bf9befb11461054a578063d05c346f146100c3578063e24663cd146100be5763e465085a146100be575f80fd5b610991565b346103455760c0366003190112610345576024356064356001600160a01b0381169081900361034557608435908115158203610345577f00000000000000000000000000000000000000000000000000000000000000006001600160a01b031633036104f057821561049f57825f541061045a578015610416576044355f52600160205260405f209160ff60026040519461015d866109d5565b805463ffffffff811687526001600160a01b03602091821c811691880191909152600182015481166040880152910154908116606086015260a01c1615801560808501526103d15760405163095ea7b360e01b81526001600160a01b037f0000000000000000000000000000000000000000000000000000000000000000811660048301819052602483018790527f000000000000000000000000000000000000000000000000000000000000000090911693909290916020816044815f895af1908115610351575f916103a2575b50156103665760209260e4911561035c575f6103e8915b61ffff63ffffffff895116936040519889978896634701287760e11b88528d600489015260248801526044870152606486015283608486015260a43560a48601521660c48401525af1908115610351575f9161030a575b505f548381039081116102f65767ffffffffffffffff9263ffffffff915f55511660405193845260208401521660408201527ff1af501d0d14e1dbb3cf9b16e85de8105b65851201576ed8b4566dd1f0d9e9dd606060043592a2005b634e487b7160e01b5f52601160045260245ffd5b90506020813d602011610349575b8161032560209383610a05565b81010312610345575167ffffffffffffffff81168103610345578361029a565b5f80fd5b3d9150610318565b6040513d5f823e3d90fd5b5f6107d091610243565b60405162461bcd60e51b81526020600482015260146024820152731554d110c8185c1c1c9bdd985b0819985a5b195960621b6044820152606490fd5b6103c4915060203d6020116103ca575b6103bc8183610a05565b810190610a27565b8761022c565b503d6103b2565b60405162461bcd60e51b815260206004820152601f60248201527f44657374696e6174696f6e20636861696e206e6f7420737570706f72746564006044820152606490fd5b606460405162461bcd60e51b815260206004820152602060248201527f526563697069656e742063616e6e6f74206265207a65726f20616464726573736044820152fd5b60405162461bcd60e51b815260206004820152601960248201527f496e73756666696369656e7420746f74616c207374616b6573000000000000006044820152606490fd5b60405162461bcd60e51b8152602060048201526024808201527f52657761726420616d6f756e74206d75737420626520677265617465722074686044820152630616e20360e41b6064820152608490fd5b60405162461bcd60e51b815260206004820152602c60248201527f4f6e6c7920617574686f72697a6564206261636b656e642063616e207472616e60448201526b73666572207265776172647360a01b6064820152608490fd5b34610345575f3660031901126103455760205f54604051908152f35b61094d565b34610345576020366003190112610345576004355f526001602052602060ff600260405f20015460a01c166040519015158152f35b34610345576020366003190112610345576004355f52600160205260a060405f2060ff8154916002600180861b03600183015416910154906040519363ffffffff81168552600180871b039060201c1660208501526040840152600180851b0381166060840152831c1615156080820152f35b610909565b6108c5565b34610345576020366003190112610345575f608060405161063d816109d5565b82815282602082015282604082015282606082015201526004355f52600160205260a060405f20604051610670816109d5565b815463ffffffff8116808352602091821c5f19600180881b9190910191821684860190815290860154821660408087019182526002909701548084166060808901918252918a1c60ff16151560809889019081528951968752935185169686019690965290518316968401969096529251169381019390935251151590820152f35b34610345576060366003190112610345576004356001600160a01b038116908190036103455760443590811561087457801561082f576040516323b872dd60e01b815260048101829052306024820152604481018390526020816064815f7f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03165af1908115610351575f91610810575b50156107d4575f548281018091116102f657805f5560405192835260208301527feb6032fe8d7a2f7003e5b33948bdc807f994be970351cd021a033784812c8baa604060243593a3005b60405162461bcd60e51b81526020600482015260146024820152731554d110c81d1c985b9cd9995c8819985a5b195960621b6044820152606490fd5b610829915060203d6020116103ca576103bc8183610a05565b8361078a565b60405162461bcd60e51b815260206004820152601d60248201527f506c6179657220616464726573732063616e6e6f74206265207a65726f0000006044820152606490fd5b60405162461bcd60e51b815260206004820152602360248201527f5374616b6520616d6f756e74206d75737420626520677265617465722074686160448201526206e20360ec1b6064820152608490fd5b34610345575f366003190112610345576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b34610345575f366003190112610345576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b34610345575f366003190112610345576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b34610345575f366003190112610345576040517f00000000000000000000000000000000000000000000000000000000000000006001600160a01b03168152602090f35b60a0810190811067ffffffffffffffff8211176109f157604052565b634e487b7160e01b5f52604160045260245ffd5b90601f8019910116810190811067ffffffffffffffff8211176109f157604052565b9081602091031261034557518015158103610345579056fea2646970667358221220891b8959758f56c4506089b56395a055a88287550b9768ba150e4c422aaf1ba464736f6c634300081e0033000000000000000000000000e737e5cebeeba77efe34d4aa090756590b1ce2750000000000000000000000008fe6b999dc680ccfdd5bf7eb0974218be2542daa",
}

// VaultcontractABI is the input ABI used to generate the binding from.
// Deprecated: Use VaultcontractMetaData.ABI instead.
var VaultcontractABI = VaultcontractMetaData.ABI

// VaultcontractBin is the compiled bytecode used for deploying new contracts.
// Deprecated: Use VaultcontractMetaData.Bin instead.
var VaultcontractBin = VaultcontractMetaData.Bin

// DeployVaultcontract deploys a new Ethereum contract, binding an instance of Vaultcontract to it.
func DeployVaultcontract(auth *bind.TransactOpts, backend bind.ContractBackend, _authorizedBackend common.Address, _usdcContractAddress common.Address, _tokenMessengerV2 common.Address, _messageTransmitterV2 common.Address) (common.Address, *types.Transaction, *Vaultcontract, error) {
	parsed, err := VaultcontractMetaData.GetAbi()
	if err != nil {
		return common.Address{}, nil, nil, err
	}
	if parsed == nil {
		return common.Address{}, nil, nil, errors.New("GetABI returned nil")
	}

	address, tx, contract, err := bind.DeployContract(auth, *parsed, common.FromHex(VaultcontractBin), backend, _authorizedBackend, _usdcContractAddress, _tokenMessengerV2, _messageTransmitterV2)
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

// AUTHORIZEDBACKEND is a free data retrieval call binding the contract method 0x160c2b4b.
//
// Solidity: function AUTHORIZED_BACKEND() view returns(address)
func (_Vaultcontract *VaultcontractCaller) AUTHORIZEDBACKEND(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "AUTHORIZED_BACKEND")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// AUTHORIZEDBACKEND is a free data retrieval call binding the contract method 0x160c2b4b.
//
// Solidity: function AUTHORIZED_BACKEND() view returns(address)
func (_Vaultcontract *VaultcontractSession) AUTHORIZEDBACKEND() (common.Address, error) {
	return _Vaultcontract.Contract.AUTHORIZEDBACKEND(&_Vaultcontract.CallOpts)
}

// AUTHORIZEDBACKEND is a free data retrieval call binding the contract method 0x160c2b4b.
//
// Solidity: function AUTHORIZED_BACKEND() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) AUTHORIZEDBACKEND() (common.Address, error) {
	return _Vaultcontract.Contract.AUTHORIZEDBACKEND(&_Vaultcontract.CallOpts)
}

// MESSAGETRANSMITTERV2 is a free data retrieval call binding the contract method 0x29513693.
//
// Solidity: function MESSAGE_TRANSMITTER_V2() view returns(address)
func (_Vaultcontract *VaultcontractCaller) MESSAGETRANSMITTERV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "MESSAGE_TRANSMITTER_V2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// MESSAGETRANSMITTERV2 is a free data retrieval call binding the contract method 0x29513693.
//
// Solidity: function MESSAGE_TRANSMITTER_V2() view returns(address)
func (_Vaultcontract *VaultcontractSession) MESSAGETRANSMITTERV2() (common.Address, error) {
	return _Vaultcontract.Contract.MESSAGETRANSMITTERV2(&_Vaultcontract.CallOpts)
}

// MESSAGETRANSMITTERV2 is a free data retrieval call binding the contract method 0x29513693.
//
// Solidity: function MESSAGE_TRANSMITTER_V2() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) MESSAGETRANSMITTERV2() (common.Address, error) {
	return _Vaultcontract.Contract.MESSAGETRANSMITTERV2(&_Vaultcontract.CallOpts)
}

// TOKENMESSENGERV2 is a free data retrieval call binding the contract method 0x085eac90.
//
// Solidity: function TOKEN_MESSENGER_V2() view returns(address)
func (_Vaultcontract *VaultcontractCaller) TOKENMESSENGERV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "TOKEN_MESSENGER_V2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// TOKENMESSENGERV2 is a free data retrieval call binding the contract method 0x085eac90.
//
// Solidity: function TOKEN_MESSENGER_V2() view returns(address)
func (_Vaultcontract *VaultcontractSession) TOKENMESSENGERV2() (common.Address, error) {
	return _Vaultcontract.Contract.TOKENMESSENGERV2(&_Vaultcontract.CallOpts)
}

// TOKENMESSENGERV2 is a free data retrieval call binding the contract method 0x085eac90.
//
// Solidity: function TOKEN_MESSENGER_V2() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) TOKENMESSENGERV2() (common.Address, error) {
	return _Vaultcontract.Contract.TOKENMESSENGERV2(&_Vaultcontract.CallOpts)
}

// USDCCONTRACTADDRESS is a free data retrieval call binding the contract method 0xe24663cd.
//
// Solidity: function USDC_CONTRACT_ADDRESS() view returns(address)
func (_Vaultcontract *VaultcontractCaller) USDCCONTRACTADDRESS(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "USDC_CONTRACT_ADDRESS")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// USDCCONTRACTADDRESS is a free data retrieval call binding the contract method 0xe24663cd.
//
// Solidity: function USDC_CONTRACT_ADDRESS() view returns(address)
func (_Vaultcontract *VaultcontractSession) USDCCONTRACTADDRESS() (common.Address, error) {
	return _Vaultcontract.Contract.USDCCONTRACTADDRESS(&_Vaultcontract.CallOpts)
}

// USDCCONTRACTADDRESS is a free data retrieval call binding the contract method 0xe24663cd.
//
// Solidity: function USDC_CONTRACT_ADDRESS() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) USDCCONTRACTADDRESS() (common.Address, error) {
	return _Vaultcontract.Contract.USDCCONTRACTADDRESS(&_Vaultcontract.CallOpts)
}

// ChainConfigs is a free data retrieval call binding the contract method 0x51bfbc32.
//
// Solidity: function chainConfigs(uint256 ) view returns(uint32 domainId, address usdcAddress, address tokenMessenger, address messageTransmitter, bool isSupported)
func (_Vaultcontract *VaultcontractCaller) ChainConfigs(opts *bind.CallOpts, arg0 *big.Int) (struct {
	DomainId           uint32
	UsdcAddress        common.Address
	TokenMessenger     common.Address
	MessageTransmitter common.Address
	IsSupported        bool
}, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "chainConfigs", arg0)

	outstruct := new(struct {
		DomainId           uint32
		UsdcAddress        common.Address
		TokenMessenger     common.Address
		MessageTransmitter common.Address
		IsSupported        bool
	})
	if err != nil {
		return *outstruct, err
	}

	outstruct.DomainId = *abi.ConvertType(out[0], new(uint32)).(*uint32)
	outstruct.UsdcAddress = *abi.ConvertType(out[1], new(common.Address)).(*common.Address)
	outstruct.TokenMessenger = *abi.ConvertType(out[2], new(common.Address)).(*common.Address)
	outstruct.MessageTransmitter = *abi.ConvertType(out[3], new(common.Address)).(*common.Address)
	outstruct.IsSupported = *abi.ConvertType(out[4], new(bool)).(*bool)

	return *outstruct, err

}

// ChainConfigs is a free data retrieval call binding the contract method 0x51bfbc32.
//
// Solidity: function chainConfigs(uint256 ) view returns(uint32 domainId, address usdcAddress, address tokenMessenger, address messageTransmitter, bool isSupported)
func (_Vaultcontract *VaultcontractSession) ChainConfigs(arg0 *big.Int) (struct {
	DomainId           uint32
	UsdcAddress        common.Address
	TokenMessenger     common.Address
	MessageTransmitter common.Address
	IsSupported        bool
}, error) {
	return _Vaultcontract.Contract.ChainConfigs(&_Vaultcontract.CallOpts, arg0)
}

// ChainConfigs is a free data retrieval call binding the contract method 0x51bfbc32.
//
// Solidity: function chainConfigs(uint256 ) view returns(uint32 domainId, address usdcAddress, address tokenMessenger, address messageTransmitter, bool isSupported)
func (_Vaultcontract *VaultcontractCallerSession) ChainConfigs(arg0 *big.Int) (struct {
	DomainId           uint32
	UsdcAddress        common.Address
	TokenMessenger     common.Address
	MessageTransmitter common.Address
	IsSupported        bool
}, error) {
	return _Vaultcontract.Contract.ChainConfigs(&_Vaultcontract.CallOpts, arg0)
}

// GetAuthorizedBackend is a free data retrieval call binding the contract method 0x3d932ea1.
//
// Solidity: function getAuthorizedBackend() view returns(address)
func (_Vaultcontract *VaultcontractCaller) GetAuthorizedBackend(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getAuthorizedBackend")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetAuthorizedBackend is a free data retrieval call binding the contract method 0x3d932ea1.
//
// Solidity: function getAuthorizedBackend() view returns(address)
func (_Vaultcontract *VaultcontractSession) GetAuthorizedBackend() (common.Address, error) {
	return _Vaultcontract.Contract.GetAuthorizedBackend(&_Vaultcontract.CallOpts)
}

// GetAuthorizedBackend is a free data retrieval call binding the contract method 0x3d932ea1.
//
// Solidity: function getAuthorizedBackend() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) GetAuthorizedBackend() (common.Address, error) {
	return _Vaultcontract.Contract.GetAuthorizedBackend(&_Vaultcontract.CallOpts)
}

// GetChainConfig is a free data retrieval call binding the contract method 0x19ed16dc.
//
// Solidity: function getChainConfig(uint256 chainId) view returns((uint32,address,address,address,bool))
func (_Vaultcontract *VaultcontractCaller) GetChainConfig(opts *bind.CallOpts, chainId *big.Int) (IVaultContractChainConfig, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getChainConfig", chainId)

	if err != nil {
		return *new(IVaultContractChainConfig), err
	}

	out0 := *abi.ConvertType(out[0], new(IVaultContractChainConfig)).(*IVaultContractChainConfig)

	return out0, err

}

// GetChainConfig is a free data retrieval call binding the contract method 0x19ed16dc.
//
// Solidity: function getChainConfig(uint256 chainId) view returns((uint32,address,address,address,bool))
func (_Vaultcontract *VaultcontractSession) GetChainConfig(chainId *big.Int) (IVaultContractChainConfig, error) {
	return _Vaultcontract.Contract.GetChainConfig(&_Vaultcontract.CallOpts, chainId)
}

// GetChainConfig is a free data retrieval call binding the contract method 0x19ed16dc.
//
// Solidity: function getChainConfig(uint256 chainId) view returns((uint32,address,address,address,bool))
func (_Vaultcontract *VaultcontractCallerSession) GetChainConfig(chainId *big.Int) (IVaultContractChainConfig, error) {
	return _Vaultcontract.Contract.GetChainConfig(&_Vaultcontract.CallOpts, chainId)
}

// GetMessageTransmitterV2 is a free data retrieval call binding the contract method 0x780e9562.
//
// Solidity: function getMessageTransmitterV2() view returns(address)
func (_Vaultcontract *VaultcontractCaller) GetMessageTransmitterV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getMessageTransmitterV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetMessageTransmitterV2 is a free data retrieval call binding the contract method 0x780e9562.
//
// Solidity: function getMessageTransmitterV2() view returns(address)
func (_Vaultcontract *VaultcontractSession) GetMessageTransmitterV2() (common.Address, error) {
	return _Vaultcontract.Contract.GetMessageTransmitterV2(&_Vaultcontract.CallOpts)
}

// GetMessageTransmitterV2 is a free data retrieval call binding the contract method 0x780e9562.
//
// Solidity: function getMessageTransmitterV2() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) GetMessageTransmitterV2() (common.Address, error) {
	return _Vaultcontract.Contract.GetMessageTransmitterV2(&_Vaultcontract.CallOpts)
}

// GetTokenMessengerV2 is a free data retrieval call binding the contract method 0x1c4cc651.
//
// Solidity: function getTokenMessengerV2() view returns(address)
func (_Vaultcontract *VaultcontractCaller) GetTokenMessengerV2(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getTokenMessengerV2")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetTokenMessengerV2 is a free data retrieval call binding the contract method 0x1c4cc651.
//
// Solidity: function getTokenMessengerV2() view returns(address)
func (_Vaultcontract *VaultcontractSession) GetTokenMessengerV2() (common.Address, error) {
	return _Vaultcontract.Contract.GetTokenMessengerV2(&_Vaultcontract.CallOpts)
}

// GetTokenMessengerV2 is a free data retrieval call binding the contract method 0x1c4cc651.
//
// Solidity: function getTokenMessengerV2() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) GetTokenMessengerV2() (common.Address, error) {
	return _Vaultcontract.Contract.GetTokenMessengerV2(&_Vaultcontract.CallOpts)
}

// GetTotalStakes is a free data retrieval call binding the contract method 0x68c33627.
//
// Solidity: function getTotalStakes() view returns(uint256)
func (_Vaultcontract *VaultcontractCaller) GetTotalStakes(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getTotalStakes")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// GetTotalStakes is a free data retrieval call binding the contract method 0x68c33627.
//
// Solidity: function getTotalStakes() view returns(uint256)
func (_Vaultcontract *VaultcontractSession) GetTotalStakes() (*big.Int, error) {
	return _Vaultcontract.Contract.GetTotalStakes(&_Vaultcontract.CallOpts)
}

// GetTotalStakes is a free data retrieval call binding the contract method 0x68c33627.
//
// Solidity: function getTotalStakes() view returns(uint256)
func (_Vaultcontract *VaultcontractCallerSession) GetTotalStakes() (*big.Int, error) {
	return _Vaultcontract.Contract.GetTotalStakes(&_Vaultcontract.CallOpts)
}

// GetUsdcContractAddress is a free data retrieval call binding the contract method 0xe465085a.
//
// Solidity: function getUsdcContractAddress() view returns(address)
func (_Vaultcontract *VaultcontractCaller) GetUsdcContractAddress(opts *bind.CallOpts) (common.Address, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "getUsdcContractAddress")

	if err != nil {
		return *new(common.Address), err
	}

	out0 := *abi.ConvertType(out[0], new(common.Address)).(*common.Address)

	return out0, err

}

// GetUsdcContractAddress is a free data retrieval call binding the contract method 0xe465085a.
//
// Solidity: function getUsdcContractAddress() view returns(address)
func (_Vaultcontract *VaultcontractSession) GetUsdcContractAddress() (common.Address, error) {
	return _Vaultcontract.Contract.GetUsdcContractAddress(&_Vaultcontract.CallOpts)
}

// GetUsdcContractAddress is a free data retrieval call binding the contract method 0xe465085a.
//
// Solidity: function getUsdcContractAddress() view returns(address)
func (_Vaultcontract *VaultcontractCallerSession) GetUsdcContractAddress() (common.Address, error) {
	return _Vaultcontract.Contract.GetUsdcContractAddress(&_Vaultcontract.CallOpts)
}

// IsChainSupported is a free data retrieval call binding the contract method 0x5221c1f0.
//
// Solidity: function isChainSupported(uint256 chainId) view returns(bool)
func (_Vaultcontract *VaultcontractCaller) IsChainSupported(opts *bind.CallOpts, chainId *big.Int) (bool, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "isChainSupported", chainId)

	if err != nil {
		return *new(bool), err
	}

	out0 := *abi.ConvertType(out[0], new(bool)).(*bool)

	return out0, err

}

// IsChainSupported is a free data retrieval call binding the contract method 0x5221c1f0.
//
// Solidity: function isChainSupported(uint256 chainId) view returns(bool)
func (_Vaultcontract *VaultcontractSession) IsChainSupported(chainId *big.Int) (bool, error) {
	return _Vaultcontract.Contract.IsChainSupported(&_Vaultcontract.CallOpts, chainId)
}

// IsChainSupported is a free data retrieval call binding the contract method 0x5221c1f0.
//
// Solidity: function isChainSupported(uint256 chainId) view returns(bool)
func (_Vaultcontract *VaultcontractCallerSession) IsChainSupported(chainId *big.Int) (bool, error) {
	return _Vaultcontract.Contract.IsChainSupported(&_Vaultcontract.CallOpts, chainId)
}

// TotalStakes is a free data retrieval call binding the contract method 0xbf9befb1.
//
// Solidity: function totalStakes() view returns(uint256)
func (_Vaultcontract *VaultcontractCaller) TotalStakes(opts *bind.CallOpts) (*big.Int, error) {
	var out []interface{}
	err := _Vaultcontract.contract.Call(opts, &out, "totalStakes")

	if err != nil {
		return *new(*big.Int), err
	}

	out0 := *abi.ConvertType(out[0], new(*big.Int)).(**big.Int)

	return out0, err

}

// TotalStakes is a free data retrieval call binding the contract method 0xbf9befb1.
//
// Solidity: function totalStakes() view returns(uint256)
func (_Vaultcontract *VaultcontractSession) TotalStakes() (*big.Int, error) {
	return _Vaultcontract.Contract.TotalStakes(&_Vaultcontract.CallOpts)
}

// TotalStakes is a free data retrieval call binding the contract method 0xbf9befb1.
//
// Solidity: function totalStakes() view returns(uint256)
func (_Vaultcontract *VaultcontractCallerSession) TotalStakes() (*big.Int, error) {
	return _Vaultcontract.Contract.TotalStakes(&_Vaultcontract.CallOpts)
}

// Stake is a paid mutator transaction binding the contract method 0x0c51b88f.
//
// Solidity: function stake(address playerAddress, uint256 gameId, uint256 amount) returns()
func (_Vaultcontract *VaultcontractTransactor) Stake(opts *bind.TransactOpts, playerAddress common.Address, gameId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.contract.Transact(opts, "stake", playerAddress, gameId, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x0c51b88f.
//
// Solidity: function stake(address playerAddress, uint256 gameId, uint256 amount) returns()
func (_Vaultcontract *VaultcontractSession) Stake(playerAddress common.Address, gameId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.Stake(&_Vaultcontract.TransactOpts, playerAddress, gameId, amount)
}

// Stake is a paid mutator transaction binding the contract method 0x0c51b88f.
//
// Solidity: function stake(address playerAddress, uint256 gameId, uint256 amount) returns()
func (_Vaultcontract *VaultcontractTransactorSession) Stake(playerAddress common.Address, gameId *big.Int, amount *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.Stake(&_Vaultcontract.TransactOpts, playerAddress, gameId, amount)
}

// TransferRewardsCrossChain is a paid mutator transaction binding the contract method 0xd05c346f.
//
// Solidity: function transferRewardsCrossChain(uint256 gameId, uint256 amount, uint256 destinationChainId, address recipient, bool useFastTransfer, uint256 maxFee) returns()
func (_Vaultcontract *VaultcontractTransactor) TransferRewardsCrossChain(opts *bind.TransactOpts, gameId *big.Int, amount *big.Int, destinationChainId *big.Int, recipient common.Address, useFastTransfer bool, maxFee *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.contract.Transact(opts, "transferRewardsCrossChain", gameId, amount, destinationChainId, recipient, useFastTransfer, maxFee)
}

// TransferRewardsCrossChain is a paid mutator transaction binding the contract method 0xd05c346f.
//
// Solidity: function transferRewardsCrossChain(uint256 gameId, uint256 amount, uint256 destinationChainId, address recipient, bool useFastTransfer, uint256 maxFee) returns()
func (_Vaultcontract *VaultcontractSession) TransferRewardsCrossChain(gameId *big.Int, amount *big.Int, destinationChainId *big.Int, recipient common.Address, useFastTransfer bool, maxFee *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.TransferRewardsCrossChain(&_Vaultcontract.TransactOpts, gameId, amount, destinationChainId, recipient, useFastTransfer, maxFee)
}

// TransferRewardsCrossChain is a paid mutator transaction binding the contract method 0xd05c346f.
//
// Solidity: function transferRewardsCrossChain(uint256 gameId, uint256 amount, uint256 destinationChainId, address recipient, bool useFastTransfer, uint256 maxFee) returns()
func (_Vaultcontract *VaultcontractTransactorSession) TransferRewardsCrossChain(gameId *big.Int, amount *big.Int, destinationChainId *big.Int, recipient common.Address, useFastTransfer bool, maxFee *big.Int) (*types.Transaction, error) {
	return _Vaultcontract.Contract.TransferRewardsCrossChain(&_Vaultcontract.TransactOpts, gameId, amount, destinationChainId, recipient, useFastTransfer, maxFee)
}

// VaultcontractCrossChainTransferInitiatedIterator is returned from FilterCrossChainTransferInitiated and is used to iterate over the raw logs and unpacked data for CrossChainTransferInitiated events raised by the Vaultcontract contract.
type VaultcontractCrossChainTransferInitiatedIterator struct {
	Event *VaultcontractCrossChainTransferInitiated // Event containing the contract specifics and raw log

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
func (it *VaultcontractCrossChainTransferInitiatedIterator) Next() bool {
	// If the iterator failed, stop iterating
	if it.fail != nil {
		return false
	}
	// If the iterator completed, deliver directly whatever's available
	if it.done {
		select {
		case log := <-it.logs:
			it.Event = new(VaultcontractCrossChainTransferInitiated)
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
		it.Event = new(VaultcontractCrossChainTransferInitiated)
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
func (it *VaultcontractCrossChainTransferInitiatedIterator) Error() error {
	return it.fail
}

// Close terminates the iteration process, releasing any pending underlying
// resources.
func (it *VaultcontractCrossChainTransferInitiatedIterator) Close() error {
	it.sub.Unsubscribe()
	return nil
}

// VaultcontractCrossChainTransferInitiated represents a CrossChainTransferInitiated event raised by the Vaultcontract contract.
type VaultcontractCrossChainTransferInitiated struct {
	GameId            *big.Int
	Amount            *big.Int
	DestinationDomain uint32
	Nonce             uint64
	Raw               types.Log // Blockchain specific contextual infos
}

// FilterCrossChainTransferInitiated is a free log retrieval operation binding the contract event 0xf1af501d0d14e1dbb3cf9b16e85de8105b65851201576ed8b4566dd1f0d9e9dd.
//
// Solidity: event CrossChainTransferInitiated(uint256 indexed gameId, uint256 amount, uint32 destinationDomain, uint64 nonce)
func (_Vaultcontract *VaultcontractFilterer) FilterCrossChainTransferInitiated(opts *bind.FilterOpts, gameId []*big.Int) (*VaultcontractCrossChainTransferInitiatedIterator, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Vaultcontract.contract.FilterLogs(opts, "CrossChainTransferInitiated", gameIdRule)
	if err != nil {
		return nil, err
	}
	return &VaultcontractCrossChainTransferInitiatedIterator{contract: _Vaultcontract.contract, event: "CrossChainTransferInitiated", logs: logs, sub: sub}, nil
}

// WatchCrossChainTransferInitiated is a free log subscription operation binding the contract event 0xf1af501d0d14e1dbb3cf9b16e85de8105b65851201576ed8b4566dd1f0d9e9dd.
//
// Solidity: event CrossChainTransferInitiated(uint256 indexed gameId, uint256 amount, uint32 destinationDomain, uint64 nonce)
func (_Vaultcontract *VaultcontractFilterer) WatchCrossChainTransferInitiated(opts *bind.WatchOpts, sink chan<- *VaultcontractCrossChainTransferInitiated, gameId []*big.Int) (event.Subscription, error) {

	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Vaultcontract.contract.WatchLogs(opts, "CrossChainTransferInitiated", gameIdRule)
	if err != nil {
		return nil, err
	}
	return event.NewSubscription(func(quit <-chan struct{}) error {
		defer sub.Unsubscribe()
		for {
			select {
			case log := <-logs:
				// New log arrived, parse the event and forward to the user
				event := new(VaultcontractCrossChainTransferInitiated)
				if err := _Vaultcontract.contract.UnpackLog(event, "CrossChainTransferInitiated", log); err != nil {
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

// ParseCrossChainTransferInitiated is a log parse operation binding the contract event 0xf1af501d0d14e1dbb3cf9b16e85de8105b65851201576ed8b4566dd1f0d9e9dd.
//
// Solidity: event CrossChainTransferInitiated(uint256 indexed gameId, uint256 amount, uint32 destinationDomain, uint64 nonce)
func (_Vaultcontract *VaultcontractFilterer) ParseCrossChainTransferInitiated(log types.Log) (*VaultcontractCrossChainTransferInitiated, error) {
	event := new(VaultcontractCrossChainTransferInitiated)
	if err := _Vaultcontract.contract.UnpackLog(event, "CrossChainTransferInitiated", log); err != nil {
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
	PlayerAddress common.Address
	GameId        *big.Int
	Amount        *big.Int
	NewTotal      *big.Int
	Raw           types.Log // Blockchain specific contextual infos
}

// FilterStakeDeposited is a free log retrieval operation binding the contract event 0xeb6032fe8d7a2f7003e5b33948bdc807f994be970351cd021a033784812c8baa.
//
// Solidity: event StakeDeposited(address indexed playerAddress, uint256 indexed gameId, uint256 amount, uint256 newTotal)
func (_Vaultcontract *VaultcontractFilterer) FilterStakeDeposited(opts *bind.FilterOpts, playerAddress []common.Address, gameId []*big.Int) (*VaultcontractStakeDepositedIterator, error) {

	var playerAddressRule []interface{}
	for _, playerAddressItem := range playerAddress {
		playerAddressRule = append(playerAddressRule, playerAddressItem)
	}
	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Vaultcontract.contract.FilterLogs(opts, "StakeDeposited", playerAddressRule, gameIdRule)
	if err != nil {
		return nil, err
	}
	return &VaultcontractStakeDepositedIterator{contract: _Vaultcontract.contract, event: "StakeDeposited", logs: logs, sub: sub}, nil
}

// WatchStakeDeposited is a free log subscription operation binding the contract event 0xeb6032fe8d7a2f7003e5b33948bdc807f994be970351cd021a033784812c8baa.
//
// Solidity: event StakeDeposited(address indexed playerAddress, uint256 indexed gameId, uint256 amount, uint256 newTotal)
func (_Vaultcontract *VaultcontractFilterer) WatchStakeDeposited(opts *bind.WatchOpts, sink chan<- *VaultcontractStakeDeposited, playerAddress []common.Address, gameId []*big.Int) (event.Subscription, error) {

	var playerAddressRule []interface{}
	for _, playerAddressItem := range playerAddress {
		playerAddressRule = append(playerAddressRule, playerAddressItem)
	}
	var gameIdRule []interface{}
	for _, gameIdItem := range gameId {
		gameIdRule = append(gameIdRule, gameIdItem)
	}

	logs, sub, err := _Vaultcontract.contract.WatchLogs(opts, "StakeDeposited", playerAddressRule, gameIdRule)
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

// ParseStakeDeposited is a log parse operation binding the contract event 0xeb6032fe8d7a2f7003e5b33948bdc807f994be970351cd021a033784812c8baa.
//
// Solidity: event StakeDeposited(address indexed playerAddress, uint256 indexed gameId, uint256 amount, uint256 newTotal)
func (_Vaultcontract *VaultcontractFilterer) ParseStakeDeposited(log types.Log) (*VaultcontractStakeDeposited, error) {
	event := new(VaultcontractStakeDeposited)
	if err := _Vaultcontract.contract.UnpackLog(event, "StakeDeposited", log); err != nil {
		return nil, err
	}
	event.Raw = log
	return event, nil
}
