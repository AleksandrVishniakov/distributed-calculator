import {APIError} from "./APIError";

interface Worker {
    id: number
    url: string
    executors: number
    lastModified: Date
}

interface Task {
    leftResult: number
    operationType: number
    rightResult: number
    status: number
}

export class WorkerAPI {
    private readonly host: string;

    constructor(host: string = "http://localhost:8000") {
        this.host = host;
    }

    public async getWorkers(): Promise<Array<Worker>> {
        const url = `${this.host}/api/workers`

        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }

        return await response.json() as Array<Worker>
    }

    public async getTasks(workerID: number): Promise<Array<Task>> {
        const url = `${this.host}/api/worker/${workerID}/tasks`

        const response = await fetch(url, {
            method: "GET",
            headers: {
                "Content-Type": "application/json",
            }
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }

        return await response.json() as Array<Task>
    }
}

function apiErrorToString(apiError: APIError): string {
    return apiError.code + ": " + apiError.message
}