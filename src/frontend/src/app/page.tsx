"use client";
import { z } from "zod";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";

const formSchema = z.object({
    customerTransactions: z.number().lt(50, {
        message: "requests of this type annot be greater than 50",
    }),
    customerTotal: z.number().lt(50, {
        message: "requests of this type annot be greater than 50",
    }),
    customerSorted: z.number().lt(40, {
        message: "requests of this type annot be greater than 40",
    })
})

export default function Home() {
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            customerTransactions: 0,
            customerTotal: 0,
            customerSorted: 0,
        }
    })

    function handleSubmit(values: z.infer<typeof formSchema>) {

    }

    return (
        <div className="w-dvw h-dvh bg-black flex items-center justify-center p-16 flex-col gap-y-8">
            <Form {...form}>
                <form onSubmit={form.handleSubmit(handleSubmit)} className="h-16 w-full flex justify-around">
                    <FormField
                        control={form.control}
                        name="customerTransactions"
                        render={({field}) => (
                            <FormItem className="flex items-center text-white">
                                <FormLabel>
                                    Customer Transactions
                                </FormLabel>
                                <FormControl>
                                    <Input
                                        type="number"
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    >
                    </FormField>
                    <FormField
                        control={form.control}
                        name="customerTotal"
                        render={({field}) => (
                            <FormItem className="flex items-center text-white">
                                <FormLabel>
                                    Customer Total
                                </FormLabel>
                                <FormControl>
                                    <Input
                                        type="number"
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    >
                    </FormField>
                    <FormField
                        control={form.control}
                        name="customerSorted"
                        render={({field}) => (
                            <FormItem className="flex items-center text-white">
                                <FormLabel>
                                    Customer Sorted
                                </FormLabel>
                                <FormControl>
                                    <Input
                                        type="number"
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </FormItem>
                        )}
                    >
                    </FormField>
                </form>
            </Form>
            <div className="flex flex-row gap-x-12 flex-grow bg-yellow w-full">
                <div className="flex-grow h-full text-white flex flex-col">
                    <h2 className="font-bold text-4xl mb-3">Database</h2>
                    <div className="flex-grow border border-gray-800">
                    </div>
                </div>
                <div className="flex-grow h-full text-white flex-col flex">
                    <h2 className="font-bold text-4xl mb-3">Cache</h2>
                    <div className="flex-grow border border-gray-800">
                    </div>
                </div>
            </div>
        </div>
    );
}
