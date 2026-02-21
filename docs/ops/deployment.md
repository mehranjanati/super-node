# Infrastructure & Deployment Guide

This guide describes the **Fractal DevOps** model: how to deploy and manage a Nexus Super Node using modern Infrastructure-as-Code (IaC).

## 1. Hardware Requirements

| Component | Minimum Spec | Recommended (Production) |
| :--- | :--- | :--- |
| **CPU** | 4 Cores (ARM64/AMD64) | 16+ Cores (EPYC/Threadripper) |
| **RAM** | 16 GB | 64 GB+ (ECC) |
| **Storage** | 500 GB NVMe | 4 TB NVMe RAID 0/1 |
| **Network** | 1 Gbps | 10 Gbps Fiber |
| **GPU** | Optional (for routing only) | NVIDIA A100/H100 or 2x 4090 |

## 2. The Tech Stack (Kubernetes/Nomad)

We support two deployment flavors:
1.  **K3s (Kubernetes)**: The industry standard. Best for scaling and tooling compatibility.
2.  **Nomad**: Simpler, single-binary. Best for edge devices and mixed workloads (Docker + Wasm binaries).

### Recommended: K3s Cluster
We use **FluxCD** for GitOps. You don't run `kubectl apply`; you push to Git.

```yaml
# cluster-config.yaml
apiVersion: k3s.cattle.io/v1
kind: HelmChart
metadata:
  name: nexus-stack
spec:
  chart: oci://ghcr.io/nexus/charts/super-node
  values:
    redpanda:
      enabled: true
      replicas: 3
    tidb:
      enabled: true
      storage: "nvme-local"
    unsloth:
      gpu:
        enabled: true
        count: 1
```

## 3. GPU Passthrough & AI Configuration

To enable **Unsloth** fine-tuning and LLM inference, the container runtime must access the GPU.

### NVIDIA Container Toolkit Setup
1.  Install drivers: `sudo apt install nvidia-driver-535`
2.  Install toolkit: `sudo apt install nvidia-container-toolkit`
3.  Configure Docker/Containerd to use the `nvidia` runtime.

### MIG (Multi-Instance GPU)
For A100/H100 cards, we enable MIG to slice one GPU into 7 smaller instances. This allows a single Super Node to serve 7 concurrent "Small Agent" workloads simultaneously.

## 4. Monitoring & Observability

A Super Node must be transparent. We use the **TIG Stack** (Telegraf, InfluxDB, Grafana) or Prometheus.

### Key Metrics to Watch
-   **Auction Win Rate**: How often are your bids accepted?
-   **Inference Latency**: Time to First Token (TTFT).
-   **Disk I/O**: Redpanda and TiDB are sensitive to disk latency.
-   **Peer Count**: Number of active DHT connections.

## 5. Security Hardening
-   **mTLS**: All internal traffic (Go <-> TiDB <-> Redpanda) is encrypted.
-   **Firewall**: Only port 443 (HTTPS) and 4001 (libp2p) should be exposed.
-   **Wasm Sandbox**: Ensure Wazero is configured to block file system access (except `/tmp`) and limit memory usage per instance.
