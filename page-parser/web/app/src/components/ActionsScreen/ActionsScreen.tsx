import React from "react";
import ActionsScreenTopBar from "./ActionsScreenTopBar/ActionsScreenTopBar";
import ActionsScreenNavigation from "./ActionsScreenNavigation/ActionsScreenNavigation";
import SettingsScreen from "../SettingsScreen/SettingsScreen";
import CalculationsScreen from "../CalculationsScreen/CalculationsScreen";
import WorkersScreen from "../WorkersScreen/WorkersScreen";
import {ExpressionsAPI} from "../../api/ExpressionsAPI";
import {OperatorsAPI} from "../../api/OperatorsAPI";
import {WorkerAPI} from "../../api/WorkersAPI";

interface ActionsScreenProps {
    expressionsAPI: ExpressionsAPI
    operatorsAPI: OperatorsAPI
    workersAPI : WorkerAPI

    user: {
        login: string
    }

    onLogout: () => void;
    onError: (msg: string) => void;
}

enum Actions {
    Calculations,
    Settings,
    Machines
}

const ActionsScreen: React.FC<ActionsScreenProps> = ({expressionsAPI, operatorsAPI, workersAPI, user, onLogout, onError}) => {
    const [currentAction, setAction] = React.useState<Actions>(Actions.Settings);

    const handleLogout = () => {
        onLogout()
    }

    const renderActions = (): React.ReactNode => {
        switch (currentAction) {
            case Actions.Calculations:
                return (
                    <CalculationsScreen
                        expressionsAPI={expressionsAPI}
                        onError={onError}
                    />
                )
            case Actions.Settings:
                return (
                    <SettingsScreen
                        operatorsAPI={operatorsAPI}
                        onError={onError}
                    />
                )
            case Actions.Machines:
                return (
                    <WorkersScreen
                        workersAPI={workersAPI}
                        onError={onError}
                    />
                )
            default:
                return null
        }
    }

    return (
        <main className="ActionsScreen">
            <ActionsScreenTopBar
                user={user}
                onLogout={handleLogout}
            />

            <ActionsScreenNavigation
                defaultItem="калькулятор"
                Items={[
                    {
                        title: "калькулятор",
                        Icon:
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
                                <path
                                    d="M320-240h60v-80h80v-60h-80v-80h-60v80h-80v60h80v80Zm200-30h200v-60H520v60Zm0-100h200v-60H520v60Zm44-152 56-56 56 56 42-42-56-58 56-56-42-42-56 56-56-56-42 42 56 56-56 58 42 42Zm-314-70h200v-60H250v60Zm-50 472q-33 0-56.5-23.5T120-200v-560q0-33 23.5-56.5T200-840h560q33 0 56.5 23.5T840-760v560q0 33-23.5 56.5T760-120H200Zm0-80h560v-560H200v560Zm0-560v560-560Z"/>
                            </svg>,
                        onClick: () => {
                            setAction(Actions.Calculations)
                        }
                    },

                    {
                        title: "настройки",
                        Icon:
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
                                <path
                                    d="m370-80-16-128q-13-5-24.5-12T307-235l-119 50L78-375l103-78q-1-7-1-13.5v-27q0-6.5 1-13.5L78-585l110-190 119 50q11-8 23-15t24-12l16-128h220l16 128q13 5 24.5 12t22.5 15l119-50 110 190-103 78q1 7 1 13.5v27q0 6.5-2 13.5l103 78-110 190-118-50q-11 8-23 15t-24 12L590-80H370Zm70-80h79l14-106q31-8 57.5-23.5T639-327l99 41 39-68-86-65q5-14 7-29.5t2-31.5q0-16-2-31.5t-7-29.5l86-65-39-68-99 42q-22-23-48.5-38.5T533-694l-13-106h-79l-14 106q-31 8-57.5 23.5T321-633l-99-41-39 68 86 64q-5 15-7 30t-2 32q0 16 2 31t7 30l-86 65 39 68 99-42q22 23 48.5 38.5T427-266l13 106Zm42-180q58 0 99-41t41-99q0-58-41-99t-99-41q-59 0-99.5 41T342-480q0 58 40.5 99t99.5 41Zm-2-140Z"/>
                            </svg>,
                        onClick: () => {
                            setAction(Actions.Settings)
                        }
                    },

                    {
                        title: "активные машины",
                        Icon:
                            <svg xmlns="http://www.w3.org/2000/svg" viewBox="0 -960 960 960">
                                <path
                                    d="M40-240q9-107 65.5-197T256-580l-74-128q-6-9-3-19t13-15q8-5 18-2t16 12l74 128q86-36 180-36t180 36l74-128q6-9 16-12t18 2q10 5 13 15t-3 19l-74 128q94 53 150.5 143T920-240H40Zm240-110q21 0 35.5-14.5T330-400q0-21-14.5-35.5T280-450q-21 0-35.5 14.5T230-400q0 21 14.5 35.5T280-350Zm400 0q21 0 35.5-14.5T730-400q0-21-14.5-35.5T680-450q-21 0-35.5 14.5T630-400q0 21 14.5 35.5T680-350Z"/>
                            </svg>,
                        onClick: () => {
                            setAction(Actions.Machines)
                        }
                    },
                ]}
            />

            {renderActions()}
        </main>
    )
}

export default ActionsScreen;