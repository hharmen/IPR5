const express = require('express');
const axios = require('axios');
const path = require('path');

const app = express();
const PORT = process.env.PORT || 3000;
const BACKEND_URL = process.env.BACKEND_URL || 'http://backend-service:5000';

app.use(express.static(path.join(__dirname, 'public')));
app.use(express.json());

app.get('/api/tasks', async (req, res) => {
  try {
    const response = await axios.get(`${BACKEND_URL}/api/tasks`);
    res.json(response.data);
  } catch (error) {
    console.error('Error fetching tasks:', error.message);
    res.status(500).json({ error: 'Failed to fetch tasks' });
  }
});

app.post('/api/tasks', async (req, res) => {
  try {
    const response = await axios.post(`${BACKEND_URL}/api/tasks`, req.body);
    res.status(201).json(response.data);
  } catch (error) {
    console.error('Error creating task:', error.message);
    res.status(500).json({ error: 'Failed to create task' });
  }
});

app.put('/api/tasks/:id', async (req, res) => {
  try {
    const response = await axios.put(`${BACKEND_URL}/api/tasks/${req.params.id}`);
    res.json(response.data);
  } catch (error) {
    console.error('Error updating task:', error.message);
    res.status(500).json({ error: 'Failed to update task' });
  }
});

app.delete('/api/tasks/:id', async (req, res) => {
  try {
    await axios.delete(`${BACKEND_URL}/api/tasks/${req.params.id}`);
    res.sendStatus(204);
  } catch (error) {
    console.error('Error deleting task:', error.message);
    res.status(500).json({ error: 'Failed to delete task' });
  }
});

app.get('/health', (req, res) => {
  res.json({ status: 'healthy', service: 'task-frontend' });
});

app.listen(PORT, () => {
  console.log(`Frontend started on port ${PORT}, backend: ${BACKEND_URL}`);
});
