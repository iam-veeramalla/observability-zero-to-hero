// service-b/index.js
require('dotenv').config();
require('./tracing'); // Add this line to initialize tracing
const express = require('express');
const morgan = require('morgan');

const app = express();
const PORT = 3002;
app.use(morgan('common'))

app.get('/hello', (req, res) => {
  res.send('Hello from Service B!');
});

app.listen(PORT, () => {
  console.log(`Service B is running on port ${PORT}`);
});