![Header](https://github.com/concurrentai/concurrentai-core/raw/master/misc/images/header.png)

Concurrent.ai is a platform that enables you to get started with machine learning in a rapid and evolvable way.

# concurrentai-core

- [Overview](https://github.com/concurrentai/concurrentai-core/tree/readme#getting-started)

## Overview

### Getting Started

This repo contains the core components for Concurrent.ai. If you're looking to start experimenting with your own Concurrent.ai stack, see **[concurrentai/concurrentai-infra](https://github.com/concurrentai/concurrentai-infra)**.

### Background

From a technical perspective, Concurrent.ai is a generalized manifestation of the Rendezvous architecture detailed by Ted Dunning and Ellen Friedman in [Machine Learning Logistics](https://www.oreilly.com/library/view/machine-learning-logistics/9781491997628/).

For a brief overview of the Rendezvous architecture and its many benefits, see [Rendezvous Architecture for Data Science in Production](https://towardsdatascience.com/rendezvous-architecture-for-data-science-in-production-79c4d48f12b) by Jan Teichmann – a highly recommended read!

### Why Concurrent.ai?

#### Reason #1

Although there are many benefits that the Rendezvous architecture offers, one major drawback is how high the initial engineering effort is to implement it. With Concurrent.ai, implementing a Rendezvous architecture can now be as simple as writing a few lines of JSON.

#### Reason #2

Concurrent.ai will extend the Rendezvous architecture concept beyond machine learning and into general business logic, allowing you to start with a simple, non-ML solution first and seamlessly iterate towards ML without rebuilding your application.

#### Reason #3

All of the benefits that come along with a Rendezvous architecture: auto-scaling, ability to validate model behavior and performance in production without impacting users, not being locked into a single ML framework, and more.

### Core Components

![Core Components](https://github.com/concurrentai/concurrentai-core/raw/master/misc/diagrams/Concurrent.ai%20Core%20Components.png)

#### Rendezvous API

The entrypoint for handling inference requests over HTTP/HTTPS. For reach request, the Rendezvous API assigns a unique request ID, publishes the request to a "model-request" topic, then waits for a model response over a unix socket that will be written by the Rendezvous Collector.

When deployed within a Concurrent.ai stack, the Rendezvous API sits behind an API gateway that handles timouts, SSL, etc.

#### Model Enricher

(_In Progress_) A background service that subscribes to the "model-request" topic, applies any configured transformations or data integrations, then forwards the request to the "model-input" topic. Currently, support for transformations and data integrations has not yet been implemented, so the Model Enricher is acting only as a pass-through service.

#### Model Executor

A background service (per model) that subscribes to the "model-input" topic, executes the inference request, then publishes to the "model-response" topic. Currently, the Model Executor only supports MLflow models built as Docker images, but additional framework support is a top priority.

#### Rendezvous Collector

A background service that subscribes to the "model-response" topic, checks if a response is for an active model, and if so writes the response back to the API over a unix socket based on the request ID.
