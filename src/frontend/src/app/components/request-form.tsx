import { Form, FormControl, FormField, FormItem, FormLabel, FormMessage } from "@/components/ui/form";
import { Input } from "@/components/ui/input";
import { zodResolver } from "@hookform/resolvers/zod";
import { useForm } from "react-hook-form";
import { formSchema } from "@/app/types/schemas";
import { z } from "zod";
import { Button } from "@/components/ui/button";

export default function RequestForm({ handler }: {
    handler: (values: z.infer<typeof formSchema>) => void;
}) {
    const form = useForm<z.infer<typeof formSchema>>({
        resolver: zodResolver(formSchema),
        defaultValues: {
            customerTransactions: "0",
            customerTotal: "0",
            customerSorted: "0",
        }
    })

  return (
        <Form {...form}>
            <form onSubmit={form.handleSubmit(handler)} className="h-16 w-full flex justify-around items-center">
                <FormField
                    control={form.control}
                    name="customerTransactions"
                    render={({field}) => (
                        <FormItem className="flex items-center gap-x-2">
                            <FormLabel className="text-white">
                                Customer Transactions
                            </FormLabel>
                            <div className="flex flex-col">
                                <FormControl>
                                    <Input
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </div>
                        </FormItem>
                    )}
                >
                </FormField>
                <FormField
                    control={form.control}
                    name="customerTotal"
                    render={({field}) => (
                        <FormItem className="flex items-center gap-x-2">
                            <FormLabel className="text-white">
                                Customer Total
                            </FormLabel>
                            <div className="flex flex-col">
                                <FormControl>
                                    <Input
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </div>
                        </FormItem>
                    )}
                >
                </FormField>
                <FormField
                    control={form.control}
                    name="customerSorted"
                    render={({field}) => (
                        <FormItem className="flex items-center gap-x-2">
                            <FormLabel className="text-white">
                                Customer Sorted
                            </FormLabel>
                            <div className="flex flex-col">
                                <FormControl>
                                    <Input
                                        {...field}
                                    />
                                </FormControl>
                                <FormMessage />
                            </div>
                        </FormItem>
                    )}
                >
                </FormField>
                <Button type="submit" className="bg-white text-black hover:text-white hover:bg-black">
                    Done
                </Button>
            </form>
        </Form>
  )
}

