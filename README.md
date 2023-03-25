Protocol
========

SWELL is a protocol for ordering events in a distributed network under adversarial conditions. 
It is inspired in tendermint but open for any number of participants. 


The events are classified into epochs. Each epoch consists of a fixed number of time 
slots. Events are incorporated into time slots by the formation of new blocks.

Framework
=========

swell package is a framework for deploying SWELL protocol on any state machine. 

The state machine must expose an interface 
