package assets

 // AssetHolderETHBinRuntime is the runtime part of the compiled bytecode used for deploying new contracts.
var AssetHolderETHBinRuntime = `6080604052600436106100555760003560e01c80631de26e161461005a5780634ed4283c1461006f57806353c2ed8e1461008f578063ae9ee18c146100ba578063d945af1d146100e7578063fc79a66d14610114575b600080fd5b61006d610068366004610976565b610134565b005b34801561007b57600080fd5b5061006d61008a366004610997565b6101ac565b34801561009b57600080fd5b506100a461034e565b6040516100b19190610a4d565b60405180910390f35b3480156100c657600080fd5b506100da6100d53660046108e7565b61035d565b6040516100b19190610e0e565b3480156100f357600080fd5b506101076101023660046108e7565b61036f565b6040516100b19190610a61565b34801561012057600080fd5b5061006d61012f3660046108ff565b610384565b61013e82826105a2565b60008281526020819052604090205461015790826105c5565b60008381526020819052604090205561017082826105c1565b817fcd2fe07293de5928c5df9505b65a8d6506f8668dfe81af09090920687edc48a9826040516101a09190610e0e565b60405180910390a25050565b823560009081526001602052604090205460ff166101e55760405162461bcd60e51b81526004016101dc90610d50565b60405180910390fd5b61024d836040516020016101f99190610dc2565b60408051601f198184030181526020601f8601819004810284018101909252848352919085908590819084018382808284376000920191909152506102489250505060408701602088016108cb565b6105f1565b6102695760405162461bcd60e51b81526004016101dc90610cd7565b6000610285843561028060408701602088016108cb565b61062c565b600081815260208190526040902054909150606085013511156102ba5760405162461bcd60e51b81526004016101dc90610c19565b6102c584848461065f565b6000818152602081905260409020546102e2906060860135610664565b6000828152602081905260409020556102fc8484846106a6565b807fd0b6e7d0170f56c62f87de6a8a47a0ccf41c86ffb5084d399d8eb62e823f2a816060860180359061033290604089016108cb565b604051610340929190610a6c565b60405180910390a250505050565b6002546001600160a01b031681565b60006020819052908152604090205481565b60016020526000908152604090205460ff1681565b6002546001600160a01b031633146103ae5760405162461bcd60e51b81526004016101dc90610d7d565b8281146103cd5760405162461bcd60e51b81526004016101dc90610bd0565b60008581526001602052604090205460ff16156103fc5760405162461bcd60e51b81526004016101dc90610c50565b60008581526020819052604081208054908290559060608567ffffffffffffffff8111801561042a57600080fd5b50604051908082528060200260200182016040528015610454578160200160208202803683370190505b50905060005b868110156104fb5760006104898a8a8a8581811061047457fe5b905060200201602081019061028091906108cb565b90508083838151811061049857fe5b6020026020010181815250506104c960008083815260200190815260200160002054866105c590919063ffffffff16565b94506104f08787848181106104da57fe5b90506020020135856105c590919063ffffffff16565b93505060010161045a565b508183106105555760005b868110156105535785858281811061051a57fe5b9050602002013560008084848151811061053057fe5b602090810291909101810151825281019190915260400160002055600101610506565b505b6000888152600160208190526040808320805460ff19169092179091555189917fef898d6cd3395b6dfe67a3c1923e5c726c1b154e979fb0a25a9c41d0093168b891a25050505050505050565b8034146105c15760405162461bcd60e51b81526004016101dc90610b99565b5050565b6000828201838110156105ea5760405162461bcd60e51b81526004016101dc90610b62565b9392505050565b60008061060485805190602001206106f8565b905060006106128286610728565b6001600160a01b0390811690851614925050509392505050565b60008282604051602001610641929190610a6c565b60405160208183030381529060405280519060200120905092915050565b505050565b60006105ea83836040518060400160405280601e81526020017f536166654d6174683a207375627472616374696f6e206f766572666c6f770000815250610856565b6106b660608401604085016108cb565b6001600160a01b03166108fc84606001359081150290604051600060405180830381858888f193505050501580156106f2573d6000803e3d6000fd5b50505050565b60008160405160200161070b9190610a1c565b604051602081830303815290604052805190602001209050919050565b6000815160411461074b5760405162461bcd60e51b81526004016101dc90610b2b565b60208201516040830151606084015160001a7f7fffffffffffffffffffffffffffffff5d576e7357a4501ddfe92f46681b20a082111561079d5760405162461bcd60e51b81526004016101dc90610c95565b8060ff16601b141580156107b557508060ff16601c14155b156107d25760405162461bcd60e51b81526004016101dc90610d0e565b6000600187838686604051600081526020016040526040516107f79493929190610a83565b6020604051602081039080840390855afa158015610819573d6000803e3d6000fd5b5050604051601f1901519150506001600160a01b03811661084c5760405162461bcd60e51b81526004016101dc90610af4565b9695505050505050565b6000818484111561087a5760405162461bcd60e51b81526004016101dc9190610aa1565b505050900390565b60008083601f840112610893578182fd5b50813567ffffffffffffffff8111156108aa578182fd5b60208301915083602080830285010111156108c457600080fd5b9250929050565b6000602082840312156108dc578081fd5b81356105ea81610e17565b6000602082840312156108f8578081fd5b5035919050565b600080600080600060608688031215610916578081fd5b85359450602086013567ffffffffffffffff80821115610934578283fd5b61094089838a01610882565b90965094506040880135915080821115610958578283fd5b5061096588828901610882565b969995985093965092949392505050565b60008060408385031215610988578182fd5b50508035926020909101359150565b600080600083850360a08112156109ac578384fd5b60808112156109b9578384fd5b50839250608084013567ffffffffffffffff808211156109d7578384fd5b818601915086601f8301126109ea578384fd5b8135818111156109f8578485fd5b876020828501011115610a09578485fd5b6020830194508093505050509250925092565b7f19457468657265756d205369676e6564204d6573736167653a0a3332000000008152601c810191909152603c0190565b6001600160a01b0391909116815260200190565b901515815260200190565b9182526001600160a01b0316602082015260400190565b93845260ff9290921660208401526040830152606082015260800190565b6000602080835283518082850152825b81811015610acd57858101830151858201604001528201610ab1565b81811115610ade5783604083870101525b50601f01601f1916929092016040019392505050565b60208082526018908201527f45434453413a20696e76616c6964207369676e61747572650000000000000000604082015260600190565b6020808252601f908201527f45434453413a20696e76616c6964207369676e6174757265206c656e67746800604082015260600190565b6020808252601b908201527f536166654d6174683a206164646974696f6e206f766572666c6f770000000000604082015260600190565b6020808252601f908201527f77726f6e6720616d6f756e74206f662045544820666f72206465706f73697400604082015260600190565b60208082526029908201527f7061727469636970616e7473206c656e6774682073686f756c6420657175616c6040820152682062616c616e63657360b81b606082015260800190565b6020808252601f908201527f696e73756666696369656e742045544820666f72207769746864726177616c00604082015260600190565b60208082526025908201527f747279696e6720746f2073657420616c726561647920736574746c6564206368604082015264185b9b995b60da1b606082015260800190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202773272076616c604082015261756560f01b606082015260800190565b6020808252601d908201527f7369676e617475726520766572696669636174696f6e206661696c6564000000604082015260600190565b60208082526022908201527f45434453413a20696e76616c6964207369676e6174757265202776272076616c604082015261756560f01b606082015260800190565b60208082526013908201527218da185b9b995b081b9bdd081cd95d1d1b1959606a1b604082015260600190565b60208082526025908201527f63616e206f6e6c792062652063616c6c6564206279207468652061646a75646960408201526431b0ba37b960d91b606082015260800190565b81358152608081016020830135610dd881610e17565b6001600160a01b039081166020840152604084013590610df782610e17565b166040830152606092830135929091019190915290565b90815260200190565b6001600160a01b0381168114610e2c57600080fd5b5056fea2646970667358221220d6ba670268c776fefbc38023b4324c5f34b1c9178cd8e7fbda42553b155384e064736f6c63430007030033`
