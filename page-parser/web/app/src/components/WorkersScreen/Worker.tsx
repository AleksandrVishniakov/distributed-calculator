import React from "react";

import './WorkersScreen.css'

interface Task {
    leftResult: number
    operationType: number
    rightResult: number
    status: number
}

interface Props {
    id: number
    url: string
    executors: number
    lastModified: Date

    getTasks: () => Promise<Array<Task>>
}

const Worker: React.FC<Props> = ({id, url, lastModified, executors, getTasks}) => {
    const [calculationsOpen, setCalculationsOpen] = React.useState(false);
    const [tasks, setTasks] = React.useState<Task[] | null>(null);

    const handleInfoClick = () => {
        if (!calculationsOpen) {
            getTasks().then((t) => setTasks(t));

            setCalculationsOpen(true);
        } else {
            setCalculationsOpen(false)
        }
    }

    return (
        <div className="Worker">
            <div className="Worker__info">
                <h3>{`Машина #${id}`}</h3>
                <h4>Характеристики: </h4>
                <ul>
                    <li>Ссылка: {url}</li>
                    <li>Количество горутин: {executors}</li>
                    <li>Последнее обновление: {new Date(lastModified).toTimeString().split(" ")[0]}</li>
                </ul>
                <button
                    onClick={handleInfoClick}
                >
                    {calculationsOpen ? "Закрыть" : "Подробнее"}
                </button>
            </div>

            {calculationsOpen && tasks ?
                <div className="Worker__calculations">
                    {tasks.map((task, index) => {
                        return (
                            <p key={index}>
                                {formatTask(task)}
                            </p>
                        )
                    })}
                </div>:null
            }
        </div>
    )
}

const formatTask = (task: Task): string => {
    let operator: string
    switch (task.operationType) {
        case 0:
            operator = "+"
            break
        case 1:
            operator = "-"
            break
        case 2:
            operator = "*"
            break
        case 3:
            operator = "/"
            break
        default:
            operator = ""
    }

    let status: string
    switch (task.status) {
        case 1:
            status = "в очереди"
            break
        case 2:
            status = "выполняется"
            break
        case 4:
            status = "ошибка вычисления"
            break
        default:
            status = ""
    }

    return task.leftResult + operator + task.rightResult + " (" + status + ")"
}

export default Worker;