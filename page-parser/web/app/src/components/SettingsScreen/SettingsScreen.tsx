import React, {useEffect} from "react";

import './SettingsScreen.css'
import {OperatorsAPI} from "../../api/OperatorsAPI";

interface SettingsScreenProps {
    operatorsAPI: OperatorsAPI
    onError: (msg: string) => void
}

const SettingsScreen: React.FC<SettingsScreenProps> = ({operatorsAPI, onError}) => {
    const [plusInputValue, setPlusInputValue] = React.useState<number>(0)
    const [minusInputValue, setMinusInputValue] = React.useState<number>(0)
    const [divideInputValue, setDivideInputValue] = React.useState<number>(0)
    const [multiplyInputValue, setMultiplyInputValue] = React.useState<number>(0)

    useEffect(() => {
        operatorsAPI.getOperators()
            .then((operators) => {
                operators.sort((a, b) => {
                    return a.operationType - b.operationType
                })

                setPlusInputValue(operators[0].durationMS)
                setMinusInputValue(operators[1].durationMS)
                setDivideInputValue(operators[2].durationMS)
                setMultiplyInputValue(operators[3].durationMS)
            })
            .catch((e) => {
                onError(e.toString())
            })
    }, []);

    const saveOperators = () => {
        operatorsAPI.setOperators(plusInputValue, minusInputValue, divideInputValue, multiplyInputValue)
            .catch((e) => {
                onError(e.toString())
            })
    }
    
    const handleSubmit = (evt: React.FormEvent<HTMLFormElement>) => {
        evt.preventDefault()
        saveOperators()
    }

    const onInputChange = (callback: React.Dispatch<React.SetStateAction<number>>) => {
        return (evt: React.ChangeEvent<HTMLInputElement>) => {
            callback(+evt.target.value)
        }
    }

    return (
        <section className="SettingsScreen">
            <h2 className="SettingsScreen__title">
                Настройка времени выполнения
            </h2>

            <form
                onSubmit={handleSubmit}
                className="SettingsScreen__form"
            >
                <div className="SettingsScreen__input-wrapper">
                    <label htmlFor="plus-time">Время сложения (мс):</label>
                    <input
                        type="number"
                        id="plus-time"
                        value={plusInputValue}
                        onChange={onInputChange(setPlusInputValue)}
                    />
                </div>

                <div className="SettingsScreen__input-wrapper">
                    <label htmlFor="minus-time">Время вычитания (мс):</label>
                    <input
                        type="number"
                        id="minus-time"
                        value={minusInputValue}
                        onChange={onInputChange(setMinusInputValue)}
                    />
                </div>

                <div className="SettingsScreen__input-wrapper">
                    <label htmlFor="divide-time">Время деления (мс):</label>
                    <input
                        type="number"
                        id="divide-time"
                        value={divideInputValue}
                        onChange={onInputChange(setDivideInputValue)}
                    />
                </div>

                <div className="SettingsScreen__input-wrapper">
                    <label htmlFor="multiply-time">Время умножения (мс):</label>
                    <input
                        type="number"
                        id="multiply-time"
                        value={multiplyInputValue}
                        onChange={onInputChange(setMultiplyInputValue)}
                    />
                </div>

                <button type="submit">Сохранить изменения</button>
            </form>
        </section>
    )
}

export default SettingsScreen