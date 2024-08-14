import { z } from "zod";

export const formSchema = z.object({
    customerTransactions: z.string().regex(new RegExp(/[0-9]+/), {
        message: "please provide a valid number",
    }).refine((val) => Number(val) <= 50, {
        message: "request type cannot be greater than 50"
    }),
    customerTotal: z.string().regex(new RegExp(/[0-9]+/), {
        message: "please provide a valid number",
    }).refine((val) => Number(val) <= 50, {
        message: "request type cannot be greater than 50"
    }),
    customerSorted: z.string().regex(new RegExp(/[0-9]+/), {
        message: "please provide a valid number",
    }).refine((val) => Number(val) <= 40, {
        message: "request type cannot be greater than 40"
    })
})
