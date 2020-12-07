typedef struct _Transaction{
	char *err;
	char *address;
	char *hash;
	double size;
	char *data;
} Transaction;

// NewLog is used for the receipt
typedef struct _NewLog{
	// The Removed field is true if this log was reverted due to a chain reorganisation.
	// You must pay attention to this field if you receive logs through a filter query.
	// 1 for true, 0 for false.
	int Removed;
	// index of the log in the block
	char *LogIndex;
	// index of the transaction in the block
	char *TransactionIndex;
	// hash of the transaction
	char *TransactionHash;
	// hash of the block in which the transaction was included
	char *BlockHash;
	// Derived fields. These fields are filled in by the node
	// but not secured by consensus.
	// block in which the transaction was included
	char *BlockNumber;
	// Consensus fields:
	// address of the contract that generated the event
	char *Address;
	// supplied by the contract, usually ABI-encoded
	char *Data;
	// Type for FISCO BCOS
	char *Type;
	// list of topics provided by the contract.
	char **Topics;
} NewLog;

// The receipt of an executed transaction.
typedef struct _Receipt{
    char *TransactionHash;
    char *TransactionIndex;
    char *BlockHash;
    char *BlockNumber;
    char *GasUsed;
    char *ContractAddress;
    char *Root;
    int Status;
    char *From;
    char *To;
    char *Input;
    char *Output;
    char *Logs;
    //NewLog *Logs;
    char *LogsBloom;
} Receipt;

typedef struct _TxExecResult{
    Transaction *tx;
    Receipt *receipt;
    char * err;
} TxExecResult;

typedef struct _callResult{
    char * result;
    char * err;
} CallResult;

typedef struct _rpcResult{
    char * result;
    char * err;
} RPCResult;

typedef struct _rpcReceiptResult{
    Receipt * receipt;
    char * err;
} RPCReceiptResult;

