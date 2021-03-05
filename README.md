# trace

This requires docker to run.  At the moment this continuously runs a traceroute
to google.com and provides a UI that in theory should provide a network graph
view of the topology.  I'm not doing anything with weights to indicate how many
packets followed one path vs. another, colors to indicate loss, or labels to add
rich data, yet.

If you're using Docker for Mac you will need to enable experimental features for
this to work as ICMP responses are not forwarded back to the host without
experimental features due to VpnKit.  See
https://github.com/docker/for-mac/issues/1381.  Go to Docker Preferences ->
Docker Engine -> set `"experimental": true`, and Apply & Restart.

I have not tried this on Docker for Windows.

To get started:

```
make up
make start
```

Load `http://localhost:50000`.  This should refresh data every 5 seconds.  I
have nothing but room for improvement in the UI, along with most other things.

To stop the test:

```
make stop
```

To stop the service, control+C in the window that's running the server.

I was thinking about adding linked graphics for latency, loss, etc. with
https://metricsgraphicsjs.org/examples.htm.
