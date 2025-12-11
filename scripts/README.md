# Scripts

This directory contains utility scripts for managing the gosources application.

## populate_sources.go

This script populates the database with source configurations including all CSS selectors.

### What it does

1. **Updates the sudbury.com source** with selectors from sources.yml:
   - Article selectors for title, body, published_time, image, link, category
   - Exclude patterns for cleaning content

2. **Creates the Mid-North Monitor source** with comprehensive selectors:
   - Article selectors (24 fields)
   - List selectors (container, article cards, article list)
   - Page selectors (container, title, content, metadata)
   - Extensive exclude patterns for ads, navigation, scripts, etc.

### Prerequisites

- The API server must be running on `http://localhost:8050`
- The database must be initialized and migrated

### Usage

1. Start your services (API server + PostgreSQL):
   ```bash
   # Option 1: Using Docker Compose
   task docker:up

   # Option 2: Start services manually
   # Start PostgreSQL, then:
   go run main.go -config config.yml
   ```

2. Run the population script:
   ```bash
   cd scripts
   go run populate_sources.go
   ```

3. Verify the sources:
   - Visit http://localhost:3000/sources
   - Edit sudbury.com: http://localhost:3000/sources/96abf22e-2577-4b4b-aea3-c524e5c8dd3a/edit
   - You should see the new Mid-North Monitor source in the list

### Expected Output

```
✓ Updated source: sudbury.com
✓ Created source: Mid-North Monitor

✓ All sources populated successfully!
```

### Troubleshooting

**Error: "connection refused"**
- Make sure the API server is running on port 8050
- Check with: `curl http://localhost:8050/api/v1/sources`

**Error: "update failed with status 404"**
- The sudbury.com source may have a different ID
- Check the current ID in the database or via the API

**Error: "create failed with status 409"**
- Mid-North Monitor source already exists
- Either delete it first or modify the script to update instead

### Manual Alternative

If you prefer to use the UI:

1. Visit http://localhost:3000/sources/96abf22e-2577-4b4b-aea3-c524e5c8dd3a/edit
2. Expand the "Article Selectors" section
3. Fill in the values from the script above
4. Click "Update Source"

For Mid-North Monitor:
1. Visit http://localhost:3000/sources/new
2. Fill in all fields from the script
3. Click "Create Source"
