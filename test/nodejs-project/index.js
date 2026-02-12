#!/usr/bin/env node

console.log('Hello, Node.js Project!');

// 简单的 Express 服务器
const express = require('express');
const app = express();
const port = 3000;

app.get('/', (req, res) => {
  res.send('Hello, Node.js Project!');
});

app.listen(port, () => {
  console.log(`Server running at http://localhost:${port}`);
});
