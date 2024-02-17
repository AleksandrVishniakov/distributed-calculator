import React, {useEffect, useState} from "react";
import TextField from "@mui/material/TextField";
import {Button} from "@mui/material";
import {v4 as uuidv4} from "uuid";
import ResponseError from "../../ResponseError";

interface OperatorResponse {
    operationType: number
    durationMS: number
}

class OperatorRequest {
    operationType: number
    durationMS: number

    constructor(operationType: number, durationMS: number) {
        this.operationType = operationType
        this.durationMS = durationMS
    }
}

interface SettingsScreenProps {
    host: string
    onError: (msg: string) => void
}

const formatResponseError = (err: ResponseError): string => {
    return err.code + ". " + err.message
}

const SettingsScreen: React.FC<SettingsScreenProps> = (props) => {
    const [operators, setOperators] = useState(new Array<OperatorResponse>)
    const [plusInputValue, setPlusInputValue] = useState("")
    const [minusInputValue, setMinusInputValue] = useState("")
    const [multiplyInputValue, setMultiplyInputValue] = useState("")
    const [divideInputValue, setDivideInputValue] = useState("")

    const handleButtonClick = () => {
        const ops = [new OperatorRequest(0, parseInt(plusInputValue)),
            new OperatorRequest(1, parseInt(minusInputValue)),
            new OperatorRequest(2, parseInt(multiplyInputValue)),
            new OperatorRequest(3, parseInt(divideInputValue))]

        setOperators(ops)

        saveOperators(ops)
    }

    useEffect(()=> {
        getOperators()
    }, [])

    useEffect(() => {
        if (operators.length != 4) return

        setPlusInputValue(operators[0].durationMS.toString())
        setMinusInputValue(operators[1].durationMS.toString())
        setMultiplyInputValue(operators[2].durationMS.toString())
        setDivideInputValue(operators[3].durationMS.toString())
    }, [operators]);

    const getOperators = async () => {
        const response = await fetch(props.host + "/api/operators", {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });

        const json = await response.json();

        if (response.ok) {
            const resp = json as OperatorResponse[]

            setOperators(resp)
        } else {
            const err = json as ResponseError
            props.onError(formatResponseError(err))

            console.error(json)
        }
    }

    const saveOperators = async (operators: OperatorRequest[]) => {
        const response = await fetch(props.host + "/api/operators", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(operators),
        });


        if (response.ok) {
        } else {
            const err = await response.json() as ResponseError
            props.onError(formatResponseError(err))

            console.error("save operator error")
        }
    }

    return (
        <div style={{display:"flex", flexDirection:"column", height:"300px", justifyContent:"space-around"}}>
            <TextField
                required
                id="standard-required"
                label="Время выполнения +, мс"
                value={plusInputValue}
                variant="standard"
                onChange={(evt: React.ChangeEvent<HTMLInputElement>) => {
                    setPlusInputValue(evt.target.value)
                }}
                style={{
                    width: "150px"
                }}
                type={"number"}
            />

            <TextField
                required
                id="standard-required"
                label="Время выполнения -, мс"
                value={minusInputValue}
                variant="standard"
                onChange={(evt: React.ChangeEvent<HTMLInputElement>) => {
                    setMinusInputValue(evt.target.value)
                }}
                style={{
                    width: "150px"
                }}
                type={"number"}
            />

            <TextField
                required
                id="standard-required"
                label="Время выполнения *, мс"
                value={multiplyInputValue}
                variant="standard"
                onChange={(evt: React.ChangeEvent<HTMLInputElement>) => {
                    setMultiplyInputValue(evt.target.value)
                }}
                style={{
                    width: "150px"
                }}
                type={"number"}
            />

            <TextField
                required
                id="standard-required"
                label="Время выполнения /, мс"
                value={divideInputValue}
                variant="standard"
                onChange={(evt: React.ChangeEvent<HTMLInputElement>) => {
                    setDivideInputValue(evt.target.value)
                }}
                style={{
                    width: "150px"
                }}
                type={"number"}
            />

            <Button variant="contained"
                    onClick={handleButtonClick}
                    style={{
                        width:"150px"
                    }}
            >
                Сохранить
            </Button>
        </div>
    )
}

export default SettingsScreen