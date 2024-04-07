Chart = G2.Chart

const chart = new Chart({
    width: window.innerWidth,
    height: window.innerHeight,
    container: "container",
});

chart
    .point()
    .data({
        type: 'fetch',
        value: '/timeline.json',
    })
    .encode('x', 'Time')
    .encode('y', 'Function')
    .encode('color', 'File')
    .axis('x', false)

chart.render();
