# TypeScript Approval Activities

This directory contains the TypeScript implementation of approval activities for the Temporal CI/CD workshop. It demonstrates how to implement the same functionality as the Go approval activities using TypeScript and the Temporal TypeScript SDK.

## Overview

The TypeScript worker implements three approval activities:
- `sendApprovalRequest`: Sends a notification for approval with workflow details
- `logApprovalDecision`: Logs approval/rejection decisions with timestamps  
- `sendApprovalNotification`: Sends formatted notification messages

## Setup

### Prerequisites
- Node.js 18.0.0 or higher
- npm or yarn package manager
- Temporal server running (default: localhost:7233)

### Installation

1. Install dependencies:
```bash
npm install
```

2. Build the TypeScript code:
```bash
npm run build
```

## Running the Worker

### Development Mode (with TypeScript compilation)
```bash
npm run dev
```

### Production Mode (compiled JavaScript)
```bash
npm run build
npm start
```

### Environment Variables
- `TEMPORAL_HOST`: Temporal server address (default: localhost:7233)

## Project Structure

```
typescript-activities/
├── package.json          # NPM dependencies and scripts
├── tsconfig.json         # TypeScript configuration
├── src/
│   ├── activities/
│   │   └── approval.ts   # Approval activity implementations
│   ├── types.ts          # TypeScript type definitions
│   ├── worker.ts         # Worker startup and configuration
│   └── workflows.ts      # Empty workflows file (required by SDK)
└── README.md
```

## Activities

### sendApprovalRequest
Logs an approval request with workflow information and provides CLI commands for approval/rejection.

### logApprovalDecision  
Records approval decisions with timestamps and approver information.

### sendApprovalNotification
Sends formatted notifications about approval outcomes.

## Integration with Go Workflow

This TypeScript worker connects to the same Temporal cluster and task queue (`cicd-task-queue`) as the Go workers. The Go workflow can call these TypeScript activities seamlessly, demonstrating Temporal's polyglot capabilities.

## Development

### Linting
```bash
npm run lint
```

### Cleaning Build Artifacts
```bash
npm run clean
```

## Notes

- Activities use structured logging compatible with Temporal's observability features
- Error handling follows TypeScript best practices with proper exception propagation
- The worker handles graceful shutdown on SIGINT/SIGTERM signals
- All activities are designed to be idempotent and retry-safe