# AI-Powered Semantic Search Setup

This application uses **OpenAI embeddings** for powerful semantic search that understands meaning, not just keywords.

## Features

- **Semantic Understanding**: Search "hair of color pink" and find "pink hair", "rose-colored hair", etc.
- **Long Description Support**: Perfect for detailed descriptions - the AI understands context and relationships
- **Real-time Search**: Results update as you type (with 500ms debounce)
- **Ranked by Relevance**: Most similar items appear first using cosine similarity

## Setup Instructions

### 1. Get an OpenAI API Key

1. Go to [OpenAI Platform](https://platform.openai.com/)
2. Sign up or log in
3. Navigate to API Keys: https://platform.openai.com/api-keys
4. Click "Create new secret key"
5. Copy the key (it starts with `sk-`)

**Cost Estimation:**
- Using `text-embedding-3-small` model (best price/performance)
- ~$0.02 per 1 million tokens
- For 1,000 items with ~100 words each: ~$0.002 (very cheap!)
- Search queries are also very cheap (~$0.00002 per query)

### 2. Configure the API Key

Create a `.env` file in the project root:

```bash
cp .env.example .env
```

Edit `.env` and add your API key:

```
OPENAI_API_KEY=sk-your-actual-api-key-here
```

**Important:** Never commit your `.env` file to git! It's already in `.gitignore`.

### 3. Set Environment Variable

Before running the app, set the environment variable:

**Windows (PowerShell):**
```powershell
$env:OPENAI_API_KEY="sk-your-actual-api-key-here"
```

**Windows (Command Prompt):**
```cmd
set OPENAI_API_KEY=sk-your-actual-api-key-here
```

**Linux/Mac:**
```bash
export OPENAI_API_KEY="sk-your-actual-api-key-here"
```

Or add it to your system environment variables permanently.

### 4. Generate Embeddings (One-Time Setup)

Before search will work, you need to generate embeddings for all your models and stages.

**Option A: Use the UI Button (Easiest - Recommended)**

1. Start the app: `wails dev`
2. Look in the sidebar for the "Generate AI Search" button (with sparkles ✨ icon)
3. Click it and wait (progress shown in console)
4. Success message will appear when complete!

The button shows:
- Loading spinner while generating
- Real-time progress in the console
- Success/error status messages
- Safe to click multiple times (skips items that already have embeddings)

**Option B: Manual generation via code**

You can also trigger it programmatically:

```go
err := GenerateAllEmbeddings()
if err != nil {
    log.Fatal(err)
}
```

**What happens during generation:**
- Scans all models and stages
- Sends each description to OpenAI
- Stores the resulting embeddings in `data/data.json` and `data/stages.json`
- Skips items that already have embeddings (safe to run multiple times)
- Shows progress in console

**Time estimation:**
- ~100-200ms per item (API call + rate limiting)
- 1,000 items ≈ 2-3 minutes
- Only needs to run once (or when you add new items)

### 5. Run the Application

```bash
wails dev
```

The TypeScript bindings will be automatically generated when you run the dev server.

## Usage

### Searching

1. Navigate to Models or Stages page
2. Type your search query in the search bar
3. Results appear automatically (debounced by 500ms)
4. Results are ranked by semantic similarity

### Search Examples

**Traditional keyword search would fail, but AI search succeeds:**

- "character with pink colored hair" → finds "pink hair"
- "female protagonist" → finds descriptions containing "girl", "woman", "heroine"
- "outdoor scene with trees" → finds "forest", "woods", "nature"
- "sad emotional moment" → finds "crying", "depressed", "melancholic"

### When to Regenerate Embeddings

Run embedding generation again when:
- You add new models or stages
- You update descriptions significantly
- Embeddings are missing or corrupted

Just run `GenerateAllEmbeddings()` - it will skip items that already have embeddings.

## How It Works

### Architecture

```
User Query
    ↓
OpenAI Embeddings API → Vector (1536 dimensions)
    ↓
Compare with all stored item embeddings
    ↓
Calculate cosine similarity for each item
    ↓
Sort by similarity score (highest = most similar)
    ↓
Return top results
```

### Files Modified

**Backend (Go):**
- `embeddings.go` - OpenAI API integration, cosine similarity
- `models.go` - Added `Embedding []float64` field
- `stages.go` - Added `Embedding []float64` field
- `app.go` - `SearchModels()` and `SearchStages()` methods
- `generate_embeddings.go` - Batch embedding generation

**Frontend (React):**
- `ModelsGrid/index.tsx` - Search UI and logic
- `StagesGrid/index.tsx` - Search UI and logic

**Data Files:**
- `data/data.json` - Now includes `embedding` arrays
- `data/stages.json` - Now includes `embedding` arrays

## Troubleshooting

### Error: "OPENAI_API_KEY environment variable not set"

Solution: Make sure you set the environment variable before running the app.

### Error: "OpenAI API error (status 401)"

Solution: Your API key is invalid. Check that you copied it correctly.

### Error: "OpenAI API error (status 429)"

Solution: Rate limit exceeded. The code already has 100ms delays between requests. If you have many items, you might need to:
- Increase the delay in `generate_embeddings.go`
- Upgrade your OpenAI API tier
- Generate in smaller batches

### Search returns no results

Possible causes:
1. Embeddings haven't been generated yet - run `GenerateAllEmbeddings()`
2. API key not set correctly
3. Network/firewall blocking OpenAI API

Check the browser console and app console for error messages.

### TypeScript errors about SearchModels/SearchStages not existing

Solution: Run `wails dev` to regenerate TypeScript bindings.

## Cost Management

### Estimating Costs

**Initial embedding generation:**
- Formula: `(number_of_items × average_description_length_in_tokens) × $0.00002 / 1000`
- Example: 1,000 items × 100 tokens × $0.00002 / 1000 = $0.002

**Search queries:**
- Each search: ~20-50 tokens × $0.00002 / 1000 ≈ $0.000001
- 1,000 searches ≈ $0.001

**Total cost for typical usage:**
- Initial setup: < $0.01
- Monthly searches (1000 queries): < $0.01
- **Total: Less than a cup of coffee per year**

### Monitoring Usage

Check your OpenAI dashboard: https://platform.openai.com/usage

## Alternative: Offline Search

If you want to avoid OpenAI costs or work completely offline, you can switch to TF-IDF search (no ML model needed) - let me know and I can implement that instead.

## Advanced Configuration

### Change Embedding Model

Edit `embeddings.go`:

```go
const EmbeddingModel = "text-embedding-3-large" // Higher quality, more expensive
// or
const EmbeddingModel = "text-embedding-3-small" // Better price/performance (default)
```

### Adjust Search Result Limit

Edit `ModelsGrid/index.tsx` or `StagesGrid/index.tsx`:

```typescript
const results = await SearchModels(searchQuery, 100); // Change from 1000 to your limit
```

### Disable Debounce (Search Immediately)

Edit the `useEffect` delay in the grid components:

```typescript
const timer = setTimeout(() => {
    // ...
}, 0); // Change from 500 to 0 for instant search
```

## Support

If you encounter issues or need help, check:
1. Browser console (F12) for frontend errors
2. Terminal console for backend errors
3. OpenAI API status: https://status.openai.com/
