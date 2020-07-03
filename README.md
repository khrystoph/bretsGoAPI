# bretsGoAPI

## Intro
I created this api mostly as an exercise in understanding how APIs work, but also to make certain tasks easier for me to execute on a regular basis.
I'm starting with simple API verbs first, such as financial verbs and building from there. Additionally, GET will be all that is initially supported.

## Planned Features
This list will change as I continue to develop more of my needs, but will follow a rough implementation path (not necessarily in order):
1. Stock Quotes for individual stocks - complete
1. Look Up Stock Porftolio value by JSON inputs (ticker, # shares)
1. Track Gain/Loss per stock (will require either database or input file)
1. Track cost changes over time (will require setting up a database)
1. Internal DNS tracking/updating
1. Track resource usages (such as NAS)
1. Track uptimes
1. Track pending package updates
1. Send notifications when there are pending package updates
1. Trigger Workflows to update packages

## Notes
Some of these features may never get implemented, this is sort of a working list of ideas to track. 
There will eventually be an auth system set up (specifically before this accepts more than just GET and POST requests, when there are mutating requests
or sensitive information, this will be set up before such features go live).

Eventually, this will be broken up into modules because this will eventually get difficult to manage.

Unit testing will also be added in the near future to ensure consistent performance and to catch error codes.

There is also a flag to use testing mode which disables SSL/TLS handling and runs the server on 8080. This makes for quick testing of the API functions
