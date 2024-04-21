import React, {useEffect} from "react";

import './CalculationsScreen.css'
import {ExpressionsAPI} from "../../api/ExpressionsAPI";

interface Expression {
    id: number
    expression: string
    createdAt: Date
    finishedAt: Date
    status: number
    result: number
}

interface Props {
    expressionsAPI: ExpressionsAPI
    onError: (msg: string) => void
}

const CalculationsScreen: React.FC<Props> = ({expressionsAPI, onError}) => {
    const [expressionInputValue, setExpressionInputValue] = React.useState('');
    const [expressions, setExpressions] = React.useState(new Array<Expression>());

    const getExpressions = () => {
        expressionsAPI.getExpressions()
            .then((expr) => {
                setExpressions(expr)
            })
            .catch((e) => {
                onError(e.toString())
            })
    }

    const newExpression = (expr: string) => {
        expressionsAPI.newExpression(expr)
            .catch((e) => {
                onError(e.toString())
            })
    }

    useEffect(() => {
        getExpressions()
        const interval = setInterval(() => {
            getExpressions()
        }, 5000)

        return () => clearInterval(interval)
    }, [])

    const handleSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
        evt.preventDefault()
        newExpression(expressionInputValue)
    }

    const onInputChange = (callback: React.Dispatch<React.SetStateAction<string>>) => {
        return (evt: React.ChangeEvent<HTMLInputElement>) => {
            callback(evt.target.value)
        }
    }

    return (
        <section className="CalculationsScreen">
            <h2 className="CalculationsScreen__title">
                Вычисление выражений
            </h2>

            <p>Введите выражение, используя цифры и символы +, -, *, /, (, )</p>

            <form
                className="CalculationsScreen__form"
                onSubmit={handleSubmit}
            >
                <input
                    type="text"
                    placeholder="Выражение"
                    value={expressionInputValue}
                    onChange={onInputChange(setExpressionInputValue)}
                />
                <button type="submit">=</button>
            </form>

            <h2 className="CalculationsScreen__title">
                Результаты вычислений
            </h2>

            <table>
                <thead>
                    <tr>
                        <th>Выражение</th>
                        <th>Результат</th>
                        <th>Создано</th>
                        <th>Закочено</th>
                        <th>Статус</th>
                    </tr>
                </thead>

                <tbody>
                    {expressions ? expressions.map((expr) => {
                        return (
                            <tr key={expr.id}>
                                <td>{expr.expression}</td>
                                <td>{expr.result}</td>
                                <td>{formatDate(expr.createdAt)}</td>
                                <td>{formatDate(expr.finishedAt)}</td>
                                <td>{formatStatus(expr.status)}</td>
                            </tr>
                        )
                    }) : null}
                </tbody>
            </table>
        </section>
    )
}

const weekDays = new Map<number, string>([
    [1, "Вс"],
    [2, "Пн"],
    [3, "Вт"],
    [4, "Ср"],
    [5, "Чт"],
    [6, "Пт"],
    [7, "Сб"],
])

const months = new Map<number, string>([
    [1, "Январь"],
    [2, "Февраль"],
    [3, "Март"],
    [4, "Апрель"],
    [5, "Май"],
    [6, "Июнь"],
    [7, "Июль"],
    [8, "Август"],
    [9, "Сентябрь"],
    [10, "Октябрь"],
    [11, "Ноябрь"],
    [12, "Декабрь"],
])

const formatDate = (d: Date): string => {
    const date = new Date(d)
    let month = months.get(date.getMonth() + 1)
    if (!month) {
        month = "???"
    }

    let weekday = weekDays.get(date.getDay() + 1)
    if (!weekday) {
        weekday = "???"
    }

    const day = date.getDate() < 10 ? "0" + date.getDate() : date.getDate().toString()
    const year = date.getFullYear().toString()

    return `${month}, ${weekday} ${day} ${year} ${date.toTimeString().split(" ")[0]}`
}

const formatStatus = (value: number): string => {
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

export default CalculationsScreen