![Header](https://github.com/concurrentai/concurrentai-core/raw/master/misc/images/header.png)

Concurrent.ai is a platform that enables you to get started with machine learning in a rapid and evolvable way.

# concurrentai-core

This repo contains the core components for Concurrent.ai.

If you're looking to start experimenting with your own Concurrent.ai stack, see **[concurrentai/concurrentai-infra](https://github.com/concurrentai/concurrentai-infra)**.

## Background

From a technical perspective, Concurrent.ai is a generalized manifestation of the Rendezvous architecture detailed by Ted Dunning and Ellen Friedman in [Machine Learning Logistics](https://www.oreilly.com/library/view/machine-learning-logistics/9781491997628/).

### What is the Rendezvous architecture?

Rendezvous architecture is essentially a framework for running multiple versions of a machine learning model behind a single API or stream endpoint. While only one model is considered "active" and used for responses, all other models can be evaluated in parallel without affecting users.

#### Benefits of a Rendezvous architecture

There are many, including:
- Faster iteration time by evaluating multiple models in parallel
- Increased confidence by validating model performance and behavior before impacting users
- Long-term flexibility by decoupling API interfaces from ML frameworks

#### Further reading

There is a phenomenal article written by Jan Teichmann called [Rendezvous Architecture for Data Science in Production](https://towardsdatascience.com/rendezvous-architecture-for-data-science-in-production-79c4d48f12b) that details the Rendezvous architecture concept and benefits further. It's a must read!

## Why Concurrent.ai?

Although there are many benefits that the Rendezvous architecture offers, one major drawback is how high the initial engineering effort is to implement it. With Concurrent.ai, implementing a Rendezvous architecture can now be as simple as writing a few lines of JSON.

## Core Components

![Core Components](https://github.com/concurrentai/concurrentai-core/raw/master/misc/diagrams/Concurrent.ai%20Core%20Components.png)
