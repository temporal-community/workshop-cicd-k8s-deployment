"""Shared types that mirror the Go implementation."""

from dataclasses import dataclass
from typing import Optional


@dataclass
class DeployToKubernetesRequest:
    ImageTag: str
    Environment: str  # staging or production


@dataclass
class DeployToKubernetesResponse:
    success: bool
    deployment_url: str
    message: str
    timestamp: str


@dataclass
class CheckDeploymentStatusRequest:
    Environment: str


@dataclass
class CheckDeploymentStatusResponse:
    ready: bool
    replicas: int
    ready_replicas: int
    message: str


@dataclass
class RollbackDeploymentRequest:
    Environment: str
    ImageTag: Optional[str] = None
    Reason: str = ""
    Timestamp: Optional[str] = None


@dataclass
class RollbackDeploymentResponse:
    success: bool
    message: str
    timestamp: Optional[str] = None


@dataclass
class GetServiceURLRequest:
    Environment: str
    ServiceName: str


@dataclass
class GetServiceURLResponse:
    url: str
    ready: bool
    message: str