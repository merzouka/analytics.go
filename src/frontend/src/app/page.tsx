"use client";
import { formSchema } from "@/app/types/schemas";
import { z } from "zod";
import RequestForm from "./components/request-form";
import { useRef } from "react";


function onMessage(dbRef: React.RefObject<HTMLDivElement>, cacheRef: React.RefObject<HTMLDivElement>, response: SSEResponse, eventSource: EventSource) {
    const { source, data, done, duration, success } = response;
    if (done) {
        eventSource.close()
        let dbResult = {}
        let cacheResult = {}
        for (let metric of Object.keys(data)) {
            // @ts-ignore
            dbResult[metric] = data[metric].database
            // @ts-ignore
            cacheResult[metric] = data[metric].cache
        }
        if (dbRef != null && dbRef.current != null) {
            dbRef.current.innerHTML += `---<br>${JSON.stringify(dbResult)}`
        }
        if (cacheRef != null && cacheRef.current != null) {
            cacheRef.current.innerHTML += `---<br>${JSON.stringify(cacheResult)}`
        }
        return
    }

    const result = `${data}: ${success} [${duration}]<br>`
    if (dbRef != null && dbRef.current != null && source == "database") {
        dbRef.current.innerHTML += result
    }
    if (cacheRef != null && cacheRef.current != null && source == "cache") {
        cacheRef.current.innerHTML += result
    }
}

export default function Home() {
    const outputRefDB = useRef<HTMLDivElement>(null);
    const outputRefCache = useRef<HTMLDivElement>(null);

    function handleSubmit(values: z.infer<typeof formSchema>) {
        const url = `http://localhost/bulk?requests=customer_transactions:${values.customerTransactions},customer_total:${values.customerTotal},customer_sorted:${values.customerSorted}`
        if (outputRefDB != null && outputRefDB.current != null) {
            outputRefDB.current.innerHTML = ""
        }
        if (outputRefCache != null && outputRefCache.current != null) {
            outputRefCache.current.innerHTML = ""
        }
        const source = new EventSource(url)
        source.onmessage = (e) => {
            onMessage(outputRefDB, outputRefCache, JSON.parse(e.data), source)
        }
    }

    return (
        <div className="w-dvw h-dvh bg-black flex items-center justify-center p-16 flex-col gap-y-8">
            <RequestForm handler={handleSubmit} />
            <div className="flex flex-row gap-x-12 flex-grow bg-yellow w-full">
                <div className="h-full text-white flex flex-col w-full">
                    <h2 className="font-bold text-4xl mb-3">Database</h2>
                    <div className="flex-grow border border-gray-800 relative">
                        <div className="absolute top-0 right-0 left-0 bottom-0 overflow-scroll" ref={outputRefDB}>
                        </div>
                    </div>
                </div>
                <div className="h-full text-white flex-col flex w-full">
                    <h2 className="font-bold text-4xl mb-3">Cache</h2>
                    <div className="flex-grow border border-gray-800" ref={outputRefCache}>
                    </div>
                </div>
            </div>
        </div>
    );
}
