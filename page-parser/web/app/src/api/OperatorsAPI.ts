import {APIError} from "./APIError";

interface Operator {
    operationType: number
    durationMS: number
}

export class OperatorsAPI {
    private readonly host: string

    constructor(host: string = "http://localhost:8000") {
        this.host = host
    }

    public async getOperators(): Promise<Operator[]> {
        const url = `${this.host}/api/operators`

        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            },
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }

        return await response.json() as Operator[]
    }

    public async setOperators(plusTime: number, minusTime: number, divideTime: number, multiplyTime: number): Promise<void> {
        const operators: Operator[] = [
            {
                operationType: 0,
                durationMS: plusTime
            }, {
                operationType: 1,
                durationMS: minusTime
            }, {
                operationType: 2,
                durationMS: divideTime
            }, {
                operationType: 3,
                durationMS: multiplyTime
            }
        ]

        const url = `${this.host}/api/operators`

        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify(operators)
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }
    }
}

function apiErrorToString(apiError: APIError): string {
    return apiError.code + ": " + apiError.message
}