import React, {useEffect} from "react";

import './WorkersScreen.css'
import Worker from "./Worker";
import {WorkerAPI} from "../../api/WorkersAPI";

interface WorkerDTO {
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

interface Props {
    workersAPI: WorkerAPI
    onError: (msg: string) => void
}

const WorkersScreen: React.FC<Props> = ({workersAPI, onError}) => {
    const [workers, setWorkers] = React.useState<Array<WorkerDTO>>(new Array<WorkerDTO>());

    const getWorkers = () => {
        workersAPI.getWorkers()
            .then((w) => {
                setWorkers(w);
            })
            .catch((e) => {
                onError(e.toString())
            })
    }

    const getTaskCallback = (workerId: number) => {
        return async ()  => {
            try {
                return await workersAPI.getTasks(workerId)
            }
            catch (e: any) {
                onError(e.toString())
            }

            return new Array<Task>()
        }
    }

    useEffect(() => {
        getWorkers()
    }, []);

    return (
        <section className="WorkersScreen">
            <h2 className="WorkersScreen__title">
                Активные машины
            </h2>

            <div className="WorkersScreen__workers-container">
                {workers.map((w)=>{
                    return (
                        <Worker
                            key={w.id}
                            id={w.id}
                            url={w.url}
                            executors={w.executors}
                            lastModified={w.lastModified}

                            getTasks={getTaskCallback(w.id)}
                        />
                    )
                })}
            </div>
        </section>
    )
}

export default WorkersScreen