# Twitch Test
A toy service based around streaming events using Kafka - Uses Twitch as an event source

## Goal
To utilize Twitch webhooks to populate events in the service and have consumers read from Kafka and perform work on the events.

## TODO
* Work out process for Twitch webhook integration
* Create ingestion service which writes to Kafka
* Create Kafka clients to process events and do $work
