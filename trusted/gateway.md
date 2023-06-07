Gateway provides an interface between the p2p validator network and the outside 
world. It is used to inform about the mint of new blocks as seen by each node,
to trasmit those blocks if requested. And finally to submit, on behalf of 
end-users, new actions to be considered and incorporated by the network. 

It is also an interface for routing information outside the network. Gateways
might connect to one-another, or to other specialized applications like data
indexers or even end-user applications. 

The basic functionality is


A trusted connection is based on signed messages:

Message Template:
    message type
    time stamp
    message data
    signature

A particular message is the confirmation of a message. It is characterized by a
message type = message confirmation and message data = hash of confirmed message.

A trusted connection has two basic funcionalities: to receive/transmit new 
blocks and to receive/transmit new actions. The associated messages are

NewBlockAvailable:
    block age
    block hash

RequestBlock:
    block age
    block hash

SendBlock:
    block data 

SendAction:
    action data

