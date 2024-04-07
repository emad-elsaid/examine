Chart = G2.Chart

const chart = new Chart({ container: "container" });

chart.options({
    type: "point",
    height: 300,
    data: {
        type: "fetch",
        value: "/timeline.json"
    },
    encode: { y: "Function", x: "Time", shape: "line", size: 10 },
    transform: [{ type: "sortX", channel: "x" }],
    scale: {
        x: { type: "point" },
        y: { },
        color: { type: "ordinal" },
    },
});

chart.render();
