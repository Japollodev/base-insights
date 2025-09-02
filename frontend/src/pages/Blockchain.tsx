import React, { useState, useEffect, ChangeEvent } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  CircularProgress,
  Alert,
  Button,
  TextField,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
} from '@mui/material';
import { Search, Refresh } from '@mui/icons-material';
import axios from 'axios';
import { Block, Transaction, GasPrice, NetworkInfo } from '../types';

function Blockchain() {
  const [currentBlock, setCurrentBlock] = useState<Block | null>(null);
  const [gasPrice, setGasPrice] = useState<GasPrice | null>(null);
  const [networkInfo, setNetworkInfo] = useState<NetworkInfo | null>(null);
  const [loading, setLoading] = useState<boolean>(true);
  const [error, setError] = useState<string | null>(null);
  const [searchBlock, setSearchBlock] = useState<string>('');
  const [searchedBlock, setSearchedBlock] = useState<Block | null>(null);

  useEffect(() => {
    fetchBlockchainData();
    const interval = setInterval(fetchBlockchainData, 10000); // Refresh every 10 seconds
    return () => clearInterval(interval);
  }, []);

  const fetchBlockchainData = async (): Promise<void> => {
    try {
      setLoading(true);
      const [blockRes, gasRes, networkRes] = await Promise.all([
        axios.get<Block>('/api/v1/blockchain/current-block'),
        axios.get<GasPrice>('/api/v1/blockchain/gas-price'),
        axios.get<NetworkInfo>('/api/v1/blockchain/network-info'),
      ]);
      
      setCurrentBlock(blockRes.data);
      setGasPrice(gasRes.data);
      setNetworkInfo(networkRes.data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch blockchain data');
      console.error('Blockchain error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleSearchBlock = async (): Promise<void> => {
    if (!searchBlock) return;
    
    try {
      const response = await axios.get<Block>(`/api/v1/blockchain/block/${searchBlock}`);
      setSearchedBlock(response.data);
    } catch (err) {
      setError('Failed to fetch block');
      console.error('Search block error:', err);
    }
  };

  if (loading) {
    return (
      <Box display="flex" justifyContent="center" alignItems="center" minHeight="400px">
        <CircularProgress />
      </Box>
    );
  }

  if (error) {
    return <Alert severity="error">{error}</Alert>;
  }

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 3 }}>
        <Typography variant="h4">
          Blockchain
        </Typography>
        <Button
          variant="outlined"
          startIcon={<Refresh />}
          onClick={fetchBlockchainData}
        >
          Refresh
        </Button>
      </Box>

      {/* Current Block Info */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Current Block
              </Typography>
              {currentBlock && (
                <Box>
                  <Typography variant="body2" color="textSecondary">
                    Block Number: {currentBlock.number}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Hash: {currentBlock.hash?.substring(0, 20)}...
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Timestamp: {new Date(currentBlock.timestamp).toLocaleString()}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Transactions: {currentBlock.transactions?.length || 0}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Gas Used: {currentBlock.gasUsed?.toLocaleString()}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Gas Limit: {currentBlock.gasLimit?.toLocaleString()}
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Gas Price
              </Typography>
              {gasPrice && (
                <Box>
                  <Typography variant="body2" color="textSecondary">
                    Gas Price: {gasPrice.gasPrice} Wei
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Gas Price: {gasPrice.gasPriceGwei} Gwei
                  </Typography>
                  {gasPrice.maxFeePerGas && (
                    <Typography variant="body2" color="textSecondary">
                      Max Fee Per Gas: {gasPrice.maxFeePerGas} Wei
                    </Typography>
                  )}
                  {gasPrice.maxPriorityFeePerGas && (
                    <Typography variant="body2" color="textSecondary">
                      Max Priority Fee: {gasPrice.maxPriorityFeePerGas} Wei
                    </Typography>
                  )}
                  {gasPrice.baseFeePerGas && (
                    <Typography variant="body2" color="textSecondary">
                      Base Fee Per Gas: {gasPrice.baseFeePerGas} Wei
                    </Typography>
                  )}
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Network Info */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Network Information
              </Typography>
              {networkInfo && (
                <Grid container spacing={2}>
                  <Grid item xs={12} sm={6} md={3}>
                    <Typography variant="body2" color="textSecondary">
                      Chain ID: {networkInfo.chainId}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={6} md={3}>
                    <Typography variant="body2" color="textSecondary">
                      Network: {networkInfo.networkName}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={6} md={3}>
                    <Typography variant="body2" color="textSecondary">
                      Status: {networkInfo.isListening ? 'Connected' : 'Disconnected'}
                    </Typography>
                  </Grid>
                  <Grid item xs={12} sm={6} md={3}>
                    <Typography variant="body2" color="textSecondary">
                      Current Block: {networkInfo.currentBlock}
                    </Typography>
                  </Grid>
                </Grid>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Block Search */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Search Block
              </Typography>
              <Box display="flex" gap={2} alignItems="center">
                <TextField
                  label="Block Number"
                  type="number"
                  value={searchBlock}
                  onChange={(e: ChangeEvent<HTMLInputElement | HTMLTextAreaElement>) => setSearchBlock(e.target.value)}
                  placeholder="Enter block number"
                  sx={{ flexGrow: 1 }}
                />
                <Button
                  variant="contained"
                  startIcon={<Search />}
                  onClick={handleSearchBlock}
                  disabled={!searchBlock}
                >
                  Search
                </Button>
              </Box>
              
              {searchedBlock && (
                <Box sx={{ mt: 2 }}>
                  <Typography variant="subtitle1" gutterBottom>
                    Block {searchedBlock.number}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Hash: {searchedBlock.hash?.substring(0, 20)}...
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Timestamp: {new Date(searchedBlock.timestamp).toLocaleString()}
                  </Typography>
                  <Typography variant="body2" color="textSecondary">
                    Transactions: {searchedBlock.transactions?.length || 0}
                  </Typography>
                </Box>
              )}
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Recent Transactions */}
      {currentBlock?.transactions && currentBlock.transactions.length > 0 && (
        <Grid container spacing={3}>
          <Grid item xs={12}>
            <Card>
              <CardContent>
                <Typography variant="h6" gutterBottom>
                  Recent Transactions
                </Typography>
                <TableContainer component={Paper} sx={{ backgroundColor: 'transparent' }}>
                  <Table>
                    <TableHead>
                      <TableRow>
                        <TableCell>Hash</TableCell>
                        <TableCell>From</TableCell>
                        <TableCell>To</TableCell>
                        <TableCell>Value (ETH)</TableCell>
                        <TableCell>Gas</TableCell>
                      </TableRow>
                    </TableHead>
                    <TableBody>
                      {currentBlock.transactions.slice(0, 10).map((tx: Transaction, index: number) => (
                        <TableRow key={index}>
                          <TableCell>{tx.hash?.substring(0, 20)}...</TableCell>
                          <TableCell>{tx.from?.substring(0, 20)}...</TableCell>
                          <TableCell>{tx.to?.substring(0, 20)}...</TableCell>
                          <TableCell>{tx.value ? (parseInt(tx.value) / 1e18).toFixed(6) : '0'}</TableCell>
                          <TableCell>{tx.gas?.toLocaleString()}</TableCell>
                        </TableRow>
                      ))}
                    </TableBody>
                  </Table>
                </TableContainer>
              </CardContent>
            </Card>
          </Grid>
        </Grid>
      )}
    </Box>
  );
};

export default Blockchain;
