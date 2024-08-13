// const eventSource = new EventSource("http://localhost:8080/bulk?requests=customer_transactions:10,customer_total:20,customer_sorted:10")
const eventSource = new EventSource("http://localhost:8080/bulk?requests=customer_sorted:40,customer_transactions:50")
const output = document.getElementById("sse-data")
eventSource.onmessage = e => {
    const { data, source, duration, done } = JSON.parse(e.data)
    if (done) {
        eventSource.close()
        console.log(data)
        return
    }
    output.innerHTML = output.innerHTML + `${source}: ${data} in ${duration} (${done})<br>`
}
