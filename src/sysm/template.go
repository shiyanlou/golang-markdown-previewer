package sysm

import (
	"fmt"
	"net/http"
	"text/template"
)

func Template(w http.ResponseWriter, port int) {
	var style string
	style = "<style>" + DefaultStyle + "</style>"
	templateStr := fmt.Sprintf(`
<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8">
    <title>实验楼项目课</title>
    <script src="http://cdn.staticfile.org/jquery/2.1.1-rc2/jquery.min.js"></script>
    <script src="http://cdn.staticfile.org/socket.io/0.9.16/socket.io.min.js"></script>
    %[2]s
</head>
<body>
<div id="container" style="min-width: 310px; height: 400px; margin: 0 auto"></div>
<script src="http://cdn.staticfile.org/highcharts/4.0.4/highcharts.js"></script>
<script src="http://cdn.staticfile.org/highcharts/4.0.4/modules/exporting.js"></script>
</body>
<script>
var monitor = function(ip){
    Highcharts.setOptions({
        global: {
            useUTC: false
        }
    });

    $('#container').highcharts({
        chart: {
            type: 'spline',
            animation: Highcharts.svg,
            marginRight: 10
        },
        title: {
            text: '实验楼项目课程：Markdown预览器练习'
        },
        credits : {
            href:'http://www.shiyanlou.com',
            position: {
                x:-30,
                y:-30
            },
            style:{
                color:'#191a37',
                fontWeight:'bold'
            },
            text:'http://www.shiyanlou.com'
        },
        xAxis: {
            maxPadding : 0.05,
            minPadding : 0.05,
            type: 'datetime',
            tickWidth:5
        },
        yAxis: {
            title: {
                text: 'Percent(%)'
            },
            plotLines: [{
                value: 0,
                width: 1,
                color: '#808080'
            }]
        },
        tooltip: {
            formatter: function() {
                    return '<b>'+ this.series.name +'</b>('+num+')<br/>'+
                    Highcharts.dateFormat('%H:%M:%S', this.x) +'<br/>'+
                    Highcharts.numberFormat(this.y, 2);
            }
        },
        legend: {
            enabled: true
        },
        exporting: {
            enabled: false
        },
        series: [{
            name: 'CPU',
            data: [
                [(new Date()).getTime(),0]
            ]
        },{
            name: 'Memory',
            data: (function() {
                var data = [];
                data.push([(new Date()).getTime(),0]);
                return data;
            })()
        }]
    });


    var num = 0;
    socket = new WebSocket(ip);
    socket.onmessage = function (evt) {
        data = JSON.parse(evt.data)
        var x = data.time;
        var y1 = data.cpu;
        var y2 = data.mem;
        console.log("time:"+x+",CPU:"+y1+",Memory:"+y2);

        var chart = $('#container').highcharts();
        chart.series[0].addPoint([x, y1], true, (++num>120?true:false));
        chart.series[1].addPoint([x, y2], true, (num>120?true:false));
    };
}


window.onload = function(){
    var ip = "ws://localhost:%[1]d";
    monitor(ip);
}
</script>
</html>
    `, port, style)

	var (
		t   *template.Template
		err error
	)

	if t, err = template.New("template").Parse(templateStr); err != nil {
		panic(err)
	}

	if err = t.Execute(w, nil); err != nil {
		panic(err)
	}
}

var DefaultStyle = `
/* default css */
body {
  padding: 50px;
  font: 14px "Lucida Grande", Helvetica, Arial, sans-serif;
}

a {
  color: #00B7FF;
}

.center {
    width: 960px;
    margin-left: auto;
    margin-right: auto;
}
`
