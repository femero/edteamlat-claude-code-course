# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A simple vanilla JavaScript todo list application built for educational purposes (Claude Code course). The app uses no frameworks or build tools - just HTML, CSS, and JavaScript served via http-server.

## Running the Application

Start the development server:
```bash
npm start
```

This runs `npx http-server public -p 5173 -c-1`, serving the `public/` directory on port 5173 with caching disabled.

## Architecture

### State Management Pattern

The application uses a simple in-memory state pattern with immutable updates:

- **State**: Single `todos` array in `public/app.js:5`
- **Updates**: All state changes use immutable methods (`map()`, `filter()`) rather than mutations
- **Rendering**: Full re-render on every state change via `render()` function

### Data Flow

1. User interaction → Event handler
2. Event handler → State transformation (immutable)
3. State change → `render()` call
4. `render()` → Full DOM rebuild via `innerHTML = ''`

See `sequence-diagram.md` for detailed interaction flows.

### File Structure

```
public/
├── index.html  - Main HTML structure
├── app.js      - Application logic and state
└── styles.css  - Styling
```

All application code is client-side in the `public/` directory.

## Key Implementation Details

- **IDs**: Generated using `crypto.randomUUID()` (requires modern browser)
- **Todo Object**: `{ id: string, text: string, done: boolean }`
- **Re-rendering**: Complete list rebuild on each change (simple but not optimized for large lists)
- **Language**: Spanish UI text (form labels, buttons)
- **No Persistence**: State resets on page reload (no localStorage or backend)

## Making Changes

When modifying functionality:
- All business logic is in `public/app.js`
- Keep the immutable update pattern (`map`/`filter` instead of direct mutations)
- Call `render()` after any state change
- Maintain the simple architecture - this is teaching code, not production
- en lo posible evita usar else en las clausulas if