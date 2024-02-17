import React, {useEffect, useState} from "react";
import TextField from '@mui/material/TextField';
import "./CalculationScreen.module.css"
import {Button} from "@mui/material";
import {v4 as uuidv4} from 'uuid';
import Paper from '@mui/material/Paper';
import Table from '@mui/material/Table';
import TableBody from '@mui/material/TableBody';
import TableCell from '@mui/material/TableCell';
import TableContainer from '@mui/material/TableContainer';
import TableHead from '@mui/material/TableHead';
import TablePagination from '@mui/material/TablePagination';
import TableRow from '@mui/material/TableRow';
import ResponseError from "../../ResponseError";

interface CalculationResponse {
    id: number
}

class CalculationRequest {
    expression: string
    idempotencyKey: string

    constructor(expression: string, key: string) {
        this.expression = expression
        this.idempotencyKey = key
    }
}

interface ExpressionResponse {
    id: number
    expression: string
    createdAt: Date
    finishedAt: Date
    status: number
    result: number
}

interface CalculationScreenProps {
    onClick: (expression: string) => void
    onError: (msg: string) => void
    host: string
}

interface Column {
    id: 'id' | 'expression' | 'createdAt' | 'finishedAt' | 'status' | 'result';
    label: string;
    minWidth?: number;
    format?: (value: number) => string;
}

const columns: readonly Column[] = [
    //{ id: 'id', label: 'Id', minWidth: 170 },
    {id: 'expression', label: 'Expression', minWidth: 100},
    {
        id: 'result',
        label: 'Result',
        minWidth: 50,
    },
    {
        id: 'createdAt',
        label: 'Created At',
        minWidth: 170,
    },
    {
        id: 'finishedAt',
        label: 'Finished At',
        minWidth: 170,
    },
    {
        id: 'status',
        label: 'Status',
        minWidth: 50,
        format: (value: number): string => {
            switch (value) {
                case 0:
                    return "Создано"
                case 1:
                    return "В очереди"
                case 2:
                    return "Выполняется"
                case 3:
                    return "Готово"
                case 4:
                    return "Ошибка вычисления"
            }

            return ""
        }
    },
];

const formatResponseError = (err: ResponseError): string => {
    return err.code + ". " + err.message
}

const CalculationScreen: React.FC<CalculationScreenProps> = (props) => {
    const [inputValue, setInputValue] = useState("")
    const [expressions, setExpressions] = useState(new Array<ExpressionResponse>)
    const [page, setPage] = React.useState(0);
    const [rowsPerPage, setRowsPerPage] = React.useState(10);
    const fetchExpression = async (value: string) => {
        const request = new CalculationRequest(value, uuidv4())

        const response = await fetch(props.host + "/api/expression", {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(request),
        });

        const json = await response.json();

        if (response.ok) {
            const resp = json as CalculationResponse

            await getExpressions()
        } else {
            const err = json as ResponseError

            props.onError(formatResponseError(err))
            console.error(json)
        }
    }

    const getExpressions = async () => {
        const response = await fetch(props.host + "/api/expressions", {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });

        const json = await response.json();

        if (response.ok) {
            const resp = json as ExpressionResponse[]

            setExpressions(resp)
        } else {
            console.error(json)
        }
    }

    useEffect(() => {
        getExpressions()

        const interval = setInterval(() => {
            getExpressions()
        }, 5000)

        return () => clearInterval(interval)
    }, [])

    const handleButtonClick = async () => {
        await fetchExpression(inputValue)
        props.onClick(inputValue)
    }

    const handleChangePage = (event: unknown, newPage: number) => {
        setPage(newPage);
    };

    const handleChangeRowsPerPage = (event: React.ChangeEvent<HTMLInputElement>) => {
        setRowsPerPage(+event.target.value);
        setPage(0);
    };

    const formatColumn = (column: Column, expr: ExpressionResponse): string => {
        const value = expr[column.id]

        if (column.format && typeof value === 'number') {
            return column.format(value)
        } else if ((column.id === "result" || column.id === "finishedAt") && expr.status != 3) {
            return ""
        }

        if (column.id === "createdAt" || column.id === "finishedAt") {
            const date = new Date(value)

            return format(date.getDate()) + "." + format(date.getMonth() + 1) + "." + date.getFullYear() + " " + format(date.getHours()) + ":" + format(date.getMinutes()) + ":" + format(date.getSeconds())
        }

        return value.toString()
    }

    return (
        <div className="CalculationScreen">
            <h2>Посчитать выражение</h2>
            <p>Введите выражение, используя цифры и символы +, -, *, /, (, )</p>

            <TextField
                required
                id="standard-required"
                label="Выражение"
                defaultValue={inputValue}
                variant="standard"
                onChange={(evt: React.ChangeEvent<HTMLInputElement>) => {
                    setInputValue(evt.target.value)
                }}
                style={{
                    width: "500px"
                }}
            />

            <Button variant="contained"
                    onClick={handleButtonClick}
            >
                Посчитать
            </Button>

            {
                expressions ?
                    <Paper sx={{width: '100%', overflow: 'hidden'}}>
                        <TableContainer sx={{maxHeight: 440}}>
                            <Table stickyHeader aria-label="sticky table">
                                <TableHead>
                                    <TableRow>
                                        {columns.map((column) => (
                                            <TableCell
                                                key={column.id}
                                                style={{minWidth: column.minWidth}}
                                            >
                                                {column.label}
                                            </TableCell>
                                        ))}
                                    </TableRow>
                                </TableHead>
                                <TableBody>
                                    {
                                        expressions
                                            .slice(page * rowsPerPage, page * rowsPerPage + rowsPerPage)
                                            .map((expr) => {
                                                return (
                                                    <TableRow hover role="checkbox" tabIndex={-1} key={expr.id}>
                                                        {columns.map((column) => {
                                                            const value = expr[column.id];
                                                            return (
                                                                <TableCell key={column.id}>
                                                                    {
                                                                        formatColumn(column, expr)
                                                                    }
                                                                </TableCell>
                                                            );
                                                        })}
                                                    </TableRow>
                                                );
                                            })
                                    }
                                </TableBody>
                            </Table>
                        </TableContainer>
                        <TablePagination
                            rowsPerPageOptions={[10, 25, 100]}
                            component="div"
                            count={expressions.length}
                            rowsPerPage={rowsPerPage}
                            page={page}
                            onPageChange={handleChangePage}
                            onRowsPerPageChange={handleChangeRowsPerPage}
                        />
                    </Paper>
                    : <span></span>
            }
        </div>
    )
}

function format(n: number): string {
    if (n > 10) {
        return n.toString()
    }

    return "0" + n.toString()
}

export default CalculationScreen