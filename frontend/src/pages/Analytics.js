import React, { useState, useEffect } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  CircularProgress,
  Alert,
  Table,
  TableBody,
  TableCell,
  TableContainer,
  TableHead,
  TableRow,
  Paper,
  TextField,
  Button,
  FormControl,
  InputLabel,
  Select,
  MenuItem,
} from '@mui/material';
import {
  Search,
  FilterList,
  Download,
  TrendingUp,
  TrendingDown,
} from '@mui/icons-material';
import { LineChart, Line, XAxis, YAxis, CartesianGrid, Tooltip, ResponsiveContainer, BarChart, Bar } from 'recharts';
import axios from 'axios';

const Analytics = () => {
  const [blocks, setBlocks] = useState([]);
  const [transactions, setTransactions] = useState([]);
  const [gasPrices, setGasPrices] = useState([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [filter, setFilter] = useState('all');
  const [searchTerm, setSearchTerm] = useState('');

  useEffect(() => {
    fetchAnalyticsData();
  }, []);

  const fetchAnalyticsData = async () => {
    try {
      setLoading(true);
      const [blocksRes, transactionsRes, gasPricesRes] = await Promise.all([
        axios.get('/api/v1/analytics/blocks'),
        axios.get('/api/v1/analytics/transactions'),
        axios.get('/api/v1/analytics/gas-prices'),
      ]);
      
      setBlocks(blocksRes.data || []);
      setTransactions(transactionsRes.data || []);
      setGasPrices(gasPricesRes.data || []);
      setError(null);
    } catch (err) {
      setError('Failed to fetch analytics data');
      console.error('Analytics error:', err);
    } finally {
      setLoading(false);
    }
  };

  // Mock data for demonstration
  const mockBlocks = [
    { number: 12345678, hash: '0x1234...', timestamp: new Date(), transactions: 150, gasUsed: 15000000, gasLimit: 30000000 },
    { number: 12345677, hash: '0x5678...', timestamp: new Date(Date.now() - 12000), transactions: 145, gasUsed: 14800000, gasLimit: 30000000 },
    { number: 12345676, hash: '0x9abc...', timestamp: new Date(Date.now() - 24000), transactions: 160, gasUsed: 15200000, gasLimit: 30000000 },
  ];

  const mockTransactions = [
    { hash: '0xabcd...', from: '0x1234...', to: '0x5678...', value: '0.1', gas: 21000, gasPrice: '20000000000' },
    { hash: '0xefgh...', from: '0x9abc...', to: '0xdef0...', value: '0.05', gas: 65000, gasPrice: '25000000000' },
    { hash: '0xijkl...', from: '0x1111...', to: '0x2222...', value: '0.2', gas: 100000, gasPrice: '22000000000' },
  ];

  const mockGasPrices = [
    { timestamp: new Date(Date.now() - 60000), gasPrice: '20000000000', gasPriceGwei: '20' },
    { timestamp: new Date(Date.now() - 30000), gasPrice: '22000000000', gasPriceGwei: '22' },
    { timestamp: new Date(), gasPrice: '25000000000', gasPriceGwei: '25' },
  ];

  const filteredBlocks = filter === 'all' ? mockBlocks : mockBlocks.filter(block => 
    block.number.toString().includes(searchTerm) || 
    block.hash.includes(searchTerm)
  );

  const filteredTransactions = filter === 'all' ? mockTransactions : mockTransactions.filter(tx => 
    tx.hash.includes(searchTerm) || 
    tx.from.includes(searchTerm) || 
    tx.to.includes(searchTerm)
  );

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
      <Typography variant="h4" gutterBottom sx={{ mb: 3 }}>
        Analytics
      </Typography>

      {/* Search and Filter */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={6}>
          <TextField
            fullWidth
            label="Search"
            value={searchTerm}
            onChange={(e) => setSearchTerm(e.target.value)}
            placeholder="Search blocks, transactions, or addresses..."
          />
        </Grid>
        <Grid item xs={12} md={3}>
          <FormControl fullWidth>
            <InputLabel>Filter</InputLabel>
            <Select
              value={filter}
              label="Filter"
              onChange={(e) => setFilter(e.target.value)}
            >
              <MenuItem value="all">All</MenuItem>
              <MenuItem value="blocks">Blocks</MenuItem>
              <MenuItem value="transactions">Transactions</MenuItem>
            </Select>
          </FormControl>
        </Grid>
        <Grid item xs={12} md={3}>
          <Button
            fullWidth
            variant="outlined"
            startIcon={<Download />}
            onClick={() => alert('Export functionality coming soon!')}
          >
            Export Data
          </Button>
        </Grid>
      </Grid>

      {/* Charts */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Gas Price Trends (24h)
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <LineChart data={mockGasPrices}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="timestamp" tickFormatter={(value) => new Date(value).toLocaleTimeString()} />
                  <YAxis />
                  <Tooltip labelFormatter={(value) => new Date(value).toLocaleString()} />
                  <Line type="monotone" dataKey="gasPriceGwei" stroke="#1976d2" strokeWidth={2} />
                </LineChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Block Gas Usage
              </Typography>
              <ResponsiveContainer width="100%" height={300}>
                <BarChart data={mockBlocks}>
                  <CartesianGrid strokeDasharray="3 3" />
                  <XAxis dataKey="number" />
                  <YAxis />
                  <Tooltip />
                  <Bar dataKey="gasUsed" fill="#1976d2" />
                </BarChart>
              </ResponsiveContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Blocks Table */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Recent Blocks
              </Typography>
              <TableContainer component={Paper} sx={{ backgroundColor: 'transparent' }}>
                <Table>
                  <TableHead>
                    <TableRow>
                      <TableCell>Block Number</TableCell>
                      <TableCell>Hash</TableCell>
                      <TableCell>Timestamp</TableCell>
                      <TableCell>Transactions</TableCell>
                      <TableCell>Gas Used</TableCell>
                      <TableCell>Gas Limit</TableCell>
                      <TableCell>Utilization</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredBlocks.map((block, index) => {
                      const utilization = ((block.gasUsed / block.gasLimit) * 100).toFixed(1);
                      return (
                        <TableRow key={index}>
                          <TableCell>{block.number.toLocaleString()}</TableCell>
                          <TableCell>{block.hash}</TableCell>
                          <TableCell>{block.timestamp.toLocaleString()}</TableCell>
                          <TableCell>{block.transactions}</TableCell>
                          <TableCell>{block.gasUsed.toLocaleString()}</TableCell>
                          <TableCell>{block.gasLimit.toLocaleString()}</TableCell>
                          <TableCell>
                            <Box display="flex" alignItems="center" gap={1}>
                              {utilization}%
                              {utilization > 80 ? (
                                <TrendingUp color="warning" fontSize="small" />
                              ) : (
                                <TrendingDown color="success" fontSize="small" />
                              )}
                            </Box>
                          </TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Transactions Table */}
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
                      <TableCell>Gas Price (Gwei)</TableCell>
                      <TableCell>Total Cost (ETH)</TableCell>
                    </TableRow>
                  </TableHead>
                  <TableBody>
                    {filteredTransactions.map((tx, index) => {
                      const gasPriceEth = (parseInt(tx.gasPrice) / 1e18).toFixed(9);
                      const totalCost = ((parseInt(tx.gas) * parseInt(tx.gasPrice)) / 1e18).toFixed(6);
                      return (
                        <TableRow key={index}>
                          <TableCell>{tx.hash}</TableCell>
                          <TableCell>{tx.from}</TableCell>
                          <TableCell>{tx.to}</TableCell>
                          <TableCell>{tx.value}</TableCell>
                          <TableCell>{tx.gas.toLocaleString()}</TableCell>
                          <TableCell>{(parseInt(tx.gasPrice) / 1e9).toFixed(0)}</TableCell>
                          <TableCell>{totalCost}</TableCell>
                        </TableRow>
                      );
                    })}
                  </TableBody>
                </Table>
              </TableContainer>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default Analytics;
