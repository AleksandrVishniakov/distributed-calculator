import React, {useEffect, useState} from "react";
import ResponseError from "../../ResponseError";

interface WorkersScreenProps {
    host: string
    onError: (msg: string) => void
}

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

const formatResponseError = (err: ResponseError): string => {
    return err.code + ". " + err.message
}

const WorkersScreen: React.FC<WorkersScreenProps> = (props) => {
    const [workers, setWorkers] = useState(new Array<Worker>())
    const [tasks, setTasks] = useState(new Array<Task>())
    const [workerId, setWorkerId] = useState(-1)

    const getWorkers = async () => {
        const response = await fetch(props.host + "/api/workers", {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });

        const json = await response.json();

        if (response.ok) {
            const resp = json as Worker[]

            setWorkers(resp)
        } else {
            const err = json as ResponseError
            props.onError(formatResponseError(err))

            console.error(json)
        }
    }

    const getWorkerTasks = async (workerId: number) => {
        const response = await fetch(props.host + "/api/worker/" + workerId + "/tasks", {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json'
            },
        });

        const json = await response.json();

        if (response.ok) {
            const resp = json as Task[]

            setTasks(resp)
        } else {
            const err = json as ResponseError
            props.onError(formatResponseError(err))
            console.error(json)
        }
    }

    useEffect(() => {
        getWorkers()
    }, []);

    const handleTaskListClick = async (workerId: number) => {
        setWorkerId(workerId)
        getWorkerTasks(workerId)
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

    return (
        <div>
            {
                workers.length > 0 ?
                    workers.map((worker, index) => {
                        return (
                            <div key={index}>
                                <h4>Машина №{worker.id}</h4>
                                <p>Url: {worker.url}</p>
                                <p>Горутины: {worker.executors}</p>

                                {
                                    workerId === worker.id ?
                                        <ol>
                                            {
                                                tasks.length > 0 ? tasks.map((task, index) => {
                                                    return (
                                                        <li>{formatTask(task)}</li>
                                                    )
                                                }) : <i>нет заданий</i>
                                            }
                                        </ol> : <span></span>
                                }

                                <button onClick={() => {
                                    handleTaskListClick(worker.id)
                                }}>Список заданий
                                </button>
                            </div>
                        )
                    }) :
                    <i>нет доступных машин</i>
            }
        </div>
    )
}

export default WorkersScreen