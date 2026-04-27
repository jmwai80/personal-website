---
title: "Decoupling synchronous code paths into Kafka, without regrets"
date: 2026-04-08
tags: [kafka, systems, essay]
description: "Notes from a year of moving a synchronous service into a Publisher → Topic → Consumer pipeline: picking partition keys, surviving rebalances, getting idempotency right, and what we actually put in our DLQ."
draft: false
---

The first version of the service was synchronous, and that was fine — until it wasn't.
A single endpoint had grown to fan out to four downstream systems on every request, and a
slow afternoon at any one of them turned into a bad afternoon for the whole API. We were
absorbing other people's tail latency, and the on-call rotation could feel it.

Pushing the work onto Kafka wasn't a controversial decision; the controversial part was
getting the migration right while production traffic kept moving. This is a write-up of
the choices that mattered, and the ones I'd happily un-make.

## 1. Pick the partition key with delivery semantics in mind

The partition key is the single most consequential design choice in a Kafka pipeline,
because it determines two things at once: **which messages stay ordered**, and
**which consumer instance sees them**. We chose `tenant_id`, which
kept everything ordered per tenant and let us scale consumers horizontally without
worrying about cross-tenant ordering.

The trap: a handful of huge tenants meant a handful of huge partitions. Once we noticed
consumer lag concentrated on one or two partitions, we added a small `shard`
suffix for the heaviest tenants — a partial-key trick that's invisible to consumers but
hugely effective at smoothing the load.

## 2. Idempotency is not optional

At-least-once delivery means duplicates are a feature of the system, not a bug. We
designed every consumer with a deterministic `dedupe_key` derived from the
producer-side message ID, and persisted "already handled" markers next to the side
effects they protected.

```
// pseudo-code
if seen(dedupeKey) {
    return ack()
}
if err := apply(message); err != nil {
    return retry(err)
}
mark(dedupeKey)
ack()
```

That last `mark / ack` ordering is load-bearing: marking after the side effect
but before the ack means a crash between mark and ack just produces a redelivery that
gets dropped at the top of the loop. It's the cheapest insurance you can buy.

## 3. The DLQ is a queue, not a graveyard

A dead-letter queue is only useful if a human ever looks at it. We took two unglamorous
steps that paid for themselves immediately:

- Every DLQ entry carries the original payload, the error class, and the consumer
  version that produced it. No reading code archaeology to debug a six-month-old failure.
- The DLQ has its own dashboard with the same SLO treatment as the main pipeline. If
  it grows faster than we drain it, somebody gets paged.

## 4. Rebalances are a load test you didn't schedule

The first time a deploy triggered a consumer-group rebalance under load, I learned more
about Kafka in ten minutes than I had in a quarter of reading. Two things made the next
rebalance boring:

- Cooperative rebalancing — partitions move incrementally instead of "stop the world."
- A small `session.timeout.ms` bump for consumers doing heavy work, paired
  with explicit `poll()` heartbeats during long handlers.

## 5. The migration pattern: dual-write, then flip

We didn't cut over in one step. The producer wrote to Kafka *and* the legacy
synchronous path for two weeks while we compared outputs in a side-by-side reconciliation
job. When the diff job hit zero for ten straight days, we flipped the read path; a week
later we removed the synchronous code.

> Boring migrations are the ones where the diff job is the most exciting page in the deck.

## What I'd do again

Treating the pipeline as *five independent decisions* — partitioning, idempotency,
retry, DLQ, and rebalance behaviour — instead of "one Kafka thing" was the unlock. Each
of those has its own happy path and its own failure mode, and naming them out loud makes
them debuggable.

The piece I'd push harder on next time: **observability before traffic.**
We had producer and consumer metrics from day one, but we added cross-pipeline tracing
weeks later. Those weeks were the most expensive ones.
