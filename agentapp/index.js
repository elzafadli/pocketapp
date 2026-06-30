import 'dotenv/config';
process.env.OPENAI_API_KEY = process.env.DEEPSEEK_API_KEY;
import express from 'express';
import { app as graphApp } from './graph.js';

const port = process.env.PORT || 3000;
const app = express();

app.use(express.json());

app.post('/summarize', async (req, res) => {
  try {
    const inputData = req.body;
    
    // Check if data and url are provided
    if (!inputData || !inputData.data || !inputData.data.url) {
      return res.status(400).json({ error: "Invalid payload. 'data' and 'data.url' are required." });
    }

    const initialState = {
      url: inputData.data.url,
      data: inputData.data,
      scrapedContent: "",
      summary: "",
      todoList: [],
      error: ""
    };

    console.log(`Processing URL: ${initialState.url}`);
    const finalState = await graphApp.invoke(initialState);

    if (finalState.error) {
      console.error("Graph Error:", finalState.error);
      return res.status(500).json({ error: finalState.error });
    }

    res.json({
      success: true,
      summary: finalState.summary,
      todoList: finalState.todoList
    });

  } catch (error) {
    console.error("Error processing request:", error);
    res.status(500).json({ error: "Internal server error" });
  }
});

const server = app.listen(port, () => {
  console.log(`LangGraph Summarizer app listening on port ${port}`);
});

// Handle port already in use — jangan crash, tampilkan pesan jelas
server.on('error', (err) => {
  if (err.code === 'EADDRINUSE') {
    console.error(`\n[ERROR] Port ${port} sudah digunakan oleh proses lain.`);
    console.error(`Hentikan proses yang berjalan di port ${port} terlebih dahulu, lalu jalankan ulang.\n`);
    process.exit(1);
  } else {
    console.error('Server error:', err);
    process.exit(1);
  }
});

// Graceful shutdown — tangkap SIGINT (Ctrl+C) dan SIGTERM
function shutdown(signal) {
  console.log(`\n[${signal}] Shutting down gracefully...`);
  server.close(() => {
    console.log('Server closed. Bye!');
    process.exit(0);
  });

  // Force exit jika tidak selesai dalam 5 detik
  setTimeout(() => {
    console.error('Forced shutdown after timeout.');
    process.exit(1);
  }, 5000);
}

process.on('SIGINT', () => shutdown('SIGINT'));
process.on('SIGTERM', () => shutdown('SIGTERM'));
