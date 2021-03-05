!function(){"use strict";const t=$(document).width()-100,e=$(document).height()-100,n=10,a=30,r=10,s=30;var l,i,o=d3.forceSimulation().force("charge",d3.forceManyBody().strength(-800)).force("link",d3.forceLink().id((t=>t.id))).force("y",d3.forceY()).alphaTarget(1).alphaDecay(.001).on("tick",(function(){if(Date.now()<i){var t=c.select(".nodes").selectAll(".node"),e=c.select(".links").selectAll(".link");t.attr("transform",(function(t){return"translate("+t.x+","+t.y+")"})).call(d3.drag().on("start",f).on("drag",u).on("end",g)),e.attr("x1",(t=>t.source.x)).attr("y1",(t=>t.source.y)).attr("x2",(t=>t.target.x)).attr("y2",(t=>t.target.y))}else o.stop()}));const c=d3.select("div").append("svg").attr("viewBox",[-t/2,-e/2,t,e]).attr("height",e+n+r).attr("width",t+s+a);function d(t){var e=c.select(".nodes").selectAll(".node").data(t.nodes,(t=>t.id)),n=e.enter().append("g").attr("class","node");n.append("circle").attr("r",20).style("fill","#69b3a2").merge(n),n.append("text").text((t=>t.id)).attr("x",10).attr("y",10).merge(n),e.exit().transition().duration("2000").remove();var a=c.select(".links").selectAll(".link").data(t.links);a.enter().append("line").attr("class","link").merge(a),a.exit().remove(),l=Date.now(),i=l+1250,o.nodes(t.nodes),o.force("link").links(t.links),o.alphaTarget(1).restart()}function f(t,e){t.active||o.alphaTarget(.3).restart(),e.fx=e.x,e.fy=e.y}function u(t,e){e.fx=t.x,e.fy=t.y}function g(t,e){t.active||o.alphaTarget(0),e.fx=null,e.fy=null}c.append("g").attr("class","links"),c.append("g").attr("class","nodes"),d3.json("files/first.json").then((t=>{d(t)})),document.getElementById("btn1").addEventListener("click",(function(){d3.json("files/first.json").then((t=>{d(t)}))})),document.getElementById("btn2").addEventListener("click",(function(){d3.json("files/second.json").then((t=>{d(t)}))})),document.getElementById("btn3").addEventListener("click",(function(){d3.json("files/third.json").then((t=>{d(t)}))}))}();
//# sourceMappingURL=bundle.js.map
