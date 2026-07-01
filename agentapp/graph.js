import { StateGraph, START, END, Annotation } from "@langchain/langgraph";
import { ChatOpenAI } from "@langchain/openai";
import { HumanMessage, SystemMessage } from "@langchain/core/messages";
import * as cheerio from "cheerio";

// Define the state
export const GraphState = Annotation.Root({
  url: Annotation(),
  data: Annotation(), // The input JSON data
  scrapedContent: Annotation(),
  summary: Annotation(),
  todoList: Annotation(),
  error: Annotation(),
});

// Node 1: Scrape URL
async function scrapeUrl(state) {
  try {
    const { url } = state;
    if (!url) {
      throw new Error("URL is missing");
    }
    
    // Fetch HTML
    const response = await fetch(url);
    if (!response.ok) {
      throw new Error(`Failed to fetch URL: ${response.statusText}`);
    }
    const html = await response.text();
    
    // Parse with cheerio to extract readable text
    const $ = cheerio.load(html);
    // Remove scripts and styles
    $('script, style, noscript, iframe, img, svg').remove();
    
    const textContent = $('body').text().replace(/\s+/g, ' ').trim();
    
    // Return scraped content (limit to first 10000 chars to avoid token limits)
    return { scrapedContent: textContent.substring(0, 10000) };
  } catch (error) {
    return { error: error.message };
  }
}

// Node 2: Summarize and Extract Todo
async function processContent(state) {
  if (state.error) {
    return state;
  }

  try {
    const llm = new ChatOpenAI({
      apiKey: process.env.DEEPSEEK_API_KEY,
      configuration: {
        baseURL: "https://api.deepseek.com",
      },
      modelName: "deepseek-v4-flash",
      temperature: 0.3,
    });

    const systemPrompt = `You are a helpful assistant. You have been given the contents of a webpage and some metadata about it.
Your task is to:
1. Provide a concise summary of the webpage in Indonesian (Bahasa Indonesia).
2. Extract a structured TO-DO list in Indonesian (Bahasa Indonesia) based on the content (e.g. actions the user should take after reading).

Respond ONLY with a JSON object in the following format:
{
  "summary": "String containing the summary",
  "todoList": ["Task 1", "Task 2"]
}`;

    const humanPrompt = `Metadata:
${JSON.stringify(state.data, null, 2)}

Webpage Content:
${state.scrapedContent}`;

    const messages = [
      new SystemMessage(systemPrompt),
      new HumanMessage(humanPrompt),
    ];

    const response = await llm.invoke(messages);
    let result;
    try {
      // Clean up markdown block if any
      const rawText = response.content.replace(/```json/g, "").replace(/```/g, "").trim();
      result = JSON.parse(rawText);
    } catch (e) {
      throw new Error("Failed to parse LLM output as JSON: " + response.content);
    }

    return {
      summary: result.summary,
      todoList: result.todoList,
    };
  } catch (error) {
    return { error: error.message };
  }
}

// Build the graph
const workflow = new StateGraph(GraphState)
  .addNode("scrape", scrapeUrl)
  .addNode("process", processContent)
  .addEdge(START, "scrape")
  .addEdge("scrape", "process")
  .addEdge("process", END);

export const app = workflow.compile();
