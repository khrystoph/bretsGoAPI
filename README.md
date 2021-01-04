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

## Usage
The API currently has only one function: to make an API query against the finnhub API and return a stock quote to you. There's not a large advantage to this function just yet, but it will be used for more complex functions (like looking up multiple stocks in a portfolio to get total value of the portfolio). If you're using this in a production-like environment, you should be ensuring you're running the non-testing version (ie. feeding it a FQDN) so that it runs in SSL mode AND you should ensure that you have your own API key that you can feed to the docker container.  

Avoid direct input of the API key (and use a secrets manager or lookup from some other secure source instead of providing it on the cli).  The API will allow for no credentials, but you will hit the limit on how many API calls you can make to finnhub without a key.

There are three flags that you can provide to either the binary or to a docker container:
1. `-testing` - This tells the go binary whether or not to run using TLS and auto-cert or to use port 8080.
1. `-d` or `-domain` - This flag tells the binary what hostname we're looking for (FQDN) to set up the SSL cert, but also validates whether or not the request is going to the valid hostname.
1. `-apikey` - Provide your Finnhub API key to the binary or to the docker container.

To run the docker container, you need to run the following command (or similar):  

```
docker run -d --rm --name api_test -p 80:8080 <name_of_container_build_name> -d <your_domain_name_here> -testing -apikey <finnhub API key>
```

Once the docker container is running, you can interact with the api like this for testing purposes:

```
$ curl -X POST -H "Content-type: application/json" -H "Accept: application/json" -d '{"ticker":"BEP"}' http://localhost/api/v1/quote
```



## Notes
Some of these features may never get implemented, this is sort of a working list of ideas to track. 
There will eventually be an auth system set up (specifically before this accepts more than just GET and POST requests, when there are mutating requests
or sensitive information, this will be set up before such features go live).

Eventually, this will be broken up into modules because this will eventually get difficult to manage.

Unit testing will also be added in the near future to ensure consistent performance and to catch error codes.

There is also a flag to use testing mode which disables SSL/TLS handling and runs the server on 8080. This makes for quick testing of the API functions
