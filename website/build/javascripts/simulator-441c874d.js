function graph_settings(){return{chart:{type:"spline"},title:{text:null},xAxis:{title:{text:"Time (sec)"}},yAxis:{title:{text:"Convergence %"},min:0,max:100},tooltip:{formatter:function(){return"<b>"+Math.round(1e3*this.y)/1e3+"%</b><br/>"}},legend:{enabled:!1},series:[{name:"Convergence Rate",data:[]}]}}function create_graph(){return $("#graph").highcharts(graph_settings()),$("#graph").highcharts()}function update_interval(t,e){var n=t.value,r=Number(n);return isNaN(r)?(alert("Gossip interval must be a number!"),void 0):0>=r?(alert("Gossip interval must be a positive value!"),void 0):(e.interval=r,e.draw(),console.log("Redraw with interval set to: "+r),void 0)}function update_fanout(t,e){var n=t.value,r=Number(n);return isNaN(r)?(alert("Gossip fanout must be a number!"),void 0):0>=r?(alert("Gossip fanout must be a positive value!"),void 0):(e.fanout=r,e.draw(),console.log("Redraw with fanout set to: "+r),void 0)}function update_nodes(t,e){var n=t.value,r=Number(n);return isNaN(r)?(alert("Node count must be a number!"),void 0):1>=r?(alert("Must have at least one node"),void 0):(e.nodes=r,e.draw(),console.log("Redraw with nodes set to: "+r),void 0)}function update_packetloss(t,e){var n=t.value,r=Number(n);return isNaN(r)?(alert("Packet loss must be a number!"),void 0):0>r||r>=100?(alert("Packet loss must be greater or equal to 0 and less than 100"),void 0):(e.packetLoss=r/100,e.draw(),console.log("Redraw with packet loss set to: "+r),void 0)}function update_failed(t,e){var n=t.value,r=Number(n);return isNaN(r)?(alert("Failure rate must be a number!"),void 0):0>r||r>=100?(alert("Failure rate must be greater or equal to 0 and less than 100"),void 0):(e.nodeFail=r/100,e.draw(),console.log("Redraw with failure rate set to: "+r),void 0)}var Simulator=Class.$extend({__init__:function(t,e,n){this.graph=t,this.bytes=e,this.maxConverge=n,this.interval=.2,this.fanout=3,this.nodes=30,this.packetLoss=0,this.nodeFail=0},convergenceAtRound:function(t){var e=.5*this.fanout/this.nodes*(1-this.packetLoss)*(1-this.nodeFail),n=this.nodes/(1+(this.nodes+1)*Math.pow(Math.E,-1*e*this.nodes*t));return n/this.nodes},roundLength:function(){return this.interval},seriesData:function(){for(var t=[],e=0,n=0,r=this.roundLength();e<this.maxConverge&&100>n;)e=this.convergenceAtRound(n),t.push([n*r,100*e]),n++;return t},bytesUsed:function(){var t=this.roundLength(),e=1/t,n=1400,r=n*this.fanout*e;return 2*r},draw:function(){var t=this.seriesData();this.graph.series[0].setData(t,!1),this.graph.redraw();var e=8*this.bytesUsed(),n=Math.round(10*(e/1024))/10;this.bytes.html(""+n)}});$(function(){var t=$("#bytes"),e=create_graph(),n=new Simulator(e,t,.9999);n.draw();var r=$("#interval");r.change(function(){update_interval(r[0],n)});var i=$("#fanout");i.change(function(){update_fanout(i[0],n)});var o=$("#nodes");o.change(function(){update_nodes(o[0],n)});var a=$("#packetloss");a.change(function(){update_packetloss(a[0],n)});var s=$("#failed");s.change(function(){update_failed(s[0],n)})});