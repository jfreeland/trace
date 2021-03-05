// import * as d3 from "d3";

const width = $(document).width() - 100;
const height = $(document).height() - 100;
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

const svg = d3
  .select("div")
  .append("svg")
  .attr("viewBox", [-width / 2, -height / 2, width, height])
  .attr("height", height + margin.top + margin.bottom)
  .attr("width", width + margin.left + margin.right);

svg.append("g").attr("class", "links");
svg.append("g").attr("class", "nodes");

function draw(data) {
  var nodeElements = svg
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

  var linkElements = svg.select(".links").selectAll(".link").data(data.links);
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
    var nodeElements = svg.select(".nodes").selectAll(".node");
    var linkElements = svg.select(".links").selectAll(".link");

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

d3.json("files/first.json").then((data) => {
  draw(data);
});

// setTimeout(function () {
//     d3.json('files/second.json')
//         .then(data => {
//             draw(data);
//         });
// }, 3000);

// setTimeout(function () {
//     d3.json('files/third.json')
//         .then(data => {
//             draw(data);
//         });
// }, 6000);

document.getElementById("btn1").addEventListener("click", function () {
  d3.json("files/first.json").then((data) => {
    draw(data);
  });
});

document.getElementById("btn2").addEventListener("click", function () {
  d3.json("files/second.json").then((data) => {
    draw(data);
  });
});

document.getElementById("btn3").addEventListener("click", function () {
  d3.json("files/third.json").then((data) => {
    draw(data);
  });
});
