import * as d3 from "d3";
import $ from "jquery";

const width = $(document).width() - 100;
const height = $(document).height() * 0.9;
const margin = { top: 10, right: 30, bottom: 10, left: 30 };
const simulationDurationMs = 1250; // 1.25 second
var startTime, endTime;

var simulation = d3
  .forceSimulation()
  .force("charge", d3.forceManyBody().strength(-800))
  .force(
    "link",
    d3.forceLink().id((d) => d.id)
  )
  .force("y", d3.forceY())
  .alphaTarget(1)
  .alphaDecay(0.001)
  .on("tick", ticked);

const graph = d3
  .select("#graph")
  .append("svg")
  .attr("viewBox", [-width / 2, -height / 2, width, height])
  .attr("height", height + margin.top + margin.bottom)
  .attr("width", width + margin.left + margin.right);

graph.append("g").attr("class", "links");
graph.append("g").attr("class", "nodes");

function drawGraph(data) {
  var nodeElements = graph
    .select(".nodes")
    .selectAll(".node")
    .data(data.nodes, (d) => d.id);

  var enterSelection = nodeElements.enter().append("g").attr("class", "node");

  var circles = enterSelection
    .append("circle")
    .attr("r", 20)
    .style("fill", "#69b3a2")
    .merge(enterSelection);

  var labels = enterSelection
    .append("text")
    .text((d) => d.id)
    .attr("x", 10)
    .attr("y", 10)
    .merge(enterSelection);

  nodeElements.exit().transition().duration("2000").remove();

  var linkElements = graph.select(".links").selectAll(".link").data(data.links);
  linkElements.enter().append("line").attr("class", "link").merge(linkElements);
  linkElements.exit().remove();

  startTime = Date.now();
  endTime = startTime + simulationDurationMs;
  simulation.nodes(data.nodes);
  simulation.force("link").links(data.links);
  simulation.alphaTarget(1).restart();
}

function ticked() {
  if (Date.now() < endTime) {
    var nodeElements = graph.select(".nodes").selectAll(".node");
    var linkElements = graph.select(".links").selectAll(".link");

    nodeElements
      .attr("transform", function (d) {
        return "translate(" + d.x + "," + d.y + ")";
      })
      .call(
        d3
          .drag()
          .on("start", dragstarted)
          .on("drag", dragged)
          .on("end", dragended)
      );

    linkElements
      .attr("x1", (d) => d.source.x)
      .attr("y1", (d) => d.source.y)
      .attr("x2", (d) => d.target.x)
      .attr("y2", (d) => d.target.y);
  } else {
    simulation.stop();
  }
}

function dragstarted(event, d) {
  if (!event.active) simulation.alphaTarget(0.3).restart();
  d.fx = d.x;
  d.fy = d.y;
}

function dragged(event, d) {
  d.fx = event.x;
  d.fy = event.y;
}

function dragended(event, d) {
  if (!event.active) simulation.alphaTarget(0);
  d.fx = null;
  d.fy = null;
}

// This is awful.  Come back to:
// https://bl.ocks.org/robyngit/89327a78e22d138cff19c6de7288c1cf
// https://bl.ocks.org/d3noob/6bd13f974d6516f3e491

var metricsHeight = 300;

const metrics = d3
  .select("#metrics")
  .append("svg")
  .attr("height", metricsHeight + margin.top + margin.bottom)
  .attr("width", width + margin.left + margin.right)
  .append("g")
  .attr("transform", "translate(" + margin.left + "," + margin.top + ")");

var x = d3.scaleLinear().range([0, width]);
var xAxis = d3.axisBottom().scale(x);
metrics
  .append("g")
  .attr("transform", "translate(0," + metricsHeight + ")")
  .attr("class", "xAxis");

var y = d3.scaleLinear().range([metricsHeight, 0]);
var yAxis = d3.axisLeft().scale(y);
metrics.append("g").attr("class", "yAxis");

function drawMetrics(data) {
  x.domain([
    0,
    d3.max(data, function (d) {
      return d.date;
    }),
  ]);
  metrics.selectAll(".xAxis").transition().duration(1000).call(xAxis);

  y.domain([
    0,
    d3.max(data, function (d) {
      return d.value;
    }),
  ]);
  metrics.selectAll(".yAxis").transition().duration(1000).call(yAxis);

  var updateMetrics = metrics.selectAll(".line").data([data], function (d) {
    return d.date;
  });

  updateMetrics
    .enter()
    .append("path")
    .attr("class", "line")
    .merge(updateMetrics)
    .transition()
    .duration(1000)
    .attr(
      "d",
      d3
        .line()
        .x(function (d) {
          return x(d.date);
        })
        .y(function (d) {
          return y(d.value);
        })
    )
    .attr("fill", "none")
    .attr("stroke", "steelblue")
    .attr("stroke-width", 2.5);
}

setInterval(function () {
  d3.json("v0/metrics", {
    headers: new Headers({
      TraceHost: "google.com",
    }),
  }).then((data) => {
    var parseTime = d3.timeParse("%Y-%m-%dT%H:%M:%S.%LZ");
    data.forEach(function (d, i) {
      d.date = parseTime(d.date);
    });
    console.log(data);
    drawMetrics(data);
  });
}, 5000);

setInterval(function () {
  d3.json("v0/graph", {
    headers: new Headers({
      TraceHost: "google.com",
    }),
  }).then((data) => {
    console.log(data);
    drawGraph(data);
  });
}, 5000);
