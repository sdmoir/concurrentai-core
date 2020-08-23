![Header](https://github.com/concurrentai/concurrentai-core/raw/master/misc/images/header.png)

Concurrent.ai is a platform that enables you to get started with machine learning in a rapid and evolvable way.

# concurrentai-core

- [Overview](https://github.com/concurrentai/concurrentai-core/tree/readme#getting-started)
- [Development](https://github.com/concurrentai/concurrentai-core/tree/readme#development)
- [Roadmap](https://github.com/concurrentai/concurrentai-core/tree/readme#roadmap)

## Getting Started

See **[concurrentai/concurrentai-infra](https://github.com/concurrentai/concurrentai-infra)** to start experimenting with your own Concurrent.ai stack.

## Concurrent.ai Overview

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

![Core Components](misc/diagrams/Concurrent.ai%20Core%20Components.png)

#### Rendezvous API

The entrypoint for handling inference requests over HTTP/HTTPS. [Read more →](https://github.com/concurrentai/concurrentai-core/tree/readme/src/rendezvous-api)

#### Model Enricher

_In Progress_. Optionally transforms input data and/or provides supplemental data before processing an inference request. [Read more →](https://github.com/concurrentai/concurrentai-core/tree/readme/src/model-enricher)

#### Model Executor(s)

Executes the acutal inference request for each model and forwards the response. [Read more →](https://github.com/concurrentai/concurrentai-core/tree/readme/src/model-executor)

#### Rendezvous Collector

Receives model responses asyncronously as they are processed and returns the "active" model's response back to the API. [Read more →](https://github.com/concurrentai/concurrentai-core/tree/readme/src/rendezvous-collector)

#### Analysis Collector

_Not yet implemented_. Receives all model responses asyncronously and forwards them to an external system such as Elasticsearch for further analysis.

## Roadmap

See the [Concurrent.ai Roadmap](https://github.com/orgs/concurrentai/projects/1) project for an up-to-date roadmap.

## Contributing

Pull requests are welcome! Many details here are still being worked out – see [CONTRIBUTING.md](CONTRIBUTING.md) to get started.

Everyone contributing to Concurrent.ai repositories or engaging in discussion is expected to follow the [Code of Conduct](CODE_OF_CONDUCT.md).

## License

Licensed under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0).