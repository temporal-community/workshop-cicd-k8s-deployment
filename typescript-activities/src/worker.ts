import { Worker } from '@temporalio/worker';
import { sendApprovalRequest, logApprovalDecision, sendApprovalNotification } from './activities/approval';

/**
 * TypeScript worker for approval activities
 * Connects to the same Temporal cluster as the Go workers
 */
async function run() {
  // Get Temporal host from environment variable or use default
  const temporalHost = process.env.TEMPORAL_HOST || 'localhost:7233';
  
  console.log('Starting TypeScript worker for approval activities');
  console.log(`Connecting to Temporal server: ${temporalHost}`);
  console.log('Worker listening on task queue: cicd-task-queue-typescript');
  console.log('Registered activities:');
  console.log('  - TypeScript Approval: SendApprovalRequest, LogApprovalDecision, SendApprovalNotification');

  // Create and configure the worker
  const worker = await Worker.create({
    // Use TypeScript-specific task queue
    taskQueue: 'cicd-task-queue-typescript',
    
    // Register approval activities with PascalCase names to match Go workflow calls
    activities: {
      SendApprovalRequest: sendApprovalRequest,
      LogApprovalDecision: logApprovalDecision,
      SendApprovalNotification: sendApprovalNotification,
    },
  });

  // Handle graceful shutdown
  const shutdownHandler = () => {
    console.log('Received shutdown signal, shutting down worker...');
    void worker.shutdown();
  };

  process.on('SIGINT', shutdownHandler);
  process.on('SIGTERM', shutdownHandler);

  // Start the worker
  console.log('Worker started successfully!');
  console.log('Press Ctrl+C to stop the worker');
  
  try {
    await worker.run();
  } catch (error) {
    console.error('Worker error:', error);
    process.exit(1);
  }
}

// Create empty workflows file to satisfy worker requirements
// We only register activities, not workflows
const emptyWorkflows = {};

// Start the worker
run().catch((error) => {
  console.error('Failed to start worker:', error);
  process.exit(1);
});