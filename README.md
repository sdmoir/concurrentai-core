# concurrentai-core

## Getting Started

This repo contains the core components for Concurrent.ai. If you're looking to get started experimenting with your own Concurrent.ai stack, see [concurrentai/concurrentai-infra](https://github.com/concurrentai/concurrentai-infra) for a tutorial on how to deploy Concurrent.ai to a Kubernetes cluster or [Minikube](https://kubernetes.io/docs/setup/learning-environment/minikube/) environment.

## Overview

Concurrent.ai is a generalized manifestation of the Rendezvous architecture detailed by Ted Dunning and Ellen Friedman in [Machine Learning Logistics](https://www.oreilly.com/library/view/machine-learning-logistics/9781491997628/). Although there are many benefits that a Rendezvous architecture for machine learning offers, one major drawback is how high the initial engineering effort is to implement it. With Concurrent.ai, implementing a Rendezvous architecture is now as simple as writing a few lines of JSON.

### Core Components

![Core Components](https://github.com/concurrentai/concurrentai-core/raw/master/Concurrent.ai%20Core%20Components%20.png)
