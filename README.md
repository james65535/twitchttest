# Twitch Test
A toy service based around streaming events using Kafka - Uses Twitch as an event source

## Goal
To utilize Twitch webhooks to populate events in the service and have consumers read from Kafka and perform work on the events.

## TODO
* Work out process for Twitch webhook integration
* Create ingestion service which writes to Kafka
* Create Kafka client workers to process events and do $work
* Integrate ingestion server and client workers into Knative
* Understand value propasition of Tanzu portfolio
  * Improve process around current YAML manipulation
  * Improve process for multiple contributors to roll out private environments
