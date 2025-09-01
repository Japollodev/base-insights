import React, { useState, useEffect } from 'react';
import {
  Grid,
  Card,
  CardContent,
  Typography,
  Box,
  CircularProgress,
  Alert,
  Button,
  Switch,
  FormControlLabel,
  LinearProgress,
  Chip,
} from '@mui/material';
import {
  PlayArrow,
  Stop,
  Refresh,
  CheckCircle,
  Error,
  Warning,
} from '@mui/icons-material';
import axios from 'axios';

const DataCollection = () => {
  const [status, setStatus] = useState(null);
  const [stats, setStats] = useState(null);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState(null);
  const [isStarting, setIsStarting] = useState(false);
  const [isStopping, setIsStopping] = useState(false);

  useEffect(() => {
    fetchDataCollectionData();
    const interval = setInterval(fetchDataCollectionData, 5000); // Refresh every 5 seconds
    return () => clearInterval(interval);
  }, []);

  const fetchDataCollectionData = async () => {
    try {
      const [statusRes, statsRes] = await Promise.all([
        axios.get('/api/v1/data/status'),
        axios.get('/api/v1/data/stats'),
      ]);
      
      setStatus(statusRes.data);
      setStats(statsRes.data);
      setError(null);
    } catch (err) {
      setError('Failed to fetch data collection data');
      console.error('Data collection error:', err);
    } finally {
      setLoading(false);
    }
  };

  const handleStartService = async () => {
    try {
      setIsStarting(true);
      await axios.post('/api/v1/data/start');
      await fetchDataCollectionData();
    } catch (err) {
      setError('Failed to start data collection service');
      console.error('Start service error:', err);
    } finally {
      setIsStarting(false);
    }
  };

  const handleStopService = async () => {
    try {
      setIsStopping(true);
      await axios.post('/api/v1/data/stop');
      await fetchDataCollectionData();
    } catch (err) {
      setError('Failed to stop data collection service');
      console.error('Stop service error:', err);
    } finally {
      setIsStopping(false);
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

  const getStatusColor = (isRunning) => {
    return isRunning ? 'success' : 'error';
  };

  const getStatusIcon = (isRunning) => {
    return isRunning ? <CheckCircle color="success" /> : <Error color="error" />;
  };

  return (
    <Box>
      <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 3 }}>
        <Typography variant="h4">
          Data Collection
        </Typography>
        <Button
          variant="outlined"
          startIcon={<Refresh />}
          onClick={fetchDataCollectionData}
        >
          Refresh
        </Button>
      </Box>

      {/* Service Control */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Box display="flex" justifyContent="space-between" alignItems="center" sx={{ mb: 2 }}>
                <Typography variant="h6">
                  Service Control
                </Typography>
                <Chip
                  icon={getStatusIcon(status?.isRunning)}
                  label={status?.isRunning ? 'Running' : 'Stopped'}
                  color={getStatusColor(status?.isRunning)}
                  variant="outlined"
                />
              </Box>
              
              <Box display="flex" gap={2}>
                <Button
                  variant="contained"
                  color="success"
                  startIcon={<PlayArrow />}
                  onClick={handleStartService}
                  disabled={status?.isRunning || isStarting}
                  sx={{ minWidth: 120 }}
                >
                  {isStarting ? 'Starting...' : 'Start'}
                </Button>
                
                <Button
                  variant="contained"
                  color="error"
                  startIcon={<Stop />}
                  onClick={handleStopService}
                  disabled={!status?.isRunning || isStopping}
                  sx={{ minWidth: 120 }}
                >
                  {isStopping ? 'Stopping...' : 'Stop'}
                </Button>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Service Status */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Service Status
              </Typography>
              <Box>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                  Status: {status?.isRunning ? 'Running' : 'Stopped'}
                </Typography>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                  Last Processed Block: {status?.lastProcessedBlock?.toLocaleString() || 0}
                </Typography>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                  Uptime: {status?.uptime || 'N/A'}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>

        <Grid item xs={12} md={6}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Processing Progress
              </Typography>
              <Box>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                  Blocks Processed: {stats?.blocksProcessed?.toLocaleString() || 0}
                </Typography>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                  Transactions Processed: {stats?.transactionsProcessed?.toLocaleString() || 0}
                </Typography>
                <Typography variant="body2" color="textSecondary" sx={{ mb: 1 }}>
                  Errors: {stats?.errors || 0}
                </Typography>
                <Typography variant="body2" color="textSecondary">
                  Last Update: {stats?.lastUpdate ? new Date(stats.lastUpdate).toLocaleString() : 'Never'}
                </Typography>
              </Box>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Performance Metrics */}
      <Grid container spacing={3} sx={{ mb: 4 }}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Performance Metrics
              </Typography>
              
              <Grid container spacing={3}>
                <Grid item xs={12} md={4}>
                  <Box>
                    <Typography variant="body2" color="textSecondary" gutterBottom>
                      Processing Rate
                    </Typography>
                    <Typography variant="h4" color="primary">
                      {stats?.blocksProcessed > 0 ? Math.round(stats.blocksProcessed / Math.max(1, Math.floor((Date.now() - new Date(stats.lastUpdate).getTime()) / 1000 / 60))) : 0}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      blocks/minute
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={12} md={4}>
                  <Box>
                    <Typography variant="body2" color="textSecondary" gutterBottom>
                      Error Rate
                    </Typography>
                    <Typography variant="h4" color="error">
                      {stats?.errors > 0 && stats?.blocksProcessed > 0 ? ((stats.errors / stats.blocksProcessed) * 100).toFixed(2) : 0}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      %
                    </Typography>
                  </Box>
                </Grid>
                
                <Grid item xs={12} md={4}>
                  <Box>
                    <Typography variant="body2" color="textSecondary" gutterBottom>
                      Efficiency
                    </Typography>
                    <Typography variant="h4" color="success">
                      {stats?.blocksProcessed > 0 ? Math.round((stats.blocksProcessed / (stats.blocksProcessed + stats.errors)) * 100) : 100}
                    </Typography>
                    <Typography variant="body2" color="textSecondary">
                      %
                    </Typography>
                  </Box>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>
      </Grid>

      {/* Configuration Info */}
      <Grid container spacing={3}>
        <Grid item xs={12}>
          <Card>
            <CardContent>
              <Typography variant="h6" gutterBottom>
                Configuration
              </Typography>
              <Grid container spacing={2}>
                <Grid item xs={12} sm={6} md={3}>
                  <Typography variant="body2" color="textSecondary">
                    Block Scan Interval: 5s
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                  <Typography variant="body2" color="textSecondary">
                    Transaction Scan: 3s
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                  <Typography variant="body2" color="textSecondary">
                    Price Updates: 1m
                  </Typography>
                </Grid>
                <Grid item xs={12} sm={6} md={3}>
                  <Typography variant="body2" color="textSecondary">
                    Batch Size: 10
                  </Typography>
                </Grid>
              </Grid>
            </CardContent>
          </Card>
        </Grid>
      </Grid>
    </Box>
  );
};

export default DataCollection;
