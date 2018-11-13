# ingestor

[![Build Status](https://travis-ci.org/the4thamigo-uk/ingestor.svg?branch=master)](https://travis-ci.org/the4thamigo-uk/ingestor?branch=master)
[![Coverage Status](https://coveralls.io/repos/the4thamigo-uk/ingestor/badge.svg?branch=master&service=github)](https://coveralls.io/github/the4thamigo-uk/ingestor?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/the4thamigo-uk/ingestor)](https://goreportcard.com/report/github.com/the4thamigo-uk/ingestor)
[![Godoc](https://godoc.org/github.com/the4thamigo-uk/ingestor?status.svg)](https://godoc.org/github.com/the4thamigo-uk/ingestor)

## Description

This is a simple project for a service that ingests data files in CSV format (containing contact details in this case). The 
service provides an API implemented with [grpc](https://github.com/grpc/grpc-go), from which clients can extract a stream
of typed _Contact_ records from data files of [this form](./testdata/data.csv).

Assumptions : 

- The data files are assumed to be static on disk.
- The ID in the CSV is assumed not to be a global identifier, and therefore is not exposed to clients of the ingestor service.
- Mobile numbers are assumed to be correct UK phone numbers and are canonicalised into the standard [international format](https://www.cm.com/blog/how-to-format-international-telephone-numbers/)
- The canonicalised mobile number is considered the id of the _Contact_ entity. Thus, a given email can have multiple mobile numbers, but a 
mobile number can have, at most, one email. (Depending on the business needs the _email_ could instead be used as the id, or we could allow for multiple emails to be stored against a given phone number).
- If there are two records for the same canonicalised phone number the contact details for the second record will overwrite the first.
- It is assumed that grpc streaming will provide the 'backpressure' flow-control required to prevent the reader service becoming overloaded with data sent in the stream (needs further research to properly understand this).

Limitations :

- The tests are not complete
- The reader service currently puts the received contact data into an in-memory store or into [Cassandra](http://cassandra.apache.org/). However, the store is de-coupled with an appropriate interface, so other
implementations can be used. The interface can support batched operations.
- The ingestor service currently creates a new UUID identifier for each new file added. The identifiers are not persisted so that a restart will
create a new UUID and hence trigger a polling reader to re-download the data.
- The reader service currently polls the ingestor service for new data files. It may be more timely to listen to an event stream from the server, but this
adds complexity, and it is thought that 'instantaneous' updates of such data are probably not required.
- The reader service considers it has downloaded the file after it completes reading it in its entirity. This means if the download is interrupted, the
client will re-download the entire file in the next poll. It could be sensible to introduce the ability to 're-join' an interrupted download, but this adds complexity.
- The CSV file is expected to have at least a header record containing the four contact fields (in any order). Each data record must contain non-blank data for each 
of these four fields.
- The grpc connection is not secure.

## Getting Started

To build the grpc protocol interfaces

    go generate ./pkg/protocol

To run unit tests :

    go test -v ./pkg/...

To build and run the ingestor service serving the given datafile :

    go build ./cmd/ingestor && ./ingestor -l :8080 -f ./testdata/data.csv

To build and run the cli tool to add a data file to the ingestor use :

    go build ./cmd/cli && ./cli -s :8080 add -f ./testdata/data.csv

To build and run the example reader service :

    go build ./cmd/reader && ./reader -c :8080

Alternatively, to run an example environment containing an ingestor service, and a reader service connected to a Cassandra database run :

    docker-compose build
    docker-compose up

Then you can add the test data file to the ingestor service using the cli tool :

    docker-compose exec ingestor /cli add -f /data/data.csv

Then you can wait for the records to be imported by the reader and verify them in Cassandra using csql :

    docker-compose exec cassandra cqlsh -e "SELECT * from reader.contacts"
