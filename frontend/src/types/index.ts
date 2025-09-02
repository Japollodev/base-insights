// Common types used across the application

export interface Stats {
  blocksProcessed: number;
  transactionsProcessed: number;
  errors: number;
  isRunning: boolean;
  lastUpdate: string;
}

export interface BlockchainStatus {
  isConnected: boolean;
  chainId: string;
  rpcUrl: string;
  lastBlockNumber: number;
}

export interface ChartData {
  time: string;
  blocks: number;
  transactions: number;
}

export interface Block {
  number: number;
  hash: string;
  timestamp: number;
  transactions: Transaction[];
  gasUsed: number;
  gasLimit: number;
}

export interface Transaction {
  hash: string;
  from: string;
  to: string;
  value: string;
  gas: number;
}

export interface GasPrice {
  gasPrice: string;
  gasPriceGwei: string;
  maxFeePerGas?: string;
  maxPriorityFeePerGas?: string;
  baseFeePerGas?: string;
}

export interface NetworkInfo {
  chainId: string;
  networkName: string;
  isListening: boolean;
  currentBlock: number;
}

export interface BlockData {
  number: number;
  hash: string;
  timestamp: Date;
  transactions: number;
  gasUsed: number;
  gasLimit: number;
}

export interface TransactionData {
  hash: string;
  from: string;
  to: string;
  value: string;
  gas: number;
  gasPrice: string;
}

export interface GasPriceData {
  timestamp: Date;
  gasPrice: string;
  gasPriceGwei: string;
}

export interface DataCollectionStatus {
  isRunning: boolean;
  lastProcessedBlock: number;
  uptime: string;
}

export interface DataCollectionStats {
  blocksProcessed: number;
  transactionsProcessed: number;
  errors: number;
  lastUpdate: string;
}

export interface NavItem {
  path: string;
  label: string;
  icon: React.ReactElement;
}
