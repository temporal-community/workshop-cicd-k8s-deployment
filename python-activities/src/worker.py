"""
Python worker for Kubernetes activities.
Connects to the same Temporal cluster as the Go workers.
"""

import asyncio
import logging
import os
import signal
from typing import Optional

from temporalio.client import Client
from temporalio.worker import Worker

from src.activities.kubernetes import KubernetesActivities

# Configure logging
logging.basicConfig(
    level=logging.INFO,
    format="%(asctime)s - %(name)s - %(levelname)s - %(message)s"
)
logger = logging.getLogger(__name__)


async def main() -> None:
    """Main function to start the Python worker."""
    # Get Temporal host from environment variable or use default
    temporal_host = os.getenv("TEMPORAL_HOST", "localhost:7233")
    
    logger.info("Starting Python worker for Kubernetes activities")
    logger.info(f"Connecting to Temporal server: {temporal_host}")
    logger.info("Worker listening on task queue: cicd-task-queue-python")
    logger.info("Registered activities:")
    logger.info("  - Python Kubernetes: DeployToKubernetes, CheckDeploymentStatus, RollbackDeployment, GetServiceURL")

    # Create Temporal client
    client = await Client.connect(temporal_host)
    
    # Create Kubernetes activities instance
    kubernetes_activities = KubernetesActivities()
    
    # Create and configure the worker
    worker = Worker(
        client,
        task_queue="cicd-task-queue-python",  # Use Python-specific task queue
        activities=[
            kubernetes_activities.deploy_to_kubernetes,
            kubernetes_activities.check_deployment_status,
            kubernetes_activities.rollback_deployment,
            kubernetes_activities.get_service_url,
        ],
    )

    await worker.run()





if __name__ == "__main__":
    asyncio.run(main())