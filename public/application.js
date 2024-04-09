const Chart = G2.Chart

const id = 'container'
const container = document.getElementById(id)
const chart = new Chart({
    width: container.clientWidth,
    height: container.clientHeight,
    container: id,
});

chart.options({
    data: {
        type: 'fetch',
        value: '/timeline.json',
    },
    type: 'point',
    encode: {
        x: 'Time',
        y: 'Function'
    },
    style: {
        stroke: 'black'
    },
    slider: {
        x: {}
    },
    scrollbar: {
        y: {}
    },
    axis: {
        x: false
    }
})

chart.render();
