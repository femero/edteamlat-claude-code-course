# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

A simple vanilla JavaScript todo list application built for educational purposes (Claude Code course).

## Architecture

See `sequence-diagram.md` for detailed interaction flows.

## Key Implementation Details

- **IDs**: Generated using `crypto.randomUUID()` (requires modern browser)
- **Language**: Spanish UI text (form labels, buttons)
- **No Persistence**: State resets on page reload (no localStorage or backend)

## Making Changes

When modifying functionality:
- All business logic is in `public/app.js`
- Keep the immutable update pattern (`map`/`filter` instead of direct mutations)
- Call `render()` after any state change
- Maintain the simple architecture - this is teaching code, not production
- en lo posible evita usar else en las clausulas if