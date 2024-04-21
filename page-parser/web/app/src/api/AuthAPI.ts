import {APIError} from "./APIError";

interface LoginResponse {
    token: string
}

interface RegisterResponse {
    id: number
}

export class AuthAPI {
    private readonly host: string

    constructor(host: string ="http://localhost:8005") {
        this.host = host
    }

    public async login(login: string, password: string) {
        const url = `${this.host}/api/v1/login`

        const response = await fetch(url, {
            method: "POST",
            headers: {
                "Content-Type": "application/json",
            },
            body: JSON.stringify({
                login: login,
                password: password,
            }),
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }

        const token = await response.json() as LoginResponse

        localStorage.setItem("token", token.token)
    }

    public async register(login: string, password: string) {
        const url = `${this.host}/api/v1/register`

        const response = await fetch(url, {
            method: "POST",
            body: JSON.stringify({
                login: login,
                password: password,
            }),
        })

        if (!response.ok) {
            const apiError = await response.json() as APIError

            throw new Error(apiErrorToString(apiError))
        }

        const id = (await response.json() as RegisterResponse).id

        console.log(id)
    }
}

function apiErrorToString(apiError: APIError): string {
    return apiError.code + ": " + apiError.message
}