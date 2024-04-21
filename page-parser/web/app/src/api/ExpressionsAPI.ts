import {APIError} from "./APIError";

interface Expression {
    id: number
    expression: string
    createdAt: Date
    finishedAt: Date
    status: number
    result: number
}

export class ExpressionsAPI {
    private readonly host: string

    private readonly getToken = () => localStorage.getItem("token") || ""

    constructor(host: string = "http://localhost:8000") {
        this.host = host
    }

    public async newExpression(expression: string) {
        const url = `${this.host}/api/expression`

        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${this.getToken()}`,
            },
            body: JSON.stringify({
                expression: expression,
                idempotencyKey: new Date()
            })
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }
    }

    public async getExpressions(): Promise<Expression[]> {
        const url = `${this.host}/api/expressions`

        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
                "Authorization": `Bearer ${this.getToken()}`,
            },
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }

        return await response.json() as Expression[]
    }
}

function apiErrorToString(apiError: APIError): string {
    return apiError.code + ": " + apiError.message
}