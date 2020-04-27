# Twitch Test
A toy service based around streaming events using Kafka - Uses Twitch as an event source

## Goal
To utilize Twitch webhooks to populate events in the service and have consumers read from Kafka and perform work on the events.

## TODO
* Create Kafka client workers to process events and do $work
* Integrate ingestion server and client workers into Knative
* Understand value proposition of Tanzu portfolio
  * Improve process around current YAML manipulation
  * Improve process for multiple contributors to roll out private environments
  * Integrate Wavefront for metrics and tracing
