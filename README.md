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

The entrypoint for handling inference requests over HTTP/HTTPS. [Read more →]()

#### Model Enricher

_In Progress_ – A background service that can transform can trasform input data or provide supplemental data before processing an inference request. [Read more →]()

#### Model Executor

A background service (per model) that executes the acutal inferences request for the model and forwards thee result. [Read more →]()

#### Rendezvous Collector

A background service that receives model responses asyncronously as they are processed and returns the response for the active model back to the API. [Read more →]()
