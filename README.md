# trace

This requires docker to run.  At the moment this continuously runs a traceroute
to google.com.

If you're using Docker for Mac you will need to enable experimental features for
this to work as ICMP responses are not forwarded back to the host without
experimental features due to VpnKit.  See
https://github.com/docker/for-mac/issues/1381.  Go to Docker Preferences ->
Docker Engine -> set `"experimental": true`, and Apply & Restart.

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
